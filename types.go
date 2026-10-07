// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"github.com/kaitencloud/sdk-go/internal/gen"
	"github.com/kaitencloud/sdk-go/internal/genplatform"
)

// AuditTrail records a single event in an instance's audit trail.
type AuditTrail = gen.AuditTrail

// Basic is a rollout strategy that targets either all users or no users.
type Basic = gen.Basic

// BasicType identifies a Basic rollout strategy.
type BasicType = gen.BasicType

// BasicTargeting wraps a Basic rollout strategy for targeting.
type BasicTargeting = gen.BasicTargeting

// BasicTargetingType identifies a BasicTargeting strategy.
type BasicTargetingType = gen.BasicTargetingType

// BooleanEntitlementValue is a boolean entitlement value.
type BooleanEntitlementValue = gen.BooleanEntitlementValue

// BooleanEntitlementValueType identifies a BooleanEntitlementValue.
type BooleanEntitlementValueType = gen.BooleanEntitlementValueType

// Component is a product component that can be included in releases.
type Component = gen.Component

// ConfigEntitlementValue is an object (config) entitlement value.
type ConfigEntitlementValue = gen.ConfigEntitlementValue

// ConfigEntitlementValueType identifies a ConfigEntitlementValue.
type ConfigEntitlementValueType = gen.ConfigEntitlementValueType

// Customer is an organization's customer.
type Customer = gen.Customer

// DefaultVariant is the variant served when no rollout targeting matches.
type DefaultVariant = gen.DefaultVariant

// DeploymentZone is a region or environment instances can be deployed to.
type DeploymentZone = gen.DeploymentZone

// Entitlement is a feature or limit that can be granted to customers via licenses.
type Entitlement = gen.Entitlement

// EntitlementAggregationMethod describes how usage of an entitlement is aggregated.
type EntitlementAggregationMethod = gen.EntitlementAggregationMethod

// EntitlementGroup is a named collection of entitlements.
type EntitlementGroup = gen.EntitlementGroup

// EntitlementGroupRef references an entitlement group.
//
// The Core spec calls this schema EntitlementGroupSummary. The Go name is kept
// as it is: it is a published identifier, and renaming it would break every
// consumer for a wire name they never see.
type EntitlementGroupRef = gen.EntitlementGroupSummary

// EntitlementGroupUsageItem reports usage of one entitlement within an entitlement group.
//
// The Core spec calls this schema EntitlementGroupUsage; the Go name is kept
// for the same reason as EntitlementGroupRef above.
type EntitlementGroupUsageItem = gen.EntitlementGroupUsage

// EntitlementType describes the value type of an entitlement.
type EntitlementType = gen.EntitlementType

// EntitlementUsage reports current usage for an entitlement.
type EntitlementUsage = gen.EntitlementUsage

// EntitlementUsageBehavior describes how a reported usage value is applied.
//
// report-entitlement-usage-metric declares its own request type
// (ReportEntitlementUsageBody) rather than reusing Entitlement or
// EntitlementUsage: its response is a computed snapshot, not an echo of the
// request, so there is no shared resource schema to fold into. oapi-codegen
// names this operation-scoped enum after that type.
type EntitlementUsageBehavior = gen.ReportEntitlementUsageBodyBehavior

// EntitlementUsageValue is the value reported for entitlement usage.
type EntitlementUsageValue = gen.EntitlementUsage_Value

// ErrorDetail describes a single validation error.
type ErrorDetail = gen.ErrorDetail

// ErrorModel is the RFC 9457 problem-details payload returned by the API on errors.
//
// The Core spec calls this schema Problem, and every error response shares it.
// The Go name is kept as it is: it is a published identifier, and renaming it
// would break every consumer for a wire name they never see.
type ErrorModel = gen.Problem

// FeatureFlag is a feature flag definition.
type FeatureFlag = gen.FeatureFlag

// FeatureFlagType describes the value type of a feature flag.
type FeatureFlagType = gen.FeatureFlagType

// Instance is a deployed instance of the product for a customer.
type Instance = gen.Instance

// InstanceStatus is the health status of an instance.
type InstanceStatus = gen.InstanceStatus

// License grants a customer a set of entitlements.
type License = gen.License

// MetadataField is a custom metadata field definition for a resource type.
type MetadataField = gen.MetadataField

// MetadataFieldResourceType identifies the resource type a metadata field applies to.
type MetadataFieldResourceType = gen.MetadataFieldResourceType

// MetadataFieldDryRunImpact reports the impact of a metadata field schema change.
type MetadataFieldDryRunImpact = gen.Impact

// MetadataFieldDryRunImpactSample is a sample resource affected by a metadata field schema change.
type MetadataFieldDryRunImpactSample = gen.ImpactSample

// LicenseFamilyView is a license family together with the version it currently
// resolves to. The family is the product; a License is one version of it.
type LicenseFamilyView = gen.LicenseFamilyView

// LicenseLifecycleState says whether a version of a license may be served.
// Several versions of one family may be Published at once.
type LicenseLifecycleState = gen.LicenseLifecycleState

// Lifecycle states a license version can be in.
const (
	LicenseLifecycleDraft     LicenseLifecycleState = "DRAFT"
	LicenseLifecyclePublished LicenseLifecycleState = "PUBLISHED"
	LicenseLifecycleArchived  LicenseLifecycleState = "ARCHIVED"
)

// LicenseEntitlement is an entitlement value granted by a license.
type LicenseEntitlement = gen.LicenseEntitlement

// LicenseEntitlementEntitlementType describes the value type of a license entitlement.
type LicenseEntitlementEntitlementType = gen.LicenseEntitlementEntitlementType

// LicenseEntitlementValue is the value granted by a license entitlement.
type LicenseEntitlementValue = gen.LicenseEntitlement_Value

// LicenseType describes the type of a license.
type LicenseType = gen.LicenseType

// NumberEntitlementValue is a numeric entitlement value.
type NumberEntitlementValue = gen.NumberEntitlementValue

// NumberEntitlementValueType identifies a NumberEntitlementValue.
type NumberEntitlementValueType = gen.NumberEntitlementValueType

// IssuedToken is an organization-scoped token as the Platform API reports it:
// identity and authority, never the secret. It is what PlatformTokens.List
// returns, and the way back from the Name a caller chose to the Slug revocation
// is addressed by -- a slug is server-generated and cannot be derived from a
// name.
type IssuedToken = genplatform.Token

// MintedToken is an organization-scoped token including its plaintext secret, returned
// only by PlatformTokens.Mint. Token is populated on that response and nowhere else;
// Slug, not Name, is what PlatformTokens.Revoke takes.
type MintedToken = genplatform.PlainToken

// Connector is a registered connector manifest. Sourced from the Platform API's schema,
// which is where connectors are registered -- the Core API's own Connector describes the
// same row to the organization reading it, and the two documents are free to diverge.
type Connector = genplatform.Connector

// Organization is a Kaiten organization. Sourced from the Platform API's schema, which is
// where organizations are administered.
type Organization = genplatform.Organization

// PlainToken is a service account token including its plaintext secret, returned only on creation.
type PlainToken = gen.PlainToken

// PlatformCredential describes the platform credential (`ksm_...`) a PlatformClient is
// authenticating with, as returned by PlatformClient.Me. Never carries a token value.
type PlatformCredential = genplatform.PlatformCredential

// Release is a versioned collection of components.
type Release = gen.Release

// RolloutDate is a rollout strategy that activates on a fixed date.
type RolloutDate = gen.RolloutDate

// RolloutDateTargeting wraps a RolloutDate strategy for targeting.
type RolloutDateTargeting = gen.RolloutDateTargeting

// RolloutDateTargetingType identifies a RolloutDateTargeting strategy.
type RolloutDateTargetingType = gen.RolloutDateTargetingType

// RolloutDateType identifies a RolloutDate strategy.
type RolloutDateType = gen.RolloutDateType

// RolloutPercentage is a rollout strategy that activates for a percentage of users.
type RolloutPercentage = gen.RolloutPercentage

// RolloutPercentageTargeting wraps a RolloutPercentage strategy for targeting.
type RolloutPercentageTargeting = gen.RolloutPercentageTargeting

// RolloutPercentageTargetingType identifies a RolloutPercentageTargeting strategy.
type RolloutPercentageTargetingType = gen.RolloutPercentageTargetingType

// RolloutPercentageType identifies a RolloutPercentage strategy.
type RolloutPercentageType = gen.RolloutPercentageType

// RolloutStep is one step of a gradual rollout.
type RolloutStep = gen.RolloutStep

// ServiceAccount is a service account used for API authentication.
type ServiceAccount = gen.ServiceAccount

// Targetings lists the rollout targeting strategies for a feature flag variant.
type Targetings = gen.Targetings

// TargetingsItem is a single rollout targeting strategy.
type TargetingsItem = gen.Targetings_Item

// Token is a service account token.
type Token = gen.Token

// User is a Kaiten user.
type User = gen.User

// Variant is a possible value of a feature flag.
type Variant = gen.Variant

// Entitlement usage behaviors, value type discriminators, and instance status values.
const (
	Append EntitlementUsageBehavior = gen.ReportEntitlementUsageBodyBehaviorAppend
	Set    EntitlementUsageBehavior = gen.ReportEntitlementUsageBodyBehaviorSet

	Boolean BooleanEntitlementValueType = "boolean"
	Number  NumberEntitlementValueType  = "number"
	Object  ConfigEntitlementValueType  = "object"

	InstanceStatusHealthy     InstanceStatus = gen.InstanceStatusHEALTHY
	InstanceStatusDegraded    InstanceStatus = gen.InstanceStatusDEGRADED
	InstanceStatusIncident    InstanceStatus = gen.InstanceStatusINCIDENT
	InstanceStatusMaintenance InstanceStatus = gen.InstanceStatusMAINTENANCE

	FeatureFlagTypeBoolean FeatureFlagType = gen.FeatureFlagTypeBoolean
	FeatureFlagTypeNumber  FeatureFlagType = gen.FeatureFlagTypeNumber
	FeatureFlagTypeObject  FeatureFlagType = gen.FeatureFlagTypeObject
	FeatureFlagTypeString  FeatureFlagType = gen.FeatureFlagTypeString
)

// The rollout-strategy discriminators are unexported: they are the tag inside a union,
// and BasicDefaultVariant and BasicVariantTargeting are how a caller reaches that arm.
const (
	basicStrategy          BasicType          = gen.BasicTypeBasic
	basicTargetingStrategy BasicTargetingType = gen.BasicTargetingTypeBasic
)
