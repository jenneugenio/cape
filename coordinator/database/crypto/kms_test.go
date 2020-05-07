package crypto

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

func TestKMSEncryptDecryptDEK(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	var key [32]byte

	_, err := rand.Read(key[:])
	gm.Expect(err).To(gm.BeNil())

	encodedKey := base64.New(key[:]).String()

	u, err := NewKeyURL("base64key://" + encodedKey)
	gm.Expect(err).To(gm.BeNil())

	kms, err := NewLocalKMS(u)
	gm.Expect(err).To(gm.BeNil())

	err = kms.Open(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer kms.Close()

	var dek [32]byte
	_, err = rand.Read(dek[:])
	gm.Expect(err).To(gm.BeNil())

	wrappedDEK, err := kms.Encrypt(ctx, dek[:])
	gm.Expect(err).To(gm.BeNil())

	newDEK, err := kms.Decrypt(ctx, wrappedDEK)
	gm.Expect(err).To(gm.BeNil())

	gm.Expect(dek[:]).To(gm.Equal(newDEK))
}
