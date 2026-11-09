// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"net/http"
	"testing"
)

func TestDashboardSummarySendsTheRange(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK, `{"asOf":"2026-11-09T12:00:00Z"}`))

	window := DashboardRange7d
	summary, err := client.Dashboard.Summary(t.Context(), &DashboardSummaryParams{Range: &window})
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}

	request := sent.only(t)
	if request.method != http.MethodGet || request.query.Get("range") != "7d" {
		t.Errorf("request = %s %s?%s, want GET with range=7d", request.method, request.path, request.query.Encode())
	}
	if summary.AsOf.IsZero() {
		t.Error("AsOf is zero, want the decoded summary")
	}
}

func TestDashboardSummaryLeavesTheRangeToKaitenWithoutParams(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK, `{}`))

	if _, err := client.Dashboard.Summary(t.Context(), nil); err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if sent.only(t).query.Has("range") {
		t.Errorf("query = %s, want no range", sent.only(t).query.Encode())
	}
}
