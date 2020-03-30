package auth

import (
	"encoding/json"
	"testing"

	gm "github.com/onsi/gomega"
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

		pkg, err := kp.Package()
		gm.Expect(err).To(gm.BeNil())

		unpkg, err := pkg.Unpackage()
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(unpkg).To(gm.Equal(kp))
	})

	t.Run("can json serialize and unserialze to same kp", func(t *testing.T) {
		kp, err := NewKeypair()
		gm.Expect(err).To(gm.BeNil())

		pkg, err := kp.Package()
		gm.Expect(err).To(gm.BeNil())

		out, err := json.Marshal(pkg)
		gm.Expect(err).To(gm.BeNil())

		in := &KeypairPackage{}
		err = json.Unmarshal(out, in)
		gm.Expect(err).To(gm.BeNil())

		new, err := in.Unpackage()
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(new).To(gm.Equal(kp))
	})
}
