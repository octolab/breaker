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
	// Close closes the Done channel and releases resources associated with it.
	Close()
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// Err returns a non-nil error if Done is closed and nil otherwise.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error
	// trigger is a private method to guarantee that the Breakers come from
	// this package and all of them return a valid Done channel.
	trigger() Interface
}
