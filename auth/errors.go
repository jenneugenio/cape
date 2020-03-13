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

	BadSaltLength      = errors.NewCause(errors.BadRequestCategory, "bad_salt_length")
	BadSecretLength    = errors.NewCause(errors.BadRequestCategory, "bad_secret_length")
	BadPublicKeyLength = errors.NewCause(errors.BadRequestCategory, "bad_public_key_length")
)
