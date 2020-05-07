package crypto

import (
	"context"
)

// Encryptable represents a primitive that can be encrypted
// and decrypted.
type Encryptable interface {
	// Encrypt uses codec to encrypt the data and then
	// marshals it into json, returning the encoded bytes
	Encrypt(context.Context, EncryptionCodec) ([]byte, error)

	// Decrypt uses codec to decrypt the data and then
	// unmarshals it into the implementing struct
	Decrypt(context.Context, EncryptionCodec, []byte) error
}
