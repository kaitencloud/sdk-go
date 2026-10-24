// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import "github.com/kaitencloud/sdk-go/internal/gen"

// These aliases make the generated AI catalogue, management and query operations on
// Client.Raw() usable without importing internal/gen. They are the wire contract,
// not a second pricing implementation.

// Model is the catalogue's create, full-replace update and read body. Keep its slug
// and modality unchanged on update. Vendor costs are informational, not credit rates.
type Model = gen.Model

// ModelModality describes a catalogue model's modality.
type ModelModality = gen.ModelModality

// ModelStatus describes the reversible ACTIVE, DEPRECATED and ARCHIVED lifecycle.
type ModelStatus = gen.ModelStatus

// VendorCosts carries exact decimal USD strings, not conversion rates.
type VendorCosts = gen.VendorCosts

// VendorCostsCurrency is the catalogue cost currency (USD).
type VendorCostsCurrency = gen.VendorCostsCurrency

// ListModelsParams filters and paginates Client.Raw().ListModelsWithResponse.
type ListModelsParams = gen.ListModelsParams

// ListModelsParamsStatus filters catalogue lifecycle state.
type ListModelsParamsStatus = gen.ListModelsParamsStatus

// ListModelsParamsModality filters catalogue modality.
type ListModelsParamsModality = gen.ListModelsParamsModality

// ConversionRates preserves decimal strings and nil (unpriced) versus "0".
// Missing cache rates fall back to the input-token rate on the server.
type ConversionRates = gen.ConversionRates

// ConversionRule belongs to one license version and entitlement, never the family.
// A nil model or request type is a wildcard; unknown models are never authorized.
type ConversionRule = gen.ConversionRule

// ConversionRuleChargeOnError selects the server's error charging policy.
type ConversionRuleChargeOnError = gen.ConversionRuleChargeOnError

// ConversionRuleRequestType selects a priced request type.
type ConversionRuleRequestType = gen.ConversionRuleRequestType

// ConversionRuleStatus is ACTIVE or ARCHIVED.
type ConversionRuleStatus = gen.ConversionRuleStatus

// CreateConversionRule is the body of Client.Raw().CreateConversionRuleWithResponse.
type CreateConversionRule = gen.CreateConversionRule

// CreateConversionRuleChargeOnError selects the create request's error policy.
type CreateConversionRuleChargeOnError = gen.CreateConversionRuleChargeOnError

// CreateConversionRuleRequestType selects the create request's request type.
type CreateConversionRuleRequestType = gen.CreateConversionRuleRequestType

// UpdateConversionRule is the body of Client.Raw().UpdateConversionRuleWithResponse.
// Rates are replaced, not merged; send the complete desired rate vector.
type UpdateConversionRule = gen.UpdateConversionRule

// UpdateConversionRuleChargeOnError selects the update request's error policy.
type UpdateConversionRuleChargeOnError = gen.UpdateConversionRuleChargeOnError

// ListConversionRulesParams filters the rules of a specific license version slug.
type ListConversionRulesParams = gen.ListConversionRulesParams

// ListConversionRulesParamsStatus filters active or archived rules.
type ListConversionRulesParamsStatus = gen.ListConversionRulesParamsStatus

// AIInference is one reserve/commit operation, including exact decimal credits.
type AIInference = gen.AiInference

// AIInferencePage is a page of operations, with its opaque next cursor.
type AIInferencePage = gen.AiInferencePage

// AIReservations lists open, unexpired operations, not the numeric usage counter.
type AIReservations = gen.AiReservations

// AIInferenceRequestType is an operation's request type.
type AIInferenceRequestType = gen.AiInferenceRequestType

// AIInferenceStatus includes RESERVED, COMMITTED, RELEASED and EXPIRED.
type AIInferenceStatus = gen.AiInferenceStatus

// ListAIInferencesParams filters and paginates Client.Raw().ListAiInferencesWithResponse.
type ListAIInferencesParams = gen.ListAiInferencesParams

// ListAIInferencesParamsStatus filters an operation's lifecycle state.
type ListAIInferencesParamsStatus = gen.ListAiInferencesParamsStatus

// AIUsageQuantities is the generated usage shape for raw reports and inference reads.
// Unlike the convenience AIUsage, pointers preserve explicit zero versus unknown.
// Token dimensions are disjoint; Units is an exact decimal string.
type AIUsageQuantities = gen.AiUsage

// ModelUpdated is the before/after catalogue snapshot used by update hooks.
type ModelUpdated = gen.ModelUpdated

// ConversionRuleUpdated is the before/after rule snapshot used by update hooks.
type ConversionRuleUpdated = gen.ConversionRuleUpdated

// OnModelCreatedJSONBody is the model-created webhook payload.
type OnModelCreatedJSONBody = gen.OnModelCreatedJSONBody

// OnModelUpdatedJSONBody is the model-updated webhook payload.
type OnModelUpdatedJSONBody = gen.OnModelUpdatedJSONBody

// OnModelDeletedJSONBody is the model-deleted webhook payload.
type OnModelDeletedJSONBody = gen.OnModelDeletedJSONBody

// OnCreateConversionRuleJSONBody is the rule-created webhook payload.
type OnCreateConversionRuleJSONBody = gen.OnCreateConversionRuleJSONBody

// OnUpdateConversionRuleJSONBody is the rule-updated webhook payload.
type OnUpdateConversionRuleJSONBody = gen.OnUpdateConversionRuleJSONBody

// OnArchiveConversionRuleJSONBody is the rule-archived webhook payload.
type OnArchiveConversionRuleJSONBody = gen.OnArchiveConversionRuleJSONBody

// OnInstanceAIModelUsageUnauthorizedJSONBody records an unauthorized AI commit.
type OnInstanceAIModelUsageUnauthorizedJSONBody = gen.OnInstanceAIModelUsageUnauthorizedJSONBody

// OnInstanceEntitlementUsageReportAcceptedJSONBody carries numeric usage events,
// including AI credit charges; its value is not an ai_usage input.
type OnInstanceEntitlementUsageReportAcceptedJSONBody = gen.OnInstanceEntitlementUsageReportAcceptedJSONBody

// AIUsageValue is the write-only raw reserve/commit variant. Use Usage pointers when
// an explicitly known zero must be distinguished from an absent quantity.
type AIUsageValue = gen.AiUsageValue

// AIUsageValuePhase selects reserve or commit on a raw input.
type AIUsageValuePhase = gen.AiUsageValuePhase

// AIUsageValueType is the write-only ai_usage discriminator.
type AIUsageValueType = gen.AiUsageValueType

// ReportEntitlementUsageBody is the raw usage report request, not a read-value union.
type ReportEntitlementUsageBody = gen.ReportEntitlementUsageBody

// ReportEntitlementUsageValue is the write union used by ReportEntitlementUsageBody.
type ReportEntitlementUsageValue = gen.ReportEntitlementUsageBody_Value
