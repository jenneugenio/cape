package connector

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// MissingAuthToken occurs when there is no auth token found on the interceptor
	// metadata
	MissingAuthToken = errors.NewCause(errors.BadRequestCategory, "missing_auth_token")

	// InvalidConfigCause occurs when the config isn't valid
	InvalidConfigCause = errors.NewCause(errors.BadRequestCategory, "invalid_config")

	// MissingIdentityID occurs when the identity id isn't on the request
	MissingIdentityID = errors.NewCause(errors.BadRequestCategory, "missing_identity_id")

	// FieldNotFound occurs when the data connector is trying to find information on
	// field for a given schema and cannot.
	FieldNotFound = errors.NewCause(errors.BadRequestCategory, "field_not_found")
)
