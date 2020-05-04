package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) CreateService(ctx context.Context, input model.CreateServiceRequest) (*primitives.Service, error) {
	creds, err := primitives.NewCredentials(&input.PublicKey, &input.Salt)
	if err != nil {
		return nil, err
	}

	service, err := primitives.NewService(input.Email, input.Type, input.Endpoint, creds)
	if err != nil {
		return nil, err
	}

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, tx)

	// We're building up a list of role labels that we will create assignments
	// for once we've gotten their objects
	roleLabels := []primitives.Label{primitives.GlobalRole}
	if input.Type == primitives.DataConnectorServiceType {
		roleLabels = append(roleLabels, primitives.DataConnectorRole)
	}

	roles, err := getRolesByLabel(ctx, enforcer, roleLabels)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	err = createAssignments(ctx, enforcer, service, roles)
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
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	err := enforcer.Delete(ctx, primitives.ServicePrimitiveType, input.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) Service(ctx context.Context, id database.ID) (*primitives.Service, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	service := &primitives.Service{}
	err := enforcer.Get(ctx, id, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *queryResolver) ServiceByEmail(ctx context.Context, email primitives.Email) (*primitives.Service, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	service := &primitives.Service{}
	err := enforcer.QueryOne(ctx, service, database.NewFilter(database.Where{"email": email.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *queryResolver) Services(ctx context.Context) ([]*primitives.Service, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var services []*primitives.Service
	err := enforcer.Query(ctx, &services, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceResolver) Roles(ctx context.Context, obj *primitives.Service) ([]*primitives.Role, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	return fw.QueryRoles(ctx, enforcer, obj.ID)
}

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

type serviceResolver struct{ *Resolver }
