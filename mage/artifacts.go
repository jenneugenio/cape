package mage

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// Tracker is a global tracker of build artifacts
var Tracker *Artifacts

func init() {
	Tracker = &Artifacts{
		artifacts: map[string]bool{},
		lock:      &sync.Mutex{},
	}
}

type Artifacts struct {
	artifacts map[string]bool
	lock      *sync.Mutex
}

func (a Artifacts) Add(path string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.artifacts[path]; ok {
		return fmt.Errorf("Path already tracked as an artifact: %s", path)
	}

	a.artifacts[path] = true
	return nil
}

func (a Artifacts) Clean(_ context.Context) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	// TODO: Introduce a multi error type
	errors := []error{}
	for path := range a.artifacts {
		err := os.Remove(path)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		return errors[0]
	}

	return nil
}
