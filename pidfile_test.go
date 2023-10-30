package pidfile

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	var pf Cleanable
	path := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10)+".pid")
	assert.NoFileExists(t, path)
	assert.NotPanics(t, func() { pf = Generate(path) })
	assert.FileExists(t, path)
	assert.NotNil(t, pf)
	bs, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(os.Getpid()), string(bs))
	assert.NotPanics(t, func() { pf.Cleanup() })
	assert.NoFileExists(t, path)
}
func TestGenerateEmpty(t *testing.T) {
	var pf Cleanable
	path := ""
	assert.NotPanics(t, func() { pf = Generate(path) })
	assert.NotNil(t, pf)
	assert.NotPanics(t, func() { pf.Cleanup() })
}
func TestGenerateNotEmpty(t *testing.T) {
	var pf Cleanable
	path := ""
	assert.Panics(t, func() { pf = Generate(path, NotEmpty) })
	assert.Nil(t, pf)
}

func ExampleGenerate() {
	defer Generate("example.pid").Cleanup()
	// ...
}
