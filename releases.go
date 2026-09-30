package sdk

import (
	"context"
	"fmt"
)

// List returns all releases.
func (s *Releases) List(ctx context.Context) ([]Release, error) {
	return listAll[Release](ctx, s.client.list, listPath("releases"), nil)
}

// Get returns the release identified by releaseSlug.
func (s *Releases) Get(ctx context.Context, releaseSlug string) (Release, error) {
	resp, err := s.client.raw.GetReleaseBySlugWithResponse(ctx, releaseSlug)
	if err != nil {
		var zero Release
		return zero, fmt.Errorf("get release: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Create creates a new release.
func (s *Releases) Create(ctx context.Context, input ReleaseInput) (Release, error) {
	body, err := jsonBody(input.payload())
	if err != nil {
		var zero Release
		return zero, fmt.Errorf("encode release request: %w", err)
	}

	resp, err := s.client.raw.CreateReleaseWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Release
		return zero, fmt.Errorf("create release: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Delete deletes the release identified by releaseSlug.
func (s *Releases) Delete(ctx context.Context, releaseSlug string) error {
	resp, err := s.client.raw.DeleteReleaseWithResponse(ctx, releaseSlug)
	if err != nil {
		return fmt.Errorf("delete release: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
