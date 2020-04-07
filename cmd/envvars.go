package main

import (
	"github.com/dropoutlabs/cape/auth"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

var (
	capeTokenVar = &EnvVar{
		Name:        "CAPE_TOKEN",
		Required:    true,
		Description: "A token the data connector can use to authenticate with a controller",
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
	capeDBPassword = &EnvVar{
		Name:     "CAPE_DB_PASSWORD",
		Required: false,
		Description: "The password for the database. This variable exists so the database password " +
			" can be passed securely without being exposed outside the current userspace.",
		Processor: func(in string) (interface{}, error) {
			return in, nil
		},
	}
)
