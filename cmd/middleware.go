package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/cmd/config"
	"github.com/dropoutlabs/cape/cmd/ui"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// ContextKey is a type alias used for storing data in a context
type ContextKey string

const (
	// ConfigContextKey is the name of the key storing configuration on the context
	ConfigContextKey ContextKey = "config"

	// ArgumentContextKey is the name of the key storing argument values for a command
	ArgumentContextKey ContextKey = "arguments"

	// UIContextKey is the name of the key storing the ui.UI struct for a command
	UIContextKey ContextKey = "ui"

	// SessionContextKey is the name of key storing the config.Session struct for a command
	SessionContextKey ContextKey = "session"
)

// Config returns the config object stored on the context
func Config(ctx context.Context) *config.Config {
	cfg := ctx.Value(ConfigContextKey)
	if cfg == nil {
		panic("config not available on context")
	}

	return cfg.(*config.Config)
}

// Arguments returns the ArgumentValues object stored on the context
func Arguments(ctx context.Context) ArgumentValues {
	args := ctx.Value(ArgumentContextKey)
	if args == nil {
		panic("argument values not available on context")
	}

	return args.(ArgumentValues)
}

// UI returns the UI object stored on the context
func UI(ctx context.Context) *ui.UI {
	u := ctx.Value(UIContextKey)
	if u == nil {
		panic("ui not available on context")
	}

	return u.(*ui.UI)
}

// Session returns the session object stored on the context
func Session(ctx context.Context) *config.Session {
	s := ctx.Value(SessionContextKey)
	if s == nil {
		panic("session not available on context")
	}

	return s.(*config.Session)
}

func retrieveConfig(next cli.ActionFunc) cli.ActionFunc {
	return cli.ActionFunc(func(c *cli.Context) error {
		cfg, err := config.Parse()
		if err != nil {
			return err
		}

		u, err := ui.NewUI(cfg)
		if err != nil {
			return err
		}

		c.Context = context.WithValue(c.Context, UIContextKey, u)
		c.Context = context.WithValue(c.Context, ConfigContextKey, cfg)

		return next(c)
	})
}

func processArguments(cmd *Command, next cli.ActionFunc) cli.ActionFunc {
	return cli.ActionFunc(func(c *cli.Context) error {
		values := ArgumentValues{}
		for i, arg := range cmd.Arguments {
			input := c.Args().Get(i)
			if input == "" && arg.Required {
				return errors.New(MissingArgCause, "The argument %s is required, but was not provided", arg.Name)
			}

			if input == "" {
				return nil
			}

			value, err := arg.Processor(input)
			if err != nil {
				return err
			}

			values[arg.Name] = value
		}

		c.Context = context.WithValue(c.Context, ArgumentContextKey, values)
		return next(c)
	})
}

func handleSessionOverrides(next cli.ActionFunc) cli.ActionFunc {
	return cli.ActionFunc(func(c *cli.Context) error {
		cfg := Config(c.Context)
		clusterStr := c.String("cluster")

		var cluster *config.Cluster
		if clusterStr != "" {
			c, err := cfg.GetCluster(clusterStr)
			if err != nil {
				return err
			}

			cluster = c
		} else {
			c, err := cfg.Cluster()
			if err != nil {
				return err
			}

			cluster = c
		}

		session := config.NewSession(cfg, cluster)
		c.Context = context.WithValue(c.Context, SessionContextKey, session)
		return next(c)
	})
}
