// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// List returns all licenses.
func (s *Licenses) List(ctx context.Context) ([]License, error) {
	return listAll[License](ctx, s.client.list, listPath("licenses"), nil)
}

// Get returns the license identified by licenseSlug.
func (s *Licenses) Get(ctx context.Context, licenseSlug string) (License, error) {
	resp, err := s.client.raw.GetLicenseWithResponse(ctx, licenseSlug)
	if err != nil {
		var zero License
		return zero, fmt.Errorf("get license: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new license.
//
// With neither input.FamilySlug nor input.FamilyID set it opens a new family, of
// which the license is version 1. With one of them it adds the next version to
// that family -- the version number is server-assigned, so input.Version is
// ignored. Two families may share a name: the family, not the name, is the
// product. The family rules are documented on LicenseInput, and each one arrives
// as an *Error whose Code names it: CreateLicense.FamilyNotFound on a 404,
// CreateLicense.FamilyMismatch on a 422, and CreateLicense.DefaultMustBePublished
// on a 409, for a default created in any state but PUBLISHED.
func (s *Licenses) Create(ctx context.Context, input LicenseInput) (License, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero License
		return zero, fmt.Errorf("encode license request: %w", err)
	}

	resp, err := s.client.raw.CreateLicenseWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero License
		return zero, fmt.Errorf("create license: %w", err)
	}

	// 404, 409 and 503 are declared on this operation and were not passed
	// through, so the errors the family fields make reachable -- an unknown
	// family, an unpublished default -- arrived as an *Error with no Problem and
	// an empty Code, the one member callers are told to branch on.
	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the license identified by licenseSlug.
//
// It never changes the version's lifecycle state: input.LifecycleState is not
// sent, and the state moves through Publish, Archive and Unarchive. Nor is
// input.FamilySlug sent -- see LicenseInput for both. IsDefault true on a version
// that is not published is refused with a 409,
// UpdateLicense.DefaultMustBePublished: publish or unarchive it first. IsDefault
// false on the family's default unsets it, which is also the way to archive the
// version afterwards.
func (s *Licenses) Update(ctx context.Context, licenseSlug string, input LicenseInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode license request: %w", err)
	}

	resp, err := s.client.raw.UpdateLicenseWithBodyWithResponse(ctx, licenseSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update license: %w", err)
	}

	// 409 was declared here and not passed through either; see Create.
	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Publish puts the draft version identified by licenseSlug on sale and returns
// it, now PUBLISHED. It can then be made its family's default, and a family
// without a default serves its highest-numbered published version, which may be
// this one.
//
// Only a draft can be published. Any other state is a 409,
// PublishLicense.NotADraft: an archived version goes back on sale through
// Unarchive. An unknown slug is a 404, PublishLicense.NotFound.
func (s *Licenses) Publish(ctx context.Context, licenseSlug string) (License, error) {
	resp, err := s.client.raw.PublishLicenseWithResponse(ctx, licenseSlug)
	if err != nil {
		var zero License
		return zero, fmt.Errorf("publish license: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Archive withdraws the published version identified by licenseSlug from sale
// and returns it, now ARCHIVED. Its family stops serving it and no instance can
// be assigned to it any more; instances already on it keep it.
//
// Only a published version can be archived. Any other state is a 409,
// ArchiveLicense.NotPublished: a draft that was never on sale is deleted
// instead. The family's default cannot be archived either, a 409 with
// ArchiveLicense.DefaultMustBePublished: make another version the default, or
// unset IsDefault through Update, first. An unknown slug is a 404,
// ArchiveLicense.NotFound.
func (s *Licenses) Archive(ctx context.Context, licenseSlug string) (License, error) {
	resp, err := s.client.raw.ArchiveLicenseWithResponse(ctx, licenseSlug)
	if err != nil {
		var zero License
		return zero, fmt.Errorf("archive license: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Unarchive puts the archived version identified by licenseSlug back on sale and
// returns it, PUBLISHED again: it can be assigned to instances and made its
// family's default once more.
//
// Only an archived version can be unarchived. Any other state is a 409,
// UnarchiveLicense.NotArchived: a draft goes on sale through Publish. An unknown
// slug is a 404, UnarchiveLicense.NotFound.
func (s *Licenses) Unarchive(ctx context.Context, licenseSlug string) (License, error) {
	resp, err := s.client.raw.UnarchiveLicenseWithResponse(ctx, licenseSlug)
	if err != nil {
		var zero License
		return zero, fmt.Errorf("unarchive license: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Delete deletes the license identified by licenseSlug.
func (s *Licenses) Delete(ctx context.Context, licenseSlug string) error {
	resp, err := s.client.raw.DeleteLicenseWithResponse(ctx, licenseSlug)
	if err != nil {
		return fmt.Errorf("delete license: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// ListEntitlements returns the entitlements granted by the license identified by licenseSlug.
func (s *Licenses) ListEntitlements(ctx context.Context, licenseSlug string) ([]LicenseEntitlement, error) {
	return listAll[LicenseEntitlement](ctx, s.client.list, listPath("licenses", licenseSlug, "entitlements"), nil)
}

// AssociateEntitlement grants the entitlement identified by entitlementSlug to the license identified by licenseSlug with the given value.
//
// A numeric grant enforces its value as a hard limit unless WithOveragePercent gives it
// an allowance; an unlimited value is never enforced. A license grants an entitlement
// once: a second grant is a 409, AssociateEntitlementToLicense.AlreadyAssociated, and
// UpdateEntitlement is how a grant changes.
func (s *Licenses) AssociateEntitlement(ctx context.Context, licenseSlug, entitlementSlug string, value LicenseEntitlementValue, opts ...LicenseEntitlementOption) error {
	options := newLicenseEntitlementOptions(opts)

	body, err := jsonBody(licenseEntitlementCreatePayload{
		EntitlementSlug:                entitlementSlug,
		Value:                          value,
		LimitCapExceededOveragePercent: options.overagePercent,
	})
	if err != nil {
		return fmt.Errorf("encode license entitlement request: %w", err)
	}

	resp, err := s.client.raw.AssociateEntitlementWithLicenseWithBodyWithResponse(ctx, licenseSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("associate entitlement with license: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503))
}

// GetEntitlement returns the entitlement identified by entitlementSlug granted by the license identified by licenseSlug.
func (s *Licenses) GetEntitlement(ctx context.Context, licenseSlug, entitlementSlug string) (LicenseEntitlement, error) {
	resp, err := s.client.raw.GetLicenseEntitlementWithResponse(ctx, licenseSlug, entitlementSlug)
	if err != nil {
		var zero LicenseEntitlement
		return zero, fmt.Errorf("get license entitlement: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// UpdateEntitlement updates the value of the entitlement identified by entitlementSlug granted by the license identified by licenseSlug.
//
// It replaces the grant rather than patching it. A numeric grant written without
// WithOveragePercent gets the allowance the server derives from the value -- none, for a
// finite one -- so changing only the value of a soft limit makes it a hard one. To keep
// the allowance, pass back the LimitCapExceededOveragePercent that GetEntitlement reads,
// unless the new value is unlimited, which takes -1.
func (s *Licenses) UpdateEntitlement(ctx context.Context, licenseSlug, entitlementSlug string, value LicenseEntitlementValue, opts ...LicenseEntitlementOption) error {
	options := newLicenseEntitlementOptions(opts)

	body, err := jsonBody(licenseEntitlementUpdatePayload{
		Value:                          value,
		LimitCapExceededOveragePercent: options.overagePercent,
	})
	if err != nil {
		return fmt.Errorf("encode license entitlement request: %w", err)
	}

	resp, err := s.client.raw.UpdateLicenseEntitlementWithBodyWithResponse(ctx, licenseSlug, entitlementSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update license entitlement: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// DeleteEntitlement removes the entitlement identified by entitlementSlug from the license identified by licenseSlug.
func (s *Licenses) DeleteEntitlement(ctx context.Context, licenseSlug, entitlementSlug string) error {
	resp, err := s.client.raw.DeleteLicenseEntitlementWithResponse(ctx, licenseSlug, entitlementSlug)
	if err != nil {
		return fmt.Errorf("delete license entitlement: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
