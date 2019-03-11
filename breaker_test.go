package breaker_test

import (
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

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
	t.Run("passed deadline", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(-delta))
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := BreakByDeadline(time.Now().Add(time.Hour))
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
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

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
	t.Run("without signal", func(t *testing.T) {
		br := BreakBySignal()
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := BreakBySignal(os.Kill)
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
}

func TestBreakByTimeout(t *testing.T) {
	t.Run("valid timeout", func(t *testing.T) {
		br := BreakByTimeout(5 * delta)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start.Add(5*delta), time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
	t.Run("passed timeout", func(t *testing.T) {
		br := BreakByTimeout(-delta)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
	t.Run("close multiple times", func(t *testing.T) {
		br := BreakByTimeout(time.Hour)
		repeat(br.Close, 5)
		start := time.Now()
		<-br.Done()
		assert.WithinDuration(t, start, time.Now(), delta)

		delay(func() { assert.True(t, br.(interface{ Released() bool }).Released()) }, delta)
	})
}

func delay(action func(), d time.Duration) {
	time.Sleep(d)
	action()
}

func repeat(action func(), times int) {
	for range make([]struct{}, times) {
		action()
	}
}
