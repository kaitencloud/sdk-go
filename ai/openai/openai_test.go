// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package openai_test

import (
	"encoding/json"
	"testing"

	"github.com/kaitencloud/sdk-go"
	kaitenopenai "github.com/kaitencloud/sdk-go/ai/openai"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// decode reads a synthetic usage object the way the OpenAI SDK reads a response, so
// each case below starts from the wire names, not from the SDK's Go field names.
func decode[T any](t *testing.T, raw string) T {
	t.Helper()

	var value T
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatalf("decode %s: %v", raw, err)
	}

	return value
}

func TestCompletionUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want sdk.AIUsage
	}{
		{
			name: "cached tokens leave the input count",
			raw:  `{"prompt_tokens":1200,"completion_tokens":300,"total_tokens":1500,"prompt_tokens_details":{"cached_tokens":1024},"completion_tokens_details":{"reasoning_tokens":120}}`,
			want: sdk.AIUsage{InputTokens: 176, CacheReadInputTokens: 1024, OutputTokens: 300},
		},
		{
			name: "no details is no cache",
			raw:  `{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}`,
			want: sdk.AIUsage{InputTokens: 10, OutputTokens: 5},
		},
		{
			name: "more cached than prompt clamps to zero",
			raw:  `{"prompt_tokens":10,"completion_tokens":0,"total_tokens":10,"prompt_tokens_details":{"cached_tokens":12}}`,
			want: sdk.AIUsage{CacheReadInputTokens: 12},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := kaitenopenai.CompletionUsage(decode[openai.CompletionUsage](t, test.raw)); got != test.want {
				t.Errorf("CompletionUsage() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestResponseUsage(t *testing.T) {
	t.Parallel()

	usage := decode[responses.ResponseUsage](t, `{"input_tokens":2048,"input_tokens_details":{"cached_tokens":2000},"output_tokens":512,"output_tokens_details":{"reasoning_tokens":256},"total_tokens":2560}`)

	want := sdk.AIUsage{InputTokens: 48, CacheReadInputTokens: 2000, OutputTokens: 512}
	if got := kaitenopenai.ResponseUsage(usage); got != want {
		t.Errorf("ResponseUsage() = %+v, want %+v", got, want)
	}
}

func TestEmbeddingUsage(t *testing.T) {
	t.Parallel()

	usage := decode[openai.CreateEmbeddingResponseUsage](t, `{"prompt_tokens":64,"total_tokens":64}`)

	want := sdk.AIUsage{InputTokens: 64}
	if got := kaitenopenai.EmbeddingUsage(usage); got != want {
		t.Errorf("EmbeddingUsage() = %+v, want %+v", got, want)
	}
}
