// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// vendorCostPattern is the format the API accepts for a commit's vendorCost.
var vendorCostPattern = regexp.MustCompile(`^[0-9]{1,12}(\.[0-9]{1,12})?$`)

// aiIdempotencyKeyPattern is the format the API accepts for a refund's Idempotency-Key.
var aiIdempotencyKeyPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)

// validateVendorCost refuses, before any request, a vendor cost the API would answer
// with a 422: a malformed decimal, or one on a NOT_SENT commit.
func (in AICommitInput) validateVendorCost() error {
	if in.VendorCost == "" {
		return nil
	}

	if in.RequestStatus == AIRequestStatusNotSent {
		return fmt.Errorf("%w: a NOT_SENT commit charges nothing and takes no vendor cost", ErrInvalidVendorCost)
	}

	if !vendorCostPattern.MatchString(in.VendorCost) {
		return fmt.Errorf("%w: %q is not a decimal of at most 12 integer and 12 fractional digits", ErrInvalidVendorCost, in.VendorCost)
	}

	return nil
}

// VendorCostFromRat formats a vendor cost the caller computed, in USD, as the decimal
// string AICommitInput.VendorCost takes: rounded half away from zero to 12 decimals,
// trailing zeros dropped, so the same cost always gives the same string. Compute it
// once per call and keep the string for the commit's retries. It refuses a negative
// cost and one of more than 12 integer digits.
//
// It holds no price: the cost is whatever the caller's own price list, or the
// provider's bill, says.
func VendorCostFromRat(usd *big.Rat) (string, error) {
	if usd == nil || usd.Sign() < 0 {
		return "", fmt.Errorf("%w: a vendor cost is a non-negative amount", ErrInvalidVendorCost)
	}

	formatted := usd.FloatString(12)
	if strings.Contains(formatted, ".") {
		formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	}

	if !vendorCostPattern.MatchString(formatted) {
		return "", fmt.Errorf("%w: %s USD has more than 12 integer digits", ErrInvalidVendorCost, formatted)
	}

	return formatted, nil
}

// AICapabilities is what AI.Capabilities read.
type AICapabilities struct {
	// Supported is false when the server predates AI usage: it does not serve the
	// capabilities operation at all.
	Supported bool

	// Enabled is whether AI usage is turned on for the organization. While it is off,
	// every AI report answers ErrAIUsageDisabled.
	Enabled bool
}

// Capabilities reads whether the server meters AI usage and whether it is turned on for
// the organization. It is the way to detect support: an AI report's refusal is not.
// The answer can change at any time and is not cached; read it again after an
// ErrAIUsageDisabled. It needs no scope beyond a valid token, and is answered while AI
// usage is off.
func (s *AI) Capabilities(ctx context.Context) (AICapabilities, error) {
	resp, err := s.client.raw.GetAiUsageCapabilitiesWithResponse(ctx)
	if err != nil {
		return AICapabilities{}, fmt.Errorf("get AI capabilities: %w", err)
	}

	if resp.JSON200 != nil {
		return AICapabilities{Supported: true, Enabled: resp.JSON200.Enabled}, nil
	}

	if resp.StatusCode() == http.StatusNotFound {
		return AICapabilities{}, nil
	}

	return AICapabilities{}, responseError(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON500))
}

// AIInferencesOptions filters an inference history. Every field is optional. Without a
// range the API reads the last 30 days.
type AIInferencesOptions struct {
	// From is the start of the range, inclusive; To its end, exclusive.
	From, To *time.Time

	// Status keeps the operations in this state. A release that closed no reservation
	// (RELEASED with NO_RESERVATION) is listed only when Status is RELEASED or
	// TransactionID is set.
	Status AIInferenceStatus

	// TransactionID keeps the operation reported under this id.
	TransactionID string

	// Authorized keeps the operations a commit authorized (true) or refused (false). An
	// operation still RESERVED matches neither.
	Authorized *bool

	// ModelSlug keeps the operations that resolved to this catalogue model; Model the
	// ones whose model, as sent, is exactly this one, which finds models the catalogue
	// does not know.
	ModelSlug, Model string

	// RequestType keeps the operations of this request type.
	RequestType AIRequestType

	// Limit caps how many operations are returned, newest first. Nil returns every one
	// in the range, walking as many pages as that takes.
	Limit *int32
}

// OrganizationAIInferencesOptions filters the organization's inference history.
type OrganizationAIInferencesOptions struct {
	AIInferencesOptions

	// InstanceSlug and EntitlementSlug keep one instance's, or one credit pool's,
	// operations. A slug the organization does not have matches nothing; it is not an
	// error.
	InstanceSlug, EntitlementSlug string
}

func (o *AIInferencesOptions) query() (url.Values, int32) {
	query := url.Values{}
	if o == nil {
		return query, 0
	}

	set := func(name, value string) {
		if value != "" {
			query.Set(name, value)
		}
	}
	set("from", stringValue(timeStringPtr(o.From)))
	set("to", stringValue(timeStringPtr(o.To)))
	set("status", string(o.Status))
	set("transactionId", o.TransactionID)
	if o.Authorized != nil {
		query.Set("authorized", strconv.FormatBool(*o.Authorized))
	}
	set("modelSlug", o.ModelSlug)
	set("model", o.Model)
	set("requestType", string(o.RequestType))

	var limit int32
	if o.Limit != nil {
		limit = *o.Limit
	}

	return query, limit
}

// ListInferences returns the AI operations of one credit pool of an instance, newest
// first, with the pricing snapshot, licence, warnings, vendor cost and refunds of
// each. It reads a pool the instance no longer holds, too.
func (s *AI) ListInferences(ctx context.Context, instanceSlug, entitlementSlug string, options *AIInferencesOptions) ([]AIInference, error) {
	query, limit := options.query()

	inferences, err := listAtMost[AIInference](ctx, s.client.list, listPath("instances", instanceSlug, "entitlements", entitlementSlug, "usage", "inferences"), query, limit)
	if err != nil {
		return nil, aiReportError(err)
	}

	return inferences, nil
}

// ListOrganizationInferences returns the AI operations of every instance and credit pool
// of the organization, newest first, each naming its instance and pool. It needs
// read:instances.
func (s *AI) ListOrganizationInferences(ctx context.Context, options *OrganizationAIInferencesOptions) ([]OrganizationAIInference, error) {
	var (
		query url.Values
		limit int32
	)
	if options != nil {
		query, limit = options.query()
		if options.InstanceSlug != "" {
			query.Set("instanceSlug", options.InstanceSlug)
		}
		if options.EntitlementSlug != "" {
			query.Set("entitlementSlug", options.EntitlementSlug)
		}
	}

	inferences, err := listAtMost[OrganizationAIInference](ctx, s.client.list, listPath("usage", "inferences"), query, limit)
	if err != nil {
		return nil, aiReportError(err)
	}

	return inferences, nil
}

// The typed errors of a refund, one per code the API answers it with. Each is returned
// wrapped together with the *Error, as the AI report's are.
var (
	// ErrRefundContextNotFound is a 404: no such instance, or no such credit pool.
	ErrRefundContextNotFound = errors.New("AI refund: instance or credit pool not found")
	// ErrRefundInferenceNotFound is a 404: no inference under that transaction id on
	// the instance and pool.
	ErrRefundInferenceNotFound = errors.New("AI refund: inference not found")
	// ErrRefundNotCommitted is a 409: the inference is RESERVED, RELEASED or EXPIRED,
	// and only a COMMITTED one is charged.
	ErrRefundNotCommitted = errors.New("AI refund: inference not committed")
	// ErrRefundNothingToRefund is a 409: the inference was charged nothing, or its
	// charge is already refunded in full.
	ErrRefundNothingToRefund = errors.New("AI refund: nothing left to refund")
	// ErrRefundExceedsRefundable is a 409: more credits than the inference has left to
	// give back. RefundableCredits reads what it has left.
	ErrRefundExceedsRefundable = errors.New("AI refund: exceeds what is left to refund")
	// ErrRefundIdempotencyKeyReused is a 409: the key already holds another refund, the
	// same refund with another body, or a numeric report of the pool. Send this refund
	// under a new key.
	ErrRefundIdempotencyKeyReused = errors.New("AI refund: idempotency key reused")
	// ErrRefundPoolNotGranted is a 409: the instance no longer holds the credit pool, and
	// a refund is a report against it.
	ErrRefundPoolNotGranted = errors.New("AI refund: credit pool no longer granted")
	// ErrRefundIdempotencyKeyRequired is a 422: the request carried no key. The SDK
	// always sends one, so it reaches a caller only through a proxy that drops it.
	ErrRefundIdempotencyKeyRequired = errors.New("AI refund: idempotency key required")
	// ErrRefundInvalidIdempotencyKey is a key that is not 1 to 128 characters of
	// [A-Za-z0-9._:-]: refused before any request, or a 422.
	ErrRefundInvalidIdempotencyKey = errors.New("AI refund: invalid idempotency key")
	// ErrRefundInvalidAmount is a 422: credits is not a positive decimal of at most 12
	// integer and 6 fractional digits.
	ErrRefundInvalidAmount = errors.New("AI refund: invalid amount")
	// ErrRefundInvalidReason is a 422: reason is empty or longer than 500 characters.
	ErrRefundInvalidReason = errors.New("AI refund: invalid reason")
)

// refundErrors maps a refund's problem codes onto their typed errors.
var refundErrors = map[string]error{
	"RefundAiInference.ContextNotFound":        ErrRefundContextNotFound,
	"RefundAiInference.InferenceNotFound":      ErrRefundInferenceNotFound,
	"RefundAiInference.NotCommitted":           ErrRefundNotCommitted,
	"RefundAiInference.NothingToRefund":        ErrRefundNothingToRefund,
	"RefundAiInference.ExceedsRefundable":      ErrRefundExceedsRefundable,
	"RefundAiInference.IdempotencyKeyReused":   ErrRefundIdempotencyKeyReused,
	"RefundAiInference.PoolNotGranted":         ErrRefundPoolNotGranted,
	"RefundAiInference.IdempotencyKeyRequired": ErrRefundIdempotencyKeyRequired,
	"RefundAiInference.InvalidIdempotencyKey":  ErrRefundInvalidIdempotencyKey,
	"RefundAiInference.InvalidAmount":          ErrRefundInvalidAmount,
	"RefundAiInference.InvalidReason":          ErrRefundInvalidReason,
	codeAIUsageDisabled:                        ErrAIUsageDisabled,
}

// AIRefundInput is the input for AI.RefundInference.
type AIRefundInput struct {
	// Credits is how many credits to give back, a positive decimal string of at most 6
	// fractional digits. Empty refunds everything the inference has left.
	Credits string

	// Reason says why. It is kept with the refund and shown to anyone who can read the
	// organization's instances. Empty gives none.
	Reason string

	// IdempotencyKey makes the refund safe to send again: the same key with the same
	// inference, credits and reason answers the original refund and moves nothing.
	// Empty generates a UUIDv7, which AIRefund.IdempotencyKey reports. Use one key per
	// decision to refund -- not per attempt -- and reuse it to retry. It must match
	// ^[A-Za-z0-9._:-]{1,128}$.
	IdempotencyKey string
}

// AIRefund is what AI.RefundInference reported.
type AIRefund struct {
	// Refund is the refund recorded, or the one first recorded under the key.
	Refund AIInferenceRefund

	// IdempotencyKey is the key the refund was sent under, the generated one included.
	// It is set on an error as well: resend under it to retry without refunding twice.
	IdempotencyKey string

	// Replayed is true when the key already held this refund: nothing moved again.
	Replayed bool
}

// RefundInference gives back credits charged for one COMMITTED inference of a credit
// pool, in full or in part. The refund is a negative usage report: it lands in the usage
// window open now and is listed on the inference (AIInference.Refunds,
// Credits.Refunded); the charge itself does not change. It needs write:instances.
//
// The refund is sent under input.IdempotencyKey, or a generated one. A network error or
// a 500, 502, 503 or 504 is retried up to three times with the same key and body, which
// the API applies once. Refusals carry their typed error: ErrRefundNotCommitted,
// ErrRefundNothingToRefund, ErrRefundExceedsRefundable, ErrRefundIdempotencyKeyReused,
// ErrRefundPoolNotGranted, ErrRefundContextNotFound, ErrRefundInferenceNotFound, the
// ErrRefundInvalid* errors, or ErrAIUsageDisabled.
func (s *AI) RefundInference(ctx context.Context, instanceSlug, entitlementSlug, transactionID string, input AIRefundInput) (AIRefund, error) {
	key := input.IdempotencyKey
	if key == "" {
		id, err := newTransactionID()
		if err != nil {
			return AIRefund{}, err
		}
		key = id
	}
	result := AIRefund{IdempotencyKey: key}

	if !aiIdempotencyKeyPattern.MatchString(key) {
		return result, fmt.Errorf("%w: %q is not 1 to 128 characters of letters, digits, '.', '_', ':' and '-'", ErrRefundInvalidIdempotencyKey, key)
	}

	// Marshalled once: every attempt sends these bytes under this key.
	body, err := json.Marshal(gen.NewAiInferenceRefund{Credits: optionalString(input.Credits), Reason: optionalString(input.Reason)})
	if err != nil {
		return result, fmt.Errorf("encode AI refund request: %w", err)
	}
	params := &gen.RefundAiInferenceParams{IdempotencyKey: &key}

	for attempt := 0; ; attempt++ {
		httpResp, err := s.client.raw.RefundAiInferenceWithBody(ctx, instanceSlug, entitlementSlug, transactionID, params, contentTypeJSON, bytes.NewReader(body))
		resp := &gen.RefundAiInferenceResponse{HTTPResponse: httpResp}

		if attempt < aiMaxRetries && retryableAIReport(ctx, resp, err) {
			if httpResp != nil {
				_ = httpResp.Body.Close()
			}
			if waitErr := s.backoff(ctx, attempt); waitErr != nil {
				return result, fmt.Errorf("refund AI inference: %w", errors.Join(retryCause(resp, err), waitErr))
			}
			continue
		}

		if err != nil {
			return result, fmt.Errorf("refund AI inference: %w", err)
		}

		resp, err = gen.ParseRefundAiInferenceResponse(httpResp)
		if err != nil {
			return result, fmt.Errorf("decode AI refund response: %w", err)
		}

		if resp.JSON200 == nil {
			converted := newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503)
			return result, refundError(responseError(converted))
		}

		result.Refund = *resp.JSON200
		result.Replayed = resp.HTTPResponse.Header.Get(idempotentReplayedHeader) == "true"

		return result, nil
	}
}

// refundError tags a refused refund with its typed error.
func refundError(err error) error {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return err
	}

	sentinel, ok := refundErrors[apiErr.Code]
	if !ok {
		return err
	}

	return fmt.Errorf("%w: %w", sentinel, err)
}

// RefundableCredits reads, from an ErrRefundExceedsRefundable, how many credits the
// inference still has to give back, as a decimal string.
func RefundableCredits(err error) (string, bool) {
	if !errors.Is(err, ErrRefundExceedsRefundable) {
		return "", false
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Problem == nil || apiErr.Problem.Errors == nil {
		return "", false
	}

	for _, detail := range *apiErr.Problem.Errors {
		if stringValue(detail.Location) != "body.credits" {
			continue
		}
		if left, ok := detail.Value.(string); ok {
			return left, true
		}
	}

	return "", false
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
