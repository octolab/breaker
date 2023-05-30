---
code: null
id: MDU6SXNzdWU0OTAzMjMxMTM=
databaseId: 490323113
number: 4
url: https://github.com/octolab/breaker/issues/4
title: "break by channel"
issueType: null
labels:
  - "type:feature"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-09-06T13:18:39Z
updatedAt: 2021-01-06T20:44:39Z
lastEditedAt: 2019-09-07T06:31:23Z
closedAt: 2021-01-06T20:44:39Z
---

# break by channel

for example

```go
func CloseByChannel(ctx context.Context, ch <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func () {
		<-ch // ensure that this event will happen
		cancel()
	}
	return ctx
}
```

<!-- 2019-09-07T06:34Z https://github.com/octolab/breaker/issues/4#issuecomment-529079174
```go
type Signal <-chan struct{}

func (s Signal) Done() <-chan struct{} { return s }

func (s Signal) Close() { close(s) /* how to close safely */ }

func (s Signal) WithContext(ctx context.Context) context.Context {
	...
}
```
-->

<!-- 2019-09-10T12:16Z https://github.com/octolab/breaker/issues/4#issuecomment-529907484
```go
type Signal struct{}

var Done Signal

type Notifier <-chan Signal

var notifier := make(Notifier)
notifier <- Done
notifier.Notify(Done)
```
-->

<!-- 2019-09-10T12:16Z https://github.com/octolab/breaker/issues/4#issuecomment-529907578
https://dave.cheney.net/2013/04/30/curious-channels
-->

<!-- 2020-03-30T19:59Z https://github.com/octolab/breaker/issues/4#issuecomment-606213625
```go
type ClosableContext interface {
	context.Context
	Close()
}

func Closable(ctx context.Context, cancel context.CancelFunc) ClosableContext {
	return &closable{ctx, cancel}
}

type closable struct {
	context.Context
	cancel context.CancelFunc
}

func (closable *closable) Close() {
	closable.cancel()
}
```
-->
