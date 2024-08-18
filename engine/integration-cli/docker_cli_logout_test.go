package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/testutil"
	"gotest.tools/v3/assert"
)

func (s *DockerRegistryAuthHtpasswdSuite) TestLogoutWithExternalAuth(c *testing.T) {
	ctx := testutil.GetContext(c)
	s.d.StartWithBusybox(ctx, c)

	workingDir, err := os.Getwd()
	assert.NilError(c, err)
	absolute, err := filepath.Abs(filepath.Join(workingDir, "fixtures", "auth"))
	assert.NilError(c, err)

	osPath := os.Getenv("PATH")
	testPath := fmt.Sprintf("%s%c%s", osPath, filepath.ListSeparator, absolute)
	c.Setenv("PATH", testPath)

	imgRepoName := fmt.Sprintf("%v/dockercli/busybox:authtest", privateRegistryURL)

	tmp, err := ioutil.TempDir("", "integration-cli-")
	assert.NilError(c, err)
	defer os.RemoveAll(tmp)

	externalAuthConfig := `{ "credsStore": "shell-test" }`

	configPath := filepath.Join(tmp, "config.json")
<<<<<<< HEAD
	err = os.WriteFile(configPath, []byte(externalAuthConfig), 0o644)
=======
	err = ioutil.WriteFile(configPath, []byte(externalAuthConfig), 0644)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(c, err)

	_, err = s.d.Cmd("--config", tmp, "login", "-u", s.reg.Username(), "-p", s.reg.Password(), privateRegistryURL)
	assert.NilError(c, err)

	b, err := ioutil.ReadFile(configPath)
	assert.NilError(c, err)
	assert.Assert(c, !strings.Contains(string(b), `"auth":`))
	assert.Assert(c, strings.Contains(string(b), privateRegistryURL))

	_, err = s.d.Cmd("--config", tmp, "tag", "busybox", imgRepoName)
	assert.NilError(c, err)
	_, err = s.d.Cmd("--config", tmp, "push", imgRepoName)
	assert.NilError(c, err)
	_, err = s.d.Cmd("--config", tmp, "logout", privateRegistryURL)
	assert.NilError(c, err)

	b, err = ioutil.ReadFile(configPath)
	assert.NilError(c, err)
	assert.Assert(c, !strings.Contains(string(b), privateRegistryURL))

	// check I cannot pull anymore
	out, err := s.d.Cmd("--config", tmp, "pull", imgRepoName)
	assert.ErrorContains(c, err, "", out)
	assert.Assert(c, strings.Contains(out, "no basic auth credentials"))
}

// #23100
func (s *DockerRegistryAuthHtpasswdSuite) TestLogoutWithWrongHostnamesStored(c *testing.T) {
	workingDir, err := os.Getwd()
	assert.NilError(c, err)
	absolute, err := filepath.Abs(filepath.Join(workingDir, "fixtures", "auth"))
	assert.NilError(c, err)

	osPath := os.Getenv("PATH")
	testPath := fmt.Sprintf("%s%c%s", osPath, filepath.ListSeparator, absolute)
	c.Setenv("PATH", testPath)

	cmd := exec.Command("docker-credential-shell-test", "store")
	stdin := bytes.NewReader([]byte(fmt.Sprintf(`{"ServerURL": "https://%s", "Username": "%s", "Secret": "%s"}`, privateRegistryURL, s.reg.Username(), s.reg.Password())))
	cmd.Stdin = stdin
	assert.NilError(c, cmd.Run())

	tmp, err := ioutil.TempDir("", "integration-cli-")
	assert.NilError(c, err)

	externalAuthConfig := fmt.Sprintf(`{ "auths": {"https://%s": {}}, "credsStore": "shell-test" }`, privateRegistryURL)

	configPath := filepath.Join(tmp, "config.json")
<<<<<<< HEAD
	err = os.WriteFile(configPath, []byte(externalAuthConfig), 0o644)
=======
	err = ioutil.WriteFile(configPath, []byte(externalAuthConfig), 0644)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(c, err)

	cli.DockerCmd(c, "--config", tmp, "login", "-u", s.reg.Username(), "-p", s.reg.Password(), privateRegistryURL)

	b, err := ioutil.ReadFile(configPath)
	assert.NilError(c, err)
	assert.Assert(c, strings.Contains(string(b), fmt.Sprintf(`"https://%s": {}`, privateRegistryURL)))
	assert.Assert(c, strings.Contains(string(b), fmt.Sprintf(`"%s": {}`, privateRegistryURL)))

	cli.DockerCmd(c, "--config", tmp, "logout", privateRegistryURL)

	b, err = ioutil.ReadFile(configPath)
	assert.NilError(c, err)
	assert.Assert(c, !strings.Contains(string(b), fmt.Sprintf(`"https://%s": {}`, privateRegistryURL)))
	assert.Assert(c, !strings.Contains(string(b), fmt.Sprintf(`"%s": {}`, privateRegistryURL)))
}
