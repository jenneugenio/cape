package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/controller/graph/generated"
	"github.com/dropoutlabs/cape/controller/graph/model"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *mutationResolver) CreateService(ctx context.Context, input model.CreateServiceRequest) (*primitives.Service, error) {
	if input.Type == primitives.DataConnectorServiceType {
		return createDataConnector(ctx, r.Backend, input)
	}

	// create a non-data connector service!!!
	creds := &primitives.Credentials{
		PublicKey: &input.PublicKey,
		Salt:      &input.Salt,
		Alg:       input.Alg,
	}

	service, err := primitives.NewService(input.Email, input.Type, nil, creds)
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

func (r *queryResolver) ServiceByEmail(ctx context.Context, email primitives.Email) (*primitives.Service, error) {
	service := &primitives.Service{}
	err := r.Backend.QueryOne(ctx, service, database.NewFilter(database.Where{"email": email.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *queryResolver) Services(ctx context.Context) ([]*primitives.Service, error) {
	services := []*primitives.Service{}
	err := r.Backend.Query(ctx, &services, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceResolver) Roles(ctx context.Context, obj *primitives.Service) ([]*primitives.Role, error) {
	assignments := []*primitives.Assignment{}
	err := r.Backend.Query(ctx, &assignments, database.NewFilter(database.Where{"identity_id": obj.ID.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	if len(assignments) == 0 {
		return nil, nil
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	roles := []*primitives.Role{}
	err = r.Backend.Query(ctx, &roles, database.NewFilter(database.Where{"id": roleIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

type serviceResolver struct{ *Resolver }
