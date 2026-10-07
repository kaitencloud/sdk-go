# Kaiten Go SDK

[![CI](https://github.com/kaitencloud/sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kaitencloud/sdk-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kaitencloud/sdk-go.svg)](https://pkg.go.dev/github.com/kaitencloud/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaitencloud/sdk-go)](https://goreportcard.com/report/github.com/kaitencloud/sdk-go)
[![Release](https://img.shields.io/github/v/tag/kaitencloud/sdk-go?label=release)](https://github.com/kaitencloud/sdk-go/tags)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

Typed Go clients for the [Kaiten](https://kaiten.sh) API: customers, instances,
licenses, entitlements, usage metering, feature flags, and the Platform API that
provisions organizations and their credentials.

- [Installation](#installation)
- [Quick start](#quick-start)
- [Two APIs, two clients](#two-apis-two-clients)
- [Configuration](#configuration)
- [Core API](#core-api)
- [Platform API](#platform-api)
- [Pagination](#pagination)
- [Error handling](#error-handling)
- [Low-level access](#low-level-access)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Installation

Requires Go 1.25 or newer.

```shell
go get github.com/kaitencloud/sdk-go
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sdk "github.com/kaitencloud/sdk-go"
)

func main() {
	client, err := sdk.NewClient(
		"https://kaiten.example.com/api",
		sdk.WithBearerToken(os.Getenv("KAITEN_TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	instances, err := client.Instances.List(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, instance := range instances {
		fmt.Println(instance.Name)
	}
}
```

The base URL includes the API's path prefix. Every Kaiten deployment serves its own
OpenAPI document at `/api/docs`; the [Kaiten documentation](https://docs.kaiten.sh)
describes the resources and fields.

## Two APIs, two clients

Kaiten exposes two contracts on two listeners, and this package mirrors that split.

|                | `sdk.NewClient`                                                         | `sdk.NewPlatformClient`                                                 |
| -------------- | ----------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| Contract       | Core API                                                                | Platform API                                                            |
| Default port   | 3000                                                                    | 3001                                                                    |
| Credential     | organization-scoped token (`ksh_...`)                                   | platform credential (`ksm_...`)                                         |
| Resources      | customers, instances, licenses, entitlements, feature flags, components | organizations, users, organization tokens, connector registry           |

Each listener refuses the other's credential, so a process that needs both builds
both, against different base URLs. `Option` values are shared between the two
constructors.

## Configuration

```go
client, err := sdk.NewClient("https://kaiten.example.com/api",
	// Authorization: Bearer <token> on every request.
	sdk.WithBearerToken(token),

	// Your own transport, timeouts and retries. Replaces the default client
	// and its 30s DefaultTimeout.
	sdk.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),

	// Edit every outgoing request: tracing headers, per-call credentials, ...
	sdk.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Request-ID", requestIDFrom(ctx))
		return nil
	}),
)
```

| Option                     | Purpose                                                          |
| -------------------------- | ---------------------------------------------------------------- |
| `WithBearerToken(token)`   | Set the bearer token sent with every request.                    |
| `WithHTTPClient(client)`   | Use a custom `*http.Client`; replaces `DefaultTimeout` (30s).    |
| `WithRequestEditorFn(fn)`  | Run a function on each request before it is sent. Repeatable.    |
| `WithBaseURL(url)`         | Override the base URL given to the constructor.                  |

## Core API

`Client` groups the Core API by resource: `Customers`, `Instances`, `Licenses`,
`LicenseFamilies`, `Entitlements`, `EntitlementGroups`, `FeatureFlags`,
`MetadataFields`, `Components`, `DeploymentZones`, `Releases` and `ServiceAccounts`.
Every namespace exposes the CRUD operations its resource supports, plus the
operations below.

### Customers and instances

```go
customer, err := client.Customers.Create(ctx, sdk.CustomerInput{
	Name: "Acme Corp",
})
if err != nil {
	return err
}

instance, err := client.Instances.Create(ctx, sdk.InstanceInput{
	Name:             "acme-production",
	CustomerID:       *customer.Id,
	LicenseID:        licenseID,
	StartLicenseDate: time.Now(),
	EndLicenseDate:   time.Now().AddDate(1, 0, 0),
	Metadata:         map[string]any{"region": "eu-west-1"},
})
if err != nil {
	return err
}

// Instances are addressed by slug from here on.
err = client.Instances.UpdateStatus(ctx, *instance.Slug, sdk.InstanceStatusHealthy)
```

`InstanceInput.Slug` is generated when nil on create and renames the instance on
update; a slug another instance holds is a 409.

### Licenses and entitlements

A license is one version of a family; `LicenseFamilies` reads the family with its
current version, `Licenses` writes individual versions.

```go
// Grant 25 seats to the license, then make it available to instances.
seats, err := sdk.NumberLicenseValue(25)
if err != nil {
	return err
}
if err := client.Licenses.AssociateEntitlement(ctx, "team-v1", "seats", seats); err != nil {
	return err
}

license, err := client.Licenses.Publish(ctx, "team-v1")
if err != nil {
	return err
}
fmt.Println(*license.LifecycleState) // PUBLISHED
```

A license is created as `LicenseLifecycleDraft` or `LicenseLifecyclePublished` and
moves only through `Licenses.Publish`, `Licenses.Archive` and `Licenses.Unarchive`.
`sdk.BooleanLicenseValue` and `sdk.ConfigLicenseValue` build the other two value
kinds.

### Usage metering

```go
// Report a delta against a metered entitlement.
err := client.Instances.ReportUsage(ctx, "acme-production", "api-calls", 1)
if errors.Is(err, sdk.ErrThresholdExceeded) {
	// The instance has reached its limit for this entitlement.
}

// Or set the absolute value and get the resulting metric back.
usage, err := client.Instances.ReportEntitlementUsageMetric(ctx, "acme-production", "seats",
	sdk.UsageReportInput{Value: 18, Behavior: sdk.Set},
)
```

`Instances.ListEntitlementUsageMetrics` and `EntitlementGroups.GetUsage` read the
current position of an instance against its limits.

### Feature flags

```go
flags, err := client.FeatureFlags.List(ctx)
for _, flag := range flags {
	fmt.Println(*flag.Slug, flag.Enabled, flag.Type)
}
```

`sdk.BasicDefaultVariant` and `sdk.BasicVariantTargeting` build the default variant
and the targeting rules a `FeatureFlag` carries on create and update.

## Platform API

`PlatformClient` reaches the Platform API with a platform credential and exposes
`Organizations`, `Users`, `Tokens` and `Connectors`. `PlatformClient.Me` returns the
credential the client is using.

```go
platform, err := sdk.NewPlatformClient(
	"https://kaiten.example.com:3001/api",
	sdk.WithBearerToken(os.Getenv("KAITEN_PLATFORM_TOKEN")),
)
if err != nil {
	log.Fatal(err)
}

// Idempotent: the organization id derives from the external id, so every run
// converges on the same organization.
org, err := platform.Organizations.Ensure(ctx, sdk.EnsureOrganizationInput{
	ExternalID: "org_2abcDEF",
	Name:       &name,
})
if err != nil {
	log.Fatal(err)
}

// Mint an organization-scoped token for a machine consumer. The plaintext is
// returned once; the API stores a hash.
ttl := 90 * 24 * time.Hour
minted, err := platform.Tokens.Mint(ctx, org.Id, sdk.MintTokenInput{
	Name:   "zone-eu-west-1",
	Scopes: []string{"write:instances", "read:feature_flags"},
	TTL:    &ttl,
})
if err != nil {
	log.Fatal(err)
}
fmt.Println(*minted.Token)

// Tokens are revoked by slug, not by name.
err = platform.Tokens.Revoke(ctx, org.Id, *minted.Slug)
```

A token minted over the Platform API is bound to the credential that minted it, and
revoking that credential revokes every token it minted. Rotate as mint-new, adopt,
revoke-old.

### Registering a connector

A connector hosted outside the Kaiten API process declares itself deployment-wide.
Registration is an upsert on the name, so a host may call it on every start.

```go
connector, err := platform.Connectors.Register(ctx, sdk.RegisterConnectorInput{
	Name:    "kaiten.integration.crm.example",
	Version: "1.0.0",
	SettingsSchema: map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []any{"apiKey"},
		"properties": map[string]any{
			// writeOnly marks a secret: accepted, never read back.
			"apiKey": map[string]any{"type": "string", "writeOnly": true},
		},
	},
	// The BOOLEAN entitlement an organization's license must grant to
	// activate this connector. Nil leaves it ungated.
	EntitlementSlug: &entitlementSlug,
})
```

## Pagination

Every list method returns the whole collection: the SDK walks the API's cursor pages
itself, and a caller never sees a cursor.

The one exception is `Instances.ListAuditTrails`, whose collection grows without
bound. `AuditTrailsOptions.Limit` caps how many entries are fetched, and
`EventName`, `After` and `Before` filter them.

```go
limit := int32(100)
trails, err := client.Instances.ListAuditTrails(ctx, "acme-production", &sdk.AuditTrailsOptions{
	EventName: "instance.updated",
	After:     &since,
	Limit:     &limit,
})
```

`InstancesListOptions.IncludeDeleted` and `AuditTrailsOptions.Offset` are still sent
but ignored by current Kaiten versions; their doc comments say why they remain.

## Error handling

Every non-2xx response is returned as an `*sdk.Error`.

```go
instance, err := client.Instances.Get(ctx, "acme-production")
if err != nil {
	var apiErr *sdk.Error
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.Code == "Instance.NotFound":
			// Branch on Code: it is stable and distinguishes, for example,
			// the two different 404s a family endpoint can answer.
		case apiErr.StatusCode == http.StatusForbidden:
			// Wrong credential class or missing scope.
		default:
			// apiErr.Problem carries the RFC 9457 body; apiErr.ErrorID
			// correlates with the server-side log entry.
			log.Printf("%v", apiErr)
		}
		return err
	}
	return err // transport error
}
```

| Field        | Contents                                                              |
| ------------ | --------------------------------------------------------------------- |
| `StatusCode` | HTTP status code.                                                     |
| `Code`       | Stable, machine-readable code such as `License.NotFound`, when given. |
| `Problem`    | The decoded problem-details body (`title`, `detail`, `errors`, ...).  |
| `ErrorID`    | Correlation id for the server-side log entry, when given.             |
| `Body`       | The raw response body.                                                |

`Instances.ReportUsage` additionally returns `sdk.ErrThresholdExceeded` when a report
would cross the entitlement's limit.

## Low-level access

`Client.Raw()` and `PlatformClient.Raw()` return the generated OpenAPI clients for
operations the typed namespaces do not cover. They are the oapi-codegen output of
Kaiten's two OpenAPI documents; their request and response types are the ones the
typed layer aliases in `types.go`.

## Development

Install [Task](https://taskfile.dev) and Go 1.25+. Tools run through
`go run <tool>@<version>` against the versions pinned in `Taskfile.yml`, so local
runs match CI.

```shell
task test           # go test -race -shuffle=on -cover ./...
task lint           # golangci-lint, the version CI runs
task fmt            # gofumpt + gci
task vuln           # govulncheck over reachable code
task release:check  # gorelease: exported API diff against the last tag
```

`internal/gen` and `internal/genplatform` are generated from Kaiten's Core and
Platform OpenAPI documents and are never edited by hand:

```shell
task generate          OPENAPI_SPEC=path/to/openapi.yaml
task generate:platform PLATFORM_OPENAPI_SPEC=path/to/platform-openapi.yaml
```

Request bodies are rendered by hand in `inputs.go`, because the API rejects unknown
keys. `TestRequestBodiesMatchTheirSpecBodies` pins each input type's key set against
the spec, so adding a field to an input means adding it to that test's expectations
too.

## Contributing

Issues and pull requests are welcome. Before opening one, run `task fmt`, `task lint`
and `task test`. A change to any exported identifier is a release concern for every
consumer of this module: run `task release:check` and say in the pull request which
semver bump it forces.

Every commit must be signed off (`git commit -s`) under the Developer Certificate of
Origin: [CONTRIBUTING.md](CONTRIBUTING.md) explains it, and how to fix a commit that
missed it.

## License

Kaiten SDK for Go is open source and licensed under the
[Apache License, Version 2.0](./LICENSE).

By contributing, you agree to certify your contribution under the
[Developer Certificate of Origin 1.1](./DCO.md).

The Kaiten name and logos are not licensed under Apache-2.0. See the
[Kaiten trademark policy](https://github.com/kaitencloud/kaiten/blob/main/TRADEMARKS.md).
