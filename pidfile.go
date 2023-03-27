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

func Must(path string, constrants ...Constrant) *PIDFile {
	pf, err := New(path, constrants...)
	if err != nil {
		panic(err)
	}
	return pf
}
func New(path string, constrants ...Constrant) (*PIDFile, error) {
	pf := PIDFile{
		path: path,
	}
	for ix := range constrants {
		if err := constrants[ix](); err != nil {
			return nil, err
		}
	}
	return &pf, nil
}
func (pf *PIDFile) Generate() (err error) {
	if pf == nil {
		return nil
	}
	pf.Do(func() {
		err := os.MkdirAll(filepath.Dir(pf.path), 0755)
		if err != nil {
			err = errors.New("mkdir failed: " + err.Error())
			return
		}
		if _, err := os.Stat(pf.path); !errors.Is(err, os.ErrNotExist) {
			err = ErrDuplicated
			return
		}
		f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			err = errors.New("open file failed: " + err.Error())
			return
		}
		defer f.Close()
		f.WriteString(strconv.Itoa(os.Getpid()))
	})
	return
}
func (pf *PIDFile) Cleanup() error {
	if pf == nil {
		return nil
	}
	return os.Remove(pf.path)
}
func Generate(path string, constrants ...Constrant) interface {
	Cleanup() error
} {
	pf := Must(path, constrants...)
	if err := pf.Generate(); err != nil {
		panic(err)
	}
	return pf
}
