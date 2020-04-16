package graph

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func queryIdentity(ctx context.Context, db database.Backend, email primitives.Email) (primitives.Identity, error) {
	filter := database.NewFilter(database.Where{"email": email.String()}, nil, nil)

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

// createSystemRoles is a helper intended to be used by the Setup graphql route.
// It creates all the roles given by the list of role labels and makes sure
// they are system roles
func createSystemRoles(ctx context.Context, tx database.Transaction) error {
	entities := make([]database.Entity, len(primitives.SystemRoles))
	for i, roleLabel := range primitives.SystemRoles {
		role, err := primitives.NewRole(roleLabel, true)
		if err != nil {
			return err
		}

		entities[i] = role
	}

	return tx.Create(ctx, entities...)
}

// getRolesByLabel is a helper to retrieve a specific role from the database. This is
// useful for getting a system role from the database.
func getRolesByLabel(ctx context.Context, tx database.Transaction, labels []primitives.Label) ([]*primitives.Role, error) {
	in := make(database.In, len(labels))
	for i, label := range labels {
		in[i] = label
	}

	f := database.NewFilter(database.Where{"label": in}, nil, nil)
	roles := []*primitives.Role{}

	err := tx.Query(ctx, &roles, f)
	if err != nil {
		return nil, err
	}

	if len(labels) != len(roles) {
		return nil, errors.New(NotFoundCause, "Could not find a role")
	}

	return roles, nil
}

// getRoles is a helper that returns all of the roles assigned to a given identity.
func getRoles(ctx context.Context, db database.Backend, identityID database.ID) ([]*primitives.Role, error) {
	assignments := []*primitives.Assignment{}
	filter := database.NewFilter(database.Where{
		"identity_id": identityID,
	}, nil, nil)
	err := db.Query(ctx, &assignments, filter)
	if err != nil {
		return nil, err
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	roles := []*primitives.Role{}
	err = db.Query(ctx, &roles, database.NewFilter(database.Where{
		"id": roleIDs,
	}, nil, nil))
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// createAssignments is a helper function that makes it easy to assign roles to
// a given identity.
func createAssignments(ctx context.Context, tx database.Transaction,
	identity primitives.Identity, roles []*primitives.Role) error {
	assignments := make([]database.Entity, len(roles))
	for i, role := range roles {
		assignment, err := primitives.NewAssignment(identity.GetID(), role.ID)
		if err != nil {
			return err
		}

		assignments[i] = assignment
	}

	return tx.Create(ctx, assignments...)
}
