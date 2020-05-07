package crypto

import (
	"context"

	"github.com/manifoldco/go-base64"
)

// EncryptionCodec represents a way to encrypt binary data
// with a symmetric key. SecretBoxCodec can be used as an
// example implementation
type EncryptionCodec interface {
	Encrypt(context.Context, *base64.Value) (*base64.Value, error)
	Decrypt(context.Context, *base64.Value) (*base64.Value, error)
}
