package system // import "github.com/docker/docker/pkg/system"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

<<<<<<< HEAD
// TestChtimesModTime tests Chtimes on a tempfile. Test only mTime, because
// aTime is OS dependent.
func TestChtimesModTime(t *testing.T) {
	file := filepath.Join(t.TempDir(), "exist")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	beforeUnixEpochTime := unixEpochTime.Add(-100 * time.Second)
	afterUnixEpochTime := unixEpochTime.Add(100 * time.Second)
=======
// prepareTempFile creates a temporary file in a temporary directory.
func prepareTempFile(t *testing.T) (string, string) {
	dir, err := ioutil.TempDir("", "docker-system-test")
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(dir, "exist")
	if err := ioutil.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	return file, dir
}

// TestChtimes tests Chtimes on a tempfile. Test only mTime, because aTime is OS dependent
func TestChtimes(t *testing.T) {
	file, dir := prepareTempFile(t)
	defer os.RemoveAll(dir)

	beforeUnixEpochTime := time.Unix(0, 0).Add(-100 * time.Second)
	unixEpochTime := time.Unix(0, 0)
	afterUnixEpochTime := time.Unix(100, 0)
	unixMaxTime := maxTime
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)

	// Test both aTime and mTime set to Unix Epoch
	t.Run("both aTime and mTime set to Unix Epoch", func(t *testing.T) {
		if err := Chtimes(file, unixEpochTime, unixEpochTime); err != nil {
			t.Error(err)
		}

		f, err := os.Stat(file)
		if err != nil {
			t.Fatal(err)
		}

		if f.ModTime() != unixEpochTime {
			t.Fatalf("Expected: %s, got: %s", unixEpochTime, f.ModTime())
		}
	})

	// Test aTime before Unix Epoch and mTime set to Unix Epoch
	t.Run("aTime before Unix Epoch and mTime set to Unix Epoch", func(t *testing.T) {
		if err := Chtimes(file, beforeUnixEpochTime, unixEpochTime); err != nil {
			t.Error(err)
		}

		f, err := os.Stat(file)
		if err != nil {
			t.Fatal(err)
		}

		if f.ModTime() != unixEpochTime {
			t.Fatalf("Expected: %s, got: %s", unixEpochTime, f.ModTime())
		}
	})

	// Test aTime set to Unix Epoch and mTime before Unix Epoch
	t.Run("aTime set to Unix Epoch and mTime before Unix Epoch", func(t *testing.T) {
		if err := Chtimes(file, unixEpochTime, beforeUnixEpochTime); err != nil {
			t.Error(err)
		}

		f, err := os.Stat(file)
		if err != nil {
			t.Fatal(err)
		}

		if f.ModTime() != unixEpochTime {
			t.Fatalf("Expected: %s, got: %s", unixEpochTime, f.ModTime())
		}
	})

	// Test both aTime and mTime set to after Unix Epoch (valid time)
	t.Run("both aTime and mTime set to after Unix Epoch (valid time)", func(t *testing.T) {
		if err := Chtimes(file, afterUnixEpochTime, afterUnixEpochTime); err != nil {
			t.Error(err)
		}

		f, err := os.Stat(file)
		if err != nil {
			t.Fatal(err)
		}

		if f.ModTime() != afterUnixEpochTime {
			t.Fatalf("Expected: %s, got: %s", afterUnixEpochTime, f.ModTime())
		}
	})

	// Test both aTime and mTime set to Unix max time
	t.Run("both aTime and mTime set to Unix max time", func(t *testing.T) {
		if err := Chtimes(file, unixMaxTime, unixMaxTime); err != nil {
			t.Error(err)
		}

		f, err := os.Stat(file)
		if err != nil {
			t.Fatal(err)
		}

		if f.ModTime().Truncate(time.Second) != unixMaxTime.Truncate(time.Second) {
			t.Fatalf("Expected: %s, got: %s", unixMaxTime.Truncate(time.Second), f.ModTime().Truncate(time.Second))
		}
	})
}
