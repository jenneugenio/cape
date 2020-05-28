package client

import errors "github.com/capeprivacy/cape/partyerrors"

var (
	SerializationCause = errors.NewCause(errors.NotImplementedCategory, "serialization_error")
)
