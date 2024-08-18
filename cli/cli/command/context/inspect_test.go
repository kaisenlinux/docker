package context

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestInspect(t *testing.T) {
<<<<<<< HEAD
	cli := makeFakeCli(t)
	createTestContext(t, cli, "current", map[string]any{
		"MyCustomMetadata": "MyCustomMetadataValue",
	})
=======
	cli, cleanup := makeFakeCli(t)
	defer cleanup()
	createTestContextWithKubeAndSwarm(t, cli, "current", "all")
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	cli.OutBuffer().Reset()
	assert.NilError(t, runInspect(cli, inspectOptions{
		refs: []string{"current"},
	}))
	expected := string(golden.Get(t, "inspect.golden"))
	si := cli.ContextStore().GetStorageInfo("current")
	expected = strings.Replace(expected, "<METADATA_PATH>", strings.ReplaceAll(si.MetadataPath, `\`, `\\`), 1)
	expected = strings.Replace(expected, "<TLS_PATH>", strings.ReplaceAll(si.TLSPath, `\`, `\\`), 1)
	assert.Equal(t, cli.OutBuffer().String(), expected)
}
