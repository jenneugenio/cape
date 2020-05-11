package crypto

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

type EncryptableImpl struct {
	NonSecret string `json:"non_secret"`
	Secret    string `json:"secret"`
}

type encrypted struct {
	*EncryptableImpl
	Secret *base64.Value `json:"secret"`
}

func (e *EncryptableImpl) Encrypt(ctx context.Context, codec EncryptionCodec) ([]byte, error) {
	data, err := codec.Encrypt(ctx, base64.New([]byte(e.Secret)))
	if err != nil {
		return nil, err
	}

	return json.Marshal(encrypted{
		EncryptableImpl: &EncryptableImpl{
			NonSecret: e.NonSecret,
		},
		Secret: data,
	})
}

func (e *EncryptableImpl) Decrypt(ctx context.Context, codec EncryptionCodec, data []byte) error {
	in := &encrypted{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Secret)
	if err != nil {
		return err
	}

	e.Secret = string(*unencrypted)
	return nil
}

func TestEncryptable(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	t.Run("test encrypt decrypt", func(t *testing.T) {
		impl := &EncryptableImpl{NonSecret: "not-a-secret", Secret: "super secret secret"}

		var key [32]byte

		_, err := rand.Read(key[:])
		gm.Expect(err).To(gm.BeNil())

		encodedKey := base64.New(key[:]).String()

		url, err := NewKeyURL("base64key://" + encodedKey)
		gm.Expect(err).To(gm.BeNil())

		kms, err := NewLocalKMS(url)
		gm.Expect(err).To(gm.BeNil())

		err = kms.Open(ctx)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := NewSecretBoxCodec(kms)

		nonEncryptedJS, err := json.Marshal(impl)
		gm.Expect(err).To(gm.BeNil())

		encryptedJS, err := impl.Encrypt(ctx, codec)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(nonEncryptedJS).ToNot(gm.Equal(encryptedJS))

		newImpl := &EncryptableImpl{}
		err = newImpl.Decrypt(ctx, codec, encryptedJS)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(newImpl.Secret).To(gm.Equal(impl.Secret))
	})
}
