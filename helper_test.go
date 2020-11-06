package breaker_test

import (
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
)

const delta = 10 * time.Millisecond

func checkBreakerIsReleased(tb testing.TB, br Interface) {
	tb.Helper()

	if time.Sleep(delta); !isReleased(br) {
		tb.Error("a breaker is not released")
	}
	if br.Err() == nil {
		tb.Error("a breaker has no error")
	}
}

func checkBreakerIsReleasedFast(tb testing.TB, br Interface) {
	tb.Helper()

	if !isReleased(br) {
		tb.Error("a breaker is not released")
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
	if br.Err() != nil {
		tb.Error("a breaker has error")
	}
}

func checkDuration(tb testing.TB, expected, actual time.Time) {
	tb.Helper()

	if dt := expected.Sub(actual); dt < -delta || dt > delta {
		tb.Errorf("max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt)
	}
}

func isReleased(br Interface) bool {
	check, is := br.(interface{ Released() bool })
	return is && check.Released()
}
