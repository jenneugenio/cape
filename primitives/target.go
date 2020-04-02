package primitives

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
	"regexp"
)

// Target of a policy
type Target string

// only records are supported right now
var targetRegex = regexp.MustCompile(`^records:[a-z]+\..*$`)

// Validate that target is valid
func (t Target) Validate() error {
	if !targetRegex.MatchString(string(t)) {
		msg := "Target must be in the form <type>:<collection>.<collection>"
		return errors.New(InvalidTargetCause, msg)
	}

	return nil
}

// NewTarget validates that the target is valid label before returning it
func NewTarget(in string) (Target, error) {
	t := Target(in)
	err := t.Validate()
	return t, err
}
