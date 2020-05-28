package sources

import (
	"context"
	"fmt"
	"github.com/capeprivacy/cape/coordinator/client"
	"sync"

	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// CoordinatorClient is an interface that represents a client for the
// Coordinator.
//
// This interface exists to make it easy to test the sources package in
// isolation of other packages.
type CoordinatorClient interface {
	GetSourceByLabel(ctx context.Context, label primitives.Label, opts *client.SourceOptions) (*client.SourceResponse, error)
}

// Cache is responsible for managing a cache of active sources. Users request a
// data source from the Cache. If one does not exist it will reach out to the
// Coordinator to attempt to create a source (if the connector has access).
//
// Once a source is created the Cache will keep it hot and ready to serve
// requests. In future, if a source is not actively being used it will be aged
// off and closed using behaviour similar to an LRUCache.
type Cache struct {
	closed   bool
	lock     *sync.RWMutex
	client   CoordinatorClient
	sources  map[primitives.Label]Source
	registry *Registry
	cfg      *Config
}

// NewCache returns a Manager if valid configuration is provided.
func NewCache(cfg *Config, c CoordinatorClient, r *Registry) (*Cache, error) {
	if r == nil {
		r = registry
	}

	return &Cache{
		lock:     &sync.RWMutex{},
		cfg:      cfg,
		client:   c,
		closed:   false,
		sources:  map[primitives.Label]Source{},
		registry: r,
	}, nil
}

// Get returns a Source for the given label or returns an error if this Cache
// is not able to fetch the credentials for the data source.
func (c *Cache) Get(ctx context.Context, label primitives.Label) (Source, error) {
	s, err := c.get(label)
	if err != nil && err != ErrCacheNotFound {
		return nil, err
	}
	if s != nil {
		return s, nil
	}

	return c.add(ctx, label)
}

func (c *Cache) get(label primitives.Label) (Source, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.closed {
		return nil, ErrCacheClosed
	}

	s, ok := c.sources[label]
	if !ok {
		return nil, ErrCacheNotFound
	}

	return s, nil
}

func (c *Cache) add(ctx context.Context, label primitives.Label) (Source, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	source, err := c.client.GetSourceByLabel(ctx, label, nil)
	if err != nil {
		return nil, err
	}

	ctor, err := c.registry.Get(source.Type)
	if err != nil {
		return nil, err
	}

	s, err := ctor(ctx, c.cfg, source.Source)
	if err != nil {
		return nil, err
	}

	c.sources[label] = s
	return s, nil
}

// Close closes all sources and returns an error if any source error'd while
// attemping to close the sources
//
// XXX: In future this should have a built-in timeout in case a Source does not
// close properly or if it gets stalled.
func (c *Cache) Close(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// It would be ideal if this was done in parallel instead of in series.
	errs := []string{}
	for _, source := range c.sources {
		err := source.Close(ctx)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Could not close %s: %s", source.Label(), err.Error()))
		}
	}

	c.closed = true
	c.sources = map[primitives.Label]Source{}

	if len(errs) > 0 {
		return errors.NewMulti(ClosingCause, errs)
	}

	return nil
}
