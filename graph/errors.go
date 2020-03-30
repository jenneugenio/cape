package graph

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// RouteNotImplemented occurs when a graphQL route that has not been implemented is invoked
	RouteNotImplemented = errors.NewCause(errors.NotImplementedCategory, "route_not_implemented")

	// MembersNotCreated occurs when members couldn't be created as a part of
	// creating a role
	MembersNotCreated = errors.NewCause(errors.BadRequestCategory, "members_not_created")

	// MustBeDataConnector occurs when linking a service to a data source. The service must be
	// of data connector type
	MustBeDataConnector = errors.NewCause(errors.BadRequestCategory, "must_be_data_connector")
)
