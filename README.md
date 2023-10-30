## Pidfile

![test workflow](https://github.com/chg1f/pidfile/actions/workflows/test.yml/badge.svg?branch=master)

## Intro

Easily generate and cleanup golang programs pidfile

## Quick Start

```go
import "github.com/chg1f/pidfile"

func main(){
  defer pidfile.Generate("/var/run/example.pid").Cleanup()
  // do anything
}
```

> Documents visit https://pkg.go.dev/github.com/chg1f/pidfile
