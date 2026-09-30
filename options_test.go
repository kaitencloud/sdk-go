package sdk

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// A request editor that fails must abort the call rather than send an unedited request:
// WithBearerToken is an editor, so "the editor errored" and "the request went out
// without its credential" have to stay distinguishable.
func TestWithRequestEditorFnAbortsTheCallWhenAnEditorFails(t *testing.T) {
	t.Parallel()

	editorErr := errors.New("no credential available")

	client, sent := givenClient(t, respondJSON(http.StatusOK, `{"items":[],"hasMore":false}`),
		WithRequestEditorFn(func(context.Context, *http.Request) error { return editorErr }))

	_, err := client.FeatureFlags.List(t.Context())
	if !errors.Is(err, editorErr) {
		t.Fatalf("List() error = %v, want %v", err, editorErr)
	}

	if sent.count() != 0 {
		t.Fatalf("requests = %d, want none: a failed editor must not send", sent.count())
	}
}

// WithBaseURL is resolved by the SDK and handed to the constructor, so its effect on
// the wire is worth pinning -- for both clients, since one Option type serves both.
func TestWithBaseURLOverridesTheConstructorArgument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		issue    func(*testing.T, string) *recorder
		override string
		wantHost string
		wantPath string
	}{
		{
			name:     "core client reaches the overridden core listener",
			override: "https://core.example.com/api",
			issue: func(t *testing.T, override string) *recorder {
				t.Helper()

				client, sent := givenClient(t, respondNoContent(), WithBaseURL(override))
				_, _ = client.ServiceAccounts.List(t.Context())

				return sent
			},
			wantHost: "core.example.com",
			wantPath: "/api/service-accounts",
		},
		{
			name:     "platform client reaches the overridden platform listener",
			override: "https://platform.example.com:3001/api",
			issue: func(t *testing.T, override string) *recorder {
				t.Helper()

				client, sent := givenPlatformClient(t, respondNoContent(), WithBaseURL(override))
				_, _ = client.Me(t.Context())

				return sent
			},
			wantHost: "platform.example.com:3001",
			wantPath: "/api/platform/me",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.issue(t, test.override).only(t)

			if got.host != test.wantHost {
				t.Errorf("host = %s, want %s", got.host, test.wantHost)
			}
			if got.path != test.wantPath {
				t.Errorf("path = %s, want %s", got.path, test.wantPath)
			}
		})
	}
}

// Asserted through a list call on purpose: the paginated lists build their requests by
// hand rather than through the generated client, so they are the one path that can lose
// the editors entirely -- and a list that skipped them would 401 against every real
// deployment while every test that checked only the URL kept passing.
func TestWithBearerTokenSetsTheAuthorizationHeader(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK, `{"items":[],"hasMore":false}`), WithBearerToken("ksh_secret"))

	if _, err := client.FeatureFlags.List(t.Context()); err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if got := sent.only(t).header.Get("Authorization"); got != "Bearer ksh_secret" {
		t.Fatalf("Authorization = %q, want %q", got, "Bearer ksh_secret")
	}
}

// A blank token sets no header at all, rather than sending "Bearer ". The two are
// different on the wire: kaiten answers the header-less request with the 401 that says
// "unauthenticated", and a malformed one with a parse failure that reads like an outage.
func TestWithBearerTokenSendsNoHeaderForABlankToken(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK, `{"items":[],"hasMore":false}`), WithBearerToken("   "))

	if _, err := client.FeatureFlags.List(t.Context()); err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if got, ok := sent.only(t).header["Authorization"]; ok {
		t.Fatalf("Authorization = %q, want the header absent", got)
	}
}
