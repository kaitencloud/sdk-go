package sdk

import (
	"context"
	"fmt"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// List returns all entitlement groups.
func (s *EntitlementGroups) List(ctx context.Context) ([]EntitlementGroup, error) {
	return listAll[EntitlementGroup](ctx, s.client.list, listPath("entitlement-groups"), nil)
}

// Get returns the entitlement group identified by entitlementGroupSlug.
func (s *EntitlementGroups) Get(ctx context.Context, entitlementGroupSlug string) (EntitlementGroup, error) {
	resp, err := s.client.raw.GetEntitlementGroupWithResponse(ctx, entitlementGroupSlug)
	if err != nil {
		var zero EntitlementGroup
		return zero, fmt.Errorf("get entitlement group: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new entitlement group.
func (s *EntitlementGroups) Create(ctx context.Context, input EntitlementGroupInput) (EntitlementGroup, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero EntitlementGroup
		return zero, fmt.Errorf("encode entitlement group request: %w", err)
	}

	resp, err := s.client.raw.CreateEntitlementGroupWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero EntitlementGroup
		return zero, fmt.Errorf("create entitlement group: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the entitlement group identified by entitlementGroupSlug.
func (s *EntitlementGroups) Update(ctx context.Context, entitlementGroupSlug string, input EntitlementGroupInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode entitlement group request: %w", err)
	}

	resp, err := s.client.raw.UpdateEntitlementGroupWithBodyWithResponse(ctx, entitlementGroupSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update entitlement group: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Delete deletes the entitlement group identified by entitlementGroupSlug.
func (s *EntitlementGroups) Delete(ctx context.Context, entitlementGroupSlug string) error {
	resp, err := s.client.raw.DeleteEntitlementGroupWithResponse(ctx, entitlementGroupSlug)
	if err != nil {
		return fmt.Errorf("delete entitlement group: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// AddEntitlement adds the entitlement identified by entitlementSlug to the entitlement group identified by entitlementGroupSlug.
func (s *EntitlementGroups) AddEntitlement(ctx context.Context, entitlementGroupSlug, entitlementSlug string) error {
	body, err := jsonBody(struct {
		EntitlementSlug string `json:"entitlementSlug"`
	}{
		EntitlementSlug: entitlementSlug,
	})
	if err != nil {
		return fmt.Errorf("encode entitlement group request: %w", err)
	}

	resp, err := s.client.raw.AddEntitlementToGroupWithBodyWithResponse(ctx, entitlementGroupSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("add entitlement to group: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// RemoveEntitlement removes the entitlement identified by entitlementSlug from the entitlement group identified by entitlementGroupSlug.
func (s *EntitlementGroups) RemoveEntitlement(ctx context.Context, entitlementGroupSlug, entitlementSlug string) error {
	resp, err := s.client.raw.RemoveEntitlementFromGroupWithResponse(ctx, entitlementGroupSlug, entitlementSlug)
	if err != nil {
		return fmt.Errorf("remove entitlement from group: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// GetUsage returns the usage of the entitlement group identified by entitlementGroupSlug for the instance identified by instanceSlug.
func (s *EntitlementGroups) GetUsage(ctx context.Context, entitlementGroupSlug, instanceSlug string) ([]EntitlementGroupUsageItem, error) {
	params := &gen.GetEntitlementGroupUsageParams{Instance: instanceSlug}
	resp, err := s.client.raw.GetEntitlementGroupUsageWithResponse(ctx, entitlementGroupSlug, params)
	if err != nil {
		return nil, fmt.Errorf("get entitlement group usage: %w", err)
	}

	return expectSlice(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}
