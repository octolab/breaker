---
code: null
id: MDU6SXNzdWU1NjU3MjE0NDg=
databaseId: 565721448
number: 8
url: https://github.com/octolab/breaker/issues/8
title: "extend interface"
issueType: null
labels:
  - "type:task"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-02-15T09:14:05Z
updatedAt: 2021-01-01T20:03:47Z
lastEditedAt: null
closedAt: 2021-01-01T20:03:47Z
---

# extend interface

```go
// Interface carries a cancellation signal to break an action execution.
//
// Example based on github.com/kamilsk/retry package:
//
//  if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action); err != nil {
//  	log.Fatal(err)
//  }
//
// Example based on github.com/kamilsk/semaphore package:
//
//  if err := semaphore.Acquire(breaker.BreakByTimeout(time.Minute), 5); err != nil {
//  	log.Fatal(err)
//  }
//
type Interface interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// Close closes the Done channel and releases resources associated with it.
	Close()
	// trigger is a private method to guarantee that the Breakers come from
	// this package and all of them return a valid Done channel.
	trigger() Interface

	Error() error
}
```
