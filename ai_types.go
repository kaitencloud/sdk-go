// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import "time"

// The types in this file are written by hand against the AI credits contract, because
// internal/gen predates it. Once the Core OpenAPI document carries AiUsageValue and
// AiUsageOutcome, the enums become aliases of the generated ones, as in types.go, and
// the response members decode from the generated EntitlementUsage.

// AIRequestType is the kind of provider call an AI usage report meters.
type AIRequestType string

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
type AIRequestStatus string

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
type AIUnitType string

// AI unit types.
const (
	AIUnitTypeImage  AIUnitType = "image"
	AIUnitTypeSecond AIUnitType = "second"
	AIUnitTypeFrame  AIUnitType = "frame"
)

// AIPhase is the phase of an AI usage report.
type AIPhase string

// AI usage report phases.
const (
	AIPhaseReserve AIPhase = "reserve"
	AIPhaseCommit  AIPhase = "commit"
)

// AIOperationStatus is the state of an AI operation after the report that answered it.
type AIOperationStatus string

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
)

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

	Warnings []AIWarning `json:"warnings"`
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
