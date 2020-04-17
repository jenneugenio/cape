package primitives

import (
	"crypto/ed25519"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/database"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestAssignment(t *testing.T) {
	gm.RegisterTestingT(t)

	pub, _, _ := ed25519.GenerateKey(nil)
	creds, err := NewCredentials(base64.New(pub), base64.New([]byte("SALTSALTSALTSALT")))
	gm.Expect(err).To(gm.BeNil())

	email, err := NewEmail("email@email.com")
	gm.Expect(err).To(gm.BeNil())

	user, err := NewUser(Name("name"), email, creds)
	gm.Expect(err).To(gm.BeNil())

	role, err := NewRole(Label("role"), false)
	gm.Expect(err).To(gm.BeNil())

	t.Run("valid assignment", func(t *testing.T) {
		assignment, err := NewAssignment(user.ID, role.ID)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(assignment.IdentityID).To(gm.Equal(user.ID))
		gm.Expect(assignment.RoleID).To(gm.Equal(role.ID))
	})

	tests := []struct {
		name       string
		identityID database.ID
		roleID     database.ID
	}{
		{
			"invalid identity id",
			database.EmptyID,
			role.ID,
		},
		{
			"invalid role id",
			user.ID,
			database.EmptyID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewAssignment(tc.identityID, tc.roleID)
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, InvalidAssignmentCause)).To(gm.BeTrue())
		})
	}
}
