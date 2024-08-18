package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/credentials"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

<<<<<<< HEAD
func setupConfigDir(t *testing.T) string {
	t.Helper()
	tmpdir := t.TempDir()
=======
var homeKey = "HOME"

func init() {
	if runtime.GOOS == "windows" {
		homeKey = "USERPROFILE"
	}
}

func setupConfigDir(t *testing.T) (string, func()) {
	tmpdir, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	oldDir := Dir()
	SetDir(tmpdir)

	return tmpdir, func() {
		SetDir(oldDir)
		os.RemoveAll(tmpdir)
	}
}

func TestEmptyConfigDir(t *testing.T) {
	tmpHome, cleanup := setupConfigDir(t)
	defer cleanup()

	config, err := Load("")
	assert.NilError(t, err)

	expectedConfigFilename := filepath.Join(tmpHome, ConfigFileName)
	assert.Check(t, is.Equal(expectedConfigFilename, config.Filename))

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestMissingFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestSaveFileToDirs(t *testing.T) {
<<<<<<< HEAD
	tmpHome := filepath.Join(t.TempDir(), ".docker")
=======
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	tmpHome += "/.docker"

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	config, err := Load(tmpHome)
	assert.NilError(t, err)

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

func TestEmptyFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
<<<<<<< HEAD
	err := os.WriteFile(fn, []byte(""), 0o600)
=======
	err = ioutil.WriteFile(fn, []byte(""), 0600)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(t, err)

	_, err = Load(tmpHome)
	assert.NilError(t, err)
}

func TestEmptyJSON(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
<<<<<<< HEAD
	err := os.WriteFile(fn, []byte("{}"), 0o600)
=======
	err = ioutil.WriteFile(fn, []byte("{}"), 0600)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(t, err)

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	// Now save it and make sure it shows up in new form
	saveConfigAndValidateNewFormat(t, config, tmpHome)
}

<<<<<<< HEAD
=======
func TestOldInvalidsAuth(t *testing.T) {
	invalids := map[string]string{
		`username = test`: "The Auth config file is empty",
		`username
password`: "Invalid Auth config file",
		`username = test
email`: "Invalid auth configuration file",
	}

	resetHomeDir()
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)
	defer env.Patch(t, homeKey, tmpHome)()

	for content, expectedError := range invalids {
		fn := filepath.Join(tmpHome, oldConfigfile)
		err := ioutil.WriteFile(fn, []byte(content), 0600)
		assert.NilError(t, err)

		_, err = Load(tmpHome)
		assert.ErrorContains(t, err, expectedError)
	}
}

func TestOldValidAuth(t *testing.T) {
	resetHomeDir()
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)
	defer env.Patch(t, homeKey, tmpHome)()

	fn := filepath.Join(tmpHome, oldConfigfile)
	js := `username = am9lam9lOmhlbGxv
	email = user@example.com`
	err = ioutil.WriteFile(fn, []byte(js), 0600)
	assert.NilError(t, err)

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	// defaultIndexserver is https://index.docker.io/v1/
	ac := config.AuthConfigs["https://index.docker.io/v1/"]
	assert.Equal(t, ac.Username, "joejoe")
	assert.Equal(t, ac.Password, "hello")

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}
}`

	assert.Check(t, is.Equal(expConfStr, configStr))
}

func TestOldJSONInvalid(t *testing.T) {
	resetHomeDir()
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)
	defer env.Patch(t, homeKey, tmpHome)()

	fn := filepath.Join(tmpHome, oldConfigfile)
	js := `{"https://index.docker.io/v1/":{"auth":"test","email":"user@example.com"}}`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	// Use Contains instead of == since the file name will change each time
	if err == nil || !strings.Contains(err.Error(), "Invalid auth configuration file") {
		t.Fatalf("Expected an error got : %v, %v", config, err)
	}
}

func TestOldJSON(t *testing.T) {
	resetHomeDir()
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)
	defer env.Patch(t, homeKey, tmpHome)()

	fn := filepath.Join(tmpHome, oldConfigfile)
	js := `{"https://index.docker.io/v1/":{"auth":"am9lam9lOmhlbGxv","email":"user@example.com"}}`
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	ac := config.AuthConfigs["https://index.docker.io/v1/"]
	assert.Equal(t, ac.Username, "joejoe")
	assert.Equal(t, ac.Password, "hello")

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "am9lam9lOmhlbGxv",
			"email": "user@example.com"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n'%s'\n not \n'%s'\n", configStr, expConfStr)
	}
}

func TestOldJSONFallbackDeprecationWarning(t *testing.T) {
	js := `{"https://index.docker.io/v1/":{"auth":"am9lam9lOmhlbGxv","email":"user@example.com"}}`
	tmpHome := fs.NewDir(t, t.Name(), fs.WithFile(oldConfigfile, js))
	defer tmpHome.Remove()
	defer env.PatchAll(t, map[string]string{homeKey: tmpHome.Path(), "DOCKER_CONFIG": ""})()

	// reset the homeDir, configDir, and its sync.Once, to force them being resolved again
	resetHomeDir()
	resetConfigDir()

	buffer := new(bytes.Buffer)
	configFile := LoadDefaultConfigFile(buffer)
	expected := configfile.New(tmpHome.Join(configFileDir, ConfigFileName))
	expected.AuthConfigs = map[string]types.AuthConfig{
		"https://index.docker.io/v1/": {
			Username:      "joejoe",
			Password:      "hello",
			Email:         "user@example.com",
			ServerAddress: "https://index.docker.io/v1/",
		},
	}
	assert.Assert(t, strings.Contains(buffer.String(), "WARNING: Support for the legacy ~/.dockercfg configuration file and file-format is deprecated and will be removed in an upcoming release"))
	assert.Check(t, is.DeepEqual(expected, configFile))
}

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
func TestNewJSON(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := ` { "auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv" } } }`
<<<<<<< HEAD
	if err := os.WriteFile(fn, []byte(js), 0o600); err != nil {
=======
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	ac := config.AuthConfigs["https://index.docker.io/v1/"]
	assert.Equal(t, ac.Username, "joejoe")
	assert.Equal(t, ac.Password, "hello")

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", configStr, expConfStr)
	}
}

func TestNewJSONNoEmail(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := ` { "auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv" } } }`
<<<<<<< HEAD
	if err := os.WriteFile(fn, []byte(js), 0o600); err != nil {
=======
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	ac := config.AuthConfigs["https://index.docker.io/v1/"]
	assert.Equal(t, ac.Username, "joejoe")
	assert.Equal(t, ac.Password, "hello")

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)

	expConfStr := `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	}
}`

	if configStr != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", configStr, expConfStr)
	}
}

func TestJSONWithPsFormat(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := `{
		"auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } },
		"psFormat": "table {{.ID}}\\t{{.Label \"com.docker.label.cpu\"}}"
}`
<<<<<<< HEAD
	if err := os.WriteFile(fn, []byte(js), 0o600); err != nil {
=======
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	if config.PsFormat != `table {{.ID}}\t{{.Label "com.docker.label.cpu"}}` {
		t.Fatalf("Unknown ps format: %s\n", config.PsFormat)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)
	if !strings.Contains(configStr, `"psFormat":`) ||
		!strings.Contains(configStr, "{{.ID}}") {
		t.Fatalf("Should have save in new form: %s", configStr)
	}
}

func TestJSONWithCredentialStore(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := `{
		"auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } },
		"credsStore": "crazy-secure-storage"
}`
<<<<<<< HEAD
	if err := os.WriteFile(fn, []byte(js), 0o600); err != nil {
=======
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	if config.CredentialsStore != "crazy-secure-storage" {
		t.Fatalf("Unknown credential store: %s\n", config.CredentialsStore)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)
	if !strings.Contains(configStr, `"credsStore":`) ||
		!strings.Contains(configStr, "crazy-secure-storage") {
		t.Fatalf("Should have save in new form: %s", configStr)
	}
}

func TestJSONWithCredentialHelpers(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	js := `{
		"auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } },
		"credHelpers": { "images.io": "images-io", "containers.com": "crazy-secure-storage" }
}`
<<<<<<< HEAD
	if err := os.WriteFile(fn, []byte(js), 0o600); err != nil {
=======
	if err := ioutil.WriteFile(fn, []byte(js), 0600); err != nil {
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		t.Fatal(err)
	}

	config, err := Load(tmpHome)
	assert.NilError(t, err)

	if config.CredentialHelpers == nil {
		t.Fatal("config.CredentialHelpers was nil")
	} else if config.CredentialHelpers["images.io"] != "images-io" ||
		config.CredentialHelpers["containers.com"] != "crazy-secure-storage" {
		t.Fatalf("Credential helpers not deserialized properly: %v\n", config.CredentialHelpers)
	}

	// Now save it and make sure it shows up in new form
	configStr := saveConfigAndValidateNewFormat(t, config, tmpHome)
	if !strings.Contains(configStr, `"credHelpers":`) ||
		!strings.Contains(configStr, "images.io") ||
		!strings.Contains(configStr, "images-io") ||
		!strings.Contains(configStr, "containers.com") ||
		!strings.Contains(configStr, "crazy-secure-storage") {
		t.Fatalf("Should have save in new form: %s", configStr)
	}
}

// Save it and make sure it shows up in new form
func saveConfigAndValidateNewFormat(t *testing.T, config *configfile.ConfigFile, configDir string) string {
	t.Helper()
	assert.NilError(t, config.Save())

	buf, err := ioutil.ReadFile(filepath.Join(configDir, ConfigFileName))
	assert.NilError(t, err)
	assert.Check(t, is.Contains(string(buf), `"auths":`))
	return string(buf)
}

func TestConfigDir(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	if Dir() == tmpHome {
		t.Fatalf("Expected ConfigDir to be different than %s by default, but was the same", tmpHome)
	}

	// Update configDir
	SetDir(tmpHome)

	if Dir() != tmpHome {
		t.Fatalf("Expected ConfigDir to %s, but was %s", tmpHome, Dir())
	}
}

func TestJSONReaderNoFile(t *testing.T) {
	js := ` { "auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } } }`

	config, err := LoadFromReader(strings.NewReader(js))
	assert.NilError(t, err)

	ac := config.AuthConfigs["https://index.docker.io/v1/"]
	assert.Equal(t, ac.Username, "joejoe")
	assert.Equal(t, ac.Password, "hello")
}

func TestJSONWithPsFormatNoFile(t *testing.T) {
	js := `{
		"auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv", "email": "user@example.com" } },
		"psFormat": "table {{.ID}}\\t{{.Label \"com.docker.label.cpu\"}}"
}`
	config, err := LoadFromReader(strings.NewReader(js))
	assert.NilError(t, err)

	if config.PsFormat != `table {{.ID}}\t{{.Label "com.docker.label.cpu"}}` {
		t.Fatalf("Unknown ps format: %s\n", config.PsFormat)
	}
}

func TestJSONSaveWithNoFile(t *testing.T) {
	js := `{
		"auths": { "https://index.docker.io/v1/": { "auth": "am9lam9lOmhlbGxv" } },
		"psFormat": "table {{.ID}}\\t{{.Label \"com.docker.label.cpu\"}}"
}`
	config, err := LoadFromReader(strings.NewReader(js))
	assert.NilError(t, err)
	err = config.Save()
	assert.ErrorContains(t, err, "with empty filename")

	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	defer f.Close()

	assert.NilError(t, config.SaveToWriter(f))
	buf, err := ioutil.ReadFile(filepath.Join(tmpHome, ConfigFileName))
	assert.NilError(t, err)
	expConfStr := `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "am9lam9lOmhlbGxv"
		}
	},
	"psFormat": "table {{.ID}}\\t{{.Label \"com.docker.label.cpu\"}}"
}`
	if string(buf) != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\nnot \n%s", string(buf), expConfStr)
	}
}

<<<<<<< HEAD
=======
func TestLegacyJSONSaveWithNoFile(t *testing.T) {
	js := `{"https://index.docker.io/v1/":{"auth":"am9lam9lOmhlbGxv","email":"user@example.com"}}`
	config, err := LegacyLoadFromReader(strings.NewReader(js))
	assert.NilError(t, err)
	err = config.Save()
	assert.ErrorContains(t, err, "with empty filename")

	tmpHome, err := ioutil.TempDir("", "config-test")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpHome)

	fn := filepath.Join(tmpHome, ConfigFileName)
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()

	assert.NilError(t, config.SaveToWriter(f))
	buf, err := ioutil.ReadFile(filepath.Join(tmpHome, ConfigFileName))
	assert.NilError(t, err)

	expConfStr := `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "am9lam9lOmhlbGxv",
			"email": "user@example.com"
		}
	}
}`

	if string(buf) != expConfStr {
		t.Fatalf("Should have save in new form: \n%s\n not \n%s", string(buf), expConfStr)
	}
}

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
func TestLoadDefaultConfigFile(t *testing.T) {
	dir, cleanup := setupConfigDir(t)
	defer cleanup()
	buffer := new(bytes.Buffer)

	filename := filepath.Join(dir, ConfigFileName)
	content := []byte(`{"PsFormat": "format"}`)
<<<<<<< HEAD
	err := os.WriteFile(filename, content, 0o644)
=======
	err := ioutil.WriteFile(filename, content, 0644)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(t, err)

	configFile := LoadDefaultConfigFile(buffer)
	credStore := credentials.DetectDefaultStore("")
	expected := configfile.New(filename)
	expected.CredentialsStore = credStore
	expected.PsFormat = "format"

	assert.Check(t, is.DeepEqual(expected, configFile))
}

func TestConfigPath(t *testing.T) {
	oldDir := Dir()

	for _, tc := range []struct {
		name        string
		dir         string
		path        []string
		expected    string
		expectedErr string
	}{
		{
			name:     "valid_path",
			dir:      "dummy",
			path:     []string{"a", "b"},
			expected: filepath.Join("dummy", "a", "b"),
		},
		{
			name:     "valid_path_absolute_dir",
			dir:      "/dummy",
			path:     []string{"a", "b"},
			expected: filepath.Join("/dummy", "a", "b"),
		},
		{
			name:        "invalid_relative_path",
			dir:         "dummy",
			path:        []string{"e", "..", "..", "f"},
			expectedErr: fmt.Sprintf("is outside of root config directory %q", "dummy"),
		},
		{
			name:        "invalid_absolute_path",
			dir:         "dummy",
			path:        []string{"/a", "..", ".."},
			expectedErr: fmt.Sprintf("is outside of root config directory %q", "dummy"),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			SetDir(tc.dir)
			f, err := Path(tc.path...)
			assert.Equal(t, f, tc.expected)
			if tc.expectedErr == "" {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}

	SetDir(oldDir)
}

// TestSetDir verifies that Dir() does not overwrite the value set through
// SetDir() if it has not been run before.
func TestSetDir(t *testing.T) {
	const expected = "my_config_dir"
	resetConfigDir()
	SetDir(expected)
	assert.Check(t, is.Equal(Dir(), expected))
}
