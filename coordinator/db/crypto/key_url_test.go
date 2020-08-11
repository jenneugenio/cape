package crypto

import (
	"strconv"
	"testing"

	"net/url"

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

	t.Run("new azure key url", func(t *testing.T) {
		key, err := NewKeyURL("azurekeyvault://mykey.com/mykey")
		g.Expect(err).To(gm.BeNil())
		g.Expect(key).ToNot(gm.BeNil())
	})

	t.Run("invalid scheme", func(t *testing.T) {
		_, err := NewKeyURL("http://haha.com")
		g.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("key url from url", func(t *testing.T) {
		u, _ := url.Parse("azurekeyvault://mykey.com/mykey")

		key, err := KeyURLFromURL(u)
		g.Expect(err).To(gm.BeNil())
		g.Expect(key).ToNot(gm.BeNil())
	})

	t.Run("unmarhsal json", func(t *testing.T) {
		str := "azurekeyvault://mykey.com/mykey"
		key := &KeyURL{}
		err := key.UnmarshalJSON([]byte(strconv.Quote(str)))
		g.Expect(err).To(gm.BeNil())

		g.Expect(key.String()).To(gm.Equal(str))
	})
}
