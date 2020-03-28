package main

import errors "github.com/dropoutlabs/cape/partyerrors"

var (
	// InvalidURLCause is when a parsed url is invalid
	InvalidURLCause = errors.NewCause(errors.BadRequestCategory, "invalid_url")

	// InvalidAPITokenCause is when a provided APIToken is not valid
	InvalidAPITokenCause = errors.NewCause(errors.BadRequestCategory, "invalid_api_token")

	// MissingArgCause is when a user has not supplied a required argument
	MissingArgCause = errors.NewCause(errors.BadRequestCategory, "missing_argument")

	// MissingEnvVarCause is when a user has not supplied a required environment variable
	MissingEnvVarCause = errors.NewCause(errors.BadRequestCategory, "missing_environment_variable")

	// InvalidLengthCause is when you enter something that does not meet specified length constraints
	// (e.g. a password that is too short)
	InvalidLengthCause = errors.NewCause(errors.BadRequestCategory, "invalid_length")

	// PasswordNoMatch happens when you confirm you password and it doesn't match your initial input
	PasswordNoMatch = errors.NewCause(errors.BadRequestCategory, "passwords_dont_match")
)
