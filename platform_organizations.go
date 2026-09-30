package sdk

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Ensure creates the organization for an identity provider's organization id, or returns
// the existing one unchanged. The row's id is derived from the external id, so the same
// external id always converges on the same organization -- including one the API created
// just-in-time from a login. Name is only applied to an organization this call creates;
// it never renames an existing one.
//
// This is the one write in the Platform API that is safe to call from a bootstrapper on
// every run, and it is why configuring a deployment does not need database access.
func (s *PlatformOrganizations) Ensure(ctx context.Context, input EnsureOrganizationInput) (Organization, error) {
	body, err := jsonBody(input.payload())
	if err != nil {
		var zero Organization
		return zero, fmt.Errorf("encode ensure organization request: %w", err)
	}

	resp, err := s.client.raw.EnsureOrganizationWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Organization
		return zero, fmt.Errorf("ensure organization: %w", err)
	}

	return expectJSON(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Get returns the organization identified by id.
func (s *PlatformOrganizations) Get(ctx context.Context, id string) (Organization, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		var zero Organization
		return zero, fmt.Errorf("parse organization id: %w", err)
	}

	resp, err := s.client.raw.GetOrganizationWithResponse(ctx, openapi_types.UUID(parsedID))
	if err != nil {
		var zero Organization
		return zero, fmt.Errorf("get organization: %w", err)
	}

	return expectJSON(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Delete soft-deletes the organization identified by id, along with its remaining
// memberships. Requires the caller's credential to hold the delete:organizations scope --
// a privileged, cross-organization operation, not implied by write:organizations.
func (s *PlatformOrganizations) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse organization id: %w", err)
	}

	resp, err := s.client.raw.DeleteOrganizationWithResponse(ctx, openapi_types.UUID(parsedID))
	if err != nil {
		return fmt.Errorf("delete organization: %w", err)
	}

	return expectNoContent(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// DeleteMembership soft-deletes a single user's membership on an organization. Requires
// delete:memberships -- like Delete, a privileged, cross-organization operation.
func (s *PlatformOrganizations) DeleteMembership(ctx context.Context, organizationID, userID string) error {
	parsedOrgID, err := uuid.Parse(organizationID)
	if err != nil {
		return fmt.Errorf("parse organization id: %w", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	resp, err := s.client.raw.DeleteMembershipWithResponse(ctx, openapi_types.UUID(parsedOrgID), openapi_types.UUID(parsedUserID))
	if err != nil {
		return fmt.Errorf("delete membership: %w", err)
	}

	return expectNoContent(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
