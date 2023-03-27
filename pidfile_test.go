package pidfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPIDFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test.pid")
	pf, err := New(path)
	assert.NoError(t, err)
	assert.NotNil(t, pf)
	assert.FileExists(t, path)
	pf.Cleanup()
	assert.NoFileExists(t, path)
}

func ExamplePIDFile() {
	pf, err := New("./test.pid")
	if err != nil {
		panic(err)
	}
	defer pf.Cleanup()
	// ...

	// Output:
}
