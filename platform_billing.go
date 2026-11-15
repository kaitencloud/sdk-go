// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ClosePeriods closes the due billing periods of the organization identified by
// organizationID now, rather than at the period-close job's next pass. Only what is due
// closes, at most the configured batch size per call: while the report's HasMore is true,
// call again. input.InstanceSlug narrows it to that instance's subscription.
func (s *PlatformBilling) ClosePeriods(ctx context.Context, organizationID string, input PlatformClosePeriodsScope) (PlatformClosePeriodsReport, error) {
	parsedOrgID, err := uuid.Parse(organizationID)
	if err != nil {
		var zero PlatformClosePeriodsReport
		return zero, fmt.Errorf("parse organization id: %w", err)
	}

	body, err := jsonBody(input)
	if err != nil {
		var zero PlatformClosePeriodsReport
		return zero, fmt.Errorf("encode close billing periods request: %w", err)
	}

	resp, err := s.client.raw.PlatformCloseBillingPeriodsWithBodyWithResponse(ctx, openapi_types.UUID(parsedOrgID), contentTypeJSON, body)
	if err != nil {
		var zero PlatformClosePeriodsReport
		return zero, fmt.Errorf("close billing periods: %w", err)
	}

	apiResp := withRetryAfter(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.HTTPResponse)
	return expectJSON(apiResp, resp.JSON200)
}
