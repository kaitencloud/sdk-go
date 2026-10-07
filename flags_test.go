// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"encoding/json"
	"testing"
)

// TestBasicDefaultVariantMarshalsTheBasicArm pins the union tag. DefaultVariant is a
// union over rollout strategies, and a missing or wrong "type" is not a validation error
// a caller sees -- the server picks a different arm and the flag serves the wrong
// fallback, which for a gate means "on" where the author wrote "off".
func TestBasicDefaultVariantMarshalsTheBasicArm(t *testing.T) {
	t.Parallel()

	defaultVariant, err := BasicDefaultVariant("off")
	if err != nil {
		t.Fatalf("BasicDefaultVariant() error = %v", err)
	}

	encoded, err := json.Marshal(defaultVariant)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded["type"] != "basic" {
		t.Errorf(`type = %v, want "basic"`, decoded["type"])
	}
	if decoded["value"] != "off" {
		t.Errorf(`value = %v, want "off"`, decoded["value"])
	}
}

func TestBasicVariantTargetingMarshalsEveryField(t *testing.T) {
	t.Parallel()

	targeting, err := BasicVariantTargeting("Demo instances", "__kaiten.instance.metadata.demo == true", "on")
	if err != nil {
		t.Fatalf("BasicVariantTargeting() error = %v", err)
	}

	encoded, err := json.Marshal(targeting)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	expected := map[string]any{
		"name":    "Demo instances",
		"rule":    "__kaiten.instance.metadata.demo == true",
		"type":    "basic",
		"variant": "on",
	}
	for key, want := range expected {
		if decoded[key] != want {
			t.Errorf("%s = %v, want %v", key, decoded[key], want)
		}
	}
}
