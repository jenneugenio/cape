package primitives

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
	"regexp"
)

// Target is in the form <target>:<collection>.<entity>

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

// Target of a policy
type Target string

// only records are supported right now
var targetRegex = regexp.MustCompile(`^records:(.*)\.(.*)+$`)

// Validate that target is valid
func (t Target) Validate() error {
	if !targetRegex.MatchString(string(t)) {
		msg := "Target must be in the form <type>:<collection>.<entity>"
		return errors.New(InvalidTargetCause, msg)
	}

	return nil
}

// Type returns what type this is targeting
func (t Target) Type() TargetType {
	return Records
}

// Collection returns which collection this target refers to
func (t Target) Collection() Collection {
	res := targetRegex.FindStringSubmatch(t.String())
	return Collection(res[1])
}

// Entity returns which entity this target refers to
func (t Target) Entity() Entity {
	res := targetRegex.FindStringSubmatch(t.String())
	return Entity(res[2])
}

// String turns the Target into a string
func (t Target) String() string {
	return string(t)
}

// NewTarget validates that the target is valid label before returning it
func NewTarget(in string) (Target, error) {
	t := Target(in)
	err := t.Validate()
	return t, err
}
