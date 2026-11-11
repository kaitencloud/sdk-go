// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"
)

// The answers of a contract that explains its charges, and of one that predates it: a
// client must read both.
const (
	// committedWithRuleBody is committedBody from a server that names the rule.
	committedWithRuleBody = `{"entitlementId":"e1","entitlementSlug":"chat-credits","licenseId":"l1","licenseSlug":"pro-v3",
		"value":{"type":"number","value":91.2},"limit":{"type":"number","value":100},
		"reserved":"0","source":"api","transactionId":null,
		"ai":{"transactionId":"tx-1","phase":"commit","status":"COMMITTED","credits":{"reserved":"6.6","computed":"1.2","charged":"1.2"},
			"remaining":"8.8","authorized":true,"rule":{"id":"3a9b8c7d-0000-4000-8000-000000000001","specificity":"EXACT_TYPED"},
			"warnings":["COUNTER_PRECISION_DEGRADED","A_WARNING_FROM_THE_FUTURE"]}}`

	// releasedWithoutReservationBody is a release of a transaction id nothing reserved.
	releasedWithoutReservationBody = `{"entitlementId":"e1","entitlementSlug":"chat-credits","licenseId":"l1","licenseSlug":"pro-v3",
		"value":{"type":"number","value":90},"reserved":"0",
		"ai":{"transactionId":"tx-1","phase":"commit","status":"RELEASED","credits":{"computed":"0","charged":"0"},
			"remaining":"10","warnings":["NO_RESERVATION"]}}`

	// newInferenceRow carries every member the explainable-charges contract adds.
	newInferenceRow = `{"transactionId":"req-42","status":"COMMITTED","reportedAt":"2026-11-11T10:00:00Z","committedAt":"2026-11-11T10:00:02Z",
		"model":"gpt-4o-2024-11-20","modelSlug":"gpt-4o","requestType":"COMPLETION","requestStatus":"SUCCESS",
		"usage":{"inputTokens":1000,"outputTokens":800},"authorized":true,"conversionRuleId":"3a9b8c7d-0000-4000-8000-000000000001",
		"credits":{"reserved":"3.4","computed":"3.4","charged":"3.4","refunded":"1"},"valueBefore":10,"valueAfter":13.4,
		"windowStart":"2026-11-01T00:00:00Z",
		"pricing":{"ruleId":"3a9b8c7d-0000-4000-8000-000000000001","specificity":"EXACT_TYPED","rates":{"perInputToken":"0.001","perOutputToken":"0.003"},
			"minimumCreditsPerRequest":"0","maximumCreditsPerRequest":null,"chargeOnError":"NEVER"},
		"licenseId":"6f1c2a9e-0000-4000-8000-000000000002","licenseSlug":"pro","warnings":["MODEL_DEPRECATED"],"reportSeq":17,
		"vendorCost":{"reported":"0.0125","computed":"0.0122","currency":"USD"},
		"refunds":[{"id":"8d0e0000-0000-4000-8000-000000000003","credits":"1","reason":"duplicate call","reportSeq":18,
			"createdAt":"2026-11-11T11:00:00Z","createdBy":{"id":"u1","name":"Ada Lovelace"}}]}`

	// oldInferenceRow is the same operation read from a server that predates them.
	oldInferenceRow = `{"transactionId":"req-42","status":"COMMITTED","reportedAt":"2026-11-11T10:00:00Z","model":"gpt-4o-2024-11-20",
		"modelSlug":"gpt-4o","requestType":"COMPLETION","usage":{"inputTokens":1000,"outputTokens":800},"authorized":true,
		"conversionRuleId":"3a9b8c7d-0000-4000-8000-000000000001","credits":{"reserved":"3.4","computed":"3.4","charged":"3.4"}}`
)

func TestAIReportReadsTheRuleAndTheNewWarnings(t *testing.T) {
	t.Parallel()

	client, _ := givenAIClient(t, respondJSON(http.StatusOK, committedWithRuleBody))

	report, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1"})
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	rule := report.Outcome.Rule
	if rule == nil || rule.ID != "3a9b8c7d-0000-4000-8000-000000000001" || rule.Specificity != AIRuleSpecificityExactTyped {
		t.Errorf("Rule = %+v, want the EXACT_TYPED rule", rule)
	}
	if !report.Outcome.HasWarning(AIWarningCounterPrecisionDegraded) {
		t.Errorf("Warnings = %v, want COUNTER_PRECISION_DEGRADED", report.Outcome.Warnings)
	}
	if !report.Outcome.HasWarning("A_WARNING_FROM_THE_FUTURE") {
		t.Error("an unknown warning was dropped; the set is open")
	}
	if report.Outcome.NoReservation() {
		t.Error("NoReservation() = true on a commit")
	}
}

// The answer of a server that predates the rule decodes, the rule nil.
func TestAIReportFromAnOlderServerHasNoRule(t *testing.T) {
	t.Parallel()

	client, _ := givenAIClient(t, respondJSON(http.StatusOK, committedBody))

	report, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1"})
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	if report.Outcome.Rule != nil || report.Outcome.Credits.Charged != "1.2" {
		t.Errorf("Outcome = %+v, want the old answer with no rule", report.Outcome)
	}
}

func TestCommitSendsTheVendorCostAsWrittenOnEveryAttempt(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondInTurn(
		respondProblem(http.StatusServiceUnavailable, ""),
		respondNetworkError(),
		respondJSON(http.StatusOK, committedBody),
	))

	_, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{
		TransactionID: "tx-1",
		Usage:         AIUsage{InputTokens: 10},
		VendorCost:    "0.010",
	})
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	want := `{"value":{"type":"ai_usage","phase":"commit","transactionId":"tx-1","requestStatus":"SUCCESS","usage":{"inputTokens":10},"vendorCost":"0.010"}}`
	if sent.count() != 3 {
		t.Fatalf("requests = %d, want 3", sent.count())
	}
	for i := range sent.count() {
		if got := string(sent.at(t, i).body); got != want {
			t.Errorf("attempt %d body = %s, want %s byte for byte", i+1, got, want)
		}
	}
}

func TestCommitRefusesAVendorCostTheAPIWould(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input AICommitInput
	}{
		{"too many decimals", AICommitInput{TransactionID: "tx-1", VendorCost: "0.0000000000001"}},
		{"not a decimal", AICommitInput{TransactionID: "tx-1", VendorCost: "1e-3"}},
		{"negative", AICommitInput{TransactionID: "tx-1", VendorCost: "-1"}},
		{"on a release", AICommitInput{TransactionID: "tx-1", RequestStatus: AIRequestStatusNotSent, VendorCost: "0.01"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, respondJSON(http.StatusOK, committedBody))

			_, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", test.input)
			if !errors.Is(err, ErrInvalidVendorCost) {
				t.Errorf("Commit() error = %v, want ErrInvalidVendorCost", err)
			}
			if sent.count() != 0 {
				t.Errorf("requests = %d, want none", sent.count())
			}
		})
	}
}

func TestVendorCostFromRatIsCanonical(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		usd  *big.Rat
		want string
		err  bool
	}{
		{"tokens times a rate", new(big.Rat).Mul(big.NewRat(1234, 1), big.NewRat(25, 10_000_000)), "0.003085", false},
		{"whole dollars", big.NewRat(3, 1), "3", false},
		{"zero", new(big.Rat), "0", false},
		{"rounded to 12 decimals", big.NewRat(1, 3), "0.333333333333", false},
		{"negative", big.NewRat(-1, 100), "", true},
		{"thirteen integer digits", big.NewRat(1_000_000_000_000, 1), "", true},
		{"nil", nil, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := VendorCostFromRat(test.usd)
			if (err != nil) != test.err || got != test.want {
				t.Errorf("VendorCostFromRat() = %q, %v; want %q, error %v", got, err, test.want, test.err)
			}
			if err != nil && !errors.Is(err, ErrInvalidVendorCost) {
				t.Errorf("error = %v, want ErrInvalidVendorCost", err)
			}
		})
	}
}

func TestRunForwardsTheVendorCostButNotOnARelease(t *testing.T) {
	t.Parallel()

	for _, status := range []AIRequestStatus{"", AIRequestStatusNotSent} {
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, respondToAIReport)

			_, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
				return AICall{Status: status, VendorCost: "0.0123"}, nil
			})
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			cost, present := aiValue(t, sent.at(t, 1))["vendorCost"]
			if wantPresent := status != AIRequestStatusNotSent; present != wantPresent || (present && cost != "0.0123") {
				t.Errorf("commit vendorCost = %v (present %v), want present %v", cost, present, wantPresent)
			}
		})
	}
}

// What Run did before: a panicking fn skipped the commit, so the reservation lapsed and
// the call went unbilled. It now commits PROVIDER_ERROR, then panics again.
func TestRunCommitsAndPanicsAgainWhenFnPanics(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondToAIReport)

	recovered := func() (recovered any) {
		defer func() { recovered = recover() }()

		_, _ = client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
			panic("provider SDK bug")
		})

		return nil
	}()

	if recovered != "provider SDK bug" {
		t.Errorf("recovered = %v, want fn's panic value", recovered)
	}
	if sent.count() != 2 {
		t.Fatalf("requests = %d, want the reserve and a commit", sent.count())
	}

	reserve, commit := aiValue(t, sent.at(t, 0)), aiValue(t, sent.at(t, 1))
	if commit["phase"] != "commit" || commit["requestStatus"] != "PROVIDER_ERROR" || commit["transactionId"] != reserve["transactionId"] {
		t.Errorf("commit = %v, want PROVIDER_ERROR under the reservation's id", commit)
	}
}

// runtime.Goexit -- t.FailNow inside fn, for one -- unwinds without a panic value: the
// commit still goes out.
func TestRunCommitsWhenFnExitsItsGoroutine(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondToAIReport)

	var wg sync.WaitGroup
	wg.Go(func() {
		_, _ = client.AI.Run(context.Background(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
			runtime.Goexit()
			return AICall{}, nil
		})
	})
	wg.Wait()

	if sent.count() != 2 || aiValue(t, sent.at(t, 1))["requestStatus"] != "PROVIDER_ERROR" {
		t.Errorf("requests = %d, want the reserve and a PROVIDER_ERROR commit", sent.count())
	}
}

// A reserve whose outcome is unknown may have been recorded: Run releases its id so the
// hold is freed now, and a reserve still in flight lands on the release and is refused.
func TestRunReleasesAReserveWhoseOutcomeIsUnknown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reserve responder
	}{
		{"gateway errors past the retries", respondProblem(http.StatusBadGateway, "")},
		{"network errors past the retries", respondNetworkError()},
		{"a 200 the client cannot decode", respondJSON(http.StatusOK, `{"ai":`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, func(req *http.Request) (*http.Response, error) {
				if requestPhase(t, req) == "reserve" {
					return test.reserve(req)
				}
				return jsonResponse(req, http.StatusOK, releasedWithoutReservationBody), nil
			})

			run, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{TransactionID: "tx-1", Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
				t.Error("fn was called after a reserve that did not answer")
				return AICall{}, nil
			})
			if err == nil {
				t.Fatal("Run() error = nil")
			}

			last := sent.at(t, sent.count()-1)
			release := aiValue(t, last)
			if release["phase"] != "commit" || release["requestStatus"] != "NOT_SENT" || release["transactionId"] != "tx-1" ||
				release["model"] != "gpt-4o" || release["requestType"] != "COMPLETION" {
				t.Errorf("last request = %v, want a NOT_SENT release of tx-1 naming its model and request type", release)
			}
			if !run.Commit.Outcome.NoReservation() {
				t.Errorf("run.Commit = %+v, want the release's NO_RESERVATION answer", run.Commit)
			}
		})
	}
}

// The API's own refusal is an answer: nothing was reserved, nothing is released.
func TestRunDoesNotReleaseARefusedReserve(t *testing.T) {
	t.Parallel()

	for _, code := range []string{"ReportEntitlementUsageMetric.UnknownModel", "AIUsage.Disabled"} {
		t.Run(code, func(t *testing.T) {
			t.Parallel()

			status := http.StatusUnprocessableEntity
			if code == "AIUsage.Disabled" {
				status = http.StatusForbidden
			}
			client, sent := givenAIClient(t, respondProblem(status, code))

			_, err := client.AI.Run(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion}, func(context.Context) (AICall, error) {
				return AICall{}, nil
			})
			if err == nil || sent.count() != 1 {
				t.Errorf("Run() error = %v after %d requests, want the refusal alone", err, sent.count())
			}
		})
	}
}

func TestReleaseOfNothingReservedIsReadAsATombstone(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondInTurn(
		respondJSON(http.StatusOK, releasedWithoutReservationBody),
		respondProblem(http.StatusConflict, "ReportEntitlementUsageMetric.TransactionIdReused"),
	))

	release, err := client.AI.Commit(t.Context(), "acme-prod", "chat-credits", AICommitInput{TransactionID: "tx-1", RequestStatus: AIRequestStatusNotSent})
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	if !release.Outcome.NoReservation() {
		t.Errorf("NoReservation() = false on %+v", release.Outcome)
	}
	if body := aiValue(t, sent.at(t, 0)); body["model"] != nil || body["requestType"] != nil {
		t.Errorf("release = %v, want no model or request type when none is given", body)
	}

	// The reserve that was in flight lands after the release: refused, never retried.
	_, err = client.AI.Reserve(t.Context(), "acme-prod", "chat-credits", AIReserveInput{TransactionID: "tx-1", Model: "gpt-4o", RequestType: AIRequestTypeCompletion})
	if !errors.Is(err, ErrTransactionIDReused) {
		t.Errorf("late Reserve() error = %v, want ErrTransactionIDReused", err)
	}
	if sent.count() != 2 {
		t.Errorf("requests = %d, want 2", sent.count())
	}
}

func TestAIUsageDisabledIsTyped(t *testing.T) {
	t.Parallel()

	calls := map[string]func(*Client) error{
		"Reserve": func(c *Client) error {
			_, err := c.AI.Reserve(t.Context(), "acme-prod", "chat-credits", AIReserveInput{Model: "gpt-4o", RequestType: AIRequestTypeCompletion})
			return err
		},
		"ListInferences": func(c *Client) error {
			_, err := c.AI.ListInferences(t.Context(), "acme-prod", "chat-credits", nil)
			return err
		},
		"ListOrganizationInferences": func(c *Client) error {
			_, err := c.AI.ListOrganizationInferences(t.Context(), nil)
			return err
		},
		"RefundInference": func(c *Client) error {
			_, err := c.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "tx-1", AIRefundInput{})
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, _ := givenAIClient(t, respondProblem(http.StatusForbidden, "AIUsage.Disabled"))

			err := call(client)
			var apiErr *Error
			if !errors.Is(err, ErrAIUsageDisabled) || !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
				t.Errorf("error = %v, want ErrAIUsageDisabled with its *Error", err)
			}
		})
	}

	// A missing scope is a 403 too, and not this.
	client, _ := givenAIClient(t, respondProblem(http.StatusForbidden, "Auth.MissingScope"))
	if _, err := client.AI.Reserve(t.Context(), "acme-prod", "chat-credits", AIReserveInput{}); errors.Is(err, ErrAIUsageDisabled) {
		t.Errorf("Auth.MissingScope read as ErrAIUsageDisabled: %v", err)
	}
}

func TestCapabilities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		respond responder
		want    AICapabilities
	}{
		{"enabled", respondJSON(http.StatusOK, `{"enabled":true}`), AICapabilities{Supported: true, Enabled: true}},
		{"disabled", respondJSON(http.StatusOK, `{"enabled":false}`), AICapabilities{Supported: true}},
		{"a server that predates AI usage", respondProblem(http.StatusNotFound, ""), AICapabilities{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, test.respond)

			got, err := client.AI.Capabilities(t.Context())
			if err != nil || got != test.want {
				t.Errorf("Capabilities() = %+v, %v; want %+v", got, err, test.want)
			}
			if req := sent.only(t); req.method != http.MethodGet || req.path != "/api/ai/capabilities" {
				t.Errorf("request = %s %s", req.method, req.path)
			}
		})
	}
}

func TestListInferencesSendsEveryFilterAndReadsTheNewMembers(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondJSON(http.StatusOK, `{"items":[`+newInferenceRow+`],"hasMore":false}`))

	from := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)
	authorized := false
	limit := int32(20)
	rows, err := client.AI.ListInferences(t.Context(), "acme-prod", "chat-credits", &AIInferencesOptions{
		From: &from, Status: "COMMITTED", TransactionID: "req-42", Authorized: &authorized,
		ModelSlug: "gpt-4o", Model: "gpt-4o-2024-11-20", RequestType: AIRequestTypeCompletion, Limit: &limit,
	})
	if err != nil {
		t.Fatalf("ListInferences() error = %v", err)
	}

	req := sent.only(t)
	if req.path != "/api/instances/acme-prod/entitlements/chat-credits/usage/inferences" {
		t.Errorf("path = %s", req.path)
	}
	for name, want := range map[string]string{
		"from": "2026-11-01T00:00:00Z", "status": "COMMITTED", "transactionId": "req-42", "authorized": "false",
		"modelSlug": "gpt-4o", "model": "gpt-4o-2024-11-20", "requestType": "COMPLETION", "limit": "20",
	} {
		if got := req.query.Get(name); got != want {
			t.Errorf("query %s = %q, want %q", name, got, want)
		}
	}
	if req.query.Has("to") || req.query.Has("cursor") {
		t.Errorf("query = %v, carries a filter nobody set", req.query)
	}

	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Pricing == nil || row.Pricing.Specificity != AIInferencePricingSpecificity(AIRuleSpecificityExactTyped) ||
		stringValue(row.Pricing.Rates.PerOutputToken) != "0.003" || row.Pricing.MaximumCreditsPerRequest != nil {
		t.Errorf("Pricing = %+v", row.Pricing)
	}
	if stringValue(row.LicenseSlug) != "pro" || row.LicenseId == nil || row.ReportSeq == nil || *row.ReportSeq != 17 {
		t.Errorf("licence and journal link = %v %v %v", row.LicenseSlug, row.LicenseId, row.ReportSeq)
	}
	if row.VendorCost == nil || stringValue(row.VendorCost.Reported) != "0.0125" || stringValue(row.VendorCost.Computed) != "0.0122" {
		t.Errorf("VendorCost = %+v", row.VendorCost)
	}
	if len(row.Refunds) != 1 || stringValue(row.Refunds[0].CreatedBy.Name) != "Ada Lovelace" || stringValue(row.Credits.Refunded) != "1" {
		t.Errorf("Refunds = %+v, Credits = %+v", row.Refunds, row.Credits)
	}
	if len(row.Warnings) != 1 || AIWarning(row.Warnings[0]) != AIWarningModelDeprecated {
		t.Errorf("Warnings = %v", row.Warnings)
	}
}

// The same operation read from a server that predates the new members decodes, each of
// them absent.
func TestListInferencesReadsAnOlderServer(t *testing.T) {
	t.Parallel()

	client, _ := givenAIClient(t, respondJSON(http.StatusOK, `{"items":[`+oldInferenceRow+`],"hasMore":false}`))

	rows, err := client.AI.ListInferences(t.Context(), "acme-prod", "chat-credits", nil)
	if err != nil {
		t.Fatalf("ListInferences() error = %v", err)
	}

	row := rows[0]
	if row.Pricing != nil || row.VendorCost != nil || row.ReportSeq != nil || row.LicenseId != nil || row.Refunds != nil || row.Warnings != nil || row.Credits.Refunded != nil {
		t.Errorf("row = %+v, want every new member absent", row)
	}
	if stringValue(row.Credits.Charged) != "3.4" {
		t.Errorf("Credits = %+v", row.Credits)
	}
}

func TestListOrganizationInferencesWalksPagesWithItsFilters(t *testing.T) {
	t.Parallel()

	rows := make([]map[string]any, 0, 250)
	for range 250 {
		rows = append(rows, map[string]any{
			"transactionId": "tx", "status": "RELEASED", "reportedAt": "2026-11-11T10:00:00Z", "model": "", "requestType": "COMPLETION",
			"usage": map[string]any{}, "credits": map[string]any{"charged": "0"}, "warnings": []string{"NO_RESERVATION"}, "refunds": []any{},
			"instanceSlug": "acme-prod", "entitlementSlug": "chat-credits",
		})
	}
	client, sent := givenAIClient(t, pagedRows(rows))

	got, err := client.AI.ListOrganizationInferences(t.Context(), &OrganizationAIInferencesOptions{
		AIInferencesOptions: AIInferencesOptions{Status: "RELEASED"},
		InstanceSlug:        "acme-prod",
		EntitlementSlug:     "chat-credits",
	})
	if err != nil {
		t.Fatalf("ListOrganizationInferences() error = %v", err)
	}

	if len(got) != 250 || stringValue(got[249].InstanceSlug) != "acme-prod" || stringValue(got[0].EntitlementSlug) != "chat-credits" {
		t.Fatalf("rows = %d, want 250 naming their instance and pool", len(got))
	}
	if sent.count() != 2 {
		t.Fatalf("requests = %d, want two pages", sent.count())
	}
	for i := range sent.count() {
		req := sent.at(t, i)
		if req.path != "/api/usage/inferences" || req.query.Get("instanceSlug") != "acme-prod" || req.query.Get("entitlementSlug") != "chat-credits" || req.query.Get("status") != "RELEASED" {
			t.Errorf("request %d = %s?%s", i, req.path, req.query.Encode())
		}
	}
}

const refundBody = `{"id":"8d0e0000-0000-4000-8000-000000000003","credits":"3.5","reason":"duplicate call","reportSeq":18,
	"createdAt":"2026-11-11T11:00:00Z","createdBy":{"id":"u1","name":"Ada Lovelace"}}`

func TestRefundInferenceSendsItsKeyAndBody(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, func(req *http.Request) (*http.Response, error) {
		resp := jsonResponse(req, http.StatusOK, refundBody)
		resp.Header.Set("Idempotent-Replayed", "true")
		return resp, nil
	})

	refund, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req/42", AIRefundInput{Credits: "3.5", Reason: "duplicate call", IdempotencyKey: "refund-req-42-1"})
	if err != nil {
		t.Fatalf("RefundInference() error = %v", err)
	}

	req := sent.only(t)
	if req.method != http.MethodPost || req.path != "/api/instances/acme-prod/entitlements/chat-credits/usage/inferences/req/42/refund" {
		t.Errorf("request = %s %s", req.method, req.path)
	}
	if got := req.header.Get("Idempotency-Key"); got != "refund-req-42-1" {
		t.Errorf("Idempotency-Key = %q", got)
	}
	if got, want := string(req.body), `{"credits":"3.5","reason":"duplicate call"}`; got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
	if refund.Refund.Credits != "3.5" || refund.Refund.ReportSeq != 18 || refund.IdempotencyKey != "refund-req-42-1" || !refund.Replayed {
		t.Errorf("refund = %+v", refund)
	}
}

// No key generates one, sent and reported; an empty input refunds everything left with
// a body of {}, which the API requires.
func TestRefundInferenceGeneratesAKey(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondJSON(http.StatusOK, refundBody))

	refund, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req-42", AIRefundInput{})
	if err != nil {
		t.Fatalf("RefundInference() error = %v", err)
	}

	req := sent.only(t)
	if key := req.header.Get("Idempotency-Key"); key == "" || key != refund.IdempotencyKey || !aiIdempotencyKeyPattern.MatchString(key) {
		t.Errorf("Idempotency-Key = %q, reported %q", key, refund.IdempotencyKey)
	}
	if string(req.body) != `{}` {
		t.Errorf("body = %s, want {}", req.body)
	}
	if refund.Replayed {
		t.Error("Replayed = true without the header")
	}
}

// Over a real socket: a dropped connection and a 503 are retried under the same key and
// body, and the key -- generated -- is reported when every attempt failed.
func TestRefundInferenceRetriesUnderOneKey(t *testing.T) {
	t.Parallel()

	var (
		mu     sync.Mutex
		keys   []string
		bodies []string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		mu.Lock()
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		bodies = append(bodies, string(body))
		attempt := len(keys)
		mu.Unlock()

		switch attempt {
		case 1:
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Errorf("hijack: %v", err)
				return
			}
			_ = conn.Close()
		case 2:
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"title":"Service Unavailable","status":503}`))
		default:
			w.Header().Set("Content-Type", contentTypeJSON)
			_, _ = w.Write([]byte(refundBody))
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL+"/api", WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.AI.retryBase = time.Microsecond

	refund, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req-42", AIRefundInput{Credits: "3.5"})
	if err != nil {
		t.Fatalf("RefundInference() error = %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(keys) != 3 {
		t.Fatalf("attempts = %d, want 3", len(keys))
	}
	for i := range keys {
		if keys[i] != refund.IdempotencyKey || bodies[i] != `{"credits":"3.5"}` {
			t.Errorf("attempt %d = key %q body %s, want key %q and the same body", i+1, keys[i], bodies[i], refund.IdempotencyKey)
		}
	}
}

func TestRefundInferenceReportsItsKeyWhenItGivesUp(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondProblem(http.StatusServiceUnavailable, ""))

	refund, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req-42", AIRefundInput{})
	if err == nil {
		t.Fatal("RefundInference() error = nil")
	}
	if refund.IdempotencyKey == "" || sent.count() != 1+aiMaxRetries {
		t.Errorf("key %q after %d requests, want the key reported after four", refund.IdempotencyKey, sent.count())
	}
}

func TestRefundInferenceRefusalsCarryTheirTypedError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status int
		code   string
		want   error
	}{
		{http.StatusNotFound, "RefundAiInference.ContextNotFound", ErrRefundContextNotFound},
		{http.StatusNotFound, "RefundAiInference.InferenceNotFound", ErrRefundInferenceNotFound},
		{http.StatusConflict, "RefundAiInference.NotCommitted", ErrRefundNotCommitted},
		{http.StatusConflict, "RefundAiInference.NothingToRefund", ErrRefundNothingToRefund},
		{http.StatusConflict, "RefundAiInference.ExceedsRefundable", ErrRefundExceedsRefundable},
		{http.StatusConflict, "RefundAiInference.IdempotencyKeyReused", ErrRefundIdempotencyKeyReused},
		{http.StatusConflict, "RefundAiInference.PoolNotGranted", ErrRefundPoolNotGranted},
		{http.StatusUnprocessableEntity, "RefundAiInference.IdempotencyKeyRequired", ErrRefundIdempotencyKeyRequired},
		{http.StatusUnprocessableEntity, "RefundAiInference.InvalidIdempotencyKey", ErrRefundInvalidIdempotencyKey},
		{http.StatusUnprocessableEntity, "RefundAiInference.InvalidAmount", ErrRefundInvalidAmount},
		{http.StatusUnprocessableEntity, "RefundAiInference.InvalidReason", ErrRefundInvalidReason},
		{http.StatusForbidden, "AIUsage.Disabled", ErrAIUsageDisabled},
		{http.StatusForbidden, "Auth.MissingScope", nil},
		{http.StatusConflict, "RefundAiInference.SomethingNew", nil},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			t.Parallel()

			client, sent := givenAIClient(t, respondProblem(test.status, test.code))

			_, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req-42", AIRefundInput{IdempotencyKey: "k-1"})

			for code, sentinel := range refundErrors {
				if got, want := errors.Is(err, sentinel), sentinel == test.want; got != want {
					t.Errorf("errors.Is(err, %s) = %v, want %v: %v", code, got, want, err)
				}
			}
			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.Code != test.code {
				t.Errorf("errors.As(*Error) lost the problem: %v", err)
			}
			if sent.count() != 1 {
				t.Errorf("requests = %d, want 1: a 4xx is never retried", sent.count())
			}
		})
	}
}

func TestRefundInferenceRefusesAMalformedKeyBeforeSending(t *testing.T) {
	t.Parallel()

	client, sent := givenAIClient(t, respondJSON(http.StatusOK, refundBody))

	_, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req-42", AIRefundInput{IdempotencyKey: "not a key"})
	if !errors.Is(err, ErrRefundInvalidIdempotencyKey) || sent.count() != 0 {
		t.Errorf("error = %v after %d requests, want ErrRefundInvalidIdempotencyKey before any", err, sent.count())
	}
}

func TestRefundableCreditsReadsWhatIsLeft(t *testing.T) {
	t.Parallel()

	client, _ := givenAIClient(t, respondJSON(http.StatusConflict, `{"title":"Conflict","status":409,"code":"RefundAiInference.ExceedsRefundable",
		"errors":[{"message":"credits left to refund on this inference","location":"body.credits","value":"1.5"}]}`))

	_, err := client.AI.RefundInference(t.Context(), "acme-prod", "chat-credits", "req-42", AIRefundInput{Credits: "9"})

	left, ok := RefundableCredits(err)
	if !ok || left != "1.5" {
		t.Errorf("RefundableCredits() = %q, %v; want 1.5", left, ok)
	}
	if _, ok := RefundableCredits(errors.New("other")); ok {
		t.Error("RefundableCredits() read an unrelated error")
	}
}
