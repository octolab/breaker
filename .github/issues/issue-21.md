---
code: null
id: MDU6SXNzdWU3OTI4NTQxNTc=
databaseId: 792854157
number: 21
url: https://github.com/octolab/breaker/issues/21
title: "Chaining and nesting"
issueType: null
labels:
  - "type:feature"
  - "help wanted"
  - "difficulty: medium"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: OPEN
stateReason: null
duplicateOf: null
createdAt: 2021-01-24T17:52:21Z
updatedAt: 2021-01-25T18:58:55Z
lastEditedAt: null
closedAt: null
---

# Chaining and nesting

```go
interrupter := breaker.WithTimeout(breaker.WithSignal(...), time.Minute)
```
