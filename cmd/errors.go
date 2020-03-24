package main

import errors "github.com/dropoutlabs/cape/partyerrors"

var (
	// InvalidURLCause is when a parsed url is invalid
	InvalidURLCause = errors.NewCause(errors.BadRequestCategory, "invalid_url")

	// MissingArgCause is when a user has not supplied a required argument
	MissingArgCause = errors.NewCause(errors.BadRequestCategory, "missing_argument")
)
