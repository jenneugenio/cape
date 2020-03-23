package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/version"
)

var commands []*cli.Command

func main() {
	app := cli.NewApp()
	app.Name = "cape"
	app.HelpName = "cape"
	app.Usage = "Cape is used to manage access to your sensitive data"
	app.Version = version.Version
	app.Commands = commands
	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		errorPrinter(err)
	}
}
