// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

// Package openai converts the usage the OpenAI Go SDK reports into the sdk.AIUsage a
// Kaiten AI credit commit takes.
//
// It is a module of its own so that github.com/kaitencloud/sdk-go never depends on a
// provider SDK: only a program that imports this package pulls in
// github.com/openai/openai-go/v3. Import it under a name of its own, since its package
// name is the provider's:
//
//	import kaitenopenai "github.com/kaitencloud/sdk-go/ai/openai"
//
// OpenAI counts cached input tokens inside its prompt tokens, where Kaiten counts them
// apart: every function here subtracts the cached tokens from the input ones. Reasoning
// tokens are already inside OpenAI's output tokens and stay there.
package openai

import (
	"github.com/kaitencloud/sdk-go"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// CompletionUsage converts the usage of a chat completion, or of the final chunk of a
// chat completion stream sent with stream_options.include_usage.
func CompletionUsage(usage openai.CompletionUsage) sdk.AIUsage {
	cached := usage.PromptTokensDetails.CachedTokens

	return sdk.AIUsage{
		InputTokens:          nonNegative(usage.PromptTokens - cached),
		CacheReadInputTokens: nonNegative(cached),
		OutputTokens:         nonNegative(usage.CompletionTokens),
	}
}

// ResponseUsage converts the usage of a Responses API response.
func ResponseUsage(usage responses.ResponseUsage) sdk.AIUsage {
	cached := usage.InputTokensDetails.CachedTokens

	return sdk.AIUsage{
		InputTokens:          nonNegative(usage.InputTokens - cached),
		CacheReadInputTokens: nonNegative(cached),
		OutputTokens:         nonNegative(usage.OutputTokens),
	}
}

// EmbeddingUsage converts the usage of an embeddings response, which has input only.
func EmbeddingUsage(usage openai.CreateEmbeddingResponseUsage) sdk.AIUsage {
	return sdk.AIUsage{InputTokens: nonNegative(usage.PromptTokens)}
}

// nonNegative clamps a count to zero: a provider reporting more cached tokens than
// input tokens is a provider bug, and a negative count is one the API refuses.
func nonNegative(count int64) int64 {
	return max(count, 0)
}
