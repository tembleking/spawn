// Package spawn: A package for spawning goroutines and waiting for their results without messing with channels.
package spawn

import (
	"context"
	"sync/atomic"
)

// JoinHandle is a handle to a spawned goroutine that can be used to wait for the result of the function.
// It is similar to the JoinHandle in Rust's std::thread module.
type JoinHandle[T any] struct {
	success  chan T
	error    chan error
	finished *atomic.Bool
}

// Func spawns a goroutine that runs the given function and returns a JoinHandle
// that can be used to wait for the result of the function.
// Use Wait() or WaitCtx() to wait for the result.
func Func[T any](f func() (T, error)) JoinHandle[T] {
	success := make(chan T)
	error := make(chan error)
	finished := &atomic.Bool{}

	go func() {
		defer close(success)
		defer close(error)

		result, err := f()
		finished.Store(true)
		if err != nil {
			error <- err
		} else {
			success <- result
		}
	}()

	return JoinHandle[T]{
		success:  success,
		error:    error,
		finished: finished,
	}
}

// Wait waits indefinitely for the result of the function and returns the result or an error.
func (j *JoinHandle[T]) Wait() (result T, err error) {
	return j.WaitCtx(context.Background())
}

// WaitCtx waits for the result of the function and returns the result or an error. It returns an error if the context is cancelled.
// If the context is cancelled, the goroutine running the function will continue to run until it finishes, so you need to ensure that the function
// is also cancellable if you want to stop it early.
func (j *JoinHandle[T]) WaitCtx(ctx context.Context) (result T, err error) {
	select {
	case result := <-j.success:
		return result, nil
	case err := <-j.error:
		return result, err
	case <-ctx.Done():
		return result, ctx.Err()
	}
}

// IsFinished returns true if the function has finished running.
func (j *JoinHandle[T]) IsFinished() bool {
	return j.finished.Load()
}
