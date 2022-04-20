package pidfile

import (
	"errors"
	"io"
	"os"
	"strconv"
	"testing"
)

func Example() {
	defer Start("/var/run/example/pidfile").Stop()
}

func TestStart(t *testing.T) {
	// path := filepath.Join(os.TempDir(), "pidfile")
	path := "./pidfile"
	c := Start(path)
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
	err = c.Stop()
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error(os.ErrExist)
	}
}
