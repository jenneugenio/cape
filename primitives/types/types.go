package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// Type is a list of all registered entity types within the Cape ecosystem.
// They are encoded within 16bits.
type Type uint16

// String returns the human-readable representation of an entity type (string)
// not (hex)
func (t Type) String() string {
	return registry[t].name
}

// Upper returns the first byte of the type
func (t Type) Upper() byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(t))

	return b[0]
}

// Lower returns the second byte of the type
func (t Type) Lower() byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(t))

	return b[1]
}

// Mutable returns whether or not the type is mutable
func (t Type) Mutable() bool {
	return registry[t].mutable
}

// ErrUnknownType occurs when an unknown type is provided
var ErrUnknownType = errors.New("Cannot decode, encountered an unknown type")

// A list of all the different object types that are stored within Cape.  They
// are explicitly given hex values which cannot change (otherwise stored data
// will not be able to be recalled out of database files).
//
// This makes it easier to figure out which types of data are stored inside
// which rows of the database (as it's all binary). Furthermore, it gives us
// the structure to consider the types of objects that exist within our
// ecosystem and their purposes/categories.
//
// To add an entity type consider the category/purpose of the the underlying
// entity. If it doesn't fit within an existing category create a new one. Each
// category should have room for approximately 100 types.
var (
	// Core System Primitives (users, services, etc)
	User       Type = 0x000
	Service    Type = 0x001
	Token      Type = 0x002
	Role       Type = 0x003
	Policy     Type = 0x004
	Attachment Type = 0x005

	// Everything greater than 0x100 is reserved for testing
	Reserved    Type = 0xD00
	Test        Type = 0xD01
	TestMutable Type = 0xD02

	// >12 bits will mess the ID encoding up!!!
	Overflow Type = 0xFFF
)

type item struct {
	name    string
	t       Type
	mutable bool
}

var registry = map[Type]*item{}
var inverse = map[string]*item{}

// Register registers a type into the global type registry. If the type or
// name already exists a panic will be thrown.
func Register(t Type, name string, mutable bool) {
	if t > Overflow {
		panic(fmt.Sprintf("Type %s with name %s has a value greater than the Overflow value", t, name))
	}

	name = strings.ToLower(name)
	item := &item{name: name, t: t, mutable: mutable}

	if _, ok := registry[t]; ok {
		panic(fmt.Sprintf("Type %s with name %s is already registered", t, name))
	}
	if _, ok := inverse[name]; ok {
		panic(fmt.Sprintf("Name %s is already registered for type %s", name, t))
	}

	registry[t] = item
	inverse[name] = item
}

// Get returns an Type for the given
func Get(name string) Type {
	name = strings.ToLower(name)
	i, ok := inverse[name]
	if !ok {
		panic(fmt.Sprintf("Unknown entity type with name: %s", name))
	}

	return i.t
}

// DecodeBytes returns an entity type for the given bytes, if a type hasn't been
// registered for the given value then an error is returned.
func DecodeBytes(upper, lower byte) (Type, error) {
	return Decode(binary.BigEndian.Uint16([]byte{upper, lower}))
}

// Decode returns an entity type for the given uint16, if a type hasn't been
// registered for the given value then an error is returned.
func Decode(in uint16) (Type, error) {
	t := Type(in)
	_, ok := registry[t]
	if !ok {
		return t, ErrUnknownType
	}

	return t, nil
}

func init() {
	Register(User, "user", false)
	Register(Service, "service", false)
	Register(Token, "token", true)
	Register(Role, "role", true)
	Register(Policy, "policy", true)
	Register(Attachment, "attachment", true)
	Register(Test, "test", false)
	Register(TestMutable, "test_mutable", true)
}
