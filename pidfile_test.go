package pidfile

import (
	"errors"
	"io"
	"os"
	"strconv"
	"testing"
)

func TestPIDFile(t *testing.T) {
	// path := filepath.Join(os.TempDir(), "pidfile")
	path := "./pidfile"
	c, err := New(path)
	if err != nil {
		t.Error(err)
	}
	err = c.Generate()
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	if string(bs) != strconv.Itoa(os.Getpid()) {
		t.Error(string(bs) + "!=" + strconv.Itoa(os.Getpid()))
	}
	f.Close()
	err = c.Cleanup()
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error(os.ErrExist)
	}
}

func ExamplePIDFile() {
	pf, err := New("/var/run/example/pidfile")
	if err != nil {
		panic(err)
	}
	if err := pf.Generate(); err != nil {
		panic(err)
	}
	defer pf.Cleanup()
}

func TestGenerate(t *testing.T) {
	path := "./pidfile"
	c := Generate(path)
	f, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	if string(bs) != strconv.Itoa(os.Getpid()) {
		t.Error(string(bs) + "!=" + strconv.Itoa(os.Getpid()))
	}
	f.Close()
	err = c.Cleanup()
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error(os.ErrExist)
	}
}

func ExampleGenerate() {
	defer Generate("/var/run/example/pidfile").Cleanup()
}
