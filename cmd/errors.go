package cmd

import errors "github.com/dropoutlabs/privacyai/partyerrors"

var (
	// InvalidURLCause is when a parsed url is invalid
	InvalidURLCause = errors.NewCause(errors.BadRequestCategory, "invalid_url")
)
