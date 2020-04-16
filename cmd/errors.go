package main

import errors "github.com/capeprivacy/cape/partyerrors"

var (
	// InvalidAPITokenCause is when a provided APIToken is not valid
	InvalidAPITokenCause = errors.NewCause(errors.BadRequestCategory, "invalid_api_token")

	// MissingArgCause is when a user has not supplied a required argument
	MissingArgCause = errors.NewCause(errors.BadRequestCategory, "missing_argument")

	// MissingEnvVarCause is when a user has not supplied a required environment variable
	MissingEnvVarCause = errors.NewCause(errors.BadRequestCategory, "missing_environment_variable")

	// PasswordNoMatch happens when you confirm you password and it doesn't match your initial input
	PasswordNoMatch = errors.NewCause(errors.BadRequestCategory, "passwords_dont_match")

	// BadCertificate happens when the server cert is bad
	BadCertificate = errors.NewCause(errors.BadRequestCategory, "bad_certificate")

	// MustSupplyEndpoint is used with the service create command to make sure an endpoint is
	// supplied when creating a data connector
	MustSupplyEndpoint = errors.NewCause(errors.BadRequestCategory, "must_supply_endpoint")

	// NoLinkedService occurs when a service has not been linked to a source
	NoLinkedService = errors.NewCause(errors.BadRequestCategory, "no_linked_service")

	// ClusterExistsCause cause happens when you try to create a new cluster that already exists
	ClusterExistsCause = errors.NewCause(errors.BadRequestCategory, "cluster_already_exists")
)
