package modelmigration

import (
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func PoliciesFromPrimitive(prims []*primitives.Policy) []*models.Policy {
	policies := make([]*models.Policy, 0, len(prims))
	for _, prim := range prims {
		p := PolicyFromPrimitive(prim)
		policies = append(policies, &p)
	}
	return policies
}

func PolicyFromPrimitive(prim *primitives.Policy) models.Policy {
	spec := PolicySpecFromPrimitive(prim.Spec)
	p := models.Policy{
		ID:        prim.Primitive.ID.String(),
		Version:   prim.Primitive.Version,
		Label:     LabelFromPrimitive(prim.Label),
		Spec:      &spec,
		CreatedAt: prim.Primitive.CreatedAt,
		UpdatedAt: prim.Primitive.UpdatedAt,
	}
	return p
}

func PolicySpecFromPrimitive(prim *primitives.PolicySpec) models.PolicySpec {
	s := models.PolicySpec{
		Version: PolicyVersionFromPrimitive(prim.Version),
		Label:   LabelFromPrimitive(prim.Label),
		Rules:   RulesFromPrimitive(prim.Rules),
	}
	return s
}

func LabelFromPrimitive(prim primitives.Label) models.Label { return models.Label(prim.String()) }

func PolicyVersionFromPrimitive(prim primitives.PolicyVersion) models.PolicyVersion {
	return models.PolicyVersion(uint8(prim))
}

func RulesFromPrimitive(prims []*primitives.Rule) []*models.Rule {
	rules := make([]*models.Rule, 0, len(prims))
	for _, prim := range prims {
		rule := RuleFromPrimitive(prim)
		rules = append(rules, &rule)
	}
	return rules
}

func RuleFromPrimitive(prim *primitives.Rule) models.Rule {
	r := models.Rule{
		Target:          TargetFromPrimitive(prim.Target),
		Action:          ActionFromPrimitive(prim.Action),
		Effect:          EffectFromPrimitive(prim.Effect),
		Fields:          FieldsFromPrimitive(prim.Fields),
		Where:           WheresFromPrimitive(prim.Where),
		Transformations: TransformationsFromPrimitive(prim.Transformations),
		Sudo:            prim.Sudo,
	}
	return r
}

func TargetFromPrimitive(prim primitives.Target) models.Target { return models.Target(string(prim)) }

func ActionFromPrimitive(prim primitives.Action) models.Action { return models.Action(string(prim)) }

func EffectFromPrimitive(prim primitives.Effect) models.Effect { return models.Effect(string(prim)) }

func FieldsFromPrimitive(prims []primitives.Field) []models.Field {
	fields := make([]models.Field, 0, len(prims))
	for _, prim := range prims {
		fields = append(fields, FieldFromPrimitive(prim))
	}
	return fields
}

func FieldFromPrimitive(prim primitives.Field) models.Field { return models.Field(string(prim)) }

func WheresFromPrimitive(prims []primitives.Where) []models.Where {
	wheres := make([]models.Where, 0, len(prims))
	for _, prim := range prims {
		wheres = append(wheres, WhereFromPrimitive(prim))
	}
	return wheres
}

func WhereFromPrimitive(prim primitives.Where) models.Where { return models.Where(map[string]string(prim)) }

func TransformationsFromPrimitive(prims []*primitives.Transformation) []*models.Transformation {
	xforms := make([]*models.Transformation, 0, len(prims))
	for _, prim := range prims {
		xform := TransformationFromPrimitive(prim)
		xforms = append(xforms, &xform)
	}
	return xforms
}

func TransformationFromPrimitive(prim *primitives.Transformation) models.Transformation {
	t := models.Transformation{
		Field: FieldFromPrimitive(prim.Field),
		Function: prim.Function,
		Args: ArgsFromPrimitive(prim.Args),
		Where: ConditionFromPrimitive(prim.Where),
	}
	return t
}

func ArgsFromPrimitive(prim primitives.Args) models.Args {
	return models.Args(prim)
}

func ConditionFromPrimitive(prim primitives.Condition) models.Condition { return models.Condition(prim) }
