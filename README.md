> # 🚧 breaker
>
> Flexible mechanism to make execution flow interruptible.

[![Awesome][icon_awesome]][page_awesome]
[![Build][icon_build]][page_build]
[![Coverage][icon_coverage]][page_coverage]
[![Quality][icon_quality]][page_quality]
[![Documentation][icon_docs]][page_docs]

## 💡 Idea

The breaker carries a cancellation signal to interrupt an action execution.

```go
interrupter := breaker.Multiplex(
	breaker.BreakByTimeout(time.Minute),
	breaker.BreakBySignal(os.Interrupt),
)
defer interrupter.Close()

<-interrupter.Done() // wait context cancellation, timeout or interrupt signal
```

Full description of the idea is available
[here](https://www.notion.so/octolab/breaker-77116e98fda74c28bd64e42bd440bbf3?r=0b753cbf767346f5a6fd51194829a2f3).

## 🏆 Motivation

I have to make [github.com/kamilsk/retry][retry] package:

```go
if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
	log.Fatal(err)
}
```

and [github.com/kamilsk/semaphore][semaphore] package:

```go
if err := semaphore.Acquire(breaker.BreakByTimeout(time.Minute), 5); err != nil {
	log.Fatal(err)
}
```

more consistent and reliable.

## 🤼‍♂️ How to

```go
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kamilsk/breaker"
)

func Handle(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	deadline, _ := time.ParseDuration(req.Header.Get("X-Timeout"))
	if deadline == 0 {
		deadline = 7 * time.Millisecond
	}
	interrupter := breaker.Multiplex(
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
			_, _ = buf.Write([]byte{b})
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
		...
	}()

	return outcome
}
```

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [dep][] or [go modules][gomod] to manage its version.

```bash
$ dep ensure -add github.com/kamilsk/breaker

$ go get -u github.com/kamilsk/breaker
```

---

made with ❤️ for everyone

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/breaker.svg?branch=master
[icon_coverage]:   https://api.codeclimate.com/v1/badges/1d703de640b4c6cfcd6f/test_coverage
[icon_docs]:       https://godoc.org/github.com/kamilsk/breaker?status.svg
[icon_quality]:    https://goreportcard.com/badge/github.com/kamilsk/breaker

[page_awesome]:    https://github.com/avelino/awesome-go#goroutines
[page_build]:      https://travis-ci.org/kamilsk/breaker
[page_coverage]:   https://codeclimate.com/github/kamilsk/breaker/test_coverage
[page_docs]:       https://godoc.org/github.com/kamilsk/breaker
[page_quality]:    https://goreportcard.com/report/github.com/kamilsk/breaker

[dep]:             https://golang.github.io/dep/
[gomod]:           https://github.com/golang/go/wiki/Modules
[promo]:           https://github.com/kamilsk/breaker
[retry]:           https://github.com/kamilsk/retry
[semaphore]:       https://github.com/kamilsk/semaphore
