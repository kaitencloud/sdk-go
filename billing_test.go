package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func invoiceRows(count int) []map[string]any {
	rows := make([]map[string]any, 0, count)
	for i := range count {
		rows = append(rows, map[string]any{"id": fmt.Sprintf("inv-%d", i)})
	}

	return rows
}

// List promises every invoice: it walks the cursor, at the largest page the
// API takes, and keeps the caller's filters on every page.
func TestBillingInvoicesListWalksEveryPage(t *testing.T) {
	t.Parallel()

	client, rec := givenClient(t, pagedRows(invoiceRows(450)))
	overdue := true
	cursor := "ignored"

	invoices, err := client.Billing.Invoices.List(t.Context(), &ListInvoicesParams{Overdue: &overdue, Cursor: &cursor})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(invoices) != 450 {
		t.Fatalf("List() = %d invoices, want 450", len(invoices))
	}
	if rec.count() != 3 {
		t.Fatalf("requests = %d, want 3 pages of 200", rec.count())
	}
	first, last := rec.at(t, 0), rec.at(t, 2)
	if first.path != "/api/invoices" || first.query.Get("limit") != "200" || first.query.Get("cursor") != "" {
		t.Fatalf("first page = %s?%s", first.path, first.query.Encode())
	}
	if last.query.Get("cursor") != "row-400" || last.query.Get("overdue") != "true" {
		t.Fatalf("last page query = %s", last.query.Encode())
	}
}

// A server handing back the cursor it was given would loop the walk forever.
func TestBillingInvoicesListRefusesARepeatedCursor(t *testing.T) {
	t.Parallel()

	client, rec := givenClient(t, respondJSON(http.StatusOK, `{"items":[{"id":"a"}],"hasMore":true,"nextCursor":"same"}`))

	if _, err := client.Billing.Invoices.List(t.Context(), nil); err == nil {
		t.Fatal("List() error = nil, want the repeated cursor reported")
	}
	if rec.count() != 2 {
		t.Fatalf("requests = %d, want 2", rec.count())
	}
}

// Retry-After reaches the caller on the *Error.
func TestBillingErrorsCarryRetryAfter(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, func(req *http.Request) (*http.Response, error) {
		resp := jsonResponse(req, http.StatusServiceUnavailable, `{"title":"Service Unavailable","status":503,"code":"BillingProvider.Unavailable"}`)
		resp.Header.Set("Retry-After", "30")
		return resp, nil
	})

	_, err := client.Billing.Health(t.Context())

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want an *Error", err)
	}
	if apiErr.Code != "BillingProvider.Unavailable" || apiErr.RetryAfter != 30*time.Second {
		t.Fatalf("Code = %q, RetryAfter = %s", apiErr.Code, apiErr.RetryAfter)
	}
}

func TestBillingHandoffClaimAndAck(t *testing.T) {
	t.Parallel()

	client, rec := givenClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/api/billing/handoff/claim" {
			return jsonResponse(req, http.StatusOK, `{"leaseId":"lease-1","leasedUntil":"2027-01-01T00:15:00Z","invoices":[{"id":"inv-1"}]}`), nil
		}
		return jsonResponse(req, http.StatusOK, `{"id":"inv-1"}`), nil
	})

	limit, lease := int32(10), int32(600)
	claim, err := client.Billing.Handoff.Claim(t.Context(), HandoffClaimSize{Limit: &limit, LeaseSeconds: &lease})
	if err != nil {
		t.Fatalf("Claim() error = %v", err)
	}
	if claim.LeaseId != "lease-1" || len(claim.Invoices) != 1 {
		t.Fatalf("Claim() = %+v", claim)
	}
	reference := "ACC-2027-001"
	if _, err := client.Billing.Handoff.Ack(t.Context(), "inv-1", HandoffBooking{LeaseId: &claim.LeaseId, ExternalReference: &reference}); err != nil {
		t.Fatalf("Ack() error = %v", err)
	}

	claimed, acked := rec.at(t, 0), rec.at(t, 1)
	if claimed.method != http.MethodPost || fmt.Sprint(claimed.decodeBody(t)) != "map[leaseSeconds:600 limit:10]" {
		t.Fatalf("claim = %s %s", claimed.method, claimed.body)
	}
	if acked.method != http.MethodPost || acked.path != "/api/billing/handoff/inv-1/ack" {
		t.Fatalf("ack = %s %s", acked.method, acked.path)
	}
	if body := acked.decodeBody(t); body["leaseId"] != "lease-1" || body["externalReference"] != reference {
		t.Fatalf("ack body = %s", acked.body)
	}
}

// Export answers CSV: the document comes back as the API sent it.
func TestBillingInvoicesExportReturnsTheDocument(t *testing.T) {
	t.Parallel()

	const csv = "invoiceId,total\ninv-1,1200\n"
	client, rec := givenClient(t, func(req *http.Request) (*http.Response, error) {
		resp := jsonResponse(req, http.StatusOK, csv)
		resp.Header.Set("Content-Type", "text/csv")
		return resp, nil
	})

	from := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)
	document, err := client.Billing.Invoices.Export(t.Context(), &ExportInvoicesParams{IssuedFrom: &from, IssuedTo: &to})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}
	if string(document) != csv {
		t.Fatalf("Export() = %q", document)
	}
	if got := rec.only(t).query.Get("issuedFrom"); got != "2027-01-01T00:00:00Z" {
		t.Fatalf("issuedFrom = %q", got)
	}
}

// Recompose answers 200 when it rewrote the draft in place and 201 when it
// replaced it: both are the invoice.
func TestBillingInvoicesRecomposeReadsBothStatuses(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusOK, http.StatusCreated} {
		client, _ := givenClient(t, respondJSON(status, `{"id":"inv-2"}`))

		invoice, err := client.Billing.Invoices.Recompose(t.Context(), "inv-1")
		if err != nil {
			t.Fatalf("%d: Recompose() error = %v", status, err)
		}
		if invoice.Id != "inv-2" {
			t.Fatalf("%d: Recompose() = %+v", status, invoice)
		}
	}
}
