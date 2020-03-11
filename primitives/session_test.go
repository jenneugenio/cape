package primitives

import (
	"testing"
	"time"

	"github.com/dropoutlabs/cape/auth"
	gm "github.com/onsi/gomega"
)

func TestNewSession(t *testing.T) {
	gm.RegisterTestingT(t)

	user, err := NewUser("bob", "bob@bob.com", &auth.Credentials{})
	gm.Expect(err).To(gm.BeNil())

	session, err := NewSession(user, time.Now(), Login, nil)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(session.GetType()).To(gm.Equal(SessionType))
	gm.Expect(session.AuthCredentials).ToNot(gm.BeNil())

	session, err = NewSession(user, time.Now(), Authenticated, nil)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(err).To(gm.BeNil())
}
