// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk_test

import (
	"reflect"
	"strings"
	"testing"

	sdk "github.com/kaitencloud/sdk-go"
)

// Use only the public module, as a consumer does. The internal generated types
// are reachable through Raw even though consumers cannot import their package.
func TestPublicClientsDoNotExposeFlagEvaluation(t *testing.T) {
	t.Parallel()

	client, err := sdk.NewClient("https://kaiten.example.com/api")
	if err != nil {
		t.Fatal(err)
	}
	platform, err := sdk.NewPlatformClient("https://platform.kaiten.example.com/api")
	if err != nil {
		t.Fatal(err)
	}
	for _, service := range []any{client, client.Raw(), client.Raw().ClientInterface, client.FeatureFlags, platform, platform.Raw(), platform.Raw().ClientInterface} {
		typ := reflect.TypeOf(service)
		for i := range typ.NumMethod() {
			name := typ.Method(i).Name
			if strings.Contains(name, "EvaluateFlag") || strings.Contains(name, "OfrepEvaluate") {
				t.Errorf("%v exposes direct flag evaluation method %s", typ, name)
			}
		}
	}
}
