package breaker_test

import (
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

const delta = 10 * time.Millisecond

func checkBreakerIsReleased(tb testing.TB, br Interface) {
	tb.Helper()

	time.Sleep(delta)
	if check, is := br.(interface{ Released() bool }); !is || !check.Released() {
		tb.Error("a breaker is not released")
	}
}

func checkDuration(tb testing.TB, expected, actual time.Time) {
	tb.Helper()

	dt := expected.Sub(actual)
	if dt < -delta || dt > delta {
		tb.Errorf("max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt)
	}
}
