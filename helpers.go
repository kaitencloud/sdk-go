// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kaitencloud/sdk-go/internal/genplatform"
)

// ErrThresholdExceeded is returned when reporting usage would exceed the configured entitlement limit.
var ErrThresholdExceeded = fmt.Errorf("usage threshold exceeded")

// ErrTransactionIDReused is returned when a usage report's TransactionID was
// already used, within Kaiten's idempotency window, for a report with another
// behavior or value. It is a bug in how the caller makes keys, not a transient
// failure: retrying the same report under the same key fails the same way. A
// correction is a new report under a new key.
var ErrTransactionIDReused = errors.New("usage report transaction ID already used for a different report")

// ErrInvalidTransactionID is returned, before any request, for a TransactionID
// Kaiten would refuse: it must be 1 to 128 characters of letters, digits, '.',
// '_', ':' and '-'.
var ErrInvalidTransactionID = errors.New("invalid usage report transaction ID")

// contentTypeJSON is the only media type this SDK sends or accepts. Every generated
// write wrapper takes it as a plain string argument, so naming it once keeps a typo
// from reaching one call site and turning into a 415 nothing else would catch.
const contentTypeJSON = "application/json"

// Error represents a non-2xx response from the Kaiten API.
type Error struct {
	StatusCode int
	Status     string
	Problem    *ErrorModel
	Body       []byte

	// Code is the stable, machine-readable error code (e.g. "License.NotFound"),
	// when the API supplied one. Branch on this rather than on Problem.Title or the
	// status code alone: a 409 from an entitlement limit and a 409 from a slug
	// collision are the same status but not the same thing, and the license family
	// endpoint answers two different 404s. Empty when the response carried no code.
	Code string

	// ErrorID correlates this response with the server-side log entry describing a
	// cause that was deliberately withheld from the client. Quote it in a bug report;
	// it is the only handle on the real error. Empty when none was supplied.
	ErrorID string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	status := strings.TrimSpace(e.Status)
	if status == "" && e.StatusCode != 0 {
		status = http.StatusText(e.StatusCode)
	}

	prefix := fmt.Sprintf("kaiten API error (%d %s)", e.StatusCode, status)
	if code := strings.TrimSpace(e.Code); code != "" {
		prefix += " [" + code + "]"
	}

	var suffix string
	if errorID := strings.TrimSpace(e.ErrorID); errorID != "" {
		suffix = " (errorId: " + errorID + ")"
	}

	if e.Problem != nil {
		title := strings.TrimSpace(stringValue(e.Problem.Title))
		detail := strings.TrimSpace(stringValue(e.Problem.Detail))

		switch {
		case title != "" && detail != "":
			return fmt.Sprintf("%s: %s: %s%s", prefix, title, detail, suffix)
		case detail != "":
			return fmt.Sprintf("%s: %s%s", prefix, detail, suffix)
		case title != "":
			return fmt.Sprintf("%s: %s%s", prefix, title, suffix)
		}
	}

	body := strings.TrimSpace(string(e.Body))
	if body != "" {
		return fmt.Sprintf("%s: %s%s", prefix, body, suffix)
	}

	if status != "" || e.Code != "" {
		return prefix + suffix
	}

	return "kaiten API error" + suffix
}

type statusResponse interface {
	StatusCode() int
	Status() string
}

type apiResponse struct {
	StatusCode int
	Status     string
	Body       []byte
	Problem    *ErrorModel
	Code       string
	ErrorID    string
}

func jsonBody(v any) (io.Reader, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

// newAPIResponse converts a generated Core response into the apiResponse the expect*
// helpers read. problems are the response's ApplicationproblemJSONnnn fields, and a
// wrapper passes all of them: a status left out reaches the caller as an *Error with
// Problem nil and Code empty, so TestClientWrappersReportEveryDeclaredProblem drives
// every wrapper through every status its response declares.
func newAPIResponse(resp statusResponse, body []byte, problems ...*ErrorModel) apiResponse {
	clonedBody := append([]byte(nil), body...)
	problem := firstProblem(problems...)

	converted := apiResponse{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Body:       clonedBody,
		Problem:    problem,
	}
	if problem == nil {
		return converted
	}

	// Lifted out of the problem for the same reason the Platform path does it:
	// Error.Code is documented as the member to branch on, and it was reaching
	// callers empty for every Core response. The Core Problem carries both
	// members now -- it did not when this was written, which is what the stale
	// note on Error.Code described -- and the two 404s the license family
	// endpoint answers with are told apart by nothing else.
	converted.Code = stringValue(problem.Code)
	converted.ErrorID = stringValue(problem.ErrorId)

	return converted
}

func firstProblem(problems ...*ErrorModel) *ErrorModel {
	for _, problem := range problems {
		if problem != nil {
			return problem
		}
	}

	return nil
}

// newPlatformAPIResponse is newAPIResponse for the Platform API, whose generated
// package has its own Problem type. It is a strict superset of the Core API's
// ErrorModel -- same six RFC 9457 fields, plus a machine-readable code and a log
// correlation id -- so the two collapse onto one apiResponse and callers keep seeing a
// single *Error. Converting this way rather than widening Error.Problem to an interface
// keeps errors.As(err, &sdk.Error{}) and every existing Problem field access working.
//
// As with newAPIResponse, problems are all of the response's ApplicationproblemJSONnnn
// fields; TestPlatformClientWrappersReportEveryDeclaredProblem holds the wrappers to it.
func newPlatformAPIResponse(resp statusResponse, body []byte, problems ...*genplatform.Problem) apiResponse {
	converted := apiResponse{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Body:       append([]byte(nil), body...),
	}

	problem := firstPlatformProblem(problems...)
	if problem == nil {
		return converted
	}

	converted.Problem = &ErrorModel{
		Detail:   problem.Detail,
		Errors:   convertPlatformErrorDetails(problem.Errors),
		Instance: problem.Instance,
		Status:   problem.Status,
		Title:    problem.Title,
		Type:     problem.Type,
	}
	converted.Code = stringValue(problem.Code)
	converted.ErrorID = stringValue(problem.ErrorId)

	return converted
}

func firstPlatformProblem(problems ...*genplatform.Problem) *genplatform.Problem {
	for _, problem := range problems {
		if problem != nil {
			return problem
		}
	}

	return nil
}

func convertPlatformErrorDetails(details *[]genplatform.ErrorDetail) *[]ErrorDetail {
	if details == nil {
		return nil
	}

	converted := make([]ErrorDetail, 0, len(*details))
	for _, detail := range *details {
		converted = append(converted, ErrorDetail(detail))
	}

	return &converted
}

func expectJSON[T any](resp apiResponse, payload *T) (T, error) {
	var zero T
	if payload != nil {
		return *payload, nil
	}

	return zero, responseError(resp)
}

func expectSlice[T any](resp apiResponse, payload *[]T) ([]T, error) {
	if payload != nil {
		return *payload, nil
	}

	return nil, responseError(resp)
}

func expectNoContent(resp apiResponse) error {
	if code := resp.StatusCode; code >= http.StatusOK && code < http.StatusMultipleChoices {
		return nil
	}

	return responseError(resp)
}

func responseError(resp apiResponse) error {
	err := &Error{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Problem:    resp.Problem,
		Body:       resp.Body,
		Code:       resp.Code,
		ErrorID:    resp.ErrorID,
	}

	return err
}

func cloneStringSlice(values []string) *[]string {
	if values == nil {
		return nil
	}

	cloned := append([]string(nil), values...)
	return &cloned
}

func cloneMap(values map[string]any) map[string]any {
	if values == nil {
		return nil
	}

	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}

func cloneMapPtr(values map[string]any) *map[string]any {
	if values == nil {
		return nil
	}

	cloned := cloneMap(values)
	return &cloned
}

func timeStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
