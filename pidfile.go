package pidfile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	ErrDuplicated = errors.New("duplicated")
)

type Cleanable interface{ Cleanup() error }

func Generate(path string, conflicts ...func() error) Cleanable {
	pf := &pidfile{path: path}
	return pf.Generate(conflicts...)
}

type pidfile struct {
	path string
	once sync.Once
}

func (pf *pidfile) Generate(conflicts ...func() error) Cleanable {
	if pf == nil || pf.path == "" {
		return nil
	}
	for ix := range conflicts {
		if err := conflicts[ix](); err != nil {
			panic(err)
		}
	}
	err := os.MkdirAll(filepath.Dir(pf.path), 0755)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(pf.path); !errors.Is(err, os.ErrNotExist) {
		f, err := os.Open(pf.path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		bs, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		if strconv.Itoa(os.Getpid()) == string(bs) {
			return pf
		}
		panic(ErrDuplicated)
	}
	f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(strconv.Itoa(os.Getpid()))
	return pf
}
func (pf *pidfile) Cleanup() error {
	if pf == nil || pf.path == "" {
		return nil
	}
	f, err := os.Open(pf.path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if strconv.Itoa(os.Getpid()) != string(bs) {
		return nil
	}
	return os.Remove(pf.path)
}
