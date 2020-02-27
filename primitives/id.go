package primitives

import (
	"crypto/rand"
	"encoding/json"
	"errors"

	"github.com/manifoldco/go-base32"
	"golang.org/x/crypto/blake2b"

	"github.com/dropoutlabs/privacyai/primitives/types"
)

// ErrInvalidID occurs when an invalid id (bad encoding or overflow) is provided
var ErrInvalidID = errors.New("Invalid ID Provided")

// ErrNotMutable occurs when a ID is attempted to be generated instead of
// derived for an Entity
var ErrNotMutable = errors.New("Must be mutable entity to generate ID; did you mean to derive?")

const (
	idVersion  = 0x01
	byteLength = 16
)

// ID represents a container for an Entities content addressable ID
type ID [byteLength]byte

// EmptyID is a comparator for checking whether or not the given ID is empty
var EmptyID = ID{}

// DeriveID returns a content-addressable ID for the given Entity
func DeriveID(e Entity) (ID, error) {
	h, err := blake2b.New(byteLength-2, nil)
	if err != nil {
		return EmptyID, err
	}

	b, err := json.Marshal(e)
	if err != nil {
		return EmptyID, err
	}
	_, err = h.Write(b)
	if err != nil {
		return EmptyID, err
	}

	t := e.GetType()
	id := ID{idVersion<<4 | t.Upper(), t.Lower()}
	copy(id[2:], h.Sum(nil))

	return id, nil
}

// GenerateID returns an ID for a type representing a mutable Entity (e.g. not
// content-addressable)
func GenerateID(t types.Type) (ID, error) {
	if !t.Mutable() {
		return ID{}, ErrNotMutable
	}

	b := make([]byte, byteLength-2)
	_, err := rand.Read(b)
	if err != nil {
		return ID{}, err
	}

	id := ID{idVersion<<4 | t.Upper(), t.Lower()}
	copy(id[2:], b)

	return id, nil
}

// DecodeFromString returns an ID encoded in the provided string
func DecodeFromString(value string) (ID, error) {
	return DecodeFromBytes([]byte(value))
}

// DecodeFromBytes returns an ID encoded in the given byte slice
func DecodeFromBytes(b []byte) (ID, error) {
	id := ID{}
	err := id.fill(b)
	if err != nil {
		return EmptyID, err
	}

	return id, nil
}

// Validate returns an error if the ID is not a valid ID
func (id ID) Validate(_ interface{}) error {
	if id == EmptyID {
		return ErrInvalidID
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for IDs.
//
// IDs are encoded in unpadded base32.
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for IDs.
func (id *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return errors.New("value is not a string")
	}

	return id.fill(b[1 : len(b)-1])
}

func (id ID) String() string {
	ret := base32.EncodeToString(id[:])

	return ret
}

// Version returns the version of the ID type
func (id ID) Version() byte {
	return id[0] & 0xF0 >> 4
}

// Type returns the underlying type
func (id ID) Type() (types.Type, error) {
	return types.DecodeBytes(id[0]&0x0F, id[1])
}

func (id *ID) fill(in []byte) error {
	out, err := decodeFromBytes(in)
	if err != nil {
		return err
	}

	copy(id[:], out)
	return nil
}

func decodeFromBytes(raw []byte) ([]byte, error) {
	out, err := base32.DecodeString(string(raw))
	if err != nil {
		return nil, err
	}

	if len(out) != byteLength {
		return nil, errors.New("Incorrect length for id")
	}

	return out, nil
}
