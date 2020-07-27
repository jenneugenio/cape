package crypto

import (
	"context"

	"gocloud.dev/secrets"

	// for azure key value
	_ "gocloud.dev/secrets/azurekeyvault"
)

// The "azurekeyvault" URL scheme is replaced with "https" to construct an Azure
// Key Vault keyID, as described in https://docs.microsoft.com/en-us/azure/key-vault/about-keys-secrets-and-certificates.
// You can add an optional "/{key-version}" to the path to use a specific
// version of the key; it defaults to the latest version.

type AzureKMS struct {
	url    *KeyURL
	keeper *secrets.Keeper
}

func NewAzureKMS(url *KeyURL) (*AzureKMS, error) {
	return &AzureKMS{
		url: url,
	}, nil
}

func (a *AzureKMS) Open(ctx context.Context) error {
	keeper, err := secrets.OpenKeeper(ctx, a.url.String())
	if err != nil {
		return err
	}
	a.keeper = keeper

	return nil
}

func (a *AzureKMS) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	return a.keeper.Encrypt(ctx, plaintext)
}

func (a *AzureKMS) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	return a.keeper.Decrypt(ctx, ciphertext)
}

func (a *AzureKMS) Close() error {
	return a.keeper.Close()
}

func (a *AzureKMS) EncryptedKeyLength() int {
	// TODO but why 342
	return 342
}
