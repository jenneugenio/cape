package transformations

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	// UnsupportedType happens when a transform encounters a type it does not support
	UnsupportedType = errors.NewCause(errors.BadRequestCategory, "unsupported_type")
	// MissingArgument happens when an argument is missing to apply the transformation
	MissingArgument = errors.NewCause(errors.BadRequestCategory, "missing_arg")
	// WrongArgument happens when a wrong combination of argument is apply to the transformation
	WrongArgument = errors.NewCause(errors.BadRequestCategory, "wrong_arg")

	// EvaluateBoolOnly occurs when a govalute expression evaluates to something other than
	// bool
	EvaluateBoolOnly = errors.NewCause(errors.BadRequestCategory, "evaluate_bool_only")

	// FieldNotFound occurs when the data connector is trying to find information on
	// field for a given schema and cannot.
	FieldNotFound = errors.NewCause(errors.BadRequestCategory, "field_not_found")

	// InvalidFieldType occurs when a fields type can not be accounted for
	InvalidFieldType = errors.NewCause(errors.BadRequestCategory, "invalid_field_type")

	// TransformationNotFound occurs when a transformation is not registered properly
	TransformationNotFound = errors.NewCause(errors.BadRequestCategory, "transformation_not_found")
)
