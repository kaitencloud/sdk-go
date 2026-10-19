// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package anthropic_test

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/kaitencloud/sdk-go"
	kaitenanthropic "github.com/kaitencloud/sdk-go/ai/anthropic"
)

// decode reads a synthetic usage object the way the Anthropic SDK reads a response, so
// each case starts from the wire names, not from the SDK's Go field names.
func decode[T any](t *testing.T, raw string) T {
	t.Helper()

	var value T
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatalf("decode %s: %v", raw, err)
	}

	return value
}

func TestUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want sdk.AIUsage
	}{
		{
			name: "the cache counts carry over apart from the input",
			raw:  `{"input_tokens":50,"cache_read_input_tokens":2000,"cache_creation_input_tokens":300,"output_tokens":420}`,
			want: sdk.AIUsage{InputTokens: 50, CacheReadInputTokens: 2000, CacheWriteInputTokens: 300, OutputTokens: 420},
		},
		{
			name: "no cache",
			raw:  `{"input_tokens":12,"output_tokens":7}`,
			want: sdk.AIUsage{InputTokens: 12, OutputTokens: 7},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := kaitenanthropic.Usage(decode[anthropic.Usage](t, test.raw)); got != test.want {
				t.Errorf("Usage() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestMessageDeltaUsage(t *testing.T) {
	t.Parallel()

	usage := decode[anthropic.MessageDeltaUsage](t, `{"input_tokens":50,"cache_read_input_tokens":2000,"cache_creation_input_tokens":0,"output_tokens":512}`)

	want := sdk.AIUsage{InputTokens: 50, CacheReadInputTokens: 2000, OutputTokens: 512}
	if got := kaitenanthropic.MessageDeltaUsage(usage); got != want {
		t.Errorf("MessageDeltaUsage() = %+v, want %+v", got, want)
	}
}
