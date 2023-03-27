package pidfile

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	ErrDuplicated = errors.New("duplicated")
)

type PIDFile struct {
	sync.Once
	path string
}

type Constrant func() error

func New(path string, constrants ...Constrant) (*PIDFile, error) {
	pf := PIDFile{
		path: path,
	}
	for ix := range constrants {
		if err := constrants[ix](); err != nil {
			return nil, err
		}
	}
	err := os.MkdirAll(filepath.Dir(pf.path), 0755)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(pf.path); !errors.Is(err, os.ErrNotExist) {
		return nil, ErrDuplicated
	}
	f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	f.WriteString(strconv.Itoa(os.Getpid()))
	return &pf, nil
}
func (pf *PIDFile) Cleanup() error {
	if pf == nil {
		return nil
	}
	pf.Do(func() {
		os.Remove(pf.path)
	})
	return nil
}
func Generate(path string, constrants ...Constrant) interface {
	Cleanup() error
} {
	pf, err := New(path, constrants...)
	if err != nil {
		panic(err)
	}
	return pf
}
