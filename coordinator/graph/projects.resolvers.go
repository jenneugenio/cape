package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/gosimple/slug"
)

func (r *contributorResolver) User(ctx context.Context, obj *models.Contributor) (*models.User, error) {
	return r.Database.Users().GetByID(ctx, obj.UserID)
}

func (r *contributorResolver) Project(ctx context.Context, obj *models.Contributor) (*models.Project, error) {
	return r.Database.Projects().GetByID(ctx, obj.ProjectID)
}

func (r *contributorResolver) Role(ctx context.Context, obj *models.Contributor) (*models.Role, error) {
	return r.Database.Roles().GetByID(ctx, obj.RoleID)
}

func (r *mutationResolver) CreateProject(ctx context.Context, project model.CreateProjectRequest) (*models.Project, error) {
	currSession := fw.Session(ctx)
	var label models.Label
	if project.Label != nil {
		label = *project.Label
	} else {
		labelStr := slug.Make(project.Name.String())
		label = models.Label(labelStr)
	}

	p := models.NewProject(project.Name, label, project.Description)
	if err := r.Database.Projects().Create(ctx, p); err != nil {
		return nil, err
	}

	// Now make the creator the project owner
	// TODO -- this should happen in a transaction
	_, err := r.Database.Contributors().Add(ctx, label, currSession.User.Email, models.ProjectOwnerRole)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *mutationResolver) UpdateProject(ctx context.Context, id *string, label *models.Label, update model.UpdateProjectRequest) (*models.Project, error) {
	var project *models.Project
	var err error
	if id != nil && *id != "" {
		project, err = r.Database.Projects().GetByID(ctx, *id)
	} else if label != nil {
		project, err = r.Database.Projects().Get(ctx, *label)
	} else {
		return nil, errs.New(fw.InvalidParametersCause, "either id or label must be supplied to updateProject")
	}

	if err != nil {
		return nil, errs.New(fw.InvalidParametersCause, "could not find the requested project")
	}

	if update.Name != nil {
		project.Name = *update.Name
	}

	if update.Description != nil {
		project.Description = *update.Description
	}

	err = r.Database.Projects().Update(ctx, *project)
	return project, err
}

func (r *mutationResolver) UpdateProjectSpec(ctx context.Context, id *string, label *models.Label, request models.ProjectSpecFile) (*models.Project, error) {
	var project *models.Project
	var err error
	if id != nil && *id != "" {
		project, err = r.Database.Projects().GetByID(ctx, *id)
	} else if label != nil {
		project, err = r.Database.Projects().Get(ctx, *label)
	} else {
		return nil, errs.New(fw.InvalidParametersCause, "either id or label must be supplied to updateProject")
	}

	if err != nil {
		return nil, errs.New(fw.InvalidParametersCause, "could not find the requested project")
	}

	// Insert the spec
	// TODO -- How do you specify the parent? This concept doesn't make sense until we have proposals & diffing
	spec := models.NewProjectSpec(project.ID, nil, request.Policy)
	err = r.Database.Projects().CreateProjectSpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	// Make this spec active on the project
	project.CurrentSpecID = spec.ID
	// A spec makes the project active!
	project.Status = models.ProjectActive
	err = r.Database.Projects().Update(ctx, *project)
	return project, err
}

func (r *mutationResolver) ArchiveProject(ctx context.Context, id *string, label *models.Label) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnarchiveProject(ctx context.Context, id *string, label *models.Label) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateContributor(ctx context.Context, projectLabel models.Label, userEmail models.Email, roleLabel models.Label) (*models.Contributor, error) {
	return r.Database.Contributors().Add(ctx, projectLabel, userEmail, roleLabel)
}

func (r *mutationResolver) RemoveContributor(ctx context.Context, projectLabel models.Label, userEmail models.Email) (*models.Contributor, error) {
	return r.Database.Contributors().Delete(ctx, projectLabel, userEmail)
}

func (r *projectResolver) CurrentSpec(ctx context.Context, obj *models.Project) (*models.ProjectSpec, error) {
	if obj.CurrentSpecID == "" {
		return nil, errs.New(NoActiveSpecCause, "Project %s has no active project spec", obj.Name)
	}

	return r.Database.Projects().GetProjectSpec(ctx, obj.CurrentSpecID)
}

func (r *projectResolver) Contributors(ctx context.Context, obj *models.Project) ([]*models.Contributor, error) {
	// TODO -- can this not be a copy of listContributors?
	contribs, err := r.Database.Contributors().List(ctx, obj.Label)
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

func (r *projectSpecResolver) Project(ctx context.Context, obj *models.ProjectSpec) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) Parent(ctx context.Context, obj *models.ProjectSpec) (*models.ProjectSpec, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) CreatedAt(ctx context.Context, obj *models.ProjectSpec) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectSpecResolver) UpdatedAt(ctx context.Context, obj *models.ProjectSpec) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Projects(ctx context.Context, status models.ProjectStatus) ([]*models.Project, error) {
	if err := status.Validate(); err != nil {
		return nil, err
	}

	var projects []models.Project
	var err error
	if status == models.Any {
		projects, err = r.Database.Projects().List(ctx)
	} else {
		projects, err = r.Database.Projects().ListByStatus(ctx, status)
	}

	if err != nil {
		return nil, err
	}

	res := make([]*models.Project, len(projects))
	for i, p := range projects {
		res[i] = &p
	}

	return res, nil
}

func (r *queryResolver) Project(ctx context.Context, id *string, label *models.Label) (*models.Project, error) {
	if id == nil && label == nil {
		return nil, errs.New(fw.InvalidParametersCause, "Must provide an id or label")
	}

	if id != nil {
		return r.Database.Projects().GetByID(ctx, *id)
	}

	// otherwise, get by label
	return r.Database.Projects().Get(ctx, *label)
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