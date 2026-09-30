package sdk

import (
	"errors"
	"net/http"
	"testing"
)

// TestLicenseFamiliesListWalksThePagedEnvelope pins the half of the family
// list the generated client cannot express: /license-families answers a paged
// envelope, while internal/gen types every list response as a bare array.
// Without the hand-written walk this call fails to decode, and with a walk that
// stops at the first page it silently returns a truncated catalogue.
func TestLicenseFamiliesListWalksThePagedEnvelope(t *testing.T) {
	t.Parallel()

	first := `{"items":[{"slug":"starter","versionCount":1}],"nextCursor":"cursor-2","hasMore":true}`
	second := `{"items":[{"slug":"premium","versionCount":3}],"hasMore":false}`

	client, sent := givenClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("cursor") == "cursor-2" {
			return respondJSON(http.StatusOK, second)(req)
		}
		return respondJSON(http.StatusOK, first)(req)
	})

	families, err := client.LicenseFamilies.List(t.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(families) != 2 {
		t.Fatalf("List() returned %d families, want 2 across both pages", len(families))
	}
	if families[0].Slug != "starter" || families[1].Slug != "premium" {
		t.Errorf("List() slugs = %q, %q, want starter, premium", families[0].Slug, families[1].Slug)
	}
	if sent.count() != 2 {
		t.Errorf("requests = %d, want 2: the second page has to be fetched", sent.count())
	}
	if got := sent.at(t, 0).path; got != "/api/license-families" {
		t.Errorf("path = %q, want /api/license-families", got)
	}
}

// TestLicenseFamiliesListKeepsADraftOnlyFamily is the reason CurrentVersion is
// a pointer. A family whose versions are all drafts or all archived is still
// part of the catalogue -- a console renders it as such -- so it must survive
// the walk rather than be filtered out for having nothing to resolve to.
func TestLicenseFamiliesListKeepsADraftOnlyFamily(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusOK,
		`{"items":[{"slug":"in-progress","versionCount":2}],"hasMore":false}`))

	families, err := client.LicenseFamilies.List(t.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(families) != 1 {
		t.Fatalf("List() returned %d families, want 1", len(families))
	}
	if families[0].CurrentVersion != nil {
		t.Errorf("CurrentVersion = %v, want nil for a family with nothing published", families[0].CurrentVersion)
	}
	if families[0].VersionCount != 2 {
		t.Errorf("VersionCount = %d, want 2: it counts the history, not what is for sale", families[0].VersionCount)
	}
}

// TestLicenseFamiliesGetResolvesTheCurrentVersionWithoutOptions covers the
// call the family address exists for: no parameters, and the answer is
// whichever version the family serves today.
func TestLicenseFamiliesGetResolvesTheCurrentVersionWithoutOptions(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK,
		`{"slug":"premium","versionCount":3,"currentVersion":{"id":"l2","slug":"premium-v2","version":"2","lifecycleState":"PUBLISHED"}}`))

	family, err := client.LicenseFamilies.Get(t.Context(), "premium", nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if family.CurrentVersion == nil {
		t.Fatal("CurrentVersion is nil, want the published version")
	}
	if family.CurrentVersion.Version == nil || *family.CurrentVersion.Version != "2" {
		t.Errorf("CurrentVersion.Version = %v, want 2", family.CurrentVersion.Version)
	}

	request := sent.only(t)
	if request.path != "/api/license-families/premium" {
		t.Errorf("path = %q, want /api/license-families/premium", request.path)
	}
	// Neither parameter may be sent when the caller asked for neither: `version`
	// alone changes the answer from "what is current" to "this exact row".
	for _, parameter := range []string{"version", "include"} {
		if request.query.Has(parameter) {
			t.Errorf("query carries %q = %q, want it absent", parameter, request.query.Get(parameter))
		}
	}
}

// TestLicenseFamiliesGetSendsTheOptionsItWasGiven pins both parameters onto
// the wire. They are the two ways a caller opts out of resolution, and an
// option dropped between the struct and the query silently answers a different
// question.
func TestLicenseFamiliesGetSendsTheOptionsItWasGiven(t *testing.T) {
	t.Parallel()

	client, sent := givenClient(t, respondJSON(http.StatusOK,
		`{"slug":"premium","versionCount":3,"currentVersion":{"id":"l1","slug":"premium-v1","version":"1","lifecycleState":"ARCHIVED"},"versions":[{"id":"l1","version":"1"}]}`))

	version := int32(1)
	family, err := client.LicenseFamilies.Get(t.Context(), "premium", &LicenseFamilyOptions{
		Version:         &version,
		IncludeVersions: true,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	request := sent.only(t)
	if got := request.query.Get("version"); got != "1" {
		t.Errorf("version = %q, want 1", got)
	}
	if got := request.query.Get("include"); got != includeVersions {
		t.Errorf("include = %q, want %q", got, includeVersions)
	}

	if family.Versions == nil {
		t.Fatal("Versions is nil, want the history include=versions asked for")
	}
	if family.CurrentVersion == nil || family.CurrentVersion.LifecycleState == nil {
		t.Fatal("CurrentVersion or its lifecycle state is nil")
	}
	if *family.CurrentVersion.LifecycleState != LicenseLifecycleArchived {
		t.Errorf("lifecycleState = %v, want ARCHIVED: naming a version reaches it whatever its state",
			*family.CurrentVersion.LifecycleState)
	}
}

// TestLicenseFamiliesGetReportsNoPublishedVersionByCode keeps the two 404s
// distinguishable. "This family has nothing to sell" and "no such family" share
// a status, so a caller that branches on the status alone cannot tell a
// draft-only product from a typo -- the code is what separates them.
func TestLicenseFamiliesGetReportsNoPublishedVersionByCode(t *testing.T) {
	t.Parallel()

	client, _ := givenClient(t, respondJSON(http.StatusNotFound,
		`{"title":"Not Found","status":404,"code":"GetLicenseFamily.NoPublishedVersion","detail":"License family \"in-progress\" has no published version (2 version(s), all draft or archived)"}`))

	_, err := client.LicenseFamilies.Get(t.Context(), "in-progress", nil)
	if err == nil {
		t.Fatal("Get() error = nil, want a not-found error")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(err, *Error) = false for %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
	if apiErr.Code != "GetLicenseFamily.NoPublishedVersion" {
		t.Errorf("Code = %q, want GetLicenseFamily.NoPublishedVersion", apiErr.Code)
	}
}
