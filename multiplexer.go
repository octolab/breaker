package breaker

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
