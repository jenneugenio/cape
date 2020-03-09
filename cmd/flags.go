package cmd

import "github.com/urfave/cli/v2"

func dbURLFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "db-url",
		Usage:    "The database url",
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

func serviceIDFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "service-id",
		Usage:   "Service ID of the component to run",
		EnvVars: []string{"CAPE_SERVICE_ID"},
	}
}
