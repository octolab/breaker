package breaker_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

func TestMultiplexTwo(t *testing.T) {
	t.Parallel()

	br := MultiplexTwo(
		BreakByDeadline(time.Now().Add(-delta)),
		BreakByTimeout(time.Hour),
	)
	checkBreakerIsReleased(t, br)
}

func TestMultiplexThree(t *testing.T) {
	t.Parallel()

	br := MultiplexThree(
		BreakByDeadline(time.Now().Add(-delta)),
		BreakBySignal(os.Kill),
		BreakByTimeout(time.Hour),
	)
	checkBreakerIsReleased(t, br)
}

func TestWithContext(t *testing.T) {
	t.Parallel()

	t.Run("cancel context", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		br, ctx := WithContext(ctx)
		checkBreakerIsNotReleased(t, br)

		cancel()
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("propagate timeout", func(t *testing.T) {
		t.Parallel()

		timeout := 5 * delta
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		br, ctx := WithContext(ctx)

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(timeout), time.Now())
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("deadline has already passed", func(t *testing.T) {
		t.Parallel()

		timeout := -delta
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		br, ctx := WithContext(ctx)
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		br, ctx := WithContext(ctx)
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		br, ctx := WithContext(ctx)
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})
}
