package main

import errors "github.com/capeprivacy/cape/partyerrors"

var (
	// MissingArgCause is when a user has not supplied a required argument
	MissingArgCause = errors.NewCause(errors.BadRequestCategory, "missing_argument")

	// MissingEnvVarCause is when a user has not supplied a required environment variable
	MissingEnvVarCause = errors.NewCause(errors.BadRequestCategory, "missing_environment_variable")

	// PasswordNoMatch happens when you confirm you password and it doesn't match your initial input
	PasswordNoMatch = errors.NewCause(errors.BadRequestCategory, "passwords_dont_match")

	// ClusterExistsCause happens when you try to create a new cluster that already exists
	ClusterExistsCause = errors.NewCause(errors.BadRequestCategory, "cluster_already_exists")

	// NoIdentityCause happens when you try to do something on an identity that you do not have access to
	// or does not exist
	NoIdentityCause = errors.NewCause(errors.BadRequestCategory, "identity_not_found")

	InvalidPortCause = errors.NewCause(errors.BadRequestCategory, "invalid_port")

	CreateFileCause = errors.NewCause(errors.BadRequestCategory, "create_file")
)
