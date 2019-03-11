package breaker_test

import (
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
	"github.com/stretchr/testify/assert"
)

func TestMultiplex(t *testing.T) {
	t.Run("with breakers", func(t *testing.T) {
		br := Multiplex(BreakByTimeout(5*delta), BreakByDeadline(time.Now().Add(time.Hour)))
		defer br.Close()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start.Add(5*delta), time.Now(), delta)
	})
	t.Run("without breakers", func(t *testing.T) {
		br := Multiplex()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := Multiplex(BreakByTimeout(time.Hour))
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
}

func TestMultiplexTwo(t *testing.T) {
	br := MultiplexTwo(
		BreakByDeadline(time.Now().Add(-delta)),
		BreakByTimeout(time.Hour),
	)
	start := time.Now()
	<-br.Done()
	assert.WithinDuration(t, start, time.Now(), delta)
}

func TestMultiplexThree(t *testing.T) {
	br := MultiplexThree(
		BreakByDeadline(time.Now().Add(-delta)),
		BreakBySignal(os.Kill),
		BreakByTimeout(time.Hour),
	)
	start := time.Now()
	<-br.Done()
	assert.WithinDuration(t, start, time.Now(), delta)
}
