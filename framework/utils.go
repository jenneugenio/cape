package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	errors "github.com/capeprivacy/cape/partyerrors"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/primitives"
	"github.com/go-openapi/runtime/middleware/header"
	"github.com/manifoldco/go-base64"
	"github.com/markbates/pkger"

	pkgerrors "errors"
)

// QueryUserPolicies is a helper function to query all roles assigned to an user and then
// all policies attached to those roles.
func QueryUserPolicies(ctx context.Context, db database.Querier, userID string) ([]*primitives.Policy, error) {
	var assignments []*primitives.Assignment
	assignmentFilter := database.NewFilter(database.Where{"user_id": userID}, nil, nil)
	err := db.Query(ctx, &assignments, assignmentFilter)
	if err != nil {
		return nil, err
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	if len(roleIDs) == 0 {
		return []*primitives.Policy{}, nil
	}

	var attachments []*primitives.Attachment
	attachmentFilter := database.NewFilter(database.Where{"role_id": roleIDs}, nil, nil)
	err = db.Query(ctx, &attachments, attachmentFilter)
	if err != nil {
		return nil, err
	}

	policyIDs := database.InFromEntities(attachments, func(e interface{}) interface{} {
		return e.(*primitives.Attachment).PolicyID
	})

	if len(policyIDs) == 0 {
		return []*primitives.Policy{}, nil
	}

	var policies []*primitives.Policy
	err = db.Query(ctx, &policies, database.NewFilter(database.Where{"id": policyIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func QueryRoles(ctx context.Context, db database.Querier, userID string) ([]*primitives.Role, error) {
	var assignments []*primitives.Assignment
	filter := database.NewFilter(database.Where{
		"user_id": userID,
	}, nil, nil)
	err := db.Query(ctx, &assignments, filter)
	if err != nil {
		return nil, err
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	var roles []*primitives.Role
	err = db.Query(ctx, &roles, database.NewFilter(database.Where{
		"id": roleIDs,
	}, nil, nil))
	if err != nil {
		return nil, err
	}

	return roles, nil
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case pkgerrors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case pkgerrors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case pkgerrors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case pkgerrors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func getCredentialProvider(ctx context.Context, q database.Querier, capedb db.Interface, input LoginRequest) (primitives.CredentialProvider, error) {
	if input.Email != nil {
		return capedb.Users().Get(ctx, *input.Email)
	}

	token := &primitives.Token{}
	err := q.Get(ctx, *input.TokenID, token)
	if err != nil {
		return nil, err
	}

	return token, nil
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

func attachDefaultPolicy(ctx context.Context, db database.Querier) error {
	adminPolicy, err := loadPolicyFile(primitives.DefaultAdminPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	globalPolicy, err := loadPolicyFile(primitives.DefaultGlobalPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	err = db.Create(ctx, adminPolicy, globalPolicy)
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

	err = db.Create(ctx, adminAttachment, globalAttachment)
	if err != nil {
		return err
	}

	return nil
}

func createAttachment(ctx context.Context, db database.Querier, policyID database.ID,
	roleLabel primitives.Label) (*primitives.Attachment, error) {
	roles, err := GetRolesByLabel(ctx, db, []primitives.Label{roleLabel})
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

// GetRolesByLabel is a helper to retrieve a specific role from the database. This is
// useful for getting a system role from the database.
func GetRolesByLabel(ctx context.Context, db database.Querier, labels []primitives.Label) ([]*primitives.Role, error) {
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

// CreateAssignments is a helper function that makes it easy to assign roles to
// a given user.
func CreateAssignments(ctx context.Context, db database.Querier,
	userID string, roles []*primitives.Role) error {
	assignments := make([]database.Entity, len(roles))
	for i, role := range roles {
		assignment, err := primitives.NewAssignment(userID, role.ID)
		if err != nil {
			return err
		}

		assignments[i] = assignment
	}

	return db.Create(ctx, assignments...)
}
