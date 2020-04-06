package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/version"
)

var commands []*cli.Command

const cliName = "cape"

var commandHelpTemplate = `NAME:
   {{.Name}} - {{.Description}}

USAGE:
   {{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} [arguments...]{{if .Aliases}}

ALIASES:
	 {{range $i, $v := .Aliases}}{{if $i}}, {{$v}}{{else}}{{$v}}{{end}}{{end}}{{end}}{{if .ArgsUsage}}   
   {{.ArgsUsage}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{else}}
{{end}}{{if .Category}}
CATEGORY:
   {{.Category}}
{{end}}{{if .UsageText}}
EXAMPLES:
	 {{.UsageText}}{{end}}
`

//nolint: lll
var appHelpTemplate = fmt.Sprintf(`NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

See '%s help <command>' to read about a specific command.

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`, cliName)

func main() {
	cli.CommandHelpTemplate = commandHelpTemplate
	cli.AppHelpTemplate = appHelpTemplate

	cli.HelpFlag = helpFlag()
	cli.VersionFlag = versionFlag()
	cli.VersionPrinter = func(c *cli.Context) {
		err := versionCmd(c)
		if err != nil {
			exitHandler(c, err)
		}
	}

	app := cli.NewApp()
	app.Name = cliName
	app.HelpName = cliName
	app.Usage = "Cape is used to manage access to your sensitive data"
	app.Version = version.Version
	app.Commands = commands
	app.EnableBashCompletion = true
	app.Copyright = "(c) 2020 Cape, Inc."

	// Before runs our global middleware for all commands including the command
	// not found middleware
	app.Before = cli.BeforeFunc(retrieveConfig)
	app.CommandNotFound = commandNotFound
	app.ExitErrHandler = exitHandler

	err := app.Run(os.Args)
	if err != nil {
		// Errors are handled by the ExitErrHandler
		os.Exit(1)
	}
}
