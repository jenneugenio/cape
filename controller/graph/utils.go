package graph

import (
	"context"

	"github.com/dropoutlabs/cape/controller/graph/model"
	"github.com/dropoutlabs/cape/database"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
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
func createSystemRoles(ctx context.Context, tx database.Transaction, roleLabels []primitives.Label) ([]*primitives.Role, error) {
	// this is silly, need to create two arrays, one to be able to easily
	// add to the database and another to return so the roles can be used
	// later!!
	entities := make([]database.Entity, len(roleLabels))
	roles := make([]*primitives.Role, len(roleLabels))

	for i, roleLabel := range roleLabels {
		role, err := primitives.NewRole(roleLabel, true)
		if err != nil {
			return nil, err
		}

		entities[i] = role
		roles[i] = role
	}

	err := tx.Create(ctx, entities...)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// createAssignmentsForUser is a helper function intended to be used by the Setup
// graphql route. Its used in tandem with the above function create a
func createAssignmentsForUser(ctx context.Context, tx database.Transaction,
	user *primitives.User, roles []*primitives.Role) error {
	assignments := make([]database.Entity, len(roles))

	for i, role := range roles {
		assignment, err := primitives.NewAssignment(user.ID, role.ID)
		if err != nil {
			return err
		}

		assignments[i] = assignment
	}

	return tx.Create(ctx, assignments...)
}

func createDataConnector(ctx context.Context, db database.Backend, input model.CreateServiceRequest) (*primitives.Service, error) {
	creds := &primitives.Credentials{
		PublicKey: &input.PublicKey,
		Salt:      &input.Salt,
		Alg:       input.Alg,
	}

	role := &primitives.Role{}
	err := db.QueryOne(ctx, role, database.NewFilter(database.Where{"label": primitives.DataConnectorRole}, nil, nil))
	if err != nil {
		return nil, err
	}

	service, err := primitives.NewService(input.Email, input.Type, input.Endpoint, creds)
	if err != nil {
		return nil, err
	}

	assignment, err := primitives.NewAssignment(service.ID, role.ID)
	if err != nil {
		return nil, err
	}

	tx, err := db.Transaction(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx) // nolint: errcheck

	err = tx.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	err = tx.Create(ctx, assignment)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return service, nil
}
