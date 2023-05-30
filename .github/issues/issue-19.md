---
code: null
id: MDU6SXNzdWU3OTI4NTM1MTU=
databaseId: 792853515
number: 19
url: https://github.com/octolab/breaker/issues/19
title: "Signal trap"
issueType: null
labels:
  - "type:feature"
  - "help wanted"
  - "difficulty: easy"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: OPEN
stateReason: null
duplicateOf: null
createdAt: 2021-01-24T17:49:27Z
updatedAt: 2021-01-25T18:58:54Z
lastEditedAt: 2021-01-24T17:50:25Z
closedAt: null
---

# Signal trap

- [ ] alias
```go
interrupter := breaker.Termination()
// breaker.BreakBySignal(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
```
- [ ] error
```go
errors.Is(interrupter.Err(), breaker.SignalTrapped)
// SignalTrapped Error = "signal trapped"
```
