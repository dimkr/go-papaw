```
  __ _  ___        _ __   __ _ _ __   __ ___      __
 / _` |/ _ \ _____| '_ \ / _` | '_ \ / _` \ \ /\ / /
| (_| | (_) |_____| |_) | (_| | |_) | (_| |\ V  V /
 \__, |\___/      | .__/ \__,_| .__/ \__,_| \_/\_/
 |___/            |_|         |_|
```

## Overview

go-papaw is a [Go](https://golang.org) package for [papaw](http://github.com/dimkr/papaw).

## Implementation

go-papaw is implemented using [cgo](https://golang.org/cmd/cgo).

## Usage

Just import the package!

```go
package main

import (
        _ "github.com/dimkr/go-papaw"
)
```

The module's init function locks the executable to RAM, then deletes the executable.

## Legal Information

go-papaw is free and unencumbered software released under the terms of the MIT license; see COPYING for the license text.

The ASCII art logo at the top was made using [FIGlet](http://www.figlet.org/).
