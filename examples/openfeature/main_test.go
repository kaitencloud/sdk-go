// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/open-feature/go-sdk-contrib/providers/ofrep"
	"github.com/open-feature/go-sdk/openfeature"
)

// Drive the actual example entry point too: the guarded value must be the one
// it uses, not only a copied assertion in a provider contract test. These tests
// are serial because main owns process-wide OpenFeature registration and stdout.
func TestExampleMain(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"targeted", http.StatusOK, `{"key":"express-delivery","value":true,"reason":"TARGETING_MATCH"}`, "express-delivery: true\n"},
		{"default-true", http.StatusOK, `{"key":"express-delivery","value":true,"reason":"DEFAULT"}`, "express-delivery: true\n"},
		{"default-false", http.StatusOK, `{"key":"express-delivery","value":false,"reason":"DEFAULT"}`, "express-delivery: false\n"},
		{"error-reason", http.StatusOK, `{"key":"express-delivery","value":true,"reason":"ERROR"}`, ""},
		{"unauthorized", http.StatusUnauthorized, "", ""},
		{"forbidden", http.StatusForbidden, "", ""},
		{"server-error", http.StatusInternalServerError, `{"errorCode":"GENERAL"}`, ""},
		{"malformed", http.StatusOK, "not-json", ""},
		{"wrong-type", http.StatusOK, `{"key":"express-delivery","value":"true","reason":"STATIC"}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPost || req.URL.Path != "/api/ofrep/v1/evaluate/flags/express-delivery" || req.Header.Get("Authorization") != "Bearer fixture-pat" {
					t.Error("example lost the official provider path or authentication")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			t.Setenv("KAITEN_BASE_URL", server.URL+"/api")
			t.Setenv("KAITEN_TOKEN", "fixture-pat")
			t.Setenv("KAITEN_TARGETING_KEY", "customer-a")
			t.Setenv("KAITEN_INSTANCE_SLUG", "customer-a-prod")
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			defer writer.Close()
			stdout := os.Stdout
			os.Stdout = writer
			defer func() { os.Stdout = stdout }()
			main()
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			output, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			if string(output) != tc.want {
				t.Fatalf("example output = %q; want %q", output, tc.want)
			}
		})
	}
}

// No live Core: the real official provider talks to an isolated protocol fixture.
// Two request contexts must stay distinct even on a reused OpenFeature client.
func TestOfficialProviderContract(t *testing.T) {
	var contexts []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost || req.URL.Path != "/api/ofrep/v1/evaluate/flags/express-delivery" {
			t.Errorf("unexpected method/path: %s %s", req.Method, req.URL.Path)
		}
		if req.Header.Get("Authorization") != "Bearer fixture-pat" || req.Header.Get("Content-Type") != "application/json" {
			t.Error("missing PAT or JSON content type")
		}
		var payload struct {
			Context map[string]any `json:"context"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		contexts = append(contexts, payload.Context)
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"key":"express-delivery","value":true,"reason":"TARGETING_MATCH","variant":"on"}`); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	httpClient := server.Client()
	httpClient.Timeout = 2 * time.Second
	provider := ofrep.NewProvider(server.URL+"/api", ofrep.WithBearerToken("fixture-pat"), ofrep.WithClient(httpClient))
	if err := openfeature.SetNamedProviderAndWait(t.Name(), provider); err != nil {
		t.Fatal(err)
	}
	defer openfeature.Shutdown()
	flags := openfeature.NewClient(t.Name())
	for _, customer := range []string{"customer-a", "customer-b"} {
		evalCtx := openfeature.NewEvaluationContext(customer, map[string]any{
			"kaiten": map[string]any{"instanceSlug": customer + "-prod"},
		})
		enabled, err := flags.BooleanValue(t.Context(), "express-delivery", false, evalCtx)
		if err != nil || !enabled {
			t.Fatalf("BooleanValue() = %v, %v; want true, nil", enabled, err)
		}
	}
	want := []map[string]any{
		{"targetingKey": "customer-a", "kaiten": map[string]any{"instanceSlug": "customer-a-prod"}},
		{"targetingKey": "customer-b", "kaiten": map[string]any{"instanceSlug": "customer-b-prod"}},
	}
	if !reflect.DeepEqual(contexts, want) {
		t.Errorf("request contexts = %v; want %v", contexts, want)
	}
	if httpClient.Timeout != 2*time.Second {
		t.Fatal("provider changed the caller-owned HTTP timeout")
	}
}

// DEFAULT is a valid resolved value, not necessarily an error fallback. Preserve
// both true and false while refusing ERROR even without an accompanying Go error.
func TestOfficialProviderReasonGuard(t *testing.T) {
	for _, tc := range []struct {
		name   string
		body   string
		closed bool
		value  bool
	}{
		{"default-true", `{"key":"express-delivery","value":true,"reason":"DEFAULT"}`, false, true},
		{"default-false", `{"key":"express-delivery","value":false,"reason":"DEFAULT"}`, false, false},
		{"error-without-go-error", `{"key":"express-delivery","value":true,"reason":"ERROR"}`, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			provider := ofrep.NewProvider(server.URL+"/api", ofrep.WithClient(server.Client()))
			if err := openfeature.SetNamedProviderAndWait(t.Name(), provider); err != nil {
				t.Fatal(err)
			}
			defer openfeature.Shutdown()
			flags := openfeature.NewClient(t.Name())
			details, err := flags.BooleanValueDetails(t.Context(), "express-delivery", false, openfeature.NewEvaluationContext("customer-a", nil))
			if err != nil || details.Value != tc.value {
				t.Fatalf("BooleanValueDetails() = %v, %v", details, err)
			}
			if closed := err != nil || details.Reason == openfeature.ErrorReason; closed != tc.closed {
				t.Fatalf("gate closed = %v; want %v (details: %v)", closed, tc.closed, details)
			}
		})
	}
}

func TestOfficialProviderOutageKeepsGateClosed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("closed server received a request")
	}))
	httpClient := server.Client()
	httpClient.Timeout = time.Second
	provider := ofrep.NewProvider(server.URL+"/api", ofrep.WithClient(httpClient))
	if err := openfeature.SetNamedProviderAndWait(t.Name(), provider); err != nil {
		server.Close()
		t.Fatal(err)
	}
	defer openfeature.Shutdown()
	server.Close()
	flags := openfeature.NewClient(t.Name())
	enabled, err := flags.BooleanValue(t.Context(), "express-delivery", false, openfeature.NewEvaluationContext("customer-a", nil))
	if err == nil || enabled {
		t.Fatalf("outage evaluation = %v, %v; want false and an error", enabled, err)
	}
}

func TestOfficialProviderFailuresKeepGateClosed(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"unauthorized", http.StatusUnauthorized, ""},
		{"forbidden", http.StatusForbidden, ""},
		{"not-found", http.StatusNotFound, `{"key":"express-delivery","errorCode":"FLAG_NOT_FOUND"}`},
		{"server-error", http.StatusInternalServerError, `{"key":"express-delivery","errorCode":"GENERAL"}`},
		{"malformed", http.StatusOK, "not-json"},
		{"wrong-type", http.StatusOK, `{"key":"express-delivery","value":"true","reason":"STATIC"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			provider := ofrep.NewProvider(server.URL+"/api", ofrep.WithClient(server.Client()))
			if err := openfeature.SetNamedProviderAndWait(t.Name(), provider); err != nil {
				t.Fatal(err)
			}
			defer openfeature.Shutdown()
			flags := openfeature.NewClient(t.Name())
			enabled, err := flags.BooleanValue(t.Context(), "express-delivery", false, openfeature.NewEvaluationContext("customer-a", nil))
			if err == nil || enabled {
				t.Fatalf("BooleanValue() = %v, %v; want false and an error", enabled, err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			enabled, err = flags.BooleanValue(ctx, "express-delivery", false, openfeature.NewEvaluationContext("customer-a", nil))
			if err == nil || enabled {
				t.Fatalf("canceled evaluation = %v, %v; want false and an error", enabled, err)
			}
		})
	}
}
