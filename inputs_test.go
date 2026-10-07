// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"net/http"
	"slices"
	"testing"
	"time"
)

// TestRequestBodiesMatchTheirSpecBodies pins the exact key set every write renders.
//
// It exists because three separate never-worked paths were found by reading the spec
// rather than by running anything: EntitlementInput rendered a meterType no body accepts,
// LicenseInput rendered a version and an isActive no body accepts, and the entitlement
// bodies' required-but-nullable description was omitted rather than sent as null. Every
// request body in the Core API is additionalProperties:false and huma enforces both
// halves, so a key too many and a required key missing are the same outcome: a 422 on
// every call, with no type error to catch it.
//
// The expected sets below are transcribed from app/openapi.yaml's Create*/Update* bodies.
// A create body and its update body routinely differ -- a slug is assigned once and then
// immutable, a license's lifecycle state is chosen on create and moved by its own
// operations afterwards -- so each direction is listed separately rather than assumed
// symmetric.
func TestRequestBodiesMatchTheirSpecBodies(t *testing.T) {
	t.Parallel()

	future := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	name := "pinned"
	scopes := []string{"read:instances"}
	draft := LicenseLifecycleDraft

	tests := []struct {
		name     string
		call     func(context.Context, *Client) error
		expected []string
	}{
		{
			name: "create component",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Components.Create(ctx, ComponentInput{Name: name, Version: "1.0.0", PreviousComponentID: &name, Slug: &name})
				return err
			},
			expected: []string{"name", "previousComponentId", "slug", "version"},
		},
		{
			name: "update component",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Components.Update(ctx, name, ComponentInput{Name: name, Version: "1.0.0", PreviousComponentID: &name, Slug: &name})
				return err
			},
			expected: []string{"name", "slug", "version"},
		},
		{
			name: "create customer",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Customers.Create(ctx, CustomerInput{Name: name, ExternalCustomerID: &name, Slug: &name})
				return err
			},
			expected: []string{"externalCustomerId", "name", "slug"},
		},
		{
			name: "update customer",
			call: func(ctx context.Context, c *Client) error {
				return c.Customers.Update(ctx, name, CustomerInput{Name: name, ExternalCustomerID: &name, Slug: &name})
			},
			expected: []string{"externalCustomerId", "name"},
		},
		{
			name: "create deployment zone",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.DeploymentZones.Create(ctx, DeploymentZoneInput{Name: name, Type: "production", Description: "d", Metadata: map[string]any{"region": "eu-west-1"}, ReleaseID: &name, Slug: &name})
				return err
			},
			expected: []string{"description", "metadata", "name", "releaseId", "slug", "type"},
		},
		{
			name: "update deployment zone",
			call: func(ctx context.Context, c *Client) error {
				return c.DeploymentZones.Update(ctx, name, DeploymentZoneInput{Name: name, Type: "production", Description: "d", Metadata: map[string]any{"region": "eu-west-1"}, ReleaseID: &name, Slug: &name})
			},
			expected: []string{"description", "metadata", "name", "releaseId", "type"},
		},
		{
			name: "create entitlement group",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.EntitlementGroups.Create(ctx, EntitlementGroupInput{Name: name, Slug: &name})
				return err
			},
			// description with no value set: required and nullable, so it renders as null.
			expected: []string{"description", "name", "slug"},
		},
		{
			name: "update entitlement group",
			call: func(ctx context.Context, c *Client) error {
				return c.EntitlementGroups.Update(ctx, name, EntitlementGroupInput{Name: name, Slug: &name})
			},
			expected: []string{"description", "name"},
		},
		{
			name: "create entitlement",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Entitlements.Create(ctx, EntitlementInput{Name: name, Slug: &name, GroupSlugs: []string{"platform"}})
				return err
			},
			expected: []string{"description", "groupSlugs", "name", "slug"},
		},
		{
			name: "update entitlement",
			call: func(ctx context.Context, c *Client) error {
				return c.Entitlements.Update(ctx, name, EntitlementInput{Name: name, Slug: &name, GroupSlugs: []string{"platform"}})
			},
			expected: []string{"description", "groupSlugs", "name"},
		},
		{
			name: "create instance",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Instances.Create(ctx, InstanceInput{Name: name, CustomerID: name, LicenseID: name, DeploymentZoneID: &name, StartLicenseDate: future, EndLicenseDate: future, Slug: &name})
				return err
			},
			expected: []string{"customerId", "deploymentZoneId", "description", "endLicenseDate", "licenseId", "metadata", "name", "slug", "startLicenseDate"},
		},
		{
			name: "update instance",
			call: func(ctx context.Context, c *Client) error {
				return c.Instances.Update(ctx, name, InstanceInput{Name: name, CustomerID: name, LicenseID: name, DeploymentZoneID: &name, StartLicenseDate: future, EndLicenseDate: future, Slug: &name})
			},
			// Slug is optional here and renames when sent, so it renders like it does on
			// create; TestUpdateRendersTheSlugOnlyWhenSet pins the omitted half.
			expected: []string{"customerId", "deploymentZoneId", "description", "endLicenseDate", "licenseId", "metadata", "name", "slug", "startLicenseDate"},
		},
		// The two PATCH writes were the only Core writes this test did not cover,
		// which is exactly why UpdateLifecycleStage shipped rendering the stale
		// generated `lifecycle_stage` key and 422'd on every call. One key each is
		// the whole assertion: PatchInstanceBody is additionalProperties:false and
		// absence is how it spells "leave the other field alone", so a second key
		// here would be as wrong as a misspelled first one.
		{
			name: "patch instance status",
			call: func(ctx context.Context, c *Client) error {
				return c.Instances.UpdateStatus(ctx, name, InstanceStatusHealthy)
			},
			expected: []string{"status"},
		},
		{
			name: "patch instance lifecycle stage",
			call: func(ctx context.Context, c *Client) error {
				return c.Instances.UpdateLifecycleStage(ctx, name, "ACTIVE")
			},
			expected: []string{"lifecycleStage"},
		},
		{
			name: "create license",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Create(ctx, LicenseInput{Name: name, Description: "d", Type: LicenseType("COMMUNITY"), Version: "1", VersionName: &name, IsDefault: true, Slug: &name, FamilySlug: &name, FamilyID: &name, LifecycleState: &draft})
				return err
			},
			// No version: the server assigns it as the next in the family. Both family
			// handles render when set -- whether they agree is the server's call, not
			// the SDK's -- and TestCreateRendersTheFamilyFieldsOnlyWhenSet pins the
			// omitted half.
			expected: []string{"description", "familyId", "familySlug", "isDefault", "lifecycleState", "name", "slug", "type", "versionName"},
		},
		{
			name: "update license",
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.Update(ctx, name, LicenseInput{Name: name, Description: "d", Type: LicenseType("COMMUNITY"), Version: "1", VersionName: &name, IsDefault: true, Slug: &name, FamilySlug: &name, FamilyID: &name, LifecycleState: &draft})
			},
			// Version renders as the read-back echo the server accepts. Slug and
			// familySlug are refused and lifecycleState moves through Publish,
			// Archive and Unarchive, so none of the three renders even though all
			// were set; familyId does, as the value a read-modify-write caller
			// carries back.
			expected: []string{"description", "familyId", "isDefault", "name", "type", "version", "versionName"},
		},
		{
			name: "create service account",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.Create(ctx, ServiceAccountInput{Name: name, Slug: &name})
				return err
			},
			expected: []string{"name", "slug"},
		},
		{
			name: "update service account",
			call: func(ctx context.Context, c *Client) error {
				return c.ServiceAccounts.Update(ctx, name, ServiceAccountInput{Name: name, Slug: &name})
			},
			expected: []string{"name"},
		},
		{
			name: "create service account token",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.CreateToken(ctx, name, TokenInput{Name: name, Slug: &name, Scopes: scopes, ExpiresAt: &future})
				return err
			},
			expected: []string{"expiresAt", "name", "scopes", "slug"},
		},
		{
			name: "create service account token inheriting scopes",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.CreateToken(ctx, name, TokenInput{Name: name})
				return err
			},
			// scopes is required and nullable: null inherits the service account's own.
			expected: []string{"name", "scopes"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondToWrite)

			if err := test.call(t.Context(), client); err != nil {
				t.Fatalf("call() error = %v", err)
			}

			if got := sent.only(t).bodyKeys(t); !slices.Equal(got, test.expected) {
				t.Errorf("request body keys = %v, want %v", got, test.expected)
			}
		})
	}
}

// respondToWrite answers a write the way the listener does: creates answer 201 and
// updates 200, which is what the generated client matches on to fill in its typed
// response fields. The body is the smallest one every wrapper here can decode.
func respondToWrite(req *http.Request) (*http.Response, error) {
	status := http.StatusOK
	if req.Method == http.MethodPost {
		status = http.StatusCreated
	}

	return jsonResponse(req, status, `{"id":"1","name":"pinned","slug":"pinned","token":"ksh_x"}`), nil
}
