package encrypt

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/models"
)

type tokensEncrypt struct {
	db    db.TokensDB
	codec crypto.EncryptionCodec
}

var _ db.TokensDB = &tokensEncrypt{}

func (t *tokensEncrypt) Get(ctx context.Context, ID string) (*models.Token, error) {
	encToken, err := t.db.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	dec, err := t.codec.Decrypt(ctx, encToken.Credentials.Secret)
	if err != nil {
		return nil, err
	}

	token := encToken
	token.Credentials.Secret = dec
	return token, nil
}

func (t *tokensEncrypt) Create(ctx context.Context, token models.Token) error {
	enc, err := t.codec.Encrypt(ctx, token.Credentials.Secret)
	if err != nil {
		return err
	}

	encToken := token
	encToken.Credentials.Secret = enc

	return t.db.Create(ctx, encToken)
}

func (t *tokensEncrypt) Delete(ctx context.Context, ID string) error {
	return t.db.Delete(ctx, ID)
}

func (t *tokensEncrypt) ListByUserID(ctx context.Context, userID string) ([]models.Token, error) {
	return t.db.ListByUserID(ctx, userID)
}
