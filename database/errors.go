package database

import errors "github.com/dropoutlabs/privacyai/partyerrors"

var (
	// NotImplementedDBCause happens when you try to connect to a db we do not support
	NotImplementedDBCause = errors.NewCause(errors.NotImplementedCategory, "not_implemented")

	// NotFoundCause happens when the entity you were trying to operate on was not found
	NotFoundCause = errors.NewCause(errors.NotFoundCategory, "not_found")

	// NotMutableCause happens when an entity is immutable but an operation
	// attempting to update it is performed.
	NotMutableCause = errors.NewCause(errors.BadRequestCategory, "not_mutable")

	// TypeMismatchCause happens when an ID is provided that doesn't match the
	// entity type
	TypeMismatchCause = errors.NewCause(errors.BadRequestCategory, "wrong_id_type")

	// DuplicateCause happens when an entity already exists due to some
	// constraint that exists in the database
	DuplicateCause = errors.NewCause(errors.ConflictCategory, "entity_already_exists")
)
