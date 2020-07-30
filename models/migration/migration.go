package modelmigration

import (
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func LabelFromPrimitive(label primitives.Label) models.Label {
	return models.Label(label.String())
}

func EmailFromPrimitive(email primitives.Email) models.Email {
	return models.Email(email.String())
}

func CredentialsFromPrimitives(prim *primitives.Credentials) *models.Credentials {
	return &models.Credentials{
		Secret: prim.Secret,
		Salt:   prim.Salt,
		Alg:    models.CredentialsAlgType(prim.Alg),
	}
}
