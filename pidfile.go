package pidfile

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
)

var (
	RaiseEmptyPath bool

	ErrNoPIDFile  = errors.New("no pidfile")
	ErrDuplicated = errors.New("duplicated")
)

type Constrant func() error
type PIDFile struct {
	generate int32
	cleanup  int32

	path string
	cons []Constrant
}

func Generate(path string, cons ...Constrant) interface{ Cleanup() error } {
	return New(path, cons...).Generate()
}
func New(path string, cons ...Constrant) *PIDFile {
	return &PIDFile{
		path: path,
		cons: cons,

		cleanup: 1,
	}
}
func (pf *PIDFile) Generate() interface{ Cleanup() error } {
	if pf == nil {
		panic(ErrNoPIDFile)
	}
	if pf.path == "" {
		if RaiseEmptyPath {
			panic(ErrNoPIDFile)
		}
		return pf
	}
	if atomic.CompareAndSwapInt32(&pf.generate, 0, 1) {
		for ix := range pf.cons {
			if err := pf.cons[ix](); err != nil {
				panic(err)
			}
		}
		err := os.MkdirAll(filepath.Dir(pf.path), 0755)
		if err != nil {
			panic(err)
		}
		if _, err := os.Stat(pf.path); !errors.Is(err, os.ErrNotExist) {
			panic(ErrDuplicated)
		}
		f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(strconv.Itoa(os.Getpid()))
		atomic.StoreInt32(&pf.cleanup, 0)
	}
	return pf
}
func (pf *PIDFile) Cleanup() error {
	if pf == nil {
		return ErrNoPIDFile
	}
	if pf.path == "" {
		if RaiseEmptyPath {
			return ErrNoPIDFile
		}
		return nil
	}
	if atomic.LoadInt32(&pf.generate) == 1 && atomic.CompareAndSwapInt32(&pf.cleanup, 0, 1) {
		return os.Remove(pf.path)
	}
	return ErrNoPIDFile
}
