package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/client"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/generated"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	fw "github.com/capeprivacy/cape/framework"
	errs "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) AddSource(ctx context.Context, input model.AddSourceRequest) (*primitives.Source, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	if input.ServiceID != nil {
		service := &primitives.Service{}
		err := enforcer.Get(ctx, *input.ServiceID, service)
		if err != nil {
			return nil, err
		}

		if service.Type != primitives.DataConnectorServiceType {
			return nil, errs.New(MustBeDataConnector, "Linking service to data source must be a data connector")
		}
	}

	source, err := primitives.NewSource(input.Label, &input.Credentials, input.ServiceID)
	if err != nil {
		return nil, err
	}

	err = enforcer.Create(ctx, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *mutationResolver) UpdateSource(ctx context.Context, input model.UpdateSourceRequest) (*primitives.Source, error) {
	session := fw.Session(ctx)

	tx, err := r.Backend.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // nolint: errcheck

	enforcer := auth.NewEnforcer(session, tx)

	source := &primitives.Source{}
	err = enforcer.QueryOne(ctx, source, database.NewFilter(database.Where{"label": input.SourceLabel}, nil, nil))
	if err != nil {
		return nil, err
	}

	if input.ServiceID != nil {
		service := &primitives.Service{}
		err := enforcer.Get(ctx, *input.ServiceID, service)
		if err != nil {
			return nil, err
		}

		if service.Type != primitives.DataConnectorServiceType {
			return nil, errs.New(MustBeDataConnector, "Linking service to data source must be a data connector")
		}
	}

	source.ServiceID = input.ServiceID
	err = enforcer.Update(ctx, source)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *mutationResolver) RemoveSource(ctx context.Context, input model.RemoveSourceRequest) (*string, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	source := primitives.Source{}
	filter := database.Filter{Where: database.Where{"label": input.Label}}
	err := enforcer.QueryOne(ctx, &source, filter)
	if err != nil {
		return nil, err
	}

	err = enforcer.Delete(ctx, primitives.SourcePrimitiveType, source.ID)
	return nil, err
}

func (r *queryResolver) Sources(ctx context.Context) ([]*primitives.Source, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	var sources []*primitives.Source
	err := enforcer.Query(ctx, &sources, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func (r *queryResolver) Source(ctx context.Context, id database.ID) (*primitives.Source, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	source := &primitives.Source{}
	err := enforcer.Get(ctx, id, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *queryResolver) SourceByLabel(ctx context.Context, label primitives.Label) (*primitives.Source, error) {
	currSession := fw.Session(ctx)
	enforcer := auth.NewEnforcer(currSession, r.Backend)

	source := &primitives.Source{}
	err := enforcer.QueryOne(ctx, source, database.NewFilter(
		database.Where{"label": label.String()},
		nil, nil))
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *sourceResolver) Credentials(ctx context.Context, obj *primitives.Source) (*primitives.DBURL, error) {
	session := fw.Session(ctx)
	identity := session.Identity

	if obj.ServiceID != nil && identity.GetID() == *obj.ServiceID {
		return obj.Credentials, nil
	}

	return nil, nil
}

func (r *sourceResolver) Service(ctx context.Context, obj *primitives.Source) (*primitives.Service, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	if obj.ServiceID == nil {
		return nil, nil
	}

	service := &primitives.Service{}
	err := enforcer.Get(ctx, *obj.ServiceID, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *sourceResolver) Schema(ctx context.Context, obj *primitives.Source, opts *client.SchemaOptions) (*primitives.Schema, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	var schema primitives.Schema
	err := enforcer.QueryOne(ctx, &schema, database.NewFilter(database.Where{"source_id": obj.ID.String()}, nil, nil))

	if opts != nil && opts.Definition != "" {
		schema.Definition = primitives.SchemaDefinition{
			opts.Definition: schema.Definition[opts.Definition],
		}
	}

	return &schema, err
}

// Source returns generated.SourceResolver implementation.
func (r *Resolver) Source() generated.SourceResolver { return &sourceResolver{r} }

type sourceResolver struct{ *Resolver }
