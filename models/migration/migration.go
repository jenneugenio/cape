package modelmigration

import (
	"github.com/capeprivacy/cape/coordinator/database"
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

func ProjectFromPrimitive(prim *primitives.Project) models.Project {
	return models.Project{
		ID:    prim.ID.String(),
		Label: models.Label(prim.Label.String()),
	}
}

func PrimitiveFromProject(p models.Project) (*primitives.Project, error) {
	id, err := database.DecodeFromString(p.ID)
	if err != nil {
		return nil, err
	}

	return &primitives.Project{
		Primitive: &database.Primitive{ID: id},
		Name:      "TODO",
		Label:     primitives.Label(p.Label),
	}, nil
}

func PrimitiveFromRole(r models.Role) (*primitives.Role, error) {
	id, err := database.DecodeFromString(r.ID)
	if err != nil {
		return nil, err
	}

	return &primitives.Role{
		Primitive: &database.Primitive{ID: id},
		Label:     primitives.Label(r.Label),
		System:    true,
	}, nil
}
