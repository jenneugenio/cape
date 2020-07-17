package encrypt

import (
	"context"
	"fmt"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
)

var _ db.UserDB = &encryptUser{}

type encryptUser struct {
	db    db.UserDB
	codec crypto.EncryptionCodec
}

func (e *encryptUser) Create(ctx context.Context, user models.User) error {
	u, err := userEncrypt(ctx, e.codec, user)
	if err != nil {
		return fmt.Errorf("error encrypting user for creation: %w", err)
	}

	err = e.db.Create(ctx, *u)
	if err != nil {
		return err
	}

	return nil
}

func (e *encryptUser) Update(ctx context.Context, id string, user models.User) error {
	u, err := userEncrypt(ctx, e.codec, user)
	if err != nil {
		return err
	}

	err = e.db.Update(ctx, id, *u)
	if err != nil {
		return err
	}

	return nil
}

func (e *encryptUser) Delete(ctx context.Context, email models.Email) (db.DeleteStatus, error) {
	status, err := e.db.Delete(ctx, email)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (e *encryptUser) Get(ctx context.Context, email models.Email) (*models.User, error) {
	user, err := e.db.Get(ctx, email)

	if err != nil {
		return nil, err
	}

	return userDecrypt(ctx, e.codec, *user)
}

func (e *encryptUser) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := e.db.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return userDecrypt(ctx, e.codec, *user)
}

func (e *encryptUser) List(ctx context.Context, opts *db.ListUserOptions) ([]models.User, error) {
	users, err := e.db.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	for i, user := range users {
		u, err := e.codec.Decrypt(ctx, user.Credentials.Secret)
		if err != nil {
			return nil, err
		}

		users[i].Credentials.Secret = u
	}

	return users, nil
}

func userEncrypt(ctx context.Context, codec crypto.EncryptionCodec, user models.User) (*models.User, error) {
	enc, err := codec.Encrypt(ctx, user.Credentials.Secret)
	if err != nil {
		return nil, err
	}

	u := user

	u.Credentials.Secret = enc

	return &u, nil
}

func userDecrypt(ctx context.Context, codec crypto.EncryptionCodec, user models.User) (*models.User, error) {
	dec, err := codec.Decrypt(ctx, user.Credentials.Secret)
	if err != nil {
		return nil, err
	}

	u := user

	u.Credentials.Secret = dec

	return &u, nil
}
