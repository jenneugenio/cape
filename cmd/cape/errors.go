package main

import errors "github.com/capeprivacy/cape/partyerrors"

var (
	// MissingArgCause is when a user has not supplied a required argument
	MissingArgCause = errors.NewCause(errors.BadRequestCategory, "missing_argument")

	// MissingEnvVarCause is when a user has not supplied a required environment variable
	MissingEnvVarCause = errors.NewCause(errors.BadRequestCategory, "missing_environment_variable")

	// NoUserCause happens when you try to do something on an user that you do not have access to
	// or does not exist
	NoUserCause = errors.NewCause(errors.BadRequestCategory, "user_not_found")

	InvalidPortCause = errors.NewCause(errors.BadRequestCategory, "invalid_port")

	CreateFileCause = errors.NewCause(errors.BadRequestCategory, "create_file")
)
