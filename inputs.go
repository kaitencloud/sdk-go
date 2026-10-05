package sdk

import "time"

// InstancesListOptions configures Instances.List.
type InstancesListOptions struct {
	// IncludeDeleted is sent as include_deleted, which the Core API no longer
	// declares: GET /instances now takes only cursor and limit, and huma
	// ignores a query parameter it has no field for. So this currently
	// changes nothing -- soft-deleted instances are excluded either way. Kept
	// and still sent rather than removed, because the alternative is a
	// breaking change to every caller for a field that costs nothing, and
	// re-adding the parameter to Core would make it work again untouched.
	IncludeDeleted bool
}

// includeVersions is the only value the family endpoint's include list accepts
// today. Named here so the wire value appears once, rather than at the call
// site and again in whatever test asserts it.
const includeVersions = "versions"

// LicenseFamilyOptions configures LicenseFamilies.Get. Nil, or a zero value,
// asks for the version the family currently serves.
type LicenseFamilyOptions struct {
	// Version resolves one version by number instead of the current one. The
	// version comes back whatever its lifecycle state, archived included --
	// addressing a version explicitly is how a caller reaches history.
	Version *int32

	// IncludeVersions fills LicenseFamilyView.Versions with every version of
	// the family, oldest first. Off by default: it is the lifecycle view, not
	// what a catalogue needs to render one product.
	IncludeVersions bool
}

// AuditTrailsOptions configures Instances.ListAuditTrails.
type AuditTrailsOptions struct {
	EventName string
	After     *time.Time
	Before    *time.Time

	// Limit caps how many entries are returned. Unlike every other list in
	// this package, audit trails are not walked to completion: the collection
	// grows without bound, so "all of them" is not a useful default and a
	// caller's ceiling is honoured as one.
	Limit *int32

	// Offset is dead. Audit trails moved from offset/limit to cursor
	// pagination, so the Core API no longer declares this parameter and
	// ignores it -- an offset that used to skip entries now returns the newest
	// ones. Left in place so callers still compile. Use Limit instead: it
	// walks as many pages as it takes to reach the count asked for, which
	// covers what offset was reached for in practice, and only resuming from
	// where a previous call stopped is no longer expressible.
	Offset *int32
}

// ComponentInput is the input for creating or updating a component.
type ComponentInput struct {
	Name                string
	Version             string
	Description         *string
	PreviousComponentID *string
	Slug                *string
}

// CustomerInput is the input for creating or updating a customer.
type CustomerInput struct {
	Name               string
	ExternalCustomerID *string
	Slug               *string
}

// DeploymentZoneInput is the input for creating or updating a deployment zone.
type DeploymentZoneInput struct {
	Name        string
	Type        string
	Description string
	Metadata    map[string]any
	ReleaseID   *string
	Slug        *string
}

// EnsureOrganizationInput is the input for ensuring a platform organization exists.
type EnsureOrganizationInput struct {
	// ExternalID is the identity provider's id for the organization (e.g. "org_2abcDEF").
	// The organization's own id is derived from it, so this value is the identity.
	ExternalID string

	// Name is applied only when this call creates the organization; it never renames an
	// existing one. Omitted, a new organization is named after its external id.
	Name *string
}

// EntitlementGroupInput is the input for creating or updating an entitlement group.
type EntitlementGroupInput struct {
	Name        string
	Description *string
	Slug        *string
}

// EntitlementInput is the input for creating or updating an entitlement.
//
// The presentation fields -- Icon, UnitSingular, UnitPlural, UserFacing and DisplayOrder --
// are what customer-facing components render: a pricing table shows the user-facing
// entitlements in DisplayOrder with their unit labels, and hides the rest. They are
// optional because most entitlements are internal counters that should stay hidden, which
// is exactly what leaving them nil produces.
//
// The API validates the unit labels as an all-or-none pair, so set both or neither.
type EntitlementInput struct {
	Name              string
	Description       *string
	Type              *EntitlementType
	AggregationMethod *EntitlementAggregationMethod
	GroupSlugs        []string
	Slug              *string

	// Icon is a provider-namespaced token such as "lucide:rocket", deliberately not tied
	// to any one icon library.
	Icon *string

	// UnitSingular and UnitPlural label the base unit a NUMBER entitlement is measured in
	// ("seat" / "seats"). Both or neither.
	UnitSingular *string
	UnitPlural   *string

	// UserFacing exposes the entitlement in customer-facing components. Defaults to false
	// server-side, which keeps raw counters out of a pricing page.
	UserFacing *bool

	// DisplayOrder sorts the entitlement within customer-facing components, ascending.
	DisplayOrder *int32
}

// InstanceInput is the input for creating or updating an instance.
type InstanceInput struct {
	Name             string
	Description      string
	CustomerID       string
	LicenseID        string
	DeploymentZoneID *string

	// Metadata is the instance's org-defined data, and the only such field.
	//
	// It absorbed `platform`, which was the same thing twice: two free-form
	// JSONB blobs on one row, both org-defined, both surfaced to feature-flag
	// targeting, with nothing to say which key belonged in which. The API stopped
	// accepting `platform` in kaiten's "replace platform with typed metadata",
	// and its request bodies refuse unknown properties -- so a client still
	// sending it gets a 422 on every create and every update, not a field
	// quietly ignored.
	Metadata         map[string]any
	StartLicenseDate time.Time
	EndLicenseDate   time.Time

	// Slug is the instance's URL-friendly identity, auto-generated when nil on create and
	// a rename on update -- nil there keeps the current slug, and a slug another instance
	// already holds is a 409.
	Slug *string
}

// LicenseInput is the input for creating or updating a license.
//
// A license is one version of a family, and the family is the product -- see
// LicenseFamilies.Get. Create opens a new family unless FamilySlug or FamilyID names
// an existing one, in which case the license becomes that family's next version.
// Update never moves a version between families.
//
// Several fields are not symmetric between the two operations, because the API's two
// bodies are not either:
//
//   - Version is read-only. A family's first version is 1 and the server assigns the
//     next number from there, so the create body declares no such field at all.
//     Update sends it only as an echo of what was read, since a version's slug is
//     built from its number.
//   - Slug is accepted on create and refused on update. A slug is the license's stable
//     identity once assigned.
//   - FamilySlug is accepted on create and refused on update, so updatePayload never
//     renders it. FamilyID is accepted on both, but on update only as the license's
//     own family.
//   - LifecycleState is create-only, and DRAFT or PUBLISHED there. After that the
//     state moves through Licenses.Publish, Archive and Unarchive, so updatePayload
//     never renders it.
//
// There is no IsActive: no request body publishes one, and both bodies are
// additionalProperties:false, so sending it was a 422 on every call.
type LicenseInput struct {
	Name        string
	Description string
	Type        LicenseType

	// Version is the license's version within its family, assigned by the server
	// and never changed. Create ignores it. Update sends it, and the API accepts it
	// empty or equal to the stored number -- so a License read back can be written
	// back as it is -- and refuses any other value with a 422,
	// UpdateLicense.VersionNotSettable.
	Version string

	VersionName *string
	IsDefault   bool

	// Slug is the license's URL-friendly identity, auto-generated when nil -- create only.
	// Update refuses it.
	Slug *string

	// FamilySlug names, by slug, the family this license becomes the next version of
	// -- create only. Nil, with FamilyID also nil, opens a new family of which the
	// license is version 1. It is the handle an integrator holds for a product it
	// published: a family's slug is served by the family endpoints, never on a
	// License.
	//
	// An unknown family is a 404, code CreateLicense.FamilyNotFound. Sent together
	// with FamilyID the two must name the same family, or the create is a 422,
	// CreateLicense.FamilyMismatch -- the server refuses to guess which was meant.
	// Update refuses the field outright (422 UpdateLicense.FamilyNotReassignable),
	// which is why updatePayload drops it rather than let an input filled once for
	// both calls fail every update.
	FamilySlug *string

	// FamilyID names the same family by identifier: the handle a caller holding a
	// License already has, since every License carries its FamilyId. On create it is
	// the alternative to FamilySlug, under the same rules -- an unknown family is
	// CreateLicense.FamilyNotFound, and one that disagrees with FamilySlug is
	// CreateLicense.FamilyMismatch. On update it is optional and must be the
	// license's own family: it is accepted there so that a caller writing back what
	// it read is not refused, and any other value is a 422,
	// UpdateLicense.FamilyNotReassignable, because a version cannot move between
	// families.
	FamilyID *string

	// LifecycleState says whether the version may be served -- create only. Nil is
	// PUBLISHED, the server's default: a license is usable the moment it exists, and
	// a vendor preparing a version says LicenseLifecycleDraft explicitly. Those are
	// the only two states a version is created in. LicenseLifecycleArchived is a
	// 422, CreateLicense.LifecycleStateNotSettable, because a version is archived by
	// withdrawing it from sale: create it, then call Licenses.Archive. Only a
	// published version can be its family's default, so IsDefault with a draft is
	// a 409, CreateLicense.DefaultMustBePublished.
	//
	// Update never sends it. A version's state then moves through Licenses.Publish
	// (DRAFT to PUBLISHED), Archive (PUBLISHED to ARCHIVED) and Unarchive (ARCHIVED
	// to PUBLISHED), which apply the lifecycle rules and record their own events.
	LifecycleState *LicenseLifecycleState
}

// MetadataFieldInput is the input for creating a metadata field.
type MetadataFieldInput struct {
	ResourceType MetadataFieldResourceType
	Key          string
	Label        string
	JSONSchema   map[string]any
	DisplayOrder int32
}

// MetadataFieldUpdateInput is the input for updating a metadata field.
type MetadataFieldUpdateInput struct {
	Label      string
	JSONSchema map[string]any
}

// ReleaseInput is the input for creating a release.
type ReleaseInput struct {
	Version      string
	Description  *string
	ComponentIDs []string
	Slug         *string
}

// ServiceAccountInput is the input for creating or updating a service account.
type ServiceAccountInput struct {
	Name string
	Slug *string
}

// TokenInput is the input for creating a service account token.
type TokenInput struct {
	Name      string
	Slug      *string
	Scopes    []string
	ExpiresAt *time.Time
}

// MintTokenInput is the input for minting an organization-scoped token over the Platform
// API.
type MintTokenInput struct {
	// Name is a human-readable label, unique among an organization's active tokens.
	Name string

	// Scopes must be a subset of the minting platform credential's scopes. Nil inherits
	// all of them, which is rarely what a machine consumer should hold.
	Scopes []string

	// TTL bounds the token's lifetime. Nil mints a non-expiring token, whose only bound
	// is revocation -- including the cascade from the platform credential that minted it.
	TTL *time.Duration
}

// RegisterConnectorInput is the input for registering a connector's manifest over the
// Platform API.
//
// Registration is deployment-wide, not per organization: it declares that a connector
// exists here and what its settings look like. Which organizations may use it is a
// separate, organization-scoped decision -- see EntitlementSlug, and the activation
// endpoints on the Core API.
type RegisterConnectorInput struct {
	// Name is the connector's stable identifier, e.g. "kaiten.integration.crm.attio".
	// It is the primary key of the registration: registering the same name again
	// updates the manifest rather than adding a second connector.
	Name string

	// Version is what an operator reads to tell which build of the connector a
	// deployment is running. It carries no behaviour of its own.
	Version string

	// SettingsSchema is the JSON schema an organization's settings are validated
	// against, and what a console renders the settings form from. A `writeOnly`
	// property is a secret: the API accepts it and never reads it back out.
	SettingsSchema map[string]any

	// EntitlementSlug names the BOOLEAN entitlement a licence must grant before an
	// organization may activate this connector. Nil leaves the connector ungated,
	// which is the open-source default -- a deployment that licenses nothing must
	// still be able to use its connectors.
	EntitlementSlug *string
}

// UsageReportInput is the input for reporting entitlement usage.
type UsageReportInput struct {
	Value    float64
	Behavior EntitlementUsageBehavior
	// Metadata is stored with the report in the usage history when its compact
	// JSON encoding is at most 4 KiB; above that the report still counts and
	// UsageReportResult.MetadataDropped is set. It must contain no personal data.
	Metadata map[string]any
	// TransactionID makes the report idempotent: Kaiten applies a report at
	// most once per key, instance and entitlement within its idempotency window
	// (35 days by default), and answers a repeat with the original result. Use a
	// UUID, or a business event ID plus the meter ("llm-call:9f2c:tokens"). A
	// report with a key is retried on transient failures; one without is never
	// retried. 1 to 128 characters of letters, digits, '.', '_', ':' and '-'.
	TransactionID string
}

type componentCreatePayload struct {
	Name                string  `json:"name"`
	Version             string  `json:"version"`
	Description         *string `json:"description,omitempty"`
	PreviousComponentID *string `json:"previousComponentId,omitempty"`
	Slug                *string `json:"slug,omitempty"`
}

func (in ComponentInput) createPayload() componentCreatePayload {
	return componentCreatePayload(in)
}

// componentUpdatePayload drops PreviousComponentID: the predecessor link is established
// when a component is created and UpdateComponentBody does not publish it.
type componentUpdatePayload struct {
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Description *string `json:"description,omitempty"`
	Slug        *string `json:"slug,omitempty"`
}

func (in ComponentInput) updatePayload() componentUpdatePayload {
	return componentUpdatePayload{
		Name:        in.Name,
		Version:     in.Version,
		Description: in.Description,
		Slug:        in.Slug,
	}
}

type customerCreatePayload struct {
	Name               string  `json:"name"`
	ExternalCustomerID *string `json:"externalCustomerId,omitempty"`
	Slug               *string `json:"slug,omitempty"`
}

func (in CustomerInput) createPayload() customerCreatePayload {
	return customerCreatePayload(in)
}

type customerUpdatePayload struct {
	Name               string  `json:"name"`
	ExternalCustomerID *string `json:"externalCustomerId,omitempty"`
}

func (in CustomerInput) updatePayload() customerUpdatePayload {
	return customerUpdatePayload{
		Name:               in.Name,
		ExternalCustomerID: in.ExternalCustomerID,
	}
}

type deploymentZoneCreatePayload struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	ReleaseID   *string        `json:"releaseId,omitempty"`
	Slug        *string        `json:"slug,omitempty"`
}

func (in DeploymentZoneInput) createPayload() deploymentZoneCreatePayload {
	return deploymentZoneCreatePayload{
		Name:        in.Name,
		Type:        in.Type,
		Description: in.Description,
		Metadata:    cloneMap(in.Metadata),
		ReleaseID:   in.ReleaseID,
		Slug:        in.Slug,
	}
}

type deploymentZoneUpdatePayload struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	ReleaseID   *string        `json:"releaseId,omitempty"`
}

func (in DeploymentZoneInput) updatePayload() deploymentZoneUpdatePayload {
	return deploymentZoneUpdatePayload{
		Name:        in.Name,
		Type:        in.Type,
		Description: in.Description,
		Metadata:    cloneMap(in.Metadata),
		ReleaseID:   in.ReleaseID,
	}
}

// The description key is rendered even when nil, as JSON null. Both entitlement-group
// bodies declare description required and nullable, so omitting the key is a 422 while
// sending null is the sanctioned "no description" -- the one place where omitempty is
// wrong. The same is true of entitlementPayload below.
type entitlementGroupCreatePayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Slug        *string `json:"slug,omitempty"`
}

func (in EntitlementGroupInput) createPayload() entitlementGroupCreatePayload {
	return entitlementGroupCreatePayload(in)
}

type entitlementGroupUpdatePayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (in EntitlementGroupInput) updatePayload() entitlementGroupUpdatePayload {
	return entitlementGroupUpdatePayload{
		Name:        in.Name,
		Description: in.Description,
	}
}

type entitlementPayload struct {
	Name              string                        `json:"name"`
	Description       *string                       `json:"description"`
	Type              *EntitlementType              `json:"type,omitempty"`
	AggregationMethod *EntitlementAggregationMethod `json:"aggregationMethod,omitempty"`
	GroupSlugs        *[]string                     `json:"groupSlugs,omitempty"`
	Slug              *string                       `json:"slug,omitempty"`
	Icon              *string                       `json:"icon,omitempty"`
	UnitSingular      *string                       `json:"unitSingular,omitempty"`
	UnitPlural        *string                       `json:"unitPlural,omitempty"`
	UserFacing        *bool                         `json:"userFacing,omitempty"`
	DisplayOrder      *int32                        `json:"displayOrder,omitempty"`
}

func (in EntitlementInput) createPayload() entitlementPayload {
	return entitlementPayload{
		Name:              in.Name,
		Description:       in.Description,
		Type:              in.Type,
		AggregationMethod: in.AggregationMethod,
		GroupSlugs:        cloneStringSlice(in.GroupSlugs),
		Slug:              in.Slug,
		Icon:              in.Icon,
		UnitSingular:      in.UnitSingular,
		UnitPlural:        in.UnitPlural,
		UserFacing:        in.UserFacing,
		DisplayOrder:      in.DisplayOrder,
	}
}

// entitlementUpdatePayload is entitlementPayload without the slug: an entitlement's slug
// is its identity in the path, and UpdateEntitlementBody does not publish one.
type entitlementUpdatePayload struct {
	Name              string                        `json:"name"`
	Description       *string                       `json:"description"`
	Type              *EntitlementType              `json:"type,omitempty"`
	AggregationMethod *EntitlementAggregationMethod `json:"aggregationMethod,omitempty"`
	GroupSlugs        *[]string                     `json:"groupSlugs,omitempty"`
	Icon              *string                       `json:"icon,omitempty"`
	UnitSingular      *string                       `json:"unitSingular,omitempty"`
	UnitPlural        *string                       `json:"unitPlural,omitempty"`
	UserFacing        *bool                         `json:"userFacing,omitempty"`
	DisplayOrder      *int32                        `json:"displayOrder,omitempty"`
}

func (in EntitlementInput) updatePayload() entitlementUpdatePayload {
	return entitlementUpdatePayload{
		Name:              in.Name,
		Description:       in.Description,
		Type:              in.Type,
		AggregationMethod: in.AggregationMethod,
		GroupSlugs:        cloneStringSlice(in.GroupSlugs),
		Icon:              in.Icon,
		UnitSingular:      in.UnitSingular,
		UnitPlural:        in.UnitPlural,
		UserFacing:        in.UserFacing,
		DisplayOrder:      in.DisplayOrder,
	}
}

type instancePayload struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	CustomerID       string         `json:"customerId"`
	LicenseID        string         `json:"licenseId"`
	DeploymentZoneID *string        `json:"deploymentZoneId,omitempty"`
	Metadata         map[string]any `json:"metadata"`
	StartLicenseDate time.Time      `json:"startLicenseDate"`
	EndLicenseDate   time.Time      `json:"endLicenseDate"`
	Slug             *string        `json:"slug,omitempty"`
}

func (in InstanceInput) createPayload() instancePayload {
	return instancePayload{
		Name:             in.Name,
		Description:      in.Description,
		CustomerID:       in.CustomerID,
		LicenseID:        in.LicenseID,
		DeploymentZoneID: in.DeploymentZoneID,
		Metadata:         cloneMap(in.Metadata),
		StartLicenseDate: in.StartLicenseDate,
		EndLicenseDate:   in.EndLicenseDate,
		Slug:             in.Slug,
	}
}

// instanceUpdatePayload carries the slug like instancePayload does, because
// UpdateInstanceBody now publishes one: the same optional field as on create, read as a
// rename. It is omitempty because absence is how the body spells "leave the slug alone",
// so a nil pointer must render no key at all rather than an empty rename.
type instanceUpdatePayload struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	CustomerID       string         `json:"customerId"`
	LicenseID        string         `json:"licenseId"`
	DeploymentZoneID *string        `json:"deploymentZoneId,omitempty"`
	Metadata         map[string]any `json:"metadata"`
	StartLicenseDate time.Time      `json:"startLicenseDate"`
	EndLicenseDate   time.Time      `json:"endLicenseDate"`
	Slug             *string        `json:"slug,omitempty"`
}

func (in InstanceInput) updatePayload() instanceUpdatePayload {
	return instanceUpdatePayload{
		Name:             in.Name,
		Description:      in.Description,
		CustomerID:       in.CustomerID,
		LicenseID:        in.LicenseID,
		DeploymentZoneID: in.DeploymentZoneID,
		Metadata:         cloneMap(in.Metadata),
		StartLicenseDate: in.StartLicenseDate,
		EndLicenseDate:   in.EndLicenseDate,
		Slug:             in.Slug,
	}
}

// instancePatchPayload is the PATCH /instances/{instanceSlug} body.
//
// Rendered by hand rather than through internal/gen's PatchInstanceBody, whose
// generated field renders `lifecycle_stage` where the spec names the key
// `lifecycleStage`. Every Core request body is additionalProperties:false, so a
// key the spec does not declare is a 422 rather than a field the API ignores,
// and no type error catches it -- which is why both PATCH writes render this
// type and are pinned by TestRequestBodiesMatchTheirSpecBodies.
//
// Both fields are omitempty because absence is the only "leave this alone"
// signal the endpoint has -- sending a null clears nothing and fails
// validation, so a nil pointer must render no key at all.
type instancePatchPayload struct {
	Status         *InstanceStatus `json:"status,omitempty"`
	LifecycleStage *string         `json:"lifecycleStage,omitempty"`
}

// The three family-era fields are omitempty because absence is what each one's nil
// means on the wire: no family named opens a new one, and no lifecycleState is the
// server's PUBLISHED. A nil pointer has to render no key at all rather than a null --
// none of the three is nullable, so a null is a 422 like any other wrong value.
type licenseCreatePayload struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Type           LicenseType            `json:"type"`
	VersionName    *string                `json:"versionName,omitempty"`
	IsDefault      bool                   `json:"isDefault"`
	Slug           *string                `json:"slug,omitempty"`
	FamilySlug     *string                `json:"familySlug,omitempty"`
	FamilyID       *string                `json:"familyId,omitempty"`
	LifecycleState *LicenseLifecycleState `json:"lifecycleState,omitempty"`
}

func (in LicenseInput) createPayload() licenseCreatePayload {
	return licenseCreatePayload{
		Name:           in.Name,
		Description:    in.Description,
		Type:           in.Type,
		VersionName:    in.VersionName,
		IsDefault:      in.IsDefault,
		Slug:           in.Slug,
		FamilySlug:     in.FamilySlug,
		FamilyID:       in.FamilyID,
		LifecycleState: in.LifecycleState,
	}
}

// licenseUpdatePayload drops Slug, FamilySlug and LifecycleState. The slug is the
// license's identity in the path. familySlug is create-only: PUT refuses it with
// UpdateLicense.FamilyNotReassignable rather than ignore it, so rendering it from an
// input filled once for both calls would fail every update. lifecycleState moves
// through Licenses.Publish, Archive and Unarchive: PUT accepts the stored state and
// refuses any other with UpdateLicense.LifecycleStateNotSettable, so the state an
// input was created with would fail every update once the version has moved on.
// FamilyID stays, as the value a read-modify-write caller carries back.
type licenseUpdatePayload struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        LicenseType `json:"type"`
	Version     string      `json:"version"`
	VersionName *string     `json:"versionName,omitempty"`
	IsDefault   bool        `json:"isDefault"`
	FamilyID    *string     `json:"familyId,omitempty"`
}

func (in LicenseInput) updatePayload() licenseUpdatePayload {
	return licenseUpdatePayload{
		Name:        in.Name,
		Description: in.Description,
		Type:        in.Type,
		Version:     in.Version,
		VersionName: in.VersionName,
		IsDefault:   in.IsDefault,
		FamilyID:    in.FamilyID,
	}
}

type metadataFieldPayload struct {
	ResourceType MetadataFieldResourceType `json:"resourceType"`
	Key          string                    `json:"key"`
	Label        string                    `json:"label"`
	JSONSchema   map[string]any            `json:"jsonSchema"`
	DisplayOrder int32                     `json:"displayOrder"`
}

func (in MetadataFieldInput) payload() metadataFieldPayload {
	return metadataFieldPayload{
		ResourceType: in.ResourceType,
		Key:          in.Key,
		Label:        in.Label,
		JSONSchema:   cloneMap(in.JSONSchema),
		DisplayOrder: in.DisplayOrder,
	}
}

type metadataFieldUpdatePayload struct {
	Label      string         `json:"label"`
	JSONSchema map[string]any `json:"jsonSchema"`
}

func (in MetadataFieldUpdateInput) payload() metadataFieldUpdatePayload {
	return metadataFieldUpdatePayload{
		Label:      in.Label,
		JSONSchema: cloneMap(in.JSONSchema),
	}
}

type metadataFieldsReorderPayload struct {
	IDs []string `json:"ids"`
}

type releasePayload struct {
	Version      string    `json:"version"`
	Description  *string   `json:"description,omitempty"`
	ComponentIDs *[]string `json:"componentIds,omitempty"`
	Slug         *string   `json:"slug,omitempty"`
}

func (in ReleaseInput) payload() releasePayload {
	return releasePayload{
		Version:      in.Version,
		Description:  in.Description,
		ComponentIDs: cloneStringSlice(in.ComponentIDs),
		Slug:         in.Slug,
	}
}

type serviceAccountCreatePayload struct {
	Name string  `json:"name"`
	Slug *string `json:"slug,omitempty"`
}

func (in ServiceAccountInput) createPayload() serviceAccountCreatePayload {
	return serviceAccountCreatePayload(in)
}

type serviceAccountUpdatePayload struct {
	Name string `json:"name"`
}

func (in ServiceAccountInput) updatePayload() serviceAccountUpdatePayload {
	return serviceAccountUpdatePayload{Name: in.Name}
}

// Scopes renders even when nil, as JSON null: CreateServiceAccountTokenBody declares it
// required and nullable, and null is how a token inherits its service account's scopes.
type tokenPayload struct {
	Name      string     `json:"name"`
	Slug      *string    `json:"slug,omitempty"`
	Scopes    *[]string  `json:"scopes"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

func (in TokenInput) payload() tokenPayload {
	return tokenPayload{
		Name:      in.Name,
		Slug:      in.Slug,
		Scopes:    cloneStringSlice(in.Scopes),
		ExpiresAt: in.ExpiresAt,
	}
}

type ensureOrganizationPayload struct {
	ExternalID string  `json:"externalId"`
	Name       *string `json:"name,omitempty"`
}

func (in EnsureOrganizationInput) payload() ensureOrganizationPayload {
	return ensureOrganizationPayload(in)
}

type mintTokenPayload struct {
	Name   string    `json:"name"`
	Scopes *[]string `json:"scopes,omitempty"`
	TTL    *string   `json:"ttl,omitempty"`
}

func (in MintTokenInput) payload() mintTokenPayload {
	payload := mintTokenPayload{
		Name:   in.Name,
		Scopes: cloneStringSlice(in.Scopes),
	}

	if in.TTL != nil {
		// The API parses this with time.ParseDuration, the same grammar String() emits.
		formatted := in.TTL.String()
		payload.TTL = &formatted
	}

	return payload
}

type registerConnectorPayload struct {
	Name            string         `json:"name"`
	Version         string         `json:"version"`
	SettingsSchema  map[string]any `json:"settings_schema"`
	EntitlementSlug *string        `json:"entitlement_slug,omitempty"`
}

func (in RegisterConnectorInput) payload() registerConnectorPayload {
	payload := registerConnectorPayload(in)

	// The body declares settings_schema required, so an unset schema has to render as
	// an empty object rather than null: a connector with no settings is a real thing,
	// and null is a 422.
	if payload.SettingsSchema == nil {
		payload.SettingsSchema = map[string]any{}
	}

	return payload
}
