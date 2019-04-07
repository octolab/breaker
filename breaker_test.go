package breaker_test

import (
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

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
		br := BreakBySignal(os.Interrupt)
		start := time.Now()
		go func() {
			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				t.Fatalf("error is not expected: %v", err)
			}
			err = proc.Signal(os.Interrupt)
			if err != nil {
				t.Fatalf("error is not expected: %v", err)
			}
		}()
		<-br.Done()

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
