package sdk

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// List returns all metadata fields for the given resourceType.
func (s *MetadataFields) List(ctx context.Context, resourceType MetadataFieldResourceType) ([]MetadataField, error) {
	return listAll[MetadataField](ctx, s.client.list, listPath("metadata-fields"), url.Values{
		"resourceType": []string{string(resourceType)},
	})
}

// Create creates a new metadata field.
func (s *MetadataFields) Create(ctx context.Context, input MetadataFieldInput) (MetadataField, error) {
	body, err := jsonBody(input.payload())
	if err != nil {
		var zero MetadataField
		return zero, fmt.Errorf("encode metadata field request: %w", err)
	}

	resp, err := s.client.raw.CreateMetadataFieldWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		var zero MetadataField
		return zero, fmt.Errorf("create metadata field: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500, resp.ApplicationproblemJSON503), resp.JSON201)
}

// Update updates the metadata field identified by id.
func (s *MetadataFields) Update(ctx context.Context, id string, input MetadataFieldUpdateInput) (MetadataField, error) {
	body, err := jsonBody(input.payload())
	if err != nil {
		var zero MetadataField
		return zero, fmt.Errorf("encode metadata field request: %w", err)
	}

	resp, err := s.client.raw.UpdateMetadataFieldWithBodyWithResponse(ctx, id, contentTypeJSON, body)
	if err != nil {
		var zero MetadataField
		return zero, fmt.Errorf("update metadata field: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Archive archives the metadata field identified by id.
func (s *MetadataFields) Archive(ctx context.Context, id string) (MetadataField, error) {
	resp, err := s.client.raw.ArchiveMetadataFieldWithResponse(ctx, id)
	if err != nil {
		var zero MetadataField
		return zero, fmt.Errorf("archive metadata field: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Unarchive unarchives the metadata field identified by id.
func (s *MetadataFields) Unarchive(ctx context.Context, id string) (MetadataField, error) {
	resp, err := s.client.raw.UnarchiveMetadataFieldWithResponse(ctx, id)
	if err != nil {
		var zero MetadataField
		return zero, fmt.Errorf("unarchive metadata field: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON409, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// DryRun previews the impact of applying jsonSchema to the field identified by
// id: it returns how many existing resources carry a value under the field's
// key that would no longer satisfy the candidate schema, plus a small sample.
// It does not mutate the field.
func (s *MetadataFields) DryRun(ctx context.Context, id string, jsonSchema map[string]any) (MetadataFieldDryRunImpact, error) {
	body, err := jsonBody(gen.MetadataFieldSchemaProposal{JsonSchema: jsonSchema})
	if err != nil {
		var zero MetadataFieldDryRunImpact
		return zero, fmt.Errorf("encode metadata field dry-run request: %w", err)
	}

	resp, err := s.client.raw.DryRunMetadataFieldWithBodyWithResponse(ctx, id, contentTypeJSON, body)
	if err != nil {
		var zero MetadataFieldDryRunImpact
		return zero, fmt.Errorf("dry-run metadata field: %w", err)
	}

	return expectJSON(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500), resp.JSON200)
}

// Reorder sets the display order of metadata fields to the given ids.
func (s *MetadataFields) Reorder(ctx context.Context, ids []string) error {
	body, err := jsonBody(metadataFieldsReorderPayload{
		IDs: append([]string(nil), ids...),
	})
	if err != nil {
		return fmt.Errorf("encode metadata fields reorder request: %w", err)
	}

	resp, err := s.client.raw.ReorderMetadataFieldsWithBodyWithResponse(ctx, contentTypeJSON, body)
	if err != nil {
		return fmt.Errorf("reorder metadata fields: %w", err)
	}

	return expectNoContent(newAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON400, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
