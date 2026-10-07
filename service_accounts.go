// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// List returns all service accounts.
func (s *ServiceAccounts) List(ctx context.Context) ([]ServiceAccount, error) {
	return listAll[ServiceAccount](ctx, s.client.list, listPath("service-accounts"), nil)
}

// Get returns the service account identified by serviceAccountSlug.
func (s *ServiceAccounts) Get(ctx context.Context, serviceAccountSlug string) (ServiceAccount, error) {
	resp, err := s.client.raw.GetServiceAccountWithResponse(ctx, serviceAccountSlug)
	if err != nil {
		var zero ServiceAccount
		return zero, fmt.Errorf("get service account: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new service account.
func (s *ServiceAccounts) Create(ctx context.Context, input ServiceAccountInput) (ServiceAccount, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero ServiceAccount
		return zero, fmt.Errorf("encode service account request: %w", err)
	}

	resp, err := s.client.raw.CreateServiceAccountWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero ServiceAccount
		return zero, fmt.Errorf("create service account: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the service account identified by serviceAccountSlug.
func (s *ServiceAccounts) Update(ctx context.Context, serviceAccountSlug string, input ServiceAccountInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode service account request: %w", err)
	}

	resp, err := s.client.raw.UpdateServiceAccountWithBodyWithResponse(ctx, serviceAccountSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update service account: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// ListTokens returns the tokens belonging to the service account identified by serviceAccountSlug.
func (s *ServiceAccounts) ListTokens(ctx context.Context, serviceAccountSlug string) ([]Token, error) {
	return listAll[Token](ctx, s.client.list, listPath("service-accounts", serviceAccountSlug, "tokens"), nil)
}

// CreateToken creates a new token for the service account identified by serviceAccountSlug.
func (s *ServiceAccounts) CreateToken(ctx context.Context, serviceAccountSlug string, input TokenInput) (PlainToken, error) {
	body, err := jsonBody(input.payload())
	if err != nil {
		var zero PlainToken
		return zero, fmt.Errorf("encode service account token request: %w", err)
	}

	resp, err := s.client.raw.CreateServiceAccountTokenWithBodyWithResponse(ctx, serviceAccountSlug, contentTypeJSON, body)
	if err != nil {
		var zero PlainToken
		return zero, fmt.Errorf("create service account token: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// DeleteToken deletes the token identified by tokenSlug belonging to the service account identified by serviceAccountSlug.
func (s *ServiceAccounts) DeleteToken(ctx context.Context, serviceAccountSlug, tokenSlug string) error {
	resp, err := s.client.raw.DeleteServiceAccountTokenWithResponse(ctx, serviceAccountSlug, tokenSlug)
	if err != nil {
		return fmt.Errorf("delete service account token: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
