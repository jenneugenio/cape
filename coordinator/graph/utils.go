package graph

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	modelmigration "github.com/capeprivacy/cape/models/migration"
	"github.com/capeprivacy/cape/primitives"
)

// buildAttachment takes a primitives attachment and builds at graphql
// model representation of it
func buildAttachment(ctx context.Context, db *auth.Enforcer,
	attachment *primitives.Attachment) (*model.Attachment, error) {
	role := &primitives.Role{}
	err := db.Get(ctx, attachment.RoleID, role)
	if err != nil {
		return nil, err
	}

	policy := &primitives.Policy{}
	err = db.Get(ctx, attachment.PolicyID, policy)
	if err != nil {
		return nil, err
	}

	modelPolicy := modelmigration.PolicyFromPrimitive(policy)

	return &model.Attachment{
		ID:        attachment.ID,
		CreatedAt: attachment.CreatedAt,
		UpdatedAt: attachment.UpdatedAt,
		Role:      role,
		Policy:    &modelPolicy,
	}, nil
}

func hasRole(roles []*primitives.Role, label primitives.Label) bool {
	found := false
	for _, role := range roles {
		if role.Label == label {
			found = true
			break
		}
	}

	return found
}
