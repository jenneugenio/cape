package main

import (
	"fmt"
	"github.com/capeprivacy/cape/models"
	modelmigration "github.com/capeprivacy/cape/models/migration"
	"github.com/capeprivacy/cape/primitives"
	"sigs.k8s.io/yaml"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

func init() {
	// Project Subcommands

	projectsCreateCmd := &Command{
		Usage:     "Creates a project in Cape.",
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
		Usage: "List your Cape projects.",
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
		Usage:     "Update a projects attributes.",
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
				clusterFlag(),
				projectSpecFlag(),
			},
		},
	}

	projectsGetCmd := &Command{
		Usage:     "Get a details of a project.",
		Arguments: []*Argument{ProjectLabelArg},
		Examples: []*Example{
			{
				Example:     `cape projects get my-project`,
				Description: `Gets all of the details of a project including the active policy`,
			},
		},
		Command: &cli.Command{
			Name:   "get",
			Action: handleSessionOverrides(projectsGet),
		},
	}

	// Policy subcommands

	policyCreateCmd := &Command{
		Usage: "Suggest a policy change in a project.",
		Examples: []*Example{
			{
				Example:     `cape projects policy create my-project "My Title" "Description for my change" --from-spec policy.yaml`,
				Description: `Create a policy suggestion in the project my-project`,
			},
		},
		Arguments: []*Argument{ProjectLabelArg, SuggestionNameArg, SuggestionDescriptionArg},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(policyCreate),
			Flags: []cli.Flag{
				clusterFlag(),
				projectSpecFlag(),
			},
		},
	}

	policyListCmd := &Command{
		Usage: "List the policy suggestions in a project.",
		Examples: []*Example{
			{
				Example:     `cape projects policy list-suggestions my-project`,
				Description: `Lists all of the policy suggestions in "my-project"`,
			},
		},
		Arguments: []*Argument{ProjectLabelArg},
		Command: &cli.Command{
			Name:   "list-suggestions",
			Action: handleSessionOverrides(policyList),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	policyApproveCmd := &Command{
		Usage: "Approve a policy suggestion.",
		Examples: []*Example{
			{
				Example:     `cape projects policy approve <suggestion-id>`,
				Description: `Makes the provided policy suggestion active on the project`,
			},
		},
		Arguments: []*Argument{SuggestionIDArg},
		Command: &cli.Command{
			Name:   "approve",
			Action: handleSessionOverrides(policyApprove),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	policyRejectCmd := &Command{
		Usage: "Reject a policy suggestion.",
		Examples: []*Example{
			{
				Example:     `cape projects policy reject <suggestion-id>`,
				Description: `Makes the provided policy suggestion active on the project`,
			},
		},
		Arguments: []*Argument{SuggestionIDArg},
		Command: &cli.Command{
			Name:   "reject",
			Action: handleSessionOverrides(policyReject),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	policyGetCmd := &Command{
		Usage: "Get the active policy for a project.",
		Examples: []*Example{
			{
				Example:     `cape projects policy get my-project`,
				Description: `Get a detailed view of the projects active policy.`,
			},
		},
		Arguments: []*Argument{ProjectLabelArg},
		Command: &cli.Command{
			Name:   "get",
			Action: handleSessionOverrides(policyGet),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	policyGetSuggestionCmd := &Command{
		Usage: "Get more information about a policy suggestion.",
		Examples: []*Example{
			{
				Example:     `cape projects policy get-suggestion <suggestion-id>`,
				Description: `Get a detailed view of the provided policy suggestion.`,
			},
		},
		Arguments: []*Argument{SuggestionIDArg},
		Command: &cli.Command{
			Name:   "get-suggestion",
			Action: handleSessionOverrides(policyGetSuggestion),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	policyCmd := &Command{
		Usage: "Commands for interacting with Cape Project Policy.",
		Command: &cli.Command{
			Name: "policy",
			Subcommands: []*cli.Command{
				policyGetCmd.Package(),
				policyGetSuggestionCmd.Package(),
				policyApproveCmd.Package(),
				policyRejectCmd.Package(),
				policyListCmd.Package(),
				policyCreateCmd.Package(),
			},
		},
	}

	// Contributor subcommands

	contributorsAddCmd := &Command{
		Usage:     "Add a contributor to your project.",
		Arguments: []*Argument{UserEmailArg, ProjectLabelArg, RoleArg},
		Examples: []*Example{
			{
				Example:     `cape projects contributors add friend@cape.com my-project project-member`,
				Description: `Add the user with email friend@cape.com to the project "my-project" with the member role`,
			},
		},
		Command: &cli.Command{
			Name:   "add",
			Action: handleSessionOverrides(contributorsAdd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	contributorsListCmd := &Command{
		Usage:     "Lists the contributors in the provided project.",
		Arguments: []*Argument{ProjectLabelArg},
		Examples: []*Example{
			{
				Example:     `cape projects contributors list "my-project"`,
				Description: `List the contributors on the project "my-project"`,
			},
		},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(contributorsList),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	contributorsCmd := &Command{
		Usage: "Commands for interacting with Cape Project Contributors.",
		Command: &cli.Command{
			Name: "contributors",
			Subcommands: []*cli.Command{
				contributorsAddCmd.Package(),
				contributorsListCmd.Package(),
			},
		},
	}

	projectsCmd := &Command{
		Usage: "Commands for interacting with Cape projects.",
		Command: &cli.Command{
			Name: "projects",
			Subcommands: []*cli.Command{
				policyCmd.Package(),
				contributorsCmd.Package(),
				projectsCreateCmd.Package(),
				projectsListCmd.Package(),
				projectsUpdateCmd.Package(),
				projectsGetCmd.Package(),
			},
		},
	}

	commands = append(commands, policyCmd.Package(), contributorsCmd.Package(), projectsCmd.Package())
}

func contributorsAdd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	depProject := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	projectLabel := modelmigration.LabelFromPrimitive(depProject)

	depUser := Arguments(c.Context, UserEmailArg).(primitives.Email)
	userEmail := modelmigration.EmailFromPrimitive(depUser)

	roleLabel := Arguments(c.Context, RoleArg).(models.Label)

	projectResp, err := client.GetProject(c.Context, "", &projectLabel)
	if err != nil {
		return err
	}

	_, err = client.AddContributor(c.Context, *projectResp.Project, userEmail, roleLabel)
	if err != nil {
		return err
	}

	args := struct {
		User    string
		Project string
		Role    string
	}{
		User:    userEmail.String(),
		Project: projectLabel.String(),
		Role:    roleLabel.String(),
	}

	u := provider.UI(c.Context)
	return u.Template("Added user {{ .User | faded }} to project {{ .Project | faded }} with role {{ .Role | faded }}\n", args)
}

func contributorsList(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	dep := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	projectLabel := modelmigration.LabelFromPrimitive(dep)

	project, err := client.GetProject(c.Context, "", &projectLabel)
	if err != nil {
		return err
	}

	contributors, err := client.ListContributors(c.Context, *project.Project)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	if len(contributors) > 0 {
		header := []string{"Name", "Email"}
		body := make([][]string, len(contributors))
		for i, c := range contributors {
			body[i] = []string{c.User.Name.String(), c.User.Email.String()}
		}

		err = u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} contributor{{ . | pluralize \"s\"}}\n", len(contributors))
}

func projectsCreate(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	name := Arguments(c.Context, ProjectNameArg).(models.ProjectDisplayName)
	desc, _ := Arguments(c.Context, ProjectDescriptionArg).(models.ProjectDescription)

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

func policyCreate(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	dep := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	project := modelmigration.LabelFromPrimitive(dep)
	suggestionName := Arguments(c.Context, SuggestionNameArg).(models.ProjectDisplayName)
	suggestionDescription := Arguments(c.Context, SuggestionDescriptionArg).(models.ProjectDescription)

	specFile := c.String("from-spec")
	bytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return err
	}

	spec, err := models.ParseProjectSpecFile(bytes)
	if err != nil {
		return err
	}

	suggestion, err := client.SuggestPolicy(c.Context, project, suggestionName, suggestionDescription, spec)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	err = u.Template("Created policy suggestion for {{ . | faded }}\n", project.String())
	if err != nil {
		return err
	}

	details := ui.Details{
		"ID":          suggestion.ID,
		"Title":       suggestion.Title,
		"Description": suggestion.Description,
		"Status":      suggestion.State.String(),
	}

	return u.Details(details)
}

func policyList(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	deprecatedLabel := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	projectLabel := modelmigration.LabelFromPrimitive(deprecatedLabel)

	suggs, err := client.GetProjectSuggestions(c.Context, projectLabel)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	if len(suggs) > 0 {
		header := []string{"ID", "Title", "Status"}
		body := make([][]string, len(suggs))
		for i, s := range suggs {
			body[i] = []string{s.ID, s.Title, s.State.String()}
		}

		err = u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} policy suggestion{{ . | pluralize \"s\"}}\n", len(suggs))
}

func policyApprove(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	id := Arguments(c.Context, SuggestionIDArg).(string)
	s := models.Suggestion{ID: id}
	err = client.ApproveSuggestion(c.Context, s)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("\nPolicy is now {{ . | faded }}\n", "active")
}

func policyReject(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	id := Arguments(c.Context, SuggestionIDArg).(string)
	s := models.Suggestion{ID: id}
	err = client.RejectSuggestion(c.Context, s)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("\nPolicy suggestion rejected\n", nil)
}

func policyGetSuggestion(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	id := Arguments(c.Context, SuggestionIDArg).(string)
	s, err := client.GetProjectSuggestion(c.Context, id)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	details := ui.Details{
		"Name":        s.Title,
		"Description": s.Description,
		"State":       s.State.String(),
	}

	err = u.Details(details)
	if err != nil {
		return err
	}

	rules, err := yaml.Marshal(s.Policy.Rules)
	if err != nil {
		return err
	}

	transformations, err := yaml.Marshal(s.Policy.Transformations)
	if err != nil {
		return err
	}

	args := struct {
		Rules           string
		Transformations string
	}{
		string(rules), string(transformations),
	}

	return u.Template("policy:\n{{ .Rules }}\ntransformations:{{ .Transformations }}\n", args)
}

func projectsList(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	projects, err := client.ListProjects(c.Context, models.Any)
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

func policyGet(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	deprecatedLabel := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	label := modelmigration.LabelFromPrimitive(deprecatedLabel)

	project, err := client.GetProject(c.Context, "", &label)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)

	if project.Policy == nil {
		return u.Template("Project has no active policy\n", nil)
	}

	rules, err := yaml.Marshal(project.Policy.Rules)
	if err != nil {
		return err
	}

	transformations, err := yaml.Marshal(project.Policy.Transformations)
	if err != nil {
		return err
	}

	args := struct {
		Rules           string
		Transformations string
	}{
		string(rules), string(transformations),
	}

	return u.Template("policy:\n{{ .Rules }}\ntransformations:{{ .Transformations }}\n", args)
}

func updateProjectSpec(c *cli.Context, specFile string) error {
	deprecatedLabel := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	projectLabel := modelmigration.LabelFromPrimitive(deprecatedLabel)

	bytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return err
	}

	spec, err := models.ParseProjectSpecFile(bytes)
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

	deprecatedLabel := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	label := modelmigration.LabelFromPrimitive(deprecatedLabel)

	var name *models.ProjectDisplayName
	var desc *models.ProjectDescription

	nameFlag := models.ProjectDisplayName(c.String("name"))
	if nameFlag != "" {
		name = &nameFlag
	}
	descFlag := models.ProjectDescription(c.String("description"))
	if descFlag != "" {
		desc = &descFlag
	}

	project, err := client.UpdateProject(c.Context, "", &label, name, desc)
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

func projectsGet(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	deprecatedLabel := Arguments(c.Context, ProjectLabelArg).(primitives.Label)
	label := modelmigration.LabelFromPrimitive(deprecatedLabel)

	project, err := client.GetProject(c.Context, "", &label)
	if err != nil {
		return err
	}

	details := ui.Details{
		"Name":        project.Name.String(),
		"Description": project.Description.String(),
		"Label":       project.Label.String(),
		"Status":      project.Status.String(),
	}

	u := provider.UI(c.Context)
	err = u.Details(details)
	if err != nil {
		return err
	}

	err = u.Template("{{ . | faded }}\n", "Contributors")
	if err != nil {
		return err
	}

	header := ui.TableHeader{"Name", "Email", "Role"}
	body := ui.TableBody{}
	for _, c := range project.Contributors {
		body = append(body, []string{c.User.Name.String(), c.User.Email.String(), c.Role.Label.String()})
	}

	err = u.Table(header, body)
	if err != nil {
		return err
	}

	// Print the policy if there is one, but there may not be, in which case we are done
	if project.Policy == nil {
		return nil
	}

	err = u.Template("\n{{ . | faded }}\n", "Policy")
	if err != nil {
		return err
	}

	rules, err := yaml.Marshal(project.Policy.Rules)
	if err != nil {
		return err
	}

	transformations, err := yaml.Marshal(project.Policy.Transformations)
	if err != nil {
		return err
	}

	args := struct {
		Rules           string
		Transformations string
	}{
		string(rules), string(transformations),
	}

	return u.Template("policy:\n{{ .Rules }}\ntransformations:{{ .Transformations }}\n", args)
}
