package main

import (
	"fmt"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

var (
	ClusterLabelArg     = LabelArg("cluster")
	CoordinatorLabelArg = LabelArg("coordinator")
	RoleLabelArg        = LabelArg("role")
	PolicyLabelArg      = LabelArg("policy")
	ProjectLabelArg     = LabelArg("project-label")

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

	ProjectNameArg = &Argument{
		Name:        "name",
		Description: "The name of your project",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return primitives.NewDisplayName(in)
		},
	}

	ProjectDescriptionArg = &Argument{
		Name:        "description",
		Description: "Describe what your project is for",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return primitives.NewDescription(in)
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
