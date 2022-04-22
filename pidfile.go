package pidfile

import (
	"errors"
	"os"
	"strconv"
	"sync/atomic"
)

type PIDFile struct {
	inited int32
	path   string
}

var (
	ErrDuplicated = errors.New("duplicated")
	ErrFileExists = errors.New("file exists")
)

func New(path string, constrants ...Constrant) (*PIDFile, error) {
	pf := PIDFile{
		path: path,
	}
	for _, constrant := range constrants {
		if err := constrant(); err != nil {
			return nil, err
		}
	}
	return &pf, nil
}
func (pf *PIDFile) Generate() error {
	if pf == nil {
		return nil
	}
	if !atomic.CompareAndSwapInt32(&pf.inited, 0, 1) {
		return ErrDuplicated
	} else if _, err := os.Stat(pf.path); !errors.Is(err, os.ErrNotExist) {
		return ErrFileExists
	}
	f, err := os.OpenFile(pf.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(strconv.Itoa(os.Getpid()))
	return nil
}
func (pf *PIDFile) Cleanup() error {
	if pf == nil {
		return nil
	}
	defer atomic.StoreInt32(&pf.inited, 0)
	return os.Remove(pf.path)
}

type Constrant func() error

const (
	ROOT = 0
)

func Permission(uid int) Constrant {
	return func() error {
		if os.Getuid() != uid {
			return errors.New("expect uid:" + strconv.Itoa(uid))
		}
		return nil
	}
}

func Generate(path string, constrants ...Constrant) interface {
	Cleanup() error
} {
	c, err := New(path, constrants...)
	if err != nil {
		panic(err)
	}
	err = c.Generate()
	if err != nil {
		panic(err)
	}
	return c
}
