package breaker_test

import (
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

func TestMultiplex(t *testing.T) {
	t.Parallel()

	t.Run("with breakers", func(t *testing.T) {
		t.Parallel()

		timeout := 5 * delta
		br := Multiplex(
			BreakByTimeout(timeout),
			BreakByDeadline(time.Now().Add(time.Hour)),
		)

		start := time.Now()
		<-br.Done()

		checkDuration(t, start.Add(timeout), time.Now())
		checkBreakerIsReleased(t, br)
	})

	t.Run("without breakers", func(t *testing.T) {
		t.Parallel()

		br := Multiplex()
		checkBreakerIsReleasedFast(t, br)
	})

	t.Run("close breaker", func(t *testing.T) {
		t.Parallel()

		br := Multiplex(BreakByTimeout(time.Hour))
		checkBreakerIsNotReleased(t, br)

		br.Close()
		checkBreakerIsReleased(t, br)
	})

	t.Run("close breaker multiple times", func(t *testing.T) {
		t.Parallel()

		br := Multiplex(BreakByTimeout(time.Hour))
		checkBreakerIsNotReleased(t, br)

		closeBreakerConcurrently(br, times)
		checkBreakerIsReleased(t, br)
	})
}

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
