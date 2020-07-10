package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestAssignment(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user, err := GenerateUser("name", "email@email.com")
	gm.Expect(err).To(gm.BeNil())

	role, err := NewRole(Label("role"), false)
	gm.Expect(err).To(gm.BeNil())

	t.Run("valid assignment", func(t *testing.T) {
		assignment, err := NewAssignment(user.ID.String(), role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment.UserID).To(gm.Equal(user.ID.String()))
		gm.Expect(assignment.RoleID).To(gm.Equal(role.ID))
	})

	tests := []struct {
		name   string
		userID string
		roleID database.ID
	}{
		{
			"invalid user id",
			"",
			role.ID,
		},
		{
			"invalid role id",
			user.ID.String(),
			database.EmptyID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewAssignment(tc.userID, tc.roleID)
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, InvalidAssignmentCause)).To(gm.BeTrue())
		})
	}
}
