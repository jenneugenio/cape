package main

import (
	"fmt"
	"net/url"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

var (
	LabelArg = func(f string) *Argument {
		return &Argument{
			Name:        "label",
			Description: fmt.Sprintf("A label for the %s", f),
			Required:    true,
			Processor: func(in string) (interface{}, error) {
				// NewLabel validates that the label meets label criteria
				return primitives.NewLabel(in)
			},
		}
	}

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

	RoleLabelArg = &Argument{
		Name:        "label",
		Description: "An label for a role",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			// validates!!
			return primitives.NewLabel(in)
		},
	}

	SourcesCredentialsArg = &Argument{
		Name:        "connection-string",
		Description: "The connection string for the database.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			if in == "" {
				return nil, errors.New(InvalidURLCause, "A valid url must be provided")
			}

			u, err := url.Parse(in)
			if err != nil {
				return nil, errors.New(InvalidURLCause, "could not parse: %s", err)
			}

			if u.Scheme != "postgres" {
				return nil, errors.New(InvalidURLCause, "Invalid database type. Currently only postgres is supported.")
			}

			if u.Host == "" {
				return nil, errors.New(InvalidURLCause, "A host must be provided")
			}

			return u, nil
		},
	}
)
