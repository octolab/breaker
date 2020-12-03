package breaker

import (
	"reflect"
	"sync/atomic"
)

// Multiplex combines multiple breakers into one.
//
//  interrupter := breaker.Multiplex(
//  	breaker.BreakByContext(req.Context()),
//  	breaker.BreakBySignal(os.Interrupt),
//  	breaker.BreakByTimeout(time.Minute),
//  )
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
func Multiplex(breakers ...Interface) Interface {
	if len(breakers) == 0 {
		return closedBreaker()
	}
	return newMultiplexedBreaker(breakers).trigger()
}

// MultiplexTwo combines two breakers into one.
// It's an optimized version of a more generic Multiplex.
//
//  interrupter := breaker.MultiplexTwo(
//  	breaker.BreakByContext(req.Context()),
//  	breaker.BreakBySignal(os.Interrupt),
//  )
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
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
func MultiplexThree(one, two, three Interface) Interface {
	return newMultiplexedBreaker([]Interface{one, two, three}).trigger()
}

func newMultiplexedBreaker(breakers []Interface) *multiplexedBreaker {
	for len(breakers) < 3 {
		breakers = append(breakers, stub{})
	}
	return &multiplexedBreaker{newBreaker(), breakers}
}

type multiplexedBreaker struct {
	*breaker
	breakers []Interface
}

// Close closes the Done channel and releases resources associated with it.
func (br *multiplexedBreaker) Close() {
	br.closer.Do(func() {
		each(br.breakers).Close()
		close(br.signal)
	})
}

// trigger starts listening to the all Done channels of multiplexed breakers.
func (br *multiplexedBreaker) trigger() Interface {
	go func() {
		if len(br.breakers) == 3 {
			select {
			case <-br.breakers[0].Done():
			case <-br.breakers[1].Done():
			case <-br.breakers[2].Done():
			}
		} else {
			brs := make([]reflect.SelectCase, 0, len(br.breakers))
			for _, br := range br.breakers {
				brs = append(brs, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(br.Done()),
				})
			}
			reflect.Select(brs)
		}
		br.Close()

		// the goroutine is done
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

type each []Interface

// Close closes all Done channels of a list of breakers
// and releases resources associated with them.
func (list each) Close() {
	for _, br := range list {
		br.Close()
	}
}

type stub struct{}

func (br stub) Close()                {}
func (br stub) Done() <-chan struct{} { return nil }
func (br stub) Err() error            { return nil }
func (br stub) trigger() Interface    { return br }
