package mage

import (
	"context"
	"time"
)

type Checker func(context.Context) (bool, error)

func WaitFor(ctx context.Context, f Checker, timeout time.Duration) error {
	newCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := make(chan error, 1)
	go func() {
		// Run checker until the result is successful or until the context has
		// been cancelled, a deadline has been reached, or another error has
		// occurred.
		for {
			success, err := f(newCtx)
			if err != nil {
				result <- err
				break
			}
			if success {
				result <- nil
				break
			}
			if newCtx.Err() != nil {
				break
			}
		}
	}()

	select {
	case <-newCtx.Done():
		return newCtx.Err()
	case err := <-result:
		return err
	}
}
