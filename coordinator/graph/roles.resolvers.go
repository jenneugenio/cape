package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) AssignRole(ctx context.Context, input model.AssignRoleRequest) (*model.Assignment, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	assignment, err := primitives.NewAssignment(input.UserID, input.RoleID)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, assignment)
	if err != nil {
		return nil, err
	}

	role := &primitives.Role{}
	err = enforcer.Get(ctx, input.RoleID, role)
	if err != nil {
		return nil, err
	}

	user, err := r.Database.Users().GetByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &model.Assignment{
		ID:        assignment.ID,
		CreatedAt: assignment.CreatedAt,
		UpdatedAt: assignment.UpdatedAt,
		Role:      role,
		User:      user,
	}, nil
}

func (r *mutationResolver) UnassignRole(ctx context.Context, input model.AssignRoleRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	assignment := &primitives.Assignment{}

	filter := database.NewFilter(database.Where{
		"role_id": input.RoleID.String(),
		"user_id": input.UserID,
	}, nil, nil)

	err := enforcer.QueryOne(ctx, assignment, filter)
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.AssignmentType, assignment.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) Role(ctx context.Context, id string) (*models.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var primitive primitives.Role
	err := enforcer.Get(ctx, id, &primitive)
	if err != nil {
		return nil, err
	}

	return &primitive, nil
}

func (r *queryResolver) RoleByLabel(ctx context.Context, label models.Label) (*models.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	role := &primitives.Role{}
	err := enforcer.QueryOne(ctx, role, database.NewFilter(database.Where{"label": label}, nil, nil))
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *queryResolver) Roles(ctx context.Context) ([]*models.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var roles []*primitives.Role
	err := enforcer.Query(ctx, &roles, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *queryResolver) RoleMembers(ctx context.Context, roleID string) ([]*models.User, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var assignments []*primitives.Assignment
	err := enforcer.Query(ctx, &assignments, database.NewFilter(database.Where{"role_id": roleID.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	var userIDs []string
	for _, assignment := range assignments {
		userIDs = append(userIDs, assignment.UserID)
	}

	var users []models.User
	if len(userIDs) > 0 {
		opts := &db.ListUserOptions{Options: nil, FilterIDs: userIDs}
		u, err := r.Database.Users().List(ctx, opts)
		if err != nil {
			return nil, err
		}
		users = u
	}

	userPtrs := make([]*models.User, len(users))
	for i, user := range users {
		u := user
		userPtrs[i] = &u
	}

	return userPtrs, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) DeleteRole(ctx context.Context, input model.DeleteRoleRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	role := &primitives.Role{}
	err := enforcer.Get(ctx, input.ID, role)
	if err != nil {
		return nil, err
	}

	if role.System {
		return nil, errs.New(CannotDeleteSystemRole, "Role %s is a system role. Cannot delete", role.Label)
	}

	err = enforcer.Delete(ctx, primitives.RoleType, input.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (r *modelRoleResolver) ID(ctx context.Context, obj *models.Role) (database.ID, error) {
	return database.DecodeFromString(obj.ID)
}
func (r *roleResolver) ID(ctx context.Context, obj *models.Role) (database.ID, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *roleResolver) Label(ctx context.Context, obj *models.Role) (primitives.Label, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *Resolver) ModelRole() generated.ModelRoleResolver { return &modelRoleResolver{r} }
func (r *Resolver) Role() generated.RoleResolver           { return &roleResolver{r} }

type modelRoleResolver struct{ *Resolver }
type roleResolver struct{ *Resolver }
