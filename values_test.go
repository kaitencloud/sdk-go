// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"encoding/json"
	"testing"
)

// The six constructors below are the only way to build an entitlement value from outside
// this package: both value types are generated unions, so a caller cannot write a struct
// literal. What each one has to get right is the discriminator, because a wrong or
// missing "type" is not a validation error a caller sees -- the server picks a different
// arm and stores something else.
//
// So the assertion is on the rendered bytes rather than on a round trip through the
// generated As* accessors: what the server reads is what decides which arm it takes.
func TestEntitlementValueConstructors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         json.Marshaler
		expectedType  string
		expectedValue string
	}{
		{
			name:          "number usage value carries the number tag",
			value:         must(NumberUsageValue(12.5)),
			expectedType:  "number",
			expectedValue: `12.5`,
		},
		{
			name:          "boolean usage value carries the boolean tag",
			value:         must(BooleanUsageValue(true)),
			expectedType:  "boolean",
			expectedValue: `true`,
		},
		{
			name:          "config usage value carries the object tag",
			value:         must(ConfigUsageValue(map[string]any{"region": "eu-west-1"})),
			expectedType:  "object",
			expectedValue: `{"region":"eu-west-1"}`,
		},
		{
			name:          "number license value carries the number tag",
			value:         must(NumberLicenseValue(50)),
			expectedType:  "number",
			expectedValue: `50`,
		},
		{
			name:          "boolean license value carries the boolean tag",
			value:         must(BooleanLicenseValue(true)),
			expectedType:  "boolean",
			expectedValue: `true`,
		},
		{
			name:          "config license value carries the object tag",
			value:         must(ConfigLicenseValue(map[string]any{"tier": "gold"})),
			expectedType:  "object",
			expectedValue: `{"tier":"gold"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := json.Marshal(test.value)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			var decoded struct {
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			}
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("json.Unmarshal(%s) error = %v", encoded, err)
			}

			if decoded.Type != test.expectedType {
				t.Errorf("type = %q, want %q", decoded.Type, test.expectedType)
			}
			if string(decoded.Value) != test.expectedValue {
				t.Errorf("value = %s, want %s", decoded.Value, test.expectedValue)
			}
		})
	}
}

// The config constructors clone the caller's map. Without that, a caller reusing one
// builder map across several entitlements would find every value it built pointing at
// whatever the last mutation left in it.
func TestConfigValueConstructorsCloneTheCallersMap(t *testing.T) {
	t.Parallel()

	settings := map[string]any{"region": "eu-west-1"}

	usage := must(ConfigUsageValue(settings))
	license := must(ConfigLicenseValue(settings))

	settings["region"] = "us-east-1"

	usageValue, err := usage.AsConfigEntitlementValue()
	if err != nil {
		t.Fatalf("AsConfigEntitlementValue() error = %v", err)
	}
	if usageValue.Value["region"] != "eu-west-1" {
		t.Errorf("usage region = %v, want eu-west-1", usageValue.Value["region"])
	}

	licenseValue, err := license.AsConfigEntitlementValue()
	if err != nil {
		t.Fatalf("AsConfigEntitlementValue() error = %v", err)
	}
	if licenseValue.Value["region"] != "eu-west-1" {
		t.Errorf("license region = %v, want eu-west-1", licenseValue.Value["region"])
	}
}

// must unwraps a value constructor's (value, error) pair so a table entry can call one
// inline. It panics rather than taking a *testing.T, because Go only expands a
// multi-value call when it is the whole argument list -- and these constructors fail
// only when the generated union cannot marshal the arm at all, which is a broken build
// rather than a case worth a table row.
func must[T any](value T, err error) T {
	if err != nil {
		panic("building the entitlement value: " + err.Error())
	}

	return value
}
