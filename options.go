// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"net/http"
	"strings"

	"github.com/kaitencloud/sdk-go/internal/gen"
	"github.com/kaitencloud/sdk-go/internal/genplatform"
)

// RequestEditorFn edits an outgoing request before it is sent.
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Option configures a Client or a PlatformClient.
//
// Declared here rather than aliased to a generated option type, because there are two
// generated clients and their option types are not interchangeable: each is a function
// over its own package's *Client. An Option collected once and passed to either
// constructor -- which is what a caller holding a token and a pair of base URLs wants
// to do -- is only possible if the type is the SDK's own.
//
// Options are applied in order, so a later one overrides an earlier one.
type Option func(*clientOptions)

// clientOptions accumulates what an Option sets, for a constructor to translate into
// whichever generated package it is building.
type clientOptions struct {
	baseURL    string
	httpClient *http.Client
	editors    []RequestEditorFn
}

// WithHTTPClient sets the *http.Client used to send requests. It replaces the client
// the constructor would otherwise build, and with it DefaultTimeout.
func WithHTTPClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithRequestEditorFn adds a function that edits each outgoing request before it is sent.
func WithRequestEditorFn(fn RequestEditorFn) Option {
	return func(o *clientOptions) {
		o.editors = append(o.editors, fn)
	}
}

// WithBaseURL overrides the client's base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

// WithBearerToken sets the Authorization header to a bearer token on every request.
//
// Which token class is accepted depends on which client it is given to: the Core API
// (Client) takes an organization-scoped token, the Platform API (PlatformClient) takes
// a platform credential, and each listener refuses the other's.
func WithBearerToken(token string) Option {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		if strings.TrimSpace(token) != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		return nil
	})
}

// resolveOptions applies opts over the defaults every client shares.
func resolveOptions(baseURL string, opts []Option) clientOptions {
	resolved := clientOptions{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: DefaultTimeout},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&resolved)
		}
	}

	return resolved
}

func (o clientOptions) coreOptions() []gen.ClientOption {
	translated := make([]gen.ClientOption, 0, len(o.editors)+1)
	translated = append(translated, gen.WithHTTPClient(o.httpClient))

	for _, editor := range o.editors {
		translated = append(translated, gen.WithRequestEditorFn(gen.RequestEditorFn(editor)))
	}

	return translated
}

func (o clientOptions) platformOptions() []genplatform.ClientOption {
	translated := make([]genplatform.ClientOption, 0, len(o.editors)+1)
	translated = append(translated, genplatform.WithHTTPClient(o.httpClient))

	for _, editor := range o.editors {
		translated = append(translated, genplatform.WithRequestEditorFn(genplatform.RequestEditorFn(editor)))
	}

	return translated
}

// listTransport hands the resolved plumbing to the paginated list methods. The
// generated client types each endpoint's page as its own struct, so one generic
// walk cannot be written over them; the list methods build their own requests
// and need the same three things the generated client was built from.
func (o clientOptions) listTransport() listTransport {
	return listTransport{
		baseURL: o.baseURL,
		doer:    o.httpClient,
		editors: o.editors,
	}
}
