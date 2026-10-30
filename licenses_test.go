// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"encoding/json"
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

// TestLicenseWritesRenderTheCommercialTermsOnlyWhenSet pins the four commercial terms
// on both license bodies. A term left nil renders no key: the server's default on
// create, the stored value on update, and the bytes a caller who sets no term sent
// before the terms existed. A term that is set reaches the wire whatever its value:
// false turns the payment method off, and 0 and "" are how an update removes the trial
// and the URL. Create refuses that 0 and reads that "" as no URL, which is the
// server's call to make, not the SDK's.
func TestLicenseWritesRenderTheCommercialTermsOnlyWhenSet(t *testing.T) {
	t.Parallel()

	paid, custom := LicensePricingTypePAID, LicensePricingTypeCUSTOM
	twoWeeks, noTrial := int32(14), int32(0)
	required, notRequired := true, false
	contactSales, noURL := "https://kaiten-sushi.example/contact-sales", ""

	tests := []struct {
		name                  string
		pricingType           *LicensePricingType
		trialPeriodDays       *int32
		requiresPaymentMethod *bool
		selfServeCtaURL       *string
		createBody            string
		updateBody            string
	}{
		{
			name:       "no term set leaves each to the server",
			createBody: `{"name":"Premium","description":"d","type":"PAID","isDefault":false}`,
			updateBody: `{"name":"Premium","description":"d","type":"PAID","version":"2","isDefault":false}`,
		},
		{
			name:            "CUSTOM, sold through a contact page",
			pricingType:     &custom,
			selfServeCtaURL: &contactSales,
			createBody:      `{"name":"Premium","description":"d","type":"PAID","isDefault":false,"pricingType":"CUSTOM","selfServeCtaUrl":"https://kaiten-sushi.example/contact-sales"}`,
			updateBody:      `{"name":"Premium","description":"d","type":"PAID","version":"2","isDefault":false,"pricingType":"CUSTOM","selfServeCtaUrl":"https://kaiten-sushi.example/contact-sales"}`,
		},
		{
			name:                  "PAID, with a trial and a payment method at sign-up",
			pricingType:           &paid,
			trialPeriodDays:       &twoWeeks,
			requiresPaymentMethod: &required,
			createBody:            `{"name":"Premium","description":"d","type":"PAID","isDefault":false,"pricingType":"PAID","trialPeriodDays":14,"requiresPaymentMethod":true}`,
			updateBody:            `{"name":"Premium","description":"d","type":"PAID","version":"2","isDefault":false,"pricingType":"PAID","trialPeriodDays":14,"requiresPaymentMethod":true}`,
		},
		{
			name:                  "false, 0 and an empty URL are sent, not dropped",
			trialPeriodDays:       &noTrial,
			requiresPaymentMethod: &notRequired,
			selfServeCtaURL:       &noURL,
			createBody:            `{"name":"Premium","description":"d","type":"PAID","isDefault":false,"trialPeriodDays":0,"requiresPaymentMethod":false,"selfServeCtaUrl":""}`,
			updateBody:            `{"name":"Premium","description":"d","type":"PAID","version":"2","isDefault":false,"trialPeriodDays":0,"requiresPaymentMethod":false,"selfServeCtaUrl":""}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := LicenseInput{
				Name:                  "Premium",
				Description:           "d",
				Type:                  LicenseType("PAID"),
				Version:               "2",
				PricingType:           test.pricingType,
				TrialPeriodDays:       test.trialPeriodDays,
				RequiresPaymentMethod: test.requiresPaymentMethod,
				SelfServeCtaURL:       test.selfServeCtaURL,
			}

			client, sent := givenClient(t, respondToWrite)

			if _, err := client.Licenses.Create(t.Context(), input); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if err := client.Licenses.Update(t.Context(), "premium-v2", input); err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			for i, want := range []struct{ method, path, body string }{
				{http.MethodPost, "/api/licenses", test.createBody},
				{http.MethodPut, "/api/licenses/premium-v2", test.updateBody},
			} {
				request := sent.at(t, i)
				if request.method != want.method || request.path != want.path {
					t.Errorf("request %d = %s %s, want %s %s", i, request.method, request.path, want.method, want.path)
				}
				if got := string(request.body); got != want.body {
					t.Errorf("%s body = %s, want %s", want.method, got, want.body)
				}
			}
		})
	}
}

// TestLicenseReadsCarryTheCommercialTerms pins the four terms on a License read
// back, and on the version a family resolves to, which is the same type. A term the
// server sends as null comes back nil rather than as a zero: written back on update
// a zero would remove the trial or the URL, and on create 0 is refused.
func TestLicenseReadsCarryTheCommercialTerms(t *testing.T) {
	t.Parallel()

	// terms projects the four fields back to JSON with no omitempty, so nil shows as
	// null and is told apart from a zero.
	type terms struct {
		PricingType           *LicensePricingType `json:"pricingType"`
		TrialPeriodDays       *int32              `json:"trialPeriodDays"`
		RequiresPaymentMethod *bool               `json:"requiresPaymentMethod"`
		SelfServeCtaURL       *string             `json:"selfServeCtaUrl"`
	}

	premium := `"pricingType":"CUSTOM","trialPeriodDays":null,"requiresPaymentMethod":false,"selfServeCtaUrl":"https://kaiten-sushi.example/contact-sales"`
	standard := `"pricingType":"PAID","trialPeriodDays":14,"requiresPaymentMethod":true,"selfServeCtaUrl":null`

	tests := []struct {
		name  string
		body  string
		read  func(context.Context, *Client) (*License, error)
		terms string
	}{
		{
			name:  "a license",
			body:  `{"slug":"premium-v2","name":"Premium","description":"d","type":"PAID","version":"2","isDefault":true,` + premium + `}`,
			terms: premium,
			read: func(ctx context.Context, c *Client) (*License, error) {
				license, err := c.Licenses.Get(ctx, "premium-v2")
				return &license, err
			},
		},
		{
			name:  "a family's current version",
			body:  `{"slug":"standard","versionCount":2,"currentVersion":{"slug":"standard","name":"Standard","description":"d","type":"DEVELOPMENT","version":"1","isDefault":true,` + standard + `}}`,
			terms: standard,
			read: func(ctx context.Context, c *Client) (*License, error) {
				family, err := c.LicenseFamilies.Get(ctx, "standard", nil)
				return family.CurrentVersion, err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, _ := givenClient(t, respondJSON(http.StatusOK, test.body))

			license, err := test.read(t.Context(), client)
			if err != nil {
				t.Fatalf("read error = %v", err)
			}
			if license == nil {
				t.Fatal("license = nil, want the version the body carries")
			}

			got, err := json.Marshal(terms{license.PricingType, license.TrialPeriodDays, license.RequiresPaymentMethod, license.SelfServeCtaUrl})
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if want := "{" + test.terms + "}"; string(got) != want {
				t.Errorf("terms = %s, want %s", got, want)
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

// TestGrantWritesRenderTheOveragePercentOnlyWhenSet pins the allowance on both grant
// bodies. Without WithOveragePercent the bytes are the ones the two methods sent before
// the option existed, so the server keeps deriving the allowance from the value. With it,
// the percentage reaches the wire whatever it is: 0, a hard limit and what the server
// would have derived for a finite value, as much as -1, which an unlimited value requires.
func TestGrantWritesRenderTheOveragePercentOnlyWhenSet(t *testing.T) {
	t.Parallel()

	finite, err := NumberLicenseValue(1000)
	if err != nil {
		t.Fatalf("NumberLicenseValue() error = %v", err)
	}
	unlimited, err := NumberLicenseValue(-1)
	if err != nil {
		t.Fatalf("NumberLicenseValue() error = %v", err)
	}

	tests := []struct {
		name          string
		value         LicenseEntitlementValue
		opts          []LicenseEntitlementOption
		associateBody string
		updateBody    string
	}{
		{
			name:          "no option leaves the allowance to the server",
			value:         finite,
			associateBody: `{"entitlementSlug":"monthly-orders","value":{"type":"number","value":1000}}`,
			updateBody:    `{"value":{"type":"number","value":1000}}`,
		},
		{
			name:          "a soft limit",
			value:         finite,
			opts:          []LicenseEntitlementOption{WithOveragePercent(20)},
			associateBody: `{"entitlementSlug":"monthly-orders","value":{"type":"number","value":1000},"limitCapExceededOveragePercent":20}`,
			updateBody:    `{"value":{"type":"number","value":1000},"limitCapExceededOveragePercent":20}`,
		},
		{
			name:          "a hard limit is sent, not dropped",
			value:         finite,
			opts:          []LicenseEntitlementOption{WithOveragePercent(0)},
			associateBody: `{"entitlementSlug":"monthly-orders","value":{"type":"number","value":1000},"limitCapExceededOveragePercent":0}`,
			updateBody:    `{"value":{"type":"number","value":1000},"limitCapExceededOveragePercent":0}`,
		},
		{
			name:          "unlimited",
			value:         unlimited,
			opts:          []LicenseEntitlementOption{WithOveragePercent(-1)},
			associateBody: `{"entitlementSlug":"monthly-orders","value":{"type":"number","value":-1},"limitCapExceededOveragePercent":-1}`,
			updateBody:    `{"value":{"type":"number","value":-1},"limitCapExceededOveragePercent":-1}`,
		},
		{
			name:          "a nil option is skipped",
			value:         finite,
			opts:          []LicenseEntitlementOption{nil},
			associateBody: `{"entitlementSlug":"monthly-orders","value":{"type":"number","value":1000}}`,
			updateBody:    `{"value":{"type":"number","value":1000}}`,
		},
		{
			name:          "the last option wins",
			value:         finite,
			opts:          []LicenseEntitlementOption{WithOveragePercent(10), WithOveragePercent(20)},
			associateBody: `{"entitlementSlug":"monthly-orders","value":{"type":"number","value":1000},"limitCapExceededOveragePercent":20}`,
			updateBody:    `{"value":{"type":"number","value":1000},"limitCapExceededOveragePercent":20}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondNoContent())

			if err := client.Licenses.AssociateEntitlement(t.Context(), "premium", "monthly-orders", test.value, test.opts...); err != nil {
				t.Fatalf("AssociateEntitlement() error = %v", err)
			}
			if err := client.Licenses.UpdateEntitlement(t.Context(), "premium", "monthly-orders", test.value, test.opts...); err != nil {
				t.Fatalf("UpdateEntitlement() error = %v", err)
			}

			for i, want := range []struct{ method, path, body string }{
				{http.MethodPost, "/api/licenses/premium/entitlements", test.associateBody},
				{http.MethodPut, "/api/licenses/premium/entitlements/monthly-orders", test.updateBody},
			} {
				request := sent.at(t, i)
				if request.method != want.method || request.path != want.path {
					t.Errorf("request %d = %s %s, want %s %s", i, request.method, request.path, want.method, want.path)
				}
				if got := string(request.body); got != want.body {
					t.Errorf("%s body = %s, want %s", want.method, got, want.body)
				}
			}
		})
	}
}
