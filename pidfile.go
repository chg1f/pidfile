package pidfile

import (
	"errors"
	"os"
	"strconv"
	"sync/atomic"
)

var (
	inited int32
)

type Option func() error

type PIDFile struct {
	inited int32
	path   string
}

func New(path string, options ...Option) (*PIDFile, error) {
	for _, option := range options {
		if err := option(); err != nil {
			return nil, err
		}
	}
	return &PIDFile{
		path: path,
	}, nil
}
func (pf *PIDFile) Generate() error {
	if !atomic.CompareAndSwapInt32(&pf.inited, 0, 1) {
		return errors.New("duplicated")
	} else if _, err := os.Stat(pf.path); errors.Is(err, os.ErrExist) {
		return os.ErrExist
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
	if pf != nil {
		defer atomic.StoreInt32(&pf.inited, 0)
		return os.Remove(pf.path)
	}
	return nil
}

const (
	ROOT = 0
)

func Permission(uid int) Option {
	return func() error {
		if os.Getuid() != uid {
			return errors.New("expect uid:" + strconv.Itoa(uid))
		}
		return nil
	}
}

var (
	emptyPF *PIDFile
)

func Generate(path string, options ...Option) interface {
	Cleanup() error
} {
	if path == "" {
		return emptyPF
	}
	c, err := New(path, options...)
	if err != nil {
		panic(err)
	}
	err = c.Generate()
	if err != nil {
		panic(err)
	}
	return c
}
