// logging is a wrapper around rs/zerolog making it easier to create & manage
package logging

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	// UnknownTypeCause represents an error when an unknown logger type is encountered
	UnknownTypeCause = errors.NewCause(errors.BadRequestCategory, "unknown_logger")

	// UnknownLevelCause represents an error when an unknown log level is provided
	UnknownLevelCause = errors.NewCause(errors.BadRequestCategory, "unknown_log_level")
)

// Type represents a type of logger
type Type string

// String returns the type in string form
func (t Type) String() string {
	return string(t)
}

var (
	// PrettyType represents a logger that prints nicely to stdout
	PrettyType Type = "pretty"

	// DiscardType represents a logger that does not log anything
	DiscardType Type = "discard"

	// JSONType represents a logger that emits logs as json to stderr
	JSONType Type = "json"
)

var typeRegistry map[Type]string
var levels []string

// DefaultLevel is the default log level
const DefaultLevel = "info"

func init() {
	typeRegistry = map[Type]string{
		PrettyType:  PrettyType.String(),
		DiscardType: DiscardType.String(),
		JSONType:    JSONType.String(),
	}

	// These levels are the string values for a zerolog.Level
	levels = []string{
		"trace",
		"debug",
		"info",
		"warn",
		"error",
		"fatal",
	}
}

// Types returns a map of a type to string representation
func Types() map[Type]string {
	return typeRegistry
}

// Levels returns a map of a level to string representation
func Levels() []string {
	return levels
}

// ParseType returns a Type for the given string or an error if the type is
// unrecognized.
func ParseType(in string) (Type, error) {
	for t, str := range typeRegistry {
		if str == in {
			return t, nil
		}
	}

	return DiscardType, errors.New(UnknownTypeCause, "no logger named %s exists", in)
}

// Logger returns a logger that can be passed into cape components. The logger
// will either discard all logs or writes logs to stderr depending on
// environment variables.
func Logger(loggerType, logLevel, instanceID string) (*zerolog.Logger, error) {
	t, err := ParseType(loggerType)
	if err != nil {
		return nil, err
	}

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		// If err exists then the log level text was not recognized!
		return nil, errors.New(UnknownLevelCause, "level %s does not exist", logLevel)
	}

	var out io.Writer
	switch t {
	case PrettyType:
		out = zerolog.NewConsoleWriter()
	case DiscardType:
		out = ioutil.Discard
	case JSONType:
		out = os.Stderr
	default:
		return nil, errors.New(UnknownTypeCause, "logger type %s is not supported", string(t))
	}

	logger := zerolog.New(out)
	logger = logger.With().Str("instance_id", instanceID).Timestamp().Logger()
	logger = logger.Level(level)

	return &logger, nil
}
