package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
)

func TestAttachment(t *testing.T) {
	gm.RegisterTestingT(t)

	data, err := loadPolicy("policy.yaml")
	gm.Expect(err).To(gm.BeNil())

	spec, err := ParsePolicySpec(data)
	gm.Expect(err).To(gm.BeNil())

	policy, err := NewPolicy(Label("cool-policy"), spec)
	gm.Expect(err).To(gm.BeNil())

	role, err := NewRole("01EC348BGA506B3FSW6VZHMSX6", Label("role"), false)
	gm.Expect(err).To(gm.BeNil())

	t.Run("valid attachment", func(t *testing.T) {
		attachment, err := NewAttachment(policy.ID, role.ID.String())
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(attachment.PolicyID).To(gm.Equal(policy.ID))
		gm.Expect(attachment.RoleID).To(gm.Equal(role.ID))
	})

	tests := []struct {
		name       string
		identityID database.ID
		roleID     string
	}{
		{
			"invalid policy id",
			database.EmptyID,
			role.ID.String(),
		},
		{
			"invalid role id",
			policy.ID,
			"not valid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewAttachment(tc.identityID, tc.roleID)
			gm.Expect(err).ToNot(gm.BeNil())
		})
	}
}
