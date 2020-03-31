package ui

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// AbortedCause represents types of errors where the user aborted the operation
	AbortedCause = errors.NewCause(errors.BadRequestCategory, "user_aborted")

	// NotAttachedCause represents types of errors that are caused when a ui
	// element is attempted to be used that are not supported by the terminal
	NotAttachedCause = errors.NewCause(errors.BadRequestCategory, "not_attached")

	// CannotDisplayCause represents types of errors that are caused when a
	// provided value cannot be displayed by a ui component
	CannotDisplayCause = errors.NewCause(errors.BadRequestCategory, "cannot_display")

	// ErrAborted is an error when a user aborts while filling out a prompt
	ErrAborted = errors.New(AbortedCause, "Aborted")

	// ErrCantDisplay happens when a value is given to a UI Component that
	// cannot be displayed (e.g. isn't a string or doesn't have a String()
	// method)
	ErrCantDisplay = errors.New(CannotDisplayCause, "Cannot display value, expected a string or fmt.Stringer")
)
