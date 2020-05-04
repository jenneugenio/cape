package graph

import (
	"context"
	"io/ioutil"

	"github.com/markbates/pkger"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func queryIdentity(ctx context.Context, db database.Querier, email primitives.Email) (primitives.Identity, error) {
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
func createSystemRoles(ctx context.Context, db database.Querier) error {
	entities := make([]database.Entity, len(primitives.SystemRoles))
	for i, roleLabel := range primitives.SystemRoles {
		role, err := primitives.NewRole(roleLabel, true)
		if err != nil {
			return err
		}

		entities[i] = role
	}

	return db.Create(ctx, entities...)
}

// getRolesByLabel is a helper to retrieve a specific role from the database. This is
// useful for getting a system role from the database.
func getRolesByLabel(ctx context.Context, db database.Querier, labels []primitives.Label) ([]*primitives.Role, error) {
	in := make(database.In, len(labels))
	for i, label := range labels {
		in[i] = label
	}

	f := database.NewFilter(database.Where{"label": in}, nil, nil)
	var roles []*primitives.Role

	err := db.Query(ctx, &roles, f)
	if err != nil {
		return nil, err
	}

	if len(labels) != len(roles) {
		return nil, errors.New(NotFoundCause, "Could not find a role")
	}

	return roles, nil
}

// createAssignments is a helper function that makes it easy to assign roles to
// a given identity.
func createAssignments(ctx context.Context, db database.Querier,
	identity primitives.Identity, roles []*primitives.Role) error {
	assignments := make([]database.Entity, len(roles))
	for i, role := range roles {
		assignment, err := primitives.NewAssignment(identity.GetID(), role.ID)
		if err != nil {
			return err
		}

		assignments[i] = assignment
	}

	return db.Create(ctx, assignments...)
}

func attachDefaultPolicy(ctx context.Context, db database.Querier) error {
	adminPolicy, err := loadPolicyFile(primitives.DefaultAdminPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	globalPolicy, err := loadPolicyFile(primitives.DefaultGlobalPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	dcPolicy, err := loadPolicyFile(primitives.DefaultDataConnectorPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	err = db.Create(ctx, adminPolicy, globalPolicy, dcPolicy)
	if err != nil {
		return err
	}

	adminAttachment, err := createAttachment(ctx, db, adminPolicy.ID, primitives.AdminRole)
	if err != nil {
		return err
	}

	globalAttachment, err := createAttachment(ctx, db, globalPolicy.ID, primitives.GlobalRole)
	if err != nil {
		return err
	}

	dcAttachment, err := createAttachment(ctx, db, dcPolicy.ID, primitives.DataConnectorRole)
	if err != nil {
		return err
	}

	err = db.Create(ctx, adminAttachment, globalAttachment, dcAttachment)
	if err != nil {
		return err
	}

	return nil
}

func createAttachment(ctx context.Context, db database.Querier, policyID database.ID,
	roleLabel primitives.Label) (*primitives.Attachment, error) {
	roles, err := getRolesByLabel(ctx, db, []primitives.Label{roleLabel})
	if err != nil {
		return nil, err
	}

	attachment, err := primitives.NewAttachment(policyID, roles[0].ID)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func loadPolicyFile(file string) (*primitives.Policy, error) {
	dir := pkger.Dir("/primitives/policies/default")

	f, err := dir.Open(file)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return primitives.ParsePolicy(b)
}
