```
  __ _  ___        _ __   __ _ _ __   __ ___      __
 / _` |/ _ \ _____| '_ \ / _` | '_ \ / _` \ \ /\ / /
| (_| | (_) |_____| |_) | (_| | |_) | (_| |\ V  V /
 \__, |\___/      | .__/ \__,_| .__/ \__,_| \_/\_/
 |___/            |_|         |_|
```

## Overview

go-papaw is a [Go](https://golang.org) module for [papaw](http://github.com/dimkr/papaw).

## Implementation

go-papaw is implemented using [cgo](https://golang.org/cmd/cgo).

## Usage

Just import the module!

```go
import (
        _ "github.com/dimkr/go-papaw"
)
```

The module's init function locks the executable to RAM, then deletes the executable.

Before:

```
$ cat << EOF > main.go
package main
import "time"
func main() {
    time.Sleep(time.Hour)
}
EOF
$ go build -o /tmp/example
$ /tmp/example &
$ pid=$!
[1] 28165
$ grep example /proc/$pid/maps
00010000-00065000 r-xp 00000000 00:26 179920                             /tmp/example <--
00070000-000e6000 r--p 00060000 00:26 179920                             /tmp/example <--
000f0000-000f4000 rw-p 000e0000 00:26 179920                             /tmp/example <--
$ ls -la /proc/$pid/exe
lrwxrwxrwx. 1 user user 0 Apr  7 12:10 /proc/28165/exe -> /tmp/example                <--
$ head -c 4 /proc/$pid/exe
ELF                                                                                   <--
$ ls example
example                                                                               <--
```

After:

```
$ cat << EOF > main.go
package main
import (
    _ "github.com/dimkr/go-papaw"                                                     <--
    "time"
)
func main() {
    time.Sleep(time.Hour)
}
EOF
$ go build -o /tmp/example
$ /tmp/example &
$ pid=$!
[2] 27812
$ grep example /proc/$pid/maps
$ ls -la /proc/$pid/exe
lrwxrwxrwx. 1 user user 0 Apr  7 12:11 /proc/27811/exe -> '/tmp/example (deleted)'
$ head -c 4 /proc/$pid/exe
$ ls example
ls: cannot access 'example': No such file or directory
```

## Legal Information

go-papaw is free and unencumbered software released under the terms of the MIT license; see COPYING for the license text.

The ASCII art logo at the top was made using [FIGlet](http://www.figlet.org/).
