package primitives

import "github.com/manifoldco/go-base64"

type Credentials struct {
	PublicKey *base64.Value      `json:"public_key"`
	Salt      *base64.Value      `json:"salt"`
	Alg       CredentialsAlgType `json:"alg"`
}
