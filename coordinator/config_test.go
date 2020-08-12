package coordinator

import (
	"github.com/capeprivacy/cape/models"
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestCoordinatorConfg(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("config validation", func(t *testing.T) {
		validURL, err := models.NewDBURL("postgres://user:pass@host.com/5432")
		gm.Expect(err).To(gm.BeNil())

		validCfg, err := NewConfig(8080, validURL)
		gm.Expect(err).To(gm.BeNil())

		tests := []struct {
			name  string
			fn    func() (*Config, error)
			cause *errors.Cause
		}{
			{
				name: "valid config",
				fn: func() (*Config, error) {
					return NewConfig(8080, validURL)
				},
			},
			{
				name: "invalid dburl",
				fn: func() (*Config, error) {
					return &Config{
						Version: 1,
						Port:    8080,
						DB:      &DBConfig{},
						RootKey: validCfg.RootKey,
					}, nil
				},
				cause: &models.InvalidConfigCause,
			},
			{
				name: "Missing dburl",
				fn: func() (*Config, error) {
					return &Config{
						Version: 1,
						Port:    8080,
						RootKey: validCfg.RootKey,
					}, nil
				},
				cause: &InvalidConfigCause,
			},
			{
				name: "invalid version",
				fn: func() (*Config, error) {
					return &Config{
						Version: 23,
						Port:    8080,
						DB:      validCfg.DB,
						RootKey: validCfg.RootKey,
					}, nil
				},
				cause: &InvalidConfigCause,
			},
			{
				name: "port too low",
				fn: func() (*Config, error) {
					return &Config{
						Version: 1,
						Port:    -23,
						DB:      validCfg.DB,
						RootKey: validCfg.RootKey,
					}, nil
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
				gm.Expect(cfg).ToNot(gm.BeNil())
			})
		}
	})
}
