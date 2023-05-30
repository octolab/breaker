---
code: null
id: MDU6SXNzdWU4MDYxMzU4NzY=
databaseId: 806135876
number: 24
url: https://github.com/octolab/breaker/issues/24
title: "strict guarantees for Close"
issueType: null
labels:
  - "help wanted"
  - "difficulty: medium"
assignees:
  - kamilsk
milestone: "[[milestone-2]]"
parent: null
issueFields: {}
state: OPEN
stateReason: null
duplicateOf: null
createdAt: 2021-02-11T07:07:44Z
updatedAt: 2021-02-11T07:09:53Z
lastEditedAt: 2021-02-11T07:09:53Z
closedAt: null
---

# strict guarantees for Close

see the context package

```go
package main

import (
	"fmt"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()

	if ctx.Err() == nil {
		panic("not strict guarantees")
	}

	fmt.Println("strict guarantees")
}
```

I need the same thing.
