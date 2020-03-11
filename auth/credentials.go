package auth

import (
	"crypto/ed25519"

	"github.com/manifoldco/go-base64"
)

// CredentialsAlgType enum holding the supported crypto algorithms
type CredentialsAlgType string

var (
	// EDDSA algorithm type
	EDDSA CredentialsAlgType = "eddsa"
)

// Credentials holds the public key and nonce for the
type Credentials struct {
	// privateKey not saved in the database or sent anywhere!
	privateKey *ed25519.PrivateKey // nolint: structcheck, unused
	PublicKey  *base64.Value       `json:"public_key"`
	Salt       *base64.Value       `json:"nonce"`
	Alg        CredentialsAlgType  `json:"alg"`
}
