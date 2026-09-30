package sdk

import (
	"context"
	"fmt"
)

// List returns all feature flags.
func (s *FeatureFlags) List(ctx context.Context) ([]FeatureFlag, error) {
	return listAll[FeatureFlag](ctx, s.client.list, listPath("feature-flags"), nil)
}

// Get returns the feature flag identified by featureFlagSlug.
func (s *FeatureFlags) Get(ctx context.Context, featureFlagSlug string) (FeatureFlag, error) {
	resp, err := s.client.raw.GetFeatureFlagWithResponse(ctx, featureFlagSlug)
	if err != nil {
		var zero FeatureFlag
		return zero, fmt.Errorf("get feature flag: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new feature flag.
func (s *FeatureFlags) Create(ctx context.Context, input FeatureFlag) (FeatureFlag, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero FeatureFlag
		return zero, fmt.Errorf("encode feature flag request: %w", err)
	}

	resp, err := s.client.raw.CreateFeatureFlagWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero FeatureFlag
		return zero, fmt.Errorf("create feature flag: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the feature flag identified by featureFlagSlug.
func (s *FeatureFlags) Update(ctx context.Context, featureFlagSlug string, input FeatureFlag) error {
	body, err := jsonBody(input)
	if err != nil {
		return fmt.Errorf("encode feature flag request: %w", err)
	}

	resp, err := s.client.raw.UpdateFeatureFlagWithBodyWithResponse(ctx, featureFlagSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update feature flag: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Delete deletes the feature flag identified by featureFlagSlug.
func (s *FeatureFlags) Delete(ctx context.Context, featureFlagSlug string) error {
	resp, err := s.client.raw.DeleteFeatureFlagWithResponse(ctx, featureFlagSlug)
	if err != nil {
		return fmt.Errorf("delete feature flag: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
