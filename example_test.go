// +build go1.13 integration

package breaker_test

import (
	"context"
	"errors"
	"fmt"
	"net"
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

	//nolint:bodyclose
	if _, err := http.DefaultClient.Do(req); errors.Is(err, context.Canceled) && errors.Is(breaker.Err(), Interrupted) {
		fmt.Println("works well")
	}
	// output: works well
}

func Example_gracefulShutdown() {
	example := make(chan struct{})

	breaker := Multiplex(
		BreakBySignal(os.Interrupt, syscall.SIGINT, syscall.SIGTERM),
		BreakByTimeout(250*time.Millisecond),
	)
	defer breaker.Close()

	server := http.Server{
		BaseContext: func(net.Listener) context.Context {
			return ToContext(breaker)
		},
	}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		close(example)
	}()

	<-breaker.Done()
	if err := server.Shutdown(context.TODO()); err == nil && errors.Is(breaker.Err(), Interrupted) {
		fmt.Println("works well")
	}
	<-example
	// output: works well
}
