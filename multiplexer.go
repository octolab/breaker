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
	br := newBreaker()
	go func() {
		defer br.Close()
		select {
		case <-one.Done():
		case <-two.Done():
		}
	}()
	return br
}

// MultiplexThree combines three breakers into one.
// It's an optimized version of a more generic Multiplex.
//
//  interrupter := breaker.MultiplexTwo(
//  	breaker.BreakByContext(req.Context()),
//  	breaker.BreakBySignal(os.Interrupt),
//  	breaker.BreakByTimeout(time.Minute),
//  )
//  defer interrupter.Close()
//
//  background.Job().Do(interrupter)
//
func MultiplexThree(one, two, three Interface) Interface {
	br := newBreaker()
	go func() {
		defer br.Close()
		select {
		case <-one.Done():
		case <-two.Done():
		case <-three.Done():
		}
	}()
	return br
}

func newMultiplexedBreaker(breakers []Interface) Interface {
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
		brs := make([]reflect.SelectCase, 0, len(br.breakers))
		for _, br := range br.breakers {
			brs = append(brs, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(br.Done()),
			})
		}
		reflect.Select(brs)
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
