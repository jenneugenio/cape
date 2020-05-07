package crypto

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	KMSDecryptCause       = errors.NewCause(errors.BadRequestCategory, "kms_decrypt")
	SecretBoxDecryptCause = errors.NewCause(errors.BadRequestCategory, "secret_box_decrypt")

	InvalidKeyURLCause = errors.NewCause(errors.BadRequestCategory, "invalid_key_url")
)
