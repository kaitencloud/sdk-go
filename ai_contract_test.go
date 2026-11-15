// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	sdk "github.com/kaitencloud/sdk-go"
)

type aiContractTransport func(*http.Request) (*http.Response, error)

func (f aiContractTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func checkAIContractResponse(response interface{ StatusCode() int }, err error) error {
	if err != nil {
		return err
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return fmt.Errorf("status = %d", response.StatusCode())
	}
	return nil
}

// This is deliberately an external package: internal/gen cannot be imported by an
// SDK consumer. All requests are answered by an in-memory transport, never the API.
func TestAIContractOperationsAreUsableOutsideTheModule(t *testing.T) {
	t.Parallel()

	status := sdk.ListModelsParamsStatus("ACTIVE")
	modality := sdk.ListModelsParamsModality("TEXT")
	cursor, transaction, entitlement, rate := "next", "tx-1", "credits", "0.0000000001"
	limit := int32(25)
	model := sdk.Model{Name: "Model", Slug: "model", Vendor: "vendor", Modality: "TEXT", Status: "ACTIVE"}
	tests := []struct {
		name, method, path, query string
		call                      func(*sdk.Client) error
	}{
		{"list models", "GET", "/api/models", "cursor=next&limit=25&modality=TEXT&status=ACTIVE", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().ListModelsWithResponse(t.Context(), &sdk.ListModelsParams{Cursor: &cursor, Limit: &limit, Status: &status, Modality: &modality}))
		}},
		{"get model", "GET", "/api/models/model", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().GetModelWithResponse(t.Context(), model.Slug))
		}},
		{"create model", "POST", "/api/models", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().CreateModelWithResponse(t.Context(), model))
		}},
		{"update model", "PUT", "/api/models/model", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().UpdateModelWithResponse(t.Context(), model.Slug, model))
		}},
		{"delete model", "DELETE", "/api/models/model", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().DeleteModelWithResponse(t.Context(), model.Slug))
		}},
		{"list rules", "GET", "/api/licenses/pro-v3/conversion-rules", "entitlementSlug=credits", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().ListConversionRulesWithResponse(t.Context(), "pro-v3", &sdk.ListConversionRulesParams{EntitlementSlug: &entitlement}))
		}},
		{"create rule", "POST", "/api/licenses/pro-v3/conversion-rules", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().CreateConversionRuleWithResponse(t.Context(), "pro-v3", sdk.CreateConversionRule{EntitlementSlug: entitlement, Rates: sdk.ConversionRates{PerInputToken: &rate}}))
		}},
		{"get rule", "GET", "/api/licenses/pro-v3/conversion-rules/rule", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().GetConversionRuleWithResponse(t.Context(), "pro-v3", "rule"))
		}},
		{"update rule", "PUT", "/api/licenses/pro-v3/conversion-rules/rule", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().UpdateConversionRuleWithResponse(t.Context(), "pro-v3", "rule", sdk.UpdateConversionRule{Rates: sdk.ConversionRates{PerInputToken: &rate}}))
		}},
		{"archive rule", "DELETE", "/api/licenses/pro-v3/conversion-rules/rule", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().ArchiveConversionRuleWithResponse(t.Context(), "pro-v3", "rule"))
		}},
		{"inferences", "GET", "/api/instances/instance/entitlements/credits/usage/inferences", "cursor=next&limit=25&transactionId=tx-1", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().ListAiInferencesWithResponse(t.Context(), "instance", entitlement, &sdk.ListAIInferencesParams{Cursor: &cursor, Limit: &limit, TransactionId: &transaction}))
		}},
		{"reservations", "GET", "/api/instances/instance/entitlements/credits/usage/reservations", "", func(c *sdk.Client) error {
			return checkAIContractResponse(c.Raw().ListAiReservationsWithResponse(t.Context(), "instance", entitlement))
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := sdk.NewClient("https://sdk.invalid/api", sdk.WithHTTPClient(&http.Client{Transport: aiContractTransport(func(r *http.Request) (*http.Response, error) {
				wantQuery, _ := url.ParseQuery(test.query)
				if r.Method != test.method || r.URL.Path != test.path || !reflect.DeepEqual(r.URL.Query(), wantQuery) {
					t.Errorf("request = %s %s, want %s %s?%s", r.Method, r.URL, test.method, test.path, test.query)
				}
				status, body := http.StatusOK, `{}`
				if test.name == "list rules" {
					body = `[]`
				}
				if r.Method == http.MethodPost {
					status = http.StatusCreated
				}
				if r.Method == http.MethodDelete {
					status, body = http.StatusNoContent, ""
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}))
			if err != nil {
				t.Fatal(err)
			}
			if err := test.call(client); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAIContractPreservesDecimalAndUnknownDimensions(t *testing.T) {
	t.Parallel()

	zero, rate, units := "0", "0.0000000001", "123456789012.123456"
	zeroTokens := int64(0)
	requestType, requestStatus := sdk.AIRequestTypeCompletion, sdk.AIRequestStatusSuccess
	transaction := "tx-zero"
	var value sdk.ReportEntitlementUsageValue
	if err := value.FromAiUsageValue(sdk.AIUsageValue{
		Type: "ai_usage", Phase: "commit", TransactionId: &transaction,
		RequestType: &requestType, RequestStatus: &requestStatus,
		Usage: &sdk.AIUsageQuantities{InputTokens: &zeroTokens, Units: &units},
	}); err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(sdk.ReportEntitlementUsageBody{Value: value})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"inputTokens":0`) || !strings.Contains(string(body), `"units":"123456789012.123456"`) || strings.Contains(string(body), "cacheReadInputTokens") {
		t.Fatalf("explicit zero, decimal units or unknown cache changed: %s", body)
	}
	encoded, err := json.Marshal(sdk.ConversionRates{PerInputToken: &rate, PerCachedInputToken: &zero})
	if err != nil || string(encoded) != `{"perCachedInputToken":"0","perInputToken":"0.0000000001"}` {
		t.Fatalf("rates lost decimal/absent semantics: %s (%v)", encoded, err)
	}
	var hook sdk.OnCreateConversionRuleJSONBody
	if err := json.Unmarshal([]byte(`{"data":{"licenseSlug":"pro-v3","modelSlug":null,"requestType":null,"rates":{"perInputToken":"0.0000000001","perCachedInputToken":"0"}}}`), &hook); err != nil {
		t.Fatal(err)
	}
	if hook.Data.ModelSlug != nil || hook.Data.RequestType != nil || hook.Data.Rates.PerInputToken == nil || *hook.Data.Rates.PerInputToken != rate || hook.Data.Rates.PerCacheWriteToken != nil {
		t.Fatalf("hook lost wildcard or rate precision: %+v", hook.Data)
	}
	var modelHook sdk.OnModelCreatedJSONBody
	if err := json.Unmarshal([]byte(`{"data":{"vendorCosts":{"perInputToken":"0.000000000001"}}}`), &modelHook); err != nil {
		t.Fatal(err)
	}
	if modelHook.Data.VendorCosts == nil || modelHook.Data.VendorCosts.PerInputToken == nil || *modelHook.Data.VendorCosts.PerInputToken != "0.000000000001" {
		t.Fatalf("vendor decimal precision lost: %+v", modelHook.Data.VendorCosts)
	}
}
