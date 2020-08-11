package models

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestToken(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user := GenerateUser("hello", "bob@hello.com")
	t.Run("create token", func(t *testing.T) {
		tests := []struct {
			name   string
			userID string
			creds  *Credentials
			cause  string
		}{
			{
				name:   "valid parameters",
				userID: user.ID,
				creds:  &user.Credentials,
			},
			{
				name:   "missing credentials",
				userID: user.ID,
				creds:  nil,
				cause:  "credentials must be non-nil",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				token := NewToken(tc.userID, tc.creds)
				err := token.Validate()
				if tc.cause != "" {
					gm.Expect(err).ToNot(gm.BeNil())
					gm.Expect(err.Error()).To(gm.Equal(tc.cause))
					return
				}

				gm.Expect(token.UserID).To(gm.Equal(tc.userID))
				gm.Expect(token.Credentials.Alg).To(gm.Equal(tc.creds.Alg))
				gm.Expect(token.Credentials.Salt).To(gm.Equal(tc.creds.Salt))
				gm.Expect(token.Credentials.Secret).To(gm.Equal(tc.creds.Secret))
			})
		}
	})
}
