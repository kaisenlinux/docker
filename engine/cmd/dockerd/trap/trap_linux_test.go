//go:build linux

package trap // import "github.com/docker/docker/cmd/dockerd/trap"

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func buildTestBinary(t *testing.T, tmpdir string, prefix string) (string, string) {
	t.Helper()
	tmpDir, err := ioutil.TempDir(tmpdir, prefix)
	assert.NilError(t, err)
	exePath := tmpDir + "/" + prefix
	wd, _ := os.Getwd()
	testHelperCode := wd + "/testfiles/main.go"
	cmd := exec.Command("go", "build", "-o", exePath, testHelperCode)
	err = cmd.Run()
	assert.NilError(t, err)
	return exePath, tmpDir
}

func TestTrap(t *testing.T) {
	sigmap := []struct {
		name     string
		signal   os.Signal
		multiple bool
	}{
		{"TERM", syscall.SIGTERM, false},
		{"INT", os.Interrupt, false},
		{"TERM", syscall.SIGTERM, true},
		{"INT", os.Interrupt, true},
	}
	exePath, tmpDir := buildTestBinary(t, "", "main")
	defer os.RemoveAll(tmpDir)

	for _, v := range sigmap {
		t.Run(v.name, func(t *testing.T) {
			cmd := exec.Command(exePath)
			cmd.Env = append(os.Environ(), "SIGNAL_TYPE="+v.name)
			if v.multiple {
				cmd.Env = append(cmd.Env, "IF_MULTIPLE=1")
			}
			err := cmd.Start()
			assert.NilError(t, err)
			err = cmd.Wait()
			e, ok := err.(*exec.ExitError)
			assert.Assert(t, ok, "expected exec.ExitError, got %T", e)

			code := e.Sys().(syscall.WaitStatus).ExitStatus()
			if v.multiple {
				assert.Check(t, is.DeepEqual(128+int(v.signal.(syscall.Signal)), code))
			} else {
				assert.Check(t, is.Equal(99, code))
			}
		})
	}
<<<<<<< HEAD:engine/cmd/dockerd/trap/trap_linux_test.go
=======

}

func TestDumpStacks(t *testing.T) {
	directory, err := ioutil.TempDir("", "test-dump-tasks")
	assert.Check(t, err)
	defer os.RemoveAll(directory)
	dumpPath, err := DumpStacks(directory)
	assert.Check(t, err)
	readFile, _ := ioutil.ReadFile(dumpPath)
	fileData := string(readFile)
	assert.Check(t, is.Contains(fileData, "goroutine"))
}

func TestDumpStacksWithEmptyInput(t *testing.T) {
	path, err := DumpStacks("")
	assert.Check(t, err)
	assert.Check(t, is.Equal(os.Stderr.Name(), path))
>>>>>>> parent of ea55db5 (Import the 20.10.24 version):engine/pkg/signal/trap_linux_test.go
}
