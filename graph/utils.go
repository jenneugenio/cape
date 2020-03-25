package graph

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/model"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

func queryIdentity(ctx context.Context, db database.Backend, email string) (primitives.Identity, error) {
	filter := database.NewFilter(database.Where{"email": email}, nil, nil)

	user := &primitives.User{}
	err := db.QueryOne(ctx, user, filter)
	if err != nil && !errors.FromCause(err, database.NotFoundCause) {
		return nil, err
	}
	if err == nil {
		return user, err
	}

	service := &primitives.Service{}
	err = db.QueryOne(ctx, service, filter)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// buildAttachment takes a primitives attachment and builds at graphql
// model representation of it
func buildAttachment(ctx context.Context, db database.Backend,
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

	return &model.Attachment{
		ID:        attachment.ID,
		CreatedAt: attachment.CreatedAt,
		UpdatedAt: attachment.UpdatedAt,
		Role:      role,
		Policy:    policy,
	}, nil
}
