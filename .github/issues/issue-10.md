---
code: null
id: MDU6SXNzdWU2MTU0NDE1NjQ=
databaseId: 615441564
number: 10
url: https://github.com/octolab/breaker/issues/10
title: "investigate timer leaks"
issueType: null
labels:
  - "type:refactoring"
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
createdAt: 2020-05-10T17:49:30Z
updatedAt: 2021-01-07T08:06:07Z
lastEditedAt: null
closedAt: 2021-01-07T08:06:07Z
---

# investigate timer leaks

Whenever a `time.Timer` is Stopped its channel must be cleared using a
`select` if the result of the `Stop()` is false.

<!-- 2020-05-10T17:52Z https://github.com/octolab/breaker/issues/10#issuecomment-626364549
as example https://github.com/go-gitea/gitea/commit/3cf73929bc87cf5a6e05753b0c6ae13a833b1d9f, go-gitea/gitea#11333
-->
