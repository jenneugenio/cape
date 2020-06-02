package primitives

import (
	errors "github.com/capeprivacy/cape/partyerrors"
	"regexp"
)

var displayNameRegex = regexp.MustCompile(`^[\p{L}\s-_\d]{3,64}$`)

type DisplayName string

func (p DisplayName) Validate() error {
	if !displayNameRegex.MatchString(string(p)) {
		return errors.New(InvalidProjectNameCause, "%s contains invalid characters (unicode only)", p)
	}

	return nil
}

func NewDisplayName(in string) (DisplayName, error) {
	p := DisplayName(in)
	return p, p.Validate()
}
