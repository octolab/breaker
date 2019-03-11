package breaker

// Interface carries a cancellation signal to break an action execution.
//
// Example based on github.com/kamilsk/retry package:
//
//  if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
//  	log.Fatal(err)
//  }
//
// Example based on github.com/kamilsk/semaphore package:
//
//  if err := semaphore.Acquire(breaker.BreakByTimeout(time.Minute), 5); err != nil {
//  	log.Fatal(err)
//  }
//
type Interface interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// Close closes the Done channel and releases resources associated with it.
	Close()
	// trigger is a private method to guarantee that the Breakers come from
	// this package and all of them return a valid Done channel.
	trigger() Interface
}
