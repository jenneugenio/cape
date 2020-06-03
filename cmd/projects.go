package main

import (
	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/primitives"
	"github.com/urfave/cli/v2"
)

func init() {
	projectsCreateCmd := &Command{
		Usage:     "Creates a project in Cape",
		Arguments: []*Argument{ProjectNameArg, ProjectDescriptionArg},
		Examples: []*Example{
			{
				Example:     `cape projects create "My Project" "I will Cape the world!"`,
				Description: `Creates a project named "My Project" with the description "I will Cape the world!"`,
			},
		},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(projectsCreate),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	projectsCmd := &Command{
		Usage: "Commands for interacting with Cape projects",
		Command: &cli.Command{
			Name:        "projects",
			Subcommands: []*cli.Command{projectsCreateCmd.Package()},
		},
	}

	commands = append(commands, projectsCmd.Package())
}

func projectsCreate(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	name := Arguments(c.Context, ProjectNameArg).(primitives.DisplayName)
	desc, _ := Arguments(c.Context, ProjectDescriptionArg).(primitives.Description)

	project, err := client.CreateProject(c.Context, name, nil, desc)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	err = u.Template("Created project {{ . | bold }}!\n\n", project.Name.String())
	if err != nil {
		return err
	}

	details := ui.Details{
		"Name":        project.Name.String(),
		"Description": project.Description.String(),
		"Label":       project.Label.String(),
		"Status":      project.Status.String(),
	}

	return u.Details(details)
}
