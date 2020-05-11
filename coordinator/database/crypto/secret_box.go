package crypto

import (
	"context"
	"crypto/rand"

	"github.com/manifoldco/go-base64"
	"golang.org/x/crypto/nacl/secretbox"

	errors "github.com/capeprivacy/cape/partyerrors"
)

const (
	NonceLength        = 24
	KeyLength          = 32
	EncryptedKeyLength = NonceLength + KeyLength + secretbox.Overhead // 72
)

func NewSecretBoxCodec(kms KMS) *SecretBoxCodec {
	return &SecretBoxCodec{kms: kms}
}

// SecretBoxCodec implements a envelope encryption scheme
// where it leverages data encryption keys (DEKs) and
// key encryption keys (KEKs) to safely encrypt the data and
// prevent leaking the keys.
// Here's a pretty good overview of envelope encryption:
// https://cloud.google.com/kms/docs/envelope-encryption
// See individual function comments for more information
type SecretBoxCodec struct {
	kms KMS
}

// Encrypt generates a random nonce and DEK which is then used to call
// secretbox.Seal. The result is appended to the nonce so the nonce can
// be used later to decrypt the data. The DEK is then encrypted and the
// result plus the nonce is appended to the wrapped DEK.
func (s *SecretBoxCodec) Encrypt(ctx context.Context, data *base64.Value) (*base64.Value, error) {
	var nonce [NonceLength]byte
	var dek [KeyLength]byte

	// generate some random data for nonce and the dek
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, err
	}

	_, err = rand.Read(dek[:])
	if err != nil {
		return nil, err
	}

	// this appends the encrypted data to the nonce
	encrypted := secretbox.Seal(nonce[:], *data, &nonce, &dek)

	// encrypt the dek
	wrappedDEK, err := s.kms.Encrypt(ctx, dek[:])
	if err != nil {
		return nil, err
	}

	// append the encrypted data to the wrapped dek
	encrypted = append(wrappedDEK, encrypted...)

	return base64.New(encrypted), nil
}

func (s *SecretBoxCodec) Decrypt(ctx context.Context, data *base64.Value) (*base64.Value, error) {
	encrypted := []byte(*data)

	// separate the wrapped dek from the encrypted payload
	var wrappedDEK [EncryptedKeyLength]byte
	copy(wrappedDEK[:], encrypted[:EncryptedKeyLength])

	encrypted = encrypted[EncryptedKeyLength:]

	key, err := s.kms.Decrypt(ctx, wrappedDEK[:])
	if err != nil {
		return nil, err
	}

	var nkey [32]byte
	copy(nkey[:], key)

	var decryptNonce [NonceLength]byte
	copy(decryptNonce[:], encrypted[:NonceLength])
	decrypted, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, &nkey)
	if !ok {
		return nil, errors.New(SecretBoxDecryptCause, "Unable to decrypt data")
	}

	return base64.New(decrypted), nil
}
