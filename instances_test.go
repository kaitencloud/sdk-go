// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

// TestUpdateRendersTheSlugOnlyWhenSet pins both halves of the optional rename.
//
// PUT /instances/{instanceSlug} takes slug as an ordinary optional field, so the omitted
// half carries as much weight as the sent one: absence is the only way the body can say
// "keep the current slug", and a key rendered from an unset field would ask every update
// that never mentioned a slug to rename its instance to the empty string. Whole bodies are
// compared rather than key sets, because "unchanged for callers who set no slug" is a claim
// about the bytes on the wire.
func TestUpdateRendersTheSlugOnlyWhenSet(t *testing.T) {
	t.Parallel()

	start := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2028, time.January, 1, 0, 0, 0, 0, time.UTC)
	renamed := "renamed"

	unchanged := InstanceInput{
		Name:             "pinned",
		Description:      "d",
		CustomerID:       "cus_1",
		LicenseID:        "lic_1",
		StartLicenseDate: start,
		EndLicenseDate:   end,
	}

	renaming := unchanged
	renaming.Slug = &renamed

	tests := []struct {
		name     string
		input    InstanceInput
		expected string
	}{
		{
			name:     "slug set renames the instance",
			input:    renaming,
			expected: `{"name":"pinned","description":"d","customerId":"cus_1","licenseId":"lic_1","metadata":null,"startLicenseDate":"2027-01-01T00:00:00Z","endLicenseDate":"2028-01-01T00:00:00Z","slug":"renamed"}`,
		},
		{
			name:     "slug unset keeps the current one",
			input:    unchanged,
			expected: `{"name":"pinned","description":"d","customerId":"cus_1","licenseId":"lic_1","metadata":null,"startLicenseDate":"2027-01-01T00:00:00Z","endLicenseDate":"2028-01-01T00:00:00Z"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenClient(t, respondNoContent())

			if err := client.Instances.Update(t.Context(), "instance-slug", test.input); err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			if got := string(sent.only(t).body); got != test.expected {
				t.Errorf("request body = %s, want %s", got, test.expected)
			}
		})
	}
}

func TestReportUsageOnConflictPreservesSentinelAndAPIError(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusConflict, `{"title":"Conflict","detail":"threshold reached"}`))

	err := client.Instances.ReportUsage(t.Context(), "instance-slug", "entitlement-slug", 1)
	if err == nil {
		t.Fatal("ReportUsage() error = nil, want an error")
	}

	// Both, not either: a caller decides what to do from the sentinel and reports what
	// happened from the problem details, and the wrapping has to keep serving both.
	if !errors.Is(err, ErrThresholdExceeded) {
		t.Errorf("errors.Is(err, ErrThresholdExceeded) = false for %v", err)
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(err, *Error) = false for %T: %v", err, err)
	}

	if apiErr.StatusCode != http.StatusConflict {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusConflict)
	}

	if apiErr.Problem == nil || stringValue(apiErr.Problem.Detail) != "threshold reached" {
		t.Errorf("problem details lost: %+v", apiErr.Problem)
	}
}

// Only a 409 means "the limit stopped this". Tagging any other failure with
// ErrThresholdExceeded would have a caller report a quota problem for what is really an
// expired credential or an outage.
func TestReportUsageLeavesNonConflictFailuresUntagged(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusInternalServerError, `{"title":"Internal Server Error"}`))

	err := client.Instances.ReportUsage(t.Context(), "instance-slug", "entitlement-slug", 1)
	if err == nil {
		t.Fatal("ReportUsage() error = nil, want an error")
	}

	if errors.Is(err, ErrThresholdExceeded) {
		t.Errorf("a 500 was reported as %v, want it untagged", ErrThresholdExceeded)
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(err, *Error) = false for %T: %v", err, err)
	}

	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}
}
