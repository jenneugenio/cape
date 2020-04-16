package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/version"
)

func versionCmd(c *cli.Context) error {
	fmt.Printf("Cape CLI - Version: %s Date: %s\n", version.Version, version.BuildDate)
	return nil
}
