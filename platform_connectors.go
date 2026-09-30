package sdk

import (
	"context"
	"fmt"
)

// Register records a connector's manifest with this deployment, or updates the manifest
// already recorded under the same name.
//
// It is an upsert on the name, which is what makes it safe to call on every start of the
// process that hosts the connector: a redeploy re-declares the same connector rather than
// adding a second one, and a manifest change ships with the build that made it.
//
// Registering is deployment-wide and says nothing about who may use the connector. A
// caller needs a platform credential (`ksm_...`) holding write:organizations; an
// organization then activates the connector for itself over the Core API, and
// RegisterConnectorInput.EntitlementSlug is what decides whether its licence lets it.
//
// Kaiten's own built-in connectors do not come through here. They are compiled into the
// API binary and register in-process at startup, where there is no wire to be on and no
// credential to present; this endpoint is for a connector hosted somewhere else.
func (s *PlatformConnectors) Register(ctx context.Context, input RegisterConnectorInput) (Connector, error) {
	body, err := jsonBody(input.payload())
	if err != nil {
		var zero Connector
		return zero, fmt.Errorf("encode register connector request: %w", err)
	}

	resp, err := s.client.raw.RegisterConnectorWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero Connector
		return zero, fmt.Errorf("register connector: %w", err)
	}

	return expectJSON(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}
