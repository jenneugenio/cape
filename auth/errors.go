package auth

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	RequiredTokenCause = errors.NewCause(errors.BadRequestCategory, "required_token")
	RequiredSigCause   = errors.NewCause(errors.BadRequestCategory, "required_signature")

	RequiredPrivateKeyCause = errors.NewCause(errors.BadRequestCategory, "required_private_key")
	RequiredPublicKeyCause  = errors.NewCause(errors.BadRequestCategory, "required_public_key")

	SignatureNotValid = errors.NewCause(errors.UnauthorizedCategory, "signature_not_valid")

	BadAPITokenVersion = errors.NewCause(errors.BadRequestCategory, "bad_apitoken_version")
	BadSaltLength      = errors.NewCause(errors.BadRequestCategory, "bad_salt_length")
	BadSecretLength    = errors.NewCause(errors.BadRequestCategory, "bad_secret_length")
	BadPackagedKeypair = errors.NewCause(errors.BadRequestCategory, "bad_packaged_keypair")
	BadAlgType         = errors.NewCause(errors.BadRequestCategory, "bad_alg_type")
	BadPublicKeyLength = errors.NewCause(errors.BadRequestCategory, "bad_public_key_length")

	// BadTokenFormat happens when an APIToken has a bad format
	BadTokenFormat = errors.NewCause(errors.BadRequestCategory, "bad_token_format")

	// InvalidAuthHeader occurs when the auth header is in the wrong format
	InvalidAuthHeader = errors.NewCause(errors.BadRequestCategory, "invalid_auth_header")

	InvalidInfo = errors.NewCause(errors.BadRequestCategory, "invalid_auth_info")

	// AuthenticationFailure is caused by authentication failing
	AuthenticationFailure = errors.NewCause(errors.UnauthorizedCategory, "authentication_failure")

	// AuthorizationFailure is caused by authorization not being given
	AuthorizationFailure = errors.NewCause(errors.UnauthorizedCategory, "authorization_failure")

	// ErrAuthentication is the error wrapping the AuthenticationFailure cause
	ErrAuthentication = errors.New(AuthenticationFailure, "Failed to authenticate")

	// ErrAuthorization is the error wrapping the AuthorizationFailure cause
	ErrAuthorization = errors.New(AuthorizationFailure, "Access denied")
)
