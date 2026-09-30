package sdk

import (
	"slices"
	"testing"
)

// TestRegisterConnectorSendsTheManifestTheSpecDeclares pins the key set, for the reason
// TestRequestBodiesMatchTheirSpecBodies pins the Core ones: RegisterConnectorBody is
// additionalProperties:false with three required keys, so a key too many and a required
// key missing are the same 422 with no type error to catch it. This body is also the one
// place in the SDK that renders snake_case -- the Platform connector document names the
// fields that way, and a camelCase settingsSchema would be refused.
func TestRegisterConnectorSendsTheManifestTheSpecDeclares(t *testing.T) {
	t.Parallel()

	slug := "connector-attio"

	tests := []struct {
		name     string
		input    RegisterConnectorInput
		expected []string
	}{
		{
			name:     "ungated connector",
			input:    RegisterConnectorInput{Name: "kaiten.integration.crm.attio", Version: "1.0.0", SettingsSchema: map[string]any{"type": "object"}},
			expected: []string{"name", "settings_schema", "version"},
		},
		{
			name:     "licence-gated connector",
			input:    RegisterConnectorInput{Name: "kaiten.integration.crm.attio", Version: "1.0.0", SettingsSchema: map[string]any{"type": "object"}, EntitlementSlug: &slug},
			expected: []string{"entitlement_slug", "name", "settings_schema", "version"},
		},
		{
			// settings_schema is required, so a connector with nothing to configure
			// renders an empty object rather than dropping the key or sending null.
			name:     "connector with no settings",
			input:    RegisterConnectorInput{Name: "kaiten.integration.crm.attio", Version: "1.0.0"},
			expected: []string{"name", "settings_schema", "version"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client, sent := givenPlatformClient(t, respondNoContent())

			_, _ = client.Connectors.Register(t.Context(), test.input)

			got := sent.only(t)

			if keys := got.bodyKeys(t); !slices.Equal(keys, test.expected) {
				t.Errorf("request body keys = %v, want %v", keys, test.expected)
			}

			schema, ok := got.decodeBody(t)["settings_schema"].(map[string]any)
			if !ok {
				t.Fatalf("settings_schema = %v, want an object", got.decodeBody(t)["settings_schema"])
			}
			if test.input.SettingsSchema == nil && len(schema) != 0 {
				t.Errorf("settings_schema = %v, want an empty object", schema)
			}
		})
	}
}
