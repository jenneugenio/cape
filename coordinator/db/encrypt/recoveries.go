package encrypt

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/models"
)

type recoveriesEncrypt struct {
	db    db.RecoveryDB
	codec crypto.EncryptionCodec
}

var _ db.RecoveryDB = &recoveriesEncrypt{}

func (r *recoveriesEncrypt) Get(ctx context.Context, ID string) (*models.Recovery, error) {
	encRecovery, err := r.db.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	dec, err := r.codec.Decrypt(ctx, encRecovery.Credentials.Secret)
	if err != nil {
		return nil, err
	}

	recovery := encRecovery
	recovery.Credentials.Secret = dec
	return recovery, nil
}

func (r *recoveriesEncrypt) Create(ctx context.Context, recovery models.Recovery) error {
	enc, err := r.codec.Encrypt(ctx, recovery.Credentials.Secret)
	if err != nil {
		return err
	}

	encRecovery := recovery
	encRecovery.Credentials.Secret = enc

	return r.db.Create(ctx, encRecovery)
}

func (r *recoveriesEncrypt) Delete(ctx context.Context, ID string) error {
	return r.db.Delete(ctx, ID)
}
