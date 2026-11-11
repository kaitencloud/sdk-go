// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"slices"
	"time"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// Enum types follow the generated Core contract. Public outcome bridge structs
// preserve the SDK's established field types while EntitlementUsage also exposes
// the complete generated ai/reserved response.

// AIRequestType is the kind of provider call an AI usage report meters.
type AIRequestType = gen.AiUsageValueRequestType

// AI request types.
const (
	AIRequestTypeCompletion      AIRequestType = "COMPLETION"
	AIRequestTypeEmbedding       AIRequestType = "EMBEDDING"
	AIRequestTypeImageGeneration AIRequestType = "IMAGE_GENERATION"
	AIRequestTypeAudio           AIRequestType = "AUDIO"
	AIRequestTypeVideo           AIRequestType = "VIDEO"
)

// AIRequestStatus is how the provider call a commit closes ended. It decides what the
// pool's conversion rule charges: NOT_SENT is never charged, and the two error statuses
// are charged as the rule's chargeOnError says.
type AIRequestStatus = gen.AiUsageValueRequestStatus

// AI request statuses.
const (
	// AIRequestStatusSuccess is a provider call that answered.
	AIRequestStatusSuccess AIRequestStatus = "SUCCESS"
	// AIRequestStatusProviderError is a provider call that failed on the provider's side.
	AIRequestStatusProviderError AIRequestStatus = "PROVIDER_ERROR"
	// AIRequestStatusCustomerError is a provider call refused for the caller's request,
	// or abandoned by the caller.
	AIRequestStatusCustomerError AIRequestStatus = "CUSTOMER_ERROR"
	// AIRequestStatusNotSent is a provider that was never reached. The reservation is
	// released and nothing is charged, whatever the rule.
	AIRequestStatusNotSent AIRequestStatus = "NOT_SENT"
)

// AIUnitType is what AIUsage.Units counts.
type AIUnitType = gen.AiUsageUnitType

// AI unit types.
const (
	AIUnitTypeImage  AIUnitType = "image"
	AIUnitTypeSecond AIUnitType = "second"
	AIUnitTypeFrame  AIUnitType = "frame"
)

// AIPhase is the phase of an AI usage report.
type AIPhase = gen.AiUsageOutcomePhase

// AI usage report phases.
const (
	AIPhaseReserve AIPhase = "reserve"
	AIPhaseCommit  AIPhase = "commit"
)

// AIOperationStatus is the state of an AI operation after the report that answered it.
type AIOperationStatus = gen.AiUsageOutcomeStatus

// AI operation statuses.
const (
	AIOperationStatusReserved  AIOperationStatus = "RESERVED"
	AIOperationStatusCommitted AIOperationStatus = "COMMITTED"
	AIOperationStatusReleased  AIOperationStatus = "RELEASED"
)

// AIWarning flags something about an AI report that did not stop it.
type AIWarning string

// AI report warnings.
const (
	AIWarningModelDeprecated        AIWarning = "MODEL_DEPRECATED"
	AIWarningEstimateIncomplete     AIWarning = "ESTIMATE_INCOMPLETE"
	AIWarningConversionInputMissing AIWarning = "CONVERSION_INPUT_MISSING"
	AIWarningReservationExpired     AIWarning = "RESERVATION_EXPIRED"
	AIWarningOverdrawn              AIWarning = "OVERDRAWN"
	AIWarningModelNotAuthorized     AIWarning = "MODEL_NOT_AUTHORIZED"
	AIWarningUnknownModel           AIWarning = "UNKNOWN_MODEL"

	// AIWarningNoReservation is a release (a NOT_SENT commit) whose transaction id held
	// no reservation. Nothing was charged; the id is closed, so a reserve or a commit
	// sent under it later is refused with ErrTransactionIDReused.
	AIWarningNoReservation AIWarning = "NO_RESERVATION"

	// AIWarningCounterPrecisionDegraded is a commit the pool's counter could not hold
	// to the credit's exact decimals. The commit was accepted and its exact charge is
	// in AICredits.Charged; the counter carries the nearest value it can represent.
	AIWarningCounterPrecisionDegraded AIWarning = "COUNTER_PRECISION_DEGRADED"

	// AIWarningPricingSnapshotUnreadable is read side only: an inference listed without
	// its pricing snapshot, because the stored snapshot could not be decoded.
	AIWarningPricingSnapshotUnreadable AIWarning = "PRICING_SNAPSHOT_UNREADABLE"
)

// AIRuleSpecificity is how the conversion rule that priced a phase matched it. The set
// is the contract's; a value this SDK does not know is passed through unchanged.
type AIRuleSpecificity string

// AI rule specificities, most specific first.
const (
	// AIRuleSpecificityExactTyped is a rule for the call's model and request type.
	AIRuleSpecificityExactTyped AIRuleSpecificity = "EXACT_TYPED"
	// AIRuleSpecificityExactAll is a rule for the call's model, every request type.
	AIRuleSpecificityExactAll AIRuleSpecificity = "EXACT_ALL"
	// AIRuleSpecificityAnyTyped is a rule for any model, the call's request type.
	AIRuleSpecificityAnyTyped AIRuleSpecificity = "ANY_TYPED"
	// AIRuleSpecificityAnyAll is a rule for any model and any request type.
	AIRuleSpecificityAnyAll AIRuleSpecificity = "ANY_ALL"
)

// AIRule names the conversion rule that priced a phase.
type AIRule struct {
	ID          string            `json:"id"`
	Specificity AIRuleSpecificity `json:"specificity"`
}

// AIReport is the API's answer to a reserve or a commit.
type AIReport struct {
	// Usage is the pool's position after the report. A reserve does not move its
	// counter; what it set aside is in Reserved.
	Usage EntitlementUsage

	// Reserved is the sum of the pool's open, unexpired reservations, as a decimal
	// string.
	Reserved string

	// Outcome is what this report did to its operation.
	Outcome AIOutcome

	// Replayed is true when the API answered a retry of a report it had already
	// recorded: the original answer, nothing recomputed.
	Replayed bool
}

// AIOutcome is what one AI report did to its operation.
type AIOutcome struct {
	TransactionID string            `json:"transactionId"`
	Phase         AIPhase           `json:"phase"`
	Status        AIOperationStatus `json:"status"`
	Credits       AICredits         `json:"credits"`

	// Remaining is the pool's limit less its value and its reservations, as a decimal
	// string: negative when overdrawn, nil when the pool is unlimited.
	Remaining *string `json:"remaining"`

	// ExpiresAt is when a reservation lapses. Set on a reserve only.
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// Authorized is false on a commit no conversion rule priced: an unknown or archived
	// model, or no rule for it. Nil on a reserve, which is refused instead.
	Authorized *bool `json:"authorized,omitempty"`

	// Rule is the conversion rule that priced this phase. Nil when no rule applied, on
	// a server that predates the member, and on the replay of a phase answered before
	// it existed.
	Rule *AIRule `json:"rule,omitempty"`

	// Warnings is an open set: compare against the AIWarning constants and tolerate
	// values this SDK does not know.
	Warnings []AIWarning `json:"warnings"`
}

// HasWarning reports whether the report carried warning.
func (o AIOutcome) HasWarning(warning AIWarning) bool {
	return slices.Contains(o.Warnings, warning)
}

// NoReservation reports whether the report was a release that closed no reservation:
// RELEASED with AIWarningNoReservation. Nothing was charged, and the transaction id is
// now closed: a reserve still in flight under it will be refused, not held.
func (o AIOutcome) NoReservation() bool {
	return o.Status == AIOperationStatusReleased && o.HasWarning(AIWarningNoReservation)
}

// AICredits are the credits an AI report moved, as decimal strings with at most six
// decimals. Empty when the report did not set the member.
type AICredits struct {
	// Reserved is what a reserve set aside, or what a commit closed.
	Reserved string `json:"reserved,omitempty"`
	// Computed is what the rule priced the commit's usage at, before its policy.
	Computed string `json:"computed,omitempty"`
	// Charged is what the commit added to the pool's counter.
	Charged string `json:"charged,omitempty"`
}
