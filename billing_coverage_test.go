package sdk

import (
	"context"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// billingProblemCoverage drives every billing wrapper for
// TestClientWrappersReportEveryDeclaredProblem.
func billingProblemCoverage() map[string]problemCoverage[*Client] {
	return map[string]problemCoverage[*Client]{
		"Billing.Settings": {
			response: gen.GetBillingSettingsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Settings(ctx)
				return err
			},
		},

		"Billing.UpdateSettings": {
			response: gen.UpdateBillingSettingsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.UpdateSettings(ctx, BillingSettings{})
				return err
			},
		},

		"Billing.Capabilities": {
			response: gen.GetBillingCapabilitiesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Capabilities(ctx)
				return err
			},
		},

		"Billing.Health": {
			response: gen.GetBillingHealthResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Health(ctx)
				return err
			},
		},

		"Billing.Sync": {
			response: gen.SyncBillingProviderResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Sync(ctx)
				return err
			},
		},

		"Billing.ClosePeriods": {
			response: gen.CloseBillingPeriodsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.ClosePeriods(ctx, ClosePeriodsScope{})
				return err
			},
		},

		"BillingHandoff.ListPage": {
			response: gen.ListHandoffResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Handoff.ListPage(ctx, nil)
				return err
			},
		},

		"BillingHandoff.Claim": {
			response: gen.ClaimHandoffResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Handoff.Claim(ctx, HandoffClaimSize{})
				return err
			},
		},

		"BillingHandoff.Ack": {
			response: gen.AckHandoffResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Handoff.Ack(ctx, "x", HandoffBooking{})
				return err
			},
		},

		"BillingInvoices.ListPage": {
			response: gen.ListInvoicesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.ListPage(ctx, nil)
				return err
			},
		},

		"BillingInvoices.Get": {
			response: gen.GetInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.Get(ctx, "x")
				return err
			},
		},

		"BillingInvoices.Export": {
			response: gen.ExportInvoicesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.Export(ctx, nil)
				return err
			},
		},

		"BillingInvoices.Void": {
			response: gen.VoidInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.Void(ctx, "x", InvoiceVoid{})
				return err
			},
		},

		"BillingInvoices.Recompose": {
			response: gen.RecomposeInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.Recompose(ctx, "x")
				return err
			},
		},

		"BillingInvoices.MarkPaid": {
			response: gen.MarkInvoicePaidResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.MarkPaid(ctx, "x", InvoicePayment{})
				return err
			},
		},

		"BillingInvoices.WriteOff": {
			response: gen.WriteOffInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.WriteOff(ctx, "x", InvoiceWriteOff{})
				return err
			},
		},

		"BillingInvoices.ReleaseHold": {
			response: gen.ReleaseInvoiceHoldResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.ReleaseHold(ctx, "x", HoldRelease{})
				return err
			},
		},

		"BillingInvoices.RetryPush": {
			response: gen.RetryInvoicePushResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.RetryPush(ctx, "x")
				return err
			},
		},

		"BillingInvoices.Sync": {
			response: gen.SyncInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.Sync(ctx, "x")
				return err
			},
		},

		"BillingInvoices.ListLineReports": {
			response: gen.ListInvoiceLineReportsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Invoices.ListLineReports(ctx, "x", "x", nil)
				return err
			},
		},

		"BillingSubscriptions.Get": {
			response: gen.GetInstanceBillingResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.Get(ctx, "x")
				return err
			},
		},

		"BillingSubscriptions.Subscribe": {
			response: gen.SubscribeInstanceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.Subscribe(ctx, "x", NewSubscription{})
				return err
			},
		},

		"BillingSubscriptions.Update": {
			response: gen.UpdateInstanceBillingResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.Update(ctx, "x", SubscriptionTerms{})
				return err
			},
		},

		"BillingSubscriptions.Cancel": {
			response: gen.CancelSubscriptionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.Cancel(ctx, "x", SubscriptionCancellation{})
				return err
			},
		},

		"BillingSubscriptions.Reactivate": {
			response: gen.ReactivateSubscriptionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.Reactivate(ctx, "x")
				return err
			},
		},

		"BillingSubscriptions.SchedulePlanChange": {
			response: gen.SchedulePlanChangeResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.SchedulePlanChange(ctx, "x", PlanChangeTarget{})
				return err
			},
		},

		"BillingSubscriptions.CancelPlanChange": {
			response: gen.CancelPlanChangeResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.CancelPlanChange(ctx, "x")
				return err
			},
		},

		"BillingSubscriptions.UpcomingInvoice": {
			response: gen.GetUpcomingInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.UpcomingInvoice(ctx, "x")
				return err
			},
		},

		"BillingSubscriptions.ListInvoicesPage": {
			response: gen.ListInstanceInvoicesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Subscriptions.ListInvoicesPage(ctx, "x", nil)
				return err
			},
		},

		"BillingInstanceAddons.List": {
			response: gen.ListInstanceAddonsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.InstanceAddons.List(ctx, "x", nil)
				return err
			},
		},

		"BillingInstanceAddons.Attach": {
			response: gen.AttachInstanceAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.InstanceAddons.Attach(ctx, "x", NewInstanceAddon{})
				return err
			},
		},

		"BillingInstanceAddons.SetQuantity": {
			response: gen.SetInstanceAddonQuantityResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.InstanceAddons.SetQuantity(ctx, "x", "x", AddonQuantity{})
				return err
			},
		},

		"BillingInstanceAddons.Detach": {
			response: gen.DetachInstanceAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.InstanceAddons.Detach(ctx, "x", "x")
			},
		},

		"BillingInstanceVouchers.List": {
			response: gen.ListInstanceVouchersResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.InstanceVouchers.List(ctx, "x", nil)
				return err
			},
		},

		"BillingInstanceVouchers.Redeem": {
			response: gen.RedeemVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.InstanceVouchers.Redeem(ctx, "x", RedeemableCode{})
				return err
			},
		},

		"BillingInstanceVouchers.Revoke": {
			response: gen.RevokeInstanceVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.InstanceVouchers.Revoke(ctx, "x", "x", RevocationReason{})
				return err
			},
		},

		"BillingLicensePrices.List": {
			response: gen.ListLicensePricesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.LicensePrices.List(ctx, "x", nil)
				return err
			},
		},

		"BillingLicensePrices.Create": {
			response: gen.CreateLicensePriceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.LicensePrices.Create(ctx, "x", NewLicensePrice{})
				return err
			},
		},

		"BillingLicensePrices.Get": {
			response: gen.GetLicensePriceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.LicensePrices.Get(ctx, "x", "x")
				return err
			},
		},

		"BillingLicensePrices.Update": {
			response: gen.UpdateLicensePriceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.LicensePrices.Update(ctx, "x", "x", LicensePriceChanges{})
				return err
			},
		},

		"BillingLicensePrices.Deprecate": {
			response: gen.DeprecateLicensePriceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.LicensePrices.Deprecate(ctx, "x", "x")
				return err
			},
		},

		"BillingLicensePrices.Preview": {
			response: gen.PreviewLicenseInvoiceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.LicensePrices.Preview(ctx, "x", InvoicePreviewScenario{})
				return err
			},
		},

		"BillingAddons.ListPage": {
			response: gen.ListAddonsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.ListPage(ctx, nil)
				return err
			},
		},

		"BillingAddons.Create": {
			response: gen.CreateAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.Create(ctx, NewAddon{})
				return err
			},
		},

		"BillingAddons.Get": {
			response: gen.GetAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.Get(ctx, "x")
				return err
			},
		},

		"BillingAddons.Update": {
			response: gen.UpdateAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.Update(ctx, "x", AddonChanges{})
				return err
			},
		},

		"BillingAddons.Delete": {
			response: gen.DeleteAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.Addons.Delete(ctx, "x")
			},
		},

		"BillingAddons.Publish": {
			response: gen.PublishAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.Publish(ctx, "x")
				return err
			},
		},

		"BillingAddons.Archive": {
			response: gen.ArchiveAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.Archive(ctx, "x")
				return err
			},
		},

		"BillingAddons.Unarchive": {
			response: gen.UnarchiveAddonResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.Unarchive(ctx, "x")
				return err
			},
		},

		"BillingAddons.ListCompatibility": {
			response: gen.ListAddonCompatibilityResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.ListCompatibility(ctx, "x")
				return err
			},
		},

		"BillingAddons.SetCompatibility": {
			response: gen.SetAddonCompatibilityResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.Addons.SetCompatibility(ctx, "x", "x")
			},
		},

		"BillingAddons.RemoveCompatibility": {
			response: gen.RemoveAddonCompatibilityResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.Addons.RemoveCompatibility(ctx, "x", "x")
			},
		},

		"BillingAddons.ListEntitlements": {
			response: gen.ListAddonEntitlementsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.ListEntitlements(ctx, "x")
				return err
			},
		},

		"BillingAddons.AssignEntitlement": {
			response: gen.AssignAddonEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.AssignEntitlement(ctx, "x", NewAddonGrant{})
				return err
			},
		},

		"BillingAddons.GetEntitlement": {
			response: gen.GetAddonEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.GetEntitlement(ctx, "x", "x")
				return err
			},
		},

		"BillingAddons.UpdateEntitlement": {
			response: gen.UpdateAddonEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.UpdateEntitlement(ctx, "x", "x", AddonGrant{})
				return err
			},
		},

		"BillingAddons.UnassignEntitlement": {
			response: gen.UnassignAddonEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.Addons.UnassignEntitlement(ctx, "x", "x")
			},
		},

		"BillingAddons.ListPrices": {
			response: gen.ListAddonPricesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.ListPrices(ctx, "x", nil)
				return err
			},
		},

		"BillingAddons.CreatePrice": {
			response: gen.CreateAddonPriceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.CreatePrice(ctx, "x", NewAddonPrice{})
				return err
			},
		},

		"BillingAddons.DeprecatePrice": {
			response: gen.DeprecateAddonPriceResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.DeprecatePrice(ctx, "x", "x")
				return err
			},
		},

		"BillingAddons.ListFamilies": {
			response: gen.ListAddonFamiliesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.ListFamilies(ctx)
				return err
			},
		},

		"BillingAddons.GetFamily": {
			response: gen.GetAddonFamilyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.GetFamily(ctx, "x")
				return err
			},
		},

		"BillingAddons.UpdateFamily": {
			response: gen.UpdateAddonFamilyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Addons.UpdateFamily(ctx, "x", AddonFamilyVisibility{})
				return err
			},
		},

		"BillingVouchers.ListPage": {
			response: gen.ListVouchersResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.ListPage(ctx, nil)
				return err
			},
		},

		"BillingVouchers.Create": {
			response: gen.CreateVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Create(ctx, VoucherDraft{})
				return err
			},
		},

		"BillingVouchers.Get": {
			response: gen.GetVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Get(ctx, "x")
				return err
			},
		},

		"BillingVouchers.Update": {
			response: gen.UpdateVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Update(ctx, "x", VoucherDraft{})
				return err
			},
		},

		"BillingVouchers.Publish": {
			response: gen.PublishVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Publish(ctx, "x")
				return err
			},
		},

		"BillingVouchers.Archive": {
			response: gen.ArchiveVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Archive(ctx, "x")
				return err
			},
		},

		"BillingVouchers.ListRedemptionsPage": {
			response: gen.ListVoucherRedemptionsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.ListRedemptionsPage(ctx, "x", nil)
				return err
			},
		},

		"BillingVouchers.Lookup": {
			response: gen.LookupVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Lookup(ctx, VoucherCode{})
				return err
			},
		},

		"BillingVouchers.Validate": {
			response: gen.ValidateVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Vouchers.Validate(ctx, VoucherCheck{})
				return err
			},
		},

		"BillingPublishableKeys.ListPage": {
			response: gen.ListPublishableKeysResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.PublishableKeys.ListPage(ctx, nil)
				return err
			},
		},

		"BillingPublishableKeys.Create": {
			response: gen.CreatePublishableKeyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.PublishableKeys.Create(ctx, PublishableKeyDraft{})
				return err
			},
		},

		"BillingPublishableKeys.Update": {
			response: gen.UpdatePublishableKeyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.PublishableKeys.Update(ctx, "x", PublishableKeyPatch{})
				return err
			},
		},

		"BillingPublishableKeys.Revoke": {
			response: gen.RevokePublishableKeyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.PublishableKeys.Revoke(ctx, "x")
				return err
			},
		},

		"BillingCustomerSessions.Create": {
			response: gen.CreateCustomerSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.CustomerSessions.Create(ctx, CustomerSessionDraft{})
				return err
			},
		},

		"BillingCustomerSessions.Revoke": {
			response: gen.RevokeCustomerSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.CustomerSessions.Revoke(ctx, "x")
			},
		},

		"BillingCustomers.Get": {
			response: gen.GetCustomerBillingResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Customers.Get(ctx, "x")
				return err
			},
		},

		"BillingCustomers.DetachPaymentMethod": {
			response: gen.DetachPaymentMethodResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Billing.Customers.DetachPaymentMethod(ctx, "x")
			},
		},

		"BillingCustomers.CreatePaymentMethodSession": {
			response: gen.CreatePaymentMethodSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Customers.CreatePaymentMethodSession(ctx, "x", NewPaymentMethodSession{})
				return err
			},
		},

		"BillingCustomers.CompletePaymentMethodSession": {
			response: gen.CompletePaymentMethodSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Customers.CompletePaymentMethodSession(ctx, "x", "x")
				return err
			},
		},

		"BillingCustomers.CreatePortalSession": {
			response: gen.CreatePortalSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Customers.CreatePortalSession(ctx, "x", NewPortalSession{})
				return err
			},
		},

		"BillingPublic.Catalog": {
			response: gen.GetPublicCatalogResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Public.Catalog(ctx, nil)
				return err
			},
		},

		"BillingSession.Portal": {
			response: gen.GetSessionPortalResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.Portal(ctx)
				return err
			},
		},

		"BillingSession.ListInvoicesPage": {
			response: gen.ListSessionInvoicesResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.ListInvoicesPage(ctx, nil)
				return err
			},
		},

		"BillingSession.Checkout": {
			response: gen.CreateSessionCheckoutResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.Checkout(ctx, SessionCheckoutOrder{})
				return err
			},
		},

		"BillingSession.Cancel": {
			response: gen.CancelSessionSubscriptionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.Cancel(ctx, SessionCancellation{})
				return err
			},
		},

		"BillingSession.Reactivate": {
			response: gen.ReactivateSessionSubscriptionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.Reactivate(ctx)
				return err
			},
		},

		"BillingSession.SetAddonQuantity": {
			response: gen.SetSessionAddonQuantityResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.SetAddonQuantity(ctx, "x", SessionAddonQuantity{})
				return err
			},
		},

		"BillingSession.CreatePaymentMethodSession": {
			response: gen.CreateSessionPaymentMethodSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.CreatePaymentMethodSession(ctx, NewSessionPaymentMethodSession{})
				return err
			},
		},

		"BillingSession.CompletePaymentMethodSession": {
			response: gen.CompleteSessionPaymentMethodSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.CompletePaymentMethodSession(ctx, "x")
				return err
			},
		},

		"BillingSession.CreatePortalSession": {
			response: gen.CreateSessionPortalSessionResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.CreatePortalSession(ctx, NewSessionPortalSession{})
				return err
			},
		},

		"BillingSession.ValidateVoucher": {
			response: gen.ValidateSessionVoucherResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Billing.Session.ValidateVoucher(ctx, SessionVoucherCheck{})
				return err
			},
		},

		"LicenseFamilies.Update": {
			response: gen.UpdateLicenseFamilyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.LicenseFamilies.Update(ctx, "x", LicenseFamilyVisibility{})
				return err
			},
		},
	}
}

// billingServiceCoverage is the reviewed surface of the billing namespaces, for
// TestClientExposesOnlyTheReviewedPublicSurface.
func billingServiceCoverage(client *Client) map[string]serviceCoverage {
	return map[string]serviceCoverage{
		"Billing": {
			service: client.Billing,
			methods: []string{"Settings", "UpdateSettings", "Capabilities", "Health", "Sync", "ClosePeriods"},
		},
		"BillingHandoff": {
			service: client.Billing.Handoff,
			methods: []string{"ListPage", "List", "Claim", "Ack"},
		},
		"BillingInvoices": {
			service: client.Billing.Invoices,
			methods: []string{"ListPage", "List", "Get", "Export", "Void", "Recompose", "MarkPaid", "WriteOff", "ReleaseHold", "RetryPush", "Sync", "ListLineReports"},
		},
		"BillingSubscriptions": {
			service: client.Billing.Subscriptions,
			methods: []string{"Get", "Subscribe", "Update", "Cancel", "Reactivate", "SchedulePlanChange", "CancelPlanChange", "UpcomingInvoice", "ListInvoicesPage", "ListInvoices"},
		},
		"BillingInstanceAddons": {
			service: client.Billing.InstanceAddons,
			methods: []string{"List", "Attach", "SetQuantity", "Detach"},
		},
		"BillingInstanceVouchers": {
			service: client.Billing.InstanceVouchers,
			methods: []string{"List", "Redeem", "Revoke"},
		},
		"BillingLicensePrices": {
			service: client.Billing.LicensePrices,
			methods: []string{"List", "Create", "Get", "Update", "Deprecate", "Preview"},
		},
		"BillingAddons": {
			service: client.Billing.Addons,
			methods: []string{"ListPage", "List", "Create", "Get", "Update", "Delete", "Publish", "Archive", "Unarchive", "ListCompatibility", "SetCompatibility", "RemoveCompatibility", "ListEntitlements", "AssignEntitlement", "GetEntitlement", "UpdateEntitlement", "UnassignEntitlement", "ListPrices", "CreatePrice", "DeprecatePrice", "ListFamilies", "GetFamily", "UpdateFamily"},
		},
		"BillingVouchers": {
			service: client.Billing.Vouchers,
			methods: []string{"ListPage", "List", "Create", "Get", "Update", "Publish", "Archive", "ListRedemptionsPage", "ListRedemptions", "Lookup", "Validate"},
		},
		"BillingPublishableKeys": {
			service: client.Billing.PublishableKeys,
			methods: []string{"ListPage", "List", "Create", "Update", "Revoke"},
		},
		"BillingCustomerSessions": {
			service: client.Billing.CustomerSessions,
			methods: []string{"Create", "Revoke"},
		},
		"BillingCustomers": {
			service: client.Billing.Customers,
			methods: []string{"Get", "DetachPaymentMethod", "CreatePaymentMethodSession", "CompletePaymentMethodSession", "CreatePortalSession"},
		},
		"BillingPublic": {
			service: client.Billing.Public,
			methods: []string{"Catalog"},
		},
		"BillingSession": {
			service: client.Billing.Session,
			methods: []string{"Portal", "ListInvoicesPage", "ListInvoices", "Checkout", "Cancel", "Reactivate", "SetAddonQuantity", "CreatePaymentMethodSession", "CompletePaymentMethodSession", "CreatePortalSession", "ValidateVoucher"},
		},
	}
}
