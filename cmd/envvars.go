package main

import (
	"github.com/dropoutlabs/cape/auth"
	errors "github.com/dropoutlabs/cape/partyerrors"
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
)
