package spawn

import (
	"context"
	"sync/atomic"
)

type JoinHandle[T any] struct {
	success  chan T
	error    chan error
	finished *atomic.Bool
}

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

func (j *JoinHandle[T]) Wait() (result T, err error) {
	return j.WaitCtx(context.Background())
}

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

func (j *JoinHandle[T]) IsFinished() bool {
	return j.finished.Load()
}
