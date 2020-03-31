package framework

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/primitives"
)

// Component represents the business logic and functionality that is served by
// the server.
//
// A component is abstraction that enables a logical separation between the
// thing being served and the thing that actually manages the transport layer.
type Component interface {
	Setup(context.Context) (http.Handler, error)
	Teardown(context.Context) error
	CertFiles() (certFile string, keyFile string)
}

// Config represents a configuration object for a component.
//
// This interface enables components to have different configuration objects as
// long as they satisfy the needs of the Server
type Config interface {
	GetPort() int
	GetInstanceID() primitives.Label
	Validate() error
}

// Server represents an http server. It's responsibility is to provide the
// transport layer for the application contain within the component.
type Server struct {
	server    *http.Server
	component Component
	cfg       Config
	logger    *zerolog.Logger
}

// NewServer returns a new server that will be able to serve the provided
// Component and configuration.
func NewServer(cfg Config, s Component, logger *zerolog.Logger) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Server{
		component: s,
		cfg:       cfg,
		logger:    logger,
	}, nil
}

// Start sets up the component and then configures the http server. It
// subsequently attempts to get the server to listen on the configured addr.
func (s *Server) Start(ctx context.Context) error {
	handler, err := s.component.Setup(ctx)
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:    s.Addr(),
		Handler: handler,
	}

	s.logger.Info().Msgf("%s attempting to listen %s", s.cfg.GetInstanceID(), s.Addr())

	certFile, keyFile := s.component.CertFiles()
	if certFile != "" && keyFile != "" {
		err = s.server.ListenAndServeTLS(certFile, keyFile)
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}

	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop tears own the component and http server. It first stops all incoming http
// requests, drains any on-going requests, and then tears down the Component
func (s *Server) Stop(ctx context.Context) error {
	server := s.server
	component := s.component

	s.server = nil
	s.component = nil

	// Shutdown the http server first so we can begin draining connections.
	server.SetKeepAlivesEnabled(false)
	err := server.Shutdown(ctx)

	// Always attempt to shutdown the component!
	componentErr := component.Teardown(ctx)
	if err != nil {
		return err
	}
	if componentErr != nil {
		return err
	}

	return nil
}

// Addr returns the configured port and host that the Server listens on
func (s *Server) Addr() string {
	return fmt.Sprintf(":%d", s.cfg.GetPort())
}
