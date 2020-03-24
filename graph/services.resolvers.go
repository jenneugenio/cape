package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/model"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *mutationResolver) CreateService(ctx context.Context, input model.CreateServiceRequest) (*primitives.Service, error) {
	creds := &primitives.Credentials{
		PublicKey: &input.PublicKey,
		Salt:      &input.Salt,
		Alg:       input.Alg,
	}

	service, err := primitives.NewService(input.Email, creds)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *mutationResolver) DeleteService(ctx context.Context, input model.DeleteServiceRequest) (*string, error) {
	err := r.Backend.Delete(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) Service(ctx context.Context, id database.ID) (*primitives.Service, error) {
	service := &primitives.Service{}
	err := r.Backend.Get(ctx, id, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *queryResolver) ServiceByEmail(ctx context.Context, email string) (*primitives.Service, error) {
	service := &primitives.Service{}
	err := r.Backend.QueryOne(ctx, service, database.NewFilter(database.Where{"email": email}, nil, nil))
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *queryResolver) Services(ctx context.Context) ([]*primitives.Service, error) {
	var s []primitives.Service
	err := r.Backend.Query(ctx, &s, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	services := make([]*primitives.Service, len(s))
	for i := 0; i < len(services); i++ {
		services[i] = &(s[i])
	}

	return services, nil
}
