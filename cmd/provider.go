package main

import (
	"context"
	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
)

// Provider is an interface that gives different cape commands/components
// access to other components that they may need
type Provider interface {
	Client(ctx context.Context) (*coordinator.Client, error)
	UI(ctx context.Context) ui.UI
}

// AppProvider implements the Provider interface for the "normal" app. (ie not in testing)
// It can provide needed components to the commands that require them
type AppProvider struct {
}

func (a *AppProvider) transport(ctx context.Context) (coordinator.Transport, error) {
	cfgSession := Session(ctx)
	cluster, err := cfgSession.Cluster()

	if err != nil {
		return nil, err
	}

	return cluster.Transport()
}

// Client implements Client on the Provider interface
func (a *AppProvider) Client(ctx context.Context) (*coordinator.Client, error) {
	transport, err := a.transport(ctx)
	if err != nil {
		return nil, err
	}

	return coordinator.NewClient(transport), nil
}

// UI implements UI on the provider interface
func (a *AppProvider) UI(ctx context.Context) ui.UI {
	cfg := Config(ctx)
	return ui.NewStdout(cfg.UI.Colors, cfg.UI.Animations)
}

// NewAppProvider returns a provider
func NewAppProvider() Provider {
	return &AppProvider{}
}
