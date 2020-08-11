package encrypt

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/models"
)

type sessionEncrypt struct {
	db    db.SessionDB
	codec crypto.EncryptionCodec
}

var _ db.SessionDB = &sessionEncrypt{}

func (s *sessionEncrypt) Get(ctx context.Context, ID string) (*models.Session, error) {
	encSession, err := s.db.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	dec, err := s.codec.Decrypt(ctx, encSession.Token)
	if err != nil {
		return nil, err
	}

	session := encSession
	session.Token = dec
	return session, nil
}

func (s *sessionEncrypt) Create(ctx context.Context, session models.Session) error {
	enc, err := s.codec.Encrypt(ctx, session.Token)
	if err != nil {
		return err
	}

	encSession := session
	encSession.Token = enc

	return s.db.Create(ctx, encSession)
}

func (s *sessionEncrypt) Delete(ctx context.Context, ID string) error {
	return s.db.Delete(ctx, ID)
}
