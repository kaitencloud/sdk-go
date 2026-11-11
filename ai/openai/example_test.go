// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package openai_test

import (
	"fmt"
	"math/big"

	"github.com/kaitencloud/sdk-go"
	kaitenopenai "github.com/kaitencloud/sdk-go/ai/openai"
	"github.com/openai/openai-go/v3"
)

// OpenAI's usage carries no cost. A caller that knows what the call cost -- from its
// provider bill or its own price list, never from this module -- formats it once with
// sdk.VendorCostFromRat and hands it to the commit with the usage.
func ExampleCompletionUsage_vendorCost() {
	var usage openai.CompletionUsage
	usage.PromptTokens, usage.CompletionTokens = 1200, 300
	usage.PromptTokensDetails.CachedTokens = 200

	// The caller's own figure for this call, in USD.
	cost, err := sdk.VendorCostFromRat(big.NewRat(4125, 1_000_000))
	if err != nil {
		panic(err)
	}

	call := sdk.AICall{Usage: kaitenopenai.CompletionUsage(usage), VendorCost: cost}
	fmt.Printf("%+v %s\n", call.Usage, call.VendorCost)
	// Output: {InputTokens:1000 CacheReadInputTokens:200 CacheWriteInputTokens:0 OutputTokens:300 Units: UnitType: Duration:0s} 0.004125
}
