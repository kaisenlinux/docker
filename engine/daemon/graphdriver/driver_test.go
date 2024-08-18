package graphdriver // import "github.com/docker/docker/daemon/graphdriver"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsEmptyDir(t *testing.T) {
	tmp, err := ioutil.TempDir("", "test-is-empty-dir")
	assert.NilError(t, err)
	defer os.RemoveAll(tmp)

	d := filepath.Join(tmp, "empty-dir")
	err = os.Mkdir(d, 0o755)
	assert.NilError(t, err)
	empty := isEmptyDir(d)
	assert.Check(t, empty)

	d = filepath.Join(tmp, "dir-with-subdir")
	err = os.MkdirAll(filepath.Join(d, "subdir"), 0o755)
	assert.NilError(t, err)
	empty = isEmptyDir(d)
	assert.Check(t, !empty)

	d = filepath.Join(tmp, "dir-with-empty-file")
	err = os.Mkdir(d, 0o755)
	assert.NilError(t, err)
<<<<<<< HEAD
	f, err := os.CreateTemp(d, "file")
=======
	_, err = ioutil.TempFile(d, "file")
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(t, err)
	defer f.Close()
	empty = isEmptyDir(d)
	assert.Check(t, !empty)
}
