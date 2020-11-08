package breaker_test

import (
	"context"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		br := New()
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		br := New()
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleasedFast(t, br)
	})
}

func TestBreakByChannel(t *testing.T) {
	t.Parallel()

	t.Run("close channel", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		br := BreakByChannel(ch)
		checkBreakerIsNotReleased(t, br)

		close(ch)
		checkBreakerIsReleased(t, br)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		br := BreakByChannel(ch)
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleased(t, br)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		br := BreakByChannel(ch)
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleased(t, br)
	})
}

func TestBreakByContext(t *testing.T) {
	t.Parallel()

	t.Run("cancel context", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		br := BreakByContext(ctx, cancel)
		checkBreakerIsNotReleased(t, br)

		cancel()
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("propagate timeout", func(t *testing.T) {
		t.Parallel()

		timeout := 5 * delta
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		br := BreakByContext(ctx, cancel)

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
		br := BreakByContext(ctx, cancel)
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		br := BreakByContext(ctx, cancel)
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		br := BreakByContext(ctx, cancel)
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleasedFast(t, br)
		checkContextIsDone(t, ctx)
	})
}

func TestBreakByDeadline(t *testing.T) {
	t.Parallel()

	t.Run("deadline has no passed", func(t *testing.T) {
		t.Parallel()

		timeout := 5 * delta
		br := BreakByDeadline(time.Now().Add(timeout))

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(timeout), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("deadline has already passed", func(t *testing.T) {
		t.Parallel()

		timeout := -delta
		br := BreakByDeadline(time.Now().Add(timeout))
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		br := BreakByDeadline(time.Now().Add(timeout))
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleased(t, br)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		br := BreakByDeadline(time.Now().Add(timeout))
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleased(t, br)
	})
}

func TestBreakBySignal(t *testing.T) {
	t.Parallel()

	t.Run("with signal", func(t *testing.T) {
		t.Parallel()

		br := BreakBySignal(syscall.SIGCHLD)

		start := time.Now()
		go func() {
			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				t.Fatal(err)
			}
			err = proc.Signal(syscall.SIGCHLD)
			if err != nil {
				t.Fatal(err)
			}
		}()
		select {
		case <-br.Done():
		case <-time.After(delta):
			t.Skip("not stable test case")
			br.Close()
		}

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("without signal", func(t *testing.T) {
		t.Parallel()

		br := BreakBySignal()
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		br := BreakBySignal(os.Kill)
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleased(t, br)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		br := BreakBySignal(os.Kill)
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleased(t, br)
	})
}

func TestBreakByTimeout(t *testing.T) {
	t.Parallel()

	t.Run("timeout has no passed", func(t *testing.T) {
		t.Parallel()

		timeout := 5 * delta
		br := BreakByTimeout(timeout)

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(timeout), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("timeout has already passed", func(t *testing.T) {
		t.Parallel()

		timeout := -delta
		br := BreakByTimeout(timeout)
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		br := BreakByTimeout(timeout)
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleased(t, br)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		timeout := time.Hour
		br := BreakByTimeout(timeout)
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleased(t, br)
	})
}

func TestToContext(t *testing.T) {
	br := BreakByTimeout(time.Hour)

	ctx := ToContext(br)
	if ctx.Done() == nil {
		t.Error("bad context")
	}
	if ctx.Err() != nil {
		t.Error("bad context")
	}

	br.Close()
	if time.Sleep(delta); ctx.Err() == nil {
		t.Error("invalid behavior")
	}
}

// helpers

const (
	delta = 10 * time.Millisecond
	times = 10
)

func checkBreakerIsReleased(tb testing.TB, br Interface) {
	tb.Helper()

	time.Sleep(delta)
	checkBreakerIsReleasedFast(tb, br)
}

func checkBreakerIsReleasedFast(tb testing.TB, br Interface) {
	tb.Helper()

	if !isReleased(br) {
		tb.Error("a breaker is not released")
	}
	if isOpened(br) {
		tb.Error("a breaker has inconsistent state")
	}
	if br.Err() == nil {
		tb.Error("a breaker has no error")
	}
}

func checkBreakerIsNotReleased(tb testing.TB, br Interface) {
	tb.Helper()

	if isReleased(br) {
		tb.Error("a breaker is released")
	}
	if !isOpened(br) {
		tb.Error("a breaker has inconsistent state")
	}
	if br.Err() != nil {
		tb.Error("a breaker has error")
	}
}

func checkContextIsDone(tb testing.TB, ctx context.Context) {
	tb.Helper()

	if err := ctx.Err(); err == nil || (err != context.Canceled && err != context.DeadlineExceeded) {
		tb.Error("a context is not done")
	}
}

func checkDuration(tb testing.TB, expected, actual time.Time) {
	tb.Helper()

	if dt := expected.Sub(actual); dt < -delta || dt > delta {
		tb.Errorf(
			"max difference between %v and %v allowed is %v, but difference was %v",
			expected, actual, delta, dt,
		)
	}
}

func closeBreakerConcurrently(br Interface, times int) {
	wg := new(sync.WaitGroup)
	wg.Add(times)

	for range make([]struct{}, times) {
		go func() {
			br.Close()
			wg.Done()
		}()
	}

	wg.Wait()
}

// trigger guarantees that all implementations in under control and
// - the Done channel is never nil
// - the IsReleased call is safe

func isOpened(br Interface) bool {
	opened := br.Done() != nil
	select {
	case _, opened = <-br.Done():
	default:
	}
	return opened
}

func isReleased(br Interface) bool {
	return br.(interface{ IsReleased() bool }).IsReleased()
}
