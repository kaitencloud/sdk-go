// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/kaitencloud/sdk-go/internal/genplatform"
)

func TestNewPlatformClientRejectsABlankBaseURL(t *testing.T) {
	t.Parallel()

	if _, err := NewPlatformClient("   "); err == nil {
		t.Fatal("NewPlatformClient() with a blank base URL returned no error")
	}
}

func TestNewPlatformClientInitializesEveryResourceNamespace(t *testing.T) {
	t.Parallel()

	client, err := NewPlatformClient(platformBaseURL)
	if err != nil {
		t.Fatalf("NewPlatformClient() error = %v", err)
	}

	namespaces := map[string]any{
		"Connectors":    client.Connectors,
		"Organizations": client.Organizations,
		"Tokens":        client.Tokens,
		"Users":         client.Users,
	}

	for name, namespace := range namespaces {
		if reflect.ValueOf(namespace).IsNil() {
			t.Errorf("%s namespace was not initialized", name)
		}
	}
}

// One Option type serves both constructors, so the shared default has to hold on both
// sides -- a PlatformClient built with no options must be as bounded as a Client is.
func TestNewPlatformClientDefaultsHTTPTimeout(t *testing.T) {
	t.Parallel()

	client, err := NewPlatformClient(platformBaseURL)
	if err != nil {
		t.Fatalf("NewPlatformClient() error = %v", err)
	}

	if timeout := platformHTTPClientTimeout(t, client); timeout != DefaultTimeout {
		t.Fatalf("default HTTP client Timeout = %v, want %v", timeout, DefaultTimeout)
	}
}

func TestNewPlatformClientLetsWithHTTPClientReplaceTheDefaultTimeout(t *testing.T) {
	t.Parallel()

	want := 5 * time.Second

	client, err := NewPlatformClient(platformBaseURL, WithHTTPClient(&http.Client{Timeout: want}))
	if err != nil {
		t.Fatalf("NewPlatformClient() error = %v", err)
	}

	if timeout := platformHTTPClientTimeout(t, client); timeout != want {
		t.Fatalf("HTTP client Timeout = %v, want %v", timeout, want)
	}
}

func TestPlatformClientExposesOnlyTheReviewedPublicSurface(t *testing.T) {
	t.Parallel()

	client, err := NewPlatformClient(platformBaseURL)
	if err != nil {
		t.Fatalf("NewPlatformClient() error = %v", err)
	}

	assertServiceCoverage(t, map[string]serviceCoverage{
		"Connectors": {
			service: client.Connectors,
			methods: []string{"Register"},
		},
		"Organizations": {
			service: client.Organizations,
			methods: []string{"Ensure", "Get", "Delete", "DeleteMembership"},
		},
		"Tokens": {
			service: client.Tokens,
			methods: []string{"Mint", "List", "Revoke"},
		},
		"Users": {
			service: client.Users,
			methods: []string{"Delete"},
		},
	})
}

// The /platform prefix is the whole reason this client exists: the previous wrappers
// rendered /organizations/{id} and were sent to a listener that refuses the credential
// they carry. Assert the wire path, not just that the call compiles.
func TestPlatformRequestsCarryPlatformPathPrefix(t *testing.T) {
	t.Parallel()

	const (
		organizationID = "3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31"
		userID         = "8c1e4d0b-9d5f-4a2c-8b0a-1c4e8d312b0a"
	)

	tests := []struct {
		name       string
		call       func(context.Context, *PlatformClient) error
		wantMethod string
		wantPath   string
	}{
		{
			name:       "describe the credential",
			call:       func(ctx context.Context, c *PlatformClient) error { _, err := c.Me(ctx); return err },
			wantMethod: http.MethodGet,
			wantPath:   "/api/platform/me",
		},
		{
			name: "register connector",
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Connectors.Register(ctx, RegisterConnectorInput{Name: "kaiten.integration.crm.attio", Version: "1.0.0"})
				return err
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/platform/connectors",
		},
		{
			name: "ensure organization",
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Organizations.Ensure(ctx, EnsureOrganizationInput{ExternalID: "org_2abcDEF"})
				return err
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/platform/organizations",
		},
		{
			name: "get organization",
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Organizations.Get(ctx, organizationID)
				return err
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/platform/organizations/" + organizationID,
		},
		{
			name: "delete membership",
			call: func(ctx context.Context, c *PlatformClient) error {
				return c.Organizations.DeleteMembership(ctx, organizationID, userID)
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/api/platform/organizations/" + organizationID + "/memberships/" + userID,
		},
		{
			name: "mint token",
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Tokens.Mint(ctx, organizationID, MintTokenInput{Name: "zone-kaiten"})
				return err
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/platform/organizations/" + organizationID + "/tokens",
		},
		{
			name: "revoke token",
			call: func(ctx context.Context, c *PlatformClient) error {
				return c.Tokens.Revoke(ctx, organizationID, "zone-kaiten")
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/api/platform/organizations/" + organizationID + "/tokens/zone-kaiten",
		},
		{
			name:       "delete user",
			call:       func(ctx context.Context, c *PlatformClient) error { return c.Users.Delete(ctx, userID) },
			wantMethod: http.MethodDelete,
			wantPath:   "/api/platform/users/" + userID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenPlatformClient(t, respondNoContent(), WithBearerToken("ksm_secret"))

			// A 204 satisfies the no-content wrappers and fails the JSON ones; either way the
			// request was already built, which is what this test is about.
			_ = test.call(t.Context(), client)

			got := sent.only(t)

			if got.method != test.wantMethod {
				t.Errorf("method = %s, want %s", got.method, test.wantMethod)
			}
			if got.path != test.wantPath {
				t.Errorf("path = %s, want %s", got.path, test.wantPath)
			}
			if authorization := got.header.Get("Authorization"); authorization != "Bearer ksm_secret" {
				t.Errorf("Authorization = %q, want %q", authorization, "Bearer ksm_secret")
			}
		})
	}
}

// The Platform API's Problem carries a machine-readable code and a log correlation id the
// Core API's ErrorModel has no field for. Both have to survive the conversion, or callers
// are left branching on prose.
func TestPlatformResponseErrorCarriesCodeAndErrorID(t *testing.T) {
	t.Parallel()

	payload := `{"title":"Conflict","detail":"a token with that name is already active","code":"Token.NameConflict","errorId":"3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31","errors":[{"location":"body.name","message":"already taken"}]}`

	client, _ := givenPlatformClient(t, respondJSON(http.StatusConflict, payload))

	_, mintErr := client.Tokens.Mint(t.Context(), "3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31", MintTokenInput{Name: "zone-kaiten"})

	var apiErr *Error
	if !errors.As(mintErr, &apiErr) {
		t.Fatalf("Mint() error is %T, want *Error", mintErr)
	}

	if apiErr.Code != "Token.NameConflict" {
		t.Errorf("Code = %q, want Token.NameConflict", apiErr.Code)
	}
	if apiErr.ErrorID != "3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31" {
		t.Errorf("ErrorID = %q, want the correlation id", apiErr.ErrorID)
	}
	if apiErr.Problem == nil || stringValue(apiErr.Problem.Detail) != "a token with that name is already active" {
		t.Fatalf("Problem = %+v, want the detail preserved", apiErr.Problem)
	}
	if apiErr.Problem.Errors == nil || len(*apiErr.Problem.Errors) != 1 {
		t.Errorf("Problem.Errors = %v, want one detail", apiErr.Problem.Errors)
	}

	want := "kaiten API error (409 409 Conflict) [Token.NameConflict]: Conflict: a token with that name is already active (errorId: 3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31)"
	if apiErr.Error() != want {
		t.Errorf("Error() = %q, want %q", apiErr.Error(), want)
	}
}

func platformHTTPClientTimeout(t *testing.T, client *PlatformClient) time.Duration {
	t.Helper()

	raw, ok := client.raw.ClientInterface.(*genplatform.Client)
	if !ok {
		t.Fatalf("underlying client is %T, want *genplatform.Client", client.raw.ClientInterface)
	}

	doer, ok := raw.Client.(*http.Client)
	if !ok {
		t.Fatalf("request doer is %T, want *http.Client", raw.Client)
	}

	return doer.Timeout
}

// TestPlatformClientWrappersReportEveryDeclaredProblem is the Platform half of
// TestClientWrappersReportEveryDeclaredProblem. Every wrapper here passed its full set
// when the Core ones were audited; this is what keeps it so once the Platform spec
// declares a status on an operation that it does not declare today.
func TestPlatformClientWrappersReportEveryDeclaredProblem(t *testing.T) {
	t.Parallel()

	const (
		organizationID = "3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31"
		userID         = "8c1e4d0b-9d5f-4a2c-8b0a-1c4e8d312b0a"
	)

	given := func(t *testing.T, respond responder) *PlatformClient {
		t.Helper()

		client, _ := givenPlatformClient(t, respond)
		return client
	}

	assertProblemCoverage(t, given, map[string]problemCoverage[*PlatformClient]{
		"PlatformClient.Me": {
			response: genplatform.GetPlatformMeResponse{},
			call:     func(ctx context.Context, c *PlatformClient) error { _, err := c.Me(ctx); return err },
		},
		"PlatformConnectors.Register": {
			response: genplatform.RegisterConnectorResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Connectors.Register(ctx, RegisterConnectorInput{Name: "kaiten.integration.crm.attio", Version: "1.0.0"})
				return err
			},
		},
		"PlatformOrganizations.Ensure": {
			response: genplatform.EnsureOrganizationResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Organizations.Ensure(ctx, EnsureOrganizationInput{ExternalID: "org_2abcDEF"})
				return err
			},
		},
		"PlatformOrganizations.Get": {
			response: genplatform.GetOrganizationResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Organizations.Get(ctx, organizationID)
				return err
			},
		},
		"PlatformOrganizations.Delete": {
			response: genplatform.DeleteOrganizationResponse{},
			call:     func(ctx context.Context, c *PlatformClient) error { return c.Organizations.Delete(ctx, organizationID) },
		},
		"PlatformOrganizations.DeleteMembership": {
			response: genplatform.DeleteMembershipResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				return c.Organizations.DeleteMembership(ctx, organizationID, userID)
			},
		},
		"PlatformTokens.Mint": {
			response: genplatform.MintOrganizationTokenResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Tokens.Mint(ctx, organizationID, MintTokenInput{Name: "zone-kaiten"})
				return err
			},
		},
		"PlatformTokens.Revoke": {
			response: genplatform.RevokeOrganizationTokenResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				return c.Tokens.Revoke(ctx, organizationID, "zone-kaiten")
			},
		},
		"PlatformTokens.List": {
			response: genplatform.ListOrganizationTokensResponse{},
			call: func(ctx context.Context, c *PlatformClient) error {
				_, err := c.Tokens.List(ctx, organizationID)
				return err
			},
		},
		"PlatformUsers.Delete": {
			response: genplatform.DeleteUserResponse{},
			call:     func(ctx context.Context, c *PlatformClient) error { return c.Users.Delete(ctx, userID) },
		},
	})
}
