package pidfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPIDFile(t *testing.T) {
	var (
		path = filepath.Join(os.TempDir(), "test.pid")
		pf1  = New(path)
		pf2  = New(path)
	)
	assert.NotNil(t, pf1)
	assert.NotNil(t, pf2)
	assert.NotPanics(t, func() { pf1.Generate() })
	assert.Panics(t, func() { pf2.Generate() })
	assert.FileExists(t, path)
	assert.NoError(t, pf1.Cleanup())
	assert.Error(t, pf2.Cleanup())
	assert.NoFileExists(t, path)
}

func ExamplePIDFile_1() {
	defer Generate("./test.pid").Cleanup()
	// ...

	// Output:
}
func ExamplePIDFile_2() {
	var (
		pf = New("./test.pid")
	)
	defer pf.Generate().Cleanup()
	// ...

	// Output:
}
