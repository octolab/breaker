package breaker

import "context"

// MultiplexTwo combines two breakers into one.
//
//  interrupter := breaker.MultiplexTwo(
//  	breaker.BreakByContext(req.Context()),
//  	breaker.BreakBySignal(os.Interrupt),
//  )
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
// Deprecated: Multiplex has the same optimization under the hood now.
// It will be removed at v2.
func MultiplexTwo(one, two Interface) Interface {
	return newMultiplexedBreaker([]Interface{one, two, stub{}}).trigger()
}

// MultiplexThree combines three breakers into one.
// It's an optimized version of a more generic Multiplex.
//
//  interrupter := breaker.MultiplexThree(
//  	breaker.BreakByContext(req.Context()),
//  	breaker.BreakBySignal(os.Interrupt),
//  	breaker.BreakByTimeout(time.Minute),
//  )
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
// Deprecated: Multiplex has the same optimization under the hood now.
// It will be removed at v2.
func MultiplexThree(one, two, three Interface) Interface {
	return newMultiplexedBreaker([]Interface{one, two, three}).trigger()
}

// WithContext returns a new breaker and an associated Context based on the passed one.
//
//  interrupter, ctx := breaker.WithContext(req.Context())
//  defer interrupter.Close()
//
//  background.Job().Run(ctx)
//
// Deprecated: use BreakByContext instead.
// It will be removed at v2.
func WithContext(ctx context.Context) (Interface, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return (&contextBreaker{ctx, cancel}).trigger(), ctx
}
