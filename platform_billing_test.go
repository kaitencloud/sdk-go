// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformClosePeriodsSendsTheScopeToTheOrganization(t *testing.T) {
	t.Parallel()

	client, sent := givenPlatformClient(t, respondJSON(http.StatusOK, `{"closed":2,"examined":3,"held":1,"skipped":0,"hasMore":false,"invoices":[]}`))

	slug := "acme-prod"
	report, err := client.Billing.ClosePeriods(t.Context(), testOrganizationID, PlatformClosePeriodsScope{InstanceSlug: &slug})
	if err != nil {
		t.Fatalf("ClosePeriods() error = %v", err)
	}

	request := sent.only(t)
	if request.method != http.MethodPost || !strings.HasSuffix(request.path, "/platform/organizations/"+testOrganizationID+"/billing/close-periods") {
		t.Errorf("request = %s %s", request.method, request.path)
	}
	if got := string(request.body); got != `{"instanceSlug":"acme-prod"}` {
		t.Errorf("body = %s", got)
	}
	if report.Closed != 2 || report.Held != 1 {
		t.Errorf("report = %+v", report)
	}
}

func TestPlatformClosePeriodsRefusesAMalformedOrganizationID(t *testing.T) {
	t.Parallel()

	client, sent := givenPlatformClient(t, respondNoContent())

	if _, err := client.Billing.ClosePeriods(t.Context(), "not-a-uuid", PlatformClosePeriodsScope{}); err == nil {
		t.Fatal("ClosePeriods() error = nil, want a parse error")
	}
	if sent.count() != 0 {
		t.Errorf("requests = %d, want none", sent.count())
	}
}
