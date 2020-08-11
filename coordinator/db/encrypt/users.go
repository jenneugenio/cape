package encrypt

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/models"
)

var _ db.UserDB = &userEncrypt{}

type userEncrypt struct {
	db    db.UserDB
	codec crypto.EncryptionCodec
}

func (e *userEncrypt) Create(ctx context.Context, user models.User) error {
	u, err := encryptUser(ctx, e.codec, user)
	if err != nil {
		return fmt.Errorf("error encrypting user for creation: %w", err)
	}

	return e.db.Create(ctx, *u)
}

func (e *userEncrypt) Update(ctx context.Context, id string, user models.User) error {
	u, err := encryptUser(ctx, e.codec, user)
	if err != nil {
		return err
	}

	return e.db.Update(ctx, id, *u)
}

func (e *userEncrypt) Delete(ctx context.Context, email models.Email) (db.DeleteStatus, error) {
	return e.db.Delete(ctx, email)
}

func (e *userEncrypt) Get(ctx context.Context, email models.Email) (*models.User, error) {
	user, err := e.db.Get(ctx, email)

	if err != nil {
		return nil, err
	}

	return decryptUser(ctx, e.codec, *user)
}

func (e *userEncrypt) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := e.db.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return decryptUser(ctx, e.codec, *user)
}

func (e *userEncrypt) List(ctx context.Context, opts *db.ListUserOptions) ([]models.User, error) {
	users, err := e.db.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	for i, user := range users {
		u, err := decryptUser(ctx, e.codec, user)
		if err != nil {
			return nil, err
		}
		users[i] = *u
	}

	return users, nil
}

func encryptUser(ctx context.Context, codec crypto.EncryptionCodec, user models.User) (*models.User, error) {
	enc, err := codec.Encrypt(ctx, user.Credentials.Secret)
	if err != nil {
		return nil, err
	}

	u := user

	u.Credentials.Secret = enc

	return &u, nil
}

func decryptUser(ctx context.Context, codec crypto.EncryptionCodec, user models.User) (*models.User, error) {
	dec, err := codec.Decrypt(ctx, user.Credentials.Secret)
	if err != nil {
		return nil, err
	}

	u := user

	u.Credentials.Secret = dec

	return &u, nil
}
