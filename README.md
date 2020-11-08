> # 🚧 breaker [![Awesome][awesome.icon]][awesome.page]
>
> Flexible mechanism to make execution flow interruptible.

[![Build][build.icon]][build.page]
[![Documentation][docs.icon]][docs.page]
[![Quality][quality.icon]][quality.page]
[![Template][template.icon]][template.page]
[![Coverage][coverage.icon]][coverage.page]
[![Mirror][mirror.icon]][mirror.page]

## 💡 Idea

The breaker carries a cancellation signal to interrupt an action execution.

```go
var NewYear = time.Time{}.AddDate(time.Now().Year(), 0, 0)

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

more consistent and reliable. Additionally, I want to implement a Graceful Shutdown
and Circuit Breaker on the same mechanism.

## 🤼‍♂️ How to

### Do HTTP request with retries

...

### Graceful Shutdown HTTP server

...

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [go modules](https://github.com/golang/go/wiki/Modules) to manage its version.

```bash
$ go get github.com/kamilsk/breaker@latest
```

## 🤲 Outcomes

### Console tool to execute commands for a limited time

The example shows how to execute console commands for ten minutes.

```bash
$ date
# Thu Jan  7 21:02:21
$ breakit after 10m -- database run --port=5432
$ breakit after 10m -- server run --port=8080
$ breakit ps
# +--------------------------+---------------------+
# | Process                  | Done                |
# +--------------------------+---------------------+
# | database run --port=5432 | Thu Jan  7 21:12:24 |
# | server run --port=8080   | Thu Jan  7 21:12:31 |
# +--------------------------+---------------------+
# |                    Total |                   2 |
# +--------------------------+---------------------+
```

See more details [here][cli].

---

made with ❤️ for everyone

[build.page]:       https://travis-ci.com/kamilsk/breaker
[build.icon]:       https://travis-ci.com/kamilsk/breaker.svg?branch=master
[coverage.page]:    https://codeclimate.com/github/kamilsk/breaker/test_coverage
[coverage.icon]:    https://api.codeclimate.com/v1/badges/1d703de640b4c6cfcd6f/test_coverage
[design.page]:      https://www.notion.so/octolab/breaker-77116e98fda74c28bd64e42bd440bbf3?r=0b753cbf767346f5a6fd51194829a2f3
[docs.page]:        https://pkg.go.dev/github.com/kamilsk/breaker
[docs.icon]:        https://img.shields.io/badge/docs-pkg.go.dev-blue
[promo.page]:       https://github.com/kamilsk/breaker
[quality.page]:     https://goreportcard.com/report/github.com/kamilsk/breaker
[quality.icon]:     https://goreportcard.com/badge/github.com/kamilsk/breaker
[template.page]:    https://github.com/octomation/go-module
[template.icon]:    https://img.shields.io/badge/template-go--module-blue
[mirror.page]:      https://bitbucket.org/kamilsk/breaker
[mirror.icon]:      https://img.shields.io/badge/mirror-bitbucket-blue

[awesome.page]:     https://github.com/avelino/awesome-go#goroutines
[awesome.icon]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg

[cli]:              https://github.com/octolab/breakit
[retry]:            https://github.com/kamilsk/retry
[semaphore]:        https://github.com/kamilsk/semaphore
