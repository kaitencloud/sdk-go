// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Delete soft-deletes the user identified by id. Requires the caller's credential to hold
// the delete:users scope -- a privileged, cross-organization operation, not implied by
// write:users.
func (s *PlatformUsers) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	resp, err := s.client.raw.DeleteUserWithResponse(ctx, openapi_types.UUID(parsedID))
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return expectNoContent(newPlatformAPIResponse(resp, resp.Body, resp.ApplicationproblemJSON401, resp.ApplicationproblemJSON403, resp.ApplicationproblemJSON404, resp.ApplicationproblemJSON422, resp.ApplicationproblemJSON500))
}
