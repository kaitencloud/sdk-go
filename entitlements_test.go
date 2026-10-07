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
