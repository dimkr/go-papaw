```
  __ _  ___        _ __   __ _ _ __   __ ___      __
 / _` |/ _ \ _____| '_ \ / _` | '_ \ / _` \ \ /\ / /
| (_| | (_) |_____| |_) | (_| | |_) | (_| |\ V  V /
 \__, |\___/      | .__/ \__,_| .__/ \__,_| \_/\_/
 |___/            |_|         |_|
```

## Overview

go-papaw is a [Go](https://golang.org) module for [papaw](http://github.com/dimkr/papaw).

go-papaw does two things:

1. It makes executables smaller and harder to reverse-engineer
2. It allows executables to run from RAM and delete themselves

## Implementation

go-papaw is implemented using [cgo](https://golang.org/cmd/cgo).

## Packing

To pack (compress and obfuscate) an executable:

```
$ go get -u github.com/dimkr/go-papaw
$ go run github.com/dimkr/go-papaw/cmd/pack -input /tmp/example -output /tmp/example-compressed
```

Before:

```
$ /tmp/example
Hello world!
$ du -h /tmp/example
2.0M    /tmp/example
```

After:
```
$ /tmp/example-compressed
Hello world!
$ du -h /tmp/example-compressed
992K    /tmp/example-compressed
```

The `pack` tool offers two compression algorithms: `-algo deflate` (the default) and `-algo lzma`.

The compression ratio will vary based on the chosen algorithm (LZMA has a better compression ratio), the input executable content (the length and number of identical sequences) and the target CPU architecture (RISC or CISC, fixed- or variable-sized instructions).

However, the compression ratio tends to be __at least 50%__:

```
$ CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-s -w" -o /tmp/example
$ go run github.com/dimkr/go-papaw/cmd/pack -input /tmp/example -output /tmp/example-packed
2020/04/07 16:07:10 Done: output is 62.62% smaller.
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /tmp/example
$ go run github.com/dimkr/go-papaw/cmd/pack -input /tmp/example -output /tmp/example-packed
2020/04/07 16:07:11 Done: output is 58.45% smaller.
$ go run github.com/dimkr/go-papaw/cmd/pack -input /tmp/example -output /tmp/example-packed -algo lzma
2020/04/07 16:07:12 Done: output is 66.89% smaller.
$ CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o /tmp/example
$ go run github.com/dimkr/go-papaw/cmd/pack -input /tmp/example -output /tmp/example-packed
2020/04/07 16:07:13 Done: output is 56.11% smaller.
$ CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o /tmp/example
$ go run github.com/dimkr/go-papaw/cmd/pack -input /tmp/example -output /tmp/example-packed
2020/04/07 16:07:14 Done: output is 62.78% smaller.
```

## Fileless Execution

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
lrwxrwxrwx. 1 user user 0 Apr  7 12:11 /proc/27812/exe -> '/tmp/example (deleted)'
$ head -c 4 /proc/$pid/exe
$ ls example
ls: cannot access 'example': No such file or directory
```

## Legal Information

go-papaw is free and unencumbered software released under the terms of the MIT license; see COPYING for the license text.

The ASCII art logo at the top was made using [FIGlet](http://www.figlet.org/).
