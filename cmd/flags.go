package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/logging"
	"github.com/dropoutlabs/cape/primitives"
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
		Name:  "version, v",
		Usage: "Display the current version of Cape",
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

func dbURLFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "db-url",
		Usage:    "The database `URL`",
		EnvVars:  []string{"CAPE_DB_URL"},
		Required: true,
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

func clusterFlag() cli.Flag {
	usage := "The cluster to login to"
	return &cli.StringFlag{
		Name:    "cluster",
		Usage:   usage,
		EnvVars: []string{"CAPE_CLUSTER"},
	}
}

func serviceTypeFlag() cli.Flag {
	options := []string{}
	for _, str := range primitives.ServiceTypes() {
		options = append(options, str)
	}

	str := "The type of the service (options: %s)"
	usage := fmt.Sprintf(str, strings.Join(options, ", "))
	return &cli.StringFlag{
		Name:    "type",
		Usage:   usage,
		Value:   primitives.UserServiceType.String(),
		EnvVars: []string{"CAPE_TYPE"},
	}
}

func dataConnectorEndpointFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "endpoint",
		Usage:   "The endpoint to connect to a data connector. Must be supplied when creating a data-connector",
		EnvVars: []string{"CAPE_ENDPOINT"},
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

func linkFlag() cli.Flag {
	usage := "Members to assign the specified role to"
	return &cli.StringFlag{
		Name:    "link",
		Aliases: []string{"l"},
		Usage:   usage,
		EnvVars: []string{"CAPE_LINK"},
	}
}

func fileFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "from-file",
		Usage:   "Loads a policy from a file and creates it",
		EnvVars: []string{"CAPE_FILEPATH"},
	}
}
