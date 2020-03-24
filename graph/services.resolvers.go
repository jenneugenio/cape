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

func (r *mutationResolver) CreateService(ctx context.Context, input model.CreateServiceRequest) (*primitives.Service, error) {
	return nil, errs.New(RouteNotImplemented, "CreateService mutation not implemented")
}

func (r *mutationResolver) DeleteService(ctx context.Context, input model.DeleteServiceRequest) (*string, error) {
	return nil, errs.New(RouteNotImplemented, "DeleteService mutation not implemented")
}

func (r *queryResolver) Service(ctx context.Context, id database.ID) (*primitives.Service, error) {
	return nil, errs.New(RouteNotImplemented, "Service query not implemented")
}

func (r *queryResolver) ServiceByEmail(ctx context.Context, email string) (*primitives.Service, error) {
	return nil, errs.New(RouteNotImplemented, "ServiceByEmail query not implemented")
}

func (r *queryResolver) Services(ctx context.Context) ([]*primitives.Service, error) {
	return nil, errs.New(RouteNotImplemented, "Services query not implemented")
}
