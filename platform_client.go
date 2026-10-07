// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/kaitencloud/sdk-go/internal/genplatform"
)

// PlatformClient is a typed client for the Kaiten Platform API. Construct one with
// NewPlatformClient.
//
// It is a separate type from Client, not a namespace on it, because the two are
// separate deployments of a contract: the Platform API listens on its own port, is
// described by its own OpenAPI document, and accepts only a platform credential
// (`ksm_...`) -- which kaiten's Core listener rejects before it parses the request.
// One Client pointed at one base URL cannot serve both.
//
// The zero value is not usable: the namespaces below are nil until NewPlatformClient
// populates them.
type PlatformClient struct {
	raw *genplatform.ClientWithResponses

	Connectors    *PlatformConnectors
	Organizations *PlatformOrganizations
	Tokens        *PlatformTokens
	Users         *PlatformUsers
}

// NewPlatformClient creates a PlatformClient for the Kaiten Platform API at baseURL,
// applying the given Options. baseURL includes the API's path prefix and points at the
// Platform listener, e.g. https://kaiten.example.com:3001/api -- the Core API's own
// base URL will answer 404 for every operation here.
func NewPlatformClient(baseURL string, opts ...Option) (*PlatformClient, error) {
	resolved := resolveOptions(baseURL, opts)
	if strings.TrimSpace(resolved.baseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	raw, err := genplatform.NewClientWithResponses(resolved.baseURL, resolved.platformOptions()...)
	if err != nil {
		return nil, fmt.Errorf("create platform SDK client: %w", err)
	}

	return newPlatformClient(raw), nil
}

func newPlatformClient(raw *genplatform.ClientWithResponses) *PlatformClient {
	client := &PlatformClient{
		raw: raw,
	}

	client.Connectors = &PlatformConnectors{client: client}
	client.Organizations = &PlatformOrganizations{client: client}
	client.Tokens = &PlatformTokens{client: client}
	client.Users = &PlatformUsers{client: client}

	return client
}

// Me describes the platform credential the client is authenticating with: its name,
// scopes, expiry and the platform identity it belongs to. The token value is never
// returned -- it is shown once, at creation.
//
// Cheap enough to use as a liveness-and-validity probe for a stored credential, which
// is what it is for: there is no other way to tell an expired or revoked credential
// from a working one without attempting a write.
func (c *PlatformClient) Me(ctx context.Context) (PlatformCredential, error) {
	resp, err := c.raw.GetPlatformMeWithResponse(ctx)
	if err != nil {
		var zero PlatformCredential
		return zero, fmt.Errorf("describe platform credential: %w", err)
	}

	return expectJSON(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Raw returns the underlying generated Platform API client for low-level access.
func (c *PlatformClient) Raw() *genplatform.ClientWithResponses {
	return c.raw
}
