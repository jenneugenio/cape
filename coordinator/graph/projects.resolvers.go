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
	"github.com/gosimple/slug"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project model.CreateProjectRequest) (*primitives.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var label primitives.Label
	if project.Label != nil {
		label = *project.Label
	} else {
		labelStr := slug.Make(project.Name.String())
		l, err := primitives.NewLabel(labelStr)
		if err != nil {
			return nil, err
		}

		label = l
	}

	p, err := primitives.NewProject(project.Name, label, project.Description)
	if err != nil {
		return nil, err
	}

	if err := enforcer.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *mutationResolver) UpdateProject(ctx context.Context, id *database.ID, label *primitives.Label, update model.UpdateProjectRequest) (*primitives.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	project := &primitives.Project{}
	if id != nil {
		err := enforcer.Get(ctx, *id, project)
		if err != nil {
			return nil, err
		}
	} else if label != nil {
		err := enforcer.QueryOne(ctx, project, database.NewFilter(database.Where{"label": label}, nil, nil))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errs.New(InvalidParametersCause, "Either id or label must be supplied to updateProject")
	}

	if update.Name != nil {
		project.Name = *update.Name
	}

	if update.Description != nil {
		project.Description = *update.Description
	}

	err := enforcer.Update(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
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
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var projects []*primitives.Project
	err := enforcer.Query(ctx, &projects, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return projects, nil
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
