package breaker_test

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

func TestNew(t *testing.T) {
	br := New()
	checkBreakerIsNotReleased(t, br)

	br.Close()
	checkBreakerIsReleasedFast(t, br)
}

func TestBreakByChannel(t *testing.T) {
	t.Run("release breaker", func(t *testing.T) {
		ch := make(chan struct{})
		br := BreakByChannel(ch)
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleased(t, br)
	})

	t.Run("close channel", func(t *testing.T) {
		ch := make(chan struct{})
		br := BreakByChannel(ch)
		checkBreakerIsNotReleased(t, br)

		close(ch)
		checkBreakerIsReleased(t, br)
	})
}

func TestBreakByContext(t *testing.T) {
	t.Run("propagate timeout", func(t *testing.T) {
		timeout := 5 * delta
		br := BreakByContext(context.WithTimeout(context.Background(), timeout))

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(timeout), time.Now())
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("deadline has already passed", func(t *testing.T) {
		br := BreakByContext(context.WithTimeout(context.Background(), -delta))
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("release breaker", func(t *testing.T) {
		br := BreakByContext(context.WithTimeout(context.Background(), time.Hour))
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("cancel context", func(t *testing.T) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), time.Hour)
			br          = BreakByContext(ctx, cancel)
		)
		checkBreakerIsNotReleased(t, br)

		cancel()
		checkBreakerIsReleasedFast(t, br)
	})
}

func TestBreakByDeadline(t *testing.T) {
	t.Run("future deadline", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(5 * delta))

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(5*delta), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("passed deadline", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(-delta))

		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("close multiple times", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(time.Hour))
		br.Close()
		br.Close()

		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})
}

func TestBreakBySignal(t *testing.T) {
	t.Run("with signal", func(t *testing.T) {
		br := BreakBySignal(syscall.SIGCHLD)

		start := time.Now()
		go func() {
			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				panic(fmt.Errorf("error is not expected: %v", err))
			}
			err = proc.Signal(syscall.SIGCHLD)
			if err != nil {
				panic(fmt.Errorf("error is not expected: %v", err))
			}
		}()
		select {
		case <-br.Done():
		case <-time.After(9 * time.Millisecond):
			t.Skip("not stable test case")
			br.Close()
		}

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("without signal", func(t *testing.T) {
		br := BreakBySignal()

		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("close multiple times", func(t *testing.T) {
		br := BreakBySignal(os.Kill)
		br.Close()
		br.Close()

		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})
}

func TestBreakByTimeout(t *testing.T) {
	t.Run("valid timeout", func(t *testing.T) {
		br := BreakByTimeout(5 * delta)

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(5*delta), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("passed timeout", func(t *testing.T) {
		br := BreakByTimeout(-delta)

		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("close multiple times", func(t *testing.T) {
		br := BreakByTimeout(time.Hour)
		br.Close()
		br.Close()

		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})
}
