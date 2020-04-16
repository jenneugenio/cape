package main

import (
	"context"
	"github.com/capeprivacy/cape/cmd/config"
	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/primitives"
	"github.com/urfave/cli/v2"
)

// Harness wraps up the cape application and replaces the default app provider with a test provider that can
// stub various components out
type Harness struct {
	app       *cli.App
	responses []interface{}
	ui        *ui.Mock
}

// NewHarness returns a Harness
// You can provide a list of Responses you want the CLI to respond with
func NewHarness(responses []interface{}) (*cli.App, *ui.Mock) {
	u := &ui.Mock{
		Calls: []*ui.Call{},
	}

	testApp := &Harness{
		app:       NewApp(),
		responses: responses,
		ui:        u,
	}

	testApp.app.Before = testApp.mockBeforeMiddleware
	return testApp.app, testApp.ui
}

func (t *Harness) mockBeforeMiddleware(c *cli.Context) error {
	l, err := primitives.NewLabel("my-cool-cluster")
	if err != nil {
		return nil
	}

	u, err := primitives.NewURL("http://cape.com")
	if err != nil {
		return err
	}

	conf := config.Default()
	_, err = conf.AddCluster(l, u, "")
	if err != nil {
		return err
	}

	err = conf.Use(l)
	if err != nil {
		return err
	}

	m := t.NewMockProvider(c)

	c.Context = context.WithValue(c.Context, ConfigContextKey, conf)
	c.Context = context.WithValue(c.Context, ProviderContextKey, m)

	return nil
}

// MockProvider is what we replace the default provider with
type MockProvider struct {
	context   *cli.Context
	responses []interface{}
	ui        ui.UI
}

// NewMockProvider returns a mock provider
func (t *Harness) NewMockProvider(context *cli.Context) Provider {
	return &MockProvider{context: context, responses: t.responses, ui: t.ui}
}

// UI implements UI from the Provider interface
func (mp *MockProvider) UI(ctx context.Context) ui.UI {
	return mp.ui
}

// Client implements Client from the Provider interface
func (mp *MockProvider) Client(ctx context.Context) (*coordinator.Client, error) {
	t := &coordinator.MockTransport{
		Responses: mp.responses,
		Counter:   0,
	}

	return coordinator.NewClient(t), nil
}
