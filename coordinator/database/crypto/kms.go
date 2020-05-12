package crypto

import (
	"context"
	"crypto/rand"

	"github.com/manifoldco/go-base64"
	"golang.org/x/crypto/nacl/secretbox"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// KMS intends to be an abstract interface
// over a Key Management System which generally is used
// to wrap DEKs using a key encryption key (KEK).
type KMS interface {
	Open(context.Context) error
	Encrypt(context.Context, []byte) ([]byte, error)
	Decrypt(context.Context, []byte) ([]byte, error)
	Close() error
}

func NewLocalKMS(url *KeyURL) (*LocalKMS, error) {
	k, err := base64.NewFromString(url.Host)
	if err != nil {
		return nil, err
	}

	var key [32]byte

	copy(key[:], *k)

	return &LocalKMS{
		key: key,
		url: url,
	}, nil
}

// LocalKMS is a simple simulated KMS that has a single key
// which is then used to encrypt other keys. This should be
// able to be expanded to something that can handle rotating
// keys.
type LocalKMS struct {
	url *KeyURL
	key [KeyLength]byte
}

// Encrypts the data encryption key (dek) returning the encrypted bytes. The
// result is appended to the nonce.
func (l *LocalKMS) Encrypt(ctx context.Context, dek []byte) ([]byte, error) {
	return Encrypt(l.key, dek)
}

// Decrypt a wrapped dek and return it
func (l *LocalKMS) Decrypt(ctx context.Context, wrappedDEK []byte) ([]byte, error) {
	return Decrypt(l.key, wrappedDEK)
}

func (l *LocalKMS) Open(ctx context.Context) error {
	return nil
}

func (l *LocalKMS) Close() error {
	return nil
}

func LoadKMS(url *KeyURL) (KMS, error) {
	switch url.Type() {
	case Base64Key:
		return NewLocalKMS(url)
	default:
		return nil, errors.New(InvalidKeyURLCause, "Could not find url type %s for loading KMS", url.Type())
	}
}

func GenerateKey() ([KeyLength]byte, error) {
	var key [32]byte

	_, err := rand.Read(key[:])
	if err != nil {
		return [32]byte{}, err
	}
	return key, nil
}

func Encrypt(key [KeyLength]byte, data []byte) ([]byte, error) {
	var nonce [NonceLength]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, err
	}

	return secretbox.Seal(nonce[:], data, &nonce, &key), nil
}

func Decrypt(key [KeyLength]byte, encrypted []byte) ([]byte, error) {
	var decryptNonce [NonceLength]byte
	copy(decryptNonce[:], encrypted[:NonceLength])

	decrypted, ok := secretbox.Open(nil, encrypted[NonceLength:], &decryptNonce, &key)
	if !ok {
		return nil, errors.New(KMSDecryptCause, "Unable to decrypt data")
	}
	return decrypted, nil
}
