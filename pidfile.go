package pidfile

import (
	"os"
	"path/filepath"
	"strconv"
)

type Pidfile struct {
	path string
	pid  int
}

func Generate(opts ...func(*Pidfile)) interface{ Cleanup() } {
	pf := new(Pidfile)
	pf.path = filepath.Join(os.TempDir(), filepath.Base(os.Args[0])+".pid")
	pf.pid = os.Getpid()
	for i := range opts {
		opts[i](pf)
	}
	if pf.path != "" {
		if err := os.MkdirAll(filepath.Dir(pf.path), 0755); err != nil {
			panic(err)
		}
		if _, err := os.Stat(pf.path); err == nil {
			panic(os.ErrExist)
		}
		f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		if _, err := f.WriteString(strconv.Itoa(pf.pid)); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
	return pf
}

func (pf *Pidfile) Cleanup() {
	if pf.path == "" {
		return
	}
	bs, err := os.ReadFile(pf.path)
	if err != nil {
		panic(err)
	}
	if strconv.Itoa(pf.pid) != string(bs) {
		return
	}
	if err := os.Remove(pf.path); err != nil {
		panic(err)
	}
}

func Path(s string) func(*Pidfile) {
	return func(pf *Pidfile) {
		pf.path = s
	}
}
