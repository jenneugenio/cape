package graph

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// RouteNotImplemented occurs when a graphQL route that has not been implemented is invoked
	RouteNotImplemented = errors.NewCause(errors.NotImplementedCategory, "route_not_implemented")
	MembersNotCreated   = errors.NewCause(errors.BadRequestCategory, "members_not_created")
)
