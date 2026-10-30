// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
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

const acceptedUsage = `{"entitlementId":"e","entitlementSlug":"seats","licenseId":"l","value":{"type":"number","value":5,"event_count":1},"limit":{"type":"number","value":10},"transactionId":null}`

// keyedUsage is a Kaiten that honours keys answering a keyed report.
const keyedUsage = `{"entitlementId":"e","entitlementSlug":"seats","licenseId":"l","value":{"type":"number","value":5,"event_count":1},"limit":{"type":"number","value":10},"transactionId":"evt-0"}`

// confirmedClient is a client whose server has already answered a keyed report
// with its transactionId, so its keyed reports are retried.
func confirmedClient(t *testing.T, then responder) (*Client, *recorder) {
	t.Helper()
	client, sent := givenClient(t, respondInTurn(respondJSON(http.StatusOK, keyedUsage), then))
	result, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 1, TransactionID: "evt-0"})
	if err != nil || result.TransactionID != "evt-0" {
		t.Fatalf("confirming report: %+v, %v", result, err)
	}
	sent.reset()
	return client, sent
}

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

func TestReportEntitlementUsageSendsTheKeyInTheHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  UsageReportInput
		header string
	}{
		{"with a key", UsageReportInput{Value: 5, TransactionID: "evt-1"}, "evt-1"},
		{"without a key", UsageReportInput{Value: 5}, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client, sent := givenClient(t, respondJSON(http.StatusOK, acceptedUsage))

			if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", test.input); err != nil {
				t.Fatalf("ReportEntitlementUsage() error = %v", err)
			}
			req := sent.only(t)
			if got := req.header.Get("Idempotency-Key"); got != test.header {
				t.Errorf("Idempotency-Key = %q, want %q", got, test.header)
			}
			// The body is closed: the key is never in it.
			if want := `{"behavior":"append","value":{"type":"number","value":5}}`; string(req.body) != want {
				t.Errorf("request body = %s, want %s", req.body, want)
			}
		})
	}
}

func TestReportEntitlementUsageReadsWhatTheServerSaidAboutTheReport(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondWithHeaders(http.StatusOK, keyedUsage, map[string]string{
		"Idempotent-Replayed":     "true",
		"Kaiten-Metadata-Dropped": "too_large",
	}))
	result, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-0"})
	if err != nil {
		t.Fatalf("ReportEntitlementUsage() error = %v", err)
	}
	if !result.Replayed || !result.MetadataDropped || result.TransactionID != "evt-0" {
		t.Errorf("result = %+v, want Replayed, MetadataDropped and the key", result)
	}
	if result.Usage.EntitlementSlug != "seats" {
		t.Errorf("Usage = %+v, want the decoded usage", result.Usage)
	}

	plain, _ := givenClient(t, respondJSON(http.StatusOK, acceptedUsage))
	result, err = plain.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5})
	if err != nil {
		t.Fatalf("ReportEntitlementUsage() error = %v", err)
	}
	if result.Replayed || result.MetadataDropped || result.TransactionID != "" {
		t.Errorf("result = %+v, want no flag and no key", result)
	}
}

// The retry rule (billing spec §24 #60): a keyed report whose outcome is
// unknown is retried only once the server has answered a report with its
// transactionId, because an older Kaiten ignores the header and would count a
// resent report twice. A report without a key is never sent twice.
func TestReportEntitlementUsageRetriesOnlyOnceKeysAreHonoured(t *testing.T) {
	t.Parallel()

	transient := []responder{
		respondJSON(http.StatusServiceUnavailable, problem(http.StatusServiceUnavailable, "ReportEntitlementUsageMetric.LedgerUnavailable")),
		respondTransportError(),
		respondJSON(http.StatusOK, keyedUsage),
	}

	t.Run("a confirmed server: transient failures are retried under the same key", func(t *testing.T) {
		t.Parallel()
		client, sent := confirmedClient(t, respondInTurn(transient...))

		if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err != nil {
			t.Fatalf("ReportEntitlementUsage() error = %v", err)
		}
		if sent.count() != 3 {
			t.Fatalf("requests = %d, want 3", sent.count())
		}
		for i := range 3 {
			if key := sent.at(t, i).header.Get("Idempotency-Key"); key != "evt-1" {
				t.Errorf("attempt %d sent Idempotency-Key %q, want evt-1", i+1, key)
			}
		}
	})

	t.Run("an unconfirmed server: a keyed report is sent once", func(t *testing.T) {
		t.Parallel()
		client, sent := givenClient(t, respondInTurn(transient...))

		_, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"})
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("error = %v, want the 503", err)
		}
		if sent.count() != 1 {
			t.Errorf("requests = %d, want 1", sent.count())
		}
	})

	t.Run("an older server that ignores the key never confirms it", func(t *testing.T) {
		t.Parallel()
		// It answers without transactionId, keyed or not.
		client, sent := givenClient(t, respondInTurn(respondJSON(http.StatusOK, `{"entitlementId":"e","entitlementSlug":"seats","licenseId":"l","value":{"type":"number","value":5}}`),
			respondJSON(http.StatusBadGateway, `{"title":"Bad Gateway"}`)))
		if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err != nil {
			t.Fatalf("first report: %v", err)
		}
		if _, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-2"}); err == nil {
			t.Fatal("second report: error = nil, want the 502")
		}
		if sent.count() != 2 {
			t.Errorf("requests = %d, want 2: the 502 is not retried", sent.count())
		}
		if key := sent.at(t, 1).header.Get("Idempotency-Key"); key != "evt-2" {
			t.Errorf("Idempotency-Key = %q, want evt-2: a report keeps its key", key)
		}
	})

	t.Run("without a key it sends once, even on a confirmed server", func(t *testing.T) {
		t.Parallel()
		client, sent := confirmedClient(t, respondInTurn(transient...))

		_, err := client.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5})
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("error = %v, want the 503", err)
		}
		if sent.count() != 1 {
			t.Errorf("requests = %d, want 1", sent.count())
		}
	})

	t.Run("the flag is per client", func(t *testing.T) {
		t.Parallel()
		confirmedClient(t, respondJSON(http.StatusOK, keyedUsage))
		other, sent := givenClient(t, respondJSON(http.StatusBadGateway, `{"title":"Bad Gateway"}`))
		if _, err := other.Instances.ReportEntitlementUsage(t.Context(), "acme-prod", "seats", UsageReportInput{Value: 5, TransactionID: "evt-1"}); err == nil {
			t.Fatal("error = nil")
		}
		if sent.count() != 1 {
			t.Errorf("requests = %d, want 1: another client's confirmation does not count", sent.count())
		}
	})

	t.Run("it gives up after three attempts", func(t *testing.T) {
		t.Parallel()
		client, sent := confirmedClient(t, respondJSON(http.StatusBadGateway, `{"title":"Bad Gateway"}`))

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
			client, sent := confirmedClient(t, respondJSON(status, problem(status, "ReportEntitlementUsageMetric.TransactionIdReused")))
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
		client, sent := confirmedClient(t, func(req *http.Request) (*http.Response, error) {
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
