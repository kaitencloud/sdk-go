// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

// Package anthropic converts the usage the Anthropic Go SDK reports into the
// sdk.AIUsage a Kaiten AI credit commit takes.
//
// It is a module of its own so that github.com/kaitencloud/sdk-go never depends on a
// provider SDK: only a program that imports this package pulls in
// github.com/anthropics/anthropic-sdk-go. Import it under a name of its own, since its
// package name is the provider's:
//
//	import kaitenanthropic "github.com/kaitencloud/sdk-go/ai/anthropic"
//
// Anthropic already counts input tokens apart from the cache, as Kaiten does, so the
// counts carry over one for one.
package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/kaitencloud/sdk-go"
)

// Usage converts the usage of a message: the answer to Messages.New, or a stream's
// message once accumulated.
func Usage(usage anthropic.Usage) sdk.AIUsage {
	return convert(usage.InputTokens, usage.CacheReadInputTokens, usage.CacheCreationInputTokens, usage.OutputTokens)
}

// MessageDeltaUsage converts the usage of a stream's message_delta event, whose counts
// are cumulative: the last one of a stream is the stream's usage.
func MessageDeltaUsage(usage anthropic.MessageDeltaUsage) sdk.AIUsage {
	return convert(usage.InputTokens, usage.CacheReadInputTokens, usage.CacheCreationInputTokens, usage.OutputTokens)
}

func convert(input, cacheRead, cacheWrite, output int64) sdk.AIUsage {
	return sdk.AIUsage{
		InputTokens:           max(input, 0),
		CacheReadInputTokens:  max(cacheRead, 0),
		CacheWriteInputTokens: max(cacheWrite, 0),
		OutputTokens:          max(output, 0),
	}
}
