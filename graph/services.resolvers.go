package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/generated"
	"github.com/dropoutlabs/cape/graph/model"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *mutationResolver) CreateService(ctx context.Context, input model.CreateServiceRequest) (*primitives.Service, error) {
	creds := &primitives.Credentials{
		PublicKey: &input.PublicKey,
		Salt:      &input.Salt,
		Alg:       input.Alg,
	}

	service, err := primitives.NewService(input.Email, input.Type, creds)
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

func (r *serviceResolver) Roles(ctx context.Context, obj *primitives.Service) ([]*primitives.Role, error) {
	assignments := []primitives.Assignment{}
	err := r.Backend.Query(ctx, &assignments, database.NewFilter(database.Where{"identity_id": obj.ID.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	if len(assignments) == 0 {
		return nil, nil
	}

	roleIDs := make(database.In, len(assignments))
	for i, assignment := range assignments {
		roleIDs[i] = assignment.RoleID
	}

	tmpR := []primitives.Role{}
	err = r.Backend.Query(ctx, &tmpR, database.NewFilter(database.Where{"id": roleIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	roles := make([]*primitives.Role, len(tmpR))
	for i := 0; i < len(roles); i++ {
		roles[i] = &(tmpR[i])
	}

	return roles, nil
}

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

type serviceResolver struct{ *Resolver }
