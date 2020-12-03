package breaker_test

import (
	"context"
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
			BreakByChannel(context.TODO().Done()),
			BreakByContext(context.WithCancel(context.TODO())),
			BreakByDeadline(time.Now().Add(time.Hour)),
			BreakBySignal(os.Interrupt),
			BreakByTimeout(timeout),
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
