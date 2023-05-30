---
code: null
id: MDU6SXNzdWU3OTE0OTQwNjY=
databaseId: 791494066
number: 16
url: https://github.com/octolab/breaker/issues/16
title: "MultiplexTwo and MultiplexThree has goroutine leak"
issueType: null
labels:
  - "type:bug"
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
createdAt: 2021-01-21T21:07:43Z
updatedAt: 2021-01-22T06:21:43Z
lastEditedAt: null
closedAt: 2021-01-22T06:21:43Z
---

# MultiplexTwo and MultiplexThree has goroutine leak

if input breakers never have Close() call the goroutine is idle forever:

- https://github.com/kamilsk/breaker/blob/5a4f777c2dc9e9bde5c1eb3ef87657caea4e32c5/multiplexer.go#L39-L45
- https://github.com/kamilsk/breaker/blob/5a4f777c2dc9e9bde5c1eb3ef87657caea4e32c5/multiplexer.go#L63-L70
