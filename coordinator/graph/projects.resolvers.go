package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/prims"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/gosimple/slug"
)

// CreateProject has been changed to represent a new flow
func (r *mutationResolver) CreateProject(ctx context.Context, project model.CreateProjectRequest) (*prims.Project, error) {
	currSession := fw.Session(ctx)

	// Enforcer is no longer a proxy to the database, it is just used to check permissions
	enforcer := auth.GetEnforcer()

	// Methods on the enforcer are more granular, so rather than getting a bunch of policies attached to a session,
	// We can ask the enforcer if we can do exactly what we are trying to do. In this case, create a project.
	// In contrast, the older system would eventually boil the check to see if the user had a write permission
	// on the projects table
	if err := enforcer.CreateProject(currSession); err != nil {
		return nil, err
	}

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

	// Other proposed change, is that primitive constructors can no longer error. The main motivation
	// for this is that writing tests with primitives this way is annoying. You could argue the ctor is safer,
	// but you could also just skip the ctor and instantiate the struct directly to get around the error
	p := prims.NewProject(project.Name, label, project.Description)
	// Still though, check that the primitive is valid
	if err := p.Validate(); err != nil {
		return nil, err
	}

	// Not using `Backend` anymore. Eventually Backend would not be used then deleted.
	if err := r.Database.CreateProject(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *mutationResolver) UpdateProject(ctx context.Context, id *database.ID, label *primitives.Label, update model.UpdateProjectRequest) (*prims.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	// The way I did this would require us to update all of the routes around a project at once
	// Alternatively, we could have a primitives.Project & a prims.Project, then delete primitives.Project once ALL
	// of the routes are fixed. For most objects we could do everything at once (kinda sketchy to have both versions
	// around for any significant amount of time).
	project := &prims.Project{}
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

func (r *mutationResolver) UpdateProjectSpec(ctx context.Context, projectLabel *primitives.Label, projectID *database.ID, request primitives.ProjectSpecFile) (*prims.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var project prims.Project
	err := enforcer.QueryOne(ctx, &project, database.NewFilter(database.Where{"label": projectLabel.String()}, nil, nil))
	if err != nil {
		return nil, err
	}
	// Insert the spec
	// TODO -- How do you specify the parent? This concept doesn't make sense until we have proposals & diffing
	spec, err := primitives.NewProjectSpec(project.ID, nil, request.Policy)
	if err != nil {
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	enforcer = auth.NewEnforcer(currSession, tx)
	err = enforcer.Create(ctx, spec)
	if err != nil {
		return nil, err
	}

	// Make this spec active on the project
	project.CurrentSpecID = &spec.ID
	// A spec makes the project active!
	project.Status = prims.ProjectActive
	err = enforcer.Update(ctx, &project)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *mutationResolver) ArchiveProject(ctx context.Context, id *database.ID, label *primitives.Label) (*prims.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnarchiveProject(ctx context.Context, id *database.ID, label *primitives.Label) (*prims.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectResolver) CurrentSpec(ctx context.Context, obj *prims.Project) (*primitives.ProjectSpec, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if obj.CurrentSpecID == nil {
		return nil, errs.New(NoActiveSpecCause, "Project %s has no active project spec", obj.Name)
	}

	var projectSpec primitives.ProjectSpec
	err := enforcer.Get(ctx, *obj.CurrentSpecID, &projectSpec)
	return &projectSpec, err
}

func (r *projectSpecResolver) Project(ctx context.Context, obj *primitives.ProjectSpec) (*prims.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) Parent(ctx context.Context, obj *primitives.ProjectSpec) (*primitives.ProjectSpec, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Projects(ctx context.Context, status []prims.ProjectStatus) ([]*prims.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var projects []*prims.Project
	err := enforcer.Query(ctx, &projects, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *queryResolver) Project(ctx context.Context, id *database.ID, label *primitives.Label) (*prims.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if id == nil && label == nil {
		return nil, errs.New(InvalidParametersCause, "Must provide an id or label")
	}

	var project prims.Project
	if id != nil {
		err := enforcer.Get(ctx, *id, &project)
		return &project, err
	}

	// otherwise, get by label
	err := enforcer.QueryOne(ctx, &project, database.NewFilter(database.Where{"label": label}, nil, nil))
	return &project, err
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

// ProjectSpec returns generated.ProjectSpecResolver implementation.
func (r *Resolver) ProjectSpec() generated.ProjectSpecResolver { return &projectSpecResolver{r} }

type projectResolver struct{ *Resolver }
type projectSpecResolver struct{ *Resolver }
