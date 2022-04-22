package pidfile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestPIDFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "pidfile")
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
func TestPIDFileExist(t *testing.T) {
	path := filepath.Join(os.TempDir(), "pidfile")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	f.WriteString("0")
	f.Close()
	defer os.Remove(path)
	c, err := New(path)
	if err != nil {
		t.Error(err)
	}
	err = c.Generate()
	if !errors.Is(err, ErrFileExists) {
		t.Error("overwrite")
	}
	f, err = os.Open(path)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	if string(bs) != "0" {
		t.Error("overwrited " + string(bs))
	}
}
func TestPIDFileDuplicated(t *testing.T) {
	path := filepath.Join(os.TempDir(), "pidfile")
	c, err := New(path)
	if err != nil {
		t.Error(err)
	}
	err = c.Generate()
	if err != nil {
		t.Error(err)
	}
	defer c.Cleanup()
	err = c.Generate()
	if !errors.Is(err, ErrDuplicated) {
		t.Error("duplicated")
	}
}

func ExamplePIDFile() {
	pf, err := New("./pidfile")
	if err != nil {
		panic(err)
	}
	if err := pf.Generate(); err != nil {
		panic(err)
	}
	defer pf.Cleanup()
	// ...

	// Output:
}

func TestGenerate(t *testing.T) {
	path := filepath.Join(os.TempDir(), "pidfile")
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
	defer Generate("./pidfile").Cleanup()
	// ...

	// Output:
}
