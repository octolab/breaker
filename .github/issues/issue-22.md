---
code: null
id: MDU6SXNzdWU3OTM2NDI4Njg=
databaseId: 793642868
number: 22
url: https://github.com/octolab/breaker/issues/22
title: "make possible to reduce fluky tests"
issueType: null
labels:
  - "type:bug"
  - "type:critical"
  - "help wanted"
  - "difficulty: medium"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: OPEN
stateReason: null
duplicateOf: null
createdAt: 2021-01-25T18:58:30Z
updatedAt: 2021-01-26T20:03:20Z
lastEditedAt: null
closedAt: null
---

# make possible to reduce fluky tests

```
--- FAIL: TestBreakBySignal (0.00s)
    --- FAIL: TestBreakBySignal/with_signal (0.02s)
        breaker_test.go:216: max difference between 2021-01-24 21:27:21.89766 +0300 MSK m=+0.017424561 and 2021-01-24 21:27:21.908426 +0300 MSK m=+0.028190544 allowed is 10ms, but difference was -10.765983ms
FAIL
```

```
$ make test
--- FAIL: TestMultiplexThree (0.01s)
    deprecated_test.go:30: a breaker is not released
    deprecated_test.go:30: a breaker has inconsistent state
    deprecated_test.go:30: a breaker has no error
FAIL
FAIL    github.com/kamilsk/breaker      0.884s
FAIL
make: *** [test] Error 1
```

<!-- 2021-01-26T20:03Z https://github.com/octolab/breaker/issues/22#issuecomment-767793899
It shows fundamental error: if Done is closed, the Err **must** return Interrupted, but it is undefined, because `released int32` stores 1 only if goroutine is released
-->
