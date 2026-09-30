package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"maps"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// This is the one _test.go file with no source file of the same name, because what it
// holds belongs to no single one: a fake listener, the client constructors that point at
// one, and the reflection helpers the public-surface and problem-coverage tests share.
// Every other test file is named after the source file it exercises.

const (
	coreBaseURL     = "https://example.com/api"
	platformBaseURL = "https://example.com:3001/api"
)

// roundTripperFunc serves a canned response without opening a socket.
type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// responder answers one request reaching the fake listener. It is the raw
// http.RoundTripper signature rather than a canned value, because several tests need
// the answer to depend on the request -- a paged endpoint reads the cursor, and a
// create answers 201 where an update answers 200.
type responder func(req *http.Request) (*http.Response, error)

// respondJSON answers every request with status and body.
func respondJSON(status int, body string) responder {
	return func(req *http.Request) (*http.Response, error) {
		return jsonResponse(req, status, body), nil
	}
}

// respondNoContent answers every request with 204 and an empty body.
//
// It is enough for the wrappers that expect no content, and for the ones that do expect
// a body it fails them after the request was already built and sent -- which is all a
// request-shape test asserts on.
func respondNoContent() responder {
	return func(req *http.Request) (*http.Response, error) {
		return jsonResponse(req, http.StatusNoContent, ""), nil
	}
}

func jsonResponse(req *http.Request, status int, body string) *http.Response {
	header := http.Header{}
	if body != "" {
		header.Set("Content-Type", contentTypeJSON)
	}

	return &http.Response{
		StatusCode: status,
		// The generated client prints this status line into *Error, and several
		// tests assert the rendered message, so it has to read like a real one.
		Status:  fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:  header,
		Body:    io.NopCloser(bytes.NewReader([]byte(body))),
		Request: req,
	}
}

// recordedRequest is what the fake listener saw, so a test can assert on the wire
// rather than only on the value that came back. The body is copied out because the
// request's own is drained by the time an assertion runs.
type recordedRequest struct {
	method string
	host   string
	path   string
	query  url.Values
	header http.Header
	body   []byte
}

// decodeBody reads the recorded body as a generic JSON object.
func (r recordedRequest) decodeBody(t *testing.T) map[string]any {
	t.Helper()

	var decoded map[string]any
	if err := json.Unmarshal(r.body, &decoded); err != nil {
		t.Fatalf("unmarshal request body %q: %v", r.body, err)
	}

	return decoded
}

// bodyKeys returns the recorded body's top-level keys, sorted, which is how the
// request-body tests compare against a key set transcribed from the spec.
func (r recordedRequest) bodyKeys(t *testing.T) []string {
	t.Helper()

	return slices.Sorted(maps.Keys(r.decodeBody(t)))
}

// recorder collects every request the fake listener received.
//
// No mutex: net/http calls RoundTrip on the goroutine that called Do, and each test
// drives its own client from its own goroutine, so these appends stay sequential even
// when the tests themselves run in parallel.
type recorder struct {
	requests []recordedRequest
}

func (r *recorder) count() int {
	return len(r.requests)
}

// at returns the nth recorded request, failing rather than panicking when the walk made
// fewer calls than the test expected.
func (r *recorder) at(t *testing.T, n int) recordedRequest {
	t.Helper()

	if n >= len(r.requests) {
		t.Fatalf("request %d was never issued: only %d requests were made", n, len(r.requests))
	}

	return r.requests[n]
}

// only returns the single request the test expected the call to make.
func (r *recorder) only(t *testing.T) recordedRequest {
	t.Helper()

	if len(r.requests) != 1 {
		t.Fatalf("requests = %d, want exactly 1", len(r.requests))
	}

	return r.requests[0]
}

// givenClient builds a Core API client whose requests never leave the process: respond
// answers each one, and the returned recorder holds what was sent.
func givenClient(t *testing.T, respond responder, opts ...Option) (*Client, *recorder) {
	t.Helper()

	rec := &recorder{}

	client, err := NewClient(coreBaseURL, append([]Option{WithHTTPClient(recordingDoer(rec, respond))}, opts...)...)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	return client, rec
}

// givenPlatformClient is givenClient for the Platform API listener, which is a separate
// deployment with its own base URL and its own credential class.
func givenPlatformClient(t *testing.T, respond responder, opts ...Option) (*PlatformClient, *recorder) {
	t.Helper()

	rec := &recorder{}

	client, err := NewPlatformClient(platformBaseURL, append([]Option{WithHTTPClient(recordingDoer(rec, respond))}, opts...)...)
	if err != nil {
		t.Fatalf("NewPlatformClient() error = %v", err)
	}

	return client, rec
}

func recordingDoer(rec *recorder, respond responder) *http.Client {
	return &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		captured := recordedRequest{
			method: req.Method,
			host:   req.URL.Host,
			path:   req.URL.Path,
			query:  req.URL.Query(),
			header: req.Header.Clone(),
		}

		if req.Body != nil {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			captured.body = body
			// Put it back: respond may read the body too, and a paged responder
			// re-reads the request it was handed.
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		rec.requests = append(rec.requests, captured)

		return respond(req)
	})}
}

// serviceCoverage is the expected public surface of one namespace: the value itself, so
// the test reflects over the real type, and the methods it is allowed to expose.
type serviceCoverage struct {
	service any
	methods []string
}

// assertServiceCoverage fails on a public method that no test lists and on a listed method
// that does not exist. Adding a wrapper without listing it here is the mistake this
// catches -- an unreviewed method is how the unreachable /organizations paths survived.
func assertServiceCoverage(t *testing.T, expected map[string]serviceCoverage) {
	t.Helper()

	// Sorted, and one subtest per namespace, so a failure names the namespace it came
	// from and the report does not shift with Go's map iteration order.
	for _, serviceName := range slices.Sorted(maps.Keys(expected)) {
		item := expected[serviceName]

		t.Run(serviceName, func(t *testing.T) {
			t.Parallel()

			serviceType := reflect.TypeOf(item.service)

			declared := make(map[string]struct{}, serviceType.NumMethod())
			for i := range serviceType.NumMethod() {
				declared[serviceType.Method(i).Name] = struct{}{}
			}

			for _, methodName := range item.methods {
				if _, ok := declared[methodName]; !ok {
					t.Errorf("%s is missing method %s", serviceName, methodName)
				}
				delete(declared, methodName)
			}

			for methodName := range declared {
				t.Errorf("%s has untracked public method %s", serviceName, methodName)
			}
		})
	}
}

// problemCoverage is one wrapper and the generated response it decodes: the value, so
// the test reflects over the real struct for the problem statuses it declares, and a
// call that drives the wrapper, so each of those statuses can be answered through it.
type problemCoverage[C any] struct {
	response any
	call     func(context.Context, C) error
}

// problemSite is what the source says about one wrapper: which of the two conversions
// it calls, and the generated response it hands over -- named after the generated
// method, since XWithResponse and XWithBodyWithResponse both return *XResponse.
type problemSite struct {
	platform bool
	response string
}

// assertProblemCoverage fails on a wrapper that hands newAPIResponse or
// newPlatformAPIResponse fewer problems than its generated response declares.
//
// The declared set is read off the response struct by reflection -- one
// ApplicationproblemJSONnnn field per `application/problem+json` response the spec puts
// on the operation -- and each status is then answered through the wrapper with a
// problem document carrying a code. A status the wrapper leaves out arrives as an *Error
// with Problem nil and Code empty, which is the failure reported here: Code is the
// member callers are told to branch on. Reading the set off the struct rather than
// listing it is what keeps the check current. The client is regenerated, a new status
// appears on the struct, and the wrapper still passing its old list fails with no test
// edited.
//
// The table is checked against the source too, by parsing the package: every method
// that calls one of the conversions must be listed, so a new wrapper is either driven
// or reported; every listed method must be one; and the response it converts must be
// the type listed, or a pasted entry would pin another operation's statuses.
func assertProblemCoverage[C any](t *testing.T, given func(*testing.T, responder) C, wrappers map[string]problemCoverage[C]) {
	t.Helper()

	// Which conversion the wrappers call follows from the client they take: a
	// PlatformClient's wrappers convert genplatform responses, a Client's convert gen.
	var zero C
	_, platform := any(zero).(*PlatformClient)

	sites := scanProblemSites(t)

	for _, name := range slices.Sorted(maps.Keys(sites)) {
		if _, listed := wrappers[name]; sites[name].platform == platform && !listed {
			t.Errorf("%s converts a generated response but is not listed, so nothing checks the statuses it declares", name)
		}
	}

	for _, name := range slices.Sorted(maps.Keys(wrappers)) {
		item := wrappers[name]

		site, ok := sites[name]
		if !ok || site.platform != platform {
			t.Errorf("%s is listed, but no such wrapper converts a generated response", name)
			continue
		}

		responseType := reflect.TypeOf(item.response)
		if responseType.Kind() == reflect.Pointer {
			responseType = responseType.Elem()
		}
		if responseType.Name() != site.response {
			t.Errorf("%s is listed with %s but converts a %s", name, responseType.Name(), site.response)
			continue
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for _, status := range declaredProblemStatuses(t, responseType) {
				// The code names the wrapper and the status, so no answer meant for
				// another case could satisfy this one.
				code := fmt.Sprintf("%s.%d", name, status)
				client := given(t, respondJSON(status,
					fmt.Sprintf(`{"title":%q,"status":%d,"code":%q,"detail":"refused"}`, http.StatusText(status), status, code)))

				err := item.call(t.Context(), client)

				var apiErr *Error
				if !errors.As(err, &apiErr) {
					t.Errorf("%d: error = %v, want an *Error", status, err)
					continue
				}
				if apiErr.StatusCode != status {
					t.Errorf("%d: StatusCode = %d", status, apiErr.StatusCode)
				}
				if apiErr.Problem == nil || apiErr.Code != code {
					t.Errorf("%d: Code = %q, want %q: ApplicationproblemJSON%d is not passed through", status, apiErr.Code, code, status)
				}
			}
		})
	}
}

// declaredProblemStatuses reads the statuses a generated response can decode a problem
// for. oapi-codegen renders one ApplicationproblemJSONnnn field per
// `application/problem+json` response the spec declares on the operation, so the struct
// is the declared set, kept current by regeneration.
func declaredProblemStatuses(t *testing.T, response reflect.Type) []int {
	t.Helper()

	var statuses []int
	for i := range response.NumField() {
		field := response.Field(i).Name

		digits, ok := strings.CutPrefix(field, "ApplicationproblemJSON")
		if !ok {
			continue
		}

		status, err := strconv.Atoi(digits)
		if err != nil {
			t.Fatalf("%s.%s: %v", response.Name(), field, err)
		}
		statuses = append(statuses, status)
	}

	if len(statuses) == 0 {
		t.Fatalf("%s declares no problem responses", response.Name())
	}

	return statuses
}

// scanProblemSites parses the package's sources and returns every method that calls
// newAPIResponse or newPlatformAPIResponse, keyed "Receiver.Method", with what it
// converts. Test files are skipped: a test calling a conversion is not a wrapper.
func scanProblemSites(t *testing.T) map[string]problemSite {
	t.Helper()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("list sources: %v", err)
	}

	sites := map[string]problemSite{}
	fset := token.NewFileSet()

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		parsed, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		for _, decl := range parsed.Decls {
			method, ok := decl.(*ast.FuncDecl)
			if !ok || method.Recv == nil || method.Body == nil {
				continue
			}

			name := receiverTypeName(t, method) + "." + method.Name.Name

			conversions := 0
			platform := false
			var responses []string
			ast.Inspect(method.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}

				switch fun := call.Fun.(type) {
				case *ast.Ident:
					switch fun.Name {
					case "newAPIResponse":
						conversions++
					case "newPlatformAPIResponse":
						conversions++
						platform = true
					}
				case *ast.SelectorExpr:
					if operation, ok := strings.CutSuffix(fun.Sel.Name, "WithResponse"); ok {
						responses = append(responses, strings.TrimSuffix(operation, "WithBody")+"Response")
					}
				}

				return true
			})

			if conversions == 0 {
				continue
			}
			// One generated call, converted once, is the shape every wrapper has and
			// the shape the table describes: one response type per name.
			if conversions != 1 || len(responses) != 1 {
				t.Fatalf("%s makes %d generated calls and converts %d responses; list one of each per wrapper", name, len(responses), conversions)
			}

			sites[name] = problemSite{platform: platform, response: responses[0]}
		}
	}

	return sites
}

// receiverTypeName is the type a method is declared on, without the pointer.
func receiverTypeName(t *testing.T, method *ast.FuncDecl) string {
	t.Helper()

	receiver := method.Recv.List[0].Type
	if star, ok := receiver.(*ast.StarExpr); ok {
		receiver = star.X
	}

	ident, ok := receiver.(*ast.Ident)
	if !ok {
		t.Fatalf("%s: receiver %T is not a named type", method.Name.Name, receiver)
	}

	return ident.Name
}
