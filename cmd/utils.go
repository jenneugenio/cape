package cmd

import (
	"github.com/urfave/cli/v2"
)

func getServiceID(c *cli.Context) string {
	serviceID := c.String("service-id")
	if serviceID == "" {
		serviceID = "unknown"
	}

	return serviceID
}
