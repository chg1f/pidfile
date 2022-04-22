PIDFile
---

![test workflow](https://github.com/chg1f/pidfile/actions/workflows/test.yml/badge.svg?branch=master)

## Intro

Handle golang pidfile generate and cleanup

## Quick Start

```go
import "github.com/chg1f/pidfile"

func main(){
  defer pidfile.Generate("/var/run/example.pid").Cleanup()
  // do anything
}
```

> Documents visit https://pkg.go.dev/github.com/chg1f/pidfile
