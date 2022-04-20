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

type Option func()

type context struct {
	stop func() error
}

func (c *context) Stop() error {
	return c.stop()
}

const (
	ROOT = 0
)

func Permission(uid int) Option {
	return func() {
		if os.Getuid() != uid {
			panic(errors.New("expect uid:" + strconv.Itoa(uid)))
		}
	}
}

func Start(path string, options ...Option) interface {
	Stop() error
} {
	if !atomic.CompareAndSwapInt32(&inited, 0, 1) {
		panic(errors.New("duplicated"))
	}
	for _, option := range options {
		option()
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	c := new(context)
	c.stop = func() error {
		defer atomic.StoreInt32(&inited, 0)
		return os.Remove(path)
	}
	f.WriteString(strconv.Itoa(os.Getpid()))
	return c
}
