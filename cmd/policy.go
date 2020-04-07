package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dropoutlabs/cape/primitives"
	"github.com/urfave/cli/v2"
)

func init() {
	attachCmd := &Command{
		Usage: "Attaches a policy to the given role",
		Examples: []*Example{
			{
				Example:     "cape attach admin-policy admin",
				Description: "This attaches a policy 'admin-policy' to the role 'admin'.",
			},
			{
				Example:     "cape attach --from-file admin-policy.yaml admin-policy admin",
				Description: "This reads the policy from a file, creates it and then attaches it to the role 'admin'.",
			},
		},
		Arguments: []*Argument{PolicyLabelArg, RoleLabelArg},
		Command: &cli.Command{
			Name:   "attach",
			Action: handleSessionOverrides(policyAttachCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				fileFlag(),
			},
		},
	}

	detachCmd := &Command{
		Usage: "Detaches a policy from the given role",
		Examples: []*Example{
			{
				Example:     "cape detach admin-policy admin",
				Description: "This detaches the policy 'admin-policy' from the role 'admin'.",
			},
			{
				Example:     "cape detach --yes admin-policy admin",
				Description: "This detaches the policy 'admin-policy' from the role 'admin' skipping the confirm prompt.",
			},
		},
		Arguments: []*Argument{PolicyLabelArg, RoleLabelArg},
		Command: &cli.Command{
			Name:   "detach",
			Action: handleSessionOverrides(policyDetachCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				yesFlag(),
			},
		},
	}

	listCmd := &Command{
		Usage: "Lists all the policies on the cluster",
		Examples: []*Example{
			{
				Example:     "cape policies list",
				Description: "Lists all policies",
			},
		},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(policiesListCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	policiesCmd := &Command{
		Usage: "Commands for querying information about policies and modifying them",
		Command: &cli.Command{
			Name: "policies",
			Subcommands: []*cli.Command{
				attachCmd.Package(),
				detachCmd.Package(),
				listCmd.Package(),
			},
		},
	}

	commands = append(commands, policiesCmd.Package(), attachCmd.Package(), detachCmd.Package())
}

func policyAttachCmd(c *cli.Context) error {
	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	roleLabel := Arguments(c.Context, RoleLabelArg).(primitives.Label)
	policyLabel := Arguments(c.Context, PolicyLabelArg).(primitives.Label)

	file := c.String("from-file")

	role, err := client.GetRoleByLabel(c.Context, roleLabel)
	if err != nil {
		return err
	}

	var policy *primitives.Policy
	if file != "" {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		policyInput, err := primitives.ParsePolicy(data)
		if err != nil {
			return err
		}

		p, err := client.CreatePolicy(c.Context, policyInput)
		if err != nil {
			return err
		}

		policy = p
	} else {
		p, err := client.GetPolicyByLabel(c.Context, policyLabel)
		if err != nil {
			return err
		}

		policy = p
	}

	_, err = client.AttachPolicy(c.Context, policy.ID, role.ID)
	if err != nil {
		return err
	}

	fmt.Printf("'%s' policy has been attached to the '%s' role.\n", policyLabel, roleLabel)

	return nil
}

func policyDetachCmd(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	u := UI(c.Context)

	cfgSession := Session(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	roleLabel := Arguments(c.Context, RoleLabelArg).(primitives.Label)
	policyLabel := Arguments(c.Context, PolicyLabelArg).(primitives.Label)

	if !skipConfirm {
		err := u.Confirm(fmt.Sprintf("Do you really want to detach policy %s from role %s?", policyLabel, roleLabel))
		if err != nil {
			return err
		}
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	role, err := client.GetRoleByLabel(c.Context, roleLabel)
	if err != nil {
		return err
	}

	policy, err := client.GetPolicyByLabel(c.Context, policyLabel)
	if err != nil {
		return err
	}

	err = client.DetachPolicy(c.Context, policy.ID, role.ID)
	if err != nil {
		return err
	}

	fmt.Printf("The policy '%s' has been detached from the role '%s'.n", policyLabel, roleLabel)

	return nil
}

func policiesListCmd(c *cli.Context) error {
	ui := UI(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	policies, err := client.ListPolicies(c.Context)
	if err != nil {
		return err
	}

	header := []string{"Label"}
	body := make([][]string, len(policies))
	for i, p := range policies {
		body[i] = []string{p.Label.String()}
	}

	return ui.Table(header, body)
}
