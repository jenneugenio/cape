package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/capeprivacy/cape/auth"
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
	user, err := r.Database.Users().GetByID(ctx, obj.UserID)
	if err != nil {
		return nil, err
	}

	return r.Database.Roles().GetProjectRole(ctx, user.Email, obj.ProjectID)
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

	if !currSession.Roles.Global.Can(models.CreateProject) {
		return nil, errs.New(auth.AuthorizationFailure, "invalid permissions to create a project")
	}

	p := models.NewProject(project.Name, label, project.Description)
	if err := r.Database.Projects().Create(ctx, p); err != nil {
		return nil, err
	}

	// Now make the creator the project owner
	// TODO -- this should happen in a transaction
	_, err := r.Database.Contributors().Add(ctx, label, currSession.User.Email)
	if err != nil {
		return nil, err
	}

	_, err = r.Database.Roles().SetProjectRole(ctx, currSession.User.Email, label, models.ProjectOwnerRole)
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

func (r *mutationResolver) UpdateProjectSpec(ctx context.Context, id *string, label *models.Label, request model.ProjectSpecFile) (*models.Project, error) {
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
	spec := models.NewPolicy(project.ID, nil, request.Rules, request.Transformations)
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

func (r *mutationResolver) SuggestProjectPolicy(ctx context.Context, label models.Label, name string, description string, request model.ProjectSpecFile) (*models.Suggestion, error) {
	session := fw.Session(ctx)
	role, err := session.Roles.Projects.Get(label)
	if err != nil {
		return nil, err
	}

	if !role.Can(models.SuggestPolicy) {
		return nil, fmt.Errorf("you must be a project contributor to suggest policy changes")
	}

	project, err := r.Database.Projects().Get(ctx, label)
	if err != nil {
		return nil, err
	}

	spec := models.Policy{
		ID:              models.NewID(),
		ProjectID:       project.ID,
		ParentID:        nil,
		Transformations: request.Transformations,
		Rules:           request.Rules,
		Version:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = r.Database.Projects().CreateProjectSpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	suggestion := models.Suggestion{
		ID:          models.NewID(),
		Title:       name,
		Description: description,
		ProjectID:   project.ID,
		PolicyID:    spec.ID,
		State:       models.SuggestionPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = r.Database.Projects().CreateSuggestion(ctx, suggestion)
	if err != nil {
		return nil, err
	}

	return &suggestion, nil
}

func (r *mutationResolver) GetProjectSuggestions(ctx context.Context, label models.Label) ([]*models.Suggestion, error) {
	session := fw.Session(ctx)
	role, err := session.Roles.Projects.Get(label)
	if err != nil {
		return nil, err
	}

	if !role.Can(models.ListPolicySuggestions) {
		return nil, fmt.Errorf("you must be a project contributor to suggest policy changes")
	}

	suggestions, err := r.Database.Projects().GetSuggestions(ctx, label)
	if err != nil {
		return nil, err
	}

	s := make([]*models.Suggestion, len(suggestions))
	for i, sugg := range suggestions {
		suggestion := sugg
		s[i] = &suggestion
	}

	return s, nil
}

func (r *mutationResolver) ApproveProjectSuggestion(ctx context.Context, id string) (*models.Project, error) {
	session := fw.Session(ctx)

	suggestion, err := r.Database.Projects().GetSuggestion(ctx, id)
	if err != nil {
		return nil, err
	}

	projectPolicy, err := r.Database.Projects().GetProjectSpec(ctx, suggestion.PolicyID)
	if err != nil {
		return nil, err
	}

	project, err := r.Database.Projects().GetByID(ctx, suggestion.ProjectID)
	if err != nil {
		return nil, err
	}

	role, err := session.Roles.Projects.Get(project.Label)
	if err != nil {
		return nil, err
	}

	if !role.Can(models.AcceptPolicy) {
		return nil, fmt.Errorf("you must be a project contributor to suggest policy changes")
	}

	// Make this spec active on the project
	project.CurrentSpecID = projectPolicy.ID
	// A spec makes the project active!
	project.Status = models.ProjectActive
	err = r.Database.Projects().Update(ctx, *project)
	if err != nil {
		return nil, err
	}

	suggestion.State = models.SuggestionApproved
	err = r.Database.Projects().UpdateSuggestion(ctx, *suggestion)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *mutationResolver) RejectProjectSuggestion(ctx context.Context, id string) (*models.Project, error) {
	session := fw.Session(ctx)

	suggestion, err := r.Database.Projects().GetSuggestion(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := r.Database.Projects().GetByID(ctx, suggestion.ProjectID)
	if err != nil {
		return nil, err
	}

	role, err := session.Roles.Projects.Get(project.Label)
	if err != nil {
		return nil, err
	}

	if !role.Can(models.RejectPolicy) {
		return nil, fmt.Errorf("you must be a project contributor to reject policy changes")
	}

	suggestion.State = models.SuggestionRejected
	err = r.Database.Projects().UpdateSuggestion(ctx, *suggestion)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *mutationResolver) ArchiveProject(ctx context.Context, id *string, label *models.Label) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnarchiveProject(ctx context.Context, id *string, label *models.Label) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateContributor(ctx context.Context, projectLabel models.Label, userEmail models.Email, roleLabel models.Label) (*models.Contributor, error) {
	_, err := r.Database.Roles().SetProjectRole(ctx, userEmail, projectLabel, roleLabel)
	if err != nil {
		return nil, err
	}

	return r.Database.Contributors().Add(ctx, projectLabel, userEmail)
}

func (r *mutationResolver) RemoveContributor(ctx context.Context, projectLabel models.Label, userEmail models.Email) (*models.Contributor, error) {
	return r.Database.Contributors().Delete(ctx, projectLabel, userEmail)
}

func (r *policyResolver) Project(ctx context.Context, obj *models.Policy) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *policyResolver) Parent(ctx context.Context, obj *models.Policy) (*models.Policy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectResolver) CurrentSpec(ctx context.Context, obj *models.Project) (*models.Policy, error) {
	if obj.CurrentSpecID == "" {
		// If there is no set spec ID, it means the project doesn't yet have a policy
		return nil, nil
	}

	return r.Database.Projects().GetProjectSpec(ctx, obj.CurrentSpecID)
}

func (r *projectResolver) Contributors(ctx context.Context, obj *models.Project) ([]*models.Contributor, error) {
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

	var project *models.Project
	var err error

	if id != nil {
		project, err = r.Database.Projects().GetByID(ctx, *id)
	} else {
		project, err = r.Database.Projects().Get(ctx, *label)
	}

	if err != nil {
		if err.Error() == "no rows" {
			return nil, fmt.Errorf("could not find %s", label)
		}

		return nil, err
	}

	return project, nil
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

func (r *suggestionResolver) Project(ctx context.Context, obj *models.Suggestion) (*models.Project, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *suggestionResolver) Policy(ctx context.Context, obj *models.Suggestion) (*models.Policy, error) {
	panic(fmt.Errorf("not implemented"))
}

// Contributor returns generated.ContributorResolver implementation.
func (r *Resolver) Contributor() generated.ContributorResolver { return &contributorResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Policy returns generated.PolicyResolver implementation.
func (r *Resolver) Policy() generated.PolicyResolver { return &policyResolver{r} }

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

// Suggestion returns generated.SuggestionResolver implementation.
func (r *Resolver) Suggestion() generated.SuggestionResolver { return &suggestionResolver{r} }

type contributorResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type policyResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type suggestionResolver struct{ *Resolver }
