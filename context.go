package breaker

import (
	"context"
	"sync/atomic"
)

// BreakByContext returns a new Breaker based on the Context.
func BreakByContext(ctx context.Context, cancel context.CancelFunc) Interface {
	return (&contextBreaker{newBreaker(), cancel, ctx}).trigger()
}

type contextBreaker struct {
	*breaker
	cancel context.CancelFunc
	ctx    context.Context
}

// Close closes the Done channel and releases resources associated with it.
func (br *contextBreaker) Close() {
	br.cancel()
}

// Done returns a channel that's closed when a cancellation signal occurred.
func (br *contextBreaker) Done() <-chan struct{} {
	return br.ctx.Done()
}

// Err returns a non-nil error if Done is closed and nil otherwise.
// After Err returns a non-nil error, successive calls to Err return the same error.
func (br *contextBreaker) Err() error {
	return br.ctx.Err()
}

func (br *contextBreaker) trigger() Interface {
	go func() {
		<-br.ctx.Done()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}
