// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// Billing groups the billing API: the organization's settings and provider, and a namespace per resource.
type Billing struct {
	client *Client

	Handoff          *BillingHandoff
	Invoices         *BillingInvoices
	Subscriptions    *BillingSubscriptions
	InstanceAddons   *BillingInstanceAddons
	InstanceVouchers *BillingInstanceVouchers
	LicensePrices    *BillingLicensePrices
	Addons           *BillingAddons
	Vouchers         *BillingVouchers
	PublishableKeys  *BillingPublishableKeys
	CustomerSessions *BillingCustomerSessions
	Customers        *BillingCustomers
	Public           *BillingPublic
	Session          *BillingSession
}

// BillingHandoff is the accounting handoff: invoices claimed in leased batches and acknowledged one by one.
type BillingHandoff struct {
	client *Client
}

// BillingInvoices are the organization's invoices.
type BillingInvoices struct {
	client *Client
}

// BillingSubscriptions are the instances' subscriptions.
type BillingSubscriptions struct {
	client *Client
}

// BillingInstanceAddons are the add-ons attached to instances.
type BillingInstanceAddons struct {
	client *Client
}

// BillingInstanceVouchers are the vouchers redeemed on instances.
type BillingInstanceVouchers struct {
	client *Client
}

// BillingLicensePrices are the prices of licence versions.
type BillingLicensePrices struct {
	client *Client
}

// BillingAddons is the add-on catalogue.
type BillingAddons struct {
	client *Client
}

// BillingVouchers are the organization's vouchers.
type BillingVouchers struct {
	client *Client
}

// BillingPublishableKeys are the pk_ keys a public page reads the catalogue with.
type BillingPublishableKeys struct {
	client *Client
}

// BillingCustomerSessions are the kst_ sessions a customer's page acts with.
type BillingCustomerSessions struct {
	client *Client
}

// BillingCustomers is the customers' billing: their payment method and the provider's portal.
type BillingCustomers struct {
	client *Client
}

// BillingPublic is what a publishable key (pk_) reaches. Build the Client with the key as its token.
type BillingPublic struct {
	client *Client
}

// BillingSession is what a customer session (kst_) reaches. Build the Client with the session secret as its token.
type BillingSession struct {
	client *Client
}

func newBilling(c *Client) *Billing {
	return &Billing{
		client:           c,
		Handoff:          &BillingHandoff{client: c},
		Invoices:         &BillingInvoices{client: c},
		Subscriptions:    &BillingSubscriptions{client: c},
		InstanceAddons:   &BillingInstanceAddons{client: c},
		InstanceVouchers: &BillingInstanceVouchers{client: c},
		LicensePrices:    &BillingLicensePrices{client: c},
		Addons:           &BillingAddons{client: c},
		Vouchers:         &BillingVouchers{client: c},
		PublishableKeys:  &BillingPublishableKeys{client: c},
		CustomerSessions: &BillingCustomerSessions{client: c},
		Customers:        &BillingCustomers{client: c},
		Public:           &BillingPublic{client: c},
		Session:          &BillingSession{client: c},
	}
}

// Settings returns the organization's billing settings.
func (s *Billing) Settings(ctx context.Context) (BillingSettings, error) {
	resp, err := s.client.raw.GetBillingSettingsWithResponse(ctx)
	if err != nil {
		var zero BillingSettings
		return zero, fmt.Errorf("get billing settings: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// UpdateSettings changes the organization's billing settings: a member left nil is left alone.
func (s *Billing) UpdateSettings(ctx context.Context, input BillingSettings) (BillingSettings, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero BillingSettings
		return zero, fmt.Errorf("encode update billing settings request: %w", err)
	}

	resp, err := s.client.raw.UpdateBillingSettingsWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero BillingSettings
		return zero, fmt.Errorf("update billing settings: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Capabilities reports what the configured billing provider can do.
func (s *Billing) Capabilities(ctx context.Context) (BillingCapabilities, error) {
	resp, err := s.client.raw.GetBillingCapabilitiesWithResponse(ctx)
	if err != nil {
		var zero BillingCapabilities
		return zero, fmt.Errorf("get billing capabilities: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Health reports the billing provider's health and the work waiting on it.
func (s *Billing) Health(ctx context.Context) (BillingHealth, error) {
	resp, err := s.client.raw.GetBillingHealthWithResponse(ctx)
	if err != nil {
		var zero BillingHealth
		return zero, fmt.Errorf("get billing health: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Sync asks Kaiten to reconcile with the billing provider; the work runs in the background.
func (s *Billing) Sync(ctx context.Context) (SyncReport, error) {
	resp, err := s.client.raw.SyncBillingProviderWithResponse(ctx)
	if err != nil {
		var zero SyncReport
		return zero, fmt.Errorf("sync billing provider: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON202)
}

// ClosePeriods closes the billing periods whose boundary has passed, as the scheduler would.
func (s *Billing) ClosePeriods(ctx context.Context, input ClosePeriodsScope) (ClosePeriodsReport, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero ClosePeriodsReport
		return zero, fmt.Errorf("encode close billing periods request: %w", err)
	}

	resp, err := s.client.raw.CloseBillingPeriodsWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero ClosePeriodsReport
		return zero, fmt.Errorf("close billing periods: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListPage returns one page of what List returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingHandoff) ListPage(ctx context.Context, params *ListHandoffParams) (PageQueuedInvoice, error) {
	resp, err := s.client.raw.ListHandoffWithResponse(ctx, params)
	if err != nil {
		var zero PageQueuedInvoice
		return zero, fmt.Errorf("list handoff: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the invoices waiting for, or leased to, the accounting handoff. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingHandoff) List(ctx context.Context, params *ListHandoffParams) ([]QueuedInvoice, error) {
	var query ListHandoffParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []QueuedInvoice
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
			return nil, fmt.Errorf("list handoff: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Claim leases a batch of invoices to hand off: each stays leased to this claim until it is acknowledged or the lease runs out.
func (s *BillingHandoff) Claim(ctx context.Context, input HandoffClaimSize) (HandoffClaim, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero HandoffClaim
		return zero, fmt.Errorf("encode claim handoff request: %w", err)
	}

	resp, err := s.client.raw.ClaimHandoffWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero HandoffClaim
		return zero, fmt.Errorf("claim handoff: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Ack acknowledges the handoff of an invoice, with the reference the accounting system gave it.
func (s *BillingHandoff) Ack(ctx context.Context, invoiceID string, input HandoffBooking) (Invoice, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("encode ack handoff request: %w", err)
	}

	resp, err := s.client.raw.AckHandoffWithBodyWithResponse(ctx, invoiceID, contentTypeJSON, body)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("ack handoff: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListPage returns one page of what List returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingInvoices) ListPage(ctx context.Context, params *ListInvoicesParams) (PageInvoiceSummary, error) {
	resp, err := s.client.raw.ListInvoicesWithResponse(ctx, params)
	if err != nil {
		var zero PageInvoiceSummary
		return zero, fmt.Errorf("list invoices: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the organization's invoices that match params. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingInvoices) List(ctx context.Context, params *ListInvoicesParams) ([]InvoiceSummary, error) {
	var query ListInvoicesParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []InvoiceSummary
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
			return nil, fmt.Errorf("list invoices: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Get returns the invoice identified by invoiceID, with its lines.
func (s *BillingInvoices) Get(ctx context.Context, invoiceID string) (Invoice, error) {
	resp, err := s.client.raw.GetInvoiceWithResponse(ctx, invoiceID)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("get invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Export returns the invoices of a range as a CSV or JSON document, exactly as the API sends it.
func (s *BillingInvoices) Export(ctx context.Context, params *ExportInvoicesParams) ([]byte, error) {
	resp, err := s.client.raw.ExportInvoicesWithResponse(ctx, params)
	if err != nil {
		var zero []byte
		return zero, fmt.Errorf("export invoices: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectBody(apiResp)
}

// Void voids an invoice.
func (s *BillingInvoices) Void(ctx context.Context, invoiceID string, input InvoiceVoid) (Invoice, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("encode void invoice request: %w", err)
	}

	resp, err := s.client.raw.VoidInvoiceWithBodyWithResponse(ctx, invoiceID, contentTypeJSON, body)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("void invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Recompose composes a draft invoice again from what is known now.
func (s *BillingInvoices) Recompose(ctx context.Context, invoiceID string) (Invoice, error) {
	resp, err := s.client.raw.RecomposeInvoiceWithResponse(ctx, invoiceID)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("recompose invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectEither(apiResp, resp.JSON200, resp.JSON201)
}

// MarkPaid records an invoice as paid outside the billing provider.
func (s *BillingInvoices) MarkPaid(ctx context.Context, invoiceID string, input InvoicePayment) (Invoice, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("encode mark invoice paid request: %w", err)
	}

	resp, err := s.client.raw.MarkInvoicePaidWithBodyWithResponse(ctx, invoiceID, contentTypeJSON, body)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("mark invoice paid: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// WriteOff writes an invoice off as uncollectable.
func (s *BillingInvoices) WriteOff(ctx context.Context, invoiceID string, input InvoiceWriteOff) (Invoice, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("encode write off invoice request: %w", err)
	}

	resp, err := s.client.raw.WriteOffInvoiceWithBodyWithResponse(ctx, invoiceID, contentTypeJSON, body)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("write off invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ReleaseHold releases an invoice held for review.
func (s *BillingInvoices) ReleaseHold(ctx context.Context, invoiceID string, input HoldRelease) (Invoice, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("encode release invoice hold request: %w", err)
	}

	resp, err := s.client.raw.ReleaseInvoiceHoldWithBodyWithResponse(ctx, invoiceID, contentTypeJSON, body)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("release invoice hold: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// RetryPush pushes an invoice to the billing provider again; the push runs in the background.
func (s *BillingInvoices) RetryPush(ctx context.Context, invoiceID string) (Invoice, error) {
	resp, err := s.client.raw.RetryInvoicePushWithResponse(ctx, invoiceID)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("retry invoice push: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON202)
}

// Sync reads an invoice's state back from the billing provider.
func (s *BillingInvoices) Sync(ctx context.Context, invoiceID string) (Invoice, error) {
	resp, err := s.client.raw.SyncInvoiceWithResponse(ctx, invoiceID)
	if err != nil {
		var zero Invoice
		return zero, fmt.Errorf("sync invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListLineReports returns the usage reports a metered invoice line was rated from.
func (s *BillingInvoices) ListLineReports(ctx context.Context, invoiceID string, lineID string, params *ListInvoiceLineReportsParams) (LineReportPage, error) {
	resp, err := s.client.raw.ListInvoiceLineReportsWithResponse(ctx, invoiceID, lineID, params)
	if err != nil {
		var zero LineReportPage
		return zero, fmt.Errorf("list invoice line reports: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Get returns the billing of the instance identified by instanceSlug: its subscription, if any.
func (s *BillingSubscriptions) Get(ctx context.Context, instanceSlug string) (InstanceBilling, error) {
	resp, err := s.client.raw.GetInstanceBillingWithResponse(ctx, instanceSlug)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("get instance billing: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Subscribe subscribes an instance to a price of its licence.
func (s *BillingSubscriptions) Subscribe(ctx context.Context, instanceSlug string, input NewSubscription) (StartedSubscription, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero StartedSubscription
		return zero, fmt.Errorf("encode subscribe instance request: %w", err)
	}

	resp, err := s.client.raw.SubscribeInstanceWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		var zero StartedSubscription
		return zero, fmt.Errorf("subscribe instance: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Update changes an instance's subscription.
func (s *BillingSubscriptions) Update(ctx context.Context, instanceSlug string, input SubscriptionTerms) (InstanceBilling, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("encode update instance billing request: %w", err)
	}

	resp, err := s.client.raw.UpdateInstanceBillingWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("update instance billing: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Cancel cancels an instance's subscription, now or at the end of its period.
func (s *BillingSubscriptions) Cancel(ctx context.Context, instanceSlug string, input SubscriptionCancellation) (CanceledSubscription, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero CanceledSubscription
		return zero, fmt.Errorf("encode cancel subscription request: %w", err)
	}

	resp, err := s.client.raw.CancelSubscriptionWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		var zero CanceledSubscription
		return zero, fmt.Errorf("cancel subscription: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Reactivate withdraws a cancellation scheduled for the end of the period.
func (s *BillingSubscriptions) Reactivate(ctx context.Context, instanceSlug string) (InstanceBilling, error) {
	resp, err := s.client.raw.ReactivateSubscriptionWithResponse(ctx, instanceSlug)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("reactivate subscription: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// SchedulePlanChange schedules a change of price for the instance's next period.
func (s *BillingSubscriptions) SchedulePlanChange(ctx context.Context, instanceSlug string, input PlanChangeTarget) (InstanceBilling, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("encode schedule plan change request: %w", err)
	}

	resp, err := s.client.raw.SchedulePlanChangeWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("schedule plan change: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CancelPlanChange withdraws the instance's scheduled change of price.
func (s *BillingSubscriptions) CancelPlanChange(ctx context.Context, instanceSlug string) (InstanceBilling, error) {
	resp, err := s.client.raw.CancelPlanChangeWithResponse(ctx, instanceSlug)
	if err != nil {
		var zero InstanceBilling
		return zero, fmt.Errorf("cancel plan change: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// UpcomingInvoice returns what the instance's next invoice would be, measured so far.
func (s *BillingSubscriptions) UpcomingInvoice(ctx context.Context, instanceSlug string) (InvoicePreview, error) {
	resp, err := s.client.raw.GetUpcomingInvoiceWithResponse(ctx, instanceSlug)
	if err != nil {
		var zero InvoicePreview
		return zero, fmt.Errorf("get upcoming invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListInvoicesPage returns one page of what ListInvoices returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingSubscriptions) ListInvoicesPage(ctx context.Context, instanceSlug string, params *ListInstanceInvoicesParams) (PageInvoiceSummary, error) {
	resp, err := s.client.raw.ListInstanceInvoicesWithResponse(ctx, instanceSlug, params)
	if err != nil {
		var zero PageInvoiceSummary
		return zero, fmt.Errorf("list instance invoices: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListInvoices returns the invoices of the instance identified by instanceSlug. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingSubscriptions) ListInvoices(ctx context.Context, instanceSlug string, params *ListInstanceInvoicesParams) ([]InvoiceSummary, error) {
	var query ListInstanceInvoicesParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []InvoiceSummary
	for {
		page, err := s.ListInvoicesPage(ctx, instanceSlug, &query)
		if err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		if !page.HasMore || page.NextCursor == nil || *page.NextCursor == "" {
			return items, nil
		}
		if query.Cursor != nil && *query.Cursor == *page.NextCursor {
			return nil, fmt.Errorf("list instance invoices: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// List returns the add-ons attached to an instance.
func (s *BillingInstanceAddons) List(ctx context.Context, instanceSlug string, params *ListInstanceAddonsParams) ([]InstanceAddon, error) {
	resp, err := s.client.raw.ListInstanceAddonsWithResponse(ctx, instanceSlug, params)
	if err != nil {
		var zero []InstanceAddon
		return zero, fmt.Errorf("list instance addons: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectSlice(apiResp, resp.JSON200)
}

// Attach attaches an add-on to an instance.
func (s *BillingInstanceAddons) Attach(ctx context.Context, instanceSlug string, input NewInstanceAddon) (InstanceAddon, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero InstanceAddon
		return zero, fmt.Errorf("encode attach instance addon request: %w", err)
	}

	resp, err := s.client.raw.AttachInstanceAddonWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		var zero InstanceAddon
		return zero, fmt.Errorf("attach instance addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// SetQuantity changes the quantity of an add-on attached to an instance.
func (s *BillingInstanceAddons) SetQuantity(ctx context.Context, instanceSlug string, addonSlug string, input AddonQuantity) (InstanceAddon, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero InstanceAddon
		return zero, fmt.Errorf("encode set instance addon quantity request: %w", err)
	}

	resp, err := s.client.raw.SetInstanceAddonQuantityWithBodyWithResponse(ctx, instanceSlug, addonSlug, contentTypeJSON, body)
	if err != nil {
		var zero InstanceAddon
		return zero, fmt.Errorf("set instance addon quantity: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Detach removes an add-on from an instance.
func (s *BillingInstanceAddons) Detach(ctx context.Context, instanceSlug string, addonSlug string) error {
	resp, err := s.client.raw.DetachInstanceAddonWithResponse(ctx, instanceSlug, addonSlug)
	if err != nil {
		return fmt.Errorf("detach instance addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// List returns the vouchers redeemed on an instance.
func (s *BillingInstanceVouchers) List(ctx context.Context, instanceSlug string, params *ListInstanceVouchersParams) ([]Redemption, error) {
	resp, err := s.client.raw.ListInstanceVouchersWithResponse(ctx, instanceSlug, params)
	if err != nil {
		var zero []Redemption
		return zero, fmt.Errorf("list instance vouchers: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectSlice(apiResp, resp.JSON200)
}

// Redeem redeems a voucher on an instance.
func (s *BillingInstanceVouchers) Redeem(ctx context.Context, instanceSlug string, input RedeemableCode) (Redemption, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Redemption
		return zero, fmt.Errorf("encode redeem voucher request: %w", err)
	}

	resp, err := s.client.raw.RedeemVoucherWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		var zero Redemption
		return zero, fmt.Errorf("redeem voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Revoke revokes a redemption.
func (s *BillingInstanceVouchers) Revoke(ctx context.Context, instanceSlug string, instanceVoucherID string, input RevocationReason) (Redemption, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Redemption
		return zero, fmt.Errorf("encode revoke instance voucher request: %w", err)
	}

	resp, err := s.client.raw.RevokeInstanceVoucherWithBodyWithResponse(ctx, instanceSlug, instanceVoucherID, contentTypeJSON, body)
	if err != nil {
		var zero Redemption
		return zero, fmt.Errorf("revoke instance voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the prices of a licence version.
func (s *BillingLicensePrices) List(ctx context.Context, licenseSlug string, params *ListLicensePricesParams) ([]Price, error) {
	resp, err := s.client.raw.ListLicensePricesWithResponse(ctx, licenseSlug, params)
	if err != nil {
		var zero []Price
		return zero, fmt.Errorf("list license prices: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectSlice(apiResp, resp.JSON200)
}

// Create creates a price on a licence version.
func (s *BillingLicensePrices) Create(ctx context.Context, licenseSlug string, input NewLicensePrice) (Price, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("encode create license price request: %w", err)
	}

	resp, err := s.client.raw.CreateLicensePriceWithBodyWithResponse(ctx, licenseSlug, contentTypeJSON, body)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("create license price: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Get returns a price of a licence version.
func (s *BillingLicensePrices) Get(ctx context.Context, licenseSlug string, priceID string) (Price, error) {
	resp, err := s.client.raw.GetLicensePriceWithResponse(ctx, licenseSlug, priceID)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("get license price: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Update changes a price of a licence version.
func (s *BillingLicensePrices) Update(ctx context.Context, licenseSlug string, priceID string, input LicensePriceChanges) (Price, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("encode update license price request: %w", err)
	}

	resp, err := s.client.raw.UpdateLicensePriceWithBodyWithResponse(ctx, licenseSlug, priceID, contentTypeJSON, body)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("update license price: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Deprecate deprecates a price: subscriptions keep it, new ones cannot take it.
func (s *BillingLicensePrices) Deprecate(ctx context.Context, licenseSlug string, priceID string) (Price, error) {
	resp, err := s.client.raw.DeprecateLicensePriceWithResponse(ctx, licenseSlug, priceID)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("deprecate license price: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Preview rates a scenario against a licence version's prices, without billing anything.
func (s *BillingLicensePrices) Preview(ctx context.Context, licenseSlug string, input InvoicePreviewScenario) (InvoicePreview, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero InvoicePreview
		return zero, fmt.Errorf("encode preview license invoice request: %w", err)
	}

	resp, err := s.client.raw.PreviewLicenseInvoiceWithBodyWithResponse(ctx, licenseSlug, contentTypeJSON, body)
	if err != nil {
		var zero InvoicePreview
		return zero, fmt.Errorf("preview license invoice: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListPage returns one page of what List returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingAddons) ListPage(ctx context.Context, params *ListAddonsParams) (PageAddon, error) {
	resp, err := s.client.raw.ListAddonsWithResponse(ctx, params)
	if err != nil {
		var zero PageAddon
		return zero, fmt.Errorf("list addons: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the organization's add-ons. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingAddons) List(ctx context.Context, params *ListAddonsParams) ([]Addon, error) {
	var query ListAddonsParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []Addon
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
			return nil, fmt.Errorf("list addons: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Create creates an add-on.
func (s *BillingAddons) Create(ctx context.Context, input NewAddon) (Addon, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("encode create addon request: %w", err)
	}

	resp, err := s.client.raw.CreateAddonWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("create addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Get returns the add-on identified by addonSlug.
func (s *BillingAddons) Get(ctx context.Context, addonSlug string) (Addon, error) {
	resp, err := s.client.raw.GetAddonWithResponse(ctx, addonSlug)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("get addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Update changes an add-on.
func (s *BillingAddons) Update(ctx context.Context, addonSlug string, input AddonChanges) (Addon, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("encode update addon request: %w", err)
	}

	resp, err := s.client.raw.UpdateAddonWithBodyWithResponse(ctx, addonSlug, contentTypeJSON, body)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("update addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Delete deletes a draft add-on.
func (s *BillingAddons) Delete(ctx context.Context, addonSlug string) error {
	resp, err := s.client.raw.DeleteAddonWithResponse(ctx, addonSlug)
	if err != nil {
		return fmt.Errorf("delete addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// Publish publishes a draft add-on.
func (s *BillingAddons) Publish(ctx context.Context, addonSlug string) (Addon, error) {
	resp, err := s.client.raw.PublishAddonWithResponse(ctx, addonSlug)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("publish addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Archive archives an add-on: attachments keep it, new ones cannot take it.
func (s *BillingAddons) Archive(ctx context.Context, addonSlug string) (Addon, error) {
	resp, err := s.client.raw.ArchiveAddonWithResponse(ctx, addonSlug)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("archive addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Unarchive brings an archived add-on back.
func (s *BillingAddons) Unarchive(ctx context.Context, addonSlug string) (Addon, error) {
	resp, err := s.client.raw.UnarchiveAddonWithResponse(ctx, addonSlug)
	if err != nil {
		var zero Addon
		return zero, fmt.Errorf("unarchive addon: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListCompatibility returns the licence families an add-on can be attached under.
func (s *BillingAddons) ListCompatibility(ctx context.Context, addonSlug string) (CompatibleLicenseFamilies, error) {
	resp, err := s.client.raw.ListAddonCompatibilityWithResponse(ctx, addonSlug)
	if err != nil {
		var zero CompatibleLicenseFamilies
		return zero, fmt.Errorf("list addon compatibility: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// SetCompatibility makes an add-on attachable under a licence family.
func (s *BillingAddons) SetCompatibility(ctx context.Context, addonSlug string, familySlug string) error {
	resp, err := s.client.raw.SetAddonCompatibilityWithResponse(ctx, addonSlug, familySlug)
	if err != nil {
		return fmt.Errorf("set addon compatibility: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// RemoveCompatibility stops an add-on being attachable under a licence family.
func (s *BillingAddons) RemoveCompatibility(ctx context.Context, addonSlug string, familySlug string) error {
	resp, err := s.client.raw.RemoveAddonCompatibilityWithResponse(ctx, addonSlug, familySlug)
	if err != nil {
		return fmt.Errorf("remove addon compatibility: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// ListEntitlements returns what an add-on grants.
func (s *BillingAddons) ListEntitlements(ctx context.Context, addonSlug string) ([]AddonEntitlement, error) {
	resp, err := s.client.raw.ListAddonEntitlementsWithResponse(ctx, addonSlug)
	if err != nil {
		var zero []AddonEntitlement
		return zero, fmt.Errorf("list addon entitlements: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectSlice(apiResp, resp.JSON200)
}

// AssignEntitlement makes an add-on grant an entitlement.
func (s *BillingAddons) AssignEntitlement(ctx context.Context, addonSlug string, input NewAddonGrant) (AddonEntitlement, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero AddonEntitlement
		return zero, fmt.Errorf("encode assign addon entitlement request: %w", err)
	}

	resp, err := s.client.raw.AssignAddonEntitlementWithBodyWithResponse(ctx, addonSlug, contentTypeJSON, body)
	if err != nil {
		var zero AddonEntitlement
		return zero, fmt.Errorf("assign addon entitlement: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// GetEntitlement returns an add-on's grant of an entitlement.
func (s *BillingAddons) GetEntitlement(ctx context.Context, addonSlug string, entitlementSlug string) (AddonEntitlement, error) {
	resp, err := s.client.raw.GetAddonEntitlementWithResponse(ctx, addonSlug, entitlementSlug)
	if err != nil {
		var zero AddonEntitlement
		return zero, fmt.Errorf("get addon entitlement: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// UpdateEntitlement changes an add-on's grant of an entitlement.
func (s *BillingAddons) UpdateEntitlement(ctx context.Context, addonSlug string, entitlementSlug string, input AddonGrant) (AddonEntitlement, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero AddonEntitlement
		return zero, fmt.Errorf("encode update addon entitlement request: %w", err)
	}

	resp, err := s.client.raw.UpdateAddonEntitlementWithBodyWithResponse(ctx, addonSlug, entitlementSlug, contentTypeJSON, body)
	if err != nil {
		var zero AddonEntitlement
		return zero, fmt.Errorf("update addon entitlement: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// UnassignEntitlement removes an add-on's grant of an entitlement.
func (s *BillingAddons) UnassignEntitlement(ctx context.Context, addonSlug string, entitlementSlug string) error {
	resp, err := s.client.raw.UnassignAddonEntitlementWithResponse(ctx, addonSlug, entitlementSlug)
	if err != nil {
		return fmt.Errorf("unassign addon entitlement: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// ListPrices returns an add-on's prices.
func (s *BillingAddons) ListPrices(ctx context.Context, addonSlug string, params *ListAddonPricesParams) ([]Price, error) {
	resp, err := s.client.raw.ListAddonPricesWithResponse(ctx, addonSlug, params)
	if err != nil {
		var zero []Price
		return zero, fmt.Errorf("list addon prices: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectSlice(apiResp, resp.JSON200)
}

// CreatePrice creates a price on an add-on.
func (s *BillingAddons) CreatePrice(ctx context.Context, addonSlug string, input NewAddonPrice) (Price, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("encode create addon price request: %w", err)
	}

	resp, err := s.client.raw.CreateAddonPriceWithBodyWithResponse(ctx, addonSlug, contentTypeJSON, body)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("create addon price: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// DeprecatePrice deprecates a price of an add-on.
func (s *BillingAddons) DeprecatePrice(ctx context.Context, addonSlug string, priceID string) (Price, error) {
	resp, err := s.client.raw.DeprecateAddonPriceWithResponse(ctx, addonSlug, priceID)
	if err != nil {
		var zero Price
		return zero, fmt.Errorf("deprecate addon price: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListFamilies returns the organization's add-on families.
func (s *BillingAddons) ListFamilies(ctx context.Context) ([]AddonFamily, error) {
	resp, err := s.client.raw.ListAddonFamiliesWithResponse(ctx)
	if err != nil {
		var zero []AddonFamily
		return zero, fmt.Errorf("list addon families: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectSlice(apiResp, resp.JSON200)
}

// GetFamily returns an add-on family.
func (s *BillingAddons) GetFamily(ctx context.Context, familySlug string) (AddonFamily, error) {
	resp, err := s.client.raw.GetAddonFamilyWithResponse(ctx, familySlug)
	if err != nil {
		var zero AddonFamily
		return zero, fmt.Errorf("get addon family: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// UpdateFamily changes an add-on family.
func (s *BillingAddons) UpdateFamily(ctx context.Context, familySlug string, input AddonFamilyVisibility) (AddonFamily, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero AddonFamily
		return zero, fmt.Errorf("encode update addon family request: %w", err)
	}

	resp, err := s.client.raw.UpdateAddonFamilyWithBodyWithResponse(ctx, familySlug, contentTypeJSON, body)
	if err != nil {
		var zero AddonFamily
		return zero, fmt.Errorf("update addon family: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListPage returns one page of what List returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingVouchers) ListPage(ctx context.Context, params *ListVouchersParams) (PageVoucher, error) {
	resp, err := s.client.raw.ListVouchersWithResponse(ctx, params)
	if err != nil {
		var zero PageVoucher
		return zero, fmt.Errorf("list vouchers: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the organization's vouchers. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingVouchers) List(ctx context.Context, params *ListVouchersParams) ([]Voucher, error) {
	var query ListVouchersParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []Voucher
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
			return nil, fmt.Errorf("list vouchers: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Create creates a draft voucher.
func (s *BillingVouchers) Create(ctx context.Context, input VoucherDraft) (Voucher, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("encode create voucher request: %w", err)
	}

	resp, err := s.client.raw.CreateVoucherWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("create voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Get returns the voucher identified by voucherID.
func (s *BillingVouchers) Get(ctx context.Context, voucherID string) (Voucher, error) {
	resp, err := s.client.raw.GetVoucherWithResponse(ctx, voucherID)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("get voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Update changes a voucher.
func (s *BillingVouchers) Update(ctx context.Context, voucherID string, input VoucherDraft) (Voucher, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("encode update voucher request: %w", err)
	}

	resp, err := s.client.raw.UpdateVoucherWithBodyWithResponse(ctx, voucherID, contentTypeJSON, body)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("update voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Delete deletes a voucher no instance ever redeemed, whatever its status, and frees its
// code. One redeemed at least once, or held by a checkout waiting for its payment, is
// refused with 409 DeleteVoucher.Redeemed: Archive retires it instead.
func (s *BillingVouchers) Delete(ctx context.Context, voucherID string) error {
	resp, err := s.client.raw.DeleteVoucherWithResponse(ctx, voucherID)
	if err != nil {
		return fmt.Errorf("delete voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// Publish publishes a draft voucher: its code can be redeemed from now on.
func (s *BillingVouchers) Publish(ctx context.Context, voucherID string) (Voucher, error) {
	resp, err := s.client.raw.PublishVoucherWithResponse(ctx, voucherID)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("publish voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Archive archives a voucher: redemptions keep it, its code is redeemed no more.
func (s *BillingVouchers) Archive(ctx context.Context, voucherID string) (Voucher, error) {
	resp, err := s.client.raw.ArchiveVoucherWithResponse(ctx, voucherID)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("archive voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListRedemptionsPage returns one page of what ListRedemptions returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingVouchers) ListRedemptionsPage(ctx context.Context, voucherID string, params *ListVoucherRedemptionsParams) (PageRedemption, error) {
	resp, err := s.client.raw.ListVoucherRedemptionsWithResponse(ctx, voucherID, params)
	if err != nil {
		var zero PageRedemption
		return zero, fmt.Errorf("list voucher redemptions: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListRedemptions returns the redemptions of a voucher. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingVouchers) ListRedemptions(ctx context.Context, voucherID string, params *ListVoucherRedemptionsParams) ([]Redemption, error) {
	var query ListVoucherRedemptionsParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []Redemption
	for {
		page, err := s.ListRedemptionsPage(ctx, voucherID, &query)
		if err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		if !page.HasMore || page.NextCursor == nil || *page.NextCursor == "" {
			return items, nil
		}
		if query.Cursor != nil && *query.Cursor == *page.NextCursor {
			return nil, fmt.Errorf("list voucher redemptions: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Lookup finds a voucher by its code.
func (s *BillingVouchers) Lookup(ctx context.Context, input VoucherCode) (Voucher, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("encode lookup voucher request: %w", err)
	}

	resp, err := s.client.raw.LookupVoucherWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Voucher
		return zero, fmt.Errorf("lookup voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Validate says whether a code could be redeemed on an instance, without redeeming it.
func (s *BillingVouchers) Validate(ctx context.Context, input VoucherCheck) (Validity, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero Validity
		return zero, fmt.Errorf("encode validate voucher request: %w", err)
	}

	resp, err := s.client.raw.ValidateVoucherWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Validity
		return zero, fmt.Errorf("validate voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListPage returns one page of what List returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingPublishableKeys) ListPage(ctx context.Context, params *ListPublishableKeysParams) (PagePublishableKey, error) {
	resp, err := s.client.raw.ListPublishableKeysWithResponse(ctx, params)
	if err != nil {
		var zero PagePublishableKey
		return zero, fmt.Errorf("list publishable keys: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// List returns the organization's publishable keys. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingPublishableKeys) List(ctx context.Context, params *ListPublishableKeysParams) ([]PublishableKey, error) {
	var query ListPublishableKeysParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []PublishableKey
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
			return nil, fmt.Errorf("list publishable keys: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Create creates a publishable key; the answer is the only time its secret is shown.
func (s *BillingPublishableKeys) Create(ctx context.Context, input PublishableKeyDraft) (PublishableKeyCreated, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero PublishableKeyCreated
		return zero, fmt.Errorf("encode create publishable key request: %w", err)
	}

	resp, err := s.client.raw.CreatePublishableKeyWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero PublishableKeyCreated
		return zero, fmt.Errorf("create publishable key: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Update changes a publishable key.
func (s *BillingPublishableKeys) Update(ctx context.Context, keyID string, input PublishableKeyPatch) (PublishableKey, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero PublishableKey
		return zero, fmt.Errorf("encode update publishable key request: %w", err)
	}

	resp, err := s.client.raw.UpdatePublishableKeyWithBodyWithResponse(ctx, keyID, contentTypeJSON, body)
	if err != nil {
		var zero PublishableKey
		return zero, fmt.Errorf("update publishable key: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Revoke revokes a publishable key.
func (s *BillingPublishableKeys) Revoke(ctx context.Context, keyID string) (PublishableKey, error) {
	resp, err := s.client.raw.RevokePublishableKeyWithResponse(ctx, keyID)
	if err != nil {
		var zero PublishableKey
		return zero, fmt.Errorf("revoke publishable key: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Create creates a customer session; the answer is the only time its secret is shown.
func (s *BillingCustomerSessions) Create(ctx context.Context, input CustomerSessionDraft) (CreatedCustomerSession, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero CreatedCustomerSession
		return zero, fmt.Errorf("encode create customer session request: %w", err)
	}

	resp, err := s.client.raw.CreateCustomerSessionWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero CreatedCustomerSession
		return zero, fmt.Errorf("create customer session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON201)
}

// Revoke revokes a customer session.
func (s *BillingCustomerSessions) Revoke(ctx context.Context, sessionID string) error {
	resp, err := s.client.raw.RevokeCustomerSessionWithResponse(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("revoke customer session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// Get returns the billing of the customer identified by customerSlug.
func (s *BillingCustomers) Get(ctx context.Context, customerSlug string) (CustomerBilling, error) {
	resp, err := s.client.raw.GetCustomerBillingWithResponse(ctx, customerSlug)
	if err != nil {
		var zero CustomerBilling
		return zero, fmt.Errorf("get customer billing: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// DetachPaymentMethod removes the customer's default payment method.
func (s *BillingCustomers) DetachPaymentMethod(ctx context.Context, customerSlug string) error {
	resp, err := s.client.raw.DetachPaymentMethodWithResponse(ctx, customerSlug)
	if err != nil {
		return fmt.Errorf("detach payment method: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectNoContent(apiResp)
}

// CreatePaymentMethodSession starts a hosted session where the customer sets a payment method.
func (s *BillingCustomers) CreatePaymentMethodSession(ctx context.Context, customerSlug string, input NewPaymentMethodSession) (PaymentMethodSession, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero PaymentMethodSession
		return zero, fmt.Errorf("encode create payment method session request: %w", err)
	}

	resp, err := s.client.raw.CreatePaymentMethodSessionWithBodyWithResponse(ctx, customerSlug, contentTypeJSON, body)
	if err != nil {
		var zero PaymentMethodSession
		return zero, fmt.Errorf("create payment method session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CompletePaymentMethodSession reads back the payment method a session set.
func (s *BillingCustomers) CompletePaymentMethodSession(ctx context.Context, customerSlug string, sessionID string) (CompletedPaymentMethodSession, error) {
	resp, err := s.client.raw.CompletePaymentMethodSessionWithResponse(ctx, customerSlug, sessionID)
	if err != nil {
		var zero CompletedPaymentMethodSession
		return zero, fmt.Errorf("complete payment method session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CreatePortalSession opens the billing provider's portal for the customer.
func (s *BillingCustomers) CreatePortalSession(ctx context.Context, customerSlug string, input NewPortalSession) (PortalSession, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero PortalSession
		return zero, fmt.Errorf("encode create portal session request: %w", err)
	}

	resp, err := s.client.raw.CreatePortalSessionWithBodyWithResponse(ctx, customerSlug, contentTypeJSON, body)
	if err != nil {
		var zero PortalSession
		return zero, fmt.Errorf("create portal session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Catalog returns the public catalogue a publishable key exposes.
func (s *BillingPublic) Catalog(ctx context.Context, params *GetPublicCatalogParams) (PublicCatalog, error) {
	resp, err := s.client.raw.GetPublicCatalogWithResponse(ctx, params)
	if err != nil {
		var zero PublicCatalog
		return zero, fmt.Errorf("get public catalog: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Portal returns what the session's customer sees of their billing.
func (s *BillingSession) Portal(ctx context.Context) (SessionPortal, error) {
	resp, err := s.client.raw.GetSessionPortalWithResponse(ctx)
	if err != nil {
		var zero SessionPortal
		return zero, fmt.Errorf("get session portal: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListInvoicesPage returns one page of what ListInvoices returns: pass the page's NextCursor as
// params.Cursor for the next one.
func (s *BillingSession) ListInvoicesPage(ctx context.Context, params *ListSessionInvoicesParams) (PageSessionInvoice, error) {
	resp, err := s.client.raw.ListSessionInvoicesWithResponse(ctx, params)
	if err != nil {
		var zero PageSessionInvoice
		return zero, fmt.Errorf("list session invoices: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ListInvoices returns the session's invoices. It walks every page; params.Cursor is ignored, and
// params.Limit is the page size (the largest the API takes when nil).
func (s *BillingSession) ListInvoices(ctx context.Context, params *ListSessionInvoicesParams) ([]SessionInvoice, error) {
	var query ListSessionInvoicesParams
	if params != nil {
		query = *params
	}
	query.Cursor = nil
	if query.Limit == nil {
		limit := int32(maxPageLimit)
		query.Limit = &limit
	}
	var items []SessionInvoice
	for {
		page, err := s.ListInvoicesPage(ctx, &query)
		if err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		if !page.HasMore || page.NextCursor == nil || *page.NextCursor == "" {
			return items, nil
		}
		if query.Cursor != nil && *query.Cursor == *page.NextCursor {
			return nil, fmt.Errorf("list session invoices: the API returned the same cursor twice")
		}
		next := *page.NextCursor
		query.Cursor = &next
	}
}

// Checkout subscribes the session's instance, or opens the provider's checkout.
func (s *BillingSession) Checkout(ctx context.Context, input SessionCheckoutOrder) (SessionCheckout, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero SessionCheckout
		return zero, fmt.Errorf("encode create session checkout request: %w", err)
	}

	resp, err := s.client.raw.CreateSessionCheckoutWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero SessionCheckout
		return zero, fmt.Errorf("create session checkout: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectEither(apiResp, resp.JSON200, resp.JSON201)
}

// Cancel cancels the session's subscription.
func (s *BillingSession) Cancel(ctx context.Context, input SessionCancellation) (SessionSubscription, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero SessionSubscription
		return zero, fmt.Errorf("encode cancel session subscription request: %w", err)
	}

	resp, err := s.client.raw.CancelSessionSubscriptionWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero SessionSubscription
		return zero, fmt.Errorf("cancel session subscription: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Reactivate withdraws the session's scheduled cancellation.
func (s *BillingSession) Reactivate(ctx context.Context) (SessionSubscription, error) {
	resp, err := s.client.raw.ReactivateSessionSubscriptionWithResponse(ctx)
	if err != nil {
		var zero SessionSubscription
		return zero, fmt.Errorf("reactivate session subscription: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// SetAddonQuantity attaches, changes or removes an add-on on the session's instance.
func (s *BillingSession) SetAddonQuantity(ctx context.Context, addonSlug string, input SessionAddonQuantity) (SessionAddon, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero SessionAddon
		return zero, fmt.Errorf("encode set session addon quantity request: %w", err)
	}

	resp, err := s.client.raw.SetSessionAddonQuantityWithBodyWithResponse(ctx, addonSlug, contentTypeJSON, body)
	if err != nil {
		var zero SessionAddon
		return zero, fmt.Errorf("set session addon quantity: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CreatePaymentMethodSession starts a hosted session where the customer sets a payment method.
func (s *BillingSession) CreatePaymentMethodSession(ctx context.Context, input NewSessionPaymentMethodSession) (PaymentMethodSession, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero PaymentMethodSession
		return zero, fmt.Errorf("encode create session payment method session request: %w", err)
	}

	resp, err := s.client.raw.CreateSessionPaymentMethodSessionWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero PaymentMethodSession
		return zero, fmt.Errorf("create session payment method session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CompletePaymentMethodSession reads back the payment method a session set.
func (s *BillingSession) CompletePaymentMethodSession(ctx context.Context, sessionID string) (CompletedPaymentMethodSession, error) {
	resp, err := s.client.raw.CompleteSessionPaymentMethodSessionWithResponse(ctx, sessionID)
	if err != nil {
		var zero CompletedPaymentMethodSession
		return zero, fmt.Errorf("complete session payment method session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// CreatePortalSession opens the billing provider's portal for the session's customer.
func (s *BillingSession) CreatePortalSession(ctx context.Context, input NewSessionPortalSession) (PortalSession, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero PortalSession
		return zero, fmt.Errorf("encode create session portal session request: %w", err)
	}

	resp, err := s.client.raw.CreateSessionPortalSessionWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero PortalSession
		return zero, fmt.Errorf("create session portal session: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// ValidateVoucher says whether a code could be redeemed on the session's instance.
func (s *BillingSession) ValidateVoucher(ctx context.Context, input SessionVoucherCheck) (SessionVoucherValidity, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero SessionVoucherValidity
		return zero, fmt.Errorf("encode validate session voucher request: %w", err)
	}

	resp, err := s.client.raw.ValidateSessionVoucherWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero SessionVoucherValidity
		return zero, fmt.Errorf("validate session voucher: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON429, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}

// Update changes a licence family.
func (s *LicenseFamilies) Update(ctx context.Context, familySlug string, input LicenseFamilyUpdate) (LicenseFamilyView, error) {
	body, err := jsonBody(input)
	if err != nil {
		var zero LicenseFamilyView
		return zero, fmt.Errorf("encode update license family request: %w", err)
	}

	resp, err := s.client.raw.UpdateLicenseFamilyWithBodyWithResponse(ctx, familySlug, contentTypeJSON, body)
	if err != nil {
		var zero LicenseFamilyView
		return zero, fmt.Errorf("update license family: %w", err)
	}

	apiResp := withRetryAfter(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}
