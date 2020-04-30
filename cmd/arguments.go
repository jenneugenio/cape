package main

import (
	"fmt"

	"github.com/capeprivacy/cape/primitives"
)

var (
	ClusterLabelArg     = LabelArg("cluster")
	CoordinatorLabelArg = LabelArg("coordinator")
	RoleLabelArg        = LabelArg("role")
	PolicyLabelArg      = LabelArg("policy")
	SourceLabelArg      = LabelArg("source")

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
		Description: "The identity for the token",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return in, nil
		},
	}
)

func LabelArg(f string) *Argument {
	return &Argument{
		Name:        f,
		Description: fmt.Sprintf("A label for the %s", f),
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			// NewLabel validates that the label meets label criteria
			return primitives.NewLabel(in)
		},
	}
}
