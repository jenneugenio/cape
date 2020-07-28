package modelmigration

import (
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func TargetFromPrimitive(prim primitives.Target) models.Target { return models.Target(string(prim)) }

func EffectFromPrimitive(prim primitives.Effect) models.Effect { return models.Effect(string(prim)) }

func FieldsFromPrimitive(prims []primitives.Field) []models.Field {
	fields := make([]models.Field, 0, len(prims))
	for _, prim := range prims {
		fields = append(fields, FieldFromPrimitive(prim))
	}
	return fields
}

func FieldFromPrimitive(prim primitives.Field) models.Field { return models.Field(string(prim)) }

func CredentialsFromModels(model *models.Credentials) *primitives.Credentials {
	return &primitives.Credentials{
		Secret: model.Secret,
		Salt:   model.Salt,
		Alg:    primitives.CredentialsAlgType(model.Alg),
	}
}

func LabelFromPrimitive(label primitives.Label) models.Label {
	return models.Label(label.String())
}

func CredentialsFromPrimitives(prim *primitives.Credentials) *models.Credentials {
	return &models.Credentials{
		Secret: prim.Secret,
		Salt:   prim.Salt,
		Alg:    models.CredentialsAlgType(prim.Alg),
	}
}
