package main

import (
	"fmt"
	"github.com/capeprivacy/cape/prims"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/primitives"
	"github.com/urfave/cli/v2"
	"io/ioutil"
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

	// TODO -- needs status filtering
	projectsListCmd := &Command{
		Usage: "List your Cape projects",
		Examples: []*Example{
			{
				Example:     `cape projects list"`,
				Description: `Creates a project named "My Project" with the description "I will Cape the world!"`,
			},
		},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(projectsList),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	projectsUpdateCmd := &Command{
		Usage:     "Update a projects attributes",
		Arguments: []*Argument{ProjectLabelArg},
		Examples: []*Example{
			{
				Example:     `cape projects update --description "This project now does xyz" my-project`,
				Description: `Changes the description of "my-project"`,
			},
			{
				Example:     `cape projects update --from-spec spec.yaml my-project`,
				Description: `Updates the current project-spec to spec.yaml on my-project`,
			},
		},
		Command: &cli.Command{
			Name:   "update",
			Action: handleSessionOverrides(projectsUpdate),
			Flags: []cli.Flag{
				projectNameFlag(),
				projectDescriptionFlag(),
				projectSpecFlag(),
				clusterFlag(),
			},
		},
	}

	projectsCmd := &Command{
		Usage: "Commands for interacting with Cape projects",
		Command: &cli.Command{
			Name: "projects",
			Subcommands: []*cli.Command{
				projectsCreateCmd.Package(),
				projectsListCmd.Package(),
				projectsUpdateCmd.Package(),
			},
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
	desc, _ := Arguments(c.Context, ProjectDescriptionArg).(prims.Description)

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

func projectsList(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	projects, err := client.ListProjects(c.Context, nil)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	if len(projects) > 0 {
		header := []string{"Name", "Label", "Status", "Description"}
		body := make([][]string, len(projects))
		for i, p := range projects {
			descLength := len(p.Description)
			suffix := ""
			if descLength > 100 {
				suffix = "..."
				descLength = 100 - len(suffix)
			}

			desc := fmt.Sprintf("%s%s", p.Description.String()[:descLength], suffix)
			body[i] = []string{p.Name.String(), p.Label.String(), p.Status.String(), desc}
		}

		err = u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} project{{ . | pluralize \"s\"}}\n", len(projects))
}

func updateProjectSpec(c *cli.Context, specFile string) error {
	projectLabel := Arguments(c.Context, ProjectLabelArg).(primitives.Label)

	bytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return err
	}

	spec, err := primitives.ParseProjectSpecFile(bytes)
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	project, _, err := client.UpdateProjectSpec(c.Context, projectLabel, spec)
	if err != nil {
		return err
	}

	args := struct {
		ProjectName string
		FileName    string
	}{
		ProjectName: project.Name.String(),
		FileName:    specFile,
	}

	u := provider.UI(c.Context)
	return u.Template("Applied {{ .FileName | bold }} to {{ .ProjectName | bold }}\n", args)
}

func projectsUpdate(c *cli.Context) error {
	updateSpec := c.String("from-spec")
	if updateSpec != "" {
		return updateProjectSpec(c, updateSpec)
	}

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	label := Arguments(c.Context, ProjectLabelArg).(primitives.Label)

	var name *primitives.DisplayName
	var desc *prims.Description

	nameFlag, err := primitives.NewDisplayName(c.String("name"))
	if err != nil {
		return err
	}
	if nameFlag != "" {
		name = &nameFlag
	}
	descFlag, err := prims.NewDescription(c.String("description"))
	if err != nil {
		return err
	}
	if descFlag != "" {
		desc = &descFlag
	}

	project, err := client.UpdateProject(c.Context, nil, &label, name, desc)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	err = u.Template("Updated {{ . | bold }}!\n\n", project.Name.String())
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
