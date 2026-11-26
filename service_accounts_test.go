// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"errors"
	"net/http"
	"testing"
)

// Delete is an authenticated, bodiless DELETE on the service account itself, and its
// 204 is success.
func TestServiceAccountsDeleteSendsABodilessDelete(t *testing.T) {
	t.Parallel()

	client, rec := givenClient(t, respondNoContent(), WithBearerToken("ksh_secret"))

	if err := client.ServiceAccounts.Delete(t.Context(), "ci-deploy"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	request := rec.only(t)
	if request.method != http.MethodDelete || request.path != "/api/service-accounts/ci-deploy" {
		t.Errorf("request = %s %s, want DELETE /api/service-accounts/ci-deploy", request.method, request.path)
	}
	if len(request.body) != 0 {
		t.Errorf("request body = %s, want none", request.body)
	}
	if got := request.header.Get("Authorization"); got != "Bearer ksh_secret" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer ksh_secret")
	}
}

// An unknown slug, an account already deleted or another organization's, is not there
// to delete, and the refusal carries the code every service-account operation answers
// with.
func TestServiceAccountsDeleteReportsNotFound(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusNotFound, problem(http.StatusNotFound, "ServiceAccount.NotFound")))

	err := client.ServiceAccounts.Delete(t.Context(), "ci-deploy")

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want an *Error", err)
	}
	if apiErr.StatusCode != http.StatusNotFound || apiErr.Code != "ServiceAccount.NotFound" || apiErr.Problem == nil {
		t.Errorf("StatusCode = %d, Code = %q, Problem = %v, want %d, %q", apiErr.StatusCode, apiErr.Code, apiErr.Problem, http.StatusNotFound, "ServiceAccount.NotFound")
	}
}
