package primitives

import (
	errors "github.com/capeprivacy/cape/partyerrors"
	"regexp"
)

// Entity is in the form <target>:<collection>.<entity>
// However, the target can also specify wildcards such as <target>:<collection>.*
// or <target>:*
//
// If the fully specified regex fails, we will try against the other wildcard regex
var fullySpecifiedRegex = regexp.MustCompile(`^records:(.*)\.(.*)+$`)
var collectionWildcardRegex = regexp.MustCompile(`^records:(\*)$`)

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

// Entity of a policy
type Target string

// Validate that target is valid
func (t Target) Validate() error {
	if !fullySpecifiedRegex.MatchString(string(t)) && !collectionWildcardRegex.MatchString(string(t)) {
		msg := "Target must be in the form <type>:<collection>.<entity>"
		return errors.New(InvalidTargetCause, msg)
	}

	return nil
}

// Checks if this target and the provided target match. This supports wildcards
func (t Target) Matches(other Target) bool {
	return t.Entity() == other.Entity()
}

// Type returns what type this is targeting
func (t Target) Type() TargetType {
	return Records
}

// Collection returns which collection this target refers to
func (t Target) Collection() Collection {
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return Collection(res[1])
	}

	res = collectionWildcardRegex.FindStringSubmatch(t.String())
	return Collection(res[1])
}

// Entity returns which entity this target refers to
func (t Target) Entity() Entity {
	// if the collection was wildcarded, then this won't match
	res := fullySpecifiedRegex.FindStringSubmatch(t.String())
	if res != nil {
		return Entity(res[2])
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
