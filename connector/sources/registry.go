package sources

import (
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// NewSourceFunc represents a function constructor for a Source
type NewSourceFunc func(*primitives.Source) (Source, error)

// Registry represents a NewSourceFunc registry
type Registry map[primitives.SourceType]NewSourceFunc

// Register enables a source func to register for a given label
func (r Registry) Register(t primitives.SourceType, f NewSourceFunc) error {
	_, ok := r[t]
	if ok {
		return errors.New(SourceAlreadyExists, "%s has already been registered", t.String())
	}

	r[t] = f
	return nil
}

// Get returns a source for the given label if the source exists
func (r Registry) Get(t primitives.SourceType) (NewSourceFunc, error) {
	f, ok := r[t]
	if !ok {
		return nil, errors.New(SourceNotSupported, "%s not supported by connector", t.String())
	}

	return f, nil
}

var registry = &Registry{}
