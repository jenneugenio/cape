package graph

import (
	"context"
	modelmigration "github.com/capeprivacy/cape/models/migration"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

// buildAttachment takes a primitives attachment and builds at graphql
// model representation of it
func buildAttachment(ctx context.Context, enforcer *auth.Enforcer, db db.Interface,
	attachment *primitives.Attachment) (*model.Attachment, error) {
	deprecatedRole := &primitives.Role{}
	err := enforcer.Get(ctx, attachment.RoleID, deprecatedRole)
	if err != nil {
		return nil, err
	}

	policy, err := db.Policies().GetByID(ctx, attachment.PolicyID)
	if err != nil {
		return nil, err
	}

	role := modelmigration.RoleFromPrimitive(*deprecatedRole)

	return &model.Attachment{
		ID:        attachment.ID,
		CreatedAt: attachment.CreatedAt,
		UpdatedAt: attachment.UpdatedAt,
		Role:      &role,
		Policy:    policy,
	}, nil
}

func hasRole(roles []*models.Role, label models.Label) bool {
	found := false
	for _, role := range roles {
		if role.Label == label {
			found = true
			break
		}
	}

	return found
}

func ruleInputsToModelRules(inputs []*model.RuleInput) []models.Rule {
	rules := make([]models.Rule, len(inputs))
	for i, input := range inputs {
		rules[i].Match.Name = input.Match.Name

		rules[i].Actions = make([]models.Action, len(input.Actions))
		for j, action := range input.Actions {
			rules[i].Actions[j].Transform = action.Transform
		}
	}

	return rules
}
