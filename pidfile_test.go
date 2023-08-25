package pidfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	var (
		path      = filepath.Join(os.TempDir(), "test.pid")
		cleanable Cleanable
	)
	assert.NoFileExists(t, path)
	assert.NotPanics(t, func() { cleanable = Generate(path) })
	assert.FileExists(t, path)
	assert.NotNil(t, cleanable)
	assert.NotPanics(t, func() { cleanable.Cleanup() })
	assert.NoFileExists(t, path)
}
func TestGenerateEmptyPath(t *testing.T) {
	var (
		path      = ""
		cleanable Cleanable
	)
	assert.NoFileExists(t, path)
	assert.NotPanics(t, func() { cleanable = Generate(path) })
	assert.NoFileExists(t, path)
	assert.NotNil(t, cleanable)
	assert.NotPanics(t, func() { cleanable.Cleanup() })
	assert.NoFileExists(t, path)
}

func ExamplePidfile() {
	defer Generate("./test.pid").Cleanup()
	// ...
}
