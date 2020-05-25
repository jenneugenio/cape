package primitives

import (
	"crypto/rand"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// MinPasswordLength represents the minimum length of a Cape password
const MinPasswordLength = 8

// MaxPasswordLength represents the maximum length of a Cape password
const MaxPasswordLength = 32

// PasswordByteLength represents the number of bytes used to generate a Cape password
const PasswordByteLength = 16

var EmptyPassword = Password("")

// Password represents a password used by a user to log into a cape account.
//
// This primitive is _only_ used by the command line tool as secrets are
// *never* passed over the wire.
type Password string

// Validate returns an error if the given password has an incorrect length.
func (p Password) Validate() error {
	s := p.String()
	if len(s) < MinPasswordLength {
		return errors.New(InvalidPasswordCause, "Passwords must be atleast %d characters long", MinPasswordLength)
	}

	if len(s) > MaxPasswordLength {
		return errors.New(InvalidPasswordCause, "Passwords cannot be more than %d characters long", MaxPasswordLength)
	}

	return nil
}

// String returns the password as a string
func (p Password) String() string {
	return string(p)
}

// Bytes returns the password as a byte array
func (p Password) Bytes() []byte {
	return []byte(p.String())
}

// NewPassword returns a new Password for the given string. If the string isn't
// a valid password an error is returned.
func NewPassword(input string) (Password, error) {
	p := Password(input)
	return p, p.Validate()
}

// GeneratePassword returns a new password using random data sourced from a
// cryptographically strong pseudorandom source.
func GeneratePassword() (Password, error) {
	bytes := make([]byte, PasswordByteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return Password(""), err
	}

	return Password(base64.New(bytes).String()), nil
}
