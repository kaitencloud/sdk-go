// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"
)

// Summary returns the console's home in one response, computed by Kaiten over the whole
// organization. params selects the period the rates and the lists cover; nil leaves the
// choice to Kaiten, whose default is 30 days.
func (s *Dashboard) Summary(ctx context.Context, params *DashboardSummaryParams) (DashboardSummary, error) {
	resp, err := s.client.raw.GetDashboardSummaryWithResponse(ctx, params)
	if err != nil {
		var zero DashboardSummary
		return zero, fmt.Errorf("get dashboard summary: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}
