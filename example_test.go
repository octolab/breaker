package breaker_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kamilsk/breaker"
)

func Example() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/work", nil)
	req.Header.Set("X-Timeout", "50ms")

	http.HandlerFunc(Handle).ServeHTTP(rec, req)

	fmt.Println("status:", http.StatusText(rec.Result().StatusCode))
	fmt.Printf("body: ")
	_, _ = io.Copy(os.Stdout, strings.NewReader(truncate(rec.Body.String(), 28)))
	// Output:
	// status: Partial Content
	// body: breaker 🚧 Flexible mechanism
}

func Handle(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	deadline, _ := time.ParseDuration(req.Header.Get("X-Timeout"))
	interrupter := breaker.Multiplex(
		func() breaker.Interface {
			br, _ := breaker.WithContext(req.Context())
			return br
		}(),
		breaker.BreakByTimeout(deadline),
		breaker.BreakBySignal(os.Interrupt),
	)

	buf, work := bytes.NewBuffer(nil), Work(ctx, struct{}{})
	for {
		select {
		case b, ok := <-work:
			if !ok {
				rw.WriteHeader(http.StatusOK)
				_, _ = io.Copy(rw, buf)
				return
			}
			_ = buf.WriteByte(b)
		case <-interrupter.Done():
			rw.WriteHeader(http.StatusPartialContent)
			rw.Header().Set("Content-Range", fmt.Sprintf("bytes=0-%d", buf.Len()))
			_, _ = io.Copy(rw, buf)
			return
		}
	}
}

func Work(ctx context.Context, _ struct{}) <-chan byte {
	outcome := make(chan byte, 1)

	go func() {
		defer close(outcome)
		for _, b := range []byte("breaker 🚧 Flexible mechanism to make execution flow interruptible.") {
			time.Sleep(time.Millisecond)
			select {
			case <-ctx.Done():
				return
			case outcome <- b:
			}
		}
	}()

	return outcome
}

func truncate(raw string, len int) string {
	if max := utf8.RuneCountInString(raw); max < len {
		len = max
	}
	var chars int
	for pos := range raw {
		if chars >= len {
			raw = raw[:pos]
			break
		}
		chars++
	}
	return raw
}
