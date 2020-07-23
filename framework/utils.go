package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/primitives"
	"github.com/go-openapi/runtime/middleware/header"
	"github.com/markbates/pkger"

	pkgerrors "errors"
)

// QueryUserPolicies is a helper function to query all roles assigned to an user and then
// all policies attached to those roles.
func QueryUserPolicies(ctx context.Context, querier database.Querier, capedb db.Interface, userID string) ([]*models.Policy, error) {
	var assignments []*primitives.Assignment
	assignmentFilter := database.NewFilter(database.Where{"user_id": userID}, nil, nil)
	err := querier.Query(ctx, &assignments, assignmentFilter)
	if err != nil {
		return nil, err
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	if len(roleIDs) == 0 {
		return nil, nil
	}

	var attachments []*primitives.Attachment
	attachmentFilter := database.NewFilter(database.Where{"role_id": roleIDs}, nil, nil)
	err = querier.Query(ctx, &attachments, attachmentFilter)
	if err != nil {
		return nil, err
	}

	var policyIDs []string
	for _, a := range attachments {
		policyIDs = append(policyIDs, a.PolicyID)
	}

	if len(policyIDs) == 0 {
		return []*models.Policy{}, nil
	}

	policies, err := capedb.Policies().List(ctx, &db.ListPolicyOptions{FilterIDs: policyIDs})
	if err != nil {
		return nil, err
	}

	policyPtrs := make([]*models.Policy, len(policies))
	for i, policy := range policies {
		p := policy
		policyPtrs[i] = &p
	}

	return policyPtrs, nil
}

// QueryUserRBAC is a helper function to query all roles assigned to an user and then
// all policies attached to those roles.
func QueryUserRBAC(ctx context.Context, querier database.Querier, capedb db.Interface, userID string) ([]*models.RBACPolicy, error) {
	var assignments []*primitives.Assignment
	assignmentFilter := database.NewFilter(database.Where{"user_id": userID}, nil, nil)
	err := querier.Query(ctx, &assignments, assignmentFilter)
	if err != nil {
		return nil, err
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	if len(roleIDs) == 0 {
		return nil, nil
	}

	var attachments []*primitives.Attachment
	attachmentFilter := database.NewFilter(database.Where{"role_id": roleIDs}, nil, nil)
	err = querier.Query(ctx, &attachments, attachmentFilter)
	if err != nil {
		return nil, err
	}

	var policyIDs []string
	for _, a := range attachments {
		policyIDs = append(policyIDs, a.PolicyID)
	}

	if len(policyIDs) == 0 {
		return []*models.RBACPolicy{}, nil
	}

	rbacs, err := capedb.RBAC().List(ctx, &db.ListRBACOptions{FilterIDs: policyIDs})
	if err != nil {
		return nil, err
	}

	rbacPtrs := make([]*models.RBACPolicy, len(rbacs))
	for i, policy := range rbacs {
		p := policy
		rbacPtrs[i] = &p
	}

	return rbacPtrs, nil
}

// TODO -- Don't think this function makes sense anymore (currently returns no roles)
func QueryRoles(ctx context.Context, db database.Querier, userID string) ([]*models.Role, error) {
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

	return nil, nil
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
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

// CreateSystemRoles is a helper intended to be used by the Setup graphql route.
// It creates all the roles given by the list of role labels and makes sure
// they are system roles
func CreateSystemRoles(ctx context.Context, db database.Querier) error {
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

func AttachDefaultPolicy(ctx context.Context, db database.Querier, capedb db.Interface) error {
	adminRBAC, err := loadRBACFile(primitives.DefaultAdminPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	globalRBAC, err := loadRBACFile(primitives.DefaultGlobalPolicy.String() + ".yaml")
	if err != nil {
		return err
	}

	err = capedb.RBAC().Create(ctx, *adminRBAC)
	if err != nil {
		return err
	}

	err = capedb.RBAC().Create(ctx, *globalRBAC)
	if err != nil {
		return err
	}

	adminAttachment, err := createAttachment(ctx, db, adminRBAC.ID, primitives.AdminRole)
	if err != nil {
		return err
	}

	globalAttachment, err := createAttachment(ctx, db, globalRBAC.ID, primitives.GlobalRole)
	if err != nil {
		return err
	}

	err = db.Create(ctx, adminAttachment, globalAttachment)
	if err != nil {
		return err
	}

	return nil
}

func createAttachment(ctx context.Context, db database.Querier, policyID string,
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

func loadRBACFile(file string) (*models.RBACPolicy, error) {
	f, err := pkger.Open("github.com/capeprivacy/cape:/primitives/policies/default/" + file)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return models.ParseRBACPolicy(b)
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
