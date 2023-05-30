---
code: null
id: MDU6SXNzdWU0NzEyODQ2MzE=
databaseId: 471284631
number: 3
url: https://github.com/octolab/breaker/issues/3
title: "BreakByContext"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2019-07-22T19:48:08Z
updatedAt: 2019-07-23T07:08:37Z
lastEditedAt: null
closedAt: 2019-07-23T07:08:37Z
---

# BreakByContext

```go
interrupter := breaker.BreakByContext(context.WithTimeout(req.Context(), time.Minute))
```
