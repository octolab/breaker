// +build go1.13

package breaker_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/kamilsk/breaker"
)

func TestError(t *testing.T) {
	if !errors.Is(fmt.Errorf("%w", fmt.Errorf("%w", Interrupted)), Interrupted) {
		t.Error("fail error assertion")
	}
}
