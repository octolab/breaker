package breaker_test

import (
	"context"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

func TestBreakByContext(t *testing.T) {
	t.Run("active breaker", func(t *testing.T) {
		br := BreakByContext(context.WithTimeout(context.Background(), 5*delta))
		if isReleased(br) {
			t.Error("a breaker is released")
		}

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(5*delta), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("closed breaker", func(t *testing.T) {
		br := BreakByContext(context.WithTimeout(context.Background(), -delta))
		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("released breaker", func(t *testing.T) {
		br := BreakByContext(context.WithTimeout(context.Background(), time.Hour))
		br.Close()
		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("canceled context", func(t *testing.T) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), time.Hour)
			br          = BreakByContext(ctx, cancel)
		)
		cancel()
		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})
}
