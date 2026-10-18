// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// reservedBody and committedBody are the API's answers to a reserve and its commit, as
// the contract's example gives them.
const (
	reservedBody = `{"entitlementId":"e1","entitlementSlug":"chat-credits","licenseId":"l1","licenseSlug":"pro-v3",
		"value":{"type":"number","value":90},"limit":{"type":"number","value":100},
		"reserved":"6.6",
		"ai":{"transactionId":"tx-1","phase":"reserve","status":"RESERVED","credits":{"reserved":"6.6"},
			"remaining":"3.4","expiresAt":"2026-10-18T10:10:00Z","warnings":[]}}`
	committedBody = `{"entitlementId":"e1","entitlementSlug":"chat-credits","licenseId":"l1","licenseSlug":"pro-v3",
		"value":{"type":"number","value":91.2},"limit":{"type":"number","value":100},
		"reserved":"0",
		"ai":{"transactionId":"tx-1","phase":"commit","status":"COMMITTED","credits":{"reserved":"6.6","charged":"1.2"},
			"remaining":"8.8","authorized":true,"warnings":["MODEL_DEPRECATED"]}}`
)

// respondInTurn answers the nth request with the nth responder, and every later one
// with the last.
func respondInTurn(responders ...responder) responder {
	var mu sync.Mutex
	next := 0

	return func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		respond := responders[min(next, len(responders)-1)]
		next++
		mu.Unlock()

		return respond(req)
	}
}

// respondProblem answers with an RFC 9457 document carrying code.
func respondProblem(status int, code string) responder {
	return respondJSON(status, fmt.Sprintf(`{"title":%q,"status":%d,"code":%q}`, http.StatusText(status), status, code))
}

// respondNetworkError fails the request before any response, as a reset connection does.
func respondNetworkError() responder {
	return func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection reset by peer")
	}
}

// givenAIClient is givenClient with the keyed retry's backoff shortened to nothing.
func givenAIClient(t *testing.T, respond responder) (*Client, *recorder) {
	t.Helper()

	client, sent := givenClient(t, respond)
	client.AI.retryBase = time.Microsecond

	return client, sent
}

// aiValue decodes the value member of a recorded AI report.
func aiValue(t *testing.T, req recordedRequest) map[string]any {
	t.Helper()

	value, ok := req.decodeBody(t)["value"].(map[string]any)
	if !ok {
		t.Fatalf("request body %s has no value object", req.body)
	}

	return value
}

func TestReserveSendsTheReservePhaseAndReadsItsOutcome(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, func(req *http.Request) (*http.Response, error) {
		resp := jsonResponse(req, http.StatusOK, reservedBody)
		resp.Header.Set("Idempotent-Replayed", "true")
		return resp, nil
	})

	report, err := client.AI.Reserve(t.Context(), "acme-prod", "chat-credits", AIReserveInput{
		TransactionID: "tx-1",
		Model:         "gpt-4o",
		RequestType:   AIRequestTypeCompletion,
		Estimate:      AIUsage{InputTokens: 1000, OutputTokens: 4000},
		TTL:           90500 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Reserve() error = %v", err)
	}

	req := sent.only(t)
	if req.method != http.MethodPost || req.path != "/api/instances/acme-prod/entitlements/chat-credits/usage" {
		t.Errorf("request = %s %s, want the usage report route", req.method, req.path)
	}
	if req.header.Get("Idempotency-Key") == "" {
		t.Error("request carries no Idempotency-Key")
	}

	want := `{"value":{"type":"ai_usage","phase":"reserve","transactionId":"tx-1","model":"gpt-4o","requestType":"COMPLETION","usage":{"inputTokens":1000,"outputTokens":4000},"ttlSeconds":91}}`
	if got := string(req.body); got != want {
		t.Errorf("request body = %s, want %s", got, want)
	}

	if report.Reserved != "6.6" || report.Outcome.Status != AIOperationStatusReserved || report.Outcome.Credits.Reserved != "6.6" {
		t.Errorf("report = %+v, want the reservation of 6.6", report)
	}
	if report.Outcome.Remaining == nil || *report.Outcome.Remaining != "3.4" {
		t.Errorf("Remaining = %v, want 3.4", report.Outcome.Remaining)
	}
	if report.Outcome.ExpiresAt == nil || !report.Outcome.ExpiresAt.Equal(time.Date(2026, time.October, 18, 10, 10, 0, 0, time.UTC)) {
		t.Errorf("ExpiresAt = %v", report.Outcome.ExpiresAt)
	}
	if report.Usage.LicenseSlug != "pro-v3" {
		t.Errorf("Usage.LicenseSlug = %q, want the embedded EntitlementUsage decoded", report.Usage.LicenseSlug)
	}
	if !report.Replayed {
		t.Error("Replayed = false on an Idempotent-Replayed answer")
	}
}

// A reserve with no transaction id gets a UUIDv7: it has to pair with the commit, so
// the default must be unique, and v7 keeps a pool's operations in creation order.
func TestReserveDefaultsTheTransactionIDToAUUIDv7(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondJSON(http.StatusOK, reservedBody))

	if _, err := client.AI.Reserve(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}); err != nil {
		t.Fatalf("Reserve() error = %v", err)
	}

	id, ok := aiValue(t, sent.only(t))["transactionId"].(string)
	if !ok {
		t.Fatal("request carries no transactionId")
	}

	parsed, err := uuid.Parse(id)
	if err != nil || parsed.Version() != 7 {
		t.Errorf("transactionId = %q, want a UUIDv7", id)
	}
}

// A commit sends only what the reservation does not already know: model and request
// type are inherited, and an unset status is a success.
func TestCommitInheritsFromTheReservation(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondJSON(http.StatusOK, committedBody))

	report, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{
		TransactionID: "tx-1",
		Usage:         AIUsage{InputTokens: 1000, OutputTokens: 400},
	})
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	want := `{"value":{"type":"ai_usage","phase":"commit","transactionId":"tx-1","requestStatus":"SUCCESS","usage":{"inputTokens":1000,"outputTokens":400}}}`
	if got := string(sent.only(t).body); got != want {
		t.Errorf("request body = %s, want %s", got, want)
	}

	if report.Outcome.Credits.Charged != "1.2" || report.Outcome.Authorized == nil || !*report.Outcome.Authorized {
		t.Errorf("Outcome = %+v, want 1.2 charged and authorized", report.Outcome)
	}
	if len(report.Outcome.Warnings) != 1 || report.Outcome.Warnings[0] != AIWarningModelDeprecated {
		t.Errorf("Warnings = %v, want [MODEL_DEPRECATED]", report.Outcome.Warnings)
	}
}

func TestAIReportRefusalsCarryTheirTypedError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status int
		code   string
		want   error // nil: no sentinel, the *Error alone
	}{
		{"estimate does not fit", http.StatusConflict, "ReportEntitlementUsageMetric.ThresholdExceeded", ErrInsufficientCredits},
		{"transaction id reused", http.StatusConflict, "ReportEntitlementUsageMetric.TransactionIdReused", ErrTransactionIDReused},
		{"model not authorized", http.StatusUnprocessableEntity, "ReportEntitlementUsageMetric.ModelNotAuthorized", ErrModelNotAuthorized},
		{"unknown model", http.StatusUnprocessableEntity, "ReportEntitlementUsageMetric.UnknownModel", ErrUnknownModel},
		{"not a credit pool", http.StatusUnprocessableEntity, "ReportEntitlementUsageMetric.AiUsageUnsupportedEntitlement", ErrAIUsageUnsupported},
		{"server predating ai_usage, generic validation", http.StatusUnprocessableEntity, "", ErrAIUsageUnsupported},
		{"server predating ai_usage, value type", http.StatusUnprocessableEntity, "ReportEntitlementUsageMetric.InvalidValueType", ErrAIUsageUnsupported},
		{"malformed AI body", http.StatusUnprocessableEntity, "ReportEntitlementUsageMetric.InvalidAiUsage", nil},
		{"malformed transaction id", http.StatusUnprocessableEntity, "ReportEntitlementUsageMetric.InvalidTransactionId", nil},
		{"pool unknown", http.StatusNotFound, "ReportEntitlementUsageMetric.ContextNotFound", nil},
	}

	sentinels := []error{ErrInsufficientCredits, ErrTransactionIDReused, ErrModelNotAuthorized, ErrUnknownModel, ErrAIUsageUnsupported}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, respondProblem(test.status, test.code))

			_, err := client.AI.Reserve(t.Context(), "acme-prod", "chat-credits", AIReserveInput{TransactionID: "tx-1", Model: "gpt-4o", RequestType: AIRequestTypeCompletion})
			if err == nil {
				t.Fatal("Reserve() error = nil")
			}

			for _, sentinel := range sentinels {
				if got, want := errors.Is(err, sentinel), sentinel == test.want; got != want {
					t.Errorf("errors.Is(err, %v) = %v, want %v: %v", sentinel, got, want, err)
				}
			}

			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.StatusCode != test.status || apiErr.Code != test.code {
				t.Errorf("errors.As(*Error) lost the problem: %v", err)
			}

			// A 4xx is the API's answer, not a transient failure: never retried.
			if sent.count() != 1 {
				t.Errorf("requests = %d, want 1", sent.count())
			}
		})
	}
}

// The retry doctrine: a network error or a 500, 502, 503 or 504 is retried at most three
// times, every attempt carrying the same Idempotency-Key and the same body -- so the
// same transaction id -- since only that lets the API tell a retry from a second report.
func TestAIReportsRetryUnderOneKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		respond  responder
		requests int
		ok       bool
	}{
		{
			name:     "gateway statuses then success",
			respond:  respondInTurn(respondProblem(http.StatusBadGateway, ""), respondProblem(http.StatusServiceUnavailable, ""), respondProblem(http.StatusGatewayTimeout, ""), respondJSON(http.StatusOK, committedBody)),
			requests: 4,
			ok:       true,
		},
		{
			name:     "network error then success",
			respond:  respondInTurn(respondNetworkError(), respondJSON(http.StatusOK, committedBody)),
			requests: 2,
			ok:       true,
		},
		{
			name:     "500 every time gives up after three retries",
			respond:  respondProblem(http.StatusInternalServerError, ""),
			requests: 4,
		},
		{
			name:     "network error every time gives up after three retries",
			respond:  respondNetworkError(),
			requests: 4,
		},
		{
			name:     "a 400 is not retried",
			respond:  respondProblem(http.StatusBadRequest, ""),
			requests: 1,
		},
		{
			name:     "a 501 is not retried",
			respond:  respondProblem(http.StatusNotImplemented, ""),
			requests: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, test.respond)

			_, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1"})
			if (err == nil) != test.ok {
				t.Fatalf("Commit() error = %v, want ok = %v", err, test.ok)
			}

			if sent.count() != test.requests {
				t.Fatalf("requests = %d, want %d", sent.count(), test.requests)
			}

			first := sent.at(t, 0)
			for i := 1; i < sent.count(); i++ {
				retry := sent.at(t, i)
				if retry.header.Get("Idempotency-Key") != first.header.Get("Idempotency-Key") {
					t.Errorf("attempt %d Idempotency-Key = %q, want %q", i, retry.header.Get("Idempotency-Key"), first.header.Get("Idempotency-Key"))
				}
				if string(retry.body) != string(first.body) {
					t.Errorf("attempt %d body = %s, want %s", i, retry.body, first.body)
				}
			}
		})
	}
}

// Two logical calls are two keys: a key reused across calls would have the API answer
// the second commit with the first one's replay.
func TestEveryAIReportGetsItsOwnKey(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondJSON(http.StatusOK, committedBody))

	for range 2 {
		if _, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1"}); err != nil {
			t.Fatalf("Commit() error = %v", err)
		}
	}

	if sent.at(t, 0).header.Get("Idempotency-Key") == sent.at(t, 1).header.Get("Idempotency-Key") {
		t.Error("two calls sent the same Idempotency-Key")
	}
}

// The backoff waits on the caller's context: cancelling it ends the retries at once
// rather than after the remaining delays.
func TestAIRetryStopsWhenTheContextIsDone(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondProblem(http.StatusServiceUnavailable, ""))
	client.AI.retryBase = time.Hour

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := client.AI.Commit(ctx, "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1"})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Commit() error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("Commit() returned after %v, want it to stop with the context", elapsed)
	}
	if sent.count() != 1 {
		t.Errorf("requests = %d, want 1", sent.count())
	}
}

// Over a real socket: the first connection is dropped before any answer, which is the
// network error a retry exists for, and the second attempt is answered.
func TestAIReportRetriesADroppedConnection(t *testing.T) {
	t.Parallel()

	var (
		mu   sync.Mutex
		keys []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)

		mu.Lock()
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		attempt := len(keys)
		mu.Unlock()

		if attempt == 1 {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Errorf("hijack: %v", err)
				return
			}
			_ = conn.Close()
			return
		}

		w.Header().Set("Content-Type", contentTypeJSON)
		_, _ = w.Write([]byte(committedBody))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL+"/api", WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.AI.retryBase = time.Microsecond

	if _, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1"}); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(keys) != 2 || keys[0] == "" || keys[0] != keys[1] {
		t.Errorf("Idempotency-Keys = %q, want two attempts under one key", keys)
	}
}

// respondToAIReport answers a reserve and a commit as the API does, echoing the
// transaction id they carry.
func respondToAIReport(req *http.Request) (*http.Response, error) {
	var body struct {
		Value struct {
			Phase         string `json:"phase"`
			TransactionID string `json:"transactionId"`
		} `json:"value"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, err
	}

	status := "RESERVED"
	if body.Value.Phase == "commit" {
		status = "COMMITTED"
	}

	return jsonResponse(req, http.StatusOK, fmt.Sprintf(
		`{"entitlementId":"e1","entitlementSlug":"chat-credits","licenseId":"l1","licenseSlug":"pro-v3","value":{"type":"number","value":1},"reserved":"0","ai":{"transactionId":%q,"phase":%q,"status":%q,"credits":{},"remaining":null,"warnings":[]}}`,
		body.Value.TransactionID, body.Value.Phase, status)), nil
}

func TestRunReservesCallsAndCommits(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondToAIReport)

	called := false
	run, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{
		Model:       "gpt-4o",
		RequestType: AIRequestTypeCompletion,
		Estimate:    AIUsage{InputTokens: 1000, OutputTokens: 4000},
	}, func(context.Context) (AICall, error) {
		if sent.count() != 1 {
			t.Errorf("fn ran after %d requests, want after the reserve", sent.count())
		}
		called = true
		return AICall{Usage: AIUsage{InputTokens: 1000, OutputTokens: 400}, Model: "gpt-4o-2024-08-06"}, nil
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !called {
		t.Fatal("fn was not called")
	}

	reserve, commit := aiValue(t, sent.at(t, 0)), aiValue(t, sent.at(t, 1))
	if reserve["phase"] != "reserve" || commit["phase"] != "commit" {
		t.Fatalf("phases = %v then %v, want reserve then commit", reserve["phase"], commit["phase"])
	}
	if reserve["transactionId"] != commit["transactionId"] {
		t.Errorf("commit transactionId = %v, want the reservation's %v", commit["transactionId"], reserve["transactionId"])
	}
	if commit["requestStatus"] != "SUCCESS" || commit["model"] != "gpt-4o-2024-08-06" {
		t.Errorf("commit = %v, want SUCCESS on the model the call used", commit)
	}
	if _, ok := commit["requestType"]; ok {
		t.Error("commit carries requestType, which the reservation already knows")
	}
	if run.Reservation.Outcome.Status != AIOperationStatusReserved || run.Commit.Outcome.Status != AIOperationStatusCommitted {
		t.Errorf("run = %+v", run)
	}
}

// A refused reserve is the gate closing: the provider is never called.
func TestRunDoesNotCallTheProviderWhenTheReserveIsRefused(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondProblem(http.StatusConflict, "ReportEntitlementUsageMetric.ThresholdExceeded"))

	_, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
		t.Error("fn was called after a refused reserve")
		return AICall{}, nil
	})

	if !errors.Is(err, ErrInsufficientCredits) {
		t.Errorf("Run() error = %v, want ErrInsufficientCredits", err)
	}
	if sent.count() != 1 {
		t.Errorf("requests = %d, want the reserve alone", sent.count())
	}
}

func TestRunCommitsAFailedCall(t *testing.T) {
	t.Parallel()

	providerErr := errors.New("provider said no")

	tests := []struct {
		name   string
		call   AICall
		status string
	}{
		{"an unset status is a provider error", AICall{Usage: AIUsage{InputTokens: 10}}, "PROVIDER_ERROR"},
		{"a customer error is committed as one", AICall{Status: AIRequestStatusCustomerError, Usage: AIUsage{InputTokens: 10}}, "CUSTOMER_ERROR"},
		{"a provider never reached releases the reservation", AICall{Status: AIRequestStatusNotSent}, "NOT_SENT"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, respondToAIReport)

			_, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
				return test.call, providerErr
			})

			if !errors.Is(err, providerErr) {
				t.Errorf("Run() error = %v, want fn's error", err)
			}
			if errors.Is(err, ErrAICommitFailed) {
				t.Errorf("Run() error = %v, but the commit succeeded", err)
			}

			commit := aiValue(t, sent.at(t, 1))
			if commit["requestStatus"] != test.status {
				t.Errorf("requestStatus = %v, want %s", commit["requestStatus"], test.status)
			}
			if _, ok := commit["usage"]; ok != (test.call.Usage != AIUsage{}) {
				t.Errorf("commit usage present = %v, want it only when the call reported some", ok)
			}
		})
	}
}

// A call cut short by its own deadline still happened: the commit goes out on a context
// the caller's cancellation does not reach.
func TestRunCommitsAfterTheCallersContextIsDone(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondToAIReport)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	_, err := client.AI.Run(ctx, "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(ctx context.Context) (AICall, error) {
		cancel()
		return AICall{Status: AIRequestStatusCustomerError}, ctx.Err()
	})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Run() error = %v, want fn's context.Canceled", err)
	}
	if sent.count() != 2 {
		t.Fatalf("requests = %d, want the reserve and the commit", sent.count())
	}
}

// A commit that fails after its retries is reported apart from the call: the provider
// answered, and the caller has that answer.
func TestRunSaysWhenOnlyTheCommitFailed(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, func(req *http.Request) (*http.Response, error) {
		if requestPhase(t, req) == "commit" {
			return jsonResponse(req, http.StatusServiceUnavailable, `{"title":"Service Unavailable","status":503}`), nil
		}
		return respondToAIReport(req)
	})

	answer := ""
	_, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
		answer = "hello"
		return AICall{}, nil
	})

	if !errors.Is(err, ErrAICommitFailed) {
		t.Errorf("Run() error = %v, want ErrAICommitFailed", err)
	}
	if answer != "hello" {
		t.Error("fn did not run")
	}
	if sent.count() != 1+1+aiMaxRetries {
		t.Errorf("requests = %d, want the reserve and four commit attempts", sent.count())
	}
}

// requestPhase reads the phase of an AI report and puts the body back for the next reader.
func requestPhase(t *testing.T, req *http.Request) string {
	t.Helper()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("read request body: %v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body))

	var decoded struct {
		Value struct {
			Phase string `json:"phase"`
		} `json:"value"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("decode request body %s: %v", body, err)
	}

	return decoded.Value.Phase
}
