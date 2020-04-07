package auth

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	"github.com/dropoutlabs/cape/primitives"
	"github.com/manifoldco/go-base64"
)

// NewFakeIdentity returns fake identity data for when an attacker
// attempts to enumerate over emails to see who has an
// account with the coordinator
func NewFakeIdentity(email primitives.Email) (primitives.Identity, error) {
	h := sha256.New()
	_, err := h.Write([]byte(email.String()))
	if err != nil {
		return nil, err
	}

	salt := make([]byte, SaltLength)

	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
	rand.Seed(int64(seed))
	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}

	// It doesn't matter whether this is a user or a service
	user, err := primitives.NewUser("", primitives.Email{Email: ""}, &primitives.Credentials{
		Salt: base64.New(salt),
		Alg:  primitives.EDDSA,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
