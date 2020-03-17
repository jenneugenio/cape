package framework

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/logging"
)

// TestLogger returns a logger for use in writing tests
func TestLogger() *zerolog.Logger {
	loggerType := os.Getenv("CAPE_LOGGING_TYPE")
	if loggerType == "" {
		loggerType = "discard"
	}

	logLevel := os.Getenv("CAPE_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logger, err := logging.Logger(loggerType, logLevel, "test")
	if err != nil {
		panic(fmt.Sprintf("could not create test logger: %s", err))
	}

	return logger
}
