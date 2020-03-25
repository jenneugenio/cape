package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/model"
	errs "github.com/dropoutlabs/cape/partyerrors"
)

func (r *mutationResolver) CreatePolicy(ctx context.Context, input model.CreatePolicyRequest) (*model.Policy, error) {
	return nil, errs.New(RouteNotImplemented, "Create policy mutation not implemented")
}

func (r *mutationResolver) DeletePolicy(ctx context.Context, input model.DeletePolicyRequest) (*string, error) {
	return nil, errs.New(RouteNotImplemented, "Delete policy mutation not implemented")
}

func (r *mutationResolver) AttachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*model.Attachment, error) {
	return nil, errs.New(RouteNotImplemented, "Attach policy mutation not implemented")
}

func (r *mutationResolver) DetachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*string, error) {
	return nil, errs.New(RouteNotImplemented, "Detach policy mutation not implemented")
}

func (r *queryResolver) Policy(ctx context.Context, id database.ID) (*model.Policy, error) {
	return nil, errs.New(RouteNotImplemented, "Policy query not implemented")
}

func (r *queryResolver) Policies(ctx context.Context) ([]*model.Policy, error) {
	return nil, errs.New(RouteNotImplemented, "Policies query not implemented")
}

func (r *queryResolver) Attachment(ctx context.Context, roleID database.ID, policyID database.ID) (*model.Attachment, error) {
	return nil, errs.New(RouteNotImplemented, "Attachment query not implemented")
}

func (r *queryResolver) Attachments(ctx context.Context) ([]*model.Attachment, error) {
	return nil, errs.New(RouteNotImplemented, "Attachments query not implemented")
}
