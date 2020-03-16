package primitives

import (
	"testing"
	"time"

	"github.com/dropoutlabs/cape/auth"
	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

func TestNewSession(t *testing.T) {
	gm.RegisterTestingT(t)

	user, err := NewUser("bob", "bob@bob.com", &auth.Credentials{})
	gm.Expect(err).To(gm.BeNil())

	ti := time.Now()
	token := base64.New([]byte("random-string"))
	session, err := NewSession(user, ti, auth.Login, token)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(session.GetType()).To(gm.Equal(SessionType))
	gm.Expect(session.AuthCredentials).ToNot(gm.BeNil())
	gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
	gm.Expect(session.Token).To(gm.Equal(token))
	gm.Expect(session.IdentityID).To(gm.Equal(user.ID))

	session, err = NewSession(user, ti, auth.Authenticated, token)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(session.GetType()).To(gm.Equal(SessionType))
	gm.Expect(session.AuthCredentials).To(gm.BeNil())
	gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
	gm.Expect(session.Token).To(gm.Equal(token))
	gm.Expect(session.IdentityID).To(gm.Equal(user.ID))
}
