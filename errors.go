package breaker

// Error defines the package errors.
type Error string

// Error returns the string representation of the error.
func (err Error) Error() string {
	return string(err)
}

// Interrupted is the error returned by the breaker
// when a cancellation signal occurred.
const Interrupted Error = "operation interrupted"
