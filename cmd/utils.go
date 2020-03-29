package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/go-base32"
	"github.com/urfave/cli/v2"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
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

func getInstanceID(c *cli.Context, serviceType string) (primitives.Label, error) {
	instanceID := c.String("instance-id")
	if instanceID != "" {
		return formatInstanceID(instanceID, serviceType)
	}

	source := make([]byte, 4)
	_, err := rand.Read(source)
	if err != nil {
		return "", err
	}

	return formatInstanceID(base32.EncodeToString(source), serviceType)
}

func formatInstanceID(serviceType, instanceID string) (primitives.Label, error) {
	return primitives.NewLabel(fmt.Sprintf("cape-%s-%s", serviceType, instanceID))
}

func getName(c *cli.Context, question string) (primitives.Name, error) {
	validateName := func(input string) error {
		_, err := primitives.NewName(input)
		if err != nil {
			return err
		}

		return nil
	}

	msg := question
	if msg == "" {
		msg = "Please enter your name"
	}

	ui := UI(c.Context)
	nameStr, err := ui.Question(msg, validateName)
	if err != nil {
		return primitives.Name(""), err
	}

	return primitives.NewName(nameStr)
}

func getEmail(c *cli.Context, in string) (primitives.Email, error) {
	if in != "" {
		return primitives.NewEmail(in)
	}

	ui := UI(c.Context)
	out, err := ui.Question("Please enter your email address", func(input string) error {
		_, err := primitives.NewEmail(input)
		return err
	})
	if err != nil {
		return primitives.Email(""), err
	}

	return primitives.NewEmail(out)
}

func getPassword(c *cli.Context) (primitives.Password, error) {
	envVars := EnvVariables(c.Context)
	ui := UI(c.Context)

	pw, ok := envVars["CAPE_PASSWORD"].(primitives.Password)
	if ok {
		return pw, nil
	}

	// XXX: It'd be nice if we didn't need to do this weird type creation
	// manipulation. If we could just reuse the `.Validate()` function that'd
	// be awesome butthat's not how the promptui ValidatorFunc works!
	out, err := ui.Secret("Please enter a password", func(input string) error {
		_, err := primitives.NewPassword(input)
		return err
	})
	if err != nil {
		return pw, err
	}

	return primitives.NewPassword(out)
}

func getConfirmedPassword(c *cli.Context) (primitives.Password, error) {
	ui := UI(c.Context)

	empty := primitives.Password("")
	password, err := ui.Secret("Please enter a password", func(input string) error {
		_, err := primitives.NewPassword(input)
		return err
	})
	if err != nil {
		return empty, err
	}

	_, err = ui.Secret("Please confirm the password you entered", func(input string) error {
		out, err := primitives.NewPassword(input)
		if err != nil {
			return err
		}

		if password != out.String() {
			return errors.New(PasswordNoMatch, "Does not previously provided password")
		}

		return nil
	})
	if err != nil {
		return empty, err
	}

	return primitives.NewPassword(password)
}
