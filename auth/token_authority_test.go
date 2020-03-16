package auth

import (
	"testing"

	gm "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2"
)

func TestTokenAuthorityGenerate(t *testing.T) {
	gm.RegisterTestingT(t)

	tokenAuth, err := NewTokenAuthority("controller@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	sig, err := tokenAuth.Generate(Login)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(sig).ToNot(gm.BeNil())
}

func TestTokenAuthorityVerify(t *testing.T) {
	gm.RegisterTestingT(t)

	tokenAuth, err := NewTokenAuthority("controller@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	sig, err := tokenAuth.Generate(Login)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(sig).ToNot(gm.BeNil())

	err = tokenAuth.Verify(sig)
	gm.Expect(err).To(gm.BeNil())
}

func TestTokenAuthorityVerifyError(t *testing.T) {
	gm.RegisterTestingT(t)

	tokenAuth, err := NewTokenAuthority("controller@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	otherTokenAuth, err := NewTokenAuthority("controller2@controller.ai")
	gm.Expect(err).To(gm.BeNil())

	sig, err := tokenAuth.Generate(Login)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(sig).ToNot(gm.BeNil())

	err = otherTokenAuth.Verify(sig)
	gm.Expect(err).To(gm.Equal(jose.ErrCryptoFailure))
}
