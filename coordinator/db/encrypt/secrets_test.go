package encrypt

import (
	"context"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"github.com/manifoldco/go-base64"

	gm "github.com/onsi/gomega"
)

var SecretArg = models.SecretArg{
	Type:  "secret",
	Name:  "my-key",
	Value: base64.New([]byte("secretsecret")),
}

func TestSecretCreate(t *testing.T) {
	tests := []struct {
		secret  models.SecretArg
		wantErr error
		err     error
	}{
		{
			secret:  SecretArg,
			wantErr: nil,
			err:     nil,
		},
		{
			secret:  SecretArg,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	pgSecret := &testPgSecret{}
	for i, test := range tests {
		secretDB := secretEncrypt{
			db:    pgSecret,
			codec: codec,
		}

		pgSecret.err = test.err

		gotErr := secretDB.Create(context.TODO(), test.secret)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Create() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}

		if reflect.DeepEqual(SecretArg, pgSecret.receivedSecret) {
			t.Errorf("secret secret not encrypted: got %v %v", SecretArg, pgSecret.receivedSecret)
		}
	}
}

func TestSecretGet(t *testing.T) {
	gm.RegisterTestingT(t)

	secret := base64.New([]byte("secretsecret"))

	key, _ := crypto.NewBase64KeyURL(nil)
	kms, _ := crypto.NewLocalKMS(key)
	codec := crypto.NewSecretBoxCodec(kms)

	enc, _ := codec.Encrypt(context.TODO(), secret)
	encryptedSecret := models.SecretArg{
		Type:  "secret",
		Name:  "my-key",
		Value: enc,
	}

	secretSpec := models.SecretArg{
		Type:  "secret",
		Name:  "my-key",
		Value: secret,
	}

	tests := []struct {
		secret     models.SecretArg
		wantSecret *models.SecretArg
		wantErr    error
		err        error
	}{
		{
			secret:     encryptedSecret,
			wantSecret: &secretSpec,
			wantErr:    nil,
			err:        nil,
		},
		{
			secret:  encryptedSecret,
			wantErr: ErrGenericDBError,
			err:     ErrGenericDBError,
		},
	}

	pgSecret := &testPgSecret{}
	for i, test := range tests {
		secretDB := secretEncrypt{
			db:    pgSecret,
			codec: codec,
		}

		pgSecret.returnSecret = test.secret
		pgSecret.err = test.err

		gotSpec, gotErr := secretDB.Get(context.TODO(), test.secret.Name)
		if (test.wantErr == nil && gotErr != nil) ||
			(test.wantErr != nil && gotErr != nil && gotErr.Error() != test.wantErr.Error()) {
			t.Errorf("unexpected error on Get() test %d of %d: got %v want %v", i+1, len(tests), gotErr, test.wantErr)
		}

		gm.Expect(gotSpec).To(gm.Equal(test.wantSecret))
	}
}

type testPgSecret struct {
	returnSecret   models.SecretArg
	receivedSecret models.SecretArg
	err            error
}

// only testing the below two for now, rest can remain unimplemented

func (t *testPgSecret) Create(ctx context.Context, secret models.SecretArg) error {
	t.receivedSecret = secret
	return t.err
}

func (t *testPgSecret) Get(ctx context.Context, name string) (*models.SecretArg, error) {
	return &t.returnSecret, t.err
}

func (t *testPgSecret) Delete(ctx context.Context, name string) (db.DeleteStatus, error) {
	panic("not implemented")
}
