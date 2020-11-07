// Package breaker provides flexible mechanism to make your code breakable.
package breaker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

// New returns a new Breaker which can interrupted only by the Close call.
func New() Interface {
	return newBreaker().trigger()
}

// BreakByChannel returns a new Breaker based on the channel.
func BreakByChannel(ch <-chan struct{}) Interface {
	return (&channelBreaker{newBreaker(), ch}).trigger()
}

// BreakByContext returns a new Breaker based on the Context.
//
//  interrupter := breaker.BreakByContext(context.WithTimeout(req.Context(), time.Minute)),
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
func BreakByContext(ctx context.Context, cancel context.CancelFunc) Interface {
	return &contextBreaker{ctx, cancel}
}

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
	return br
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

// Err returns a non-nil error if the Done channel is closed and nil otherwise.
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

type channelBreaker struct {
	*breaker
	relay <-chan struct{}
}

// Close closes the Done channel and releases resources associated with it.
func (br *channelBreaker) Close() {
	br.closer.Do(func() {
		close(br.signal)
	})
}

// trigger starts listening internal signal to close the Done channel.
func (br *channelBreaker) trigger() Interface {
	go func() {
		select {
		case <-br.relay:
		case <-br.signal:
		}
		br.Close()

		// the goroutine is done
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

type contextBreaker struct {
	context.Context
	cancel context.CancelFunc
}

// Close closes the Done channel and releases resources associated with it.
func (br *contextBreaker) Close() {
	br.cancel()
}

// Released returns true if resources associated with the Breaker were released.
func (br *contextBreaker) Released() bool {
	select {
	case <-br.Done():
		return true
	default:
		return false
	}
}

func (br *contextBreaker) trigger() Interface {
	return br
}

func newSignaledBreaker(signals []os.Signal) *signaledBreaker {
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

		// the goroutine is done
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

func newTimedBreaker(timeout time.Duration) *timedBreaker {
	return &timedBreaker{newBreaker(), time.NewTimer(timeout)}
}

type timedBreaker struct {
	*breaker
	*time.Timer
}

// Close closes the Done channel and releases resources associated with it.
func (br *timedBreaker) Close() {
	br.closer.Do(func() {
		stop(br.Timer)
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

		// the goroutine is done
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

func stop(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}
