package worker

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	MissingEnvCause     = errors.NewCause(errors.BadRequestCategory, "missing_env")
	BadCertificateCause = errors.NewCause(errors.BadRequestCategory, "invalid_cert")
)
