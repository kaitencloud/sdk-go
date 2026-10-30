// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

// This server-side example uses OpenFeature and its official OFREP provider,
// independently of the Kaiten management SDK. Never send the PAT to a browser.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/open-feature/go-sdk-contrib/providers/ofrep"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	baseURL := os.Getenv("KAITEN_BASE_URL") // Core gateway including /api, not /ofrep.
	token := os.Getenv("KAITEN_TOKEN")      // Organization PAT with read:feature_flags.
	targetingKey := os.Getenv("KAITEN_TARGETING_KEY")
	instanceSlug := os.Getenv("KAITEN_INSTANCE_SLUG")
	if baseURL == "" || token == "" || targetingKey == "" || instanceSlug == "" {
		log.Fatal("KAITEN_BASE_URL, KAITEN_TOKEN, KAITEN_TARGETING_KEY and KAITEN_INSTANCE_SLUG are required")
	}

	// The application owns the client, its timeout, transport and auth policy.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()
	httpClient := &http.Client{Timeout: 5 * time.Second, Transport: transport}
	provider := ofrep.NewProvider(baseURL, ofrep.WithBearerToken(token), ofrep.WithClient(httpClient))
	if err := openfeature.SetNamedProviderAndWait("kaiten", provider); err != nil {
		log.Fatal("register OFREP provider")
	}
	defer openfeature.Shutdown()
	flags := openfeature.NewClient("kaiten")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// In a server, derive these identifiers from authenticated, authorized
	// request state. Do not trust browser-supplied tenant or instance selectors.
	evaluationContext := openfeature.NewEvaluationContext(targetingKey, map[string]any{
		"kaiten": map[string]any{"instanceSlug": instanceSlug},
	})
	details, err := flags.BooleanValueDetails(ctx, "express-delivery", false, evaluationContext)
	if err != nil || details.Reason == openfeature.ErrorReason {
		// Check details too: an ERROR reason must never open the gate. Keep it
		// closed and record a sanitized failure, not credentials/context.
		log.Print("express-delivery evaluation failed; gate remains closed")
		return
	}
	fmt.Println("express-delivery:", details.Value)
}
