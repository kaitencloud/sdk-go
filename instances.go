package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// List returns all instances, optionally including deleted ones per options.
func (s *Instances) List(ctx context.Context, options *InstancesListOptions) ([]Instance, error) {
	query := url.Values{}
	if options != nil && options.IncludeDeleted {
		query.Set("include_deleted", "true")
	}

	return listAll[Instance](ctx, s.client.list, listPath("instances"), query)
}

// Get returns the instance identified by instanceSlug.
func (s *Instances) Get(ctx context.Context, instanceSlug string) (Instance, error) {
	resp, err := s.client.raw.GetInstanceWithResponse(ctx, instanceSlug)
	if err != nil {
		var zero Instance
		return zero, fmt.Errorf("get instance: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new instance.
func (s *Instances) Create(ctx context.Context, input InstanceInput) (Instance, error) {
	body, err := jsonBody(input.createPayload())
	if err != nil {
		var zero Instance
		return zero, fmt.Errorf("encode instance request: %w", err)
	}

	resp, err := s.client.raw.CreateInstanceWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Instance
		return zero, fmt.Errorf("create instance: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the instance identified by instanceSlug.
func (s *Instances) Update(ctx context.Context, instanceSlug string, input InstanceInput) error {
	body, err := jsonBody(input.updatePayload())
	if err != nil {
		return fmt.Errorf("encode instance request: %w", err)
	}

	resp, err := s.client.raw.UpdateInstanceWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update instance: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// UpdateStatus sets the health status of the instance identified by instanceSlug.
func (s *Instances) UpdateStatus(ctx context.Context, instanceSlug string, status InstanceStatus) error {
	body, err := jsonBody(instancePatchPayload{Status: &status})
	if err != nil {
		return fmt.Errorf("encode instance patch request: %w", err)
	}

	resp, err := s.client.raw.PatchInstanceWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update instance status: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// UpdateLifecycleStage sets the lifecycle stage of the instance identified by instanceSlug.
func (s *Instances) UpdateLifecycleStage(ctx context.Context, instanceSlug string, lifecycleStage string) error {
	body, err := jsonBody(instancePatchPayload{LifecycleStage: &lifecycleStage})
	if err != nil {
		return fmt.Errorf("encode instance patch request: %w", err)
	}

	resp, err := s.client.raw.PatchInstanceWithBodyWithResponse(ctx, instanceSlug, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("update instance lifecycle stage: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// Delete deletes the instance identified by instanceSlug.
func (s *Instances) Delete(ctx context.Context, instanceSlug string) error {
	resp, err := s.client.raw.DeleteInstanceWithResponse(ctx, instanceSlug)
	if err != nil {
		return fmt.Errorf("delete instance: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}

// ListAuditTrails returns the audit trail entries for the instance identified by instanceSlug.
func (s *Instances) ListAuditTrails(ctx context.Context, instanceSlug string, options *AuditTrailsOptions) ([]AuditTrail, error) {
	query := url.Values{}
	var limit int32
	if options != nil {
		if options.EventName != "" {
			query.Set("event_name", options.EventName)
		}
		if after := timeStringPtr(options.After); after != nil {
			query.Set("after", *after)
		}
		if before := timeStringPtr(options.Before); before != nil {
			query.Set("before", *before)
		}
		if options.Limit != nil {
			limit = *options.Limit
		}
	}

	return listAtMost[AuditTrail](ctx, s.client.list, listPath("instances", instanceSlug, "audit-trails"), query, limit)
}

// ListEntitlementUsageMetrics returns usage metrics for all entitlements of the instance identified by instanceSlug.
func (s *Instances) ListEntitlementUsageMetrics(ctx context.Context, instanceSlug string) ([]EntitlementUsage, error) {
	resp, err := s.client.raw.GetEntitlementsUsageMetricsWithResponse(ctx, instanceSlug)
	if err != nil {
		return nil, fmt.Errorf("list entitlement usage metrics: %w", err)
	}

	return expectSlice(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// GetEntitlementUsageMetric returns the usage metric for the entitlement identified by entitlementSlug on the instance identified by instanceSlug.
func (s *Instances) GetEntitlementUsageMetric(ctx context.Context, instanceSlug, entitlementSlug string) (EntitlementUsage, error) {
	resp, err := s.client.raw.GetEntitlementUsageMetricsWithResponse(ctx, instanceSlug, entitlementSlug)
	if err != nil {
		var zero EntitlementUsage
		return zero, fmt.Errorf("get entitlement usage metric: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// ReportEntitlementUsageMetric reports usage of the entitlement identified by entitlementSlug for the instance identified by instanceSlug.
func (s *Instances) ReportEntitlementUsageMetric(ctx context.Context, instanceSlug, entitlementSlug string, input UsageReportInput) (EntitlementUsage, error) {
	usageValue, err := NumberUsageValue(input.Value)
	if err != nil {
		var zero EntitlementUsage
		return zero, fmt.Errorf("encode usage value: %w", err)
	}

	behavior := input.Behavior
	if behavior == "" {
		behavior = Append
	}

	body, err := jsonBody(struct {
		Behavior *EntitlementUsageBehavior `json:"behavior,omitempty"`
		Metadata *map[string]any           `json:"metadata,omitempty"`
		Value    EntitlementUsageValue     `json:"value"`
	}{
		Behavior: &behavior,
		Metadata: cloneMapPtr(input.Metadata),
		Value:    usageValue,
	})
	if err != nil {
		var zero EntitlementUsage
		return zero, fmt.Errorf("encode usage report request: %w", err)
	}

	resp, err := s.client.raw.ReportEntitlementUsageMetricWithBodyWithResponse(ctx, instanceSlug, entitlementSlug, contentTypeJSON, body)
	if err != nil {
		var zero EntitlementUsage
		return zero, fmt.Errorf("report entitlement usage metric: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// ReportUsage reports an incremental usage delta for the entitlement identified by entitlementSlug
// on the instance identified by instanceSlug, returning ErrThresholdExceeded if the limit is reached.
func (s *Instances) ReportUsage(ctx context.Context, instanceSlug, entitlementSlug string, value float64) error {
	_, err := s.ReportEntitlementUsageMetric(ctx, instanceSlug, entitlementSlug, UsageReportInput{
		Value:    value,
		Behavior: Append,
	})
	if err == nil {
		return nil
	}

	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict {
		// Wrap both so callers keep the sentinel and the problem details: errors.Is finds
		// ErrThresholdExceeded and errors.As still reaches the *Error, as everywhere else.
		return fmt.Errorf("report usage failed: %w: %w", ErrThresholdExceeded, apiErr)
	}

	return err
}
