package ui

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	AbortedCause     = errors.NewCause(errors.BadRequestCategory, "user_aborted")
	ErrAborted       = errors.New(AbortedCause, "Aborted")
	NotAttachedCause = errors.NewCause(errors.BadRequestCategory, "not_attached")
)
