// Package breaker provides flexible mechanism to make your code breakable.
package breaker

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

// BreakByDeadline closes the Done channel when the deadline occurs.
func BreakByDeadline(deadline time.Time) Interface {
	timeout := time.Until(deadline)
	if timeout < 0 {
		return closedBreaker()
	}
	return newTimedBreaker(timeout).trigger()
}

// BreakBySignal closes the Done channel when signals will be received.
func BreakBySignal(sig ...os.Signal) Interface {
	if len(sig) == 0 {
		return closedBreaker()
	}
	return newSignaledBreaker(sig).trigger()
}

// BreakByTimeout closes the Done channel when the timeout happens.
func BreakByTimeout(timeout time.Duration) Interface {
	if timeout < 0 {
		return closedBreaker()
	}
	return newTimedBreaker(timeout).trigger()
}

func closedBreaker() Interface {
	br := newBreaker()
	br.Close()
	return br.trigger()
}

func newBreaker() *breaker {
	return &breaker{signal: make(chan struct{})}
}

type breaker struct {
	closer   sync.Once
	signal   chan struct{}
	released int32
}

// Close closes the Done channel and releases resources associated with it.
func (br *breaker) Close() {
	br.closer.Do(func() {
		close(br.signal)
		atomic.StoreInt32(&br.released, 1)
	})
}

// Done returns a channel that's closed when a cancellation signal occurred.
func (br *breaker) Done() <-chan struct{} {
	return br.signal
}

// Err returns a non-nil error if Done is closed and nil otherwise.
// After Err returns a non-nil error, successive calls to Err return the same error.
func (br *breaker) Err() error {
	if atomic.LoadInt32(&br.released) == 1 {
		return Interrupted
	}
	return nil
}

// Released returns true if resources associated with the Breaker were released.
func (br *breaker) Released() bool {
	return atomic.LoadInt32(&br.released) == 1
}

func (br *breaker) trigger() Interface {
	return br
}

func newSignaledBreaker(signals []os.Signal) Interface {
	return &signaledBreaker{newBreaker(), make(chan os.Signal, len(signals)), signals}
}

type signaledBreaker struct {
	*breaker
	relay   chan os.Signal
	signals []os.Signal
}

// Close closes the Done channel and releases resources associated with it.
func (br *signaledBreaker) Close() {
	br.closer.Do(func() {
		signal.Stop(br.relay)
		close(br.signal)
	})
}

// trigger starts listening required signals to close the Done channel.
func (br *signaledBreaker) trigger() Interface {
	go func() {
		signal.Notify(br.relay, br.signals...)
		select {
		case <-br.relay:
		case <-br.signal:
		}
		br.Close()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

func newTimedBreaker(timeout time.Duration) Interface {
	return &timedBreaker{newBreaker(), time.NewTimer(timeout)}
}

type timedBreaker struct {
	*breaker
	*time.Timer
}

// Close closes the Done channel and releases resources associated with it.
func (br *timedBreaker) Close() {
	br.closer.Do(func() {
		br.Timer.Stop()
		close(br.signal)
	})
}

// trigger starts listening internal timer to close the Done channel.
func (br *timedBreaker) trigger() Interface {
	go func() {
		select {
		case <-br.Timer.C:
		case <-br.signal:
		}
		br.Close()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}
