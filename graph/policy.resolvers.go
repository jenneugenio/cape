package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/model"
	"github.com/dropoutlabs/cape/primitives"
)

func (r *mutationResolver) CreatePolicy(ctx context.Context, input model.CreatePolicyRequest) (*primitives.Policy, error) {
	policy, err := primitives.NewPolicy(input.Label)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *mutationResolver) DeletePolicy(ctx context.Context, input model.DeletePolicyRequest) (*string, error) {
	err := r.Backend.Delete(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mutationResolver) AttachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*model.Attachment, error) {
	attachment, err := primitives.NewAttachment(input.PolicyID, input.RoleID)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Create(ctx, attachment)
	if err != nil {
		return nil, err
	}

	return buildAttachment(ctx, r.Backend, attachment)
}

func (r *mutationResolver) DetachPolicy(ctx context.Context, input model.AttachPolicyRequest) (*string, error) {
	attachment := &primitives.Attachment{}

	filter := database.NewFilter(database.Where{
		"role_id":   input.RoleID.String(),
		"policy_id": input.PolicyID.String(),
	}, nil, nil)

	err := r.Backend.QueryOne(ctx, attachment, filter)
	if err != nil {
		return nil, err
	}

	err = r.Backend.Delete(ctx, attachment.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *queryResolver) Policy(ctx context.Context, id database.ID) (*primitives.Policy, error) {
	policy := &primitives.Policy{}
	err := r.Backend.Get(ctx, id, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *queryResolver) Policies(ctx context.Context) ([]*primitives.Policy, error) {
	var s []primitives.Policy
	err := r.Backend.Query(ctx, &s, database.NewEmptyFilter())
	if err != nil {
		return nil, err
	}

	policies := make([]*primitives.Policy, len(s))
	for i := 0; i < len(policies); i++ {
		policies[i] = &(s[i])
	}

	return policies, nil
}

func (r *queryResolver) RolePolicies(ctx context.Context, roleID database.ID) ([]*primitives.Policy, error) {
	var a []primitives.Attachment
	err := r.Backend.Query(ctx, &a, database.NewFilter(database.Where{"role_id": roleID.String()}, nil, nil))
	if err != nil {
		return nil, err
	}

	policyIDs := make(database.In, len(a))
	for i, attachment := range a {
		policyIDs[i] = attachment.PolicyID
	}

	var tmpP []primitives.Policy
	if len(policyIDs) == 0 {
		return []*primitives.Policy{}, nil
	}

	err = r.Backend.Query(ctx, &tmpP, database.NewFilter(database.Where{"id": policyIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	policies := make([]*primitives.Policy, len(a))
	for i := 0; i < len(policies); i++ {
		policies[i] = &(tmpP[i])
	}

	return policies, nil
}

func (r *queryResolver) IdentityPolicies(ctx context.Context, identityID database.ID) ([]*primitives.Policy, error) {
	assignmentFilter := database.NewFilter(database.Where{"identity_id": identityID.String()}, nil, nil)
	var assignments []primitives.Assignment
	err := r.Backend.Query(ctx, &assignments, assignmentFilter)
	if err != nil {
		return nil, err
	}

	roleIDs := make(database.In, len(assignments))
	for i, assignment := range assignments {
		roleIDs[i] = assignment.RoleID
	}

	attachmentFilter := database.NewFilter(database.Where{"role_id": roleIDs}, nil, nil)
	var attachments []primitives.Attachment
	err = r.Backend.Query(ctx, &attachments, attachmentFilter)
	if err != nil {
		return nil, err
	}

	policyIDs := make(database.In, len(attachments))
	for i, attachment := range attachments {
		policyIDs[i] = attachment.PolicyID
	}

	var tmpP []primitives.Policy
	err = r.Backend.Query(ctx, &tmpP, database.NewFilter(database.Where{"id": policyIDs}, nil, nil))
	if err != nil {
		return nil, err
	}

	policies := make([]*primitives.Policy, len(tmpP))
	for i := 0; i < len(policies); i++ {
		policies[i] = &(tmpP[i])
	}

	return policies, nil
}

func (r *queryResolver) Attachment(ctx context.Context, roleID database.ID, policyID database.ID) (*model.Attachment, error) {
	attachment := &primitives.Attachment{}

	filter := database.NewFilter(database.Where{
		"role_id":   roleID,
		"policy_id": policyID,
	}, nil, nil)

	err := r.Backend.QueryOne(ctx, attachment, filter)
	if err != nil {
		return nil, err
	}

	return buildAttachment(ctx, r.Backend, attachment)
}
