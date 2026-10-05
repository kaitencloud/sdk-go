package sdk

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"
	"time"

	"github.com/kaitencloud/sdk-go/internal/gen"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Problem codes a usage report can answer with, which ReportEntitlementUsage maps
// to sentinels.
const (
	codeThresholdExceeded     = "ReportEntitlementUsageMetric.ThresholdExceeded"
	codeTransactionIDReused   = "ReportEntitlementUsageMetric.TransactionIdReused"
	codeInvalidTransactionID  = "ReportEntitlementUsageMetric.InvalidTransactionId"
	transactionIDBodyLocation = "body.transactionId"
)

// maxUsageReportAttempts is how many times a report with a TransactionID is sent
// before its failure is returned: the first attempt and two retries.
const maxUsageReportAttempts = 3

// usageReportRetryDelay is the wait before the first retry; each later one waits
// twice as long. A variable so tests can shorten it.
var usageReportRetryDelay = 200 * time.Millisecond

// transactionIDPattern is the format Kaiten accepts for a TransactionID.
var transactionIDPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)

// UsageReportResult is what an accepted usage report answers.
type UsageReportResult struct {
	// Usage is the entitlement's usage after the report.
	Usage EntitlementUsage
	// Replayed is true when the report carried a TransactionID Kaiten had already
	// accepted: Usage is that earlier report's answer and nothing was counted
	// again. A replay can describe a usage window that has since closed, so do
	// not read it as the current gauge.
	Replayed bool
	// MetadataDropped is true when the report's Metadata was above 4 KiB and was
	// not stored. The report itself was counted.
	MetadataDropped bool
}

// ReportEntitlementUsage reports usage of the entitlement identified by
// entitlementSlug for the instance identified by instanceSlug, and returns the
// resulting usage with what Kaiten said about the report.
//
// Without a TransactionID the report is sent once and never retried: a report
// that timed out may or may not have been counted, and sending it again could
// count it twice. With one, Kaiten applies the report at most once, so it is
// retried on transport errors and on 500, 502, 503 and 504 -- three attempts
// in all, every one under the same key. A retry of a report the server did
// count answers with Replayed set.
//
// A conflict wraps ErrThresholdExceeded (the limit refused the report) or
// ErrTransactionIDReused (the key was already used for a different report) as
// well as the *Error. A malformed TransactionID is refused before any request,
// with ErrInvalidTransactionID.
//
// A Kaiten server older than TransactionID support refuses the field. The
// client then sends the report again without it, and stops sending keys and
// retrying reports for its lifetime; it logs that once.
func (s *Instances) ReportEntitlementUsage(ctx context.Context, instanceSlug, entitlementSlug string, input UsageReportInput) (UsageReportResult, error) {
	var result UsageReportResult

	usageValue, err := NumberUsageValue(input.Value)
	if err != nil {
		return result, fmt.Errorf("encode usage value: %w", err)
	}
	behavior := input.Behavior
	if behavior == "" {
		behavior = Append
	}

	key := input.TransactionID
	if key != "" && !transactionIDPattern.MatchString(key) {
		return result, fmt.Errorf("%w: %q is not 1 to 128 characters of letters, digits, '.', '_', ':' and '-'", ErrInvalidTransactionID, key)
	}
	if s.client.transactionIDsUnsupported.Load() {
		key = ""
	}
	attempts := 1
	if key != "" {
		attempts = maxUsageReportAttempts
	}

	for attempt := 1; ; attempt++ {
		body, err := jsonBody(usageReportPayload{
			Behavior:      &behavior,
			Metadata:      cloneMapPtr(input.Metadata),
			TransactionID: key,
			Value:         usageValue,
		})
		if err != nil {
			return result, fmt.Errorf("encode usage report request: %w", err)
		}

		resp, err := s.client.raw.ReportEntitlementUsageMetricWithBodyWithResponse(ctx, instanceSlug, entitlementSlug, contentTypeJSON, body)
		if err != nil {
			if attempt < attempts && ctx.Err() == nil {
				if waitErr := waitBeforeRetry(ctx, attempt); waitErr != nil {
					return result, fmt.Errorf("report entitlement usage: %w", err)
				}
				continue
			}
			return result, fmt.Errorf("report entitlement usage: %w", err)
		}

		apiResp := newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503)
		if resp.JSON200 != nil {
			result.Usage = *resp.JSON200
			if headers := resp.Headers200; headers != nil {
				result.Replayed = stringValue(headers.IdempotentReplayed) == "true"
				result.MetadataDropped = stringValue(headers.KaitenMetadataDropped) != ""
			}
			return result, nil
		}

		if key != "" && refusesTransactionID(apiResp) {
			s.client.disableTransactionIDs()
			key, attempts = "", 1
			continue
		}
		if attempt < attempts && retryableStatus(apiResp.StatusCode) {
			if waitErr := waitBeforeRetry(ctx, attempt); waitErr != nil {
				return result, responseError(apiResp)
			}
			continue
		}

		return result, usageReportError(apiResp)
	}
}

// usageReportPayload is the report body. TransactionID is omitted when empty:
// the body is closed, and a server without the field refuses any request that
// names it.
type usageReportPayload struct {
	Behavior      *EntitlementUsageBehavior `json:"behavior,omitempty"`
	Metadata      *map[string]any           `json:"metadata,omitempty"`
	TransactionID string                    `json:"transactionId,omitempty"`
	Value         EntitlementUsageValue     `json:"value"`
}

// usageReportError is the *Error of a refused report, with the sentinel its
// problem code stands for.
func usageReportError(resp apiResponse) error {
	err := responseError(resp)
	switch {
	case resp.StatusCode == http.StatusConflict && resp.Code == codeTransactionIDReused:
		return fmt.Errorf("report usage failed: %w: %w", ErrTransactionIDReused, err)
	case resp.StatusCode == http.StatusConflict && resp.Code == codeThresholdExceeded:
		return fmt.Errorf("report usage failed: %w: %w", ErrThresholdExceeded, err)
	default:
		return err
	}
}

// refusesTransactionID reports whether resp is a server without TransactionID
// support refusing the field: a validation problem located on it, from a
// server that does not know the code newer ones answer a malformed key with.
// The key was checked before sending, so it is not malformed.
func refusesTransactionID(resp apiResponse) bool {
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusBadRequest {
		return false
	}
	if resp.Code == codeInvalidTransactionID || resp.Problem == nil || resp.Problem.Errors == nil {
		return false
	}
	for _, detail := range *resp.Problem.Errors {
		if stringValue(detail.Location) == transactionIDBodyLocation {
			return true
		}
	}
	return false
}

func retryableStatus(status int) bool {
	switch status {
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// waitBeforeRetry waits before retry number attempt: usageReportRetryDelay,
// doubled for every retry before it. It returns early with the context's error.
func waitBeforeRetry(ctx context.Context, attempt int) error {
	timer := time.NewTimer(usageReportRetryDelay << (attempt - 1))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// disableTransactionIDs stops this client sending keys, logging the first time.
func (c *Client) disableTransactionIDs() {
	if c.transactionIDsUnsupported.CompareAndSwap(false, true) {
		c.logger().Warn("kaiten: the server does not support usage report transaction IDs; " +
			"reports are sent without them and are not retried for the rest of this client's life. " +
			"Upgrade Kaiten to make retries safe.")
	}
}

// UsageReport is one accepted usage report in an instance's usage history: the
// counter before and after it and the limit it was gated on. Decimals are
// strings, exact to the digit Kaiten stores.
type UsageReport = gen.UsageReport

// UsageReportBehavior is a recorded report's behavior: append or set.
type UsageReportBehavior = gen.UsageReportBehavior

// UsageReportsOptions narrows ListUsageReports.
type UsageReportsOptions struct {
	// From and To bound the reports' acceptance time, From inclusive and To
	// exclusive. To defaults to now and From to 30 days before To. A From earlier
	// than the start of the organization's usage history is refused with a 422
	// whose Code is ListUsageReports.OutsideRetention.
	From, To *time.Time
	// TransactionID selects the report sent with this key.
	TransactionID string
	// Limit caps how many reports are fetched; nil fetches every report in the
	// range.
	Limit *int32
}

// maxUsageReportsPage is the largest page the usage history endpoint accepts.
const maxUsageReportsPage = 500

// ListUsageReports returns the usage history of the entitlement identified by
// entitlementSlug for the instance identified by instanceSlug: every report in
// the range, in the order they were accepted, up to options.Limit. It follows
// the API's pages itself.
func (s *Instances) ListUsageReports(ctx context.Context, instanceSlug, entitlementSlug string, options *UsageReportsOptions) ([]UsageReport, error) {
	params := gen.ListUsageReportsParams{}
	var ceiling int
	if options != nil {
		params.From, params.To = options.From, options.To
		if options.TransactionID != "" {
			params.TransactionId = &options.TransactionID
		}
		if options.Limit != nil {
			ceiling = int(*options.Limit)
		}
	}

	reports := []UsageReport{}
	for {
		size := int32(maxUsageReportsPage)
		if ceiling > 0 && ceiling-len(reports) < maxUsageReportsPage {
			size = int32(ceiling - len(reports)) //nolint:gosec // bounded by maxUsageReportsPage
		}
		params.Limit = &size

		resp, err := s.client.raw.ListUsageReportsWithResponse(ctx, instanceSlug, entitlementSlug, &params)
		if err != nil {
			return nil, fmt.Errorf("list usage reports: %w", err)
		}
		page, err := expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
		if err != nil {
			return nil, err
		}

		reports = append(reports, page.Items...)
		if page.NextAfterSeq == nil || (ceiling > 0 && len(reports) >= ceiling) {
			return reports, nil
		}
		params.AfterSeq = page.NextAfterSeq
	}
}

// UsageExportFormat is an export's encoding.
type UsageExportFormat string

const (
	// UsageExportCSV is RFC 4180 CSV with a header row.
	UsageExportCSV UsageExportFormat = "csv"
	// UsageExportJSON is NDJSON: one UsageReport per line.
	UsageExportJSON UsageExportFormat = "json"
)

// UsageExportOptions narrows ExportUsageReports.
type UsageExportOptions struct {
	// From and To bound the export as in UsageReportsOptions; at most 366 days
	// apart.
	From, To *time.Time
	// Format defaults to UsageExportCSV.
	Format UsageExportFormat
}

// OrganizationUsageExportOptions narrows ExportOrganizationUsageReports.
type OrganizationUsageExportOptions struct {
	// From and To bound the export as in UsageReportsOptions; at most 31 days
	// apart.
	From, To *time.Time
	// Format defaults to UsageExportCSV.
	Format UsageExportFormat
	// InstanceSlug and EntitlementSlug select one live instance or entitlement.
	InstanceSlug, EntitlementSlug string
	// InstanceID and EntitlementID select one instance or entitlement by ID,
	// deleted ones included: their reports are kept after they are gone.
	InstanceID, EntitlementID string
}

// UsageExport is a usage export as it streams in. Read Body to its end and
// close it.
//
// The client's HTTP timeout covers reading Body, and DefaultTimeout is 30
// seconds: give a large export a client of its own through WithHTTPClient.
type UsageExport struct {
	// Filename is the attachment's name, which names the range.
	Filename string
	// ContentType is text/csv or application/x-ndjson.
	ContentType string
	Body        io.ReadCloser
}

// ExportUsageReports streams the usage history of the entitlement identified by
// entitlementSlug for the instance identified by instanceSlug, as CSV or NDJSON.
func (s *Instances) ExportUsageReports(ctx context.Context, instanceSlug, entitlementSlug string, options *UsageExportOptions) (*UsageExport, error) {
	params := gen.ExportUsageReportsParams{}
	if options != nil {
		params.From, params.To = options.From, options.To
		if options.Format != "" {
			format := string(options.Format)
			params.Format = &format
		}
	}

	resp, err := s.client.raw.ExportUsageReports(ctx, instanceSlug, entitlementSlug, &params)
	if err != nil {
		return nil, fmt.Errorf("export usage reports: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, exportUsageReportsError(resp)
	}
	return newUsageExport(resp), nil
}

// ExportOrganizationUsageReports streams the usage history of the whole
// organization, across instances and entitlements, as CSV or NDJSON.
func (s *Instances) ExportOrganizationUsageReports(ctx context.Context, options *OrganizationUsageExportOptions) (*UsageExport, error) {
	params := gen.ExportOrganizationUsageReportsParams{}
	if options != nil {
		params.From, params.To = options.From, options.To
		if options.Format != "" {
			format := string(options.Format)
			params.Format = &format
		}
		if options.InstanceSlug != "" {
			params.InstanceSlug = &options.InstanceSlug
		}
		if options.EntitlementSlug != "" {
			params.EntitlementSlug = &options.EntitlementSlug
		}
		var err error
		if params.InstanceId, err = parseUUIDOption("InstanceID", options.InstanceID); err != nil {
			return nil, err
		}
		if params.EntitlementId, err = parseUUIDOption("EntitlementID", options.EntitlementID); err != nil {
			return nil, err
		}
	}

	resp, err := s.client.raw.ExportOrganizationUsageReports(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("export organization usage reports: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, exportOrganizationUsageReportsError(resp)
	}
	return newUsageExport(resp), nil
}

func newUsageExport(resp *http.Response) *UsageExport {
	export := &UsageExport{ContentType: resp.Header.Get("Content-Type"), Body: resp.Body}
	if _, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition")); err == nil {
		export.Filename = params["filename"]
	}
	return export
}

// exportUsageReportsError reads a refused export's problem, as every other
// wrapper reports it. The exports send their request through the raw client
// so that a successful body streams instead of being read into memory; only a
// refusal is parsed.
func exportUsageReportsError(resp *http.Response) error {
	parsed, err := gen.ParseExportUsageReportsResponse(resp)
	if err != nil {
		return fmt.Errorf("export usage reports: read the error response: %w", err)
	}
	return responseError(newAPIResponse(parsed, parsed.Body, parsed.ApplicationproblemJSON400, parsed.ApplicationproblemJSON401, parsed.ApplicationproblemJSON403, parsed.ApplicationproblemJSON404, parsed.ApplicationproblemJSON422, parsed.ApplicationproblemJSON500))
}

func exportOrganizationUsageReportsError(resp *http.Response) error {
	parsed, err := gen.ParseExportOrganizationUsageReportsResponse(resp)
	if err != nil {
		return fmt.Errorf("export organization usage reports: read the error response: %w", err)
	}
	return responseError(newAPIResponse(parsed, parsed.Body, parsed.ApplicationproblemJSON400, parsed.ApplicationproblemJSON401, parsed.ApplicationproblemJSON403, parsed.ApplicationproblemJSON404, parsed.ApplicationproblemJSON422, parsed.ApplicationproblemJSON500))
}

func parseUUIDOption(name, value string) (*openapi_types.UUID, error) {
	if value == "" {
		return nil, nil
	}
	var id openapi_types.UUID
	if err := id.UnmarshalText([]byte(value)); err != nil {
		return nil, fmt.Errorf("%s %q is not a UUID: %w", name, value, err)
	}
	return &id, nil
}
