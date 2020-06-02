package connector

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestConnectorConfig(t *testing.T) {
	gm.RegisterTestingT(t)

	validToken, err := auth.GenerateToken()
	gm.Expect(err).To(gm.BeNil())

	validCoordURL, err := primitives.NewURL("https://localhost:5000")
	gm.Expect(err).To(gm.BeNil())

	tests := []struct {
		name  string
		cfg   *Config
		cause *errors.Cause
	}{
		{
			name: "valid config",
			cfg: &Config{
				InstanceID:     primitives.Label("hisdfsf"),
				Port:           2313,
				Token:          validToken,
				CoordinatorURL: validCoordURL,
			},
		},
		{
			name: "port too high",
			cfg: &Config{
				InstanceID:     primitives.Label("hisf"),
				Port:           2313232,
				Token:          validToken,
				CoordinatorURL: validCoordURL,
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "port too low",
			cfg: &Config{
				InstanceID:     primitives.Label("hisdf"),
				Port:           -2,
				Token:          validToken,
				CoordinatorURL: validCoordURL,
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "token is missing",
			cfg: &Config{
				InstanceID:     primitives.Label("hisdf"),
				Port:           323,
				CoordinatorURL: validCoordURL,
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "coord url is missing",
			cfg: &Config{
				InstanceID: primitives.Label("hisdf"),
				Port:       323,
				Token:      validToken,
			},
			cause: &InvalidConfigCause,
		},
		{
			name: "instance id is invalid",
			cfg: &Config{
				InstanceID:     primitives.Label("1"),
				Port:           323,
				Token:          validToken,
				CoordinatorURL: validCoordURL,
			},
			cause: &primitives.InvalidLabelCause,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.cause != nil {
				gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(err).To(gm.BeNil())
		})
	}
}
