package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/go-base32"
	"github.com/urfave/cli/v2"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

func errorPrinter(err error) {
	msg := ""
	switch e := err.(type) {
	case *errors.Error:
		// We don't want to display the cause string to the user, it's not
		// important and just adds clutter at this level of display.
		msg = strings.Join(e.Messages, ", ")
	default:
		msg = err.Error()
	}

	fmt.Fprintf(os.Stderr, "\nError: %s\n", msg)
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
