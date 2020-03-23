package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/manifoldco/go-base32"
	"github.com/urfave/cli/v2"
)

func errorPrinter(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func getInstanceID(c *cli.Context, serviceType string) (string, error) {
	instanceID := c.String("instance-id")
	if instanceID != "" {
		return formatInstanceID(instanceID, serviceType), nil
	}

	source := make([]byte, 4)
	_, err := rand.Read(source)
	if err != nil {
		return "", err
	}

	return formatInstanceID(base32.EncodeToString(source), serviceType), nil
}

func formatInstanceID(serviceType, instanceID string) string {
	return fmt.Sprintf("cape-%s-%s", serviceType, instanceID)
}
