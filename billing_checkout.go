// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// BillingCheckoutSessions is the organization's view of its anonymous checkouts.
type BillingCheckoutSessions struct {
	client *Client
}

// ListPage returns one page of what List returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingCheckoutSessions) ListPage(ctx context.Context, params *ListCheckoutSessionsParams) (PageCheckoutSessionAdmin, error) {
	resp, err := s.client.raw.ListCheckoutSessionsWithResponse(ctx, params)
	if err != nil {
		var zero PageCheckoutSessionAdmin
		return zero, fmt.Errorf("list checkout sessions: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the organization's anonymous checkouts, newest first, following the
// pagination to its end. Never a ksc_ token or a voucher code.
func (s *BillingCheckoutSessions) List(ctx context.Context, params *ListCheckoutSessionsParams) ([]CheckoutSessionAdmin, error) {
	var query ListCheckoutSessionsParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []CheckoutSessionAdmin
	for {
		page, err := s.ListPage(ctx, &query)
		if err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		if !page.HasMore || page.NextCursor == nil || *page.NextCursor == "" {
			return items, nil
		}
		if query.Cursor != nil && *query.Cursor == *page.NextCursor {
			return nil, fmt.Errorf("list checkout sessions: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Get returns one anonymous checkout, as the organization reads it.
func (s *BillingCheckoutSessions) Get(ctx context.Context, sessionID string) (CheckoutSessionAdmin, error) {
	resp, err := s.client.raw.GetCheckoutSessionAdminWithResponse(ctx, sessionID)
	if err != nil {
		var zero CheckoutSessionAdmin
		return zero, fmt.Errorf("get checkout session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CreateCheckoutSession starts a checkout for a visitor the vendor does not know yet.
// It takes the publishable key and a Turnstile token, and writes only the session: the
// customer, the instance and the subscription come when it completes. The returned Token
// is the session's ksc_, shown once: GetCheckoutSession and ConfirmCheckoutSession take a
// Client built with it as the bearer token.
func (s *BillingPublic) CreateCheckoutSession(ctx context.Context, input CheckoutSessionDraft) (CheckoutSessionStart, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero CheckoutSessionStart
		return zero, fmt.Errorf("encode create checkout session request: %w", err)
	}

	resp, err := s.client.raw.CreateCheckoutSessionWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero CheckoutSessionStart
		return zero, fmt.Errorf("create checkout session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectEither(apiResp, resp.JSON200, resp.JSON201)
}

// GetCheckoutSession returns the session the client's ksc_ token was minted for: its
// status, what it bought and what is due and, once completed, the customer and instance
// it created. The token reads it for 24 hours after it closed.
func (s *BillingPublic) GetCheckoutSession(ctx context.Context, sessionID string) (CheckoutSession, error) {
	resp, err := s.client.raw.GetCheckoutSessionWithResponse(ctx, sessionID)
	if err != nil {
		var zero CheckoutSession
		return zero, fmt.Errorf("get checkout session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ConfirmCheckoutSession completes the session when nothing has to be captured, or says
// what is missing: requires_payment_method comes with the hosted page where the visitor
// saves one, and the next confirmation completes it. It takes the session's ksc_ token.
func (s *BillingPublic) ConfirmCheckoutSession(ctx context.Context, sessionID string) (CheckoutConfirmation, error) {
	resp, err := s.client.raw.ConfirmCheckoutSessionWithResponse(ctx, sessionID)
	if err != nil {
		var zero CheckoutConfirmation
		return zero, fmt.Errorf("confirm checkout session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}
