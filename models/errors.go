package models

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	InvalidAlgType          = errors.NewCause(errors.BadRequestCategory, "invalid_alg_type")
	InvalidLabelCause       = errors.NewCause(errors.BadRequestCategory, "invalid_label")
	InvalidNameCause        = errors.NewCause(errors.BadRequestCategory, "invalid_name")
	InvalidURLCause         = errors.NewCause(errors.BadRequestCategory, "invalid_url")
	InvalidEmail            = errors.NewCause(errors.BadRequestCategory, "invalid_email")
	InvalidPasswordCause    = errors.NewCause(errors.BadRequestCategory, "invalid_password")
	InvalidDBURLCause       = errors.NewCause(errors.BadRequestCategory, "invalid_db_url")
	InvalidTargetCause      = errors.NewCause(errors.BadRequestCategory, "invalid_target")
	InvalidPolicySpecCause  = errors.NewCause(errors.BadRequestCategory, "invalid_policy_spec")
	InvalidPolicyCause      = errors.NewCause(errors.BadRequestCategory, "invalid_policy")
	InvalidFieldCause       = errors.NewCause(errors.BadRequestCategory, "invalid_field")
	InvalidConfigCause      = errors.NewCause(errors.BadRequestCategory, "invalid_config")
	InvalidSessionCause     = errors.NewCause(errors.BadRequestCategory, "invalid_session")
	InvalidTokenCause       = errors.NewCause(errors.BadRequestCategory, "invalid_token")
	InvalidCredentialsCause = errors.NewCause(errors.BadRequestCategory, "invalid_credentials")
	InvalidUserCause        = errors.NewCause(errors.BadRequestCategory, "invalid_user")
	SystemErrorCause        = errors.NewCause(errors.InternalServerErrorCategory, "system_error")
	InvalidProjectNameCause = errors.NewCause(errors.BadRequestCategory, "invalid_project_name")
	InvalidRecoveryCause    = errors.NewCause(errors.BadRequestCategory, "invalid_recovery")
)
