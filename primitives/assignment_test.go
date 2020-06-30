package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
)

func TestAssignment(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user, err := GenerateUser("name", "email@email.com")
	gm.Expect(err).To(gm.BeNil())

	role, err := NewRole("01EC348BGA506B3FSW6VZHMSX6", Label("role"), false)
	gm.Expect(err).To(gm.BeNil())

	t.Run("valid assignment", func(t *testing.T) {
		assignment, err := NewAssignment(user.ID, role.ID.String())
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment.IdentityID).To(gm.Equal(user.ID))
		gm.Expect(assignment.RoleID).To(gm.Equal(role.ID))
	})

	tests := []struct {
		name       string
		identityID database.ID
		roleID     string
	}{
		{
			"invalid identity id",
			database.EmptyID,
			role.ID.String(),
		},
		{
			"invalid role id",
			user.ID,
			"not valid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewAssignment(tc.identityID, tc.roleID)
			gm.Expect(err).ToNot(gm.BeNil())
		})
	}
}
