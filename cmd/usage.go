package main

import (
	"fmt"
	"regexp"

	"github.com/urfave/cli/v2"
)

var argNameRegex = regexp.MustCompile("^[a-z/-]{3,64}$")

// ArgumentValues represents a map of arguments to their values
type ArgumentValues map[string]interface{}

// ArgumentProcessorFunc represents a function that given a value, validates
// the value, parses the string, and then returns the value to store in the
// ArgumentValues map for this invocation of the Command.
type ArgumentProcessorFunc func(string) (interface{}, error)

// Argument represents a command line argument
type Argument struct {
	Name        string
	Required    bool
	Description string
	Processor   ArgumentProcessorFunc
}

// String returns a string representation of the argument
func (a *Argument) String() string {
	if a.Required {
		return fmt.Sprintf("<%s>", a.Name)
	}

	return fmt.Sprintf("[%s]", a.Name)
}

// Returns a string specifying the argument's usage
func (a *Argument) Usage() string {
	required := ""
	if a.Required {
		required = " (Required)"
	}

	return fmt.Sprintf("%s%s", a.Description, required)
}

// Example represents a example of how to use a command
type Example struct {
	Example     string
	Description string
}

// Command is a wrapper around urfave/cli.Command
type Command struct {
	Arguments   []*Argument
	Examples    []*Example
	Usage       string
	Description string
	Command     *cli.Command
}

// UsageText compiles the usage information for the Command.UsageText field
func (c *Command) UsageText() string {
	var str string
	for i, e := range c.Examples {
		if i > 0 {
			str += "   "
		}
		str += fmt.Sprintf("%s\n", e.Description)
		str += fmt.Sprintf("\t\t\t%s", e.Example)

		if i < len(c.Examples)-1 {
			str += "\n\n"
		}
	}
	return str
}

// ArgsUsageText compiles the usage information for the Command.ArgsUsage field
func (c *Command) ArgsUsageText() string {
	var str string
	for i, arg := range c.Arguments {
		if i > 0 {
			str += "   "
		}

		required := ""
		if arg.Required {
			required = " (required)"
		}
		str += fmt.Sprintf("%s\t%s%s", arg.String(), arg.Description, required)

		if i < len(c.Arguments)-1 {
			str += "\n"
		}
	}
	return str
}

// Update manipulates the provided cli.Command to set the appropriate values.
func (c *Command) Package() *cli.Command {
	cmd := c.Command

	// Since we're using literal struct declarations to define args, usage,
	// etc., we have to leverage the `.Package()` func to test that any
	// required properties are set and all properties are also valid.
	//
	// Ideally, these would be checked at compile time but it's not entirely
	// possible. Perhaps in future we can introduce a linter :)
	if c.Usage == "" {
		panic("All commands must have usage text.")
	}

	for _, arg := range c.Arguments {
		if !argNameRegex.MatchString(arg.Name) {
			msg := fmt.Sprintf("Incorrect argument name %s: Argument names must only contain a-z, or -", arg.Name)
			panic(msg)
		}
	}

	cmd.Description = c.Description
	cmd.Usage = c.Usage
	cmd.ArgsUsage = c.ArgsUsageText()
	cmd.UsageText = c.UsageText()

	// Apply our middleware!
	cmd.Action = retrieveConfig(processArguments(c, cmd.Action))

	return cmd
}
