// +build go1.13 integration

package breaker_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	. "github.com/kamilsk/breaker"
)

func Example_httpRequest() {
	const url = "http://example.com/"

	example := make(chan struct{})
	close(example)

	breaker := Multiplex(
		BreakByChannel(example),
		BreakBySignal(os.Interrupt, syscall.SIGINT, syscall.SIGTERM),
		BreakByTimeout(time.Hour),
	)
	defer breaker.Close()

	req, err := http.NewRequestWithContext(ToContext(breaker), http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	if _, err := http.DefaultClient.Do(req); errors.Is(err, context.Canceled) {
		fmt.Println("works well")
	}
	// output: works well
}
