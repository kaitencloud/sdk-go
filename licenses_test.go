// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

// TestCreateRendersTheFamilyFieldsOnlyWhenSet pins both halves of each family
// field on the create body. Absence is meaningful for all three -- no family
// named opens a new one, and no lifecycleState is the server's PUBLISHED -- so a
// key rendered from an unset field would turn every create that never mentioned
// a family into a lookup of the family "" and a 404. Whole bodies are compared
// rather than key sets, because "unchanged for callers who name no family" is a
// claim about the bytes on the wire.
func TestCreateRendersTheFamilyFieldsOnlyWhenSet(t *testing.T) {
	t.Parallel()

	familySlug := "premium"
	familyID := "fam_1"
	draft := LicenseLifecycleDraft

	newFamily := LicenseInput{
		Name:        "Premium",
		Description: "d",
		Type:        LicenseType("PAID"),
	}

	bySlug := newFamily
	bySlug.FamilySlug = &familySlug

	byIDAsDraft := newFamily
	byIDAsDraft.FamilyID = &familyID
	byIDAsDraft.LifecycleState = &draft

	tests := []struct {
		name     string
		input    LicenseInput
		expected string
	}{
		{
			name:     "no family named opens a new one, published",
			input:    newFamily,
			expected: `{"name":"Premium","description":"d","type":"PAID","isDefault":false}`,
		},
		{
			name:     "family named by slug",
			input:    bySlug,
			expected: `{"name":"Premium","description":"d","type":"PAID","isDefault":false,"familySlug":"premium"}`,
		},
		{
			name:     "family named by id, as a draft",
			input:    byIDAsDraft,
			expected: `{"name":"Premium","description":"d","type":"PAID","isDefault":false,"familyId":"fam_1","lifecycleState":"DRAFT"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondToWrite)

			if _, err := client.Licenses.Create(t.Context(), test.input); err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			request := sent.only(t)
			if request.method != http.MethodPost || request.path != "/api/licenses" {
				t.Errorf("request = %s %s, want POST /api/licenses", request.method, request.path)
			}
			if got := string(request.body); got != test.expected {
				t.Errorf("request body = %s, want %s", got, test.expected)
			}
		})
	}
}

// TestUpdateNeverSendsFamilySlugOrLifecycleState pins the asymmetry that makes
// updatePayload its own type. familySlug is create-only, and PUT refuses it with
// UpdateLicense.FamilyNotReassignable rather than ignore it. lifecycleState
// moves through Publish, Archive and Unarchive, and PUT refuses any state but
// the stored one with UpdateLicense.LifecycleStateNotSettable -- so an input
// created as a draft would fail every update once the version is published.
// Both have to be lost on the way out. familyId does go through when set, and
// stays absent when not.
func TestUpdateNeverSendsFamilySlugOrLifecycleState(t *testing.T) {
	t.Parallel()

	familySlug := "premium"
	familyID := "fam_1"
	draft := LicenseLifecycleDraft

	untouched := LicenseInput{
		Name:        "Premium",
		Description: "d",
		Type:        LicenseType("PAID"),
		Version:     "2",
	}

	createdAsDraft := untouched
	createdAsDraft.FamilySlug = &familySlug
	createdAsDraft.FamilyID = &familyID
	createdAsDraft.LifecycleState = &draft

	tests := []struct {
		name     string
		input    LicenseInput
		expected string
	}{
		{
			name:     "family slug and lifecycle state dropped, family id sent",
			input:    createdAsDraft,
			expected: `{"name":"Premium","description":"d","type":"PAID","version":"2","isDefault":false,"familyId":"fam_1"}`,
		},
		{
			name:     "nothing set leaves the family alone",
			input:    untouched,
			expected: `{"name":"Premium","description":"d","type":"PAID","version":"2","isDefault":false}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondNoContent())

			if err := client.Licenses.Update(t.Context(), "premium-v2", test.input); err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			request := sent.only(t)
			if request.method != http.MethodPut || request.path != "/api/licenses/premium-v2" {
				t.Errorf("request = %s %s, want PUT /api/licenses/premium-v2", request.method, request.path)
			}
			if got := string(request.body); got != test.expected {
				t.Errorf("request body = %s, want %s", got, test.expected)
			}
		})
	}
}

// The license writes have to pass their Problem and Code through, not just a
// status: a bare 409 cannot tell "already published", which a release script
// tolerates, from "you may not archive the default", which it must not.
func TestLicenseWritesReportTheirCodes(t *testing.T) {
	t.Parallel()

	familySlug := "premium"
	draft := LicenseLifecycleDraft

	tests := []struct {
		name   string
		status int
		code   string
		call   func(context.Context, *Client) error
	}{
		{
			name:   "create names an unknown family",
			status: http.StatusNotFound,
			code:   "CreateLicense.FamilyNotFound",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Create(ctx, LicenseInput{Name: "Premium", Description: "d", Type: LicenseType("PAID"), FamilySlug: &familySlug})
				return err
			},
		},
		{
			name:   "create makes an unpublished default",
			status: http.StatusConflict,
			code:   "CreateLicense.DefaultMustBePublished",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Create(ctx, LicenseInput{Name: "Premium", Description: "d", Type: LicenseType("PAID"), IsDefault: true, LifecycleState: &draft})
				return err
			},
		},
		{
			name:   "update makes a draft the default",
			status: http.StatusConflict,
			code:   "UpdateLicense.DefaultMustBePublished",
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.Update(ctx, "premium-v2", LicenseInput{Name: "Premium", Description: "d", Type: LicenseType("PAID"), Version: "2", IsDefault: true})
			},
		},
		{
			name:   "publish finds the version already published",
			status: http.StatusConflict,
			code:   "PublishLicense.NotADraft",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Publish(ctx, "premium-v2")
				return err
			},
		},
		{
			name:   "archive targets the family's default",
			status: http.StatusConflict,
			code:   "ArchiveLicense.DefaultMustBePublished",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Archive(ctx, "premium-v2")
				return err
			},
		},
		{
			name:   "unarchive names an unknown version",
			status: http.StatusNotFound,
			code:   "UnarchiveLicense.NotFound",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Unarchive(ctx, "premium-v9")
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, _ := givenClient(t, respondJSON(test.status,
				fmt.Sprintf(`{"title":%q,"status":%d,"code":%q,"detail":"refused"}`, http.StatusText(test.status), test.status, test.code)))

			err := test.call(t.Context(), client)
			if err == nil {
				t.Fatalf("call() error = nil, want a %d", test.status)
			}

			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("errors.As(err, *Error) = false for %T: %v", err, err)
			}
			if apiErr.StatusCode != test.status {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, test.status)
			}
			if apiErr.Code != test.code {
				t.Errorf("Code = %q, want %q", apiErr.Code, test.code)
			}
			if apiErr.Problem == nil {
				t.Error("Problem = nil, want the problem the response carried")
			}
		})
	}
}

// TestLifecycleTransitionsPostToTheirOwnOperation pins the three calls a
// version's state moves through: each is a bodiless POST on the version's own
// sub-resource, and each hands back the version as the server left it, so a
// caller learns the new state without a second read.
func TestLifecycleTransitionsPostToTheirOwnOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		call  func(context.Context, *Client, string) (License, error)
		path  string
		state LicenseLifecycleState
	}{
		{
			name: "publish",
			call: func(ctx context.Context, c *Client, slug string) (License, error) {
				return c.Licenses.Publish(ctx, slug)
			},
			path:  "/api/licenses/premium-v3/publish",
			state: LicenseLifecyclePublished,
		},
		{
			name: "archive",
			call: func(ctx context.Context, c *Client, slug string) (License, error) {
				return c.Licenses.Archive(ctx, slug)
			},
			path:  "/api/licenses/premium-v3/archive",
			state: LicenseLifecycleArchived,
		},
		{
			name: "unarchive",
			call: func(ctx context.Context, c *Client, slug string) (License, error) {
				return c.Licenses.Unarchive(ctx, slug)
			},
			path:  "/api/licenses/premium-v3/unarchive",
			state: LicenseLifecyclePublished,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondJSON(http.StatusOK, fmt.Sprintf(
				`{"slug":"premium-v3","name":"Premium","description":"d","type":"PAID","version":"3","isDefault":false,"lifecycleState":%q}`,
				test.state)))

			license, err := test.call(t.Context(), client, "premium-v3")
			if err != nil {
				t.Fatalf("%s() error = %v", test.name, err)
			}

			request := sent.only(t)
			if request.method != http.MethodPost || request.path != test.path {
				t.Errorf("request = %s %s, want POST %s", request.method, request.path, test.path)
			}
			if len(request.body) != 0 {
				t.Errorf("request body = %s, want none: the path says everything", request.body)
			}
			if license.LifecycleState == nil || *license.LifecycleState != test.state {
				t.Errorf("lifecycleState = %v, want %s", license.LifecycleState, test.state)
			}
		})
	}
}
