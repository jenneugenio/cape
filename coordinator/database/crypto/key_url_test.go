package crypto

import (
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

func TestKeyURL(t *testing.T) {
	g := gm.NewWithT(t)

	t.Run("new base64 key, generated", func(t *testing.T) {
		u, err := NewBase64KeyURL(nil)
		g.Expect(err).To(gm.BeNil())

		g.Expect(len(u.Host)).To(gm.Equal(43))
		g.Expect(u.Scheme).To(gm.Equal("base64key"))
	})

	t.Run("new base64 key, not generated", func(t *testing.T) {
		bytes := make([]byte, KeyLength)
		u, err := NewBase64KeyURL(bytes)
		g.Expect(err).To(gm.BeNil())

		g.Expect(len(u.Host)).To(gm.Equal(43))
		g.Expect(u.Scheme).To(gm.Equal("base64key"))

		v, err := base64.NewFromString(u.Host)
		g.Expect(err).To(gm.BeNil())

		g.Expect([]byte(*v)).To(gm.Equal(bytes))
	})

	t.Run("incorrect length", func(t *testing.T) {
		bytes := make([]byte, 31)
		_, err := NewBase64KeyURL(bytes)
		g.Expect(err).NotTo(gm.BeNil())
	})
}
