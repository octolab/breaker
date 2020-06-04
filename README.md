> # 🚧 breaker
>
> Flexible mechanism to make execution flow interruptible.

[![Build][build.icon]][build.page]
[![Documentation][docs.icon]][docs.page]
[![Quality][quality.icon]][quality.page]
[![Template][template.icon]][template.page]
[![Coverage][coverage.icon]][coverage.page]
[![Awesome][awesome.icon]][awesome.page]

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

A full description of the idea is available [here][design.page].

## 🏆 Motivation

I have to make [retry][] package:

```go
if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
	log.Fatal(err)
}
```

and [semaphore][] package:

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

var NewYear = time.Time{}.AddDate(time.Now().Year(), 0, 0)

func Handle(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	deadline, _ := time.ParseDuration(req.Header.Get("X-Timeout"))
	interrupter := breaker.Multiplex(
		breaker.BreakByContext(context.WithTimeout(ctx, deadline)),
		breaker.BreakByDeadline(NewYear),
		breaker.BreakBySignal(os.Interrupt),
	)
	defer interrupter.Close()

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
You can use [go modules](https://github.com/golang/go/wiki/Modules) to manage its version.

```bash
$ go get github.com/kamilsk/breaker@latest
```

---

made with ❤️ for everyone

[awesome.icon]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[awesome.page]:     https://github.com/avelino/awesome-go#goroutines
[build.icon]:       https://travis-ci.org/kamilsk/breaker.svg?branch=master
[build.page]:       https://travis-ci.org/kamilsk/breaker
[coverage.icon]:    https://api.codeclimate.com/v1/badges/1d703de640b4c6cfcd6f/test_coverage
[coverage.page]:    https://codeclimate.com/github/kamilsk/breaker/test_coverage
[design.page]:      https://www.notion.so/octolab/breaker-77116e98fda74c28bd64e42bd440bbf3?r=0b753cbf767346f5a6fd51194829a2f3
[docs.page]:        https://pkg.go.dev/github.com/kamilsk/breaker
[docs.icon]:        https://img.shields.io/badge/docs-pkg.go.dev-blue
[promo.page]:       https://github.com/kamilsk/breaker
[quality.icon]:     https://goreportcard.com/badge/github.com/kamilsk/breaker
[quality.page]:     https://goreportcard.com/report/github.com/kamilsk/breaker
[template.page]:    https://github.com/octomation/go-module
[template.icon]:    https://img.shields.io/badge/template-go--module-blue

[retry]:            https://github.com/kamilsk/retry
[semaphore]:        https://github.com/kamilsk/semaphore

[tmp.docs]:         https://nicedoc.io/kamilsk/breaker?theme=dark
[tmp.history]:      https://github.githistory.xyz/kamilsk/breaker/blob/master/README.md
