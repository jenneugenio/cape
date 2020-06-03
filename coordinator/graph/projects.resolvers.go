package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project model.CreateProjectRequest) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateProject(ctx context.Context, project model.UpdateSourceRequest) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateProjectSpec(ctx context.Context, projectSpec model.CreateProjectSpecRequest) (*primitives.ProjectSpec, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ArchiveProject(ctx context.Context, id *database.ID, label *primitives.Label) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnarchiveProject(ctx context.Context, id *database.ID, label *primitives.Label) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectResolver) CurrentSpec(ctx context.Context, obj *primitives.Project) (*primitives.ProjectSpec, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) Project(ctx context.Context, obj *primitives.ProjectSpec) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) Parent(ctx context.Context, obj *primitives.ProjectSpec) (*primitives.ProjectSpec, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) Sources(ctx context.Context, obj *primitives.ProjectSpec) ([]*primitives.Source, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Projects(ctx context.Context, status []primitives.ProjectStatus) ([]*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Project(ctx context.Context, id *database.ID, label *primitives.Label) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

// ProjectSpec returns generated.ProjectSpecResolver implementation.
func (r *Resolver) ProjectSpec() generated.ProjectSpecResolver { return &projectSpecResolver{r} }

type projectResolver struct{ *Resolver }
type projectSpecResolver struct{ *Resolver }
