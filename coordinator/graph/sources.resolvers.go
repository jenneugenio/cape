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
	"github.com/capeprivacy/cape/primitives"
)

func (r *mutationResolver) AddSource(ctx context.Context, input model.AddSourceRequest) (*primitives.Source, error) {
	session := fw.Session(ctx)
	enforcer := auth.NewEnforcer(session, r.Backend)

	source, err := primitives.NewSource(input.Label, &input.Credentials)
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
	panic(fmt.Errorf("not implemented"))
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
	err := enforcer.QueryOne(ctx, source, database.NewFilter(database.Where{"label": label.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *sourceResolver) Credentials(ctx context.Context, obj *primitives.Source) (*primitives.DBURL, error) {
	return nil, nil
}

// Source returns generated.SourceResolver implementation.
func (r *Resolver) Source() generated.SourceResolver { return &sourceResolver{r} }

type sourceResolver struct{ *Resolver }
