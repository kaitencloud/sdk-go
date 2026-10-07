// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/kaitencloud/sdk-go/internal/genplatform"
)

const testOrganizationID = "3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31"

func TestMintSendsNameScopesAndTTL(t *testing.T) {
	t.Parallel()

	ttl := time.Hour

	client, sent := givenPlatformClient(t, respondNoContent())

	_, _ = client.Tokens.Mint(t.Context(), testOrganizationID, MintTokenInput{
		Name:   "zone-kaiten",
		Scopes: []string{"write:instances", "read:feature_flags"},
		TTL:    &ttl,
	})

	var body genplatform.MintOrganizationTokenBody
	raw := sent.only(t).body
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal request body %q: %v", raw, err)
	}

	if body.Name != "zone-kaiten" {
		t.Errorf("name = %q, want %q", body.Name, "zone-kaiten")
	}

	if body.Scopes == nil || !slices.Equal(*body.Scopes, []string{"write:instances", "read:feature_flags"}) {
		t.Errorf("scopes = %v, want [write:instances read:feature_flags]", body.Scopes)
	}

	// time.ParseDuration reads back what Duration.String() emits, so the server accepts
	// "1h0m0s" as readily as "1h".
	if body.Ttl == nil || *body.Ttl != "1h0m0s" {
		t.Errorf("ttl = %v, want 1h0m0s", body.Ttl)
	}
}

func TestMintOmitsTTLWhenNotSet(t *testing.T) {
	t.Parallel()

	client, sent := givenPlatformClient(t, respondNoContent())

	_, _ = client.Tokens.Mint(t.Context(), testOrganizationID, MintTokenInput{Name: "zone-kaiten"})

	// An omitted ttl is what mints a non-expiring credential. Sending "0s" instead would
	// mint one that expires immediately, so the absence has to be preserved on the wire.
	if body := sent.only(t).body; bytes.Contains(body, []byte("ttl")) {
		t.Fatalf("body = %s, want no ttl key", body)
	}
}

func TestMintReturnsPlaintextToken(t *testing.T) {
	t.Parallel()

	payload := `{"name":"zone-kaiten","slug":"zone-kaiten","token":"ksh_plaintext","createdAt":"2026-01-01T00:00:00Z","createdBy":{"id":"3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31"},"scopes":["write:instances"]}`

	client, _ := givenPlatformClient(t, respondJSON(http.StatusCreated, payload))

	minted, err := client.Tokens.Mint(t.Context(), testOrganizationID, MintTokenInput{Name: "zone-kaiten"})
	if err != nil {
		t.Fatalf("Mint() error = %v", err)
	}

	// The plaintext token is shown once, on this response and nowhere else, so a
	// caller that does not receive it here has no way to obtain the credential at all.
	if minted.Token == nil || *minted.Token != "ksh_plaintext" {
		t.Errorf("Token = %v, want ksh_plaintext", minted.Token)
	}

	// Revoke takes the slug, not the name, so losing it strands the credential.
	if minted.Slug == nil || *minted.Slug != "zone-kaiten" {
		t.Errorf("Slug = %v, want zone-kaiten", minted.Slug)
	}
}

// The whole reason List exists: a caller holds the NAME it minted under, and
// revocation is addressed by the SLUG, which is server-generated with a random
// suffix and cannot be derived from the name.
func TestListReturnsTheSlugRevocationIsAddressedBy(t *testing.T) {
	t.Parallel()

	payload := `[{"name":"zone-onboarding","slug":"system-kaiten-2c8293",` +
		`"createdAt":"2026-01-01T00:00:00Z","createdBy":{"id":"3f6b1a2c-7c1e-4d0b-9d5f-2b0a1c4e8d31"},` +
		`"scopes":["write:instances","read:deployment_zones"]}]`

	client, rec := givenPlatformClient(t, respondJSON(http.StatusOK, payload))

	tokens, err := client.Tokens.List(t.Context(), testOrganizationID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	sent := rec.only(t)
	if sent.method != http.MethodGet {
		t.Errorf("method = %s, want GET", sent.method)
	}
	if !strings.HasSuffix(sent.path, "/platform/organizations/"+testOrganizationID+"/tokens") {
		t.Errorf("path = %s, want the organization's tokens collection", sent.path)
	}

	if len(tokens) != 1 {
		t.Fatalf("len(tokens) = %d, want 1", len(tokens))
	}

	if tokens[0].Name != "zone-onboarding" {
		t.Errorf("Name = %q, want zone-onboarding", tokens[0].Name)
	}
	if tokens[0].Slug == nil || *tokens[0].Slug != "system-kaiten-2c8293" {
		t.Errorf("Slug = %v, want system-kaiten-2c8293", tokens[0].Slug)
	}

	// The scopes are what a caller compares against to decide whether the
	// credential it holds still matches what its consumer needs.
	if tokens[0].Scopes == nil || !slices.Equal(*tokens[0].Scopes, []string{"write:instances", "read:deployment_zones"}) {
		t.Errorf("Scopes = %v, want the two granted scopes", tokens[0].Scopes)
	}
}

// Holding nothing in an organization is an answer, not a failure: it means the
// name is free and a mint will be accepted.
func TestListReturnsEmptyRatherThanErroringWhenNothingIsHeld(t *testing.T) {
	t.Parallel()

	client, _ := givenPlatformClient(t, respondJSON(http.StatusOK, `[]`))

	tokens, err := client.Tokens.List(t.Context(), testOrganizationID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(tokens) != 0 {
		t.Errorf("len(tokens) = %d, want 0", len(tokens))
	}
}

// A refusal must not read as "holds nothing" -- a caller that treated it that
// way would mint a duplicate credential every time the API was unreachable.
func TestListReportsARefusalRatherThanAnEmptyResult(t *testing.T) {
	t.Parallel()

	client, _ := givenPlatformClient(t, respondJSON(http.StatusForbidden,
		`{"title":"Forbidden","status":403,"detail":"missing required scope: read:tokens"}`))

	if _, err := client.Tokens.List(t.Context(), testOrganizationID); err == nil {
		t.Fatal("List() error = nil, want the refusal to surface")
	}
}
