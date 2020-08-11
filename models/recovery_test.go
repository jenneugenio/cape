package models

import (
	"testing"
	"time"

	gm "github.com/onsi/gomega"
)

func TestRecovery(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user := GenerateUser("hi", "hi@hi.hi")

	creds := GenerateCredentials()

	t.Run("Validate", func(t *testing.T) {
		tests := []struct {
			name  string
			fn    func() Recovery
			cause string
		}{
			{
				name: "valid recovery",
				fn: func() Recovery {
					return NewRecovery(user.ID, creds)
				},
			},
			{
				name: "invalid id",
				fn: func() Recovery {
					r := NewRecovery(user.ID, creds)
					r.ID = ""
					return r
				},
				cause: "id must not be empty",
			},
			{
				name: "invalid user id",
				fn: func() Recovery {
					r := NewRecovery(user.ID, creds)

					r.UserID = ""
					return r
				},
				cause: "user id must not be empty",
			},
			{
				name: "missing credentials",
				fn: func() Recovery {
					r := NewRecovery(user.ID, creds)
					r.Credentials = nil
					return r
				},
				cause: "missing credentials",
			},
			{
				name: "expires at is zero value",
				fn: func() Recovery {
					r := NewRecovery(user.ID, creds)
					r.ExpiresAt = time.Time{}
					return r
				},
				cause: "missing expires at",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				r := tc.fn()
				err := r.Validate()
				if tc.cause != "" {
					gm.Expect(err).ToNot(gm.BeNil())
					gm.Expect(err.Error()).To(gm.Equal(tc.cause))
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})

	t.Run("expired returns true if expire at exceeded", func(t *testing.T) {
		r := GenerateRecovery()

		r.ExpiresAt = time.Now().UTC().Add(-1 * time.Minute)
		gm.Expect(r.Expired()).To(gm.BeTrue())
	})

	t.Run("expired returns true if expire at exceeded", func(t *testing.T) {
		r := GenerateRecovery()
		gm.Expect(r.Expired()).To(gm.BeFalse())
	})
}
