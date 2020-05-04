package framework

import (
	"context"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

// QueryIdentityPolicies is a helper function to query all roles assigned to an identity and then
// all policies attached to those roles.
func QueryIdentityPolicies(ctx context.Context, db database.Querier, identityID database.ID) ([]*primitives.Policy, error) {
	var assignments []*primitives.Assignment
	assignmentFilter := database.NewFilter(database.Where{"identity_id": identityID.String()}, nil, nil)
	err := db.Query(ctx, &assignments, assignmentFilter)
	if err != nil {
		return nil, err
	}

	roleIDs := database.InFromEntities(assignments, func(e interface{}) interface{} {
		return e.(*primitives.Assignment).RoleID
	})

	if len(roleIDs) == 0 {
		return []*primitives.Policy{}, nil
	}

	var attachments []*primitives.Attachment
	attachmentFilter := database.NewFilter(database.Where{"role_id": roleIDs}, nil, nil)
	err = db.Query(ctx, &attachments, attachmentFilter)
	if err != nil {
		return nil, err
	}

	policyIDs := database.InFromEntities(attachments, func(e interface{}) interface{} {
		return e.(*primitives.Attachment).PolicyID
	})

	if len(policyIDs) == 0 {
		return []*primitives.Policy{}, nil
	}

	var policies []*primitives.Policy
	err = db.Query(ctx, &policies, database.NewFilter(database.Where{"id": policyIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	return policies, nil
}
