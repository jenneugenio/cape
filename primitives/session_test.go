package primitives

import (
	"testing"
	"time"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

func TestNewSession(t *testing.T) {
	gm.RegisterTestingT(t)

	email, err := NewEmail("bob@bob.com")
	gm.Expect(err).To(gm.BeNil())

	user, err := NewUser("bob", email, &Credentials{})
	gm.Expect(err).To(gm.BeNil())

	ti := time.Now()
	token := base64.New([]byte("random-string"))
	session, err := NewSession(user, ti, Login, token)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(session.GetType()).To(gm.Equal(SessionType))
	gm.Expect(session.Credentials).ToNot(gm.BeNil())
	gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
	gm.Expect(session.Token).To(gm.Equal(token))
	gm.Expect(session.IdentityID).To(gm.Equal(user.ID))

	session, err = NewSession(user, ti, Authenticated, token)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(session.GetType()).To(gm.Equal(SessionType))
	gm.Expect(session.Credentials).To(gm.BeNil())
	gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
	gm.Expect(session.Token).To(gm.Equal(token))
	gm.Expect(session.IdentityID).To(gm.Equal(user.ID))
}
