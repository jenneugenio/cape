package crypto

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestSecretBox(t *testing.T) {
	gm.RegisterTestingT(t)

	var key [32]byte

	_, err := rand.Read(key[:])
	gm.Expect(err).To(gm.BeNil())

	ctx := context.Background()

	encodedKey := base64.New(key[:]).String()
	u, err := NewKeyURL("base64key://" + encodedKey)
	gm.Expect(err).To(gm.BeNil())

	t.Run("test encrypt decrypt", func(t *testing.T) {
		kms, err := NewLocalKMS(u)
		gm.Expect(err).To(gm.BeNil())

		err = kms.Open(ctx)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := NewSecretBoxCodec(kms)

		data := base64.New([]byte("super secret data"))
		encrypted, err := codec.Encrypt(ctx, data)
		gm.Expect(err).To(gm.BeNil())

		decrypted, err := codec.Decrypt(ctx, encrypted)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(decrypted.String()).To(gm.Equal(data.String()))
	})

	t.Run("corrupted wrapped dek", func(t *testing.T) {
		kms, err := NewLocalKMS(u)
		gm.Expect(err).To(gm.BeNil())

		err = kms.Open(ctx)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := NewSecretBoxCodec(kms)

		data := base64.New([]byte("super secret data"))
		encrypted, err := codec.Encrypt(ctx, data)
		gm.Expect(err).To(gm.BeNil())

		corruptedBytes := []byte(*encrypted)
		corruptedBytes[0] = 'f'

		_, err = codec.Decrypt(ctx, base64.New(corruptedBytes))
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("corrupted encrypted payload", func(t *testing.T) {
		kms, err := NewLocalKMS(u)
		gm.Expect(err).To(gm.BeNil())

		err = kms.Open(ctx)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := NewSecretBoxCodec(kms)

		data := base64.New([]byte("super secret data"))
		encrypted, err := codec.Encrypt(ctx, data)
		gm.Expect(err).To(gm.BeNil())

		corruptedBytes := []byte(*encrypted)

		corruptedBytes[EncryptedKeyLength+1] = 'f'

		_, err = codec.Decrypt(ctx, base64.New(corruptedBytes))
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, SecretBoxDecryptCause)).To(gm.BeTrue())
	})
}
