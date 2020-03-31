package primitives

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	InvalidServiceCause  = errors.NewCause(errors.BadRequestCategory, "invalid_service")
	InvalidServiceType   = errors.NewCause(errors.BadRequestCategory, "invalid_service_type")
	InvalidTokenType     = errors.NewCause(errors.BadRequestCategory, "invalid_token_type")
	InvalidAlgType       = errors.NewCause(errors.BadRequestCategory, "invalid_alg_type")
	InvalidLabelCause    = errors.NewCause(errors.BadRequestCategory, "invalid_label")
	InvalidURLCause      = errors.NewCause(errors.BadRequestCategory, "invalid_url")
	InvalidEmail         = errors.NewCause(errors.BadRequestCategory, "invalid_email")
	InvalidName          = errors.NewCause(errors.BadRequestCategory, "invalid_name")
	InvalidPasswordCause = errors.NewCause(errors.BadRequestCategory, "invalid_password")
	InvalidDBURLCause    = errors.NewCause(errors.BadRequestCategory, "invalid_db_url")
)
