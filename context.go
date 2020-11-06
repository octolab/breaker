package breaker

import (
	"context"
)

// BreakByContext returns a new Breaker based on the Context.
//
//  interrupter := breaker.BreakByContext(context.WithTimeout(req.Context(), time.Minute)),
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
func BreakByContext(ctx context.Context, cancel context.CancelFunc) Interface {
	return (&contextBreaker{ctx, cancel}).trigger()
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
