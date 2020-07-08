package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/logging"
)

func portFlag(name string, value int) cli.Flag {
	return &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Usage:   fmt.Sprintf("The `PORT` the %s will attempt to listen on", name),
		Value:   value,
		EnvVars: []string{"CAPE_PORT"},
	}
}

func versionFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Display the current version of Cape",
	}
}

func helpFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Display documentation and examples for this command",
	}
}

func yesFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    "yes",
		Aliases: []string{"y"},
		Usage:   "If specified, the user will not be prompted to confirm their action before proceeding",
		Value:   false,
		EnvVars: []string{"CAPE_YES"},
	}
}

func useClusterFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    "use",
		Aliases: []string{"u"},
		Usage:   "If provided, the cluster being created will also be set as the current cluster",
		EnvVars: []string{"CAPE_USE"},
	}
}

func instanceIDFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "instance-id",
		Usage:   "An identifier to provide for uniquely identifying this process. One will be generated if not provided.",
		EnvVars: []string{"CAPE_INSTANCE_ID"},
	}
}

func loggingTypeFlag() cli.Flag {
	options := []string{}
	for _, str := range logging.Types() {
		options = append(options, str)
	}

	str := "The type of logger to use for disptaching logs (options: %s)"
	usage := fmt.Sprintf(str, strings.Join(options, ", "))
	return &cli.StringFlag{
		Name:    "logger",
		Usage:   usage,
		Value:   logging.JSONType.String(),
		EnvVars: []string{"CAPE_LOGGING_TYPE"},
	}
}

func loggingLevelFlag() cli.Flag {
	str := "The level of log to report when dispatching logs (options: %s)"
	usage := fmt.Sprintf(str, strings.Join(logging.Levels(), ", "))
	return &cli.StringFlag{
		Name:    "log-level",
		Usage:   usage,
		Value:   logging.DefaultLevel,
		EnvVars: []string{"CAPE_LOGGING_LEVEL"},
	}
}

func emailFlag() cli.Flag {
	usage := "The email used to log into the cluster"
	return &cli.StringFlag{
		Name:    "email",
		Usage:   usage,
		EnvVars: []string{"CAPE_EMAIL"},
	}
}

func nameFlag() cli.Flag {
	usage := "The name of a user"
	return &cli.StringFlag{
		Name:    "name",
		Usage:   usage,
		EnvVars: []string{"CAPE_NAME"},
	}
}

func clusterFlag() cli.Flag {
	usage := "The cluster to login to"
	return &cli.StringFlag{
		Name:    "cluster",
		Usage:   usage,
		EnvVars: []string{"CAPE_CLUSTER"},
	}
}

func membersFlag() cli.Flag {
	usage := "Members to assign the specified role to"
	return &cli.StringSliceFlag{
		Name:    "member",
		Aliases: []string{"m"},
		Usage:   usage,
		EnvVars: []string{"CAPE_MEMBERS"},
	}
}

func fileFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "from-file",
		Usage:   "Loads a policy from a file and creates it",
		EnvVars: []string{"CAPE_FILEPATH"},
	}
}

func configFilesFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "Configuration file locations. Multiple files can be passed where the previous file is overridden by the next.",
		EnvVars: []string{"CAPE_COORDINATOR_CONFIG"},
	}
}

func configFileOutFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "out",
		Aliases: []string{"o"},
		Usage:   "Where to output the config file",
		Value:   "config.yaml",
		EnvVars: []string{"CAPE_CONFIG_OUTPUT"},
	}
}

func formatFlag() cli.Flag {
	var options []string
	for _, str := range FormatTypes() {
		options = append(options, str)
	}

	str := "The format of the configuration file (options: %s)"
	usage := fmt.Sprintf(str, strings.Join(options, ", "))
	return &cli.StringFlag{
		Name:    "format",
		Usage:   usage,
		EnvVars: []string{"CAPE_CONFIG_FORMAT"},
	}
}

func projectNameFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "name",
		Usage:   "The name of the project",
		EnvVars: []string{"CAPE_PROJECT_NAME"},
	}
}

func projectDescriptionFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "description",
		Usage:   "The description of the project",
		EnvVars: []string{"CAPE_PROJECT_DESCRIPTION"},
	}
}

func projectSpecFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "from-spec",
		Usage:   "The spec file you wish to use",
		EnvVars: []string{"CAPE_PROJECT_SPEC"},
	}
}
