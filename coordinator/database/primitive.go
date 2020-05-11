package database

import (
	"time"

	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Primitive a low level object that can be inherited from for getting access
// to common values and methods.
type Primitive struct {
	ID        ID        `json:"id"`
	Version   uint8     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Primitive) GetEncryptable() bool {
	return false
}

// GetID satisfies the Entity interface to return an ID
func (p *Primitive) GetID() ID {
	return p.ID
}

// GetVersion satisfies the Entity interface to return a Version.
//
// If a primitive has a non-1 version it should implement this version itself
func (p *Primitive) GetVersion() uint8 {
	return 1
}

// GetCreatedAt returns the CreatedAt value for this struct
func (p *Primitive) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt returns the UpdateAt value for this struct
func (p *Primitive) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// SetUpatedAt sets the UpdateAt value for this struct
func (p *Primitive) SetUpdatedAt(t time.Time) error {
	if t.Before(p.UpdatedAt) {
		return errors.New(InvalidTimeCause, "cannot set time before current UpdatedAt value")
	}

	if t.Before(p.CreatedAt) {
		return errors.New(InvalidTimeCause, "cannot set time before CreatedAt value")
	}

	p.UpdatedAt = t
	return nil
}

// Validate returns whether or not the Primitive is valid
func (p *Primitive) Validate() error {
	if err := p.ID.Validate(); err != nil {
		return err
	}

	if p.Version < 1 {
		return errors.New(InvalidVersionCause, "Version must be greater than zero")
	}

	if p.UpdatedAt.Before(p.CreatedAt) {
		return errors.New(InvalidTimeCause, "UpdatedAt cannot be before CreatedAt")
	}

	return nil
}

// NewPrimitive returns a new primitive entity object
//
// If the type is mutable then this function will also generate an ID. For
// immutable types the caller must manually lock the struct by deriving an ID.
func NewPrimitive(t types.Type) (*Primitive, error) {
	p := &Primitive{
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if t.Mutable() {
		ID, err := GenerateID(t)
		if err != nil {
			return nil, err
		}

		p.ID = ID
	}

	return p, nil
}
