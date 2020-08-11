package encrypt

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"
)

var _ db.SecretDB = &secretEncrypt{}

type secretEncrypt struct {
	db    db.SecretDB
	codec crypto.EncryptionCodec
}

func (p *secretEncrypt) Create(ctx context.Context, secret models.SecretArg) error {
	// generate random bytes for secret value
	b := make([]byte, auth.SecretLength)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	secret.Value = base64.New(b)

	s, err := encryptSecret(ctx, p.codec, secret)
	if err != nil {
		return fmt.Errorf("error encrypting user for creation: %w", err)
	}

	return p.db.Create(ctx, *s)
}

func (p *secretEncrypt) Delete(ctx context.Context, name string) (db.DeleteStatus, error) {
	return p.db.Delete(ctx, name)
}

func (p *secretEncrypt) Get(ctx context.Context, name string) (*models.SecretArg, error) {
	secret, err := p.db.Get(ctx, name)

	if err != nil {
		return nil, err
	}

	return decryptSecret(ctx, p.codec, *secret)
}

func encryptSecret(ctx context.Context, codec crypto.EncryptionCodec, secret models.SecretArg) (*models.SecretArg, error) {
	enc, err := codec.Encrypt(ctx, secret.Value)
	if err != nil {
		return nil, err
	}

	s := secret

	s.Value = enc

	return &s, nil
}

func decryptSecret(ctx context.Context, codec crypto.EncryptionCodec, secret models.SecretArg) (*models.SecretArg, error) {
	dec, err := codec.Decrypt(ctx, secret.Value)
	if err != nil {
		return nil, err
	}

	s := secret

	s.Value = dec

	return &s, nil
}
