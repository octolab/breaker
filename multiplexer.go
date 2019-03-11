package breaker

import (
	"reflect"
	"sync/atomic"
)

// Multiplex combines multiple Breakers into one.
func Multiplex(breakers ...Interface) Interface {
	if len(breakers) == 0 {
		return closedBreaker()
	}
	return newMultiplexedBreaker(breakers).trigger()
}

// MultiplexTwo combines two Breakers into one.
// This is the optimized version of more generic Multiplex.
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

// MultiplexThree combines three Breakers into one.
// This is the optimized version of more generic Multiplex.
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

func newMultiplexedBreaker(entries []Interface) Interface {
	return &multiplexedBreaker{newBreaker(), entries}
}

type multiplexedBreaker struct {
	*breaker
	entries []Interface
}

// Close closes the Done channel and releases resources associated with it.
func (br *multiplexedBreaker) Close() {
	br.closer.Do(func() {
		each(br.entries).Close()
		close(br.signal)
	})
}

// trigger starts listening all Done channels of multiplexed Breakers.
func (br *multiplexedBreaker) trigger() Interface {
	go func() {
		brs := make([]reflect.SelectCase, 0, len(br.entries))
		for _, br := range br.entries {
			brs = append(brs, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(br.Done())})
		}
		reflect.Select(brs)
		br.Close()
		atomic.StoreInt32(&br.released, 1)
	}()
	return br
}

type each []Interface

// Close closes all Done channels of a list of Breakers
// and releases resources associated with them.
func (list each) Close() {
	for _, br := range list {
		br.Close()
	}
}
