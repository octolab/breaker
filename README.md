> # 🚧 breaker
>
> Flexible mechanism to make your code breakable.

[![Awesome][icon_awesome]][awesome]
[![Patreon][icon_patreon]][support]
[![Build][icon_build]][build]
[![Quality][icon_quality]][quality]
[![Coverage][icon_coverage]][quality]
[![GoDoc][icon_docs]][docs]
[![License][icon_license]][license]

A Breaker carries a cancellation signal to break an action execution.

Example based on [github.com/kamilsk/retry][retry] package:

```go
if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
	log.Fatal(err)
}
```

Example based on [github.com/kamilsk/semaphore][semaphore] package:

```go
if err := semaphore.Acquire(breaker.BreakByTimeout(time.Minute), 5); err != nil {
	log.Fatal(err)
}
```

Complex example:

```go
interrupter := breaker.Multiplex(
	func () breaker.Interface {
		br, _ := breaker.WithContext(request.Context())
		return br
	}()
	breaker.BreakByTimeout(time.Minute),
	breaker.BreakBySignal(os.Interrupt),
)
defer interrupter.Close()

<-interrupter.Done() // wait context cancellation, timeout or interrupt signal
```

---

[![@kamilsk][icon_tw_author]][author]
[![@octolab][icon_tw_sponsor]][sponsor]

made with ❤️ by [OctoLab][octolab]

[awesome]:         https://github.com/avelino/awesome-go#goroutines
[build]:           https://travis-ci.org/kamilsk/breaker
[docs]:            https://godoc.org/github.com/kamilsk/breaker
[license]:         LICENSE
[promo]:           https://github.com/kamilsk/breaker
[quality]:         https://scrutinizer-ci.com/g/kamilsk/breaker/?branch=master

[retry]:           https://github.com/kamilsk/retry
[semaphore]:       https://github.com/kamilsk/semaphore

[author]:          https://twitter.com/ikamilsk
[octolab]:         https://www.octolab.org/
[sponsor]:         https://twitter.com/octolab_inc
[support]:         https://www.patreon.com/octolab

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/breaker.svg?branch=master
[icon_coverage]:   https://scrutinizer-ci.com/g/kamilsk/breaker/badges/coverage.png?b=master
[icon_docs]:       https://godoc.org/github.com/kamilsk/breaker?status.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_quality]:    https://scrutinizer-ci.com/g/kamilsk/breaker/badges/quality-score.png?b=master
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
[icon_twitter]:    https://img.shields.io/twitter/url/http/shields.io.svg?style=social
