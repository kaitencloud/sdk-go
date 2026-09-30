package sdk

import (
	"context"
	"fmt"
)

// List returns all deployment zones.
func (s *DeploymentZones) List(ctx context.Context) ([]DeploymentZone, error) {
	return listAll[DeploymentZone](ctx, s.client.list, listPath("deployment-zones"), nil)
}

// Get returns the deployment zone identified by deploymentZoneSlug.
func (s *DeploymentZones) Get(ctx context.Context, deploymentZoneSlug string) (DeploymentZone, error) {
	resp, err := s.client.raw.GetDeploymentZoneBySlugWithResponse(ctx, deploymentZoneSlug)
	if err != nil {
		var zero DeploymentZone
		return zero, fmt.Errorf("get deployment zone: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new deployment zone.
func (s *DeploymentZones) Create(ctx context.Context, input DeploymentZoneInput) (DeploymentZone, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero DeploymentZone
		return zero, fmt.Errorf("encode deployment zone request: %w", err)
	}

	resp, err := s.client.raw.CreateDeploymentZoneWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero DeploymentZone
		return zero, fmt.Errorf("create deployment zone: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the deployment zone identified by deploymentZoneSlug.
func (s *DeploymentZones) Update(ctx context.Context, deploymentZoneSlug string, input DeploymentZoneInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode deployment zone request: %w", err)
	}

	resp, err := s.client.raw.UpdateDeploymentZoneWithBodyWithResponse(ctx, deploymentZoneSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update deployment zone: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Delete deletes the deployment zone identified by deploymentZoneSlug.
func (s *DeploymentZones) Delete(ctx context.Context, deploymentZoneSlug string) error {
	resp, err := s.client.raw.DeleteDeploymentZoneWithResponse(ctx, deploymentZoneSlug)
	if err != nil {
		return fmt.Errorf("delete deployment zone: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
