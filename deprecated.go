package breaker

import "context"

// WithContext returns a new breaker and an associated Context based on the passed one.
//
//  interrupter, ctx := breaker.WithContext(req.Context())
//  defer interrupter.Close()
//
//  background.Job().Run(ctx)
//
// Deprecated: use BreakByContext instead. TODO:v2 will be removed.
func WithContext(ctx context.Context) (Interface, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return (&contextBreaker{ctx, cancel}).trigger(), ctx
}
