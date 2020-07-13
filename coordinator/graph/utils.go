package graph

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/primitives"
)

// buildAttachment takes a primitives attachment and builds at graphql
// model representation of it
func buildAttachment(ctx context.Context, enforcer *auth.Enforcer, db db.Interface,
	attachment *primitives.Attachment) (*model.Attachment, error) {
	role := &primitives.Role{}
	err := enforcer.Get(ctx, attachment.RoleID, role)
	if err != nil {
		return nil, err
	}

	policy, err := db.Policies().GetByID(ctx, attachment.PolicyID)
	if err != nil {
		return nil, err
	}

	return &model.Attachment{
		ID:        attachment.ID,
		CreatedAt: attachment.CreatedAt,
		UpdatedAt: attachment.UpdatedAt,
		Role:      role,
		Policy:    policy,
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
