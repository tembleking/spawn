package spawn_test

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/spawn"
)

var _ = Describe("Spawn", MustPassRepeatedly(100), func() {
	It("spawns a function in a different goroutine", func(ctx context.Context) {
		spawned := &atomic.Bool{}
		spawn.Func(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			spawned.Store(true)
			return 0, nil
		})

		Expect(spawned.Load()).To(BeFalse())
		Eventually(spawned.Load).Should(BeTrue())
	})

	It("joins the result", func(ctx context.Context) {
		joinHandle := spawn.Func(func() (int, error) {
			return 42, nil
		})

		result, err := joinHandle.Wait()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(42))
	})

	It("joins the error", func(ctx context.Context) {
		joinHandle := spawn.Func(func() (int, error) {
			return 0, errors.New("some error")
		})

		result, err := joinHandle.Wait()
		Expect(err).To(MatchError("some error"))
		Expect(result).To(Equal(0))
	})

	It("waits the result until the context is done", func(ctx context.Context) {
		joinHandle := spawn.Func(func() (int, error) {
			time.Sleep(1 * time.Hour)
			return 42, nil
		})

		ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		result, err := joinHandle.WaitCtx(ctx)
		Expect(err).To(MatchError(context.DeadlineExceeded))
		Expect(result).To(Equal(0))
	})

	It("returns finished after it's done", func(ctx context.Context) {
		joinHandle := spawn.Func(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			return 42, nil
		})

		Expect(joinHandle.IsFinished()).To(BeFalse())
		Eventually(joinHandle.IsFinished).Should(BeTrue())
	})

	It("executes a function that returns nothing", func(ctx context.Context) {
		funcThatReturnsNothing := func() {
			time.Sleep(10 * time.Millisecond)
		}

		joinHandle := spawn.Func(func() (struct{}, error) {
			funcThatReturnsNothing()
			return struct{}{}, nil
		})

		_, _ = joinHandle.WaitCtx(ctx)
		Expect(joinHandle.IsFinished()).To(BeTrue())
	})

	It("executes a function that returns an error", func(ctx context.Context) {
		funcThatReturnsError := func() error {
			time.Sleep(10 * time.Millisecond)
			return errors.New("some error")
		}

		joinHandle := spawn.Func(func() (struct{}, error) {
			return struct{}{}, funcThatReturnsError()
		})

		_, err := joinHandle.WaitCtx(ctx)
		Expect(err).To(MatchError("some error"))
	})
})
