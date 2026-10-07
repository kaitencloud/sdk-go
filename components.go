// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// List returns all components.
func (s *Components) List(ctx context.Context) ([]Component, error) {
	return listAll[Component](ctx, s.client.list, listPath("components"), nil)
}

// Get returns the component identified by componentSlug.
func (s *Components) Get(ctx context.Context, componentSlug string) (Component, error) {
	resp, err := s.client.raw.GetComponentWithResponse(ctx, componentSlug)
	if err != nil {
		var zero Component
		return zero, fmt.Errorf("get component: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new component.
func (s *Components) Create(ctx context.Context, input ComponentInput) (Component, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero Component
		return zero, fmt.Errorf("encode component request: %w", err)
	}

	resp, err := s.client.raw.CreateComponentWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Component
		return zero, fmt.Errorf("create component: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the component identified by componentSlug.
func (s *Components) Update(ctx context.Context, componentSlug string, input ComponentInput) (Component, error) {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		var zero Component
		return zero, fmt.Errorf("encode component request: %w", err)
	}

	resp, err := s.client.raw.UpdateComponentWithBodyWithResponse(ctx, componentSlug, contentTypeJSON, body)
	if err != nil {
		var zero Component
		return zero, fmt.Errorf("update component: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Delete deletes the component identified by componentSlug.
func (s *Components) Delete(ctx context.Context, componentSlug string) error {
	resp, err := s.client.raw.DeleteComponentWithResponse(ctx, componentSlug)
	if err != nil {
		return fmt.Errorf("delete component: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
