// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// List returns all customers.
func (s *Customers) List(ctx context.Context) ([]Customer, error) {
	return listAll[Customer](ctx, s.client.list, listPath("customers"), nil)
}

// Get returns the customer identified by customerSlug.
func (s *Customers) Get(ctx context.Context, customerSlug string) (Customer, error) {
	resp, err := s.client.raw.GetCustomerWithResponse(ctx, customerSlug)
	if err != nil {
		var zero Customer
		return zero, fmt.Errorf("get customer: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new customer.
func (s *Customers) Create(ctx context.Context, input CustomerInput) (Customer, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero Customer
		return zero, fmt.Errorf("encode customer request: %w", err)
	}

	resp, err := s.client.raw.CreateCustomerWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Customer
		return zero, fmt.Errorf("create customer: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the customer identified by customerSlug.
func (s *Customers) Update(ctx context.Context, customerSlug string, input CustomerInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode customer request: %w", err)
	}

	resp, err := s.client.raw.UpdateCustomerWithBodyWithResponse(ctx, customerSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Delete deletes the customer identified by customerSlug.
func (s *Customers) Delete(ctx context.Context, customerSlug string) error {
	resp, err := s.client.raw.DeleteCustomerWithResponse(ctx, customerSlug)
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
