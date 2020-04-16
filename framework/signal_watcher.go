package framework

import (
	"context"
	"os"
	"os/signal"
	"time"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// ErrTimeout occurs when a shutdown request takes longer than the given timeframe
var ErrTimeout = errors.New(errors.TimeoutCause, "shutdown_timeout")

// DefaultSignalTimeout represents the default time the watcher will wait for a
// shutdown func to exit before the ExitFunc will be called.
var DefaultSignalTimeout = 10 * time.Second

// ShutdownFunc represents a function that is called when a signal is caught.
type ShutdownFunc func(context.Context, os.Signal) error

// ExitFunc represents a function that is called when a signal handler times out or the shutdown func returns an error
type ExitFunc func(context.Context, error)

// SignalWatcher provides functionality for watching for os.Interrupt signals
// and then shutting down resources. If the resource shutdown takes longer then
// the provided duration then the watcher calls the exit handler.
type SignalWatcher struct {
	triggered bool
	shutdown  ShutdownFunc
	exit      ExitFunc
	timeout   time.Duration

	notify chan os.Signal
	stop   chan bool
}

// NewSignalWatcher returns a new signal watcher that can start listening for
// os signals.
func NewSignalWatcher(f ShutdownFunc, exitfunc ExitFunc, timeout *time.Duration) (*SignalWatcher, error) {
	if timeout == nil {
		timeout = &DefaultSignalTimeout
	}

	return &SignalWatcher{
		triggered: false,
		shutdown:  f,
		exit:      exitfunc,
		timeout:   *timeout,
		notify:    make(chan os.Signal, 1),
		stop:      make(chan bool, 1),
	}, nil
}

// Start initiates the watching of signals.
func (w *SignalWatcher) Start() error {
	if w.triggered {
		return errors.New(errors.InvalidStateCause, "Cannot start watching if notify has been triggered")
	}

	signal.Notify(w.notify, os.Interrupt)
	go w.watch()

	return nil
}

// Stop prevents the watcher from noticing any incoming signals and acting on
// them.
func (w *SignalWatcher) Stop() {
	if w.triggered {
		return
	}

	signal.Stop(w.notify)
	w.stop <- true

	close(w.notify)
	close(w.stop)
}

func (w *SignalWatcher) watch() {
	select {
	case s := <-w.notify:
		w.handle(s)
	case <-w.stop:
		return
	}
}

func (w *SignalWatcher) handle(s os.Signal) {
	if w.triggered {
		return
	}

	// The watcher has been triggered so lets note that we've gone ahead and
	// asked the thing we're managing to shutdown.
	//
	// At the same time, we stop "watching" so no new notifications come
	// through. This ensures that if a signal is issued twice we fall back to
	// the default "go" behaviour of just exiting on an interrupt
	w.triggered = true
	w.Stop()

	// We create a context with a timeout so if shutting down takes longer than
	// our defined time we can fallback to just exiting. This prevents our
	// processes from "getting stuck"
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	result := make(chan error, 1)
	go func() {
		err := w.shutdown(ctx, s)
		result <- err
	}()

	select {
	case <-ctx.Done():
		w.exit(ctx, ErrTimeout)
	case err := <-result:
		w.exit(ctx, err)
	}
}
