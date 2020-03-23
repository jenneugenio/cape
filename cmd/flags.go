package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/logging"
)

func versionFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  "version, v",
		Usage: "Display the current version of Cape",
	}
}

func yesFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:    "yes, y",
		Usage:   "If specified, the user will not be prompted to confirm their action before proceeding",
		Value:   false,
		EnvVars: []string{"CAPE_YES"},
	}
}

func useClusterFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "use, u",
		Usage:   "Specify the `LABEL` of the Cape cluster to use",
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

func dbPasswordFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "db-password",
		Usage:   "The database password",
		EnvVars: []string{"CAPE_DB_PASSWORD"},
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
