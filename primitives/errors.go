package primitives

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	InvalidServiceCause    = errors.NewCause(errors.BadRequestCategory, "invalid_service")
	InvalidServiceType     = errors.NewCause(errors.BadRequestCategory, "invalid_service_type")
	InvalidTokenType       = errors.NewCause(errors.BadRequestCategory, "invalid_token_type")
	InvalidSourceType      = errors.NewCause(errors.BadRequestCategory, "invalid_source_type")
	InvalidAlgType         = errors.NewCause(errors.BadRequestCategory, "invalid_alg_type")
	InvalidLabelCause      = errors.NewCause(errors.BadRequestCategory, "invalid_label")
	InvalidURLCause        = errors.NewCause(errors.BadRequestCategory, "invalid_url")
	InvalidEmail           = errors.NewCause(errors.BadRequestCategory, "invalid_email")
	InvalidPasswordCause   = errors.NewCause(errors.BadRequestCategory, "invalid_password")
	InvalidDBURLCause      = errors.NewCause(errors.BadRequestCategory, "invalid_db_url")
	InvalidTargetCause     = errors.NewCause(errors.BadRequestCategory, "invalid_target")
	InvalidPolicySpecCause = errors.NewCause(errors.BadRequestCategory, "invalid_policy_spec")
	InvalidPolicyCause     = errors.NewCause(errors.BadRequestCategory, "invalid_policy")
	InvalidFieldCause      = errors.NewCause(errors.BadRequestCategory, "invalid_field")
)
