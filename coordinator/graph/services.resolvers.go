package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateService(ctx context.Context, input model.CreateServiceRequest) (*primitives.Service, error) {
	service, err := primitives.NewService(input.Email, input.Type, input.Endpoint)
	if err != nil {
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	// We're building up a list of role labels that we will create assignments
	// for once we've gotten their objects
	roleLabels := []primitives.Label{primitives.GlobalRole}
	if input.Type == primitives.DataConnectorServiceType {
		roleLabels = append(roleLabels, primitives.DataConnectorRole)
	}

	roles, err := getRolesByLabel(ctx, tx, roleLabels)
	if err != nil {
		return nil, err
	}

	err = tx.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	err = createAssignments(ctx, tx, service, roles)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *mutationResolver) DeleteService(ctx context.Context, input model.DeleteServiceRequest) (*string, error) {
	err := r.Backend.Delete(ctx, primitives.ServicePrimitiveType, input.ID)
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
	return getRoles(ctx, r.Backend, obj.ID)
}

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

type serviceResolver struct{ *Resolver }
