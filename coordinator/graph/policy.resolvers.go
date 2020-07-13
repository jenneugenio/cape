package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

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

func (r *mutationResolver) CreatePolicy(ctx context.Context, input model.CreatePolicyRequest) (*models.Policy, error) {
	label := models.Label(input.Label.String())

	policy := models.NewPolicy(label, &input.Spec)
	err := r.Database.Policies().Create(ctx, policy)
	if err == db.ErrDuplicateKey {
		return nil, ErrDuplicateKey
	}
	if err != nil {
		// TODO: Log this error and update metrics
		return nil, errs.New(errs.UnknownCause, "error saving policy")
	}

	return &policy, nil
}

func (r *mutationResolver) DeletePolicy(ctx context.Context, input model.DeletePolicyRequest) (*string, error) {
	id := models.Label(input.ID)

	err := r.Database.Policies().Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) AttachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*model.Attachment, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	attachment, err := primitives.NewAttachment(input.PolicyID, input.RoleID)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, attachment)
	if err != nil {
		return nil, err
	}

	return buildAttachment(ctx, enforcer, r.Database, attachment)
}

func (r *mutationResolver) DetachPolicy(ctx context.Context, input model.DetachPolicyRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	attachment := &primitives.Attachment{}

	filter := database.NewFilter(database.Where{
		"role_id":   input.RoleID.String(),
		"policy_id": input.PolicyID,
	}, nil, nil)

	err := enforcer.QueryOne(ctx, attachment, filter)
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.AttachmentType, attachment.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) Policy(ctx context.Context, id string) (*models.Policy, error) {
	policy, err := r.Database.Policies().GetByID(ctx, id)
	if err != nil {
		// TODO(thor): Log internal error messages and swallow them before
		// returning to the client.
		return nil, err
	}

	return policy, nil
}

func (r *queryResolver) PolicyByLabel(ctx context.Context, label string) (*models.Policy, error) {
	policy, err := r.Database.Policies().Get(ctx, models.Label(label))
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *queryResolver) Policies(ctx context.Context) ([]*models.Policy, error) {
	pols, err := r.Database.Policies().List(ctx, nil)
	if err != nil {
		return nil, err
	}

	policies := make([]*models.Policy, 0, len(pols))
	for _, pol := range pols {
		p := pol
		policies = append(policies, &p)
	}

	return policies, nil
}

func (r *queryResolver) RolePolicies(ctx context.Context, roleID database.ID) ([]*models.Policy, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var attachments []*primitives.Attachment
	err := enforcer.Query(ctx, &attachments, database.NewFilter(database.Where{"role_id": roleID.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	if len(attachments) == 0 {
		return []*models.Policy{}, nil
	}

	var policyIDs []string
	for _, a := range attachments {
		policyIDs = append(policyIDs, a.PolicyID)
	}

	policies, err := r.Database.Policies().List(ctx, &db.ListPolicyOptions{FilterIDs: policyIDs})
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

func (r *queryResolver) UserPolicies(ctx context.Context, userID string) ([]*models.Policy, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	policies, err := fw.QueryUserPolicies(ctx, enforcer, r.Database, userID)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (r *queryResolver) Attachment(ctx context.Context, roleID database.ID, policyID string) (*model.Attachment, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	attachment := &primitives.Attachment{}

	filter := database.NewFilter(database.Where{
		"role_id":   roleID,
		"policy_id": policyID,
	}, nil, nil)

	err := enforcer.QueryOne(ctx, attachment, filter)
	if err != nil {
		return nil, err
	}

	return buildAttachment(ctx, enforcer, r.Database, attachment)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
