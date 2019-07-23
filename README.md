> # 🚧 breaker
>
> Flexible mechanism to make execution flow interruptible.

[![Build][icon_build]][page_build]
[![Quality][icon_quality]][page_quality]
[![Documentation][icon_docs]][page_docs]
[![Coverage][icon_coverage]][page_coverage]
[![Awesome][icon_awesome]][page_awesome]

## 💡 Idea

The breaker carries a cancellation signal to interrupt an action execution.

```go
interrupter := breaker.Multiplex(
	breaker.BreakByContext(context.WithTimeout(req.Context(), time.Minute)),
	breaker.BreakByDeadline(NewYear),
	breaker.BreakBySignal(os.Interrupt),
)
defer interrupter.Close()

<-interrupter.Done() // wait context cancellation, timeout or interrupt signal
```

Full description of the idea is available [here][design].

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
				io.Copy(rw, buf)
				return
			}
			buf.WriteByte(b)
		case <-interrupter.Done():
			rw.WriteHeader(http.StatusPartialContent)
			rw.Header().Set("Content-Range", fmt.Sprintf("bytes=0-%d", buf.Len()))
			io.Copy(rw, buf)
			return
		}
	}
}

func Work(ctx context.Context, _ struct{}) <-chan byte {
	outcome := make(chan byte, 1)

	go func() { ... }()

	return outcome
}
```

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [go modules](https://github.com/golang/go/wiki/Modules) or
[dep](https://golang.github.io/dep/) to manage its version.

```bash
$ go get -u github.com/kamilsk/breaker

$ dep ensure -add github.com/kamilsk/breaker
```

---

made with ❤️ for everyone

[icon_awesome]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:       https://travis-ci.org/kamilsk/breaker.svg?branch=master
[icon_coverage]:    https://api.codeclimate.com/v1/badges/1d703de640b4c6cfcd6f/test_coverage
[icon_docs]:        https://godoc.org/github.com/kamilsk/breaker?status.svg
[icon_quality]:     https://goreportcard.com/badge/github.com/kamilsk/breaker

[page_awesome]:     https://github.com/avelino/awesome-go#goroutines
[page_build]:       https://travis-ci.org/kamilsk/breaker
[page_coverage]:    https://codeclimate.com/github/kamilsk/breaker/test_coverage
[page_docs]:        https://godoc.org/github.com/kamilsk/breaker
[page_quality]:     https://goreportcard.com/report/github.com/kamilsk/breaker

[design]:           https://www.notion.so/octolab/breaker-77116e98fda74c28bd64e42bd440bbf3?r=0b753cbf767346f5a6fd51194829a2f3
[egg]:              https://github.com/kamilsk/egg
[promo]:            https://github.com/kamilsk/breaker
[retry]:            https://github.com/kamilsk/retry
[semaphore]:        https://github.com/kamilsk/semaphore

[tmp.docs]:         https://nicedoc.io/kamilsk/breaker?theme=dark
[tmp.history]:      https://github.githistory.xyz/kamilsk/breaker/blob/master/README.md
