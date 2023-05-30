---
code: null
id: MDU6SXNzdWU3ODA4NTczODI=
databaseId: 780857382
number: 12
url: https://github.com/octolab/breaker/issues/12
title: "breaker.ToContext"
issueType: null
labels:
  - "type:feature"
  - "help wanted"
  - "difficulty: easy"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2021-01-06T21:28:43Z
updatedAt: 2021-01-07T08:52:54Z
lastEditedAt: null
closedAt: 2021-01-07T08:52:54Z
---

# breaker.ToContext

- investigate the "hook based" model: when the breaker is released it calls cancel function
- if it's not possible add overhead: run goroutine above Done channel
