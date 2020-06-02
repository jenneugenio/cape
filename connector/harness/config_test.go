package harness

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestHarnessConfig(t *testing.T) {
	gm.RegisterTestingT(t)

	validToken, err := auth.GenerateToken()
	gm.Expect(err).To(gm.BeNil())

	validCoordURL, err := primitives.NewURL("https://localhost:5000")
	gm.Expect(err).To(gm.BeNil())

	validDBURL, err := primitives.NewDBURL("postgres://user:pass@host:5432/db")
	gm.Expect(err).To(gm.BeNil())

	tests := []struct {
		name  string
		fn    func() (*Config, error)
		cause *errors.Cause
	}{
		{
			name: "valid config",
			fn: func() (*Config, error) {
				return &Config{
					coordinatorURL:      validCoordURL,
					token:               validToken,
					sourceMigrationsDir: "tools/seed",
					dbURL:               validDBURL,
				}, nil
			},
		},
		{
			name: "missing coord url",
			fn: func() (*Config, error) {
				return &Config{
					token:               validToken,
					sourceMigrationsDir: "tools/seed",
					dbURL:               validDBURL,
				}, nil
			},
			cause: &MissingConfig,
		},
		{
			name: "missing token",
			fn: func() (*Config, error) {
				return &Config{
					coordinatorURL:      validCoordURL,
					sourceMigrationsDir: "tools/seed",
					dbURL:               validDBURL,
				}, nil
			},
			cause: &MissingConfig,
		},
		{
			name: "missing db url",
			fn: func() (*Config, error) {
				return &Config{
					coordinatorURL:      validCoordURL,
					token:               validToken,
					sourceMigrationsDir: "tools/seed",
				}, nil
			},
			cause: &MissingConfig,
		},
		{
			name: "missing source migrations dir",
			fn: func() (*Config, error) {
				return &Config{
					coordinatorURL: validCoordURL,
					token:          validToken,
					dbURL:          validDBURL,
				}, nil
			},
			cause: &MissingConfig,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := tc.fn()
			gm.Expect(err).To(gm.BeNil())

			err = cfg.Validate()
			if tc.cause != nil {
				gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(err).To(gm.BeNil())
		})
	}
}
