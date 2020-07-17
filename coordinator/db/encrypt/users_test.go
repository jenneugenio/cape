package encrypt

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"
)

var ErrGenericDBError = errors.New("generic db error")

var SecretUser = models.User{
	Email: models.Email("hey@email.com"),
	Credentials: &models.Credentials{
		Secret: base64.New([]byte("HEYEYEYEYYE")),
	},
}

func TestUsersCreate(t *testing.T) {
	tests := []struct {
		user    models.User
		wantErr error
		err     error
	}{
		{
			user:    SecretUser,
			wantErr: nil,
			err:     nil,
		},
		{
			user:    SecretUser,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	pgUser := &testPgUser{}
	for i, test := range tests {
		userDB := encryptUser{
			db:    pgUser,
			codec: codec,
		}

		pgUser.err = test.err

		gotErr := userDB.Create(context.TODO(), test.user)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

func TestUsersUpdate(t *testing.T) {
	tests := []struct {
		user    models.User
		id      string
		wantErr error
		err     error
	}{
		{
			user:    SecretUser,
			id:      "idididid",
			wantErr: nil,
			err:     nil,
		},
		{
			user:    SecretUser,
			id:      "idididid",
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	pgUser := &testPgUser{}
	for i, test := range tests {
		userDB := encryptUser{
			db:    pgUser,
			codec: codec,
		}

		pgUser.err = test.err

		gotErr := userDB.Update(context.TODO(), test.id, test.user)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Update() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

func TestUsersDelete(t *testing.T) {
	tests := []struct {
		email   models.Email
		wantErr error
		err     error
	}{
		{
			email:   models.Email("foo"),
			wantErr: nil,
		},
		{
			email:   models.Email("foo"),
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	pgUser := &testPgUser{}
	for i, test := range tests {
		userDB := encryptUser{
			db:    pgUser,
			codec: codec,
		}

		pgUser.err = test.err
		_, gotErr := userDB.Delete(context.TODO(), test.email)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Delete() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
	}
}

func TestUserGet(t *testing.T) {
	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	encryptedSecret, _ := codec.Encrypt(context.TODO(), base64.New([]byte("HEYEYEYEYYE")))
	encryptedUser := models.User{
		Credentials: &models.Credentials{
			Secret: encryptedSecret,
		},
		Email: models.Email("hey@email.com"),
	}

	secretUser := models.User{
		Email: models.Email("hey@email.com"),
		Credentials: &models.Credentials{
			Secret: base64.New([]byte("HEYEYEYEYYE")),
		},
	}

	tests := []struct {
		user     models.User
		wantUser *models.User
		wantErr  error
		err      error
	}{
		{
			user:     encryptedUser,
			wantUser: &secretUser,
			wantErr:  nil,
			err:      nil,
		},
		{
			user:    encryptedUser,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
		{
			user:    encryptedUser,
			wantErr: db.ErrNoRows,
			err:     db.ErrNoRows,
		},
	}

	pgUser := &testPgUser{}
	for i, test := range tests {
		userDB := encryptUser{
			db:    pgUser,
			codec: codec,
		}

		pgUser.user = test.user
		pgUser.err = test.err

		gotUser, gotErr := userDB.Get(context.TODO(), test.user.Email)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Get() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(gotUser, test.wantUser) {
			t.Errorf("incorrect user returned on Get() test %d of %d: got %v want %v", i+1, len(tests), gotUser, test.wantUser)
		}
	}
}

func TestUserGetByID(t *testing.T) {
	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	encryptedSecret, _ := codec.Encrypt(context.TODO(), base64.New([]byte("HEYEYEYEYYE")))
	encryptedUser := models.User{
		ID: "idididid",
		Credentials: &models.Credentials{
			Secret: encryptedSecret,
		},
		Email: models.Email("hey@email.com"),
	}

	secretUser := models.User{
		ID:    "idididid",
		Email: models.Email("hey@email.com"),
		Credentials: &models.Credentials{
			Secret: base64.New([]byte("HEYEYEYEYYE")),
		},
	}

	tests := []struct {
		user     models.User
		wantUser *models.User
		wantErr  error
		err      error
	}{
		{
			user:     encryptedUser,
			wantUser: &secretUser,
			wantErr:  nil,
			err:      nil,
		},
		{
			user:    encryptedUser,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
		{
			user:    encryptedUser,
			wantErr: db.ErrNoRows,
			err:     db.ErrNoRows,
		},
	}

	pgUser := &testPgUser{}
	for i, test := range tests {
		userDB := encryptUser{
			db:    pgUser,
			codec: codec,
		}

		pgUser.user = test.user
		pgUser.err = test.err

		gotUser, gotErr := userDB.GetByID(context.TODO(), test.user.ID)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on GetByID() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(gotUser, test.wantUser) {
			t.Errorf("incorrect user returned on GetByID() test %d of %d: got %v want %v", i+1, len(tests), gotUser, test.wantUser)
		}
	}
}

func TestUserstList(t *testing.T) {
	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	encryptedSecret, _ := codec.Encrypt(context.TODO(), base64.New([]byte("HEYEYEYEYYE")))
	encryptedUser := models.User{
		ID: "idididid",
		Credentials: &models.Credentials{
			Secret: encryptedSecret,
		},
		Email: models.Email("hey@email.com"),
	}

	secretUser := models.User{
		ID:    "idididid",
		Email: models.Email("hey@email.com"),
		Credentials: &models.Credentials{
			Secret: base64.New([]byte("HEYEYEYEYYE")),
		},
	}

	tests := []struct {
		opt *db.ListUserOptions

		wantUsers []models.User
		wantErr   error
		err       error
		user      models.User
	}{
		{
			user:      encryptedUser,
			opt:       nil,
			wantUsers: []models.User{secretUser},
			wantErr:   nil,
			err:       nil,
		},
	}

	pgUser := &testPgUser{}
	for i, test := range tests {
		userDB := encryptUser{
			db:    pgUser,
			codec: codec,
		}

		pgUser.err = test.err
		pgUser.user = test.user
		gotUsers, gotErr := userDB.List(context.TODO(), test.opt)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on List() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(gotUsers, test.wantUsers) {
			t.Errorf("incorrect user returned on List() test %d of %d: got %v want %v", i+1, len(tests), gotUsers, test.wantUsers)
		}
	}
}

type testPgUser struct {
	user models.User
	err  error
}

func (t *testPgUser) Create(_ context.Context, _ models.User) error {
	if t.err != nil {
		return t.err
	}

	return nil
}

func (t *testPgUser) Update(_ context.Context, _ string, _ models.User) error {
	if t.err != nil {
		return t.err
	}

	return nil
}

func (t *testPgUser) Delete(_ context.Context, _ models.Email) (db.DeleteStatus, error) {
	if t.err != nil {
		return db.DeleteStatusError, t.err
	}

	return db.DeleteStatusDeleted, nil
}

func (t *testPgUser) Get(_ context.Context, _ models.Email) (*models.User, error) {
	if t.err != nil {
		return nil, t.err
	}

	return &t.user, nil
}

func (t *testPgUser) GetByID(_ context.Context, _ string) (*models.User, error) {
	if t.err != nil {
		return nil, t.err
	}

	return &t.user, nil
}

func (t *testPgUser) List(_ context.Context, _ *db.ListUserOptions) ([]models.User, error) {
	if t.err != nil {
		return nil, t.err
	}

	return []models.User{t.user}, nil
}
