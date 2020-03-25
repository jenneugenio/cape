package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/model"
	errs "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *mutationResolver) CreatePolicy(ctx context.Context, input model.CreatePolicyRequest) (*primitives.Policy, error) {
	policy, err := primitives.NewPolicy(input.Label)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *mutationResolver) DeletePolicy(ctx context.Context, input model.DeletePolicyRequest) (*string, error) {
	err := r.Backend.Delete(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) AttachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*model.Attachment, error) {
	return nil, errs.New(RouteNotImplemented, "Attach policy mutation not implemented")
}

func (r *mutationResolver) DetachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*string, error) {
	return nil, errs.New(RouteNotImplemented, "Detach policy mutation not implemented")
}

func (r *queryResolver) Policy(ctx context.Context, id database.ID) (*primitives.Policy, error) {
	policy := &primitives.Policy{}
	err := r.Backend.Get(ctx, id, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *queryResolver) Policies(ctx context.Context) ([]*primitives.Policy, error) {
	var s []primitives.Policy
	err := r.Backend.Query(ctx, &s, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	policies := make([]*primitives.Policy, len(s))
	for i := 0; i < len(policies); i++ {
		policies[i] = &(s[i])
	}

	return policies, nil
}

func (r *queryResolver) Attachment(ctx context.Context, roleID database.ID, policyID database.ID) (*model.Attachment, error) {
	return nil, errs.New(RouteNotImplemented, "Attachment query not implemented")
}

func (r *queryResolver) Attachments(ctx context.Context) ([]*model.Attachment, error) {
	return nil, errs.New(RouteNotImplemented, "Attachments query not implemented")
}
