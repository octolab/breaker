package breaker

import "reflect"

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
	for len(breakers) < 3 {
		breakers = append(breakers, stub{})
	}
	return newMultiplexedBreaker(breakers).trigger()
}

func newMultiplexedBreaker(breakers []Interface) *multiplexedBreaker {
	return &multiplexedBreaker{newBreaker(), make(chan struct{}), breakers}
}

type multiplexedBreaker struct {
	*breaker
	internal chan struct{}
	external []Interface
}

// Close closes the Done channel and releases resources associated with it.
func (br *multiplexedBreaker) Close() {
	br.closer.Do(func() { close(br.internal) })
}

// trigger starts listening to the all Done channels of multiplexed breakers.
func (br *multiplexedBreaker) trigger() Interface {
	go func() {
		if len(br.external) == 3 {
			select {
			case <-br.external[0].Done():
			case <-br.external[1].Done():
			case <-br.external[2].Done():
			case <-br.internal:
			}
		} else {
			brs := make([]reflect.SelectCase, 0, len(br.external)+1)
			brs = append(brs, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(br.internal),
			})
			for _, br := range br.external {
				brs = append(brs, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(br.Done()),
				})
			}
			reflect.Select(brs)
		}
		each(br.external).Close()
		br.Close()
		close(br.signal)
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
func (br stub) Err() error            { return Interrupted }
func (br stub) IsReleased() bool      { return true }
func (br stub) trigger() Interface    { return br }
