package harness

import (
	"context"

	"github.com/capeprivacy/cape/coordinator"
	coordHarness "github.com/capeprivacy/cape/coordinator/harness"
)

// Stack runs an entire cape stack -- The coordinator harness, and the connector harness
type Stack struct {
	Manager      *coordHarness.Manager
	CoordHarness *coordHarness.Harness
	CoordClient  *coordinator.Client
	ConnHarness  *Harness
}

func (s *Stack) Teardown(ctx context.Context) error {
	err := s.CoordHarness.Teardown(ctx)
	if err != nil {
		return err
	}

	err = s.ConnHarness.Teardown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewStack(ctx context.Context) (*Stack, error) {
	cfg, err := coordHarness.NewConfig()
	if err != nil {
		return nil, err
	}

	coordH, err := coordHarness.NewHarness(cfg)
	if err != nil {
		return nil, err
	}

	err = coordH.Setup(ctx)
	if err != nil {
		return nil, err
	}

	m := coordH.Manager()
	c, err := m.Setup(ctx)
	if err != nil {
		return nil, err
	}

	coordinatorURL, err := m.URL()
	if err != nil {
		return nil, err
	}

	err = m.CreateService(ctx, ConnectorEmail, coordinatorURL)
	if err != nil {
		return nil, err
	}

	connCfg, err := NewConfig(coordinatorURL, m.Connector.Token)
	if err != nil {
		return nil, err
	}

	connH, err := NewHarness(connCfg)
	if err != nil {
		return nil, err
	}

	err = connH.Setup(ctx)
	if err != nil {
		return nil, err
	}

	s := &Stack{
		Manager:      m,
		CoordClient:  c,
		CoordHarness: coordH,
		ConnHarness:  connH,
	}

	return s, err
}
