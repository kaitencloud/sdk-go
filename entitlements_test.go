// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"net/http"
	"testing"
)

// TestCreateEntitlementSendsPresentationFields pins the presentation half of
// EntitlementInput onto the wire. It is the storefront contract: a pricing table renders
// exactly the entitlements marked userFacing, in displayOrder, with these unit labels --
// so a field silently dropped between the input struct and the request body is a catalog
// that looks complete in Go and renders as a flat, unlabelled list to a customer.
func TestCreateEntitlementSendsPresentationFields(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusCreated, `{"id":"e1","name":"Customers","slug":"customers"}`))

	entitlementType := EntitlementType("NUMBER")
	aggregation := EntitlementAggregationMethod("SUM")
	name, singular, plural, icon := "Customers", "customer", "customers", "lucide:users"
	userFacing := true
	displayOrder := int32(10)

	if _, err := client.Entitlements.Create(t.Context(), EntitlementInput{
		Name:              name,
		Slug:              &name,
		Type:              &entitlementType,
		AggregationMethod: &aggregation,
		GroupSlugs:        []string{"platform"},
		Icon:              &icon,
		UnitSingular:      &singular,
		UnitPlural:        &plural,
		UserFacing:        &userFacing,
		DisplayOrder:      &displayOrder,
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	body := sent.only(t).decodeBody(t)

	expected := map[string]any{
		"icon":         "lucide:users",
		"unitSingular": "customer",
		"unitPlural":   "customers",
		"userFacing":   true,
		"displayOrder": float64(10),
	}
	for key, want := range expected {
		got, ok := body[key]
		if !ok {
			t.Errorf("request body is missing %q", key)
			continue
		}
		if got != want {
			t.Errorf("%s = %v, want %v", key, got, want)
		}
	}
}

// TestCreateEntitlementOmitsUnsetPresentationFields keeps the raw counters raw. The API
// defaults userFacing to false and displayOrder to 0, and validates the unit labels as an
// all-or-none pair -- so sending zero values for fields the caller never set would put
// every internal counter into the pricing table with an empty unit label.
func TestCreateEntitlementOmitsUnsetPresentationFields(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusCreated, `{"id":"e1","name":"Customers Read","slug":"customers-read"}`))

	if _, err := client.Entitlements.Create(t.Context(), EntitlementInput{Name: "Customers Read"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	body := sent.only(t).decodeBody(t)

	for _, key := range []string{"icon", "unitSingular", "unitPlural", "userFacing", "displayOrder"} {
		if _, ok := body[key]; ok {
			t.Errorf("request body carries %q, want it omitted", key)
		}
	}

	// CreateEntitlementBody is additionalProperties:false and the Core API dropped
	// meterType entirely, so any key beyond what it publishes is a 422 on every create,
	// not a field the server ignores.
	if _, ok := body["meterType"]; ok {
		t.Error(`request body carries "meterType", which the Core API rejects outright`)
	}
}

// TestEntitlementWritesRenderTheWindowAndThresholdOnlyWhenSet pins both halves of the
// usage window and the warning threshold on the wire. Absence means something on each --
// no period is a lifetime counter, no anchor is CALENDAR, no threshold is 0 -- and none
// of them is nullable, so a key rendered from an unset field would be a null the bodies
// refuse. A zero threshold a caller does set has to reach the API, though: it is how a
// warning is switched off on purpose. Whole bodies are compared, because "unchanged for
// callers who set none of them" is a claim about the bytes, and the same body is
// expected from Create and Update, whose payloads differ in nothing here.
func TestEntitlementWritesRenderTheWindowAndThresholdOnlyWhenSet(t *testing.T) {
	t.Parallel()

	month := EntitlementResetPeriodMonth
	licenseStart := EntitlementResetAnchorLicenseStart
	off := int32(0)

	lifetime := EntitlementInput{Name: "Locations"}

	monthly := lifetime
	monthly.ResetPeriod = &month

	anchored := monthly
	anchored.ResetAnchor = &licenseStart

	silenced := lifetime
	silenced.WarningThresholdPercent = &off

	tests := []struct {
		name     string
		input    EntitlementInput
		expected string
	}{
		{
			name:     "nothing set renders none of them",
			input:    lifetime,
			expected: `{"name":"Locations","description":null}`,
		},
		{
			name:     "a period alone leaves the anchor to the server",
			input:    monthly,
			expected: `{"name":"Locations","description":null,"resetPeriod":"MONTH"}`,
		},
		{
			name:     "a period and its anchor",
			input:    anchored,
			expected: `{"name":"Locations","description":null,"resetPeriod":"MONTH","resetAnchor":"LICENSE_START"}`,
		},
		{
			name:     "a zero threshold is sent, not dropped",
			input:    silenced,
			expected: `{"name":"Locations","description":null,"warningThresholdPercent":0}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondToWrite)

			if _, err := client.Entitlements.Create(t.Context(), test.input); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if err := client.Entitlements.Update(t.Context(), "locations", test.input); err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			for i, want := range []struct{ method, path string }{
				{http.MethodPost, "/api/entitlements"},
				{http.MethodPut, "/api/entitlements/locations"},
			} {
				request := sent.at(t, i)
				if request.method != want.method || request.path != want.path {
					t.Errorf("request %d = %s %s, want %s %s", i, request.method, request.path, want.method, want.path)
				}
				if got := string(request.body); got != test.expected {
					t.Errorf("%s body = %s, want %s", want.method, got, test.expected)
				}
			}
		})
	}
}

// TestEntitlementUpdateCarriesThePeriodBack pins the read-modify-write the API requires
// of a periodic entitlement. Once an entitlement has a reset period, an update that does
// not send it back is refused -- the API reads the omission as an attempt to remove it --
// and an update that leaves the warning threshold out turns the warning off. Both are
// fields of Entitlement with the types EntitlementInput declares, so what Get returns
// can be carried over as it is, which is what this test does.
func TestEntitlementUpdateCarriesThePeriodBack(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet {
			return jsonResponse(req, http.StatusOK, `{"id":"e1","name":"Monthly Orders","description":null,"slug":"monthly-orders","type":"NUMBER","aggregationMethod":"SUM","resetPeriod":"MONTH","resetAnchor":"CALENDAR","warningThresholdPercent":80}`), nil
		}

		return jsonResponse(req, http.StatusNoContent, ""), nil
	})

	current, err := client.Entitlements.Get(t.Context(), "monthly-orders")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if err := client.Entitlements.Update(t.Context(), "monthly-orders", EntitlementInput{
		Name:                    "Orders per month",
		Description:             current.Description,
		Slug:                    current.Slug,
		Type:                    current.Type,
		AggregationMethod:       current.AggregationMethod,
		ResetPeriod:             current.ResetPeriod,
		ResetAnchor:             current.ResetAnchor,
		WarningThresholdPercent: current.WarningThresholdPercent,
	}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	update := sent.at(t, 1)
	if update.method != http.MethodPut || update.path != "/api/entitlements/monthly-orders" {
		t.Errorf("request = %s %s, want PUT /api/entitlements/monthly-orders", update.method, update.path)
	}

	expected := `{"name":"Orders per month","description":null,"type":"NUMBER","aggregationMethod":"SUM","slug":"monthly-orders","resetPeriod":"MONTH","resetAnchor":"CALENDAR","warningThresholdPercent":80}`
	if got := string(update.body); got != expected {
		t.Errorf("request body = %s, want %s", got, expected)
	}
}
