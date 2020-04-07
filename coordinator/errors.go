package coordinator

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	InvalidConfigCause = errors.NewCause(errors.BadRequestCategory, "invalid_config")
)

func errorPresenter(ctx context.Context, e error) *gqlerror.Error {
	pErr, ok := e.(*errors.Error)
	if !ok {
		return graphql.DefaultErrorPresenter(ctx, e)
	}

	return gqlerror.ErrorPathf(graphql.GetFieldContext(ctx).Path(), strings.Join(pErr.Messages, ","))
}
