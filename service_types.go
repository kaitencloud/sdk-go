package sdk

// Components provides access to the components API.
type Components struct{ client *Client }

// Customers provides access to the customers API.
type Customers struct{ client *Client }

// DeploymentZones provides access to the deployment zones API.
type DeploymentZones struct{ client *Client }

// EntitlementGroups provides access to the entitlement groups API.
type EntitlementGroups struct{ client *Client }

// Entitlements provides access to the entitlements API.
type Entitlements struct{ client *Client }

// FeatureFlags provides access to the feature flags API.
type FeatureFlags struct{ client *Client }

// Instances provides access to the instances API.
type Instances struct{ client *Client }

// LicenseFamilies provides access to the license families API.
type LicenseFamilies struct{ client *Client }

// Licenses provides access to the licenses API.
type Licenses struct{ client *Client }

// MetadataFields provides access to the metadata fields API.
type MetadataFields struct{ client *Client }

// Releases provides access to the releases API.
type Releases struct{ client *Client }

// ServiceAccounts provides access to the service accounts API.
type ServiceAccounts struct{ client *Client }

// The namespaces below belong to PlatformClient, not Client: they are served by the
// Platform API listener and reject an organization-scoped token. Organizations and users
// live there because they are the objects a platform credential administers -- a tenant's
// own token has no business enumerating or deleting them.

// PlatformConnectors provides access to the platform connector registry: the manifests
// of the connectors a deployment can offer its organizations. Registration is the
// deployment-wide half -- "this connector exists here" -- and an organization activating
// one is the other half, on the Core API and gated by its licence.
type PlatformConnectors struct{ client *PlatformClient }

// PlatformOrganizations provides access to the platform organizations API.
type PlatformOrganizations struct{ client *PlatformClient }

// PlatformTokens provides access to the organization token API: minting and revoking the
// organization-scoped credentials (`ksh_...`) that Core API callers authenticate with.
type PlatformTokens struct{ client *PlatformClient }

// PlatformUsers provides access to the platform users API. Minimal by design: users are
// JIT-provisioned from JWT claims, not created or updated through this API -- Delete is
// the only operation, for privileged (delete:users) callers.
type PlatformUsers struct{ client *PlatformClient }
