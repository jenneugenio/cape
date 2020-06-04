package main

import (
	"fmt"
	"regexp"

	"github.com/urfave/cli/v2"
)

var argNameRegex = regexp.MustCompile("^[a-z/-]{3,64}$")
var envVarNameRegex = regexp.MustCompile("^[A-Z][A-Z_]{2,64}$")

// EnvVar represents a command line environment that is not a flag. Only use
// this type *if* you do not a flag. In accordance with our style guide, this
// should only be used for variables that contain secrets or sensitive
// information.
type EnvVar struct {
	Name        string
	Required    bool
	Description string
	Processor   VariableProcessorFunc
}

// String returns the string representation of the environment variable
func (e *EnvVar) String() string {
	// XXX: Should we auto namespace variables?!
	return e.Name
}

// Usage returns a string specifying the environnment variables usage
func (e *EnvVar) Usage() string {
	required := ""
	if e.Required {
		required = " (Required)"
	}

	return fmt.Sprintf("%s%s", e.Description, required)
}

// ArgumentValues contains values from arguments specified for the command
type ArgumentValues map[*Argument]interface{}

// EnvVarValues contains values from environment variables specified for a
// command
type EnvVarValues map[*EnvVar]interface{}

// VariableProcessorFunc represents a function that given a value, validates
// the value, parses the string, and then returns the value to store in the
// VariableValues map for this invocation of the Command.
type VariableProcessorFunc func(string) (interface{}, error)

// Argument represents a command line argument
type Argument struct {
	Name        string
	Required    bool
	Description string
	Processor   VariableProcessorFunc
}

// String returns a string representation of the argument
func (a *Argument) String() string {
	return a.Name
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
	Variables   []*EnvVar
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
		str += fmt.Sprintf("%s\n", e.Example)
		str += fmt.Sprintf("\t\t\t- %s", e.Description)

		if i < len(c.Examples)-1 {
			str += "\n\n"
		}
	}
	return str
}

// ArgsUsageText compiles the usage information for the Command.ArgsUsage field
// which contains information about the environment variales and arguments
// accepted by this command.
func (c *Command) ArgsUsageText() string {
	str := ""

	if len(c.Variables) > 0 {
		str += "\nENVIRONMENT VARIABLES:\n"
		for i, e := range c.Variables {
			required := ""
			if e.Required {
				required = " (required)"
			}

			str += fmt.Sprintf("  %s\t%s%s", e.String(), e.Description, required)
			if i < len(c.Variables)-1 {
				str += "\n"
			}
		}

		if len(c.Arguments) > 0 {
			str += "\n"
		}
	}

	if len(c.Arguments) > 0 {
		str += "\nARGUMENTS:\n"
		for i, arg := range c.Arguments {
			required := ""
			if arg.Required {
				required = " (required)"
			}

			str += fmt.Sprintf("  %s\t%s%s", arg.String(), arg.Description, required)
			if i < len(c.Arguments)-1 {
				str += "\n"
			}
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

	for _, e := range c.Variables {
		if !envVarNameRegex.MatchString(e.Name) {
			msg := fmt.Sprintf("Incorrect variable name: %s: Env vars must"+
				" be all caps and only contain a-z, and _.", e.Name)
			panic(msg)
		}
	}

	for _, arg := range c.Arguments {
		if !argNameRegex.MatchString(arg.Name) {
			msg := fmt.Sprintf("Incorrect argument name %s: Argument names must only contain a-z, or -", arg.Name)
			panic(msg)
		}
	}

	// We fall back to the usage string if a description is not provided.
	description := c.Description
	if description == "" {
		description = c.Usage
	}

	cmd.Description = description
	cmd.Usage = c.Usage
	cmd.ArgsUsage = c.ArgsUsageText()
	cmd.UsageText = c.UsageText()

	// Apply our middleware!
	if cmd.Action != nil {
		cmd.Action = processVariables(c, cmd.Action)
	}

	return cmd
}
