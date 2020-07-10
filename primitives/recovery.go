package primitives

import (
	"context"
	"encoding/json"
	"time"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// RecoveryExpiration is the amount of time that has passed since a recovery
// was created before it's no longer valid.
var RecoveryExpiration time.Duration = 30 * time.Minute

type Recovery struct {
	*database.Primitive
	UserID      string       `json:"user_id"`
	Credentials *Credentials `json:"-" gqlgen:"-"`
	ExpiresAt   time.Time    `json:"expires_at"`
}

type encryptedRecovery struct {
	*Recovery
	Credentials *base64.Value `json:"credentials"`
}

func (r *Recovery) Validate() error {
	if err := r.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidRecoveryCause, err)
	}

	if r.UserID == "" {
		return errors.New(InvalidRecoveryCause, "UserID must not be empty")
	}

	if r.Credentials == nil {
		return errors.New(InvalidRecoveryCause, "Missing credentials")
	}

	if err := r.Credentials.Validate(); err != nil {
		return errors.Wrap(InvalidRecoveryCause, err)
	}

	if r.ExpiresAt.IsZero() {
		return errors.New(InvalidRecoveryCause, "Missing expires at")
	}

	return nil
}

func (r *Recovery) Expired() bool {
	return time.Now().UTC().After(r.ExpiresAt)
}

func (r *Recovery) GetType() types.Type {
	return RecoveryType
}

func (r *Recovery) GetEncryptable() bool {
	return true
}

func (r *Recovery) Encrypt(ctx context.Context, codec crypto.EncryptionCodec) ([]byte, error) {
	creds, err := json.Marshal(r.Credentials)
	if err != nil {
		return nil, err
	}

	data, err := codec.Encrypt(ctx, base64.New(creds))
	if err != nil {
		return nil, err
	}

	return json.Marshal(&encryptedRecovery{
		Recovery:    r,
		Credentials: data,
	})
}

func (r *Recovery) Decrypt(ctx context.Context, codec crypto.EncryptionCodec, data []byte) error {
	in := &encryptedRecovery{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Credentials)
	if err != nil {
		return err
	}

	creds := &Credentials{}
	err = json.Unmarshal([]byte(*unencrypted), creds)
	if err != nil {
		return err
	}

	r.Primitive = in.Primitive
	r.UserID = in.UserID
	r.Credentials = creds
	r.ExpiresAt = in.ExpiresAt

	return nil
}

func NewRecovery(userID string, creds *Credentials) (*Recovery, error) {
	p, err := database.NewPrimitive(RecoveryType)
	if err != nil {
		return nil, err
	}

	r := &Recovery{
		Primitive:   p,
		UserID:      userID,
		Credentials: creds,
		ExpiresAt:   time.Now().UTC().Add(RecoveryExpiration),
	}

	id, err := database.DeriveID(r)
	if err != nil {
		return nil, err
	}

	r.ID = id
	return r, r.Validate()
}

func GenerateRecovery() (*Recovery, error) {
	userID := "thisisanid"

	return NewRecovery(userID, GenerateCredentials())
}
