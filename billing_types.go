package sdk

import "github.com/kaitencloud/sdk-go/internal/gen"

// Addon is the Core API's Addon.
type Addon = gen.Addon

// AddonAmount is the Core API's AddonAmount.
type AddonAmount = gen.AddonAmount

// AddonChanges is the Core API's AddonChanges.
type AddonChanges = gen.AddonChanges

// AddonEntitlement is the Core API's AddonEntitlement.
type AddonEntitlement = gen.AddonEntitlement

// AddonEntitlementEntitlementType is the Core API's AddonEntitlementEntitlementType.
type AddonEntitlementEntitlementType = gen.AddonEntitlementEntitlementType

// AddonEntitlementOverrideBehavior How a number grant combines with the licence's: ADD sums value × quantity, OVERRIDE replaces it with the latest attached add-on's, MAX keeps the larger. A boolean always ORs, a config always overrides.
type AddonEntitlementOverrideBehavior = gen.AddonEntitlementOverrideBehavior

// AddonFamily is the Core API's AddonFamily.
type AddonFamily = gen.AddonFamily

// AddonFamilyVisibility is the Core API's AddonFamilyVisibility.
type AddonFamilyVisibility = gen.AddonFamilyVisibility

// AddonGrant is the Core API's AddonGrant.
type AddonGrant = gen.AddonGrant

// AddonGrantOverrideBehavior How a number grant combines with the licence's
type AddonGrantOverrideBehavior = gen.AddonGrantOverrideBehavior

// AddonLifecycleState A DRAFT can be attached to try it; an ARCHIVED version stays on the instances holding it but cannot be attached again.
type AddonLifecycleState = gen.AddonLifecycleState

// AddonPricingType is the Core API's AddonPricingType.
type AddonPricingType = gen.AddonPricingType

// AddonQuantity is the Core API's AddonQuantity.
type AddonQuantity = gen.AddonQuantity

// BillingCapabilities is the Core API's BillingCapabilities.
type BillingCapabilities = gen.BillingCapabilities

// BillingCapabilitiesDisabledReason Why they do not: billing is off on this deployment, or the organization's plan does not include it
type BillingCapabilitiesDisabledReason = gen.BillingCapabilitiesDisabledReason

// BillingFeatures is the Core API's BillingFeatures.
type BillingFeatures = gen.BillingFeatures

// BillingHealth is the Core API's BillingHealth.
type BillingHealth = gen.BillingHealth

// BillingProvider is the Core API's BillingProvider.
type BillingProvider = gen.BillingProvider

// BillingProviderKind is the Core API's BillingProviderKind.
type BillingProviderKind = gen.BillingProviderKind

// BillingProviderUnavailableReason Why it is not available; absent when it is
type BillingProviderUnavailableReason = gen.BillingProviderUnavailableReason

// BillingSettings is the Core API's BillingSettings.
type BillingSettings = gen.BillingSettings

// BillingSettingsDefaultCollectionMethod How an invoice is collected when its subscription does not say: SEND_INVOICE sends it for payment; CHARGE_AUTOMATICALLY charges the customer's payment method, on a provider that can. A NOOP subscription sends its invoices whatever the default.
//
// Examples: SEND_INVOICE
type BillingSettingsDefaultCollectionMethod = gen.BillingSettingsDefaultCollectionMethod

// CanceledSubscription is the Core API's CanceledSubscription.
type CanceledSubscription = gen.CanceledSubscription

// CanceledSubscriptionBillingPeriod is the Core API's CanceledSubscriptionBillingPeriod.
type CanceledSubscriptionBillingPeriod = gen.CanceledSubscriptionBillingPeriod

// CanceledSubscriptionCollectionMethod Effective: the subscription's own, else the organization's default
type CanceledSubscriptionCollectionMethod = gen.CanceledSubscriptionCollectionMethod

// CanceledSubscriptionProviderKind Who collects its invoices. NOOP: the organization itself, through the handoff queue
type CanceledSubscriptionProviderKind = gen.CanceledSubscriptionProviderKind

// CanceledSubscriptionStatus is the Core API's CanceledSubscriptionStatus.
type CanceledSubscriptionStatus = gen.CanceledSubscriptionStatus

// Cancellation is the Core API's Cancellation.
type Cancellation = gen.Cancellation

// CancellationBillingPeriod is the Core API's CancellationBillingPeriod.
type CancellationBillingPeriod = gen.CancellationBillingPeriod

// CancellationCollectionMethod Effective: the subscription's own, else the organization's default
type CancellationCollectionMethod = gen.CancellationCollectionMethod

// CancellationMode How it was canceled; a trial is canceled at once in either mode
type CancellationMode = gen.CancellationMode

// CancellationProviderKind Who collects its invoices. NOOP: the organization itself, through the handoff queue
type CancellationProviderKind = gen.CancellationProviderKind

// CancellationStatus is the Core API's CancellationStatus.
type CancellationStatus = gen.CancellationStatus

// CloseBacklog is the Core API's CloseBacklog.
type CloseBacklog = gen.CloseBacklog

// ClosePeriodsReport is the Core API's ClosePeriodsReport.
type ClosePeriodsReport = gen.ClosePeriodsReport

// ClosePeriodsScope is the Core API's ClosePeriodsScope.
type ClosePeriodsScope = gen.ClosePeriodsScope

// ClosedInvoice is the Core API's ClosedInvoice.
type ClosedInvoice = gen.ClosedInvoice

// ClosedInvoiceKind is the Core API's ClosedInvoiceKind.
type ClosedInvoiceKind = gen.ClosedInvoiceKind

// ClosedInvoiceStatus is the Core API's ClosedInvoiceStatus.
type ClosedInvoiceStatus = gen.ClosedInvoiceStatus

// CompatibleLicenseFamilies is the Core API's CompatibleLicenseFamilies.
type CompatibleLicenseFamilies = gen.CompatibleLicenseFamilies

// CompletedPaymentMethodSession is the Core API's CompletedPaymentMethodSession.
type CompletedPaymentMethodSession = gen.CompletedPaymentMethodSession

// CreatedCustomerSession is the Core API's CreatedCustomerSession.
type CreatedCustomerSession = gen.CreatedCustomerSession

// CustomerBilling is the Core API's CustomerBilling.
type CustomerBilling = gen.CustomerBilling

// CustomerInProvider is the Core API's CustomerInProvider.
type CustomerInProvider = gen.CustomerInProvider

// CustomerInProviderProviderKind is the Core API's CustomerInProviderProviderKind.
type CustomerInProviderProviderKind = gen.CustomerInProviderProviderKind

// CustomerIntegration is the Core API's CustomerIntegration.
type CustomerIntegration = gen.CustomerIntegration

// CustomerSessionDraft is the Core API's CustomerSessionDraft.
type CustomerSessionDraft = gen.CustomerSessionDraft

// DiscountAllocation is the Core API's DiscountAllocation.
type DiscountAllocation = gen.DiscountAllocation

// DiscountDifference is the Core API's DiscountDifference.
type DiscountDifference = gen.DiscountDifference

// DueDateDifference is the Core API's DueDateDifference.
type DueDateDifference = gen.DueDateDifference

// DueToday is the Core API's DueToday.
type DueToday = gen.DueToday

// ExportInvoicesParams defines parameters for ExportInvoices.
type ExportInvoicesParams = gen.ExportInvoicesParams

// ExportInvoicesParamsHandoffStatus defines parameters for ExportInvoices.
type ExportInvoicesParamsHandoffStatus = gen.ExportInvoicesParamsHandoffStatus

// ExportInvoicesParamsKind defines parameters for ExportInvoices.
type ExportInvoicesParamsKind = gen.ExportInvoicesParamsKind

// ExportInvoicesParamsProviderKind defines parameters for ExportInvoices.
type ExportInvoicesParamsProviderKind = gen.ExportInvoicesParamsProviderKind

// ExportInvoicesParamsStatus defines parameters for ExportInvoices.
type ExportInvoicesParamsStatus = gen.ExportInvoicesParamsStatus

// ExtraDiscount is the Core API's ExtraDiscount.
type ExtraDiscount = gen.ExtraDiscount

// GetPublicCatalogParams defines parameters for GetPublicCatalog.
type GetPublicCatalogParams = gen.GetPublicCatalogParams

// GetPublicCatalogParamsInclude defines parameters for GetPublicCatalog.
type GetPublicCatalogParamsInclude = gen.GetPublicCatalogParamsInclude

// Grant is the Core API's Grant.
type Grant = gen.Grant

// GrantModifierType SET replaces the value, ADD adds to it, MULTIPLY multiplies it, UNLIMITED lifts the limit
type GrantModifierType = gen.GrantModifierType

// HandoffBacklog is the Core API's HandoffBacklog.
type HandoffBacklog = gen.HandoffBacklog

// HandoffBooking is the Core API's HandoffBooking.
type HandoffBooking = gen.HandoffBooking

// HandoffClaim is the Core API's HandoffClaim.
type HandoffClaim = gen.HandoffClaim

// HandoffClaimSize is the Core API's HandoffClaimSize.
type HandoffClaimSize = gen.HandoffClaimSize

// HeldInvoices is the Core API's HeldInvoices.
type HeldInvoices = gen.HeldInvoices

// HeldPair is the Core API's HeldPair.
type HeldPair = gen.HeldPair

// HoldDetail is the Core API's HoldDetail.
type HoldDetail = gen.HoldDetail

// HoldRecord is the Core API's HoldRecord.
type HoldRecord = gen.HoldRecord

// HoldRelease is the Core API's HoldRelease.
type HoldRelease = gen.HoldRelease

// InstanceAddon is the Core API's InstanceAddon.
type InstanceAddon = gen.InstanceAddon

// InstanceBilling is the Core API's InstanceBilling.
type InstanceBilling = gen.InstanceBilling

// InstanceBillingBillingPeriod is the Core API's InstanceBillingBillingPeriod.
type InstanceBillingBillingPeriod = gen.InstanceBillingBillingPeriod

// InstanceBillingCollectionMethod Effective: the subscription's own, else the organization's default
type InstanceBillingCollectionMethod = gen.InstanceBillingCollectionMethod

// InstanceBillingProviderKind Who collects its invoices. NOOP: the organization itself, through the handoff queue
type InstanceBillingProviderKind = gen.InstanceBillingProviderKind

// InstanceBillingStatus is the Core API's InstanceBillingStatus.
type InstanceBillingStatus = gen.InstanceBillingStatus

// InstanceIntegration is the Core API's InstanceIntegration.
type InstanceIntegration = gen.InstanceIntegration

// Invoice is the Core API's Invoice.
type Invoice = gen.Invoice

// InvoiceCollectionMethod is the Core API's InvoiceCollectionMethod.
type InvoiceCollectionMethod = gen.InvoiceCollectionMethod

// InvoiceHandoff is the Core API's InvoiceHandoff.
type InvoiceHandoff = gen.InvoiceHandoff

// InvoiceHandoffState is the Core API's InvoiceHandoffState.
type InvoiceHandoffState = gen.InvoiceHandoffState

// InvoiceHandoffStatus PENDING: waiting in the handoff queue for the organization's accounting system
type InvoiceHandoffStatus = gen.InvoiceHandoffStatus

// InvoiceHold is the Core API's InvoiceHold.
type InvoiceHold = gen.InvoiceHold

// InvoiceHoldReason Why a DRAFT is held: the usage journal it was measured from failed a consistency check
type InvoiceHoldReason = gen.InvoiceHoldReason

// InvoiceKind ACTIVATION bills a subscription's first period, RENEWAL the period that ended and the one that starts, FINAL what is left when a subscription ends
type InvoiceKind = gen.InvoiceKind

// InvoiceLine is the Core API's InvoiceLine.
type InvoiceLine = gen.InvoiceLine

// InvoiceLineBillingModel The price's billing model; absent on a DISCOUNT line
type InvoiceLineBillingModel = gen.InvoiceLineBillingModel

// InvoiceLineBillingTiming ADVANCE lines bill the period that starts at the boundary, ARREARS lines the one that ends there; absent on a DISCOUNT line
type InvoiceLineBillingTiming = gen.InvoiceLineBillingTiming

// InvoiceLineDiscount is the Core API's InvoiceLineDiscount.
type InvoiceLineDiscount = gen.InvoiceLineDiscount

// InvoiceLineDiscountAppliesTo is the Core API's InvoiceLineDiscountAppliesTo.
type InvoiceLineDiscountAppliesTo = gen.InvoiceLineDiscountAppliesTo

// InvoiceLineDiscountDiscountType is the Core API's InvoiceLineDiscountDiscountType.
type InvoiceLineDiscountDiscountType = gen.InvoiceLineDiscountDiscountType

// InvoiceLineLedger is the Core API's InvoiceLineLedger.
type InvoiceLineLedger = gen.InvoiceLineLedger

// InvoiceLineMetering is the Core API's InvoiceLineMetering.
type InvoiceLineMetering = gen.InvoiceLineMetering

// InvoiceLineOverage is the Core API's InvoiceLineOverage.
type InvoiceLineOverage = gen.InvoiceLineOverage

// InvoiceLineProvider is the Core API's InvoiceLineProvider.
type InvoiceLineProvider = gen.InvoiceLineProvider

// InvoiceLineType BASE: the subscription's FLAT_FEE price. ADDON: an add-on's FLAT_FEE price times the quantity held at the boundary. USAGE: a USAGE_BASED price's metered usage. OVERAGE: an OVERAGE price's usage above the licence's limit. DISCOUNT: a PRICE voucher's discount, negative.
type InvoiceLineType = gen.InvoiceLineType

// InvoicePayment is the Core API's InvoicePayment.
type InvoicePayment = gen.InvoicePayment

// InvoicePreview is the Core API's InvoicePreview.
type InvoicePreview = gen.InvoicePreview

// InvoicePreviewKind The boundary previewed
type InvoicePreviewKind = gen.InvoicePreviewKind

// InvoicePreviewScenario is the Core API's InvoicePreviewScenario.
type InvoicePreviewScenario = gen.InvoicePreviewScenario

// InvoicePreviewStatus Always PREVIEW
type InvoicePreviewStatus = gen.InvoicePreviewStatus

// InvoiceProviderKind Who collects it. NOOP: the organization itself, through the handoff queue
type InvoiceProviderKind = gen.InvoiceProviderKind

// InvoiceStatus MANUAL: issued, for the organization to collect. PAID, UNCOLLECTIBLE and VOID are final. DRAFT: held, waiting for a release or a recompose
type InvoiceStatus = gen.InvoiceStatus

// InvoiceSummary is the Core API's InvoiceSummary.
type InvoiceSummary = gen.InvoiceSummary

// InvoiceSummaryCollectionMethod is the Core API's InvoiceSummaryCollectionMethod.
type InvoiceSummaryCollectionMethod = gen.InvoiceSummaryCollectionMethod

// InvoiceSummaryHandoffStatus PENDING: waiting in the handoff queue for the organization's accounting system
type InvoiceSummaryHandoffStatus = gen.InvoiceSummaryHandoffStatus

// InvoiceSummaryHoldReason Why a DRAFT is held: the usage journal it was measured from failed a consistency check
type InvoiceSummaryHoldReason = gen.InvoiceSummaryHoldReason

// InvoiceSummaryKind ACTIVATION bills a subscription's first period, RENEWAL the period that ended and the one that starts, FINAL what is left when a subscription ends
type InvoiceSummaryKind = gen.InvoiceSummaryKind

// InvoiceSummaryProviderKind Who collects it. NOOP: the organization itself, through the handoff queue
type InvoiceSummaryProviderKind = gen.InvoiceSummaryProviderKind

// InvoiceSummaryStatus MANUAL: issued, for the organization to collect. PAID, UNCOLLECTIBLE and VOID are final. DRAFT: held, waiting for a release or a recompose
type InvoiceSummaryStatus = gen.InvoiceSummaryStatus

// InvoiceVoid is the Core API's InvoiceVoid.
type InvoiceVoid = gen.InvoiceVoid

// InvoiceWriteOff is the Core API's InvoiceWriteOff.
type InvoiceWriteOff = gen.InvoiceWriteOff

// LicenseFamilyVisibility is the Core API's LicenseFamilyVisibility.
type LicenseFamilyVisibility = gen.LicenseFamilyVisibility

// LicensePriceChanges is the Core API's LicensePriceChanges.
type LicensePriceChanges = gen.LicensePriceChanges

// LicensePriceChangesBillingPeriod See createLicensePrice.
type LicensePriceChangesBillingPeriod = gen.LicensePriceChangesBillingPeriod

// LicensePriceChangesBillingTiming See createLicensePrice.
type LicensePriceChangesBillingTiming = gen.LicensePriceChangesBillingTiming

// LicensePricingType How this version is sold. FREE and PAID can be bought self-serve; CUSTOM sends a buyer to selfServeCtaUrl or a conversation. Always present in responses. On create it defaults to CUSTOM; on update, omit it to keep the stored value.
//
// Examples: PAID
type LicensePricingType = gen.LicensePricingType

// LineDifference is the Core API's LineDifference.
type LineDifference = gen.LineDifference

// LineReportPage is the Core API's LineReportPage.
type LineReportPage = gen.LineReportPage

// List is the Core API's List.
type List = gen.List

// ListAddonPricesParams defines parameters for ListAddonPrices.
type ListAddonPricesParams = gen.ListAddonPricesParams

// ListAddonPricesParamsStatus defines parameters for ListAddonPrices.
type ListAddonPricesParamsStatus = gen.ListAddonPricesParamsStatus

// ListAddonsParams defines parameters for ListAddons.
type ListAddonsParams = gen.ListAddonsParams

// ListAddonsParamsLifecycleState defines parameters for ListAddons.
type ListAddonsParamsLifecycleState = gen.ListAddonsParamsLifecycleState

// ListHandoffParams defines parameters for ListHandoff.
type ListHandoffParams = gen.ListHandoffParams

// ListHandoffParamsStatus defines parameters for ListHandoff.
type ListHandoffParamsStatus = gen.ListHandoffParamsStatus

// ListInstanceAddonsParams defines parameters for ListInstanceAddons.
type ListInstanceAddonsParams = gen.ListInstanceAddonsParams

// ListInstanceInvoicesParams defines parameters for ListInstanceInvoices.
type ListInstanceInvoicesParams = gen.ListInstanceInvoicesParams

// ListInstanceInvoicesParamsHandoffStatus defines parameters for ListInstanceInvoices.
type ListInstanceInvoicesParamsHandoffStatus = gen.ListInstanceInvoicesParamsHandoffStatus

// ListInstanceInvoicesParamsKind defines parameters for ListInstanceInvoices.
type ListInstanceInvoicesParamsKind = gen.ListInstanceInvoicesParamsKind

// ListInstanceInvoicesParamsProviderKind defines parameters for ListInstanceInvoices.
type ListInstanceInvoicesParamsProviderKind = gen.ListInstanceInvoicesParamsProviderKind

// ListInstanceInvoicesParamsStatus defines parameters for ListInstanceInvoices.
type ListInstanceInvoicesParamsStatus = gen.ListInstanceInvoicesParamsStatus

// ListInstanceVouchersParams defines parameters for ListInstanceVouchers.
type ListInstanceVouchersParams = gen.ListInstanceVouchersParams

// ListInstanceVouchersParamsStatus defines parameters for ListInstanceVouchers.
type ListInstanceVouchersParamsStatus = gen.ListInstanceVouchersParamsStatus

// ListInvoiceLineReportsParams defines parameters for ListInvoiceLineReports.
type ListInvoiceLineReportsParams = gen.ListInvoiceLineReportsParams

// ListInvoicesParams defines parameters for ListInvoices.
type ListInvoicesParams = gen.ListInvoicesParams

// ListInvoicesParamsHandoffStatus defines parameters for ListInvoices.
type ListInvoicesParamsHandoffStatus = gen.ListInvoicesParamsHandoffStatus

// ListInvoicesParamsKind defines parameters for ListInvoices.
type ListInvoicesParamsKind = gen.ListInvoicesParamsKind

// ListInvoicesParamsProviderKind defines parameters for ListInvoices.
type ListInvoicesParamsProviderKind = gen.ListInvoicesParamsProviderKind

// ListInvoicesParamsStatus defines parameters for ListInvoices.
type ListInvoicesParamsStatus = gen.ListInvoicesParamsStatus

// ListLicensePricesParams defines parameters for ListLicensePrices.
type ListLicensePricesParams = gen.ListLicensePricesParams

// ListLicensePricesParamsBillingModel defines parameters for ListLicensePrices.
type ListLicensePricesParamsBillingModel = gen.ListLicensePricesParamsBillingModel

// ListLicensePricesParamsStatus defines parameters for ListLicensePrices.
type ListLicensePricesParamsStatus = gen.ListLicensePricesParamsStatus

// ListPublishableKeysParams defines parameters for ListPublishableKeys.
type ListPublishableKeysParams = gen.ListPublishableKeysParams

// ListSessionInvoicesParams defines parameters for ListSessionInvoices.
type ListSessionInvoicesParams = gen.ListSessionInvoicesParams

// ListVoucherRedemptionsParams defines parameters for ListVoucherRedemptions.
type ListVoucherRedemptionsParams = gen.ListVoucherRedemptionsParams

// ListVouchersParams defines parameters for ListVouchers.
type ListVouchersParams = gen.ListVouchersParams

// ListVouchersParamsStatus defines parameters for ListVouchers.
type ListVouchersParamsStatus = gen.ListVouchersParamsStatus

// ListVouchersParamsVoucherType defines parameters for ListVouchers.
type ListVouchersParamsVoucherType = gen.ListVouchersParamsVoucherType

// MinimumAmount is the Core API's MinimumAmount.
type MinimumAmount = gen.MinimumAmount

// NewAddon is the Core API's NewAddon.
type NewAddon = gen.NewAddon

// NewAddonGrant is the Core API's NewAddonGrant.
type NewAddonGrant = gen.NewAddonGrant

// NewAddonGrantOverrideBehavior How a number grant combines with the licence's
type NewAddonGrantOverrideBehavior = gen.NewAddonGrantOverrideBehavior

// NewAddonLifecycleState is the Core API's NewAddonLifecycleState.
type NewAddonLifecycleState = gen.NewAddonLifecycleState

// NewAddonPrice is the Core API's NewAddonPrice.
type NewAddonPrice = gen.NewAddonPrice

// NewAddonPriceBillingModel is the Core API's NewAddonPriceBillingModel.
type NewAddonPriceBillingModel = gen.NewAddonPriceBillingModel

// NewAddonPriceBillingPeriod Required on a FLAT_FEE price: it bills subscriptions of that period. Refused on a metered one.
type NewAddonPriceBillingPeriod = gen.NewAddonPriceBillingPeriod

// NewAddonPriceBillingTiming Defaults to ADVANCE for a FLAT_FEE price and ARREARS for a metered one, which must be ARREARS.
type NewAddonPriceBillingTiming = gen.NewAddonPriceBillingTiming

// NewAddonPricingType is the Core API's NewAddonPricingType.
type NewAddonPricingType = gen.NewAddonPricingType

// NewInstanceAddon is the Core API's NewInstanceAddon.
type NewInstanceAddon = gen.NewInstanceAddon

// NewLicensePrice is the Core API's NewLicensePrice.
type NewLicensePrice = gen.NewLicensePrice

// NewLicensePriceBillingModel FLAT_FEE: an amount per billing period. USAGE_BASED: per sale unit of usage. OVERAGE: per sale unit of usage above the licence's limit.
type NewLicensePriceBillingModel = gen.NewLicensePriceBillingModel

// NewLicensePriceBillingPeriod Required on a FLAT_FEE price, refused on a metered one (CreateLicensePrice.InvalidPeriod).
type NewLicensePriceBillingPeriod = gen.NewLicensePriceBillingPeriod

// NewLicensePriceBillingTiming Defaults to ADVANCE for a FLAT_FEE price and ARREARS for a metered one, which must be ARREARS (CreateLicensePrice.InvalidTiming).
type NewLicensePriceBillingTiming = gen.NewLicensePriceBillingTiming

// NewPaymentMethodSession is the Core API's NewPaymentMethodSession.
type NewPaymentMethodSession = gen.NewPaymentMethodSession

// NewPortalSession is the Core API's NewPortalSession.
type NewPortalSession = gen.NewPortalSession

// NewSessionPaymentMethodSession is the Core API's NewSessionPaymentMethodSession.
type NewSessionPaymentMethodSession = gen.NewSessionPaymentMethodSession

// NewSessionPortalSession is the Core API's NewSessionPortalSession.
type NewSessionPortalSession = gen.NewSessionPortalSession

// NewSubscription is the Core API's NewSubscription.
type NewSubscription = gen.NewSubscription

// NewSubscriptionCollectionMethod Omitted: the organization's default, read at each invoice. CHARGE_AUTOMATICALLY needs a provider that charges automatically (SubscribeInstance.CollectionMethodUnsupported).
type NewSubscriptionCollectionMethod = gen.NewSubscriptionCollectionMethod

// NewSubscriptionProviderKind Who issues and collects the invoices. NOOP (the default): the organization itself, through the handoff queue. A provider must be connected (SubscribeInstance.ProviderNotConnected); one that pushes invoices needs the customer's billing e-mail for SEND_INVOICE (.BillingEmailMissing) and the price's currency (.UnsupportedCurrency), and registers the customer first.
type NewSubscriptionProviderKind = gen.NewSubscriptionProviderKind

// Notification is the Core API's Notification.
type Notification = gen.Notification

// Outcome is the Core API's Outcome.
type Outcome = gen.Outcome

// OutcomeProviderKind is the Core API's OutcomeProviderKind.
type OutcomeProviderKind = gen.OutcomeProviderKind

// OutcomeStatus is the Core API's OutcomeStatus.
type OutcomeStatus = gen.OutcomeStatus

// OverageLimit is the Core API's OverageLimit.
type OverageLimit = gen.OverageLimit

// PageAddon is the Core API's PageAddon.
type PageAddon = gen.PageAddon

// PageInvoiceSummary is the Core API's PageInvoiceSummary.
type PageInvoiceSummary = gen.PageInvoiceSummary

// PagePublishableKey is the Core API's PagePublishableKey.
type PagePublishableKey = gen.PagePublishableKey

// PageQueuedInvoice is the Core API's PageQueuedInvoice.
type PageQueuedInvoice = gen.PageQueuedInvoice

// PageRedemption is the Core API's PageRedemption.
type PageRedemption = gen.PageRedemption

// PageSessionInvoice is the Core API's PageSessionInvoice.
type PageSessionInvoice = gen.PageSessionInvoice

// PageVoucher is the Core API's PageVoucher.
type PageVoucher = gen.PageVoucher

// Payment is the Core API's Payment.
type Payment = gen.Payment

// PaymentMethodLabels is the Core API's PaymentMethodLabels.
type PaymentMethodLabels = gen.PaymentMethodLabels

// PaymentMethodLabelsStatus ACTIVE: charged automatically. EXPIRED, FAILED: a charge said it is no longer usable; save another one
type PaymentMethodLabelsStatus = gen.PaymentMethodLabelsStatus

// PaymentMethodSession is the Core API's PaymentMethodSession.
type PaymentMethodSession = gen.PaymentMethodSession

// PaymentStatus paid. requires_action: the customer must confirm the payment on hostedInvoiceUrl. failed: the charge was refused; save another payment method, or pay on hostedInvoiceUrl. processing: the outcome is not known yet; poll GET /public/session/invoices. not_required: nothing is due today.
type PaymentStatus = gen.PaymentStatus

// PlanChangeTarget is the Core API's PlanChangeTarget.
type PlanChangeTarget = gen.PlanChangeTarget

// PortalCapabilities is the Core API's PortalCapabilities.
type PortalCapabilities = gen.PortalCapabilities

// PortalCustomer is the Core API's PortalCustomer.
type PortalCustomer = gen.PortalCustomer

// PortalEntitlement is the Core API's PortalEntitlement.
type PortalEntitlement = gen.PortalEntitlement

// PortalEntitlementSource license: the licence grants it; addon: only add-ons do
type PortalEntitlementSource = gen.PortalEntitlementSource

// PortalEntitlementType is the Core API's PortalEntitlementType.
type PortalEntitlementType = gen.PortalEntitlementType

// PortalEntitlementLimitValue The effective value: a number's limit (-1 unlimited), a BOOLEAN's flag, a CONFIG's object
type PortalEntitlementLimitValue = gen.PortalEntitlement_LimitValue

// PortalInstance is the Core API's PortalInstance.
type PortalInstance = gen.PortalInstance

// PortalInstanceSummary is the Core API's PortalInstanceSummary.
type PortalInstanceSummary = gen.PortalInstanceSummary

// PortalInstanceSummaryBillingStatus Its subscription's status; null when it was never subscribed
type PortalInstanceSummaryBillingStatus = gen.PortalInstanceSummaryBillingStatus

// PortalPaymentMethod is the Core API's PortalPaymentMethod.
type PortalPaymentMethod = gen.PortalPaymentMethod

// PortalPaymentMethodStatus EXPIRED, FAILED: save another one
type PortalPaymentMethodStatus = gen.PortalPaymentMethodStatus

// PortalSession is the Core API's PortalSession.
type PortalSession = gen.PortalSession

// PortalUpcomingInvoice is the Core API's PortalUpcomingInvoice.
type PortalUpcomingInvoice = gen.PortalUpcomingInvoice

// PortalVoucher is the Core API's PortalVoucher.
type PortalVoucher = gen.PortalVoucher

// PortalVoucherStatus is the Core API's PortalVoucherStatus.
type PortalVoucherStatus = gen.PortalVoucherStatus

// PortalVoucherVoucherType is the Core API's PortalVoucherVoucherType.
type PortalVoucherVoucherType = gen.PortalVoucherVoucherType

// Price is the Core API's Price.
type Price = gen.Price

// PriceBillingModel FLAT_FEE: an amount per billing period. USAGE_BASED: an amount per sale unit of the metered entitlement's usage. OVERAGE: an amount per sale unit of the usage above the licence's limit.
type PriceBillingModel = gen.PriceBillingModel

// PriceBillingPeriod The period a FLAT_FEE price is charged for. Null on metered prices, which bill on the subscription's period.
type PriceBillingPeriod = gen.PriceBillingPeriod

// PriceBillingTiming ADVANCE bills the period that starts; ARREARS the period that ended. Metered prices are always ARREARS.
type PriceBillingTiming = gen.PriceBillingTiming

// PriceMeter is the Core API's PriceMeter.
type PriceMeter = gen.PriceMeter

// PriceStatus A DEPRECATED price keeps billing what is already pinned to it but is no longer offered.
type PriceStatus = gen.PriceStatus

// Provenance is the Core API's Provenance.
type Provenance = gen.Provenance

// ProvenanceAddon is the Core API's ProvenanceAddon.
type ProvenanceAddon = gen.ProvenanceAddon

// ProvenanceAddonOverrideBehavior is the Core API's ProvenanceAddonOverrideBehavior.
type ProvenanceAddonOverrideBehavior = gen.ProvenanceAddonOverrideBehavior

// ProvenanceAddonValue One unit's grant
type ProvenanceAddonValue = gen.ProvenanceAddon_Value

// ProvenanceBoost is the Core API's ProvenanceBoost.
type ProvenanceBoost = gen.ProvenanceBoost

// ProvenanceBoostModifierType is the Core API's ProvenanceBoostModifierType.
type ProvenanceBoostModifierType = gen.ProvenanceBoostModifierType

// ProvenanceLicense is the Core API's ProvenanceLicense.
type ProvenanceLicense = gen.ProvenanceLicense

// ProvenanceLicenseValue is the Core API's ProvenanceLicense_Value.
type ProvenanceLicenseValue = gen.ProvenanceLicense_Value

// ProvenanceNumber is the Core API's ProvenanceNumber.
type ProvenanceNumber = gen.ProvenanceNumber

// ProviderCapabilities is the Core API's ProviderCapabilities.
type ProviderCapabilities = gen.ProviderCapabilities

// ProviderRecord is the Core API's ProviderRecord.
type ProviderRecord = gen.ProviderRecord

// ProviderRecordReconciliationStatus is the Core API's ProviderRecordReconciliationStatus.
type ProviderRecordReconciliationStatus = gen.ProviderRecordReconciliationStatus

// ProviderRecordStatus The provider's status, mirrored verbatim
type ProviderRecordStatus = gen.ProviderRecordStatus

// ProviderSync is the Core API's ProviderSync.
type ProviderSync = gen.ProviderSync

// ProviderSyncLastSyncStatus is the Core API's ProviderSyncLastSyncStatus.
type ProviderSyncLastSyncStatus = gen.ProviderSyncLastSyncStatus

// ProviderSyncProviderKind is the Core API's ProviderSyncProviderKind.
type ProviderSyncProviderKind = gen.ProviderSyncProviderKind

// PublicAddOn is the Core API's PublicAddOn.
type PublicAddOn = gen.PublicAddOn

// PublicAddOnEntitlement is the Core API's PublicAddOnEntitlement.
type PublicAddOnEntitlement = gen.PublicAddOnEntitlement

// PublicAddOnEntitlementOverrideBehavior How the grant combines with the plan's: added to it, replacing it, or the larger of the two
type PublicAddOnEntitlementOverrideBehavior = gen.PublicAddOnEntitlementOverrideBehavior

// PublicAddOnEntitlementType is the Core API's PublicAddOnEntitlementType.
type PublicAddOnEntitlementType = gen.PublicAddOnEntitlementType

// PublicAddOnPricingType is the Core API's PublicAddOnPricingType.
type PublicAddOnPricingType = gen.PublicAddOnPricingType

// PublicCatalog is the Core API's PublicCatalog.
type PublicCatalog = gen.PublicCatalog

// PublicCatalogCapabilities is the Core API's PublicCatalogCapabilities.
type PublicCatalogCapabilities = gen.PublicCatalogCapabilities

// PublicMeter is the Core API's PublicMeter.
type PublicMeter = gen.PublicMeter

// PublicPlan is the Core API's PublicPlan.
type PublicPlan = gen.PublicPlan

// PublicPlanEntitlement is the Core API's PublicPlanEntitlement.
type PublicPlanEntitlement = gen.PublicPlanEntitlement

// PublicPlanEntitlementType is the Core API's PublicPlanEntitlementType.
type PublicPlanEntitlementType = gen.PublicPlanEntitlementType

// PublicPlanLifecycleState Always PUBLISHED: only the default version is listed, and a default version is published
type PublicPlanLifecycleState = gen.PublicPlanLifecycleState

// PublicPlanPricingType CUSTOM: priced on request; render selfServeCtaUrl
type PublicPlanPricingType = gen.PublicPlanPricingType

// PublicPrice is the Core API's PublicPrice.
type PublicPrice = gen.PublicPrice

// PublicPriceBillingModel FLAT_FEE: per billing period. USAGE_BASED: per sale unit of usage. OVERAGE: per sale unit of usage above the limit. Never sum metered prices into a headline.
type PublicPriceBillingModel = gen.PublicPriceBillingModel

// PublicPriceBillingPeriod The period a FLAT_FEE price is charged for; absent on metered prices
type PublicPriceBillingPeriod = gen.PublicPriceBillingPeriod

// PublicPriceBillingTiming is the Core API's PublicPriceBillingTiming.
type PublicPriceBillingTiming = gen.PublicPriceBillingTiming

// PublicSurface is the Core API's PublicSurface.
type PublicSurface = gen.PublicSurface

// PublishableKey is the Core API's PublishableKey.
type PublishableKey = gen.PublishableKey

// PublishableKeyCreated is the Core API's PublishableKeyCreated.
type PublishableKeyCreated = gen.PublishableKeyCreated

// PublishableKeyDraft is the Core API's PublishableKeyDraft.
type PublishableKeyDraft = gen.PublishableKeyDraft

// PublishableKeyPatch is the Core API's PublishableKeyPatch.
type PublishableKeyPatch = gen.PublishableKeyPatch

// PushFailures is the Core API's PushFailures.
type PushFailures = gen.PushFailures

// QueuedInvoice is the Core API's QueuedInvoice.
type QueuedInvoice = gen.QueuedInvoice

// QueuedInvoiceCollectionMethod is the Core API's QueuedInvoiceCollectionMethod.
type QueuedInvoiceCollectionMethod = gen.QueuedInvoiceCollectionMethod

// QueuedInvoiceHandoffStatus PENDING: waiting in the handoff queue for the organization's accounting system
type QueuedInvoiceHandoffStatus = gen.QueuedInvoiceHandoffStatus

// QueuedInvoiceHoldReason Why a DRAFT is held: the usage journal it was measured from failed a consistency check
type QueuedInvoiceHoldReason = gen.QueuedInvoiceHoldReason

// QueuedInvoiceKind ACTIVATION bills a subscription's first period, RENEWAL the period that ended and the one that starts, FINAL what is left when a subscription ends
type QueuedInvoiceKind = gen.QueuedInvoiceKind

// QueuedInvoiceProviderKind Who collects it. NOOP: the organization itself, through the handoff queue
type QueuedInvoiceProviderKind = gen.QueuedInvoiceProviderKind

// QueuedInvoiceStatus MANUAL: issued, for the organization to collect. PAID, UNCOLLECTIBLE and VOID are final. DRAFT: held, waiting for a release or a recompose
type QueuedInvoiceStatus = gen.QueuedInvoiceStatus

// Reconciliation is the Core API's Reconciliation.
type Reconciliation = gen.Reconciliation

// RedeemableCode is the Core API's RedeemableCode.
type RedeemableCode = gen.RedeemableCode

// Redemption is the Core API's Redemption.
type Redemption = gen.Redemption

// RedemptionRules is the Core API's RedemptionRules.
type RedemptionRules = gen.RedemptionRules

// RedemptionStatus is the Core API's RedemptionStatus.
type RedemptionStatus = gen.RedemptionStatus

// RedemptionVoucherType is the Core API's RedemptionVoucherType.
type RedemptionVoucherType = gen.RedemptionVoucherType

// RevocationReason is the Core API's RevocationReason.
type RevocationReason = gen.RevocationReason

// SampleUsageAmount is the Core API's SampleUsageAmount.
type SampleUsageAmount = gen.SampleUsageAmount

// ScheduledChange is the Core API's ScheduledChange.
type ScheduledChange = gen.ScheduledChange

// SessionAddon is the Core API's SessionAddon.
type SessionAddon = gen.SessionAddon

// SessionAddonQuantity is the Core API's SessionAddonQuantity.
type SessionAddonQuantity = gen.SessionAddonQuantity

// SessionBoost is the Core API's SessionBoost.
type SessionBoost = gen.SessionBoost

// SessionBoostModifierType is the Core API's SessionBoostModifierType.
type SessionBoostModifierType = gen.SessionBoostModifierType

// SessionCancellation is the Core API's SessionCancellation.
type SessionCancellation = gen.SessionCancellation

// SessionCheckout is the Core API's SessionCheckout.
type SessionCheckout = gen.SessionCheckout

// SessionCheckoutOrder is the Core API's SessionCheckoutOrder.
type SessionCheckoutOrder = gen.SessionCheckoutOrder

// SessionCheckoutStatus is the Core API's SessionCheckoutStatus.
type SessionCheckoutStatus = gen.SessionCheckoutStatus

// SessionDiscount is the Core API's SessionDiscount.
type SessionDiscount = gen.SessionDiscount

// SessionDiscountAppliesTo is the Core API's SessionDiscountAppliesTo.
type SessionDiscountAppliesTo = gen.SessionDiscountAppliesTo

// SessionDiscountType is the Core API's SessionDiscountType.
type SessionDiscountType = gen.SessionDiscountType

// SessionInvoice is the Core API's SessionInvoice.
type SessionInvoice = gen.SessionInvoice

// SessionInvoiceKind is the Core API's SessionInvoiceKind.
type SessionInvoiceKind = gen.SessionInvoiceKind

// SessionInvoiceStatus PUSHED: issued and awaiting payment. PAYMENT_FAILED: a charge failed; pay on hostedInvoiceUrl. MANUAL: issued, the vendor collects it.
type SessionInvoiceStatus = gen.SessionInvoiceStatus

// SessionLine is the Core API's SessionLine.
type SessionLine = gen.SessionLine

// SessionLineType is the Core API's SessionLineType.
type SessionLineType = gen.SessionLineType

// SessionPortal is the Core API's SessionPortal.
type SessionPortal = gen.SessionPortal

// SessionScheduledChange is the Core API's SessionScheduledChange.
type SessionScheduledChange = gen.SessionScheduledChange

// SessionSubscription is the Core API's SessionSubscription.
type SessionSubscription = gen.SessionSubscription

// SessionSubscriptionBillingPeriod is the Core API's SessionSubscriptionBillingPeriod.
type SessionSubscriptionBillingPeriod = gen.SessionSubscriptionBillingPeriod

// SessionSubscriptionCollectionMethod is the Core API's SessionSubscriptionCollectionMethod.
type SessionSubscriptionCollectionMethod = gen.SessionSubscriptionCollectionMethod

// SessionSubscriptionProviderKind is the Core API's SessionSubscriptionProviderKind.
type SessionSubscriptionProviderKind = gen.SessionSubscriptionProviderKind

// SessionSubscriptionStatus is the Core API's SessionSubscriptionStatus.
type SessionSubscriptionStatus = gen.SessionSubscriptionStatus

// SessionVoucher is the Core API's SessionVoucher.
type SessionVoucher = gen.SessionVoucher

// SessionVoucherCheck is the Core API's SessionVoucherCheck.
type SessionVoucherCheck = gen.SessionVoucherCheck

// SessionVoucherDuration PRICE: how many invoices it discounts (one, durationInPeriods, or every one). ENTITLEMENT_BOOST: how long it lasts, in billing periods.
type SessionVoucherDuration = gen.SessionVoucherDuration

// SessionVoucherValidity is the Core API's SessionVoucherValidity.
type SessionVoucherValidity = gen.SessionVoucherValidity

// SessionVoucherVoucherType is the Core API's SessionVoucherVoucherType.
type SessionVoucherVoucherType = gen.SessionVoucherVoucherType

// StartedSubscription is the Core API's StartedSubscription.
type StartedSubscription = gen.StartedSubscription

// StartedSubscriptionBillingPeriod is the Core API's StartedSubscriptionBillingPeriod.
type StartedSubscriptionBillingPeriod = gen.StartedSubscriptionBillingPeriod

// StartedSubscriptionCollectionMethod Effective: the subscription's own, else the organization's default
type StartedSubscriptionCollectionMethod = gen.StartedSubscriptionCollectionMethod

// StartedSubscriptionProviderKind Who collects its invoices. NOOP: the organization itself, through the handoff queue
type StartedSubscriptionProviderKind = gen.StartedSubscriptionProviderKind

// StartedSubscriptionStatus is the Core API's StartedSubscriptionStatus.
type StartedSubscriptionStatus = gen.StartedSubscriptionStatus

// SubscriptionAddon is the Core API's SubscriptionAddon.
type SubscriptionAddon = gen.SubscriptionAddon

// SubscriptionCancellation is the Core API's SubscriptionCancellation.
type SubscriptionCancellation = gen.SubscriptionCancellation

// SubscriptionCancellationMode AT_PERIOD_END (the default): the period paid for runs out, then the FINAL invoice is issued. IMMEDIATE: the FINAL invoice is issued now, flat fees billed in full and nothing refunded.
type SubscriptionCancellationMode = gen.SubscriptionCancellationMode

// SubscriptionTerms is the Core API's SubscriptionTerms.
type SubscriptionTerms = gen.SubscriptionTerms

// SubscriptionTermsCollectionMethod null: the organization's default. CHARGE_AUTOMATICALLY needs a provider that charges automatically (UpdateInstanceBilling.CollectionMethodUnsupported).
type SubscriptionTermsCollectionMethod = gen.SubscriptionTermsCollectionMethod

// SubscriptionTermsProviderKind Who issues and collects the invoices, from the next one on. NOOP: the organization itself, through the handoff queue. A provider must be connected (UpdateInstanceBilling.ProviderNotConnected); one that pushes invoices needs the customer's billing e-mail for SEND_INVOICE (.BillingEmailMissing) and the subscription's currency (.UnsupportedCurrency).
type SubscriptionTermsProviderKind = gen.SubscriptionTermsProviderKind

// SyncReport is the Core API's SyncReport.
type SyncReport = gen.SyncReport

// TotalsDifference is the Core API's TotalsDifference.
type TotalsDifference = gen.TotalsDifference

// Validity is the Core API's Validity.
type Validity = gen.Validity

// ValidityReason The first failing check
type ValidityReason = gen.ValidityReason

// ValidityRule With NOT_ELIGIBLE: the eligibility rule
type ValidityRule = gen.ValidityRule

// Voucher is the Core API's Voucher.
type Voucher = gen.Voucher

// VoucherCheck is the Core API's VoucherCheck.
type VoucherCheck = gen.VoucherCheck

// VoucherCode is the Core API's VoucherCode.
type VoucherCode = gen.VoucherCode

// VoucherDraft is the Core API's VoucherDraft.
type VoucherDraft = gen.VoucherDraft

// VoucherDraftDuration is the Core API's VoucherDraftDuration.
type VoucherDraftDuration = gen.VoucherDraftDuration

// VoucherDraftPriceAppliesTo is the Core API's VoucherDraftPriceAppliesTo.
type VoucherDraftPriceAppliesTo = gen.VoucherDraftPriceAppliesTo

// VoucherDraftPriceDiscountType is the Core API's VoucherDraftPriceDiscountType.
type VoucherDraftPriceDiscountType = gen.VoucherDraftPriceDiscountType

// VoucherDraftVoucherType is the Core API's VoucherDraftVoucherType.
type VoucherDraftVoucherType = gen.VoucherDraftVoucherType

// VoucherDuration PRICE: how many invoices it discounts (one, durationInPeriods, or every one). BOOST: how long a redemption lasts, in the subscription's billing periods.
type VoucherDuration = gen.VoucherDuration

// VoucherPriceAppliesTo is the Core API's VoucherPriceAppliesTo.
type VoucherPriceAppliesTo = gen.VoucherPriceAppliesTo

// VoucherPriceDiscountType is the Core API's VoucherPriceDiscountType.
type VoucherPriceDiscountType = gen.VoucherPriceDiscountType

// VoucherStatus is the Core API's VoucherStatus.
type VoucherStatus = gen.VoucherStatus

// VoucherVoucherType is the Core API's VoucherVoucherType.
type VoucherVoucherType = gen.VoucherVoucherType

// Values of the billing enums.
const (
	AddonEntitlementEntitlementTypeBOOLEAN                    = gen.AddonEntitlementEntitlementTypeBOOLEAN
	AddonEntitlementEntitlementTypeCONFIG                     = gen.AddonEntitlementEntitlementTypeCONFIG
	AddonEntitlementEntitlementTypeNUMBER                     = gen.AddonEntitlementEntitlementTypeNUMBER
	AddonEntitlementEntitlementTypeNUMBERAICREDIT             = gen.AddonEntitlementEntitlementTypeNUMBERAICREDIT
	AddonEntitlementOverrideBehaviorADD                       = gen.AddonEntitlementOverrideBehaviorADD
	AddonEntitlementOverrideBehaviorMAX                       = gen.AddonEntitlementOverrideBehaviorMAX
	AddonEntitlementOverrideBehaviorOVERRIDE                  = gen.AddonEntitlementOverrideBehaviorOVERRIDE
	AddonGrantOverrideBehaviorADD                             = gen.AddonGrantOverrideBehaviorADD
	AddonGrantOverrideBehaviorMAX                             = gen.AddonGrantOverrideBehaviorMAX
	AddonGrantOverrideBehaviorOVERRIDE                        = gen.AddonGrantOverrideBehaviorOVERRIDE
	AddonLifecycleStateARCHIVED                               = gen.AddonLifecycleStateARCHIVED
	AddonLifecycleStateDRAFT                                  = gen.AddonLifecycleStateDRAFT
	AddonLifecycleStatePUBLISHED                              = gen.AddonLifecycleStatePUBLISHED
	AddonPricingTypeCUSTOM                                    = gen.AddonPricingTypeCUSTOM
	AddonPricingTypeFREE                                      = gen.AddonPricingTypeFREE
	AddonPricingTypePAID                                      = gen.AddonPricingTypePAID
	BillingCapabilitiesDisabledReasonDEPLOYMENTDISABLED       = gen.BillingCapabilitiesDisabledReasonDEPLOYMENTDISABLED
	BillingCapabilitiesDisabledReasonLessThannil              = gen.BillingCapabilitiesDisabledReasonLessThannil
	BillingCapabilitiesDisabledReasonNOTENTITLED              = gen.BillingCapabilitiesDisabledReasonNOTENTITLED
	BillingProviderKindNOOP                                   = gen.BillingProviderKindNOOP
	BillingProviderKindSTRIPE                                 = gen.BillingProviderKindSTRIPE
	BillingProviderUnavailableReasonNOTENTITLED               = gen.BillingProviderUnavailableReasonNOTENTITLED
	BillingProviderUnavailableReasonVAULTNOTCONFIGURED        = gen.BillingProviderUnavailableReasonVAULTNOTCONFIGURED
	BillingSettingsDefaultCollectionMethodCHARGEAUTOMATICALLY = gen.BillingSettingsDefaultCollectionMethodCHARGEAUTOMATICALLY
	BillingSettingsDefaultCollectionMethodSENDINVOICE         = gen.BillingSettingsDefaultCollectionMethodSENDINVOICE
	CanceledSubscriptionBillingPeriodANNUAL                   = gen.CanceledSubscriptionBillingPeriodANNUAL
	CanceledSubscriptionBillingPeriodMONTHLY                  = gen.CanceledSubscriptionBillingPeriodMONTHLY
	CanceledSubscriptionBillingPeriodQUARTERLY                = gen.CanceledSubscriptionBillingPeriodQUARTERLY
	CanceledSubscriptionBillingPeriodSEMIANNUAL               = gen.CanceledSubscriptionBillingPeriodSEMIANNUAL
	CanceledSubscriptionCollectionMethodCHARGEAUTOMATICALLY   = gen.CanceledSubscriptionCollectionMethodCHARGEAUTOMATICALLY
	CanceledSubscriptionCollectionMethodSENDINVOICE           = gen.CanceledSubscriptionCollectionMethodSENDINVOICE
	CanceledSubscriptionProviderKindNOOP                      = gen.CanceledSubscriptionProviderKindNOOP
	CanceledSubscriptionProviderKindSTRIPE                    = gen.CanceledSubscriptionProviderKindSTRIPE
	CanceledSubscriptionStatusACTIVE                          = gen.CanceledSubscriptionStatusACTIVE
	CanceledSubscriptionStatusCANCELED                        = gen.CanceledSubscriptionStatusCANCELED
	CanceledSubscriptionStatusPASTDUE                         = gen.CanceledSubscriptionStatusPASTDUE
	CanceledSubscriptionStatusTRIAL                           = gen.CanceledSubscriptionStatusTRIAL
	CancellationBillingPeriodANNUAL                           = gen.CancellationBillingPeriodANNUAL
	CancellationBillingPeriodMONTHLY                          = gen.CancellationBillingPeriodMONTHLY
	CancellationBillingPeriodQUARTERLY                        = gen.CancellationBillingPeriodQUARTERLY
	CancellationBillingPeriodSEMIANNUAL                       = gen.CancellationBillingPeriodSEMIANNUAL
	CancellationCollectionMethodCHARGEAUTOMATICALLY           = gen.CancellationCollectionMethodCHARGEAUTOMATICALLY
	CancellationCollectionMethodSENDINVOICE                   = gen.CancellationCollectionMethodSENDINVOICE
	CancellationModeATPERIODEND                               = gen.CancellationModeATPERIODEND
	CancellationModeIMMEDIATE                                 = gen.CancellationModeIMMEDIATE
	CancellationProviderKindNOOP                              = gen.CancellationProviderKindNOOP
	CancellationProviderKindSTRIPE                            = gen.CancellationProviderKindSTRIPE
	CancellationStatusACTIVE                                  = gen.CancellationStatusACTIVE
	CancellationStatusCANCELED                                = gen.CancellationStatusCANCELED
	CancellationStatusPASTDUE                                 = gen.CancellationStatusPASTDUE
	CancellationStatusTRIAL                                   = gen.CancellationStatusTRIAL
	ClosedInvoiceKindACTIVATION                               = gen.ClosedInvoiceKindACTIVATION
	ClosedInvoiceKindFINAL                                    = gen.ClosedInvoiceKindFINAL
	ClosedInvoiceKindRENEWAL                                  = gen.ClosedInvoiceKindRENEWAL
	ClosedInvoiceStatusDRAFT                                  = gen.ClosedInvoiceStatusDRAFT
	ClosedInvoiceStatusMANUAL                                 = gen.ClosedInvoiceStatusMANUAL
	ClosedInvoiceStatusPAID                                   = gen.ClosedInvoiceStatusPAID
	ClosedInvoiceStatusPAYMENTFAILED                          = gen.ClosedInvoiceStatusPAYMENTFAILED
	ClosedInvoiceStatusPUSHED                                 = gen.ClosedInvoiceStatusPUSHED
	ClosedInvoiceStatusPUSHFAILED                             = gen.ClosedInvoiceStatusPUSHFAILED
	ClosedInvoiceStatusUNCOLLECTIBLE                          = gen.ClosedInvoiceStatusUNCOLLECTIBLE
	ClosedInvoiceStatusVOID                                   = gen.ClosedInvoiceStatusVOID
	CustomerInProviderProviderKindSTRIPE                      = gen.CustomerInProviderProviderKindSTRIPE
	EntitlementTypeBOOLEAN                                    = gen.EntitlementTypeBOOLEAN
	EntitlementTypeCONFIG                                     = gen.EntitlementTypeCONFIG
	EntitlementTypeNUMBER                                     = gen.EntitlementTypeNUMBER
	EntitlementTypeNUMBERAICREDIT                             = gen.EntitlementTypeNUMBERAICREDIT
	ExportInvoicesParamsHandoffStatusACKNOWLEDGED             = gen.ExportInvoicesParamsHandoffStatusACKNOWLEDGED
	ExportInvoicesParamsHandoffStatusNOTREQUIRED              = gen.ExportInvoicesParamsHandoffStatusNOTREQUIRED
	ExportInvoicesParamsHandoffStatusPENDING                  = gen.ExportInvoicesParamsHandoffStatusPENDING
	ExportInvoicesParamsKindACTIVATION                        = gen.ExportInvoicesParamsKindACTIVATION
	ExportInvoicesParamsKindFINAL                             = gen.ExportInvoicesParamsKindFINAL
	ExportInvoicesParamsKindRENEWAL                           = gen.ExportInvoicesParamsKindRENEWAL
	ExportInvoicesParamsProviderKindNOOP                      = gen.ExportInvoicesParamsProviderKindNOOP
	ExportInvoicesParamsProviderKindSTRIPE                    = gen.ExportInvoicesParamsProviderKindSTRIPE
	ExportInvoicesParamsStatusDRAFT                           = gen.ExportInvoicesParamsStatusDRAFT
	ExportInvoicesParamsStatusMANUAL                          = gen.ExportInvoicesParamsStatusMANUAL
	ExportInvoicesParamsStatusPAID                            = gen.ExportInvoicesParamsStatusPAID
	ExportInvoicesParamsStatusPAYMENTFAILED                   = gen.ExportInvoicesParamsStatusPAYMENTFAILED
	ExportInvoicesParamsStatusPUSHED                          = gen.ExportInvoicesParamsStatusPUSHED
	ExportInvoicesParamsStatusPUSHFAILED                      = gen.ExportInvoicesParamsStatusPUSHFAILED
	ExportInvoicesParamsStatusUNCOLLECTIBLE                   = gen.ExportInvoicesParamsStatusUNCOLLECTIBLE
	ExportInvoicesParamsStatusVOID                            = gen.ExportInvoicesParamsStatusVOID
	Addons                                                    = gen.Addons
	None                                                      = gen.None
	GrantModifierTypeADD                                      = gen.GrantModifierTypeADD
	GrantModifierTypeMULTIPLY                                 = gen.GrantModifierTypeMULTIPLY
	GrantModifierTypeSET                                      = gen.GrantModifierTypeSET
	GrantModifierTypeUNLIMITED                                = gen.GrantModifierTypeUNLIMITED
	InstanceBillingBillingPeriodANNUAL                        = gen.InstanceBillingBillingPeriodANNUAL
	InstanceBillingBillingPeriodMONTHLY                       = gen.InstanceBillingBillingPeriodMONTHLY
	InstanceBillingBillingPeriodQUARTERLY                     = gen.InstanceBillingBillingPeriodQUARTERLY
	InstanceBillingBillingPeriodSEMIANNUAL                    = gen.InstanceBillingBillingPeriodSEMIANNUAL
	InstanceBillingCollectionMethodCHARGEAUTOMATICALLY        = gen.InstanceBillingCollectionMethodCHARGEAUTOMATICALLY
	InstanceBillingCollectionMethodSENDINVOICE                = gen.InstanceBillingCollectionMethodSENDINVOICE
	InstanceBillingProviderKindNOOP                           = gen.InstanceBillingProviderKindNOOP
	InstanceBillingProviderKindSTRIPE                         = gen.InstanceBillingProviderKindSTRIPE
	InstanceBillingStatusACTIVE                               = gen.InstanceBillingStatusACTIVE
	InstanceBillingStatusCANCELED                             = gen.InstanceBillingStatusCANCELED
	InstanceBillingStatusPASTDUE                              = gen.InstanceBillingStatusPASTDUE
	InstanceBillingStatusTRIAL                                = gen.InstanceBillingStatusTRIAL
	InstanceStatusDEGRADED                                    = gen.InstanceStatusDEGRADED
	InstanceStatusHEALTHY                                     = gen.InstanceStatusHEALTHY
	InstanceStatusINCIDENT                                    = gen.InstanceStatusINCIDENT
	InstanceStatusMAINTENANCE                                 = gen.InstanceStatusMAINTENANCE
	InvoiceCollectionMethodCHARGEAUTOMATICALLY                = gen.InvoiceCollectionMethodCHARGEAUTOMATICALLY
	InvoiceCollectionMethodSENDINVOICE                        = gen.InvoiceCollectionMethodSENDINVOICE
	InvoiceHandoffStateACKNOWLEDGED                           = gen.InvoiceHandoffStateACKNOWLEDGED
	InvoiceHandoffStateNOTREQUIRED                            = gen.InvoiceHandoffStateNOTREQUIRED
	InvoiceHandoffStatePENDING                                = gen.InvoiceHandoffStatePENDING
	InvoiceHandoffStatusACKNOWLEDGED                          = gen.InvoiceHandoffStatusACKNOWLEDGED
	InvoiceHandoffStatusNOTREQUIRED                           = gen.InvoiceHandoffStatusNOTREQUIRED
	InvoiceHandoffStatusPENDING                               = gen.InvoiceHandoffStatusPENDING
	InvoiceHoldReasonLEDGERCHAINBREAK                         = gen.InvoiceHoldReasonLEDGERCHAINBREAK
	InvoiceHoldReasonLEDGERCOUNTERMISMATCH                    = gen.InvoiceHoldReasonLEDGERCOUNTERMISMATCH
	InvoiceHoldReasonLEDGERSEQUENCEGAP                        = gen.InvoiceHoldReasonLEDGERSEQUENCEGAP
	InvoiceHoldReasonLessThannil                              = gen.InvoiceHoldReasonLessThannil
	InvoiceKindACTIVATION                                     = gen.InvoiceKindACTIVATION
	InvoiceKindFINAL                                          = gen.InvoiceKindFINAL
	InvoiceKindRENEWAL                                        = gen.InvoiceKindRENEWAL
	InvoiceLineBillingModelFLATFEE                            = gen.InvoiceLineBillingModelFLATFEE
	InvoiceLineBillingModelLessThannil                        = gen.InvoiceLineBillingModelLessThannil
	InvoiceLineBillingModelOVERAGE                            = gen.InvoiceLineBillingModelOVERAGE
	InvoiceLineBillingModelUSAGEBASED                         = gen.InvoiceLineBillingModelUSAGEBASED
	InvoiceLineBillingTimingADVANCE                           = gen.InvoiceLineBillingTimingADVANCE
	InvoiceLineBillingTimingARREARS                           = gen.InvoiceLineBillingTimingARREARS
	InvoiceLineBillingTimingLessThannil                       = gen.InvoiceLineBillingTimingLessThannil
	InvoiceLineDiscountAppliesToADDONS                        = gen.InvoiceLineDiscountAppliesToADDONS
	InvoiceLineDiscountAppliesToBOTH                          = gen.InvoiceLineDiscountAppliesToBOTH
	InvoiceLineDiscountAppliesToLICENSEBASE                   = gen.InvoiceLineDiscountAppliesToLICENSEBASE
	InvoiceLineDiscountAppliesToSELECTEDPRICES                = gen.InvoiceLineDiscountAppliesToSELECTEDPRICES
	InvoiceLineDiscountDiscountTypeFIXEDAMOUNT                = gen.InvoiceLineDiscountDiscountTypeFIXEDAMOUNT
	InvoiceLineDiscountDiscountTypePERCENTAGE                 = gen.InvoiceLineDiscountDiscountTypePERCENTAGE
	InvoiceLineTypeADDON                                      = gen.InvoiceLineTypeADDON
	InvoiceLineTypeBASE                                       = gen.InvoiceLineTypeBASE
	InvoiceLineTypeDISCOUNT                                   = gen.InvoiceLineTypeDISCOUNT
	InvoiceLineTypeOVERAGE                                    = gen.InvoiceLineTypeOVERAGE
	InvoiceLineTypeUSAGE                                      = gen.InvoiceLineTypeUSAGE
	InvoicePreviewKindACTIVATION                              = gen.InvoicePreviewKindACTIVATION
	InvoicePreviewKindFINAL                                   = gen.InvoicePreviewKindFINAL
	InvoicePreviewKindRENEWAL                                 = gen.InvoicePreviewKindRENEWAL
	PREVIEW                                                   = gen.PREVIEW
	InvoiceProviderKindNOOP                                   = gen.InvoiceProviderKindNOOP
	InvoiceProviderKindSTRIPE                                 = gen.InvoiceProviderKindSTRIPE
	InvoiceStatusDRAFT                                        = gen.InvoiceStatusDRAFT
	InvoiceStatusMANUAL                                       = gen.InvoiceStatusMANUAL
	InvoiceStatusPAID                                         = gen.InvoiceStatusPAID
	InvoiceStatusPAYMENTFAILED                                = gen.InvoiceStatusPAYMENTFAILED
	InvoiceStatusPUSHED                                       = gen.InvoiceStatusPUSHED
	InvoiceStatusPUSHFAILED                                   = gen.InvoiceStatusPUSHFAILED
	InvoiceStatusUNCOLLECTIBLE                                = gen.InvoiceStatusUNCOLLECTIBLE
	InvoiceStatusVOID                                         = gen.InvoiceStatusVOID
	InvoiceSummaryCollectionMethodCHARGEAUTOMATICALLY         = gen.InvoiceSummaryCollectionMethodCHARGEAUTOMATICALLY
	InvoiceSummaryCollectionMethodSENDINVOICE                 = gen.InvoiceSummaryCollectionMethodSENDINVOICE
	InvoiceSummaryHandoffStatusACKNOWLEDGED                   = gen.InvoiceSummaryHandoffStatusACKNOWLEDGED
	InvoiceSummaryHandoffStatusNOTREQUIRED                    = gen.InvoiceSummaryHandoffStatusNOTREQUIRED
	InvoiceSummaryHandoffStatusPENDING                        = gen.InvoiceSummaryHandoffStatusPENDING
	InvoiceSummaryHoldReasonLEDGERCHAINBREAK                  = gen.InvoiceSummaryHoldReasonLEDGERCHAINBREAK
	InvoiceSummaryHoldReasonLEDGERCOUNTERMISMATCH             = gen.InvoiceSummaryHoldReasonLEDGERCOUNTERMISMATCH
	InvoiceSummaryHoldReasonLEDGERSEQUENCEGAP                 = gen.InvoiceSummaryHoldReasonLEDGERSEQUENCEGAP
	InvoiceSummaryHoldReasonLessThannil                       = gen.InvoiceSummaryHoldReasonLessThannil
	InvoiceSummaryKindACTIVATION                              = gen.InvoiceSummaryKindACTIVATION
	InvoiceSummaryKindFINAL                                   = gen.InvoiceSummaryKindFINAL
	InvoiceSummaryKindRENEWAL                                 = gen.InvoiceSummaryKindRENEWAL
	InvoiceSummaryProviderKindNOOP                            = gen.InvoiceSummaryProviderKindNOOP
	InvoiceSummaryProviderKindSTRIPE                          = gen.InvoiceSummaryProviderKindSTRIPE
	InvoiceSummaryStatusDRAFT                                 = gen.InvoiceSummaryStatusDRAFT
	InvoiceSummaryStatusMANUAL                                = gen.InvoiceSummaryStatusMANUAL
	InvoiceSummaryStatusPAID                                  = gen.InvoiceSummaryStatusPAID
	InvoiceSummaryStatusPAYMENTFAILED                         = gen.InvoiceSummaryStatusPAYMENTFAILED
	InvoiceSummaryStatusPUSHED                                = gen.InvoiceSummaryStatusPUSHED
	InvoiceSummaryStatusPUSHFAILED                            = gen.InvoiceSummaryStatusPUSHFAILED
	InvoiceSummaryStatusUNCOLLECTIBLE                         = gen.InvoiceSummaryStatusUNCOLLECTIBLE
	InvoiceSummaryStatusVOID                                  = gen.InvoiceSummaryStatusVOID
	LicenseLifecycleStateARCHIVED                             = gen.LicenseLifecycleStateARCHIVED
	LicenseLifecycleStateDRAFT                                = gen.LicenseLifecycleStateDRAFT
	LicenseLifecycleStatePUBLISHED                            = gen.LicenseLifecycleStatePUBLISHED
	LicensePriceChangesBillingPeriodANNUAL                    = gen.LicensePriceChangesBillingPeriodANNUAL
	LicensePriceChangesBillingPeriodMONTHLY                   = gen.LicensePriceChangesBillingPeriodMONTHLY
	LicensePriceChangesBillingPeriodQUARTERLY                 = gen.LicensePriceChangesBillingPeriodQUARTERLY
	LicensePriceChangesBillingPeriodSEMIANNUAL                = gen.LicensePriceChangesBillingPeriodSEMIANNUAL
	LicensePriceChangesBillingTimingADVANCE                   = gen.LicensePriceChangesBillingTimingADVANCE
	LicensePriceChangesBillingTimingARREARS                   = gen.LicensePriceChangesBillingTimingARREARS
	LicensePricingTypeCUSTOM                                  = gen.LicensePricingTypeCUSTOM
	LicensePricingTypeFREE                                    = gen.LicensePricingTypeFREE
	LicensePricingTypePAID                                    = gen.LicensePricingTypePAID
	LicenseTypeCOMMUNITY                                      = gen.LicenseTypeCOMMUNITY
	LicenseTypeDEVELOPMENT                                    = gen.LicenseTypeDEVELOPMENT
	LicenseTypePAID                                           = gen.LicenseTypePAID
	LicenseTypeTRIAL                                          = gen.LicenseTypeTRIAL
	ListAddonPricesParamsStatusACTIVE                         = gen.ListAddonPricesParamsStatusACTIVE
	ListAddonPricesParamsStatusDEPRECATED                     = gen.ListAddonPricesParamsStatusDEPRECATED
	ListAddonsParamsLifecycleStateARCHIVED                    = gen.ListAddonsParamsLifecycleStateARCHIVED
	ListAddonsParamsLifecycleStateDRAFT                       = gen.ListAddonsParamsLifecycleStateDRAFT
	ListAddonsParamsLifecycleStatePUBLISHED                   = gen.ListAddonsParamsLifecycleStatePUBLISHED
	ListHandoffParamsStatusACKNOWLEDGED                       = gen.ListHandoffParamsStatusACKNOWLEDGED
	ListHandoffParamsStatusPENDING                            = gen.ListHandoffParamsStatusPENDING
	ListInstanceInvoicesParamsHandoffStatusACKNOWLEDGED       = gen.ListInstanceInvoicesParamsHandoffStatusACKNOWLEDGED
	ListInstanceInvoicesParamsHandoffStatusNOTREQUIRED        = gen.ListInstanceInvoicesParamsHandoffStatusNOTREQUIRED
	ListInstanceInvoicesParamsHandoffStatusPENDING            = gen.ListInstanceInvoicesParamsHandoffStatusPENDING
	ListInstanceInvoicesParamsKindACTIVATION                  = gen.ListInstanceInvoicesParamsKindACTIVATION
	ListInstanceInvoicesParamsKindFINAL                       = gen.ListInstanceInvoicesParamsKindFINAL
	ListInstanceInvoicesParamsKindRENEWAL                     = gen.ListInstanceInvoicesParamsKindRENEWAL
	ListInstanceInvoicesParamsProviderKindNOOP                = gen.ListInstanceInvoicesParamsProviderKindNOOP
	ListInstanceInvoicesParamsProviderKindSTRIPE              = gen.ListInstanceInvoicesParamsProviderKindSTRIPE
	ListInstanceInvoicesParamsStatusDRAFT                     = gen.ListInstanceInvoicesParamsStatusDRAFT
	ListInstanceInvoicesParamsStatusMANUAL                    = gen.ListInstanceInvoicesParamsStatusMANUAL
	ListInstanceInvoicesParamsStatusPAID                      = gen.ListInstanceInvoicesParamsStatusPAID
	ListInstanceInvoicesParamsStatusPAYMENTFAILED             = gen.ListInstanceInvoicesParamsStatusPAYMENTFAILED
	ListInstanceInvoicesParamsStatusPUSHED                    = gen.ListInstanceInvoicesParamsStatusPUSHED
	ListInstanceInvoicesParamsStatusPUSHFAILED                = gen.ListInstanceInvoicesParamsStatusPUSHFAILED
	ListInstanceInvoicesParamsStatusUNCOLLECTIBLE             = gen.ListInstanceInvoicesParamsStatusUNCOLLECTIBLE
	ListInstanceInvoicesParamsStatusVOID                      = gen.ListInstanceInvoicesParamsStatusVOID
	ListInstanceVouchersParamsStatusACTIVE                    = gen.ListInstanceVouchersParamsStatusACTIVE
	ListInstanceVouchersParamsStatusEXPIRED                   = gen.ListInstanceVouchersParamsStatusEXPIRED
	ListInstanceVouchersParamsStatusREVOKED                   = gen.ListInstanceVouchersParamsStatusREVOKED
	ListInvoicesParamsHandoffStatusACKNOWLEDGED               = gen.ListInvoicesParamsHandoffStatusACKNOWLEDGED
	ListInvoicesParamsHandoffStatusNOTREQUIRED                = gen.ListInvoicesParamsHandoffStatusNOTREQUIRED
	ListInvoicesParamsHandoffStatusPENDING                    = gen.ListInvoicesParamsHandoffStatusPENDING
	ListInvoicesParamsKindACTIVATION                          = gen.ListInvoicesParamsKindACTIVATION
	ListInvoicesParamsKindFINAL                               = gen.ListInvoicesParamsKindFINAL
	ListInvoicesParamsKindRENEWAL                             = gen.ListInvoicesParamsKindRENEWAL
	ListInvoicesParamsProviderKindNOOP                        = gen.ListInvoicesParamsProviderKindNOOP
	ListInvoicesParamsProviderKindSTRIPE                      = gen.ListInvoicesParamsProviderKindSTRIPE
	ListInvoicesParamsStatusDRAFT                             = gen.ListInvoicesParamsStatusDRAFT
	ListInvoicesParamsStatusMANUAL                            = gen.ListInvoicesParamsStatusMANUAL
	ListInvoicesParamsStatusPAID                              = gen.ListInvoicesParamsStatusPAID
	ListInvoicesParamsStatusPAYMENTFAILED                     = gen.ListInvoicesParamsStatusPAYMENTFAILED
	ListInvoicesParamsStatusPUSHED                            = gen.ListInvoicesParamsStatusPUSHED
	ListInvoicesParamsStatusPUSHFAILED                        = gen.ListInvoicesParamsStatusPUSHFAILED
	ListInvoicesParamsStatusUNCOLLECTIBLE                     = gen.ListInvoicesParamsStatusUNCOLLECTIBLE
	ListInvoicesParamsStatusVOID                              = gen.ListInvoicesParamsStatusVOID
	ListLicensePricesParamsBillingModelFLATFEE                = gen.ListLicensePricesParamsBillingModelFLATFEE
	ListLicensePricesParamsBillingModelOVERAGE                = gen.ListLicensePricesParamsBillingModelOVERAGE
	ListLicensePricesParamsBillingModelUSAGEBASED             = gen.ListLicensePricesParamsBillingModelUSAGEBASED
	ListLicensePricesParamsStatusACTIVE                       = gen.ListLicensePricesParamsStatusACTIVE
	ListLicensePricesParamsStatusDEPRECATED                   = gen.ListLicensePricesParamsStatusDEPRECATED
	ListVouchersParamsStatusACTIVE                            = gen.ListVouchersParamsStatusACTIVE
	ListVouchersParamsStatusARCHIVED                          = gen.ListVouchersParamsStatusARCHIVED
	ListVouchersParamsStatusDRAFT                             = gen.ListVouchersParamsStatusDRAFT
	ListVouchersParamsStatusEXHAUSTED                         = gen.ListVouchersParamsStatusEXHAUSTED
	ListVouchersParamsStatusEXPIRED                           = gen.ListVouchersParamsStatusEXPIRED
	ListVouchersParamsVoucherTypeENTITLEMENTBOOST             = gen.ListVouchersParamsVoucherTypeENTITLEMENTBOOST
	ListVouchersParamsVoucherTypePRICE                        = gen.ListVouchersParamsVoucherTypePRICE
	NewAddonGrantOverrideBehaviorADD                          = gen.NewAddonGrantOverrideBehaviorADD
	NewAddonGrantOverrideBehaviorMAX                          = gen.NewAddonGrantOverrideBehaviorMAX
	NewAddonGrantOverrideBehaviorOVERRIDE                     = gen.NewAddonGrantOverrideBehaviorOVERRIDE
	NewAddonLifecycleStateDRAFT                               = gen.NewAddonLifecycleStateDRAFT
	NewAddonLifecycleStatePUBLISHED                           = gen.NewAddonLifecycleStatePUBLISHED
	NewAddonPriceBillingModelFLATFEE                          = gen.NewAddonPriceBillingModelFLATFEE
	NewAddonPriceBillingModelOVERAGE                          = gen.NewAddonPriceBillingModelOVERAGE
	NewAddonPriceBillingModelUSAGEBASED                       = gen.NewAddonPriceBillingModelUSAGEBASED
	NewAddonPriceBillingPeriodANNUAL                          = gen.NewAddonPriceBillingPeriodANNUAL
	NewAddonPriceBillingPeriodMONTHLY                         = gen.NewAddonPriceBillingPeriodMONTHLY
	NewAddonPriceBillingPeriodQUARTERLY                       = gen.NewAddonPriceBillingPeriodQUARTERLY
	NewAddonPriceBillingPeriodSEMIANNUAL                      = gen.NewAddonPriceBillingPeriodSEMIANNUAL
	NewAddonPriceBillingTimingADVANCE                         = gen.NewAddonPriceBillingTimingADVANCE
	NewAddonPriceBillingTimingARREARS                         = gen.NewAddonPriceBillingTimingARREARS
	NewAddonPricingTypeCUSTOM                                 = gen.NewAddonPricingTypeCUSTOM
	NewAddonPricingTypeFREE                                   = gen.NewAddonPricingTypeFREE
	NewAddonPricingTypePAID                                   = gen.NewAddonPricingTypePAID
	NewLicensePriceBillingModelFLATFEE                        = gen.NewLicensePriceBillingModelFLATFEE
	NewLicensePriceBillingModelOVERAGE                        = gen.NewLicensePriceBillingModelOVERAGE
	NewLicensePriceBillingModelUSAGEBASED                     = gen.NewLicensePriceBillingModelUSAGEBASED
	NewLicensePriceBillingPeriodANNUAL                        = gen.NewLicensePriceBillingPeriodANNUAL
	NewLicensePriceBillingPeriodMONTHLY                       = gen.NewLicensePriceBillingPeriodMONTHLY
	NewLicensePriceBillingPeriodQUARTERLY                     = gen.NewLicensePriceBillingPeriodQUARTERLY
	NewLicensePriceBillingPeriodSEMIANNUAL                    = gen.NewLicensePriceBillingPeriodSEMIANNUAL
	NewLicensePriceBillingTimingADVANCE                       = gen.NewLicensePriceBillingTimingADVANCE
	NewLicensePriceBillingTimingARREARS                       = gen.NewLicensePriceBillingTimingARREARS
	NewSubscriptionCollectionMethodCHARGEAUTOMATICALLY        = gen.NewSubscriptionCollectionMethodCHARGEAUTOMATICALLY
	NewSubscriptionCollectionMethodSENDINVOICE                = gen.NewSubscriptionCollectionMethodSENDINVOICE
	NewSubscriptionProviderKindNOOP                           = gen.NewSubscriptionProviderKindNOOP
	NewSubscriptionProviderKindSTRIPE                         = gen.NewSubscriptionProviderKindSTRIPE
	OutcomeProviderKindNOOP                                   = gen.OutcomeProviderKindNOOP
	OutcomeProviderKindSTRIPE                                 = gen.OutcomeProviderKindSTRIPE
	OutcomeStatusFAILED                                       = gen.OutcomeStatusFAILED
	OutcomeStatusPARTIAL                                      = gen.OutcomeStatusPARTIAL
	OutcomeStatusSUCCESS                                      = gen.OutcomeStatusSUCCESS
	PaymentMethodLabelsStatusACTIVE                           = gen.PaymentMethodLabelsStatusACTIVE
	PaymentMethodLabelsStatusEXPIRED                          = gen.PaymentMethodLabelsStatusEXPIRED
	PaymentMethodLabelsStatusFAILED                           = gen.PaymentMethodLabelsStatusFAILED
	PaymentStatusFailed                                       = gen.PaymentStatusFailed
	PaymentStatusNotRequired                                  = gen.PaymentStatusNotRequired
	PaymentStatusPaid                                         = gen.PaymentStatusPaid
	PaymentStatusProcessing                                   = gen.PaymentStatusProcessing
	PaymentStatusRequiresAction                               = gen.PaymentStatusRequiresAction
	PortalEntitlementSourceAddon                              = gen.PortalEntitlementSourceAddon
	PortalEntitlementSourceLicense                            = gen.PortalEntitlementSourceLicense
	PortalEntitlementTypeBOOLEAN                              = gen.PortalEntitlementTypeBOOLEAN
	PortalEntitlementTypeCONFIG                               = gen.PortalEntitlementTypeCONFIG
	PortalEntitlementTypeNUMBER                               = gen.PortalEntitlementTypeNUMBER
	PortalEntitlementTypeNUMBERAICREDIT                       = gen.PortalEntitlementTypeNUMBERAICREDIT
	PortalInstanceSummaryBillingStatusACTIVE                  = gen.PortalInstanceSummaryBillingStatusACTIVE
	PortalInstanceSummaryBillingStatusCANCELED                = gen.PortalInstanceSummaryBillingStatusCANCELED
	PortalInstanceSummaryBillingStatusLessThannil             = gen.PortalInstanceSummaryBillingStatusLessThannil
	PortalInstanceSummaryBillingStatusPASTDUE                 = gen.PortalInstanceSummaryBillingStatusPASTDUE
	PortalInstanceSummaryBillingStatusTRIAL                   = gen.PortalInstanceSummaryBillingStatusTRIAL
	PortalPaymentMethodStatusACTIVE                           = gen.PortalPaymentMethodStatusACTIVE
	PortalPaymentMethodStatusEXPIRED                          = gen.PortalPaymentMethodStatusEXPIRED
	PortalPaymentMethodStatusFAILED                           = gen.PortalPaymentMethodStatusFAILED
	PortalVoucherStatusACTIVE                                 = gen.PortalVoucherStatusACTIVE
	PortalVoucherStatusEXPIRED                                = gen.PortalVoucherStatusEXPIRED
	PortalVoucherStatusREVOKED                                = gen.PortalVoucherStatusREVOKED
	PortalVoucherVoucherTypeENTITLEMENTBOOST                  = gen.PortalVoucherVoucherTypeENTITLEMENTBOOST
	PortalVoucherVoucherTypePRICE                             = gen.PortalVoucherVoucherTypePRICE
	PriceBillingModelFLATFEE                                  = gen.PriceBillingModelFLATFEE
	PriceBillingModelOVERAGE                                  = gen.PriceBillingModelOVERAGE
	PriceBillingModelUSAGEBASED                               = gen.PriceBillingModelUSAGEBASED
	PriceBillingPeriodANNUAL                                  = gen.PriceBillingPeriodANNUAL
	PriceBillingPeriodLessThannil                             = gen.PriceBillingPeriodLessThannil
	PriceBillingPeriodMONTHLY                                 = gen.PriceBillingPeriodMONTHLY
	PriceBillingPeriodQUARTERLY                               = gen.PriceBillingPeriodQUARTERLY
	PriceBillingPeriodSEMIANNUAL                              = gen.PriceBillingPeriodSEMIANNUAL
	PriceBillingTimingADVANCE                                 = gen.PriceBillingTimingADVANCE
	PriceBillingTimingARREARS                                 = gen.PriceBillingTimingARREARS
	PriceStatusACTIVE                                         = gen.PriceStatusACTIVE
	PriceStatusDEPRECATED                                     = gen.PriceStatusDEPRECATED
	ProvenanceAddonOverrideBehaviorADD                        = gen.ProvenanceAddonOverrideBehaviorADD
	ProvenanceAddonOverrideBehaviorMAX                        = gen.ProvenanceAddonOverrideBehaviorMAX
	ProvenanceAddonOverrideBehaviorOVERRIDE                   = gen.ProvenanceAddonOverrideBehaviorOVERRIDE
	ProvenanceBoostModifierTypeADD                            = gen.ProvenanceBoostModifierTypeADD
	ProvenanceBoostModifierTypeMULTIPLY                       = gen.ProvenanceBoostModifierTypeMULTIPLY
	ProvenanceBoostModifierTypeSET                            = gen.ProvenanceBoostModifierTypeSET
	ProvenanceBoostModifierTypeUNLIMITED                      = gen.ProvenanceBoostModifierTypeUNLIMITED
	MATCHED                                                   = gen.MATCHED
	MISMATCH                                                  = gen.MISMATCH
	ProviderRecordStatusDraft                                 = gen.ProviderRecordStatusDraft
	ProviderRecordStatusOpen                                  = gen.ProviderRecordStatusOpen
	ProviderRecordStatusPaid                                  = gen.ProviderRecordStatusPaid
	ProviderRecordStatusUncollectible                         = gen.ProviderRecordStatusUncollectible
	ProviderRecordStatusVoid                                  = gen.ProviderRecordStatusVoid
	ProviderSyncLastSyncStatusFAILED                          = gen.ProviderSyncLastSyncStatusFAILED
	ProviderSyncLastSyncStatusPARTIAL                         = gen.ProviderSyncLastSyncStatusPARTIAL
	ProviderSyncLastSyncStatusSUCCESS                         = gen.ProviderSyncLastSyncStatusSUCCESS
	ProviderSyncProviderKindNOOP                              = gen.ProviderSyncProviderKindNOOP
	ProviderSyncProviderKindSTRIPE                            = gen.ProviderSyncProviderKindSTRIPE
	PublicAddOnEntitlementOverrideBehaviorADD                 = gen.PublicAddOnEntitlementOverrideBehaviorADD
	PublicAddOnEntitlementOverrideBehaviorMAX                 = gen.PublicAddOnEntitlementOverrideBehaviorMAX
	PublicAddOnEntitlementOverrideBehaviorOVERRIDE            = gen.PublicAddOnEntitlementOverrideBehaviorOVERRIDE
	PublicAddOnEntitlementTypeBOOLEAN                         = gen.PublicAddOnEntitlementTypeBOOLEAN
	PublicAddOnEntitlementTypeCONFIG                          = gen.PublicAddOnEntitlementTypeCONFIG
	PublicAddOnEntitlementTypeNUMBER                          = gen.PublicAddOnEntitlementTypeNUMBER
	PublicAddOnPricingTypeCUSTOM                              = gen.PublicAddOnPricingTypeCUSTOM
	PublicAddOnPricingTypeFREE                                = gen.PublicAddOnPricingTypeFREE
	PublicAddOnPricingTypePAID                                = gen.PublicAddOnPricingTypePAID
	PublicPlanEntitlementTypeBOOLEAN                          = gen.PublicPlanEntitlementTypeBOOLEAN
	PublicPlanEntitlementTypeCONFIG                           = gen.PublicPlanEntitlementTypeCONFIG
	PublicPlanEntitlementTypeNUMBER                           = gen.PublicPlanEntitlementTypeNUMBER
	PublicPlanLifecycleStatePUBLISHED                         = gen.PublicPlanLifecycleStatePUBLISHED
	PublicPlanPricingTypeCUSTOM                               = gen.PublicPlanPricingTypeCUSTOM
	PublicPlanPricingTypeFREE                                 = gen.PublicPlanPricingTypeFREE
	PublicPlanPricingTypePAID                                 = gen.PublicPlanPricingTypePAID
	PublicPriceBillingModelFLATFEE                            = gen.PublicPriceBillingModelFLATFEE
	PublicPriceBillingModelOVERAGE                            = gen.PublicPriceBillingModelOVERAGE
	PublicPriceBillingModelUSAGEBASED                         = gen.PublicPriceBillingModelUSAGEBASED
	PublicPriceBillingPeriodANNUAL                            = gen.PublicPriceBillingPeriodANNUAL
	PublicPriceBillingPeriodMONTHLY                           = gen.PublicPriceBillingPeriodMONTHLY
	PublicPriceBillingPeriodQUARTERLY                         = gen.PublicPriceBillingPeriodQUARTERLY
	PublicPriceBillingPeriodSEMIANNUAL                        = gen.PublicPriceBillingPeriodSEMIANNUAL
	PublicPriceBillingTimingADVANCE                           = gen.PublicPriceBillingTimingADVANCE
	PublicPriceBillingTimingARREARS                           = gen.PublicPriceBillingTimingARREARS
	QueuedInvoiceCollectionMethodCHARGEAUTOMATICALLY          = gen.QueuedInvoiceCollectionMethodCHARGEAUTOMATICALLY
	QueuedInvoiceCollectionMethodSENDINVOICE                  = gen.QueuedInvoiceCollectionMethodSENDINVOICE
	QueuedInvoiceHandoffStatusACKNOWLEDGED                    = gen.QueuedInvoiceHandoffStatusACKNOWLEDGED
	QueuedInvoiceHandoffStatusNOTREQUIRED                     = gen.QueuedInvoiceHandoffStatusNOTREQUIRED
	QueuedInvoiceHandoffStatusPENDING                         = gen.QueuedInvoiceHandoffStatusPENDING
	QueuedInvoiceHoldReasonLEDGERCHAINBREAK                   = gen.QueuedInvoiceHoldReasonLEDGERCHAINBREAK
	QueuedInvoiceHoldReasonLEDGERCOUNTERMISMATCH              = gen.QueuedInvoiceHoldReasonLEDGERCOUNTERMISMATCH
	QueuedInvoiceHoldReasonLEDGERSEQUENCEGAP                  = gen.QueuedInvoiceHoldReasonLEDGERSEQUENCEGAP
	QueuedInvoiceHoldReasonLessThannil                        = gen.QueuedInvoiceHoldReasonLessThannil
	QueuedInvoiceKindACTIVATION                               = gen.QueuedInvoiceKindACTIVATION
	QueuedInvoiceKindFINAL                                    = gen.QueuedInvoiceKindFINAL
	QueuedInvoiceKindRENEWAL                                  = gen.QueuedInvoiceKindRENEWAL
	QueuedInvoiceProviderKindNOOP                             = gen.QueuedInvoiceProviderKindNOOP
	QueuedInvoiceProviderKindSTRIPE                           = gen.QueuedInvoiceProviderKindSTRIPE
	QueuedInvoiceStatusDRAFT                                  = gen.QueuedInvoiceStatusDRAFT
	QueuedInvoiceStatusMANUAL                                 = gen.QueuedInvoiceStatusMANUAL
	QueuedInvoiceStatusPAID                                   = gen.QueuedInvoiceStatusPAID
	QueuedInvoiceStatusPAYMENTFAILED                          = gen.QueuedInvoiceStatusPAYMENTFAILED
	QueuedInvoiceStatusPUSHED                                 = gen.QueuedInvoiceStatusPUSHED
	QueuedInvoiceStatusPUSHFAILED                             = gen.QueuedInvoiceStatusPUSHFAILED
	QueuedInvoiceStatusUNCOLLECTIBLE                          = gen.QueuedInvoiceStatusUNCOLLECTIBLE
	QueuedInvoiceStatusVOID                                   = gen.QueuedInvoiceStatusVOID
	RedemptionStatusACTIVE                                    = gen.RedemptionStatusACTIVE
	RedemptionStatusEXPIRED                                   = gen.RedemptionStatusEXPIRED
	RedemptionStatusREVOKED                                   = gen.RedemptionStatusREVOKED
	RedemptionVoucherTypeENTITLEMENTBOOST                     = gen.RedemptionVoucherTypeENTITLEMENTBOOST
	RedemptionVoucherTypePRICE                                = gen.RedemptionVoucherTypePRICE
	SessionBoostModifierTypeADD                               = gen.SessionBoostModifierTypeADD
	SessionBoostModifierTypeMULTIPLY                          = gen.SessionBoostModifierTypeMULTIPLY
	SessionBoostModifierTypeSET                               = gen.SessionBoostModifierTypeSET
	SessionBoostModifierTypeUNLIMITED                         = gen.SessionBoostModifierTypeUNLIMITED
	Preview                                                   = gen.Preview
	RequiresPaymentMethod                                     = gen.RequiresPaymentMethod
	Subscribed                                                = gen.Subscribed
	SessionDiscountAppliesToADDONS                            = gen.SessionDiscountAppliesToADDONS
	SessionDiscountAppliesToBOTH                              = gen.SessionDiscountAppliesToBOTH
	SessionDiscountAppliesToLICENSEBASE                       = gen.SessionDiscountAppliesToLICENSEBASE
	SessionDiscountAppliesToSELECTEDPRICES                    = gen.SessionDiscountAppliesToSELECTEDPRICES
	SessionDiscountTypeFIXEDAMOUNT                            = gen.SessionDiscountTypeFIXEDAMOUNT
	SessionDiscountTypePERCENTAGE                             = gen.SessionDiscountTypePERCENTAGE
	SessionInvoiceKindACTIVATION                              = gen.SessionInvoiceKindACTIVATION
	SessionInvoiceKindFINAL                                   = gen.SessionInvoiceKindFINAL
	SessionInvoiceKindRENEWAL                                 = gen.SessionInvoiceKindRENEWAL
	SessionInvoiceStatusMANUAL                                = gen.SessionInvoiceStatusMANUAL
	SessionInvoiceStatusPAID                                  = gen.SessionInvoiceStatusPAID
	SessionInvoiceStatusPAYMENTFAILED                         = gen.SessionInvoiceStatusPAYMENTFAILED
	SessionInvoiceStatusPUSHED                                = gen.SessionInvoiceStatusPUSHED
	SessionInvoiceStatusUNCOLLECTIBLE                         = gen.SessionInvoiceStatusUNCOLLECTIBLE
	SessionInvoiceStatusVOID                                  = gen.SessionInvoiceStatusVOID
	SessionLineTypeADDON                                      = gen.SessionLineTypeADDON
	SessionLineTypeBASE                                       = gen.SessionLineTypeBASE
	SessionLineTypeDISCOUNT                                   = gen.SessionLineTypeDISCOUNT
	SessionLineTypeOVERAGE                                    = gen.SessionLineTypeOVERAGE
	SessionLineTypeUSAGE                                      = gen.SessionLineTypeUSAGE
	SessionSubscriptionBillingPeriodANNUAL                    = gen.SessionSubscriptionBillingPeriodANNUAL
	SessionSubscriptionBillingPeriodMONTHLY                   = gen.SessionSubscriptionBillingPeriodMONTHLY
	SessionSubscriptionBillingPeriodQUARTERLY                 = gen.SessionSubscriptionBillingPeriodQUARTERLY
	SessionSubscriptionBillingPeriodSEMIANNUAL                = gen.SessionSubscriptionBillingPeriodSEMIANNUAL
	SessionSubscriptionCollectionMethodCHARGEAUTOMATICALLY    = gen.SessionSubscriptionCollectionMethodCHARGEAUTOMATICALLY
	SessionSubscriptionCollectionMethodSENDINVOICE            = gen.SessionSubscriptionCollectionMethodSENDINVOICE
	SessionSubscriptionProviderKindNOOP                       = gen.SessionSubscriptionProviderKindNOOP
	SessionSubscriptionProviderKindSTRIPE                     = gen.SessionSubscriptionProviderKindSTRIPE
	SessionSubscriptionStatusACTIVE                           = gen.SessionSubscriptionStatusACTIVE
	SessionSubscriptionStatusCANCELED                         = gen.SessionSubscriptionStatusCANCELED
	SessionSubscriptionStatusPASTDUE                          = gen.SessionSubscriptionStatusPASTDUE
	SessionSubscriptionStatusTRIAL                            = gen.SessionSubscriptionStatusTRIAL
	SessionVoucherDurationFOREVER                             = gen.SessionVoucherDurationFOREVER
	SessionVoucherDurationONETIME                             = gen.SessionVoucherDurationONETIME
	SessionVoucherDurationREPEATING                           = gen.SessionVoucherDurationREPEATING
	SessionVoucherVoucherTypeENTITLEMENTBOOST                 = gen.SessionVoucherVoucherTypeENTITLEMENTBOOST
	SessionVoucherVoucherTypePRICE                            = gen.SessionVoucherVoucherTypePRICE
	StartedSubscriptionBillingPeriodANNUAL                    = gen.StartedSubscriptionBillingPeriodANNUAL
	StartedSubscriptionBillingPeriodMONTHLY                   = gen.StartedSubscriptionBillingPeriodMONTHLY
	StartedSubscriptionBillingPeriodQUARTERLY                 = gen.StartedSubscriptionBillingPeriodQUARTERLY
	StartedSubscriptionBillingPeriodSEMIANNUAL                = gen.StartedSubscriptionBillingPeriodSEMIANNUAL
	StartedSubscriptionCollectionMethodCHARGEAUTOMATICALLY    = gen.StartedSubscriptionCollectionMethodCHARGEAUTOMATICALLY
	StartedSubscriptionCollectionMethodSENDINVOICE            = gen.StartedSubscriptionCollectionMethodSENDINVOICE
	StartedSubscriptionProviderKindNOOP                       = gen.StartedSubscriptionProviderKindNOOP
	StartedSubscriptionProviderKindSTRIPE                     = gen.StartedSubscriptionProviderKindSTRIPE
	StartedSubscriptionStatusACTIVE                           = gen.StartedSubscriptionStatusACTIVE
	StartedSubscriptionStatusCANCELED                         = gen.StartedSubscriptionStatusCANCELED
	StartedSubscriptionStatusPASTDUE                          = gen.StartedSubscriptionStatusPASTDUE
	StartedSubscriptionStatusTRIAL                            = gen.StartedSubscriptionStatusTRIAL
	SubscriptionCancellationModeATPERIODEND                   = gen.SubscriptionCancellationModeATPERIODEND
	SubscriptionCancellationModeIMMEDIATE                     = gen.SubscriptionCancellationModeIMMEDIATE
	SubscriptionTermsCollectionMethodCHARGEAUTOMATICALLY      = gen.SubscriptionTermsCollectionMethodCHARGEAUTOMATICALLY
	SubscriptionTermsCollectionMethodSENDINVOICE              = gen.SubscriptionTermsCollectionMethodSENDINVOICE
	SubscriptionTermsProviderKindNOOP                         = gen.SubscriptionTermsProviderKindNOOP
	SubscriptionTermsProviderKindSTRIPE                       = gen.SubscriptionTermsProviderKindSTRIPE
	UsageReportBehaviorAppend                                 = gen.UsageReportBehaviorAppend
	UsageReportBehaviorSet                                    = gen.UsageReportBehaviorSet
	ValidityReasonALREADYREDEEMED                             = gen.ValidityReasonALREADYREDEEMED
	ValidityReasonCURRENCYMISMATCH                            = gen.ValidityReasonCURRENCYMISMATCH
	ValidityReasonEXHAUSTED                                   = gen.ValidityReasonEXHAUSTED
	ValidityReasonEXPIRED                                     = gen.ValidityReasonEXPIRED
	ValidityReasonNOTACTIVE                                   = gen.ValidityReasonNOTACTIVE
	ValidityReasonNOTELIGIBLE                                 = gen.ValidityReasonNOTELIGIBLE
	ValidityReasonNOTFOUND                                    = gen.ValidityReasonNOTFOUND
	ValidityReasonNOTYETVALID                                 = gen.ValidityReasonNOTYETVALID
	ADDONNOTAPPLICABLE                                        = gen.ADDONNOTAPPLICABLE
	ANNUALONLY                                                = gen.ANNUALONLY
	FIRSTTIMEONLY                                             = gen.FIRSTTIMEONLY
	LICENSENOTAPPLICABLE                                      = gen.LICENSENOTAPPLICABLE
	MINIMUMSUBSCRIPTIONAMOUNT                                 = gen.MINIMUMSUBSCRIPTIONAMOUNT
	NOTHINGTOBOOST                                            = gen.NOTHINGTOBOOST
	RESTRICTEDCUSTOMER                                        = gen.RESTRICTEDCUSTOMER
	VoucherDraftDurationFOREVER                               = gen.VoucherDraftDurationFOREVER
	VoucherDraftDurationONETIME                               = gen.VoucherDraftDurationONETIME
	VoucherDraftDurationREPEATING                             = gen.VoucherDraftDurationREPEATING
	VoucherDraftPriceAppliesToADDONS                          = gen.VoucherDraftPriceAppliesToADDONS
	VoucherDraftPriceAppliesToBOTH                            = gen.VoucherDraftPriceAppliesToBOTH
	VoucherDraftPriceAppliesToLICENSEBASE                     = gen.VoucherDraftPriceAppliesToLICENSEBASE
	VoucherDraftPriceAppliesToSELECTEDPRICES                  = gen.VoucherDraftPriceAppliesToSELECTEDPRICES
	VoucherDraftPriceDiscountTypeFIXEDAMOUNT                  = gen.VoucherDraftPriceDiscountTypeFIXEDAMOUNT
	VoucherDraftPriceDiscountTypePERCENTAGE                   = gen.VoucherDraftPriceDiscountTypePERCENTAGE
	VoucherDraftVoucherTypeENTITLEMENTBOOST                   = gen.VoucherDraftVoucherTypeENTITLEMENTBOOST
	VoucherDraftVoucherTypePRICE                              = gen.VoucherDraftVoucherTypePRICE
	VoucherDurationFOREVER                                    = gen.VoucherDurationFOREVER
	VoucherDurationONETIME                                    = gen.VoucherDurationONETIME
	VoucherDurationREPEATING                                  = gen.VoucherDurationREPEATING
	VoucherPriceAppliesToADDONS                               = gen.VoucherPriceAppliesToADDONS
	VoucherPriceAppliesToBOTH                                 = gen.VoucherPriceAppliesToBOTH
	VoucherPriceAppliesToLICENSEBASE                          = gen.VoucherPriceAppliesToLICENSEBASE
	VoucherPriceAppliesToSELECTEDPRICES                       = gen.VoucherPriceAppliesToSELECTEDPRICES
	VoucherPriceDiscountTypeFIXEDAMOUNT                       = gen.VoucherPriceDiscountTypeFIXEDAMOUNT
	VoucherPriceDiscountTypePERCENTAGE                        = gen.VoucherPriceDiscountTypePERCENTAGE
	VoucherStatusACTIVE                                       = gen.VoucherStatusACTIVE
	VoucherStatusARCHIVED                                     = gen.VoucherStatusARCHIVED
	VoucherStatusDRAFT                                        = gen.VoucherStatusDRAFT
	VoucherStatusEXHAUSTED                                    = gen.VoucherStatusEXHAUSTED
	VoucherStatusEXPIRED                                      = gen.VoucherStatusEXPIRED
	VoucherVoucherTypeENTITLEMENTBOOST                        = gen.VoucherVoucherTypeENTITLEMENTBOOST
	VoucherVoucherTypePRICE                                   = gen.VoucherVoucherTypePRICE
)
