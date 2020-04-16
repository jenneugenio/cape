package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestLabel(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := map[string]struct {
		in    string
		valid bool
		cause errors.Cause
	}{
		"valid": {
			in:    "assdf049",
			valid: true,
		},
		"too short": {
			in:    "sd",
			valid: false,
			cause: InvalidLabelCause,
		},
		"too long": {
			in:    "sdfsfsfsfasfsfdsdfsdfsdfsdfsdfsfsdfsdfsdfsdfsdfsdfsdfdsfsdfsdfdsfsdfsdfsdfdsfsdfsdfsdfsdfsfdsffsd",
			valid: false,
			cause: InvalidLabelCause,
		},
		"starts with a -": {
			in:    "-sfsfsfsfsf",
			valid: false,
			cause: InvalidLabelCause,
		},
		"can start with a number": {
			in:    "0sfasf",
			valid: true,
		},
		"can contain - and numbers": {
			in:    "sdfs-sdf0",
			valid: true,
		},
		"cannot contain capitals": {
			in:    "sdfdfsFDs",
			valid: false,
			cause: InvalidLabelCause,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewLabel(test.in)
			gm.Expect(err == nil).To(gm.Equal(test.valid))
			if test.valid {
				gm.Expect(errors.FromCause(err, test.cause))
			}
		})
	}
}
