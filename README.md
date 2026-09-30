# sdk-go

Typed Go clients for the [Kaiten](https://kaiten.sh) API.

## Installation

```shell
go get github.com/kaitencloud/sdk-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/kaitencloud/sdk-go"
)

func main() {
	client, err := sdk.NewClient(
		"https://kaiten.example.com/api",
		sdk.WithBearerToken("your-api-token"),
	)
	if err != nil {
		log.Fatal(err)
	}

	instances, err := client.Instances.List(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, instance := range instances {
		if instance.Slug != nil {
			fmt.Println(*instance.Slug)
		}
	}
}
```

See the [Kaiten documentation](https://docs.kaiten.sh) for details on available
resources and fields. Every deployment also serves its own contract at
`/api/docs`.

## The Platform API

Kaiten publishes two contracts on two listeners, and this package mirrors that split
rather than hiding it:

| | `sdk.NewClient` | `sdk.NewPlatformClient` |
|---|---|---|
| Contract | Core API (`app/openapi.yaml`) | Platform API (`app/platform-openapi.yaml`) |
| Default port | 3000 | 3001 |
| Credential | organization-scoped token, `ksh_...` | platform credential, `ksm_...` |
| Covers | customers, instances, entitlements, licenses, feature flags | organizations, users, the organization tokens themselves, and the connector registry |

The two are not interchangeable in either direction: the Core listener refuses a platform
credential before it parses the request, and the Platform listener serves none of the Core
paths. A process that needs both builds both, against different base URLs. `Option` values
are shared -- one `WithBearerToken` or `WithHTTPClient` works with either constructor.

```go
platform, err := sdk.NewPlatformClient(
	"https://kaiten.example.com:3001/api",
	sdk.WithBearerToken(os.Getenv("KAITEN_PLATFORM_TOKEN")),
)
if err != nil {
	log.Fatal(err)
}

// Idempotent: the organization id is derived from the external id, so the same call
// converges on the same organization on every run.
org, err := platform.Organizations.Ensure(ctx, sdk.EnsureOrganizationInput{
	ExternalID: externalOrganizationID,
	Name:       &name,
})
if err != nil {
	log.Fatal(err)
}

// The plaintext is returned once and never again -- the API stores a hash.
minted, err := platform.Tokens.Mint(ctx, org.Id, sdk.MintTokenInput{
	Name:   "zone-kaiten",
	Scopes: []string{"write:instances", "read:feature_flags"},
})
```

Two things about minted tokens are easy to get wrong:

- **Revocation cascades.** A token minted over the Platform API is bound to the platform
  credential that minted it; revoking the parent revokes every child. Rotate as
  mint-new, adopt, revoke-old -- never revoke-then-mint.
- **`Slug`, not `Name`, is the handle.** `Tokens.Revoke` takes the slug from
  `MintedToken.Slug`.

### Registering a connector

A connector hosted outside the Kaiten API process declares itself over the Platform API:

```go
schema := map[string]any{
	"type":                 "object",
	"additionalProperties": false,
	"required":             []any{"apiKey"},
	"properties": map[string]any{
		// writeOnly marks a secret: the API accepts it and never reads it back out.
		"apiKey": map[string]any{"type": "string", "writeOnly": true},
	},
}

connector, err := platform.Connectors.Register(ctx, sdk.RegisterConnectorInput{
	Name:            "kaiten.integration.crm.example",
	Version:         "1.0.0",
	SettingsSchema:  schema,
	EntitlementSlug: &slug, // nil leaves the connector ungated
})
```

Two halves, and they are deliberately not the same call. `Register` is deployment-wide --
"this connector exists here, and this is what its settings look like" -- and is an upsert
on the name, so a process that hosts a connector can call it on every start. Whether a
given organization may *use* it is the other half: the organization activates the
connector for itself over the Core API, and `EntitlementSlug` names the BOOLEAN
entitlement its licence must grant for that to be allowed. Nil means ungated, which is
the open-source default.

Kaiten's own built-in connectors never come through here: they are compiled into the API
binary and register in-process at startup, where there is no wire to be on and no
credential to present.

## The generated client, and what is hand-written

`internal/gen` and `internal/genplatform` are oapi-codegen output from Kaiten's two
OpenAPI documents, and carry a `DO NOT EDIT` header: a defect in either is fixed in the
spec or by a generator bump, never by editing the file. `task generate` and
`task generate:platform` regenerate them; both need oapi-codegen v2.8.0 or newer, and
the Core config must keep `output-options: skip-prune: true`. Without it, schemas
reachable only from webhook definitions (`Deployment`, `FeatureFlagEvaluated`,
`SystemTokenIssuance`, ...) are pruned as unused while the webhook body types that
reference them are still emitted, and the package does not compile.

`types.go` aliases the generated models under this package's own names. Where a Go name
and a spec name disagree -- `EntitlementGroupRef` for `EntitlementGroupSummary`,
`ErrorModel` for `Problem` -- the Go name wins: it is a published identifier, and
renaming it would break every consumer for a wire name they never see.

Two things are deliberately **not** taken from the generator, and one thing the
generator no longer offers.

**Request bodies are rendered by hand.** Every Core request body is
`additionalProperties: false`, so a key the spec does not declare is a 422 rather than a
field the API ignores, and a required-but-nullable key the client omits is the same 422
from the other side. Neither shows up as a type error. So each input type renders its own
payload, and `TestRequestBodiesMatchTheirSpecBodies` pins the exact key set of every
write against the spec's own bodies, in both directions.

A create body and its update body routinely differ -- a slug is assigned once and then
immutable; a license's lifecycle state is chosen on create (draft or published) and then
moved only by `Licenses.Publish`, `Archive` and `Unarchive` -- so one input type renders
two payloads, `createPayload()` and `updatePayload()`. `instancePatchPayload` is
hand-written for the same reason: the generated `PatchInstanceBody` renders
`lifecycle_stage` where the spec names the key `lifecycleStage`.

**The inputs are narrower than the spec, on purpose.** `EntitlementInput` carries the
presentation fields (`icon`, `unitSingular`, `unitPlural`, `userFacing`, `displayOrder`)
but not `enforcementMode`, `resetPeriod`, `resetAnchor`, `warningThresholdPercent` or the
`saleUnit*` trio. Fields are added when a caller needs one, so that the wire behaviour of
each is pinned by a test rather than assumed.

**Registration is a platform operation.** `internal/gen` declares no
`RegisterConnector`: the Core document dropped that path, and
`PlatformConnectors.Register` is the one way in.

### Every list is paginated, and this package walks it

Core paginates all fourteen of its list endpoints: they answer
`{"items": [...], "nextCursor": ..., "hasMore": ...}`, not a bare array.

`pagination.go` decodes that envelope by hand and **walks the pages itself**, so the list
methods keep their signatures and their "returns all X" promise stays true. A caller never
sees a cursor. Two things are worth knowing:

- **`limit=200` is a hard ceiling, not a hint.** The endpoints declare `maximum: 200`, and
  huma rejects `201` with a 422 before Core's own clamp would run. `maxPageLimit` cannot be
  raised without Core raising that maximum first.
- **`Instances.ListAuditTrails` is the one exception.** That collection grows without bound,
  so "all of them" is not a useful default and `AuditTrailsOptions.Limit` is honoured as a
  ceiling: it walks only as many pages as it takes to reach the count asked for.

Walking rather than exposing pages is also what closes the silent-truncation trap. Core's
default page is 50 rows, so a list method returning one page would keep working right up
to the 51st row and then quietly omit it.

Two query parameters are now **dead**, and Core ignores rather than rejects an unknown one,
so both fail silently:

| Field | Sent as | State |
| --- | --- | --- |
| `InstancesListOptions.IncludeDeleted` | `include_deleted` | `GET /instances` takes only `cursor` and `limit`. Soft-deleted instances are excluded either way. |
| `AuditTrailsOptions.Offset` | `offset` | The endpoint moved from offset/limit to cursor pagination. An offset that used to skip entries now returns the newest ones. |

Both are kept and still sent: removing them is a breaking change to every caller, and
re-adding the parameters to Core would make them work again untouched. `Offset` has no
replacement -- `Limit` covers what it was reached for in practice, and only resuming from
where a previous call stopped is no longer expressible.

`pagination_test.go` pins all of this against a fake server that pages the way Core does:
that the envelope decodes at all, that a 451-row walk makes exactly three requests carrying
`limit=200` and the right `cursor`, that each page keeps the endpoint's own filters, that the
`Authorization` editor still runs on a hand-built request, that every request path is the
one the endpoint expects, and that a non-2xx still arrives as an `*Error` carrying its
status. That last one matters to any caller that branches on the status -- a
credential-provisioning loop reading 401 and 403 as "this credential is finished" would
mint a replacement on every transient failure if the status were lost.

## Development

Install [Task](https://taskfile.dev), then run `task --list` to see every
available workflow. The ones you will use:

```bash
task test      # go test -race -shuffle=on -cover ./...
task lint      # golangci-lint, pinned to the version CI runs
task fmt       # gofumpt + gci, the formatters .golangci.yml declares
task fix       # lint with --fix
```

Tools are fetched with `go run <tool>@<version>` against versions pinned in
`Taskfile.yml`, so there is nothing to install beyond Task and Go, and local
runs match CI exactly.

Before proposing a release:

```bash
task vuln           # govulncheck: advisories this code actually reaches
task deadcode       # functions no entry point or test can reach
task release:check  # gorelease: the exported API diff against the last tag
```

`task release:check` is the important one. This module is published, so every
tag is somebody's build: a changed or removed exported identifier is a breaking
change for every consumer, and `gorelease` is the only reliable way to spot one.

## License

This project is licensed under the [Apache License, Version 2.0](LICENSE).
