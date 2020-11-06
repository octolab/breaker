package breaker

import "context"

// WithContext returns a new Breaker and an associated Context derived from ctx.
// Deprecated: use BreakByContext instead.
// TODO:v2 will be removed
func WithContext(ctx context.Context) (Interface, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return (&contextBreaker{ctx, cancel}).trigger(), ctx
}
