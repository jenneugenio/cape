package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/version"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		err := versionCmd(c)
		if err != nil {
			errorPrinter(err)
		}
	}

	cli.VersionFlag = versionFlag()
}

func versionCmd(c *cli.Context) error {
	fmt.Printf("Cape CLI - Version: %s Date: %s\n", version.Version, version.BuildDate)
	return nil
}
