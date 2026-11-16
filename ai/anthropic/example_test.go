// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package anthropic_test

import (
	"fmt"
	"math/big"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/kaitencloud/sdk-go"
	kaitenanthropic "github.com/kaitencloud/sdk-go/ai/anthropic"
)

// Anthropic's usage carries no cost. A caller that knows what the call cost -- from its
// provider bill or its own price list, never from this module -- formats it once with
// sdk.VendorCostFromRat and hands it to the commit with the usage.
func ExampleUsage_vendorCost() {
	usage := anthropic.Usage{InputTokens: 900, CacheReadInputTokens: 100, OutputTokens: 250}

	// The caller's own figure for this call, in USD.
	cost, err := sdk.VendorCostFromRat(big.NewRat(37, 10_000))
	if err != nil {
		panic(err)
	}

	call := sdk.AICall{Usage: kaitenanthropic.Usage(usage), VendorCost: cost}
	fmt.Printf("%+v %s\n", call.Usage, call.VendorCost)
	// Output: {InputTokens:900 CacheReadInputTokens:100 CacheWriteInputTokens:0 OutputTokens:250 Units: UnitType: Duration:0s} 0.0037
}
