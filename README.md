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

The base URL includes the API's path prefix: `https://app.kaiten.sh/api` on Kaiten
Cloud, `http://localhost:6060/api` on a local Community Edition stack. The API serves
its own OpenAPI document at `/api/openapi.json` and a reference rendered from it at
`/api/docs`; the [Kaiten documentation](https://docs.kaiten.sh) describes the resources
and fields.

## Two APIs, two clients

Kaiten exposes two contracts on two listeners, and this package mirrors that split.

|                | `sdk.NewClient`                                                         | `sdk.NewPlatformClient`                                                 |
| -------------- | ----------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| Contract       | Core API                                                                | Platform API                                                            |
| Local stack    | `http://localhost:6060/api`                                             | `http://localhost:6001/api`                                             |
| Credential     | organization-scoped token (`ksh_...`)                                   | platform credential (`ksm_...`)                                         |
| Resources      | customers, instances, licenses, entitlements, feature flags, components | organizations, users, organization tokens, connector registry           |

The Platform API is an internal listener, for the operator of a self-hosted Kaiten: it
is on port 3001 unless configured otherwise, and the local Community Edition stack
publishes it on 6001.

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
	Name: "Sakura Tokyo",
})
if err != nil {
	return err
}

slug := "sakura-tokyo-prod"
instance, err := client.Instances.Create(ctx, sdk.InstanceInput{
	Name:             "Sakura Tokyo Production",
	Slug:             &slug,
	CustomerID:       *customer.Id,
	LicenseID:        licenseID,
	StartLicenseDate: time.Now(),
	EndLicenseDate:   time.Now().AddDate(1, 0, 0),
	Metadata:         map[string]any{"environment": "production"},
})
if err != nil {
	return err
}

// Instances are addressed by slug from here on.
err = client.Instances.UpdateStatus(ctx, *instance.Slug, sdk.InstanceStatusHealthy)
```

`InstanceInput.Slug` is generated when nil on create and renames the instance on
update; a slug another instance holds is a 409.

Customers, deployment zones, entitlements, entitlement groups, licenses and service
accounts are never renamed. Their `Update` sends `Slug` for the API to check: the
current slug is accepted, and any other is a 422 whose `*sdk.Error` carries the code
`<Operation>.SlugNotRenameable` (for example `UpdateEntitlement.SlugNotRenameable`).

### Licenses and entitlements

A license is one version of a family; `LicenseFamilies` reads the family with its
current version, `Licenses` writes individual versions.

```go
// Grant 10 delivery drivers to the draft license, then make it available to instances.
drivers, err := sdk.NumberLicenseValue(10)
if err != nil {
	return err
}
if err := client.Licenses.AssociateEntitlement(ctx, "premium-v3", "delivery-drivers", drivers); err != nil {
	return err
}

license, err := client.Licenses.Publish(ctx, "premium-v3")
if err != nil {
	return err
}
fmt.Println(*license.LifecycleState) // PUBLISHED
```

A license is created as `LicenseLifecycleDraft` or `LicenseLifecyclePublished` and
moves only through `Licenses.Publish`, `Licenses.Archive` and `Licenses.Unarchive`.
`sdk.BooleanLicenseValue` and `sdk.ConfigLicenseValue` build the other two value
kinds.

An entitlement says how usage is measured, and a grant how much of it a license allows
and how strictly. A periodic entitlement counts per window, and a grant may tolerate
usage past its value:

```go
// Orders count per calendar month, with a warning at 80% of what each license grants.
slug, number, sum := "monthly-orders", sdk.EntitlementType("NUMBER"), sdk.EntitlementAggregationMethod("SUM")
month, calendar, warning := sdk.EntitlementResetPeriodMonth, sdk.EntitlementResetAnchorCalendar, int32(80)
if _, err := client.Entitlements.Create(ctx, sdk.EntitlementInput{
	Name:                    "Monthly Orders",
	Slug:                    &slug,
	Type:                    &number,
	AggregationMethod:       &sum,
	ResetPeriod:             &month,
	ResetAnchor:             &calendar,
	WarningThresholdPercent: &warning,
}); err != nil {
	return err
}

// 1000 orders a month, refused past 1200: a soft limit with a 20% allowance.
orders, err := sdk.NumberLicenseValue(1000)
if err != nil {
	return err
}
err = client.Licenses.AssociateEntitlement(ctx, "premium-v2", slug, orders, sdk.WithOveragePercent(20))
```

`Entitlements.Update` and `Licenses.UpdateEntitlement` replace rather than patch. A
field left out is reset -- a warning threshold to 0, a grant's allowance to a hard
limit -- and a periodic entitlement's reset period and anchor have to be sent back as
they are, or the update is refused. Read the resource first and carry back what should
stay.

### Usage metering

```go
// Report a delta against a metered entitlement.
err := client.Instances.ReportUsage(ctx, "sakura-tokyo-prod", "monthly-orders", 1)
if errors.Is(err, sdk.ErrThresholdExceeded) {
	// Refused: usage would pass what the license grants, allowance included.
}

// Or set the absolute value and get the resulting metric back.
usage, err := client.Instances.ReportEntitlementUsageMetric(ctx, "sakura-tokyo-prod", "delivery-drivers",
	sdk.UsageReportInput{Value: 8, Behavior: sdk.Set},
)
```

`Instances.ListEntitlementUsageMetrics` and `EntitlementGroups.GetUsage` read the
current position of an instance against its limits.

### AI credits

An AI credit pool is metered in two phases: `AI.Reserve` sets credits aside on an
estimate before the provider is called, and `AI.Commit` charges what the call actually
consumed. The server prices both from the license version's current conversion rules (a
reservation does not lock rates); the SDK carries no pricing. `AI.Run` brackets one call between the two:

```go
var answer string
_, err := client.AI.Run(ctx, "sakura-tokyo-prod", "chat-credits", sdk.AIReserveInput{
	Model:       "gpt-4o",
	RequestType: sdk.AIRequestTypeCompletion,
	Estimate:    sdk.AIUsage{InputTokens: 1000, OutputTokens: 800},
}, func(ctx context.Context) (sdk.AICall, error) {
	// Call the provider here. Report sdk.AIRequestStatusNotSent when it was never
	// reached, and the usage known so far when it failed.
	answer = "..."
	return sdk.AICall{Usage: sdk.AIUsage{InputTokens: 1000, OutputTokens: 412}}, nil
})
switch {
case errors.Is(err, sdk.ErrInsufficientCredits):
	// Refused before the provider was called: the estimate does not fit the pool.
case errors.Is(err, sdk.ErrAICommitFailed):
	// The call succeeded and answer is good; only the accounting failed.
}
```

A refused reserve is one of `sdk.ErrInsufficientCredits`, `sdk.ErrModelNotAuthorized`,
`sdk.ErrUnknownModel`, `sdk.ErrTransactionIDReused` or `sdk.ErrAIUsageUnsupported` (the
pool is not an AI credit pool, or the server predates AI reports). The transaction id
defaults to a UUIDv7. Both phases retry a network error or a 500, 502, 503 or 504 up to
three times under the identical phase/body and transaction id; a 4xx is never retried.
AI helpers remove even inherited `Idempotency-Key` headers: that header is for numeric
adjustments only and is rejected for `ai_usage`. Unknown problem codes remain `*sdk.Error`, not a
guessed unsupported-capability error.
`ErrAICommitFailed` means settlement could not be confirmed, not that nothing was
charged: a lost response may hide a successful commit. Reconcile inference history
or retry the identical commit body under the same transaction id.

The provider SDKs count tokens their own way. The `ai/openai` and `ai/anthropic`
modules convert their usage structs into an `sdk.AIUsage`, and are separate modules so
this one never depends on a provider SDK:

```shell
go get github.com/kaitencloud/sdk-go/ai/openai
go get github.com/kaitencloud/sdk-go/ai/anthropic
```

### Generated AI management and history

The generated operations on `client.Raw()` cover model list/get/create/full-replace
update/delete, version-scoped conversion-rule list/get/create/update/archive, and AI
reservation/inference reads. Their input and query types are exported by this package:

```go
// Catalogue vendor costs are informational USD strings, not credit prices.
models, err := client.Raw().ListModelsWithResponse(ctx, &sdk.ListModelsParams{})

// This slug identifies a LICENSE VERSION, not its family. Rates are exact strings;
// nil is unpriced, while a pointer to "0" is explicitly zero-priced. Missing cache
// rates fall back to the input-token rate. Update replaces the rate vector.
rate := "0.0000000001"
rule, err := client.Raw().CreateConversionRuleWithResponse(ctx, "pro-v3", sdk.CreateConversionRule{
    EntitlementSlug: "chat-credits",
    Rates: sdk.ConversionRates{PerInputToken: &rate},
})

history, err := client.Raw().ListAiInferencesWithResponse(ctx, "sakura-tokyo-prod", "chat-credits", &sdk.ListAIInferencesParams{})
reservations, err := client.Raw().ListAiReservationsWithResponse(ctx, "sakura-tokyo-prod", "chat-credits")
```

Raw responses expose status codes and problem details; callers must check them as well
as transport errors. Inference credits, rates, vendor costs and generated units stay
strings without float conversion. `AIUsageQuantities` uses pointers: absent dimensions
are unknown, not known zero. The convenience `AIUsage` omits zero fields; use
`AIUsageValue` with `ReportEntitlementUsageValue.FromAiUsageValue` and the raw report
operation when an explicit zero matters. Unsupported provider dimensions must remain
absent; do not invent quantities. OpenAI input tokens include cache reads, so the
normalizer subtracts them; Anthropic's input/cache dimensions are already disjoint.

An entitlement of type `NUMBER_AI_CREDIT` uses `SUM` aggregation and numeric allowances
and read counters (`type: "number"`). `ai_usage` is only a calculation input: read the
counter through `Instances.GetEntitlementUsageMetric`, not from reservation/history totals.
Numeric adjustments/refunds use the ordinary number report with append behavior.

Model and conversion-rule webhook payload aliases are also exported, together with
`OnInstanceAIModelUsageUnauthorizedJSONBody` and
`OnInstanceEntitlementUsageReportAcceptedJSONBody`. Unauthorized commits are accounted
with their authorization/warning outcome; accepted usage hooks still carry numeric
credit values rather than the write-only `ai_usage` payload.

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

`Instances.ListAuditTrails` is the one list that takes a ceiling, because its
collection grows without bound. `AuditTrailsOptions.Limit` caps how many entries are
fetched -- without one the whole trail is walked, like any other list -- and
`EventName`, `After` and `Before` filter them.

```go
limit := int32(100)
trails, err := client.Instances.ListAuditTrails(ctx, "sakura-tokyo-prod", &sdk.AuditTrailsOptions{
	EventName: "INSTANCE_UPDATED",
	After:     &since,
	Limit:     &limit,
})
```

Two options remain only so that callers still compile, and their doc comments say
why: `InstancesListOptions.IncludeDeleted` is still sent, and ignored by current Kaiten
versions; `AuditTrailsOptions.Offset` is deprecated and never sent, since the API
paginates by cursor only.

## Error handling

Every non-2xx response is returned as an `*sdk.Error`.

```go
instance, err := client.Instances.Get(ctx, "sakura-tokyo-prod")
if err != nil {
	var apiErr *sdk.Error
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.Code == "GetInstance.NotFound":
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
| `Code`       | Stable, machine-readable code such as `GetLicense.NotFound`, if any.  |
| `Problem`    | The decoded problem-details body (`title`, `detail`, `errors`, ...).  |
| `ErrorID`    | Correlation id for the server-side log entry, when given.             |
| `Body`       | The raw response body.                                                |

`Instances.ReportUsage` additionally returns `sdk.ErrThresholdExceeded` when a report
would cross the entitlement's limit. The `AI` reports return the typed errors listed
under [AI credits](#ai-credits), each wrapping the `*sdk.Error` as well.

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
