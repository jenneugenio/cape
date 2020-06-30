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
	"github.com/capeprivacy/cape/coordinator/database2"
	"github.com/capeprivacy/cape/coordinator/graph/model"
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
func buildAttachment(ctx context.Context, enforcer *auth.Enforcer, db *database2.Database,
	attachment *primitives.Attachment) (*model.Attachment, error) {
	role, err := queryRole(ctx, db, attachment.RoleID.String())
	if err != nil {
		return nil, err
	}

	policy := &primitives.Policy{}
	err = enforcer.Get(ctx, attachment.PolicyID, policy)
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
func createSystemRoles(ctx context.Context, db *database2.Database) ([]*primitives.Role, error) {
	var roles []*primitives.Role
	for _, roleLabel := range primitives.SystemRoles {
		id := db.GetID()
		_, err := db.Pool.Exec(ctx, "INSERT INTO roles (id, label, system) VALUES ($1, $2, $3)", id.String(), roleLabel, true)
		if err != nil {
			return nil, err
		}

		role, err := primitives.NewRole(id.String(), roleLabel, true)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// createAssignments is a helper function that makes it easy to assign roles to
// a given identity.
func createAssignments(ctx context.Context, db database.Querier,
	identity primitives.Identity, roles []*primitives.Role) error {
	assignments := make([]database.Entity, len(roles))
	for i, role := range roles {
		assignment, err := primitives.NewAssignment(identity.GetID(), role.ID.String())
		if err != nil {
			return err
		}

		assignments[i] = assignment
	}

	return db.Create(ctx, assignments...)
}

func attachDefaultPolicy(ctx context.Context, enforcer database.Querier, db *database2.Database) error {
	adminPolicy, err := loadPolicyFile(primitives.DefaultAdminPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	globalPolicy, err := loadPolicyFile(primitives.DefaultGlobalPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	err = enforcer.Create(ctx, adminPolicy, globalPolicy)
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

	err = enforcer.Create(ctx, adminAttachment, globalAttachment)
	if err != nil {
		return err
	}

	return nil
}

func createAttachment(ctx context.Context, db *database2.Database, policyID database.ID,
	roleLabel primitives.Label) (*primitives.Attachment, error) {
	role, err := queryRoleByLabel(ctx, db, roleLabel)
	if err != nil {
		return nil, err
	}

	attachment, err := primitives.NewAttachment(policyID, role.ID.String())
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

func queryRole(ctx context.Context, db *database2.Database, id string) (*primitives.Role, error) {
	row := db.Pool.QueryRow(ctx, "SELECT label, system FROM roles WHERE id = $1", id)

	var label primitives.Label
	var system bool
	err := row.Scan(&label, &system)
	if err != nil {
		return nil, err
	}

	return primitives.NewRole(id, label, system)
}

func queryRoleByLabel(ctx context.Context, db *database2.Database, label primitives.Label) (*primitives.Role, error) {
	row := db.Pool.QueryRow(ctx, "SELECT id, system FROM roles WHERE label = $1", label)

	var id string
	var system bool
	err := row.Scan(&id, &system)
	if err != nil {
		return nil, err
	}

	return primitives.NewRole(id, label, system)
}
