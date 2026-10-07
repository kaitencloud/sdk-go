// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

// pagedRows answers with the envelope Core actually sends, cutting rows into pages of
// the size the request asked for and issuing a cursor while any remain.
func pagedRows(rows []map[string]any) responder {
	return func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()

		start := 0
		if cursor := query.Get("cursor"); cursor != "" {
			if _, err := fmt.Sscanf(cursor, "row-%d", &start); err != nil {
				return nil, fmt.Errorf("bad cursor %q", cursor)
			}
		}

		limit := len(rows)
		if raw := query.Get("limit"); raw != "" {
			if _, err := fmt.Sscanf(raw, "%d", &limit); err != nil {
				return nil, err
			}
		}

		end := min(start+limit, len(rows))

		body := map[string]any{"items": rows[start:end], "hasMore": end < len(rows)}
		if end < len(rows) {
			body["nextCursor"] = fmt.Sprintf("row-%d", end)
		}

		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return jsonResponse(req, http.StatusOK, string(encoded)), nil
	}
}

func flagRows(count int) []map[string]any {
	rows := make([]map[string]any, 0, count)
	for i := range count {
		rows = append(rows, map[string]any{"slug": fmt.Sprintf("flag-%d", i)})
	}

	return rows
}

// A fake server that answers the envelope, not a bare array: a test fed a
// hand-written array would pass against a decoder that cannot read what Core
// actually sends.
func TestListDecodesThePaginatedEnvelope(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, pagedRows(flagRows(2)))

	flags, err := client.FeatureFlags.List(t.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(flags) != 2 {
		t.Fatalf("len(flags) = %d, want 2", len(flags))
	}

	// The rows themselves, not just the count: a decode that produced the right
	// number of zero-valued structs would satisfy a length check alone.
	for i, flag := range flags {
		want := fmt.Sprintf("flag-%d", i)
		if flag.Slug == nil {
			t.Errorf("flags[%d].Slug = nil, want %q", i, want)

			continue
		}

		if *flag.Slug != want {
			t.Errorf("flags[%d].Slug = %q, want %q", i, *flag.Slug, want)
		}
	}
}

// TestListWalksEveryPage pins the walk itself. Core's default page is 50 rows,
// so a list method that returned one page would keep working right up to the
// day somebody created the 51st row and then silently omit it.
func TestListWalksEveryPage(t *testing.T) {
	t.Parallel()

	const total = 451

	client, sent := givenClient(t, pagedRows(flagRows(total)))

	flags, err := client.FeatureFlags.List(t.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(flags) != total {
		t.Fatalf("len(flags) = %d, want %d", len(flags), total)
	}

	// 451 rows at the API's maximum page size: three full pages and a
	// remainder. Asserted so a walk that quietly stopped early, or one that
	// requested pages of 50, shows up as a failure here rather than as a
	// short answer nobody checked.
	if sent.count() != 3 {
		t.Fatalf("requests = %d, want 3", sent.count())
	}
	for i := range sent.count() {
		if got := sent.at(t, i).query.Get("limit"); got != "200" {
			t.Errorf("request %d limit = %q, want %q", i, got, "200")
		}
	}
	if got := sent.at(t, 0).query.Get("cursor"); got != "" {
		t.Errorf("first request sent cursor %q, want none", got)
	}
	if got := sent.at(t, 1).query.Get("cursor"); got != "row-200" {
		t.Errorf("second request cursor = %q, want %q", got, "row-200")
	}
}

// TestListPreservesEndpointFiltersAcrossPages guards the merge in walkPages: a
// filter dropped on page two turns one query into two different ones and
// returns rows the caller excluded.
func TestListPreservesEndpointFiltersAcrossPages(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, pagedRows(flagRows(300)))

	if _, err := client.MetadataFields.List(t.Context(), "INSTANCE"); err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if sent.count() != 2 {
		t.Fatalf("requests = %d, want 2", sent.count())
	}
	for i := range sent.count() {
		if got := sent.at(t, i).query.Get("resourceType"); got != "INSTANCE" {
			t.Errorf("request %d resourceType = %q, want %q", i, got, "INSTANCE")
		}
	}
}

// TestListAuditTrailsHonoursLimit covers the one list that is deliberately not
// walked to completion: the collection is unbounded, so a caller's ceiling is
// the answer rather than a starting point.
func TestListAuditTrailsHonoursLimit(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, pagedRows(flagRows(1000)))

	limit := int32(250)
	entries, err := client.Instances.ListAuditTrails(t.Context(), "instance-slug", &AuditTrailsOptions{Limit: &limit})
	if err != nil {
		t.Fatalf("ListAuditTrails() error = %v", err)
	}

	if len(entries) != 250 {
		t.Fatalf("len(entries) = %d, want 250", len(entries))
	}

	// Two pages, and the second asks for exactly the 50 still wanted rather
	// than another 200 to throw away.
	if sent.count() != 2 {
		t.Fatalf("requests = %d, want 2", sent.count())
	}
	if got := sent.at(t, 1).query.Get("limit"); got != "50" {
		t.Errorf("second request limit = %q, want %q", got, "50")
	}
}

// TestListPathsAndSubresourcePaths pins the paths the hand-built requests
// produce. Nothing but this test checks them, and a wrong path is a 404 at
// runtime rather than a compile error.
func TestListPathsAndSubresourcePaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(context.Context, *Client) error
		want string
	}{
		{
			name: "components",
			want: "/api/components",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Components.List(ctx)

				return err
			},
		},
		{
			name: "customers",
			want: "/api/customers",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Customers.List(ctx)

				return err
			},
		},
		{
			name: "deployment zones",
			want: "/api/deployment-zones",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.DeploymentZones.List(ctx)

				return err
			},
		},
		{
			name: "entitlement groups",
			want: "/api/entitlement-groups",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.EntitlementGroups.List(ctx)

				return err
			},
		},
		{
			name: "entitlements",
			want: "/api/entitlements",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Entitlements.List(ctx)

				return err
			},
		},
		{
			name: "feature flags",
			want: "/api/feature-flags",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.FeatureFlags.List(ctx)

				return err
			},
		},
		{
			name: "instances",
			want: "/api/instances",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Instances.List(ctx, nil)

				return err
			},
		},
		{
			name: "license families",
			want: "/api/license-families",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.LicenseFamilies.List(ctx)

				return err
			},
		},
		{
			name: "licenses",
			want: "/api/licenses",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.List(ctx)

				return err
			},
		},
		{
			name: "metadata fields",
			want: "/api/metadata-fields",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.MetadataFields.List(ctx, "INSTANCE")

				return err
			},
		},
		{
			name: "releases",
			want: "/api/releases",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Releases.List(ctx)

				return err
			},
		},
		{
			name: "service accounts",
			want: "/api/service-accounts",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.List(ctx)

				return err
			},
		},
		{
			name: "service account tokens",
			want: "/api/service-accounts/sa-slug/tokens",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.ListTokens(ctx, "sa-slug")

				return err
			},
		},
		{
			name: "license entitlements",
			want: "/api/licenses/beta-tester/entitlements",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.ListEntitlements(ctx, "beta-tester")

				return err
			},
		},
		{
			name: "audit trails",
			want: "/api/instances/instance-slug/audit-trails",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Instances.ListAuditTrails(ctx, "instance-slug", nil)

				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, pagedRows(nil))

			if err := test.call(t.Context(), client); err != nil {
				t.Fatalf("call error = %v", err)
			}

			if got := sent.only(t).path; got != test.want {
				t.Errorf("path = %q, want %q", got, test.want)
			}
		})
	}
}

// TestListReportsAPIErrorsAsError keeps the failure shape the generated calls
// produced. Callers branch on the status: a credential-provisioning loop reads
// 401 and 403 as "this credential is finished" and anything else as "the API
// could not be asked", so an error that lost its status would have it mint a
// replacement on every transient failure.
func TestListReportsAPIErrorsAsError(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusUnauthorized, `{"title":"Unauthorized","detail":"token revoked"}`))

	_, err := client.FeatureFlags.List(t.Context())
	if err == nil {
		t.Fatal("List() error = nil, want an error")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(err, *Error) = false for %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
	if apiErr.Problem == nil || stringValue(apiErr.Problem.Detail) != "token revoked" {
		t.Errorf("problem details lost: %+v", apiErr.Problem)
	}
}

// A body that is not a problem document at all still has to surface as an *Error
// carrying the status: gateways and load balancers answer with HTML, and a caller
// branching on 502 must not be handed a JSON decode failure instead.
func TestListReportsNonProblemErrorBodiesAsError(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusBadGateway, "<html>502 Bad Gateway</html>"))

	_, err := client.FeatureFlags.List(t.Context())

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(err, *Error) = false for %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadGateway)
	}
	if apiErr.Problem != nil {
		t.Errorf("Problem = %+v, want nil for a body that is not a problem document", apiErr.Problem)
	}
	if string(apiErr.Body) != "<html>502 Bad Gateway</html>" {
		t.Errorf("Body = %s, want the raw body preserved", apiErr.Body)
	}
}

// TestListRejectsAnUnusableCursor covers the two ways a server can claim more
// rows without a way to reach them. Both end as an error, because the
// alternative is either an infinite loop or a partial list returned as if it
// were complete.
func TestListRejectsAnUnusableCursor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{
			name: "more rows promised with no next cursor",
			body: `{"items":[{"slug":"a"}],"hasMore":true}`,
		},
		{
			name: "more rows promised with a cursor that never moves",
			body: `{"items":[{"slug":"a"}],"hasMore":true,"nextCursor":"stuck"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// A walk that failed to detect the stall would spin here forever, so the
			// responder gives up rather than letting the test hang the suite.
			calls := 0
			client, _ := givenClient(t, func(req *http.Request) (*http.Response, error) {
				calls++
				if calls > 10 {
					return nil, errors.New("looped")
				}

				return jsonResponse(req, http.StatusOK, test.body), nil
			})

			if _, err := client.FeatureFlags.List(t.Context()); err == nil {
				t.Fatal("List() error = nil, want an error")
			}
		})
	}
}

// TestListReturnsAnEmptySliceNotNil keeps the generated decode's contract: the
// endpoints promise "always an array, possibly empty", and a caller that
// distinguishes the two should not start seeing nil.
func TestListReturnsAnEmptySliceNotNil(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, pagedRows(nil))

	flags, err := client.FeatureFlags.List(t.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if flags == nil {
		t.Fatal("List() = nil, want an empty slice")
	}
	if len(flags) != 0 {
		t.Fatalf("len(flags) = %d, want 0", len(flags))
	}
}
