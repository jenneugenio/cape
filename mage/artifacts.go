package mage

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// Tracker is a global tracker of build artifacts
var Tracker *Artifacts

type Cleaner func(context.Context) error

func CleanFile(path string) Cleaner {
	return func(_ context.Context) error {
		return os.Remove(path)
	}
}

func init() {
	Tracker = &Artifacts{
		artifacts: map[string]Cleaner{},
		lock:      &sync.Mutex{},
	}
}

type Artifacts struct {
	artifacts map[string]Cleaner
	lock      *sync.Mutex
}

func (a Artifacts) Add(path string, f Cleaner) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.artifacts[path]; ok {
		return fmt.Errorf("Path already tracked as an artifact: %s", path)
	}

	a.artifacts[path] = f
	return nil
}

func (a Artifacts) Clean(ctx context.Context) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	// TODO: Introduce a multi error type
	errors := []error{}
	for _, f := range a.artifacts {
		err := f(ctx)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		return errors[0]
	}

	return nil
}
