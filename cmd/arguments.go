package main

import (
	"fmt"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

var (
	ClusterLabelArg     = LabelArg("cluster", true)
	CoordinatorLabelArg = LabelArg("coordinator", true)
	RoleLabelArg        = LabelArg("role", true)
	PolicyLabelArg      = LabelArg("policy", true)
	SourceLabelArg      = LabelArg("source", true)
	CollectionLabelArg  = LabelArg("collection", false)

	ClusterURLArg = &Argument{
		Name:        "url",
		Description: "A url for the cluster",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return primitives.NewURL(in)
		},
	}

	UserEmailArg = &Argument{
		Name:        "email",
		Description: "An email for a user",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			// validates!!
			return primitives.NewEmail(in)
		},
	}

	ServiceIdentifierArg = &Argument{
		Name:        "identifier",
		Description: "An identifier for a service in the form of an email",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			// validates!!
			return primitives.NewEmail(in)
		},
	}

	SourcesCredentialsArg = &Argument{
		Name:        "connection-string",
		Description: "The connection string for the database.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return primitives.NewDBURL(in)
		},
	}

	PullQueryArgument = &Argument{
		Name:        "query",
		Description: "The SQL query to query the data with.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return in, nil
		},
	}

	TokenIdentityArg = &Argument{
		Name:        "identity",
		Description: "The identity for the owner of the token in the form of an email",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return primitives.NewEmail(in)
		},
	}

	TokenIDArg = &Argument{
		Name:        "token-id",
		Description: "The ID for the token",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return database.DecodeFromString(in)
		},
	}
)

func LabelArg(f string, required bool) *Argument {
	return &Argument{
		Name:        f,
		Description: fmt.Sprintf("A label for the %s", f),
		Required:    required,
		Processor: func(in string) (interface{}, error) {
			if in != "" && !required {
				// NewLabel validates that the label meets label criteria
				return primitives.NewLabel(in)
			} else if in == "" {
				return nil, nil
			}

			return primitives.NewLabel(in)
		},
	}
}
