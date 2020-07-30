package main

import (
	"github.com/capeprivacy/cape/primitives"
)

var (
	capePasswordVar = &EnvVar{
		Name:        "CAPE_PASSWORD",
		Required:    false,
		Description: "The password used by a human to log into their Cape account.",
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
