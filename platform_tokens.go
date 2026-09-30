package sdk

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Mint issues an organization-scoped token (`ksh_...`) for the organization identified by
// organizationID. The returned MintedToken is the only time its Token field is populated:
// the API stores a hash, so a lost plaintext cannot be recovered, only superseded.
//
// The minted token's scopes must be a subset of the calling platform credential's, and
// its lifetime is bounded by that credential in a way TTL does not express: revoking the
// platform credential revokes everything it minted. Plan rotation as mint-new, adopt,
// then revoke-old -- never revoke-then-mint.
func (s *PlatformTokens) Mint(ctx context.Context, organizationID string, input MintTokenInput) (MintedToken, error) {
	parsedOrgID, err := uuid.Parse(organizationID)
	if err != nil {
		var zero MintedToken
		return zero, fmt.Errorf("parse organization id: %w", err)
	}

	body, err := jsonBody(input.payload())
	if err != nil {
		var zero MintedToken
		return zero, fmt.Errorf("encode mint token request: %w", err)
	}

	resp, err := s.client.raw.MintOrganizationTokenWithBodyWithResponse(ctx, openapi_types.UUID(parsedOrgID), contentTypeJSON, body)
	if err != nil {
		var zero MintedToken
		return zero, fmt.Errorf("mint organization token: %w", err)
	}

	return expectJSON(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON201)
}

// Revoke revokes the organization token identified by tokenSlug, which is the Slug of a
// MintedToken -- not its plaintext and not its name. Revocation is immediate and final;
// a replacement is a fresh Mint.
func (s *PlatformTokens) Revoke(ctx context.Context, organizationID, tokenSlug string) error {
	parsedOrgID, err := uuid.Parse(organizationID)
	if err != nil {
		return fmt.Errorf("parse organization id: %w", err)
	}

	resp, err := s.client.raw.RevokeOrganizationTokenWithResponse(ctx, openapi_types.UUID(parsedOrgID), tokenSlug)
	if err != nil {
		return fmt.Errorf("revoke organization token: %w", err)
	}

	return expectNoContent(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// List names the still-active organization tokens this platform credential minted
// inside organizationID, so one can be addressed for revocation.
//
// It exists because Revoke takes a Slug and a slug is server-generated with a
// random suffix: a caller holds the Name it chose and nothing else, and minting
// a replacement beside a live predecessor of the same name is refused. Without
// this, a credential was creatable by its owner and retirable by nobody.
//
// Only credentials issued by the CALLING platform credential, in the named
// organization, and not yet revoked. No value is returned and none could be --
// the plaintext exists once, at Mint. These identify a credential; they do not
// authenticate as one.
func (s *PlatformTokens) List(ctx context.Context, organizationID string) ([]IssuedToken, error) {
	parsedOrgID, err := uuid.Parse(organizationID)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	resp, err := s.client.raw.ListOrganizationTokensWithResponse(ctx, openapi_types.UUID(parsedOrgID))
	if err != nil {
		return nil, fmt.Errorf("list organization tokens: %w", err)
	}

	tokens, err := expectJSON(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
