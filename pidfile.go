package pidfile

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

type Option func(*Pidfile) error

func NotEmpty(pf *Pidfile) error {
	if pf.path == "" {
		return errors.New("empty")
	}
	return nil
}

type Pidfile struct {
	path string
	pid  int
}

func (pf *Pidfile) Path() string {
	return pf.path
}

func New(path string) *Pidfile {
	return &Pidfile{
		path: path,
		pid:  os.Getpid(),
	}
}
func (pf *Pidfile) Generate(opts ...Option) error {
	for ix := range opts {
		if err := opts[ix](pf); err != nil {
			return err
		}
	}
	if pf.path != "" {
		if err := os.MkdirAll(filepath.Dir(pf.path), 0755); err != nil {
			return err
		}
		if _, err := os.Stat(pf.path); err == nil {
			return errors.New("file exists")
		}
		f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(strconv.Itoa(pf.pid)); err != nil {
			return err
		}
	}
	return nil
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
	return
}

type Cleanable interface{ Cleanup() }

func Generate(path string, opts ...Option) Cleanable {
	pf := New(path)
	if err := pf.Generate(opts...); err != nil {
		panic(err)
	}
	return pf
}
