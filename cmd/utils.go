package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/go-base32"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func setupSignalWatcher(server *framework.Server, logger *zerolog.Logger) (*framework.SignalWatcher, error) {
	return framework.NewSignalWatcher(func(ctx context.Context, signal os.Signal) error {
		logger.Info().Msgf("Received signal %s, attempting to shutdown", signal)

		return server.Stop(ctx)
	}, func(_ context.Context, err error) {
		if err != nil {
			logger.Error().Err(err).Msg("Encountered error while trying to shutdown")
		}

		logger.Info().Msg("Shutdown")
		os.Exit(1)
	}, nil)
}

func exitHandler(c *cli.Context, err error) {
	// This is required because this function is called for every command
	// invocation independent of whether or not it errored.
	if err == nil {
		return
	}

	provider := GetProvider(c.Context)

	msg := ""
	switch e := err.(type) {
	case *errors.Error:
		// We don't want to display the cause string to the user, it's not
		// important and just adds clutter at this level of display.
		msg = strings.Join(e.Messages, ", ")
	default:
		msg = err.Error()
	}

	u := provider.UI(c.Context)
	// We don't check the error here because its done in `cmd/main.go`
	u.Notify(ui.Error, msg) //nolint: errcheck
}

func commandNotFound(c *cli.Context, command string) {
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	msg := "Oops! Unfortunately, the '%s %s' command doesn't exist. You can list all commands using '%s help'."

	// We don't check the error here because we immediately exit
	u.Notify(ui.Error, msg, cliName, command, cliName) // nolint: errcheck
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

func getName(c *cli.Context, question string) (models.Name, error) {
	nameStr := c.String("name")
	if nameStr != "" {
		return models.Name(nameStr), nil
	}

	validateName := func(input string) error {
		// TODO validate??
		return nil
	}

	msg := question
	if msg == "" {
		msg = "Please enter your name"
	}

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	nameStr, err := u.Question(msg, validateName)
	if err != nil {
		return "", err
	}

	return models.Name(nameStr), nil
}

func getEmail(c *cli.Context, in string) (primitives.Email, error) {
	if in != "" {
		return primitives.NewEmail(in)
	}

	emailStr := c.String("email")
	if emailStr != "" {
		return primitives.NewEmail(emailStr)
	}

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	out, err := u.Question("Please enter your email address", func(input string) error {
		_, err := primitives.NewEmail(input)
		return err
	})
	if err != nil {
		return primitives.Email{Email: ""}, err
	}

	return primitives.NewEmail(out)
}

func getPassword(c *cli.Context) (primitives.Password, error) {
	// Password can be nil as it's an _optional_ environment variable. Nil
	// cannot be cast to a primitives.Password so we need to check here to see
	// if the casting worked.
	pw, ok := EnvVariables(c.Context, capePasswordVar).(primitives.Password)
	if ok && pw != "" {
		return pw, nil
	}

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	// XXX: It'd be nice if we didn't need to do this weird type creation
	// manipulation. If we could just reuse the `.Validate()` function that'd
	// be awesome butthat's not how the promptui ValidatorFunc works!
	out, err := u.Secret("Please enter a password", func(input string) error {
		_, err := primitives.NewPassword(input)
		return err
	})
	if err != nil {
		return pw, err
	}

	return primitives.NewPassword(out)
}
