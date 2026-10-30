# Kaiten Go SDK

[![CI](https://github.com/kaitencloud/sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kaitencloud/sdk-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kaitencloud/sdk-go.svg)](https://pkg.go.dev/github.com/kaitencloud/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaitencloud/sdk-go)](https://goreportcard.com/report/github.com/kaitencloud/sdk-go)
[![Release](https://img.shields.io/github/v/tag/kaitencloud/sdk-go?label=release)](https://github.com/kaitencloud/sdk-go/tags)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

Typed Go clients for the [Kaiten](https://kaiten.sh) API: customers, instances,
licenses, entitlements, usage metering, feature-flag management, and the Platform API
that provisions organizations and their credentials. **Flag evaluation requires the
OpenFeature Go SDK with its official OFREP provider**, not this management SDK.

- [Installation](#installation)
- [Quick start](#quick-start)
- [Two APIs, two clients](#two-apis-two-clients)
- [Configuration](#configuration)
- [Core API](#core-api)
- [Flag evaluation with OpenFeature](#flag-evaluation-with-openfeature)
- [Breaking migration](#breaking-migration)
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

A version also carries its commercial terms. `LicenseInput` writes them as
`PricingType` (`LicensePricingTypeFREE`, `LicensePricingTypePAID` or
`LicensePricingTypeCUSTOM`), `TrialPeriodDays`, `RequiresPaymentMethod` and
`SelfServeCtaURL`, and `License` reads them back under the same names, the last one
spelled `SelfServeCtaUrl`. A term left nil takes the server's default on create --
CUSTOM, no trial, no payment method, no URL -- and keeps its value on update, where
`TrialPeriodDays` 0 removes the trial and an empty `SelfServeCtaURL` the URL.

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

#### Safe retries with a transaction ID

A report without a key is sent once and never retried: if its request fails, it
may or may not have been counted, and sending it again could count it twice. Give
it a `TransactionID`, sent as the `Idempotency-Key` header, and Kaiten applies it
at most once (per instance, entitlement and key, within 35 days by default). The
SDK then retries it on network errors and on 500, 502, 503 and 504, up to three
attempts, all under the same key -- once the server has shown it honours keys, by
answering a report with its `transactionId`. Until then a failed keyed report is
returned to you, like one without a key.

```go
result, err := client.Instances.ReportEntitlementUsage(ctx, "acme-production", "tokens",
	sdk.UsageReportInput{
		Value:         1200,
		TransactionID: "llm-call:9f2c:tokens", // a UUID, or the business event ID plus the meter
		Metadata:      map[string]any{"model": "large"},
	},
)
switch {
case errors.Is(err, sdk.ErrTransactionIDReused):
	// The key was already used for a different report: a bug in how keys are
	// made. Not retryable. Send a correction as a new report under a new key.
case err != nil:
	return err
case result.Replayed:
	// Kaiten had already counted this report; result.Usage is its original
	// answer, possibly for a window that has since closed.
case result.MetadataDropped:
	// Counted, but the metadata was above 4 KiB and was not stored.
}
```

That proof is kept per client. A Kaiten older than keys ignores the header, so a
report resent to it would count twice: its answers never carry a
`transactionId`, and the client never retries against it. The SDK never makes a
key up, and never resends a report without its key.

#### Usage history

Every accepted report is kept in the usage history for the organization's
retention: the counter before and after it, and the limit it was gated on.

```go
from := time.Now().AddDate(0, 0, -7)
reports, err := client.Instances.ListUsageReports(ctx, "acme-production", "tokens",
	&sdk.UsageReportsOptions{From: &from},
)

// Or stream it, as CSV or NDJSON. Close the body.
export, err := client.Instances.ExportUsageReports(ctx, "acme-production", "tokens",
	&sdk.UsageExportOptions{From: &from, Format: sdk.UsageExportCSV},
)
if err != nil {
	return err
}
defer export.Body.Close()
_, err = io.Copy(file, export.Body) // export.Filename names the range
```

`Instances.ExportOrganizationUsageReports` exports the whole organization, 31 days
at a time, optionally for one instance or entitlement; `InstanceID` and
`EntitlementID` reach deleted ones. A `From` before the start of the
organization's history answers a 422 whose `Code` ends in `OutsideRetention`.
The client's timeout covers reading an export's body: give a large export a client
with a longer one through `WithHTTPClient`.

### Feature flags

```go
flags, err := client.FeatureFlags.List(ctx)
for _, flag := range flags {
	fmt.Println(*flag.Slug, flag.Enabled, flag.Type)
}
```

`sdk.BasicDefaultVariant` and `sdk.BasicVariantTargeting` build the default variant
and the targeting rules a `FeatureFlag` carries on create and update.

### Billing

`client.Billing` holds the billing settings and provider (`Settings`,
`UpdateSettings`, `Capabilities`, `Health`, `Sync`, `ClosePeriods`) and a
namespace per resource:

| Namespace           | Covers                                                                  |
| ------------------- | ----------------------------------------------------------------------- |
| `Subscriptions`     | an instance's subscription: subscribe, change, cancel, reactivate, scheduled plan change, upcoming invoice, its invoices |
| `Invoices`          | list, get, export, void, recompose, mark paid, write off, release a hold, retry a push, sync, a line's usage reports |
| `Handoff`           | the accounting handoff: claim a leased batch, acknowledge each invoice  |
| `LicensePrices`     | a licence version's prices and the invoice preview                      |
| `Addons`            | the add-on catalogue: versions, grants, prices, compatibility, families |
| `InstanceAddons`    | attach, set the quantity of, detach an add-on                           |
| `Vouchers`, `InstanceVouchers` | vouchers, their redemptions, redeem and revoke on an instance |
| `PublishableKeys`, `CustomerSessions` | the `pk_` keys and `kst_` sessions a customer's page uses |
| `Customers`         | a customer's billing, payment method and provider portal                |
| `Public`, `Session` | what a `pk_` key and a `kst_` session reach: build the client with that token |

```go
claim, err := client.Billing.Handoff.Claim(ctx, sdk.HandoffClaimSize{})
for _, invoice := range claim.Invoices {
	reference := book(invoice) // your accounting system
	_, err = client.Billing.Handoff.Ack(ctx, invoice.Id, sdk.HandoffBooking{
		LeaseId:           &claim.LeaseId,
		ExternalReference: &reference,
	})
}
```

The list methods walk every page. Each paginated one (`Invoices.List`,
`Handoff.List`, `Vouchers.ListRedemptions`, ...) also has a `...Page` variant that
returns one page and its cursor.

## Flag evaluation with OpenFeature

The Kaiten SDK only manages flag definitions. Evaluate flags with the
[OpenFeature Go SDK](https://github.com/open-feature/go-sdk) and the
[official OFREP provider](https://github.com/open-feature/go-sdk-contrib/tree/main/providers/ofrep).
There is no direct evaluation API, including through `Client.Raw()`.

```shell
go get github.com/open-feature/go-sdk@v1.17.0
go get github.com/open-feature/go-sdk-contrib/providers/ofrep@v0.1.7
```

These versions support the Go 1.25 floor of this module. In your server process,
register the provider once at startup and reuse the OpenFeature client:

```go
// Imports: fmt, net/http, time, github.com/open-feature/go-sdk/openfeature,
// and github.com/open-feature/go-sdk-contrib/providers/ofrep.
httpClient := &http.Client{Timeout: 5 * time.Second}
provider := ofrep.NewProvider(
	"https://app.kaiten.sh/api",
	ofrep.WithBearerToken(token), // Organization-scoped PAT: read:feature_flags.
	ofrep.WithClient(httpClient),
)
if err := openfeature.SetNamedProviderAndWait("kaiten", provider); err != nil {
	return err
}
// At process shutdown: openfeature.Shutdown(). The application still owns httpClient.
flags := openfeature.NewClient("kaiten")

// Per request: identifiers come from authenticated, authorized server state.
evaluationContext := openfeature.NewEvaluationContext(customerSlug, map[string]any{
	"kaiten": map[string]any{"instanceSlug": instanceSlug},
})
details, err := flags.BooleanValueDetails(ctx, "express-delivery", false, evaluationContext)
if err != nil || details.Reason == openfeature.ErrorReason {
	// Fail closed, even if an ERROR reason arrives without a Go error. Record
	// only a sanitized failure, never the token or private context.
	return fmt.Errorf("express-delivery evaluation failed")
}
enabled := details.Value // Use only after checking both error and reason.
```

The base URL is the **Core gateway including its API prefix**, not the Platform
listener and not the `/ofrep` path: Cloud `https://app.kaiten.sh/api`, local
Community Edition `http://localhost:6060/api`. The official provider appends
`/ofrep/v1/evaluate/flags/{key}`. Keep the PAT on the server; it needs
`read:feature_flags` in the target organization. SDK authentication options do
not automatically configure the separate OpenFeature provider.

Build a new evaluation context for each request; never put request-specific
identity on a global provider or client. A targeting key can be a customer slug
(for customer license facts) or a user/device identifier. Instance targeting
requires the explicit `kaiten.instanceSlug` (or `kaiten.instanceId`) context. The
backend must authorize these selectors against the authenticated tenant before
sending them; browser-supplied tenant or instance selectors are not trusted.
Kaiten resolves server-owned `__kaiten` facts within the PAT organization; callers
must not supply their own license or entitlement facts. OpenFeature gates are
not a replacement for backend authorization or metering enforcement.

[The runnable example](examples/openfeature/main.go) and its isolated contract
tests use the real official provider, without a custom evaluation HTTP wrapper.
It is a separate module so these dependencies do not ship in the management SDK.
`task example:check` tests failures, request-context isolation, auth and the API
prefix against a local protocol fixture, not a live Core deployment.

## Breaking migration

The next release removes the generated `EvaluateFlag` and `EvaluateFlagsBulk`
operation families exposed through `Client.Raw()`: typed/raw-body methods,
`WithResponse` variants, request builders, response parsers and operation-specific
request/response types. There is no compatibility alias or deprecated direct
evaluator. Replace these calls with OpenFeature evaluations configured as above.

`Client.Raw()` remains available for management operations, including feature-flag
CRUD; typed management, licensing and usage-reporting contracts are unchanged.
Core continues serving OFREP for official providers: this is a breaking **SDK
surface** change, not a removal of the server protocol. Previously published
versions and tags remain immutable, including the already published `v0.2.0`.
This break requires the next minor line (`v0.3.0`), not a patch; the module path
remains `github.com/kaitencloud/sdk-go`. For a stable `v1+` line, the same removal
would require a major version. Review the release check against the published
line when landing this scoped correction: do not replace newer billing, licensing
or usage-management APIs with an older SDK snapshot.

Go API-diff tools may not enumerate methods of an internal generated type returned
by a public method. Review the raw-client changes as breaking even if
`task release:check` does not list them. Consumers pinned to an earlier version
are not changed automatically; they need an explicit dependency upgrade and
migration if they used these raw operations.

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
| `RetryAfter` | How long to wait before retrying, when the API said (429, 503).       |

The usage report methods also return sentinels for `errors.Is`, alongside the
`*sdk.Error`:

| Sentinel                      | Meaning                                                                 |
| ----------------------------- | ----------------------------------------------------------------------- |
| `ErrThresholdExceeded`        | The report would cross the entitlement's limit.                         |
| `ErrTransactionIDReused`      | The key was already used for a different report. Do not retry it.       |
| `ErrInvalidTransactionID`     | The key is not 1 to 128 characters of `[A-Za-z0-9._:-]`; nothing was sent. |

## Low-level access

`Client.Raw()` and `PlatformClient.Raw()` return the generated OpenAPI clients for
management operations the typed namespaces do not cover. Core generation excludes
the two OFREP evaluation operations by exact operation ID; the server contract and
its protocol models are not removed. Their request and response types are the ones
the typed layer aliases in `types.go`.

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
task generate:check # fixture: no OFREP evaluator, management/overlay preserved
task example:check  # official OpenFeature + OFREP, isolated contract tests
```

`internal/gen` and `internal/genplatform` are generated from Kaiten's Core and
Platform OpenAPI documents and are never edited by hand:

```shell
task generate          OPENAPI_SPEC=path/to/openapi.yaml
task generate:platform PLATFORM_OPENAPI_SPEC=path/to/platform-openapi.yaml
```

### Core generation provenance

The published `v0.2.0` Core client was generated from Core commit
`131368cf8a6c73e51dfd846a41ecf1f752a2c13c`, `app/openapi.yaml`
([exact source](https://github.com/kaitencloud/kaiten/blob/131368cf8a6c73e51dfd846a41ecf1f752a2c13c/app/openapi.yaml)).
Its source SHA-256 is
`b0283bac70fe33abb46f6b4655896e896f26eb5fda47ef9e09107a5db06c0ee1`.
Using oapi-codegen v2.8.0, `skip-prune` and the existing invoice naming overlay
reproduces that published artifact byte-for-byte (SHA-256
`730becd1a5ff00a1fd939679c70119910ffc26d0cb54e50b678b79c9c91047b4`).

This correction regenerates from that same exact source, adding only the
`evaluateFlag` and `evaluateFlagsBulk` operation-ID exclusions. The resulting
artifact SHA-256 is
`a5cb1ba30f07aa97d79e6bd890ca28390dd78c659fdf1103a65cf50124bad8f2`:
709 deleted lines, no additions, 62 evaluation exports removed and all 9,312
other exports retained. Billing, commercial license terms, usage management and
webhook/protocol models are preserved. No historical AST fallback or Core
server-spec edit is needed on this base. `task generate:check` checks the artifact
and an overlay-aware minimal fixture; negative public API tests guard `Raw()`
and the management namespaces against reintroducing a direct evaluator.

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
