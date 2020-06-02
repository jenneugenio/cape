package main

import (
	"github.com/capeprivacy/cape/primitives"
	gm "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestLabelArg(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		Name     string
		Param    string
		Input    string
		Expected primitives.Label
		Required bool
	}{
		{
			Name:     "Gets an expected label",
			Param:    "my-label",
			Input:    "a-label",
			Expected: primitives.Label("a-label"),
			Required: true,
		},

		{
			Name:     "Returns nil if no label is passed on an optional param",
			Param:    "my-label",
			Input:    "",
			Expected: primitives.Label(""),
			Required: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			arg := LabelArg(test.Param, test.Required)
			cmd := &Command{
				Arguments: []*Argument{arg},
				Usage:     "A cool app!",
				Command: &cli.Command{
					Name: "coolcmd",
					Action: func(c *cli.Context) error {
						arg, _ := Arguments(c.Context, arg).(primitives.Label)
						gm.Expect(arg).To(gm.Equal(test.Expected))
						return nil
					},
				},
			}

			app := cli.NewApp()
			app.Commands = []*cli.Command{cmd.Package()}

			err := app.Run([]string{"cape", "coolcmd", test.Input})
			gm.Expect(err).To(gm.BeNil())
		})
	}
}
