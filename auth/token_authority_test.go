package auth

import (
	"testing"
	"time"

	gm "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2"

	"github.com/dropoutlabs/cape/primitives"
)

func TestTokenAuthorityGenerate(t *testing.T) {
	gm.RegisterTestingT(t)

	tokenAuth, err := NewTokenAuthority("controller@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(sig).ToNot(gm.BeNil())
	gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))
}

func TestTokenAuthorityVerify(t *testing.T) {
	gm.RegisterTestingT(t)

	tokenAuth, err := NewTokenAuthority("controller@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(sig).ToNot(gm.BeNil())
	gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

	err = tokenAuth.Verify(sig)
	gm.Expect(err).To(gm.BeNil())
}

func TestTokenAuthorityVerifyError(t *testing.T) {
	gm.RegisterTestingT(t)

	tokenAuth, err := NewTokenAuthority("controller@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	otherTokenAuth, err := NewTokenAuthority("controller2@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	sig, expiresIn, err := tokenAuth.Generate(primitives.Login)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(sig).ToNot(gm.BeNil())
	gm.Expect(expiresIn).To(gm.BeTemporally(">", time.Now().UTC()))

	err = otherTokenAuth.Verify(sig)
	gm.Expect(err).To(gm.Equal(jose.ErrCryptoFailure))
}
