// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/kaitencloud/sdk-go/internal/gen"
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
// A create body and its update body routinely differ -- a component's predecessor is
// linked once, a license's lifecycle state is chosen on create and moved by its own
// operations afterwards -- so each direction is listed separately rather than assumed
// symmetric.
func TestRequestBodiesMatchTheirSpecBodies(t *testing.T) {
	t.Parallel()

	future := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	name := "pinned"
	scopes := []string{"read:instances"}
	draft := LicenseLifecycleDraft
	periodic := periodicEntitlementInput(name)

	grantValue, err := NumberLicenseValue(1000)
	if err != nil {
		t.Fatalf("NumberLicenseValue() error = %v", err)
	}

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
			expected: []string{"externalCustomerId", "name", "slug"},
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
			expected: []string{"description", "metadata", "name", "releaseId", "slug", "type"},
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
			expected: []string{"description", "name", "slug"},
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
			expected: []string{"description", "groupSlugs", "name", "slug"},
		},
		{
			name: "create periodic entitlement",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Entitlements.Create(ctx, periodic)
				return err
			},
			expected: []string{"aggregationMethod", "description", "name", "resetAnchor", "resetPeriod", "saleUnitFactor", "saleUnitPlural", "saleUnitSingular", "slug", "type", "unitPlural", "unitSingular", "warningThresholdPercent"},
		},
		{
			name: "update periodic entitlement",
			call: func(ctx context.Context, c *Client) error {
				return c.Entitlements.Update(ctx, name, periodic)
			},
			// The window renders on update too: once set, the API refuses an update
			// that does not carry it back.
			expected: []string{"aggregationMethod", "description", "name", "resetAnchor", "resetPeriod", "saleUnitFactor", "saleUnitPlural", "saleUnitSingular", "slug", "type", "unitPlural", "unitSingular", "warningThresholdPercent"},
		},
		{
			name: "associate license entitlement",
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.AssociateEntitlement(ctx, name, name, grantValue)
			},
			expected: []string{"entitlementSlug", "value"},
		},
		{
			name: "associate license entitlement with an allowance",
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.AssociateEntitlement(ctx, name, name, grantValue, WithOveragePercent(20))
			},
			expected: []string{"entitlementSlug", "limitCapExceededOveragePercent", "value"},
		},
		{
			name: "update license entitlement",
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.UpdateEntitlement(ctx, name, name, grantValue)
			},
			expected: []string{"value"},
		},
		{
			name: "update license entitlement with an allowance",
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.UpdateEntitlement(ctx, name, name, grantValue, WithOveragePercent(20))
			},
			expected: []string{"limitCapExceededOveragePercent", "value"},
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
			// Version renders as the read-back echo the server accepts. familySlug is
			// refused whatever its value and lifecycleState moves through Publish,
			// Archive and Unarchive, so neither renders even though both were set;
			// familyId and slug do, as the values a read-modify-write caller carries
			// back.
			expected: []string{"description", "familyId", "isDefault", "name", "slug", "type", "version", "versionName"},
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
			expected: []string{"name", "slug"},
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

// TestEntitlementAndGrantWritesReachEveryWritableKey fails when the spec gives the
// entitlement or grant bodies a field no input can send.
//
// TestRequestBodiesMatchTheirSpecBodies cannot catch that: its expected sets are
// transcribed by hand, so a field the spec adds is missing from the expectation exactly
// as it is missing from the input. That is how resetPeriod, resetAnchor,
// warningThresholdPercent, the sale-unit trio and limitCapExceededOveragePercent stayed
// out of reach while internal/gen already declared them -- and because both PUTs replace
// rather than patch, every Update reset what it could not send. This test reads the
// body's keys off the generated request type instead, so a regenerated client with a new
// key fails here until the input renders it or the key is listed below as read-only.
func TestEntitlementAndGrantWritesReachEveryWritableKey(t *testing.T) {
	t.Parallel()

	name := "pinned"
	complete := completeEntitlementInput(name)

	grantValue, err := NumberLicenseValue(1000)
	if err != nil {
		t.Fatalf("NumberLicenseValue() error = %v", err)
	}

	// The fields the spec marks readOnly: they are in the generated request type because
	// the bodies reuse the resource schemas, and the API ignores them on a write.
	entitlementReadOnly := []string{"createdAt", "entitlementGroups", "id", "updatedAt"}
	grantReadOnly := []string{"createdAt", "createdBy", "entitlementGroups", "entitlementName", "entitlementType", "licenseId", "licenseSlug", "updatedAt", "updatedBy"}

	tests := []struct {
		name     string
		body     any
		readOnly []string
		// unsent lists writable keys an input leaves out on purpose, each with its
		// reason next to the payload type.
		unsent []string
		call   func(context.Context, *Client) error
	}{
		{
			name:     "create entitlement",
			body:     gen.CreateEntitlementJSONRequestBody{},
			readOnly: entitlementReadOnly,
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Entitlements.Create(ctx, complete)
				return err
			},
		},
		{
			name:     "update entitlement",
			body:     gen.UpdateEntitlementJSONRequestBody{},
			readOnly: entitlementReadOnly,
			call: func(ctx context.Context, c *Client) error {
				return c.Entitlements.Update(ctx, name, complete)
			},
		},
		{
			name:     "associate license entitlement",
			body:     gen.AssociateEntitlementWithLicenseJSONRequestBody{},
			readOnly: grantReadOnly,
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.AssociateEntitlement(ctx, name, name, grantValue, WithOveragePercent(20))
			},
		},
		{
			name:     "update license entitlement",
			body:     gen.UpdateLicenseEntitlementJSONRequestBody{},
			readOnly: grantReadOnly,
			// See licenseEntitlementUpdatePayload: the path names the grant.
			unsent: []string{"entitlementSlug"},
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.UpdateEntitlement(ctx, name, name, grantValue, WithOveragePercent(20))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var writable []string
			for _, key := range jsonKeys(t, reflect.TypeOf(test.body)) {
				if !slices.Contains(test.readOnly, key) && !slices.Contains(test.unsent, key) {
					writable = append(writable, key)
				}
			}

			client, sent := givenClient(t, respondToWrite)

			if err := test.call(t.Context(), client); err != nil {
				t.Fatalf("call() error = %v", err)
			}

			if got := sent.only(t).bodyKeys(t); !slices.Equal(got, writable) {
				t.Errorf("request body keys = %v, want every writable key of %T: %v", got, test.body, writable)
			}
		})
	}
}

// periodicEntitlementInput is a NUMBER entitlement counted per calendar month, with a
// warning threshold and a sale unit, and none of the presentation fields.
func periodicEntitlementInput(name string) EntitlementInput {
	number := EntitlementType("NUMBER")
	sum := EntitlementAggregationMethod("SUM")
	month := EntitlementResetPeriodMonth
	calendar := EntitlementResetAnchorCalendar
	singular, plural := "order", "orders"
	saleSingular, salePlural := "pack", "packs"
	factor := 10.0
	warning := int32(80)

	return EntitlementInput{
		Name:                    name,
		Slug:                    &name,
		Type:                    &number,
		AggregationMethod:       &sum,
		UnitSingular:            &singular,
		UnitPlural:              &plural,
		SaleUnitSingular:        &saleSingular,
		SaleUnitPlural:          &salePlural,
		SaleUnitFactor:          &factor,
		ResetPeriod:             &month,
		ResetAnchor:             &calendar,
		WarningThresholdPercent: &warning,
	}
}

// completeEntitlementInput sets every field of EntitlementInput.
func completeEntitlementInput(name string) EntitlementInput {
	input := periodicEntitlementInput(name)

	description, icon := "Orders placed in the month", "lucide:receipt"
	userFacing := true
	displayOrder := int32(2)

	input.Description = &description
	input.GroupSlugs = []string{"restaurant-operations"}
	input.Icon = &icon
	input.UserFacing = &userFacing
	input.DisplayOrder = &displayOrder

	return input
}

// jsonKeys returns the JSON names of a generated struct's fields, sorted.
func jsonKeys(t *testing.T, structType reflect.Type) []string {
	t.Helper()

	if structType.Kind() != reflect.Struct {
		t.Fatalf("%s is a %s, want a struct", structType, structType.Kind())
	}

	keys := make([]string, 0, structType.NumField())
	for i := range structType.NumField() {
		name, _, _ := strings.Cut(structType.Field(i).Tag.Get("json"), ",")
		if name == "" || name == "-" {
			continue
		}
		keys = append(keys, name)
	}

	slices.Sort(keys)

	return keys
}

// TestUpdateLetsTheAPIRefuseARename pins the six updates that never rename.
//
// Their PUTs check a slug in the body against the path: the current one is accepted and
// any other is a 422, <Operation>.SlugNotRenameable. The update payloads used to drop
// Slug, so a caller asking for a rename got a success and kept the old slug without being
// told. The responder applies the API's rule, so this test fails on a payload that drops
// the slug again. The omitted half matters too: an input with no slug renders none.
func TestUpdateLetsTheAPIRefuseARename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		code   string
		update func(ctx context.Context, c *Client, slug *string) error
	}{
		{
			name: "customer",
			code: "UpdateCustomer.SlugNotRenameable",
			update: func(ctx context.Context, c *Client, slug *string) error {
				return c.Customers.Update(ctx, "current", CustomerInput{Name: "n", Slug: slug})
			},
		},
		{
			name: "deployment zone",
			code: "UpdateDeploymentZone.SlugNotRenameable",
			update: func(ctx context.Context, c *Client, slug *string) error {
				return c.DeploymentZones.Update(ctx, "current", DeploymentZoneInput{Name: "n", Type: "production", Slug: slug})
			},
		},
		{
			name: "entitlement group",
			code: "UpdateEntitlementGroup.SlugNotRenameable",
			update: func(ctx context.Context, c *Client, slug *string) error {
				return c.EntitlementGroups.Update(ctx, "current", EntitlementGroupInput{Name: "n", Slug: slug})
			},
		},
		{
			name: "entitlement",
			code: "UpdateEntitlement.SlugNotRenameable",
			update: func(ctx context.Context, c *Client, slug *string) error {
				return c.Entitlements.Update(ctx, "current", EntitlementInput{Name: "n", Slug: slug})
			},
		},
		{
			name: "license",
			code: "UpdateLicense.SlugNotRenameable",
			update: func(ctx context.Context, c *Client, slug *string) error {
				return c.Licenses.Update(ctx, "current", LicenseInput{Name: "n", Type: LicenseType("COMMUNITY"), Slug: slug})
			},
		},
		{
			name: "service account",
			code: "UpdateServiceAccount.SlugNotRenameable",
			update: func(ctx context.Context, c *Client, slug *string) error {
				return c.ServiceAccounts.Update(ctx, "current", ServiceAccountInput{Name: "n", Slug: slug})
			},
		},
	}

	current, renamed := "current", "renamed"

	for _, test := range tests {
		t.Run(test.name+" refuses another slug", func(t *testing.T) {
			t.Parallel()

			client, _ := givenClient(t, refuseRenames(test.code))

			err := test.update(t.Context(), client, &renamed)

			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("Update() error = %v, want the API's refusal", err)
			}

			if apiErr.StatusCode != http.StatusUnprocessableEntity || apiErr.Code != test.code {
				t.Errorf("Update() error = %d %s, want 422 %s", apiErr.StatusCode, apiErr.Code, test.code)
			}
		})

		t.Run(test.name+" accepts the current slug", func(t *testing.T) {
			t.Parallel()

			client, _ := givenClient(t, refuseRenames(test.code))

			if err := test.update(t.Context(), client, &current); err != nil {
				t.Fatalf("Update() error = %v", err)
			}
		})

		t.Run(test.name+" renders no slug when none is set", func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, refuseRenames(test.code))

			if err := test.update(t.Context(), client, nil); err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			if _, ok := sent.only(t).decodeBody(t)["slug"]; ok {
				t.Errorf("request body = %s, want no slug key", sent.only(t).body)
			}
		})
	}
}

// refuseRenames answers an update the way the PUTs that never rename do: a slug in the
// body other than the path's last segment is a 422 carrying code, anything else a 204.
func refuseRenames(code string) responder {
	return func(req *http.Request) (*http.Response, error) {
		var body struct {
			Slug string `json:"slug"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, err
		}

		if body.Slug != "" && body.Slug != path.Base(req.URL.Path) {
			return jsonResponse(req, http.StatusUnprocessableEntity,
				fmt.Sprintf(`{"title":"Unprocessable Entity","status":422,"code":%q,"detail":"slug cannot be changed through this endpoint; omit it or send the current slug"}`, code)), nil
		}

		return jsonResponse(req, http.StatusNoContent, ""), nil
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
