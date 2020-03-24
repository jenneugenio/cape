package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/cmd/config"
)

// ContextKey is a type alias used for storing data in a context
type ContextKey string

const (
	// ConfigContextKey is the name of the key storing configuration on the context
	ConfigContextKey ContextKey = "config"
)

// Config returns the config object stored on the context
func Config(ctx context.Context) *config.Config {
	cfg := ctx.Value(ConfigContextKey)
	if cfg == nil {
		panic("config not available on context")
	}

	return cfg.(*config.Config)
}

func retrieveConfig(next cli.ActionFunc) cli.ActionFunc {
	return cli.ActionFunc(func(c *cli.Context) error {
		cfg, err := config.Parse()
		if err != nil {
			return err
		}

		c.Context = context.WithValue(c.Context, ConfigContextKey, cfg)
		return next(c)
	})
}
