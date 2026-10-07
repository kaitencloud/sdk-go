// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"errors"
	"net/http"
	"testing"
)

// Error.Error is what reaches a log line or a CLI's stderr when a call fails, and it is
// assembled from five optional parts. Each case below is one way a real response arrives
// -- the Core API sends no code and no errorId, the Platform API sends both, and a
// gateway in front of either sends neither a problem document nor a useful status text.
func TestErrorMessage(t *testing.T) {
	t.Parallel()

	problem := func(title, detail string) *ErrorModel {
		model := &ErrorModel{}
		if title != "" {
			model.Title = &title
		}
		if detail != "" {
			model.Detail = &detail
		}

		return model
	}

	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "nil error renders empty rather than panicking",
			err:      nil,
			expected: "",
		},
		{
			name: "title and detail are both shown",
			err: &Error{
				StatusCode: http.StatusConflict,
				Status:     "409 Conflict",
				Problem:    problem("Conflict", "threshold reached"),
			},
			expected: "kaiten API error (409 409 Conflict): Conflict: threshold reached",
		},
		{
			name: "detail alone stands in for the title",
			err: &Error{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Problem:    problem("", "no such license"),
			},
			expected: "kaiten API error (404 404 Not Found): no such license",
		},
		{
			name: "title alone stands in for the detail",
			err: &Error{
				StatusCode: http.StatusForbidden,
				Status:     "403 Forbidden",
				Problem:    problem("Forbidden", ""),
			},
			expected: "kaiten API error (403 403 Forbidden): Forbidden",
		},
		{
			name: "a missing status line is filled in from the status code",
			err: &Error{
				StatusCode: http.StatusBadGateway,
				Problem:    problem("Bad Gateway", ""),
			},
			expected: "kaiten API error (502 Bad Gateway): Bad Gateway",
		},
		{
			name: "the machine-readable code and the correlation id both survive",
			err: &Error{
				StatusCode: http.StatusConflict,
				Status:     "409 Conflict",
				Code:       "Token.NameConflict",
				ErrorID:    "3f6b1a2c",
				Problem:    problem("Conflict", "already active"),
			},
			expected: "kaiten API error (409 409 Conflict) [Token.NameConflict]: Conflict: already active (errorId: 3f6b1a2c)",
		},
		{
			name: "a body with no problem document is printed raw",
			err: &Error{
				StatusCode: http.StatusInternalServerError,
				Status:     "500 Internal Server Error",
				Body:       []byte("upstream connect error"),
			},
			expected: "kaiten API error (500 500 Internal Server Error): upstream connect error",
		},
		{
			name: "a response carrying nothing but a status still names it",
			err: &Error{
				StatusCode: http.StatusTeapot,
				Status:     "418 I'm a teapot",
			},
			expected: "kaiten API error (418 418 I'm a teapot)",
		},
		{
			name:     "a response carrying nothing at all still identifies kaiten",
			err:      &Error{},
			expected: "kaiten API error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.err.Error(); got != test.expected {
				t.Errorf("Error() = %q, want %q", got, test.expected)
			}
		})
	}
}

func TestResponseErrorPreservesTheProblemDetails(t *testing.T) {
	t.Parallel()

	title, detail := "Conflict", "threshold reached"

	err := responseError(apiResponse{
		StatusCode: http.StatusConflict,
		Status:     "409 Conflict",
		Body:       []byte(`{"detail":"threshold reached"}`),
		Problem:    &ErrorModel{Title: &title, Detail: &detail},
	})

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("responseError() returned %T, want *Error", err)
	}

	if apiErr.StatusCode != http.StatusConflict {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusConflict)
	}

	if stringValue(apiErr.Problem.Detail) != detail {
		t.Errorf("Problem.Detail = %q, want %q", stringValue(apiErr.Problem.Detail), detail)
	}
}

// newAPIResponse copies the body rather than aliasing it. The generated clients hand
// over a slice they own and go on reusing, so an *Error holding the original would print
// whatever the next response left there.
func TestNewAPIResponseCopiesTheBody(t *testing.T) {
	t.Parallel()

	body := []byte(`{"detail":"first"}`)

	resp := newAPIResponse(stubStatusResponse{code: http.StatusConflict, status: "409 Conflict"}, body)
	copy(body, []byte(`{"detail":"BBBBB"}`))

	if string(resp.Body) != `{"detail":"first"}` {
		t.Fatalf("Body = %s, want the value captured at the time of the call", resp.Body)
	}
}

type stubStatusResponse struct {
	code   int
	status string
}

func (s stubStatusResponse) StatusCode() int { return s.code }
func (s stubStatusResponse) Status() string  { return s.status }
