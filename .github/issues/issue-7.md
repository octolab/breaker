---
code: null
id: MDU6SXNzdWU1MzgwNTc3MDU=
databaseId: 538057705
number: 7
url: https://github.com/octolab/breaker/issues/7
title: "simplify implementation"
issueType: null
labels:
  - "type:refactoring"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-12-15T14:21:48Z
updatedAt: 2021-01-06T19:30:39Z
lastEditedAt: null
closedAt: 2021-01-06T19:30:39Z
---

# simplify implementation

- remove trigger
- remove released
- ? do not run background goroutine
- use simplified version of context wrapping
```go

// Wrap wraps the context and its cancel function into BreakCloser.
func Wrap(ctx context.Context, cancel context.CancelFunc) BreakCloser {
	return &wrapper{ctx, cancel}
}

type wrapper struct {
	context.Context
	cancel context.CancelFunc
}

// Close closes the Done channel and releases resources associated with it.
func (breaker *wrapper) Close() {
	breaker.cancel()
}
```
