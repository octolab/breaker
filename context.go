package breaker

import (
	"context"
	"sync/atomic"
)

// WithContext returns a new Breaker and an associated Context derived from ctx.
func WithContext(ctx context.Context) (Interface, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return (&contextBreaker{newBreaker(), cancel, ctx.Done()}).trigger(), ctx
}

type contextBreaker struct {
	*breaker
	cancel context.CancelFunc
	signal <-chan struct{}
}

// Done returns a channel that's closed when a cancellation signal occurred.
func (br *contextBreaker) Done() <-chan struct{} {
	return br.signal
}

// Close closes the Done channel and releases resources associated with it.
func (br *contextBreaker) Close() {
	br.cancel()
}

func (br *contextBreaker) trigger() Interface {
	go func() {
		<-br.signal
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}
