// Package pidfile provides structure and helper functions to create and remove
// PID file. A PID file is usually a file used to store the process ID of a
// running process.
package pidfile // import "github.com/docker/docker/pkg/pidfile"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/docker/docker/pkg/process"
)

// Read reads the "PID file" at path, and returns the PID if it contains a
// valid PID of a running process, or 0 otherwise. It returns an error when
// failing to read the file, or if the file doesn't exist, but malformed content
// is ignored. Consumers should therefore check if the returned PID is a non-zero
// value before use.
func Read(path string) (pid int, err error) {
	pidByte, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	pid, err = strconv.Atoi(string(bytes.TrimSpace(pidByte)))
	if err != nil {
		return 0, nil
	}
	if pid != 0 && process.Alive(pid) {
		return pid, nil
	}
	return 0, nil
}

<<<<<<< HEAD
// Write writes a "PID file" at the specified path. It returns an error if the
// file exists and contains a valid PID of a running process, or when failing
// to write the file.
func Write(path string, pid int) error {
	if pid < 1 {
		// We might be running as PID 1 when running docker-in-docker,
		// but 0 or negative PIDs are not acceptable.
		return fmt.Errorf("invalid PID (%d): only positive PIDs are allowed", pid)
	}
	oldPID, err := Read(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if oldPID != 0 {
		return fmt.Errorf("process with PID %d is still running", oldPID)
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o644)
=======
func checkPIDFileAlreadyExists(path string) error {
	if pidByte, err := ioutil.ReadFile(path); err == nil {
		pidString := strings.TrimSpace(string(pidByte))
		if pid, err := strconv.Atoi(pidString); err == nil {
			if processExists(pid) {
				return fmt.Errorf("pid file found, ensure docker is not running or delete %s", path)
			}
		}
	}
	return nil
}

// New creates a PIDfile using the specified path.
func New(path string) (*PIDFile, error) {
	if err := checkPIDFileAlreadyExists(path); err != nil {
		return nil, err
	}
	// Note MkdirAll returns nil if a directory already exists
	if err := system.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return nil, err
	}

	return &PIDFile{path: path}, nil
}

// Remove removes the PIDFile.
func (file PIDFile) Remove() error {
	return os.Remove(file.path)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}
