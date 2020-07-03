package graph

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// MembersNotCreated occurs when members couldn't be created as a part of
	// creating a role
	MembersNotCreated = errors.NewCause(errors.BadRequestCategory, "members_not_created")

	// CannotDeleteSystemRole occurs when deletion of a system role is attempted
	CannotDeleteSystemRole = errors.NewCause(errors.ForbiddenCategory, "cannot_delete_system_role")

	// PolicyNotSupplied occurs when a policy has not been supplied for attachPolicy route.
	// Must either supply policy ID or a policy input object
	PolicyNotSupplied = errors.NewCause(errors.BadRequestCategory, "policy_not_supplied")

	InvalidSource = errors.NewCause(errors.BadRequestCategory, "invalid_source")

	NoActiveSpecCause = errors.NewCause(errors.BadRequestCategory, "no_active_spec")

	RecoveryFailedCause = errors.NewCause(errors.UnauthorizedCategory, "recovery_failed")
	ErrRecoveryFailed   = errors.New(RecoveryFailedCause, "recovery_failed")
)
