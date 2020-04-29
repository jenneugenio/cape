package database

import (
	"time"

	"github.com/capeprivacy/cape/coordinator/database/types"
)

// Entity represents any primitive data structure stored inside the Coordinator.
// All primitives must satisfy this interface to be stored in the ata layer.
type Entity interface {
	GetID() ID
	GetType() types.Type
	GetVersion() uint8
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time) error
}
