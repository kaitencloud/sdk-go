package sdk

import (
	"context"
	"fmt"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// List returns all license families, each with the version it currently
// resolves to.
//
// A family is the product; the licenses Licenses.List returns are its versions.
// What makes this the useful list of the two is that a family's slug survives a
// rename and a new version, so it is the identifier worth storing -- a version
// slug names the row a customer bought, and changes the next time the vendor
// publishes.
//
// A family with nothing published (every version still a draft, or all of them
// archived) is included, with CurrentVersion nil.
func (s *LicenseFamilies) List(ctx context.Context) ([]LicenseFamilyView, error) {
	return listAll[LicenseFamilyView](ctx, s.client.list, listPath("license-families"), nil)
}

// Get resolves the family identified by familySlug to one version.
//
// With no options it answers the version the family currently serves: its
// default version, or its highest-numbered published one. That is the point of
// addressing a family -- publishing a new version changes what this returns,
// with no change here.
//
// options.Version asks for one version by number instead, whatever its
// lifecycle state, archived included: pinned access stays pinned.
// options.IncludeVersions fills Versions with the family's whole history.
//
// A family that exists but has no published version is an error, not an empty
// answer: the code is GetLicenseFamily.NoPublishedVersion, on a 404. Callers
// telling that apart from an unknown family should branch on
// Error.Problem.Code rather than on the status.
func (s *LicenseFamilies) Get(ctx context.Context, familySlug string, options *LicenseFamilyOptions) (LicenseFamilyView, error) {
	params := &gen.GetLicenseFamilyParams{}
	if options != nil {
		params.Version = options.Version
		if options.IncludeVersions {
			// include is a list on the wire, comma-separated, so that the API can
			// add representations without changing the parameter's type.
			include := []gen.GetLicenseFamilyParamsInclude{includeVersions}
			params.Include = &include
		}
	}

	resp, err := s.client.raw.GetLicenseFamilyWithResponse(ctx, familySlug, params)
	if err != nil {
		var zero LicenseFamilyView
		return zero, fmt.Errorf("get license family: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}
