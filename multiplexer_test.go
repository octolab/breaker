package breaker_test

import (
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

func TestMultiplex(t *testing.T) {
	t.Run("with breakers", func(t *testing.T) {
		br := Multiplex(
			BreakByTimeout(5*delta),
			BreakByDeadline(time.Now().Add(time.Hour)),
		)
		defer br.Close()
		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(5*delta), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("without breakers", func(t *testing.T) {
		br := Multiplex()
		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("close multiple times", func(t *testing.T) {
		br := Multiplex(BreakByTimeout(time.Hour))
		br.Close()
		br.Close()
		start := time.Now()
		<-br.Done()

		checkDuration(t, start, time.Now())
		checkBreakerIsReleased(t, br)
	})
}

func TestMultiplexTwo(t *testing.T) {
	br := MultiplexTwo(
		BreakByDeadline(time.Now().Add(-delta)),
		BreakByTimeout(time.Hour),
	)
	start := time.Now()
	<-br.Done()

	checkDuration(t, start, time.Now())
	checkBreakerIsReleased(t, br)
}

func TestMultiplexThree(t *testing.T) {
	br := MultiplexThree(
		BreakByDeadline(time.Now().Add(-delta)),
		BreakBySignal(os.Kill),
		BreakByTimeout(time.Hour),
	)
	start := time.Now()
	<-br.Done()

	checkDuration(t, start, time.Now())
	checkBreakerIsReleased(t, br)
}
