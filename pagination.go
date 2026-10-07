// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// maxPageLimit is the largest page the Core API's list endpoints accept, and
// so the fewest round trips a walk can cover a collection in. It is a hard
// ceiling, not a hint: the endpoints declare maximum:"200" on the limit
// parameter, so 201 comes back a 422 rather than clamped. Raising this
// constant without Core raising that maximum breaks every list call.
const maxPageLimit = 200

// page is the envelope every Core list endpoint returns. Unexported on
// purpose: the SDK's list methods promise "all of them" and walk the pages
// themselves, so a caller never holds one of these and never has to remember
// that a full answer takes more than one call.
//
// Decoded here rather than through internal/gen's own envelope types, because
// the walking below is what needs it and the generated per-endpoint types
// cannot be walked generically.
type page[T any] struct {
	Items      []T     `json:"items"`
	NextCursor *string `json:"nextCursor"`
	HasMore    bool    `json:"hasMore"`
}

// listAll walks path's pages and returns every row.
//
// Following the cursor here rather than exposing it is what keeps the list
// methods' signatures honest: Core's default page is 50 rows, so a single
// unpaged GET is a silently truncated answer waiting for the 51st row to be
// created.
//
// query carries the endpoint's own filters; limit and cursor are this
// function's to set.
func listAll[T any](ctx context.Context, transport listTransport, path string, query url.Values) ([]T, error) {
	return walkPages[T](ctx, transport, path, query, 0)
}

// listAtMost walks pages until it holds ceiling rows, for the one collection
// that grows without bound and whose callers ask for a slice of it rather than
// all of it. A ceiling of zero means no ceiling, so it behaves as listAll.
func listAtMost[T any](ctx context.Context, transport listTransport, path string, query url.Values, ceiling int32) ([]T, error) {
	return walkPages[T](ctx, transport, path, query, ceiling)
}

func walkPages[T any](ctx context.Context, transport listTransport, path string, query url.Values, ceiling int32) ([]T, error) {
	// Widened once here so the arithmetic below is all one type: a row count
	// can outgrow an int32 on paper, and every comparison against a ceiling
	// converted per iteration is another place for it to wrap.
	wanted := int(ceiling)

	// Non-nil even when empty: the endpoints promise "always an array,
	// possibly empty" and the generated decode this replaces returned an
	// empty slice, not nil.
	all := []T{}

	var cursor string
	for {
		current, err := getJSON[page[T]](ctx, transport, path, pageQuery(query, pageSize(wanted, len(all)), cursor))
		if err != nil {
			return nil, err
		}

		all = append(all, current.Items...)
		if wanted > 0 && len(all) >= wanted {
			return all[:wanted], nil
		}
		if !current.HasMore {
			return all, nil
		}

		cursor, err = advanceCursor(current.NextCursor, cursor, path)
		if err != nil {
			return nil, err
		}
	}
}

// pageSize is the limit to ask for, given a ceiling and the rows already held.
//
// Never more than is still wanted: a capped walk's last page is usually a
// partial one, and having Core read 200 rows to discard 190 of them is work
// nobody needed.
func pageSize(wanted, collected int) int {
	if wanted > 0 && wanted-collected < maxPageLimit {
		return wanted - collected
	}

	return maxPageLimit
}

// pageQuery merges the endpoint's own filters with the paging parameters this
// walk owns. The filters are copied rather than written through: the same
// url.Values backs every page of a walk, so a limit set on it would leak into
// the caller's map and into the next request.
func pageQuery(filters url.Values, limit int, cursor string) url.Values {
	values := url.Values{}
	maps.Copy(values, filters)

	values.Set("limit", strconv.Itoa(limit))
	if cursor != "" {
		values.Set("cursor", cursor)
	}

	return values
}

// advanceCursor returns the cursor for the next page, and fails when the server
// claimed more rows without a usable way to reach them.
//
// hasMore without a usable next cursor is a broken contract, and both ways it
// can break end the same way: a caller told it has every row when it does not.
// An error is the only answer that is not silently wrong -- the pagination
// package promises NextCursor is non-nil exactly when HasMore, and a cursor
// that does not move would loop forever.
func advanceCursor(next *string, current, path string) (string, error) {
	advanced := stringValue(next)
	if advanced == "" || advanced == current {
		return "", fmt.Errorf("list %s: server reported more rows but returned no usable next cursor", path)
	}

	return advanced, nil
}

// getJSON performs one GET against the Core API and decodes a 2xx body into
// T, mapping anything else onto the same *Error the generated calls produce
// -- callers branch on Error.StatusCode, so a hand-written request must not
// report failure any differently.
func getJSON[T any](ctx context.Context, transport listTransport, path string, query url.Values) (T, error) {
	var zero T

	requestURL, err := transport.resolve(path, query)
	if err != nil {
		return zero, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Accept", contentTypeJSON)

	for _, editor := range transport.editors {
		if err := editor(ctx, req); err != nil {
			return zero, err
		}
	}

	resp, err := transport.doer.Do(req)
	if err != nil {
		return zero, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return zero, responseError(apiResponse{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
			Problem:    decodeProblem(body),
		})
	}

	var decoded T
	if err := json.Unmarshal(body, &decoded); err != nil {
		return zero, err
	}

	return decoded, nil
}

// decodeProblem reads an RFC 9457 body, or reports nil when the body is not
// one. Only a problem that actually says something is attached: *Error falls
// back to printing the raw body, which is more use than an ErrorModel whose
// every field is empty.
func decodeProblem(body []byte) *ErrorModel {
	var problem ErrorModel
	if err := json.Unmarshal(body, &problem); err != nil {
		return nil
	}
	if stringValue(problem.Title) == "" && stringValue(problem.Detail) == "" {
		return nil
	}

	return &problem
}

// listPath builds an operation path from segments, escaping each one. Slugs
// reach these calls from callers rather than from the API, so they are
// escaped rather than trusted.
func listPath(segments ...string) string {
	escaped := make([]string, 0, len(segments))
	for _, segment := range segments {
		escaped = append(escaped, url.PathEscape(segment))
	}

	return "/" + strings.Join(escaped, "/")
}
