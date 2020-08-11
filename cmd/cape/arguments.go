package main

import (
	"fmt"
	"github.com/capeprivacy/cape/models"
)

var (
	ClusterLabelArg = LabelArg("cluster")
	ProjectLabelArg = LabelArg("project-label")

	ClusterURLArg = &Argument{
		Name:        "url",
		Description: "A url for the cluster.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return models.NewURL(in)
		},
	}

	UserEmailArg = &Argument{
		Name:        "email",
		Description: "An email for a user.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			// validates!!
			return models.Email(in), nil
		},
	}

	TokenUserArg = &Argument{
		Name:        "user",
		Description: "The email of the user for this token.",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return models.Email(in), nil
		},
	}

	TokenIDArg = &Argument{
		Name:        "token-id",
		Description: "The ID for the token.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return in, nil
		},
	}

	ProjectNameArg = &Argument{
		Name:        "name",
		Description: "The name of your project.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return models.ProjectDisplayName(in), nil
		},
	}

	ProjectDescriptionArg = &Argument{
		Name:        "description",
		Description: "Describe what your project is for.",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return models.ProjectDescription(in), nil
		},
	}

	SuggestionNameArg = &Argument{
		Name:        "name",
		Description: "The title for your suggestion.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			return models.ProjectDisplayName(in), nil
		},
	}

	SuggestionDescriptionArg = &Argument{
		Name:        "description",
		Description: "Describe your policy suggestion.",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return models.ProjectDescription(in), nil
		},
	}

	SuggestionIDArg = &Argument{
		Name:        "suggestion-id",
		Description: "The ID for your policy suggestion.",
		Required:    false,
		Processor: func(in string) (interface{}, error) {
			return in, nil
		},
	}

	RoleArg = &Argument{
		Name:        "role",
		Description: "The role you wish to assign.",
		Required:    true,
		Processor: func(in string) (interface{}, error) {
			for _, r := range models.SystemRoles {
				if in == r.String() {
					return models.Label(in), nil
				}
			}

			return nil, fmt.Errorf("invalid role specified")
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
			return models.Label(in), nil
		},
	}
}
