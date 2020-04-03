package sources

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	ClosingCause        = errors.NewCause(errors.BadRequestCategory, "could_not_close")
	SourceNotSupported  = errors.NewCause(errors.NotFoundCategory, "source_not_supported")
	SourceAlreadyExists = errors.NewCause(errors.BadRequestCategory, "source_already_exists")
	ClosedCause         = errors.NewCause(errors.BadRequestCategory, "closed")
	NotFoundCause       = errors.NewCause(errors.NotFoundCategory, "source_not_found")
	WrongSourceCause    = errors.NewCause(errors.BadRequestCategory, "wrong_source")
	InvalidConfig       = errors.NewCause(errors.BadRequestCategory, "invalid_config")

	// UnknownFieldType occurs when an unknown field type is encountered
	UnknownFieldType = errors.NewCause(errors.InternalServerErrorCategory, "unknown_field_type")

	ErrCacheClosed   = errors.New(ClosedCause, "Cache has been closed")
	ErrCacheNotFound = errors.New(NotFoundCause, "Source not found in cache")
	ErrWrongSource   = errors.New(WrongSourceCause, "Query made against wrong source")
)
