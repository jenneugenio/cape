package primitives

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	InvalidServiceCause     = errors.NewCause(errors.BadRequestCategory, "invalid_service")
	InvalidServiceType      = errors.NewCause(errors.BadRequestCategory, "invalid_service_type")
	InvalidTokenType        = errors.NewCause(errors.BadRequestCategory, "invalid_token_type")
	InvalidSourceType       = errors.NewCause(errors.BadRequestCategory, "invalid_source_type")
	InvalidSourceCause      = errors.NewCause(errors.BadRequestCategory, "invalid_source")
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
	InvalidAssignmentCause  = errors.NewCause(errors.BadRequestCategory, "invalid_assignment")
	InvalidAttachmentCause  = errors.NewCause(errors.BadRequestCategory, "invalid_attachment")
	InvalidConfigCause      = errors.NewCause(errors.BadRequestCategory, "invalid_config")
	InvalidRoleCause        = errors.NewCause(errors.BadRequestCategory, "invalid_role")
	InvalidSessionCause     = errors.NewCause(errors.BadRequestCategory, "invalid_session")
	InvalidTokenCause       = errors.NewCause(errors.BadRequestCategory, "invalid_token")
	InvalidCredentialsCause = errors.NewCause(errors.BadRequestCategory, "invalid_credentials")
	InvalidUserCause        = errors.NewCause(errors.BadRequestCategory, "invalid_user")
	UnsupportedSchemaCause  = errors.NewCause(errors.BadRequestCategory, "invalid_schema")
)
