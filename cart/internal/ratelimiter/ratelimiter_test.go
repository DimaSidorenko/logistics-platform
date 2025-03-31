package ratelimiter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestLimiter_Wait(t *testing.T) {
	t.Run("should allow requests within RPS limit", func(t *testing.T) {
		t.Parallel()
		const rpsLimit = 5
		limiter := New(rpsLimit)
		defer limiter.Stop()

		ctx := context.Background()
		for i := 0; i < rpsLimit; i++ {
			err := limiter.Wait(ctx)
			require.NoError(t, err, "request within limit should pass")
		}
	})

	t.Run("should block after exceeding RPS limit", func(t *testing.T) {
		t.Parallel()
		const rpsLimit = 2
		limiter := New(rpsLimit)
		defer limiter.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Exhaust the limit
		for i := 0; i < rpsLimit; i++ {
			err := limiter.Wait(context.Background())
			require.NoError(t, err, "initial requests should pass")
		}

		// Next request should block
		err := limiter.Wait(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded, "should return DeadlineExceeded when blocked")
	})

	t.Run("should respect context cancellation", func(t *testing.T) {
		t.Parallel()
		limiter := New(1)
		defer limiter.Stop()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Immediate cancellation

		_ = limiter.Wait(ctx)
		err := limiter.Wait(ctx)
		assert.ErrorIs(t, err, context.Canceled, "should return Canceled when context is canceled")
	})

	t.Run("should refill tokens after 1 second", func(t *testing.T) {
		t.Parallel()
		const rpsLimit = 1
		limiter := New(rpsLimit)
		defer limiter.Stop()

		// First request should pass
		err := limiter.Wait(context.Background())
		require.NoError(t, err, "first request should pass")

		// Second should block
		fastCtx, fastCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer fastCancel()
		err = limiter.Wait(fastCtx)
		assert.ErrorIs(t, err, context.DeadlineExceeded, "should block when limit exhausted")

		// Wait for refill
		time.Sleep(1100 * time.Millisecond)

		// Now should pass again
		err = limiter.Wait(context.Background())
		require.NoError(t, err, "request after refill should pass")
	})

	t.Run("should handle Stop correctly", func(t *testing.T) {
		t.Parallel()
		limiter := New(1)
		limiter.Stop()

		err := limiter.Wait(context.Background())
		assert.Error(t, err, "should return error after Stop")
	})

	t.Run("should not panic on double Stop", func(t *testing.T) {
		t.Parallel()
		limiter := New(1)
		limiter.Stop()
		assert.NotPanics(t, func() { limiter.Stop() }, "second Stop should not panic")
	})

	t.Run("concurrent usage", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		const (
			rpsLimit    = 10
			goroutines  = 10
			totalCalls  = 30
			minimumTime = 2 * time.Second
		)

		limiter := New(rpsLimit)
		// can comment defer and check goleak
		defer limiter.Stop()

		var (
			wg     sync.WaitGroup
			passed int
			mu     sync.Mutex
		)

		start := time.Now()

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < totalCalls/goroutines; j++ {
					if err := limiter.Wait(context.Background()); err == nil {
						mu.Lock()
						passed++
						mu.Unlock()
					}
				}
			}()
		}

		wg.Wait()
		elapsed := time.Since(start)

		assert.Equal(t, totalCalls, passed, "all requests should complete")
		assert.GreaterOrEqual(t, elapsed, minimumTime, "should finish after minimumTime")
	})
}
