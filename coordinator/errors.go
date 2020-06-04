package coordinator

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	InvalidConfigCause   = errors.NewCause(errors.BadRequestCategory, "invalid_config")
	InvalidArgumentCause = errors.NewCause(errors.BadRequestCategory, "invalid_argument")
)

func errorPresenter(ctx context.Context, e error) *gqlerror.Error {
	pErr, ok := e.(*errors.Error)
	if !ok {
		pErr = errors.New(errors.UnknownCause, e.Error())
	}

	gErr := gqlerror.ErrorPathf(graphql.GetFieldContext(ctx).Path(), strings.Join(pErr.Messages, ","))

	if gErr.Extensions == nil {
		gErr.Extensions = make(map[string]interface{})
	}
	gErr.Extensions["cause"] = pErr.Cause

	return gErr
}
