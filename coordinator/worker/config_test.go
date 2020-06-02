package worker

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestWorkerConfig(t *testing.T) {
	gm.RegisterTestingT(t)

	logger := framework.TestLogger()

	validToken, err := auth.GenerateToken()
	gm.Expect(err).To(gm.BeNil())

	validDBURL, err := primitives.NewDBURL("postgres://user:pass@host:5432/db")
	gm.Expect(err).To(gm.BeNil())

	validCoordURL, err := primitives.NewURL("https://localhost:5000")
	gm.Expect(err).To(gm.BeNil())

	tests := []struct {
		name  string
		fn    func() (*Config, error)
		cause *errors.Cause
	}{
		{
			name: "valid config",
			fn: func() (*Config, error) {
				return NewConfig(validToken, validDBURL, validCoordURL, logger)
			},
		},
		{
			name: "token missing",
			fn: func() (*Config, error) {
				cfg, err := NewConfig(validToken, validDBURL, validCoordURL, logger)
				if err != nil {
					return nil, err
				}

				cfg.Token = nil
				return cfg, nil
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "database url missing",
			fn: func() (*Config, error) {
				cfg, err := NewConfig(validToken, validDBURL, validCoordURL, logger)
				if err != nil {
					return nil, err
				}

				cfg.DatabaseURL = nil
				return cfg, nil
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "coordinator url missing",
			fn: func() (*Config, error) {
				cfg, err := NewConfig(validToken, validDBURL, validCoordURL, logger)
				if err != nil {
					return nil, err
				}

				cfg.CoordinatorURL = nil
				return cfg, nil
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "logger missing",
			fn: func() (*Config, error) {
				cfg, err := NewConfig(validToken, validDBURL, validCoordURL, logger)
				if err != nil {
					return nil, err
				}

				cfg.Logger = nil
				return cfg, nil
			},
			cause: &InvalidConfigCause,
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
