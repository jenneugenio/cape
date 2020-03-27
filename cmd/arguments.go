package main

import (
	"net/url"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

var (
	ClusterLabelArg = &Argument{
		Name:        "label",
		Description: "A label for the cluster",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			// NewLabel validates that the label meets label criteria
			return primitives.NewLabel(in)
		},
	}
	ClusterURLArg = &Argument{
		Name:        "url",
		Description: "A url for the cluster",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			if in == "" {
				return nil, errors.New(InvalidURLCause, "A valid url must be provided")
			}

			u, err := url.Parse(in)
			if err != nil {
				return nil, errors.New(InvalidURLCause, "could not parse: %s", err)
			}

			if u.Scheme != "https" && u.Scheme != "http" {
				return nil, errors.New(InvalidURLCause, "invalid scheme, must be http or https")
			}

			if u.Host == "" {
				return nil, errors.New(InvalidURLCause, "A host must be provided")
			}

			return u, nil
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
)
