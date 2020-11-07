package breaker

type extended interface {
	Interface
	IsReleased() bool
}

var (
	_ extended = new(breaker)
	_ extended = new(signalBreaker)
	_ extended = new(channelBreaker)
	_ extended = new(contextBreaker)
	_ extended = new(timeoutBreaker)
)
