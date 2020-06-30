package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateRole(ctx context.Context, input model.CreateRoleRequest) (*primitives.Role, error) {
	fmt.Println("HELLOOOOOO??????")
	id := r.Database.GetID().String()
	role, err := primitives.NewRole(id, input.Label, false)
	if err != nil {
		return nil, err
	}

	_, err = r.Database.Pool.Exec(ctx, "INSERT INTO roles (id, label, system) VALUES ($1, $2, $3)", id, role.Label, false)
	if err != nil {
		fmt.Println("BOO", err)
		return nil, err
	}

	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if len(input.IdentityIds) > 0 {
		assignments := make([]database.Entity, len(input.IdentityIds))
		for i, id := range input.IdentityIds {
			assignment, err := primitives.NewAssignment(id, role.ID.String())
			if err != nil {
				return nil, err
			}
			assignments[i] = assignment
		}
		err = enforcer.Create(ctx, assignments...)
		if err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (r *mutationResolver) DeleteRole(ctx context.Context, input model.DeleteRoleRequest) (*string, error) {
	tag, err := r.Database.Pool.Exec(ctx, "DELETE FROM roles WHERE id = $1 AND system != true", input.ID)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() != 1 {
		return nil, errs.New(CannotDeleteSystemRole, "Role %s is a system role or doesn't exist. Cannot delete", input.ID)
	}

	return nil, nil
}

func (r *mutationResolver) AssignRole(ctx context.Context, input model.AssignRoleRequest) (*model.Assignment, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	assignment, err := primitives.NewAssignment(input.IdentityID, input.RoleID)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, assignment)
	if err != nil {
		return nil, err
	}

	row := r.Database.Pool.QueryRow(ctx, "SELECT label, system FROM roles WHERE id = $1", input.RoleID)

	var label primitives.Label
	var system bool
	err = row.Scan(&label, &system)
	if err != nil {
		return nil, err
	}

	role, err := primitives.NewRole(input.RoleID, label, system)
	if err != nil {
		return nil, err
	}

	typ, err := input.IdentityID.Type()
	if err != nil {
		return nil, err
	}

	var identity primitives.Identity
	if typ == primitives.UserType {
		identity = &primitives.User{}
	} else if typ == primitives.ServicePrimitiveType {
		identity = &primitives.Service{}
	}

	err = enforcer.Get(ctx, input.IdentityID, identity)
	if err != nil {
		return nil, err
	}

	return &model.Assignment{
		ID:        assignment.ID,
		CreatedAt: assignment.CreatedAt,
		UpdatedAt: assignment.UpdatedAt,
		Role:      role,
		Identity:  identity,
	}, nil
}

func (r *mutationResolver) UnassignRole(ctx context.Context, input model.AssignRoleRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	assignment := &primitives.Assignment{}

	filter := database.NewFilter(database.Where{
		"role_id":     input.RoleID,
		"identity_id": input.IdentityID.String(),
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

func (r *queryResolver) Role(ctx context.Context, id string) (*primitives.Role, error) {
	row := r.Database.Pool.QueryRow(ctx, "SELECT label, system FROM roles WHERE id = $1", id)

	var label primitives.Label
	var system bool
	err := row.Scan(&label, &system)
	if err != nil {
		return nil, err
	}

	role, err := primitives.NewRole(id, label, system)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *queryResolver) RoleByLabel(ctx context.Context, label primitives.Label) (*primitives.Role, error) {
	return nil, nil
}

func (r *queryResolver) Roles(ctx context.Context) ([]*primitives.Role, error) {
	return nil, nil
}

func (r *queryResolver) RoleMembers(ctx context.Context, roleID string) ([]primitives.Identity, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var assignments []*primitives.Assignment
	err := enforcer.Query(ctx, &assignments, database.NewFilter(database.Where{"role_id": roleID}, nil, nil))
	if err != nil {
		return nil, err
	}

	userIDs := database.In{}
	serviceIDs := database.In{}
	for _, assignment := range assignments {
		typ, err := assignment.IdentityID.Type()
		if err != nil {
			return nil, err
		}

		if typ == primitives.UserType {
			userIDs = append(userIDs, assignment.IdentityID)
		} else if typ == primitives.ServicePrimitiveType {
			serviceIDs = append(serviceIDs, assignment.IdentityID)
		}
	}

	var users []*primitives.User
	if len(userIDs) > 0 {
		err = enforcer.Query(ctx, &users, database.NewFilter(database.Where{"id": userIDs}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	var services []*primitives.Service
	if len(serviceIDs) > 0 {
		err = enforcer.Query(ctx, &services, database.NewFilter(database.Where{"id": serviceIDs}, nil, nil))
		if err != nil {
			return nil, err
		}
	}

	identities := make([]primitives.Identity, len(assignments))
	for i, user := range users {
		identities[i] = user
	}

	for i, service := range services {
		identities[i+len(users)] = service
	}

	return identities, nil
}

func (r *roleResolver) ID(ctx context.Context, obj *primitives.Role) (string, error) {
	return obj.ID.String(), nil
}

// Role returns generated.RoleResolver implementation.
func (r *Resolver) Role() generated.RoleResolver { return &roleResolver{r} }

type roleResolver struct{ *Resolver }
