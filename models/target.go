package models

import (
	"errors"
	"regexp"
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
	typeIndex       = 1
	collectionIndex = 2
	entityIndex     = 3
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

	return nil
}

func (t TargetType) String() string {
	return string(t)
}

// Target of a rule
type Target string

// Validate that target is valid
func (t Target) Validate() error {
	if !fullySpecifiedRegex.MatchString(string(t)) && !collectionWildcardRegex.MatchString(string(t)) {
		return errors.New("Target must be in the form <type>:<collection>.<entity>")
	}

	return t.Type().Validate()
}

// Matches checks if this target and the provided target match. This supports wildcards
func (t Target) Matches(other Target) bool {
	return t.Entity() == other.Entity()
}

// Type returns what type this is targeting
func (t Target) Type() TargetType {
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return TargetType(res[typeIndex])
	}

	res = collectionWildcardRegex.FindStringSubmatch(t.String())
	return TargetType(res[typeIndex])
}

// Collection returns which collection this target refers to
func (t Target) Collection() Collection {
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return Collection(res[collectionIndex])
	}

	res = collectionWildcardRegex.FindStringSubmatch(t.String())
	return Collection(res[collectionIndex])
}

// Entity returns which entity this target refers to
func (t Target) Entity() Entity {
	// if the collection was wildcarded, then this won't match
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return Entity(res[entityIndex])
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
