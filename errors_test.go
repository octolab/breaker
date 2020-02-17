package breaker_test

import (
	"testing"

	. "github.com/kamilsk/breaker"
)

func TestError(t *testing.T) {
	var err error = Interrupted
	if err.Error() != string(Interrupted) {
		t.Error("fail error assertion")
	}
}
