// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// List returns all entitlements.
func (s *Entitlements) List(ctx context.Context) ([]Entitlement, error) {
	return listAll[Entitlement](ctx, s.client.list, listPath("entitlements"), nil)
}

// Get returns the entitlement identified by entitlementSlug.
func (s *Entitlements) Get(ctx context.Context, entitlementSlug string) (Entitlement, error) {
	resp, err := s.client.raw.GetEntitlementWithResponse(ctx, entitlementSlug)
	if err != nil {
		var zero Entitlement
		return zero, fmt.Errorf("get entitlement: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new entitlement.
func (s *Entitlements) Create(ctx context.Context, input EntitlementInput) (Entitlement, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero Entitlement
		return zero, fmt.Errorf("encode entitlement request: %w", err)
	}

	resp, err := s.client.raw.CreateEntitlementWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Entitlement
		return zero, fmt.Errorf("create entitlement: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the entitlement identified by entitlementSlug.
//
// It replaces the entitlement rather than patching it: an optional field left nil is
// reset, and a periodic entitlement's ResetPeriod and ResetAnchor have to be sent back
// as they are. EntitlementInput says which fields that covers; reading the entitlement
// with Get and carrying its fields over is the safe way to change one of them.
func (s *Entitlements) Update(ctx context.Context, entitlementSlug string, input EntitlementInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode entitlement request: %w", err)
	}

	resp, err := s.client.raw.UpdateEntitlementWithBodyWithResponse(ctx, entitlementSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update entitlement: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Delete deletes the entitlement identified by entitlementSlug.
func (s *Entitlements) Delete(ctx context.Context, entitlementSlug string) error {
	resp, err := s.client.raw.DeleteEntitlementWithResponse(ctx, entitlementSlug)
	if err != nil {
		return fmt.Errorf("delete entitlement: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
