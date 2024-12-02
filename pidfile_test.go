package pidfile

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPidfile(t *testing.T) {
	var pf interface{ Cleanup() }
	path := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10)+".pid")
	assert.NoFileExists(t, path)
	assert.NotPanics(t, func() { pf = Generate(Path(path)) })
	assert.FileExists(t, path)
	assert.NotNil(t, pf)
	bs, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(os.Getpid()), string(bs))
	assert.NotPanics(t, func() { pf.Cleanup() })
	assert.NoFileExists(t, path)
}

func ExampleGenerate() {
	defer Generate().Cleanup()
	// ...
}
