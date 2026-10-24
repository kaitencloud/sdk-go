// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kaitencloud/sdk-go/internal/gen"
)

// The typed errors of an AI report. Each is returned wrapped together with the *Error
// the API answered, so errors.Is finds the sentinel and errors.As still reaches the
// problem details and their code.
var (
	// ErrInsufficientCredits is a reserve refused because its estimate does not fit
	// what remains in the pool. Nothing was set aside; do not call the provider.
	ErrInsufficientCredits = errors.New("insufficient AI credits")

	// ErrModelNotAuthorized is a reserve of a model that is archived, or that no
	// conversion rule of the instance's license prices.
	ErrModelNotAuthorized = errors.New("AI model not authorized")

	// ErrUnknownModel is a reserve of a model the catalogue does not know, by slug or
	// by provider id.
	ErrUnknownModel = errors.New("unknown AI model")

	// ErrTransactionIDReused is a transaction id already used for a different report:
	// the same phase with another body, or a reserve after its commit. It is a bug in
	// the caller and never retried.
	ErrTransactionIDReused = errors.New("AI transaction id reused")

	// ErrAIUsageUnsupported is an AI report the API cannot take: the pool is not an AI
	// credit pool, or the server predates AI usage reports.
	ErrAIUsageUnsupported = errors.New("AI usage reports unsupported")

	// ErrAICommitFailed is returned by AI.Run when the provider call succeeded but its
	// commit could not be confirmed, retries included. The call's result is good;
	// settlement may already have succeeded if its response was lost. Reconcile or
	// retry the identical commit body under the same transaction id before assuming
	// nothing was charged. An open reservation lapses on its own.
	ErrAICommitFailed = errors.New("AI usage commit failed")
)

// The report's codes an AI error is told apart by.
const (
	codeThresholdExceeded      = "ReportEntitlementUsageMetric.ThresholdExceeded"
	codeTransactionIDReused    = "ReportEntitlementUsageMetric.TransactionIdReused"
	codeModelNotAuthorized     = "ReportEntitlementUsageMetric.ModelNotAuthorized"
	codeUnknownModel           = "ReportEntitlementUsageMetric.UnknownModel"
	codeAIUsageUnsupportedPool = "ReportEntitlementUsageMetric.AiUsageUnsupportedEntitlement"
	codeTransactionIDRequired  = "ReportEntitlementUsageMetric.TransactionIdRequired"
	codeInvalidTransactionID   = "ReportEntitlementUsageMetric.InvalidTransactionId"
	codeInvalidAIUsage         = "ReportEntitlementUsageMetric.InvalidAiUsage"
	codeAIUsageRequiresAppend  = "ReportEntitlementUsageMetric.AiUsageRequiresAppend"
	codeAICreditPoolAppendOnly = "ReportEntitlementUsageMetric.AiCreditPoolAppendOnly"
	idempotentReplayedHeader   = "Idempotent-Replayed"
	aiMaxRetries               = 3
	aiRetryBase                = 250 * time.Millisecond
	aiCommitTimeout            = DefaultTimeout
)

// Reserve sets credits aside for a provider call before it is made, priced on
// input.Estimate. The pool's counter does not move until the commit.
//
// A refusal is one of the typed errors: ErrInsufficientCredits, ErrModelNotAuthorized,
// ErrUnknownModel, ErrTransactionIDReused or ErrAIUsageUnsupported. Do not call the
// provider after any error. A network error or a 500, 502, 503 or 504 is retried up to
// three times under the same transaction id and body, which the API deduplicates.
// AI reports never send the numeric-report Idempotency-Key header.
func (s *AI) Reserve(ctx context.Context, instanceSlug, entitlementSlug string, input AIReserveInput) (AIReport, error) {
	if input.TransactionID == "" {
		id, err := newTransactionID()
		if err != nil {
			return AIReport{}, err
		}
		input.TransactionID = id
	}

	return s.report(ctx, instanceSlug, entitlementSlug, input.payload())
}

// Commit records what a provider call consumed and charges it, closing the reservation
// input.TransactionID names. A commit is never refused for a commercial reason: an
// unknown model, a missing rule or an exhausted pool is recorded and answered, flagged
// on AIReport.Outcome. It is retried like Reserve; a retry the API already recorded is
// answered with the original, AIReport.Replayed set.
func (s *AI) Commit(ctx context.Context, instanceSlug, entitlementSlug string, input AICommitInput) (AIReport, error) {
	if input.TransactionID == "" {
		id, err := newTransactionID()
		if err != nil {
			return AIReport{}, err
		}
		input.TransactionID = id
	}

	return s.report(ctx, instanceSlug, entitlementSlug, input.payload())
}

// AIRun is what AI.Run reported: the reservation it opened and the commit that closed it.
type AIRun struct {
	Reservation AIReport
	Commit      AIReport
}

// Run brackets one provider call: it reserves as Reserve does with input, calls fn, and
// commits what fn reports in its AICall under the same transaction id.
//
// When the reserve is refused, fn is not called and the typed error is returned. When
// fn returns an error, the commit still goes out -- with the AICall's Status, or
// PROVIDER_ERROR when it sets none -- and fn's error is returned, joined with the
// commit's if that failed too. Set AIRequestStatusNotSent when the provider was never
// reached: the reservation is released and nothing is charged.
//
// The commit runs on a context that ctx cancelling does not stop, bounded by
// DefaultTimeout, so a call cut short by its own deadline is still accounted for. When
// fn succeeded but the commit could not be confirmed, Run returns an error wrapping
// ErrAICommitFailed: the provider's answer, which fn keeps, is good.
func (s *AI) Run(ctx context.Context, instanceSlug, entitlementSlug string, input AIReserveInput, fn func(context.Context) (AICall, error)) (AIRun, error) {
	var run AIRun

	// Resolved here rather than read back from the reservation's answer: the commit
	// must carry the id the reserve sent, whatever the answer echoes.
	if input.TransactionID == "" {
		id, err := newTransactionID()
		if err != nil {
			return run, err
		}
		input.TransactionID = id
	}

	reservation, err := s.Reserve(ctx, instanceSlug, entitlementSlug, input)
	if err != nil {
		return run, err
	}
	run.Reservation = reservation

	call, callErr := fn(ctx)

	status := call.Status
	if status == "" {
		status = AIRequestStatusSuccess
		if callErr != nil {
			status = AIRequestStatusProviderError
		}
	}

	commitCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), aiCommitTimeout)
	defer cancel()

	run.Commit, err = s.Commit(commitCtx, instanceSlug, entitlementSlug, AICommitInput{
		TransactionID: input.TransactionID,
		RequestStatus: status,
		Model:         call.Model,
		Usage:         call.Usage,
	})
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrAICommitFailed, err)
	}

	if callErr != nil {
		return run, errors.Join(callErr, err)
	}

	return run, err
}

// report sends one AI usage report, retrying the identical phase and transaction id.
func (s *AI) report(ctx context.Context, instanceSlug, entitlementSlug string, payload aiReportPayload) (AIReport, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return AIReport{}, fmt.Errorf("encode AI usage report request: %w", err)
	}

	// Run after client-wide editors: an inherited numeric key is invalid for AI,
	// even when empty or inserted under a non-canonical header spelling.
	withoutNumericKey := func(_ context.Context, req *http.Request) error {
		for name := range req.Header {
			if strings.EqualFold(name, "Idempotency-Key") {
				delete(req.Header, name)
			}
		}
		return nil
	}

	for attempt := 0; ; attempt++ {
		httpResp, err := s.client.raw.ReportEntitlementUsageMetricWithBody(ctx, instanceSlug, entitlementSlug, contentTypeJSON, bytes.NewReader(body), withoutNumericKey)
		resp := &gen.ReportEntitlementUsageMetricResponse{HTTPResponse: httpResp}

		// Decide retries from transport errors and status before decoding JSON:
		// a malformed 4xx (or 200) body is not a network failure.
		if attempt < aiMaxRetries && retryableAIReport(ctx, resp, err) {
			if httpResp != nil {
				_ = httpResp.Body.Close()
			}
			if waitErr := s.backoff(ctx, attempt); waitErr != nil {
				return AIReport{}, fmt.Errorf("report AI usage: %w", errors.Join(retryCause(resp, err), waitErr))
			}
			continue
		}

		if err != nil {
			return AIReport{}, fmt.Errorf("report AI usage: %w", err)
		}

		resp, err = gen.ParseReportEntitlementUsageMetricResponse(httpResp)
		if err != nil {
			return AIReport{}, fmt.Errorf("decode AI usage report response: %w", err)
		}

		converted := newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500)
		if resp.StatusCode() != http.StatusOK {
			return AIReport{}, aiReportError(responseError(converted))
		}

		return decodeAIReport(resp.Body, resp.HTTPResponse)
	}
}

// retryableAIReport is the retry doctrine: a network error or a gateway-class status,
// never a 4xx, and never once the caller's context is done.
func retryableAIReport(ctx context.Context, resp interface{ StatusCode() int }, err error) bool {
	if ctx.Err() != nil {
		return false
	}

	if err != nil {
		return true
	}

	switch resp.StatusCode() {
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}

	return false
}

// retryCause is what the last attempt failed with, for an error that gives up waiting.
func retryCause(resp interface{ Status() string }, err error) error {
	if err != nil {
		return err
	}

	return fmt.Errorf("last attempt answered %s", resp.Status())
}

// backoff waits before retry attempt+1: the base doubled per attempt, half of it fixed
// and half random, so concurrent callers do not retry in step.
func (s *AI) backoff(ctx context.Context, attempt int) error {
	base := s.retryBase
	if base <= 0 {
		base = aiRetryBase
	}

	delay := base << attempt
	delay = delay/2 + rand.N(delay/2+1)

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// aiReportError tags a refused AI report with its typed error.
func aiReportError(err error) error {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return err
	}

	var sentinel error

	switch apiErr.StatusCode {
	case http.StatusConflict:
		switch apiErr.Code {
		case codeTransactionIDReused:
			sentinel = ErrTransactionIDReused
		case codeThresholdExceeded:
			sentinel = ErrInsufficientCredits
		}
	case http.StatusUnprocessableEntity:
		switch apiErr.Code {
		case codeModelNotAuthorized:
			sentinel = ErrModelNotAuthorized
		case codeUnknownModel:
			sentinel = ErrUnknownModel
		case codeAIUsageUnsupportedPool:
			sentinel = ErrAIUsageUnsupported
		case codeTransactionIDRequired, codeInvalidTransactionID, codeInvalidAIUsage, codeAIUsageRequiresAppend, codeAICreditPoolAppendOnly:
			// A malformed body: the caller's bug, reported as the API worded it.
		case "", "ReportEntitlementUsageMetric.InvalidValueType":
			// Known legacy validation responses from servers predating ai_usage.
			sentinel = ErrAIUsageUnsupported
		default:
			// Unknown codes retain their problem details, not a guessed capability error.
		}
	}

	if sentinel == nil {
		return err
	}

	return fmt.Errorf("%w: %w", sentinel, err)
}

// aiReportResponse is EntitlementUsage with the two members an AI report adds.
type aiReportResponse struct {
	EntitlementUsage

	Reserved string    `json:"reserved"`
	AI       AIOutcome `json:"ai"`
}

func decodeAIReport(body []byte, resp *http.Response) (AIReport, error) {
	var decoded aiReportResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return AIReport{}, fmt.Errorf("decode AI usage report response: %w", err)
	}

	return AIReport{
		Usage:    decoded.EntitlementUsage,
		Reserved: decoded.Reserved,
		Outcome:  decoded.AI,
		Replayed: resp != nil && resp.Header.Get(idempotentReplayedHeader) == "true",
	}, nil
}

// newTransactionID is the default transaction id: a UUIDv7, so ids sort by creation.
func newTransactionID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("generate AI transaction id: %w", err)
	}

	return id.String(), nil
}
