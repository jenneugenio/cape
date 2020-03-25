package main

import errors "github.com/dropoutlabs/cape/partyerrors"

var (
	// InvalidURLCause is when a parsed url is invalid
	InvalidURLCause = errors.NewCause(errors.BadRequestCategory, "invalid_url")

	// MissingArgCause is when a user has not supplied a required argument
	MissingArgCause = errors.NewCause(errors.BadRequestCategory, "missing_argument")

	// InvalidLengthCause is when you enter something that does not meet specified length constraints
	// (e.g. a password that is too short)
	InvalidLengthCause = errors.NewCause(errors.BadRequestCategory, "invalid_length")

	// PasswordNoMatch happens when you confirm you password and it doesn't match your initial input
	PasswordNoMatch = errors.NewCause(errors.BadRequestCategory, "passwords_dont_match")
)
