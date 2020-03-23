package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"regexp"

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

type usage struct {
	Args     []argument
	Examples []example
}

type argument struct {
	Name        string
	Required    bool
	Description string
}

type example struct {
	Example     string
	Description string
}

var argNameRegex = regexp.MustCompile("^[a-z/-]{3,64}$")

func newArgument(name string, required bool, description string) argument {
	matched := argNameRegex.MatchString(name)

	if !matched {
		msg := fmt.Sprintf("Incorrect argument name %s: Argument names must only contain a-z, /, or -", name)
		panic(msg)
	}

	return argument{Name: name, Required: required, Description: description}
}

func (u usage) UsageText() string {
	var str string
	for i, e := range u.Examples {
		if i > 0 {
			str += "   "
		}
		str += fmt.Sprintf("%s\n", e.Description)
		str += fmt.Sprintf("\t\t\t%s", e.Example)

		if i < len(u.Examples)-1 {
			str += "\n\n"
		}
	}
	return str
}

func (u usage) UsageStr() string {
	var str string
	for i, arg := range u.Args {
		if i > 0 {
			str += " "
		}

		if !arg.Required {
			str += fmt.Sprintf("[%s]", arg.Name)
		} else {
			str += fmt.Sprintf("<%s>", arg.Name)
		}
	}
	return str
}

func (u usage) ArgsUsageText() string {
	var str string
	for i, arg := range u.Args {
		if i > 0 {
			str += "   "
		}

		str += fmt.Sprintf("%s - %s", arg.Name, arg.Description)

		if i < len(u.Args)-1 {
			str += "\n"
		}
	}
	return str
}
