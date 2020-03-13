package database

import (
	"time"

	"github.com/dropoutlabs/cape/database/types"
)

// Entity represents any primitive data structure stored inside the Controller.
// All primitives must satisfy this interface to be stored in the ata layer.
type Entity interface {
	GetID() ID
	GetType() types.Type
	GetVersion() uint8
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time) error
}