package auth

import (
	"encoding/json"
	"github.com/capeprivacy/cape/models"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestKeypair(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("can generate a keypair", func(t *testing.T) {
		_, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can package and unpackage a keypair", func(t *testing.T) {
		kp, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		pkg := kp.Package()
		unpkg, err := pkg.Unpackage()
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(unpkg).To(gm.Equal(kp))
	})

	t.Run("can json serialize and unserialze to same kp", func(t *testing.T) {
		kp, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		pkg := kp.Package()
		out, err := json.Marshal(pkg)
		gm.Expect(err).To(gm.BeNil())

		in := &KeypairPackage{}
		err = json.Unmarshal(out, in)
		gm.Expect(err).To(gm.BeNil())

		new, err := in.Unpackage()
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(new).To(gm.Equal(kp))
	})

	t.Run("derivekeypair", func(t *testing.T) {
		kp, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		tests := []struct {
			name   string
			secret []byte
			salt   []byte
			cause  *errors.Cause
		}{
			{
				name:   "valid",
				secret: kp.secret,
				salt:   kp.salt,
			},
			{
				name:   "bad secret length",
				secret: kp.secret[:3],
				salt:   kp.salt,
				cause:  &BadSecretLength,
			},
			{
				name:   "bad salt length",
				secret: kp.secret,
				salt:   kp.salt[:3],
				cause:  &BadSaltLength,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				kp, err := DeriveKeypair(tc.secret, tc.salt)
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					gm.Expect(kp).To(gm.BeNil())
					return
				}

				gm.Expect(len(kp.salt)).To(gm.Equal(models.SaltLength))
				gm.Expect(len(kp.secret)).To(gm.Equal(SecretLength))
				gm.Expect(len(kp.PrivateKey)).To(gm.Equal(64))
				gm.Expect(len(kp.PublicKey)).To(gm.Equal(32))
				gm.Expect(kp.Alg).To(gm.Equal(models.EDDSA))
			})
		}
	})

	t.Run("keypair package validation", func(t *testing.T) {
		kp, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		tests := []struct {
			name   string
			secret []byte
			salt   []byte
			alg    models.CredentialsAlgType
			cause  *errors.Cause
		}{
			{
				name:   "valid keypair",
				secret: kp.secret,
				salt:   kp.salt,
				alg:    kp.Alg,
			},
			{
				name:   "missing salt",
				secret: kp.secret,
				alg:    kp.Alg,
				cause:  &BadSaltLength,
			},
			{
				name:  "missing secret",
				salt:  kp.salt,
				alg:   kp.Alg,
				cause: &BadSecretLength,
			},
			{
				name:   "wrong salt length",
				salt:   kp.salt[:4],
				secret: kp.secret,
				alg:    kp.Alg,
				cause:  &BadSaltLength,
			},
			{
				name:   "wrong secret length",
				salt:   kp.salt,
				secret: kp.secret[:2],
				alg:    kp.Alg,
				cause:  &BadSecretLength,
			},
			{
				name:   "wrong alg",
				salt:   kp.salt,
				secret: kp.secret,
				alg:    models.SHA256,
				cause:  &BadAlgType,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				pkg := &KeypairPackage{
					Secret: base64.New(tc.secret),
					Salt:   base64.New(tc.salt),
					Alg:    tc.alg,
				}

				err := pkg.Validate()
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
				}
			})
		}
	})

	t.Run("(un)marshal keypair", func(t *testing.T) {
		kp, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		by, err := kp.MarshalJSON()
		gm.Expect(err).To(gm.BeNil())

		newKp := &Keypair{}
		err = newKp.UnmarshalJSON(by)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(newKp).To(gm.Equal(kp))
	})
}
