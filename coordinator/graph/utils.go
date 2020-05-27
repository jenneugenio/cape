package graph

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/manifoldco/go-base64"
	"github.com/markbates/pkger"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func getCredentialProvider(ctx context.Context, q database.Querier, input model.SessionRequest) (primitives.CredentialProvider, error) {
	if input.Email != nil {
		filter := database.NewFilter(database.Where{"email": input.Email.String()}, nil, nil)
		user := &primitives.User{}
		err := q.QueryOne(ctx, user, filter)
		return user, err
	}

	token := &primitives.Token{}
	err := q.Get(ctx, *input.TokenID, token)
	if err != nil {
		return nil, err
	}

	return token, nil
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

	wPolicy, err := loadPolicyFile(primitives.DefaultWorkerPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	err = db.Create(ctx, adminPolicy, globalPolicy, dcPolicy, wPolicy)
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

	wAttachment, err := createAttachment(ctx, db, wPolicy.ID, primitives.WorkerRole)
	if err != nil {
		return err
	}

	err = db.Create(ctx, adminAttachment, globalAttachment, dcAttachment, wAttachment)
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
	f, err := pkger.Open("github.com/capeprivacy/cape:/primitives/policies/default/" + file)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return primitives.ParsePolicy(b)
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

func createConfig(rootKey [32]byte) (*primitives.Config, *crypto.KeyURL, *auth.Keypair, error) {
	encryptionKey, err := crypto.NewBase64KeyURL(nil)
	if err != nil {
		return nil, nil, nil, err
	}

	encryptedKey, err := crypto.Encrypt(rootKey, []byte(encryptionKey.String()))
	if err != nil {
		return nil, nil, nil, err
	}

	keypair, err := auth.NewKeypair()
	if err != nil {
		return nil, nil, nil, err
	}

	by, err := json.Marshal(keypair.Package())
	if err != nil {
		return nil, nil, nil, err
	}

	encryptedAuth, err := crypto.Encrypt(rootKey, by)
	if err != nil {
		return nil, nil, nil, err
	}

	config, err := primitives.NewConfig(base64.New(encryptedKey), base64.New(encryptedAuth))
	if err != nil {
		return nil, nil, nil, err
	}

	return config, encryptionKey, keypair, nil
}
