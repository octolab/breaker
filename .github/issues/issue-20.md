---
code: null
id: MDU6SXNzdWU3OTI4NTM5NjY=
databaseId: 792853966
number: 20
url: https://github.com/octolab/breaker/issues/20
title: "Wait signal"
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
createdAt: 2021-01-24T17:51:21Z
updatedAt: 2021-01-25T18:58:55Z
lastEditedAt: null
closedAt: null
---

# Wait signal

```go
breaker.Wait(interrupter)
// <-interrupter.Done()
```
