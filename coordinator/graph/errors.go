package graph

import (
	errors "github.com/capeprivacy/cape/partyerrors"
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

	// CannotDeleteSystemRole occurs when deletion of a system role is attempted
	CannotDeleteSystemRole = errors.NewCause(errors.ForbiddenCategory, "cannot_delete_system_role")

	// PolicyNotSupplied occurs when a policy has not been supplied for attachPolicy route.
	// Must either supply policy ID or a policy input object
	PolicyNotSupplied = errors.NewCause(errors.BadRequestCategory, "policy_not_supplied")

	// NotFoundCause occurs when an entity we were looking for could not be found
	NotFoundCause = errors.NewCause(errors.BadRequestCategory, "not_found")
)
