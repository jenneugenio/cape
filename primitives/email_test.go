package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestEmail(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("create valid user email", func(t *testing.T) {
		email, err := NewEmail("email@email.com")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(email.String()).To(gm.Equal("email@email.com"))
	})

	t.Run("Can gql unmarshal", func(t *testing.T) {
		var email Email

		err := email.UnmarshalGQL("email@email.com")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(email.String()).To(gm.Equal("email@email.com"))
	})

	t.Run("GQL unmarshal throws validate error", func(t *testing.T) {
		var email Email

		err := email.UnmarshalGQL("notanemail")
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(errors.CausedBy(err, InvalidNameCause))
	})
}
