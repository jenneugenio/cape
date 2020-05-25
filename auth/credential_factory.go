package auth

import (
	"github.com/capeprivacy/cape/primitives"
)

// CredentialFactory manages the different types of CredentialProducers which
// implement different algorithms for generating and comparing credentials.
type CredentialFactory struct {
	registry ProducerRegistry

	Alg primitives.CredentialsAlgType
}

func NewCredentialFactory(alg primitives.CredentialsAlgType) (*CredentialFactory, error) {
	err := alg.Validate()
	if err != nil {
		return nil, err
	}

	registry := ProducerRegistry{
		primitives.SHA256: &SHA256Producer{},
		primitives.Argon2ID: &Argon2IDProducer{
			Time:      1,
			Memory:    64 * 1024,
			Threads:   4,
			KeyLength: primitives.SecretLength,
		},
	}

	_, err = registry.Get(alg)
	if err != nil {
		return nil, err
	}

	return &CredentialFactory{
		registry: registry,
		Alg:      alg,
	}, nil
}

func (cf *CredentialFactory) Generate(secret primitives.Password) (*primitives.Credentials, error) {
	if err := secret.Validate(); err != nil {
		return nil, err
	}

	producer, err := cf.registry.Get(cf.Alg)
	if err != nil {
		return nil, err
	}

	creds, err := producer.Generate(secret)
	if err != nil {
		return nil, err
	}

	if err := creds.Validate(); err != nil {
		return nil, err
	}

	return creds, nil
}

func (cf *CredentialFactory) Compare(secret primitives.Password, creds *primitives.Credentials) error {
	if err := secret.Validate(); err != nil {
		return err
	}

	producer, err := cf.registry.Get(creds.Alg)
	if err != nil {
		return err
	}

	return producer.Compare(secret, creds)
}
