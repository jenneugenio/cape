package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	fw "github.com/capeprivacy/cape/framework"

	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/models"
)

func (r *assignmentResolver) Role(ctx context.Context, obj *models.Assignment) (*models.Role, error) {
	return r.Database.Roles().GetByID(ctx, obj.ID)
}

func (r *assignmentResolver) User(ctx context.Context, obj *models.Assignment) (*models.User, error) {
	return r.Database.Users().GetByID(ctx, obj.UserID)
}

func (r *mutationResolver) SetOrgRole(ctx context.Context, userEmail models.Email, roleLabel models.Label) (*models.Assignment, error) {
	if !models.ValidOrgRole(roleLabel) {
		return nil, fmt.Errorf("invalid role: %s", roleLabel)
	}

	return r.Database.Assignments().SetOrg(ctx, userEmail, roleLabel)
}

func (r *mutationResolver) SetProjectRole(ctx context.Context, userEmail models.Email, projectLabel models.Label, roleLabel models.Label) (*models.Assignment, error) {
	if !models.ValidProjectRole(roleLabel) {
		return nil, fmt.Errorf("invalid role: %s", roleLabel)
	}

	return r.Database.Assignments().SetProject(ctx, userEmail, projectLabel, roleLabel)
}

func (r *queryResolver) MyRole(ctx context.Context, projectLabel *models.Label) (*models.Role, error) {
	currSession := fw.Session(ctx)
	return &currSession.Roles[0], nil
}

// Assignment returns generated.AssignmentResolver implementation.
func (r *Resolver) Assignment() generated.AssignmentResolver { return &assignmentResolver{r} }

type assignmentResolver struct{ *Resolver }
