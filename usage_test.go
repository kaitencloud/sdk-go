package sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// Retries wait for real; a millisecond keeps the suite fast. Set once, before any
// test runs, so parallel tests never race on it.
func init() {
	usageReportRetryDelay = time.Millisecond
}

const acceptedUsage = `{"entitlementId":"e","entitlementSlug":"seats","licenseId":"l","value":{"type":"number","value":5,"event_count":1},"limit":{"type":"number","value":10}}`

// respondInTurn answers the nth request with the nth responder, and the last
// one for every request after.
func respondInTurn(responders ...responder) responder {
	var mu sync.Mutex
	n := 0
	return func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		i := min(n, len(responders)-1)
		n++
		mu.Unlock()
		return responders[i](req)
	}
}

func respondWithHeaders(status int, body string, headers map[string]string) responder {
	return func(req *http.Request) (*http.Response, error) {
		resp := jsonResponse(req, status, body)
		for name, value := range headers {
			resp.Header.Set(name, value)
		}
		return resp, nil
	}
}

func respondTransportError() responder {
	return func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection reset by peer")
	}
}

func problem(status int, code string) string {
	return fmt.Sprintf(`{"title":%q,"status":%d,"code":%q}`, http.StatusText(status), status, code)
}

func TestReportEntitlementUsageSendsTheTransactionIDOnlyWhenSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input UsageReportInput
		want  string
	}{
		{"with a key", UsageReportInput{Value: 5, TransactionID: "evt-1"}, `{"behavior":"append","transactionId":"evt-1","value":{"type":"number","value":5}}`},
		{"without a key", UsageReportInput{Value: 5}, `{"behavior":"append","value":{"type":"number","value":5}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client, sent := givenClient(t, respondJSON(http.StatusOK, acceptedUsage))

			if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", test.input); err != nil {
				t.Fatalf("ReportEntitlementUsage() error = %v", err)
			}
			if got := string(sent.only(t).body); got != test.want {
				t.Errorf("request body = %s, want %s", got, test.want)
			}
		})
	}
}

func TestReportEntitlementUsageReadsWhatTheServerSaidAboutTheReport(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondWithHeaders(http.StatusOK, acceptedUsage, map[string]string{
		"Idempotent-Replayed":     "true",
		"Kaiten-Metadata-Dropped": "too_large",
	}))
	result, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"})
	if err != nil {
		t.Fatalf("ReportEntitlementUsage() error = %v", err)
	}
	if !result.Replayed || !result.MetadataDropped {
		t.Errorf("result = %+v, want Replayed and MetadataDropped", result)
	}
	if result.Usage.EntitlementSlug != "seats" {
		t.Errorf("Usage = %+v, want the decoded usage", result.Usage)
	}

	plain, _ := givenClient(t, respondJSON(http.StatusOK, acceptedUsage))
	result, err = plain.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5})
	if err != nil {
		t.Fatalf("ReportEntitlementUsage() error = %v", err)
	}
	if result.Replayed || result.MetadataDropped {
		t.Errorf("result = %+v, want neither flag without the headers", result)
	}
}

// A report without a key may have been counted when its request failed, so it is
// never sent twice; a report with one is applied at most once, so it is retried.
func TestReportEntitlementUsageRetriesOnlyWithAKey(t *testing.T) {
	t.Parallel()

	transient := []responder{
		respondJSON(http.StatusServiceUnavailable, problem(http.StatusServiceUnavailable, "ReportEntitlementUsageMetric.LedgerUnavailable")),
		respondTransportError(),
		respondJSON(http.StatusOK, acceptedUsage),
	}

	t.Run("with a key it retries transient failures under the same key", func(t *testing.T) {
		t.Parallel()
		client, sent := givenClient(t, respondInTurn(transient...))

		if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err != nil {
			t.Fatalf("ReportEntitlementUsage() error = %v", err)
		}
		if sent.count() != 3 {
			t.Fatalf("requests = %d, want 3", sent.count())
		}
		for i := range 3 {
			if key := sent.at(t, i).decodeBody(t)["transactionId"]; key != "evt-1" {
				t.Errorf("attempt %d sent transactionId %v, want evt-1", i+1, key)
			}
		}
	})

	t.Run("without a key it sends once", func(t *testing.T) {
		t.Parallel()
		client, sent := givenClient(t, respondInTurn(transient...))

		_, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5})
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("error = %v, want the 503", err)
		}
		if sent.count() != 1 {
			t.Errorf("requests = %d, want 1", sent.count())
		}
	})

	t.Run("it gives up after three attempts", func(t *testing.T) {
		t.Parallel()
		client, sent := givenClient(t, respondJSON(http.StatusBadGateway, `{"title":"Bad Gateway"}`))

		_, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"})
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadGateway {
			t.Fatalf("error = %v, want the 502", err)
		}
		if sent.count() != maxUsageReportAttempts {
			t.Errorf("requests = %d, want %d", sent.count(), maxUsageReportAttempts)
		}
	})

	t.Run("it never retries a refusal", func(t *testing.T) {
		t.Parallel()
		for _, status := range []int{http.StatusConflict, http.StatusUnprocessableEntity, http.StatusNotFound, http.StatusForbidden} {
			client, sent := givenClient(t, respondJSON(status, problem(status, "Some.Code")))
			if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err == nil {
				t.Fatalf("%d: error = nil", status)
			}
			if sent.count() != 1 {
				t.Errorf("%d: requests = %d, want 1", status, sent.count())
			}
		}
	})

	t.Run("a cancelled context stops the retries", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(t.Context())
		client, sent := givenClient(t, func(req *http.Request) (*http.Response, error) {
			cancel()
			return jsonResponse(req, http.StatusServiceUnavailable, `{"title":"Service Unavailable"}`), nil
		})

		if _, err := client.Instances.ReportEntitlementUsage(ctx, "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err == nil {
			t.Fatal("error = nil")
		}
		if sent.count() != 1 {
			t.Errorf("requests = %d, want 1", sent.count())
		}
	})
}

func TestReportEntitlementUsageNamesTheConflict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		code      string
		want      error
		notWanted error
	}{
		{"ReportEntitlementUsageMetric.TransactionIdReused", ErrTransactionIDReused, ErrThresholdExceeded},
		{"ReportEntitlementUsageMetric.ThresholdExceeded", ErrThresholdExceeded, ErrTransactionIDReused},
	}
	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			t.Parallel()
			client, _ := givenClient(t, respondJSON(http.StatusConflict, problem(http.StatusConflict, test.code)))

			_, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"})
			if !errors.Is(err, test.want) || errors.Is(err, test.notWanted) {
				t.Errorf("error = %v, want %v and not %v", err, test.want, test.notWanted)
			}
			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.Code != test.code {
				t.Errorf("errors.As(*Error) lost the problem: %v", err)
			}

			// ReportUsage keeps the code's meaning: a reused key is not a quota.
			err = client.Instances.ReportUsage(t.Context(), "acme-prod", "seats", 5)
			if !errors.Is(err, test.want) || errors.Is(err, test.notWanted) {
				t.Errorf("ReportUsage error = %v, want %v and not %v", err, test.want, test.notWanted)
			}
		})
	}
}

func TestReportEntitlementUsageRefusesAMalformedKeyBeforeSending(t *testing.T) {
	t.Parallel()

	for _, key := range []string{"has space", strings.Repeat("k", 129), "émoji", "a/b"} {
		client, sent := givenClient(t, respondJSON(http.StatusOK, acceptedUsage))
		_, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: key})
		if !errors.Is(err, ErrInvalidTransactionID) {
			t.Errorf("%q: error = %v, want ErrInvalidTransactionID", key, err)
		}
		if sent.count() != 0 {
			t.Errorf("%q: requests = %d, want none", key, sent.count())
		}
	}
}

// A server older than transaction IDs refuses the body member with a generic
// validation problem. The report must still be counted -- sent again without
// the key -- and the client stops sending keys, and so stops retrying, for good.
func TestReportEntitlementUsageFallsBackOnAServerWithoutTransactionIDs(t *testing.T) {
	t.Parallel()

	refused := `{"title":"Unprocessable Entity","status":422,"code":"UNPROCESSABLE","errors":[{"message":"unexpected property","location":"body.transactionId"}]}`
	var logged bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logged, nil))

	client, sent := givenClient(t, func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		if strings.Contains(string(body), "transactionId") {
			return jsonResponse(req, http.StatusUnprocessableEntity, refused), nil
		}
		if strings.Contains(string(body), `"value":7`) {
			return jsonResponse(req, http.StatusServiceUnavailable, `{"title":"Service Unavailable"}`), nil
		}
		return jsonResponse(req, http.StatusOK, acceptedUsage), nil
	}, WithLogger(logger))

	if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err != nil {
		t.Fatalf("first report: error = %v, want it counted without the key", err)
	}
	if sent.count() != 2 || strings.Contains(string(sent.at(t, 1).body), "transactionId") {
		t.Fatalf("requests = %d, want the refused one then one without the key", sent.count())
	}

	// Later reports go out without a key at once, and are not retried.
	if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 7, TransactionID: "evt-2"}); err == nil {
		t.Fatal("second report: error = nil, want the 503")
	}
	if sent.count() != 3 || strings.Contains(string(sent.at(t, 2).body), "transactionId") {
		t.Errorf("requests = %d, want one more, without the key and not retried", sent.count())
	}

	if n := strings.Count(logged.String(), "does not support usage report transaction IDs"); n != 1 {
		t.Errorf("logged %d times, want once:\n%s", n, logged.String())
	}
}

func TestRefusesTransactionIDIgnoresTheNewServersOwnCode(t *testing.T) {
	t.Parallel()

	location := "body.transactionId"
	details := []ErrorDetail{{Location: &location}}
	resp := apiResponse{StatusCode: http.StatusUnprocessableEntity, Problem: &ErrorModel{Errors: &details}}
	if !refusesTransactionID(resp) {
		t.Error("a validation problem on body.transactionId is an old server")
	}
	resp.Code = "ReportEntitlementUsageMetric.InvalidTransactionId"
	if refusesTransactionID(resp) {
		t.Error("a malformed key on a new server is not an old server")
	}
}

func TestListUsageReportsWalksThePages(t *testing.T) {
	t.Parallel()

	page := func(seqs []int, next string) string {
		items := make([]string, len(seqs))
		for i, seq := range seqs {
			items[i] = fmt.Sprintf(`{"reportSeq":%d,"reportedAt":"2026-10-05T08:00:00Z","behavior":"append","aggregationMethod":"SUM","reportedValue":"1","valueBefore":"0","valueAfter":"1","delta":"1","overageDelta":"0","eventCountAfter":1,"instanceId":"11111111-1111-1111-1111-111111111111","entitlementId":"22222222-2222-2222-2222-222222222222","licenseId":"33333333-3333-3333-3333-333333333333"}`, seq)
		}
		return fmt.Sprintf(`{"items":[%s]%s}`, strings.Join(items, ","), next)
	}

	from := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	client, sent := givenClient(t, respondInTurn(
		respondJSON(http.StatusOK, page([]int{1, 2}, `,"nextAfterSeq":2`)),
		respondJSON(http.StatusOK, page([]int{3}, "")),
	))

	reports, err := client.Instances.ListUsageReports(t.Context(), "acme-prod", "seats", &UsageReportsOptions{From: &from, TransactionID: "evt-1"})
	if err != nil {
		t.Fatalf("ListUsageReports() error = %v", err)
	}
	if len(reports) != 3 || reports[2].ReportSeq != 3 {
		t.Fatalf("reports = %+v, want seqs 1, 2, 3", reports)
	}

	first, second := sent.at(t, 0), sent.at(t, 1)
	if first.path != "/api/instances/acme-prod/entitlements/seats/usage/reports" {
		t.Errorf("path = %s", first.path)
	}
	if first.query.Get("afterSeq") != "" || second.query.Get("afterSeq") != "2" {
		t.Errorf("afterSeq = %q then %q, want none then 2", first.query.Get("afterSeq"), second.query.Get("afterSeq"))
	}
	if first.query.Get("limit") != "500" || first.query.Get("transactionId") != "evt-1" || first.query.Get("from") == "" {
		t.Errorf("query = %v", first.query)
	}

	capped, sent := givenClient(t, respondJSON(http.StatusOK, page([]int{1, 2}, `,"nextAfterSeq":2`)))
	limit := int32(2)
	reports, err = capped.Instances.ListUsageReports(t.Context(), "acme-prod", "seats", &UsageReportsOptions{Limit: &limit})
	if err != nil || len(reports) != 2 {
		t.Fatalf("capped: reports = %d, error = %v", len(reports), err)
	}
	if sent.count() != 1 || sent.only(t).query.Get("limit") != "2" {
		t.Errorf("capped: requests = %d, want one page of 2", sent.count())
	}
}

func TestExportUsageReportsStreamsTheAttachment(t *testing.T) {
	t.Parallel()

	csv := "organization_id,instance_id\no,i\n"
	client, sent := givenClient(t, func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header: http.Header{
				"Content-Type":        {"text/csv; charset=utf-8"},
				"Content-Disposition": {`attachment; filename="usage-acme-prod-seats-20261001-20261005.csv"`},
			},
			Body:    io.NopCloser(strings.NewReader(csv)),
			Request: req,
		}, nil
	})

	export, err := client.Instances.ExportUsageReports(t.Context(), "acme-prod", "seats", &UsageExportOptions{Format: UsageExportCSV})
	if err != nil {
		t.Fatalf("ExportUsageReports() error = %v", err)
	}
	defer func() { _ = export.Body.Close() }()
	body, err := io.ReadAll(export.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != csv {
		t.Errorf("body = %q", body)
	}
	if export.Filename != "usage-acme-prod-seats-20261001-20261005.csv" || export.ContentType != "text/csv; charset=utf-8" {
		t.Errorf("export = %+v", export)
	}
	if got := sent.only(t); got.path != "/api/instances/acme-prod/entitlements/seats/usage/reports/export" || got.query.Get("format") != "csv" {
		t.Errorf("request = %s?%s", got.path, got.query.Encode())
	}
}

func TestExportsReportTheirProblems(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusUnprocessableEntity, http.StatusInternalServerError} {
		code := fmt.Sprintf("ExportUsageReports.%d", status)
		client, _ := givenClient(t, respondJSON(status, problem(status, code)))

		_, err := client.Instances.ExportUsageReports(t.Context(), "acme-prod", "seats", nil)
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != status || apiErr.Code != code {
			t.Errorf("pair export, %d: error = %v", status, err)
		}

		_, err = client.Instances.ExportOrganizationUsageReports(t.Context(), nil)
		if !errors.As(err, &apiErr) || apiErr.StatusCode != status || apiErr.Code != code {
			t.Errorf("organization export, %d: error = %v", status, err)
		}
	}
}

func TestExportOrganizationUsageReportsSendsItsFilters(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondWithHeaders(http.StatusOK, `{"reportSeq":1}`+"\n", map[string]string{"Content-Type": "application/x-ndjson"}))
	export, err := client.Instances.ExportOrganizationUsageReports(t.Context(), &OrganizationUsageExportOptions{
		Format:        UsageExportJSON,
		InstanceID:    "11111111-1111-1111-1111-111111111111",
		EntitlementID: "22222222-2222-2222-2222-222222222222",
	})
	if err != nil {
		t.Fatalf("ExportOrganizationUsageReports() error = %v", err)
	}
	_ = export.Body.Close()

	query := sent.only(t).query
	if sent.only(t).path != "/api/usage/reports/export" || query.Get("format") != "json" ||
		query.Get("instanceId") != "11111111-1111-1111-1111-111111111111" ||
		query.Get("entitlementId") != "22222222-2222-2222-2222-222222222222" {
		t.Errorf("request = %s?%s", sent.only(t).path, query.Encode())
	}

	if _, err := client.Instances.ExportOrganizationUsageReports(t.Context(), &OrganizationUsageExportOptions{InstanceID: "not-a-uuid"}); err == nil {
		t.Error("a malformed InstanceID was sent")
	}
}
