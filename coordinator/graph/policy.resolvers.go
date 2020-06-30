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
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreatePolicy(ctx context.Context, input model.CreatePolicyRequest) (*models.Policy, error) {
	currSession := fw.Session(ctx)

	label := models.Label(input.Label.String())

	policy, err := r.Database.Policies().Create(ctx, models.NewPolicy(label, input.Spec))
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (r *mutationResolver) DeletePolicy(ctx context.Context, input model.DeletePolicyRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	err := enforcer.Delete(ctx, primitives.PolicyType, input.ID)
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

	return buildAttachment(ctx, enforcer, attachment)
}

func (r *mutationResolver) DetachPolicy(ctx context.Context, input model.DetachPolicyRequest) (*string, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	attachment := &primitives.Attachment{}

	filter := database.NewFilter(database.Where{
		"role_id":   input.RoleID.String(),
		"policy_id": input.PolicyID.String(),
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

func (r *policyResolver) ID(ctx context.Context, obj *models.Policy) (database.ID, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *policyResolver) Label(ctx context.Context, obj *models.Policy) (primitives.Label, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Policy(ctx context.Context, id database.ID) (*models.Policy, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	policy := &primitives.Policy{}
	err := enforcer.Get(ctx, id, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *queryResolver) PolicyByLabel(ctx context.Context, label primitives.Label) (*models.Policy, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	policy := &primitives.Policy{}
	err := enforcer.QueryOne(ctx, policy, database.NewFilter(database.Where{"label": label}, nil, nil))
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *queryResolver) Policies(ctx context.Context) ([]*models.Policy, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var policies []*primitives.Policy
	err := enforcer.Query(ctx, &policies, database.NewEmptyFilter())
	if err != nil {
		return nil, err
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

	var policies []*primitives.Policy
	if len(attachments) == 0 {
		return policies, nil
	}

	policyIDs := database.InFromEntities(attachments, func(e interface{}) interface{} {
		return e.(*primitives.Attachment).PolicyID
	})
	err = enforcer.Query(ctx, &policies, database.NewFilter(database.Where{"id": policyIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (r *queryResolver) IdentityPolicies(ctx context.Context, identityID database.ID) ([]*models.Policy, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	return fw.QueryIdentityPolicies(ctx, enforcer, identityID)
}

func (r *queryResolver) Attachment(ctx context.Context, roleID database.ID, policyID database.ID) (*model.Attachment, error) {
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

	return buildAttachment(ctx, enforcer, attachment)
}

// Policy returns generated.PolicyResolver implementation.
func (r *Resolver) Policy() generated.PolicyResolver { return &policyResolver{r} }

type policyResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *policyResolver) Spec(ctx context.Context, obj *primitives.Policy) (*models.PolicySpec, error) {
	panic(fmt.Errorf("not implemented"))
}
