package main

import (
	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

var (
	capeTokenVar = &EnvVar{
		Name:        "CAPE_TOKEN",
		Required:    true,
		Description: "A token the data connector can use to authenticate with a coordinator",
		Processor: func(in string) (interface{}, error) {
			if in == "" {
				return nil, errors.New(InvalidAPITokenCause, "A token must be provided.")
			}

			return auth.ParseAPIToken(in)
		},
	}
	capePasswordVar = &EnvVar{
		Name:        "CAPE_PASSWORD",
		Required:    false,
		Description: "The password used by a human to log into their Cape account",
		Processor: func(in string) (interface{}, error) {
			if in == "" {
				return in, nil
			}

			return primitives.NewPassword(in)
		},
	}

	capeDBURL            = DBURLEnvVar(true)
	capeDBURLNotRequired = DBURLEnvVar(false)
)

func DBURLEnvVar(required bool) *EnvVar {
	return &EnvVar{
		Name:        "CAPE_DB_URL",
		Required:    required,
		Description: "The URL for the database.",
		Processor: func(in string) (interface{}, error) {
			return primitives.NewDBURL(in)
		},
	}
}
