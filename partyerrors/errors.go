package partyerrors

import (
	"fmt"
	"strings"
)

// HTTPError represents an error that is returned to a callee as an HTTP response
type HTTPError interface {
	error
	StatusCode() int
	Code() int32
}

// Error represents an Error returned by a middleware or route handler to a
// requestor.
//
// Completes the HTTPError interface
//
// The ID may or may not be a NullID depending on the source of the error - if
// its unmarshaled from the old error format then it will be a NullID.
// Error is an error
type Error struct {
	Cause    Cause    `json:"cause"`
	Messages []string `json:"messages"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Cause.Name, strings.Join(e.Messages, ","))
}

// Code returns the HTTP Status Code associated with this Error
func (e *Error) Code() int32 {
	code, ok := statusCodeForCategory(e.Cause.Category)
	if !ok {
		return int32(500)
	}
	return int32(code)
}

// Validate returns an error if the error is an invalid error
func (e *Error) Validate(_ interface{}) error {
	return nil
}

// StatusCode returns the HTTP Status Code associated with this Error
func (e *Error) StatusCode() int {
	code, ok := statusCodeForCategory(e.Cause.Category)
	if !ok {
		return int(500)
	}
	return code
}

// New creates a new error
func New(c Cause, msg string, args ...interface{}) *Error {
	return &Error{
		Cause:    c,
		Messages: []string{msg},
	}
}

// NewMulti creates a new error containing multiple error messages
func NewMulti(c Cause, msgs []string, args ...interface{}) *Error {
	err := Error{Cause: c}
	for _, msg := range msgs {
		err.Messages = append(err.Messages, fmt.Sprintf(msg, args...))
	}
	return &err
}

// Wrap an error inside an error associated with a Cause
func Wrap(c Cause, err error) *Error {
	return New(c, err.Error())
}

// WrapMulti an multiple errors inside an error associated with a Cause
func WrapMulti(c Cause, errs []error) *Error {
	msgs := []string{}
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return NewMulti(c, msgs)
}

// ToError mutates a given error into an *Error, if the error is not a *Error
// then it mutates it into an InternalServerError with a corresponding message.
func ToError(e error) *Error {
	switch err := e.(type) {
	case *Error:
		return err
	default:
		return New(UnsupportedErrorCause, "Encountered an unknown error")
	}
}

// FromCause returns whether or not an error was caused by the provided cause
func FromCause(e error, c Cause) bool {
	switch err := e.(type) {
	case *Error:
		return err.Cause == c
	default:
		return false
	}
}
