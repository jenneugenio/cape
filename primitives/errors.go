package primitives

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	InvalidTokenType  = errors.NewCause(errors.BadRequestCategory, "invalid_token_type")
	InvalidAlgType    = errors.NewCause(errors.BadRequestCategory, "invalid_alg_type")
	InvalidLabelCause = errors.NewCause(errors.BadRequestCategory, "invalid_label")
	InvalidEmail      = errors.NewCause(errors.BadRequestCategory, "invalid_email")
	InvalidName       = errors.NewCause(errors.BadRequestCategory, "invalid_name")
)
