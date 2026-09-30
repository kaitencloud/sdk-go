package sdk

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

func TestNewClientRejectsABlankBaseURL(t *testing.T) {
	t.Parallel()

	if _, err := NewClient("   "); err == nil {
		t.Fatal("NewClient() with a blank base URL returned no error")
	}
}

func TestNewClientInitializesEveryResourceNamespace(t *testing.T) {
	t.Parallel()

	client, err := NewClient(coreBaseURL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// A Client's zero value leaves these nil and calling through one panics, so
	// "the constructor populated all of them" is the type's whole usability contract.
	namespaces := map[string]any{
		"Components":        client.Components,
		"Customers":         client.Customers,
		"DeploymentZones":   client.DeploymentZones,
		"EntitlementGroups": client.EntitlementGroups,
		"Entitlements":      client.Entitlements,
		"FeatureFlags":      client.FeatureFlags,
		"Instances":         client.Instances,
		"LicenseFamilies":   client.LicenseFamilies,
		"Licenses":          client.Licenses,
		"MetadataFields":    client.MetadataFields,
		"Releases":          client.Releases,
		"ServiceAccounts":   client.ServiceAccounts,
	}

	for name, namespace := range namespaces {
		if reflect.ValueOf(namespace).IsNil() {
			t.Errorf("%s namespace was not initialized", name)
		}
	}
}

func TestNewClientDefaultsHTTPTimeout(t *testing.T) {
	t.Parallel()

	client, err := NewClient(coreBaseURL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if timeout := httpClientTimeout(t, client); timeout != DefaultTimeout {
		t.Fatalf("default HTTP client Timeout = %v, want %v", timeout, DefaultTimeout)
	}
}

func TestNewClientLetsWithHTTPClientReplaceTheDefaultTimeout(t *testing.T) {
	t.Parallel()

	want := 5 * time.Second

	client, err := NewClient(coreBaseURL, WithHTTPClient(&http.Client{Timeout: want}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if timeout := httpClientTimeout(t, client); timeout != want {
		t.Fatalf("HTTP client Timeout = %v, want %v", timeout, want)
	}
}

func TestClientExposesOnlyTheReviewedPublicSurface(t *testing.T) {
	t.Parallel()

	client, err := NewClient(coreBaseURL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	assertServiceCoverage(t, map[string]serviceCoverage{
		"Components": {
			service: client.Components,
			methods: []string{"List", "Get", "Create", "Update", "Delete"},
		},
		"Customers": {
			service: client.Customers,
			methods: []string{"List", "Get", "Create", "Update", "Delete"},
		},
		"DeploymentZones": {
			service: client.DeploymentZones,
			methods: []string{"List", "Get", "Create", "Update", "Delete"},
		},
		"EntitlementGroups": {
			service: client.EntitlementGroups,
			methods: []string{"List", "Get", "Create", "Update", "Delete", "AddEntitlement", "RemoveEntitlement", "GetUsage"},
		},
		"Entitlements": {
			service: client.Entitlements,
			methods: []string{"List", "Get", "Create", "Update", "Delete"},
		},
		"FeatureFlags": {
			service: client.FeatureFlags,
			methods: []string{"List", "Get", "Create", "Update", "Delete"},
		},
		"Instances": {
			service: client.Instances,
			methods: []string{"List", "Get", "Create", "Update", "UpdateStatus", "UpdateLifecycleStage", "Delete", "ListAuditTrails", "ListEntitlementUsageMetrics", "GetEntitlementUsageMetric", "ReportEntitlementUsageMetric", "ReportUsage"},
		},
		"LicenseFamilies": {
			service: client.LicenseFamilies,
			methods: []string{"List", "Get"},
		},
		"Licenses": {
			service: client.Licenses,
			methods: []string{"List", "Get", "Create", "Update", "Publish", "Archive", "Unarchive", "Delete", "ListEntitlements", "AssociateEntitlement", "GetEntitlement", "UpdateEntitlement", "DeleteEntitlement"},
		},
		"MetadataFields": {
			service: client.MetadataFields,
			methods: []string{"List", "Create", "Update", "Archive", "Unarchive", "DryRun", "Reorder"},
		},
		"Releases": {
			service: client.Releases,
			methods: []string{"List", "Get", "Create", "Delete"},
		},
		"ServiceAccounts": {
			service: client.ServiceAccounts,
			methods: []string{"List", "Get", "Create", "Update", "ListTokens", "CreateToken", "DeleteToken"},
		},
	})
}

func httpClientTimeout(t *testing.T, client *Client) time.Duration {
	t.Helper()

	raw, ok := client.raw.ClientInterface.(*gen.Client)
	if !ok {
		t.Fatalf("underlying client is %T, want *gen.Client", client.raw.ClientInterface)
	}

	doer, ok := raw.Client.(*http.Client)
	if !ok {
		t.Fatalf("request doer is %T, want *http.Client", raw.Client)
	}

	return doer.Timeout
}

// TestClientWrappersReportEveryDeclaredProblem drives every Core wrapper through each
// problem status its operation declares and expects the *Error to carry the document and
// its code. Twenty-four wrappers were dropping one when this was written -- every create
// its 503, most of them a 404 or a 409 as well, the four feature-flag wrappers a 403 --
// and each dropped status reached callers as an *Error with no Problem and an empty
// Code, the member they are told to branch on. The declared set comes off the generated
// struct, so a status the spec adds arrives here with the next regeneration. What this
// table has to keep up with is the wrappers, and assertProblemCoverage checks it
// against the source.
func TestClientWrappersReportEveryDeclaredProblem(t *testing.T) {
	t.Parallel()

	start := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	instance := InstanceInput{Name: "Prod", Description: "d", CustomerID: "cus_1", LicenseID: "lic_1", StartLicenseDate: start, EndLicenseDate: start.AddDate(1, 0, 0)}
	license := LicenseInput{Name: "Premium", Description: "d", Type: LicenseType("PAID"), Version: "1"}

	seats, err := NumberLicenseValue(5)
	if err != nil {
		t.Fatalf("NumberLicenseValue() error = %v", err)
	}

	given := func(t *testing.T, respond responder) *Client {
		t.Helper()

		client, _ := givenClient(t, respond)
		return client
	}

	assertProblemCoverage(t, given, map[string]problemCoverage[*Client]{
		"Components.Get": {
			response: gen.GetComponentResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Components.Get(ctx, "web"); return err },
		},
		"Components.Create": {
			response: gen.CreateComponentResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Components.Create(ctx, ComponentInput{Name: "Web", Version: "1.0.0"})
				return err
			},
		},
		"Components.Update": {
			response: gen.UpdateComponentResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Components.Update(ctx, "web", ComponentInput{Name: "Web", Version: "1.0.1"})
				return err
			},
		},
		"Components.Delete": {
			response: gen.DeleteComponentResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Components.Delete(ctx, "web") },
		},
		"Customers.Get": {
			response: gen.GetCustomerResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Customers.Get(ctx, "acme"); return err },
		},
		"Customers.Create": {
			response: gen.CreateCustomerResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Customers.Create(ctx, CustomerInput{Name: "Acme"})
				return err
			},
		},
		"Customers.Update": {
			response: gen.UpdateCustomerResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Customers.Update(ctx, "acme", CustomerInput{Name: "Acme"})
			},
		},
		"Customers.Delete": {
			response: gen.DeleteCustomerResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Customers.Delete(ctx, "acme") },
		},
		"DeploymentZones.Get": {
			response: gen.GetDeploymentZoneBySlugResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.DeploymentZones.Get(ctx, "eu-west")
				return err
			},
		},
		"DeploymentZones.Create": {
			response: gen.CreateDeploymentZoneResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.DeploymentZones.Create(ctx, DeploymentZoneInput{Name: "EU West", Type: "CLOUD", Description: "d"})
				return err
			},
		},
		"DeploymentZones.Update": {
			response: gen.UpdateDeploymentZoneResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.DeploymentZones.Update(ctx, "eu-west", DeploymentZoneInput{Name: "EU West", Type: "CLOUD", Description: "d"})
			},
		},
		"DeploymentZones.Delete": {
			response: gen.DeleteDeploymentZoneResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.DeploymentZones.Delete(ctx, "eu-west") },
		},
		"EntitlementGroups.Get": {
			response: gen.GetEntitlementGroupResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.EntitlementGroups.Get(ctx, "platform")
				return err
			},
		},
		"EntitlementGroups.Create": {
			response: gen.CreateEntitlementGroupResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.EntitlementGroups.Create(ctx, EntitlementGroupInput{Name: "Platform"})
				return err
			},
		},
		"EntitlementGroups.Update": {
			response: gen.UpdateEntitlementGroupResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.EntitlementGroups.Update(ctx, "platform", EntitlementGroupInput{Name: "Platform"})
			},
		},
		"EntitlementGroups.Delete": {
			response: gen.DeleteEntitlementGroupResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.EntitlementGroups.Delete(ctx, "platform") },
		},
		"EntitlementGroups.AddEntitlement": {
			response: gen.AddEntitlementToGroupResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.EntitlementGroups.AddEntitlement(ctx, "platform", "seats")
			},
		},
		"EntitlementGroups.RemoveEntitlement": {
			response: gen.RemoveEntitlementFromGroupResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.EntitlementGroups.RemoveEntitlement(ctx, "platform", "seats")
			},
		},
		"EntitlementGroups.GetUsage": {
			response: gen.GetEntitlementGroupUsageResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.EntitlementGroups.GetUsage(ctx, "platform", "acme-prod")
				return err
			},
		},
		"Entitlements.Get": {
			response: gen.GetEntitlementResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Entitlements.Get(ctx, "seats"); return err },
		},
		"Entitlements.Create": {
			response: gen.CreateEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Entitlements.Create(ctx, EntitlementInput{Name: "Seats"})
				return err
			},
		},
		"Entitlements.Update": {
			response: gen.UpdateEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Entitlements.Update(ctx, "seats", EntitlementInput{Name: "Seats"})
			},
		},
		"Entitlements.Delete": {
			response: gen.DeleteEntitlementResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Entitlements.Delete(ctx, "seats") },
		},
		"FeatureFlags.Get": {
			response: gen.GetFeatureFlagResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.FeatureFlags.Get(ctx, "dark-mode"); return err },
		},
		"FeatureFlags.Create": {
			response: gen.CreateFeatureFlagResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.FeatureFlags.Create(ctx, FeatureFlag{})
				return err
			},
		},
		"FeatureFlags.Update": {
			response: gen.UpdateFeatureFlagResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.FeatureFlags.Update(ctx, "dark-mode", FeatureFlag{})
			},
		},
		"FeatureFlags.Delete": {
			response: gen.DeleteFeatureFlagResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.FeatureFlags.Delete(ctx, "dark-mode") },
		},
		"Instances.Get": {
			response: gen.GetInstanceResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Instances.Get(ctx, "acme-prod"); return err },
		},
		"Instances.Create": {
			response: gen.CreateInstanceResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Instances.Create(ctx, instance); return err },
		},
		"Instances.Update": {
			response: gen.UpdateInstanceResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Instances.Update(ctx, "acme-prod", instance) },
		},
		"Instances.UpdateStatus": {
			response: gen.PatchInstanceResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Instances.UpdateStatus(ctx, "acme-prod", InstanceStatusHealthy)
			},
		},
		"Instances.UpdateLifecycleStage": {
			response: gen.PatchInstanceResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Instances.UpdateLifecycleStage(ctx, "acme-prod", "production")
			},
		},
		"Instances.Delete": {
			response: gen.DeleteInstanceResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Instances.Delete(ctx, "acme-prod") },
		},
		"Instances.ListEntitlementUsageMetrics": {
			response: gen.GetEntitlementsUsageMetricsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Instances.ListEntitlementUsageMetrics(ctx, "acme-prod")
				return err
			},
		},
		"Instances.GetEntitlementUsageMetric": {
			response: gen.GetEntitlementUsageMetricsResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Instances.GetEntitlementUsageMetric(ctx, "acme-prod", "seats")
				return err
			},
		},
		"Instances.ReportEntitlementUsageMetric": {
			response: gen.ReportEntitlementUsageMetricResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Instances.ReportEntitlementUsageMetric(ctx, "acme-prod", "seats", UsageReportInput{Value: 1})
				return err
			},
		},
		"LicenseFamilies.Get": {
			response: gen.GetLicenseFamilyResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.LicenseFamilies.Get(ctx, "premium", nil)
				return err
			},
		},
		"Licenses.Get": {
			response: gen.GetLicenseResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Licenses.Get(ctx, "premium-v1"); return err },
		},
		"Licenses.Create": {
			response: gen.CreateLicenseResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Licenses.Create(ctx, license); return err },
		},
		"Licenses.Update": {
			response: gen.UpdateLicenseResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Licenses.Update(ctx, "premium-v1", license) },
		},
		"Licenses.Publish": {
			response: gen.PublishLicenseResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Publish(ctx, "premium-v1")
				return err
			},
		},
		"Licenses.Archive": {
			response: gen.ArchiveLicenseResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Archive(ctx, "premium-v1")
				return err
			},
		},
		"Licenses.Unarchive": {
			response: gen.UnarchiveLicenseResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.Unarchive(ctx, "premium-v1")
				return err
			},
		},
		"Licenses.Delete": {
			response: gen.DeleteLicenseResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Licenses.Delete(ctx, "premium-v1") },
		},
		"Licenses.AssociateEntitlement": {
			response: gen.AssociateEntitlementWithLicenseResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.AssociateEntitlement(ctx, "premium-v1", "seats", seats)
			},
		},
		"Licenses.GetEntitlement": {
			response: gen.GetLicenseEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Licenses.GetEntitlement(ctx, "premium-v1", "seats")
				return err
			},
		},
		"Licenses.UpdateEntitlement": {
			response: gen.UpdateLicenseEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.UpdateEntitlement(ctx, "premium-v1", "seats", seats)
			},
		},
		"Licenses.DeleteEntitlement": {
			response: gen.DeleteLicenseEntitlementResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.Licenses.DeleteEntitlement(ctx, "premium-v1", "seats")
			},
		},
		"MetadataFields.Create": {
			response: gen.CreateMetadataFieldResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.MetadataFields.Create(ctx, MetadataFieldInput{ResourceType: gen.MetadataFieldResourceTypeINSTANCE, Key: "region", Label: "Region"})
				return err
			},
		},
		"MetadataFields.Update": {
			response: gen.UpdateMetadataFieldResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.MetadataFields.Update(ctx, "mf_1", MetadataFieldUpdateInput{Label: "Region"})
				return err
			},
		},
		"MetadataFields.Archive": {
			response: gen.ArchiveMetadataFieldResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.MetadataFields.Archive(ctx, "mf_1")
				return err
			},
		},
		"MetadataFields.Unarchive": {
			response: gen.UnarchiveMetadataFieldResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.MetadataFields.Unarchive(ctx, "mf_1")
				return err
			},
		},
		"MetadataFields.DryRun": {
			response: gen.DryRunMetadataFieldResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.MetadataFields.DryRun(ctx, "mf_1", map[string]any{"type": "string"})
				return err
			},
		},
		"MetadataFields.Reorder": {
			response: gen.ReorderMetadataFieldsResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.MetadataFields.Reorder(ctx, []string{"mf_1"}) },
		},
		"Releases.Get": {
			response: gen.GetReleaseBySlugResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.Releases.Get(ctx, "v1"); return err },
		},
		"Releases.Create": {
			response: gen.CreateReleaseResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.Releases.Create(ctx, ReleaseInput{Version: "1.0.0"})
				return err
			},
		},
		"Releases.Delete": {
			response: gen.DeleteReleaseResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.Releases.Delete(ctx, "v1") },
		},
		"ServiceAccounts.Get": {
			response: gen.GetServiceAccountResponse{},
			call:     func(ctx context.Context, c *Client) error { _, err := c.ServiceAccounts.Get(ctx, "ci"); return err },
		},
		"ServiceAccounts.Create": {
			response: gen.CreateServiceAccountResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.Create(ctx, ServiceAccountInput{Name: "CI"})
				return err
			},
		},
		"ServiceAccounts.Update": {
			response: gen.UpdateServiceAccountResponse{},
			call: func(ctx context.Context, c *Client) error {
				return c.ServiceAccounts.Update(ctx, "ci", ServiceAccountInput{Name: "CI"})
			},
		},
		"ServiceAccounts.CreateToken": {
			response: gen.CreateServiceAccountTokenResponse{},
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ServiceAccounts.CreateToken(ctx, "ci", TokenInput{Name: "deploy"})
				return err
			},
		},
		"ServiceAccounts.DeleteToken": {
			response: gen.DeleteServiceAccountTokenResponse{},
			call:     func(ctx context.Context, c *Client) error { return c.ServiceAccounts.DeleteToken(ctx, "ci", "deploy") },
		},
	})
}
