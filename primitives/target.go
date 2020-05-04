package primitives

import (
	"regexp"

	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Entity is in the form <target>:<collection>.<entity>
// However, the target can also specify wildcards such as <target>:<collection>.*
// or <target>:*
//
// If the fully specified regex fails, we will try against the other wildcard regex
var fullySpecifiedRegex = regexp.MustCompile(`^(.*):(.*)\.(.*)+$`)
var collectionWildcardRegex = regexp.MustCompile(`^(.*):(\*)$`)

const (
	// These are the indices for the string slice returned
	// by the above regular expressions
	TypeIndex       = 1
	CollectionIndex = 2
	EntityIndex     = 3
)

// Collection for this target
type Collection string

// String returns the string form of the collection
func (c Collection) String() string {
	return string(c)
}

// Entity for a collection
type Entity string

// String returns the string form of the entity
func (e Entity) String() string {
	return string(e)
}

// TargetType is the record type this target points at (e.g. records)
type TargetType string

const (
	Records TargetType = "records"
)

func (t TargetType) Validate() error {
	if t == Records {
		return nil
	}

	_, ok := types.Get(t.String())
	if !ok {
		return errors.New(InvalidTargetCause, "Type %s is not a valid target must be records or primitive type", t)
	}

	return nil
}

func (t TargetType) String() string {
	return string(t)
}

// Entity of a policy
type Target string

// Validate that target is valid
func (t Target) Validate() error {
	if !fullySpecifiedRegex.MatchString(string(t)) && !collectionWildcardRegex.MatchString(string(t)) {
		msg := "Target must be in the form <type>:<collection>.<entity>"
		return errors.New(InvalidTargetCause, msg)
	}

	return t.Type().Validate()
}

// Checks if this target and the provided target match. This supports wildcards
func (t Target) Matches(other Target) bool {
	return t.Entity() == other.Entity()
}

// Type returns what type this is targeting
func (t Target) Type() TargetType {
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return TargetType(res[TypeIndex])
	}

	res = collectionWildcardRegex.FindStringSubmatch(t.String())
	return TargetType(res[TypeIndex])
}

// Collection returns which collection this target refers to
func (t Target) Collection() Collection {
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return Collection(res[CollectionIndex])
	}

	res = collectionWildcardRegex.FindStringSubmatch(t.String())
	return Collection(res[CollectionIndex])
}

// Entity returns which entity this target refers to
func (t Target) Entity() Entity {
	// if the collection was wildcarded, then this won't match
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return Entity(res[EntityIndex])
	}

	return "*"
}

// String turns the Entity into a string
func (t Target) String() string {
	return string(t)
}

// NewTarget validates that the target is valid label before returning it
func NewTarget(in string) (Target, error) {
	t := Target(in)
	return t, t.Validate()
}
