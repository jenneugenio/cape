package auth

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	RequiredMsgCause = errors.NewCause(errors.BadRequestCategory, "required_msg")
	RequiredSigCause = errors.NewCause(errors.BadRequestCategory, "required_signature")

	RequiredPrivateKeyCause = errors.NewCause(errors.BadRequestCategory, "required_private_key")
	RequiredPublicKeyCause  = errors.NewCause(errors.BadRequestCategory, "required_public_key")

	SignatureNotValid = errors.NewCause(errors.UnauthorizedCategory, "signature_not_valid")

	BadAPITokenVersion = errors.NewCause(errors.BadRequestCategory, "bad_apitoken_version")
	BadSaltLength      = errors.NewCause(errors.BadRequestCategory, "bad_salt_length")
	BadSecretLength    = errors.NewCause(errors.BadRequestCategory, "bad_secret_length")
	BadPublicKeyLength = errors.NewCause(errors.BadRequestCategory, "bad_public_key_length")

	// BadTokenFormat happens when an APIToken has a bad format
	BadTokenFormat = errors.NewCause(errors.BadRequestCategory, "bad_token_format")
)
