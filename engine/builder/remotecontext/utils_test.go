package remotecontext // import "github.com/docker/docker/builder/remotecontext"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// createTestTempDir creates a temporary directory for testing.
// It returns the created path and a cleanup function which is meant to be used as deferred call.
// When an error occurs, it terminates the test.
func createTestTempDir(t *testing.T, dir, prefix string) (string, func()) {
<<<<<<< HEAD
	path, err := os.MkdirTemp(dir, prefix)
=======
	path, err := ioutil.TempDir(dir, prefix)

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	if err != nil {
		t.Fatalf("Error when creating directory %s with prefix %s: %s", dir, prefix, err)
	}

	return path, func() {
		err = os.RemoveAll(path)

		if err != nil {
			t.Fatalf("Error when removing directory %s: %s", path, err)
		}
	}
}

// createTestTempSubdir creates a temporary directory for testing.
// It returns the created path but doesn't provide a cleanup function,
// so createTestTempSubdir should be used only for creating temporary subdirectories
// whose parent directories are properly cleaned up.
// When an error occurs, it terminates the test.
func createTestTempSubdir(t *testing.T, dir, prefix string) string {
<<<<<<< HEAD
	path, err := os.MkdirTemp(dir, prefix)
=======
	path, err := ioutil.TempDir(dir, prefix)

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	if err != nil {
		t.Fatalf("Error when creating directory %s with prefix %s: %s", dir, prefix, err)
	}

	return path
}

// createTestTempFile creates a temporary file within dir with specific contents and permissions.
// When an error occurs, it terminates the test
func createTestTempFile(t *testing.T, dir, filename, contents string, perm os.FileMode) string {
	filePath := filepath.Join(dir, filename)
<<<<<<< HEAD
	err := os.WriteFile(filePath, []byte(contents), perm)
=======
	err := ioutil.WriteFile(filePath, []byte(contents), perm)

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	if err != nil {
		t.Fatalf("Error when creating %s file: %s", filename, err)
	}

	return filePath
}
