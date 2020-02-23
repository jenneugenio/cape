package cmd

import (
	"os"
)

func getServiceID() string {
	serviceID := os.Getenv("CAPE_SERVICE_ID")
	if serviceID == "" {
		serviceID = "unknown"
	}

	return serviceID
}
