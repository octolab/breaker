package breaker

import "testing"

type extended interface {
	Interface
	// IsReleased returns true if resources associated with the breaker were released.
	//
	// Deprecated: the original implementation contained a bug,
	// and there is no longer any need for this method. It will be removed at v2.
	IsReleased() bool
}

var (
	_ extended = new(breaker)
	_ extended = new(signalBreaker)
	_ extended = new(channelBreaker)
	_ extended = new(contextBreaker)
	_ extended = new(timeoutBreaker)
	_ extended = stub{}
)

func TestStub_internals(t *testing.T) {
	var breaker stub

	if breaker.Done() != nil {
		t.Error("stub's Done channel must be nil")
	}

	if breaker.Err() != Interrupted {
		t.Error("stub must be interrupted")
	}

	if !breaker.IsReleased() {
		t.Error("stub must be always released")
	}

	if breaker != breaker.trigger() {
		t.Error("unexpected behavior of stub's trigger method")
	}
}
