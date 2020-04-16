package framework

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	gm "github.com/onsi/gomega"
)

func TestSignalWatcher(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("calls handler on signal", func(t *testing.T) {
		shutdownWasCalled := false
		exitWasCalled := false

		w, err := NewSignalWatcher(func(_ context.Context, _ os.Signal) error {
			shutdownWasCalled = true
			return nil
		}, func(_ context.Context, err error) {
			exitWasCalled = true
		}, nil)
		gm.Expect(err).To(gm.BeNil())

		err = w.Start()
		gm.Expect(err).To(gm.BeNil())
		w.notify <- os.Interrupt

		gm.Eventually(func() bool {
			return shutdownWasCalled
		}).Should(gm.BeTrue())
		gm.Eventually(func() bool {
			return exitWasCalled
		}).Should(gm.BeTrue())
	})

	t.Run("handles error from shutdown", func(t *testing.T) {
		shutdownWasCalled := false
		exitWasCalled := false

		w, err := NewSignalWatcher(func(_ context.Context, _ os.Signal) error {
			shutdownWasCalled = true
			return errors.New("hi")
		}, func(_ context.Context, err error) {
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(err.Error()).To(gm.Equal("hi"))

			exitWasCalled = true
		}, nil)
		gm.Expect(err).To(gm.BeNil())

		err = w.Start()
		gm.Expect(err).To(gm.BeNil())
		w.notify <- os.Interrupt

		gm.Eventually(func() bool {
			return shutdownWasCalled
		}).Should(gm.BeTrue())
		gm.Eventually(func() bool {
			return exitWasCalled
		}).Should(gm.BeTrue())
	})

	t.Run("exits if timeout exceeded", func(t *testing.T) {
		shutdownWasCalled := false
		exitWasCalled := false

		done := make(chan bool, 1)
		defer func() {
			done <- true
		}()

		timeout := 1 * time.Millisecond
		w, err := NewSignalWatcher(func(_ context.Context, _ os.Signal) error {
			shutdownWasCalled = true
			<-done
			return nil
		}, func(_ context.Context, err error) {
			gm.Expect(err).To(gm.Equal(ErrTimeout))
			exitWasCalled = true
		}, &timeout)
		gm.Expect(err).To(gm.BeNil())

		err = w.Start()
		gm.Expect(err).To(gm.BeNil())
		w.notify <- os.Interrupt

		gm.Eventually(func() bool {
			return shutdownWasCalled
		}).Should(gm.BeTrue())
		gm.Eventually(func() bool {
			return exitWasCalled
		}).Should(gm.BeTrue())
	})
}
