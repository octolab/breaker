package breaker_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/breaker"
	"github.com/stretchr/testify/assert"
)

var delta = 10 * time.Millisecond

func TestBreakByDeadline(t *testing.T) {
	t.Run("future deadline", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(5 * delta))
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start.Add(5*delta), time.Now(), delta)
	})
	t.Run("passed deadline", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(-delta))
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(time.Hour))
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
}

func TestBreakBySignal(t *testing.T) {
	t.Run("with signal", func(t *testing.T) {
		br := BreakBySignal(os.Interrupt)
		start := time.Now()
		go func() {
			proc, err := os.FindProcess(os.Getpid())
			assert.NoError(t, err)
			assert.NoError(t, proc.Signal(os.Interrupt))
		}()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("without signal", func(t *testing.T) {
		br := BreakBySignal()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := BreakBySignal(os.Kill)
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
}

func TestBreakByTimeout(t *testing.T) {
	t.Run("valid timeout", func(t *testing.T) {
		br := BreakByTimeout(5 * delta)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start.Add(5*delta), time.Now(), delta)
	})
	t.Run("passed timeout", func(t *testing.T) {
		br := BreakByTimeout(-delta)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := BreakByTimeout(time.Hour)
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
}

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

func TestWithContext(t *testing.T) {
	t.Run("active breaker", func(t *testing.T) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), 5*delta)
			br, _       = WithContext(ctx)
		)
		defer cancel()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start.Add(5*delta), time.Now(), delta)
	})
	t.Run("closed breaker", func(t *testing.T) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), -delta)
			br, _       = WithContext(ctx)
		)
		defer cancel()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("released breaker", func(t *testing.T) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), time.Hour)
			br, _       = WithContext(ctx)
		)
		defer cancel()
		br.Close()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
	t.Run("canceled parent", func(t *testing.T) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), time.Hour)
			br, _       = WithContext(ctx)
		)
		cancel()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)
	})
}

func repeat(action func(), times int) {
	for range make([]struct{}, times) {
		action()
	}
}
