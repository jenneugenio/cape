package models

import (
	"testing"
	"time"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

func TestSession(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user := GenerateUser("bob", "test@email.com")
	_, token := GenerateToken(user)

	ti := time.Now().UTC().Add(time.Minute * 5)
	sessionToken := base64.New([]byte("random-string"))

	t.Run("validate", func(t *testing.T) {
		tests := []struct {
			name  string
			fn    func() Session
			cause string
		}{
			{
				name: "valid session owned by user",
				fn: func() Session {
					return NewSession(&user)
				},
			},
			{
				name: "valid session owned by token",
				fn: func() Session {
					return NewSession(&token)
				},
			},
			{
				name: "invalid id",
				fn: func() Session {
					session := NewSession(&user)
					session.ID = ""
					return session
				},
				cause: "id is empty",
			},
			{
				name: "invalid user id",
				fn: func() Session {
					session := NewSession(&user)

					session.UserID = ""
					return session
				},
				cause: "user id is empty",
			},
			{
				name: "invalid owner id",
				fn: func() Session {
					session := NewSession(&user)

					session.OwnerID = ""
					return session
				},
				cause: "owner id is empty",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				s := tc.fn()
				err := s.Validate()
				if tc.cause != "" {
					gm.Expect(err).ToNot(gm.BeNil())
					gm.Expect(err.Error()).To(gm.Equal(tc.cause))
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})

	t.Run("new session", func(t *testing.T) {
		session := NewSession(&user)
		session.SetToken(sessionToken, ti)

		gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
		gm.Expect(session.Token).To(gm.Equal(sessionToken))
		gm.Expect(session.UserID).To(gm.Equal(user.ID))
		gm.Expect(session.OwnerID).To(gm.Equal(user.ID))
	})
}
