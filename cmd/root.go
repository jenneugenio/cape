package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

// Execute creates the CLI App and runs its
func Execute() {
	app := cli.NewApp()
	app.Name = "cape"
	app.HelpName = "cape"
	app.Usage = "Cape is used to manage access to your sensitive data"
	app.Version = "0.0.1"
	app.Commands = commands
	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
