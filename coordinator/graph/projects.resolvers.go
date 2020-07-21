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
	"github.com/capeprivacy/cape/models"
	modelmigration "github.com/capeprivacy/cape/models/migration"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/gosimple/slug"
)

func (r *contributorResolver) User(ctx context.Context, obj *models.Contributor) (*models.User, error) {
	return r.Database.Users().GetByID(ctx, obj.UserID)
}

func (r *contributorResolver) Project(ctx context.Context, obj *models.Contributor) (*primitives.Project, error) {
	p, err := r.Database.Projects().GetByID(ctx, obj.ProjectID)
	if err != nil {
		return nil, err
	}

	return modelmigration.PrimitiveFromProject(*p)
}

func (r *contributorResolver) Role(ctx context.Context, obj *models.Contributor) (*primitives.Role, error) {
	role, err := r.Database.Roles().GetByID(ctx, obj.RoleID)
	if err != nil {
		return nil, err
	}

	return modelmigration.PrimitiveFromRole(*role)
}

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

	// Now make the creator the project owner
	// TODO -- this should happen in a transaction
	l := modelmigration.LabelFromPrimitive(label)
	_, err = r.Database.Contributors().Add(ctx, l, currSession.User.Email, models.ProjectOwnerRole)
	if err != nil {
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
		return nil, errs.New(fw.InvalidParametersCause, "Either id or label must be supplied to updateProject")
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

func (r *mutationResolver) UpdateProjectSpec(ctx context.Context, projectLabel *primitives.Label, projectID *database.ID, request primitives.ProjectSpecFile) (*primitives.Project, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var project primitives.Project
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
	project.Status = primitives.ProjectActive
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

func (r *mutationResolver) ArchiveProject(ctx context.Context, id *database.ID, label *primitives.Label) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnarchiveProject(ctx context.Context, id *database.ID, label *primitives.Label) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateContributor(ctx context.Context, projectLabel models.Label, userEmail models.Email, roleLabel models.Label) (*models.Contributor, error) {
	return r.Database.Contributors().Add(ctx, projectLabel, userEmail, roleLabel)
}

func (r *mutationResolver) RemoveContributor(ctx context.Context, projectLabel models.Label, userEmail models.Email) (*models.Contributor, error) {
	return r.Database.Contributors().Delete(ctx, projectLabel, userEmail)
}

func (r *projectResolver) CurrentSpec(ctx context.Context, obj *primitives.Project) (*primitives.ProjectSpec, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if obj.CurrentSpecID == nil {
		return nil, errs.New(NoActiveSpecCause, "Project %s has no active project spec", obj.Name)
	}

	var projectSpec primitives.ProjectSpec
	err := enforcer.Get(ctx, *obj.CurrentSpecID, &projectSpec)
	return &projectSpec, err
}

func (r *projectResolver) Contributors(ctx context.Context, obj *primitives.Project) ([]*models.Contributor, error) {
	// TODO -- can this not be a copy of listContributors?
	contribs, err := r.Database.Contributors().List(ctx, modelmigration.LabelFromPrimitive(obj.Label))
	if err != nil {
		return nil, err
	}

	contributors := make([]*models.Contributor, 0, len(contribs))
	for _, con := range contribs {
		c := con
		contributors = append(contributors, &c)
	}

	return contributors, nil
}

func (r *projectSpecResolver) Project(ctx context.Context, obj *primitives.ProjectSpec) (*primitives.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) Parent(ctx context.Context, obj *primitives.ProjectSpec) (*primitives.ProjectSpec, error) {
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
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	if id == nil && label == nil {
		return nil, errs.New(fw.InvalidParametersCause, "Must provide an id or label")
	}

	var project primitives.Project
	if id != nil {
		err := enforcer.Get(ctx, *id, &project)
		return &project, err
	}

	// otherwise, get by label
	err := enforcer.QueryOne(ctx, &project, database.NewFilter(database.Where{"label": label}, nil, nil))
	return &project, err
}

func (r *queryResolver) ListContributors(ctx context.Context, projectLabel models.Label) ([]*models.Contributor, error) {
	contribs, err := r.Database.Contributors().List(ctx, projectLabel)
	if err != nil {
		return nil, err
	}

	contributors := make([]*models.Contributor, 0, len(contribs))
	for _, con := range contribs {
		c := con
		contributors = append(contributors, &c)
	}

	return contributors, nil
}

// Contributor returns generated.ContributorResolver implementation.
func (r *Resolver) Contributor() generated.ContributorResolver { return &contributorResolver{r} }

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

// ProjectSpec returns generated.ProjectSpecResolver implementation.
func (r *Resolver) ProjectSpec() generated.ProjectSpecResolver { return &projectSpecResolver{r} }

type contributorResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type projectSpecResolver struct{ *Resolver }
