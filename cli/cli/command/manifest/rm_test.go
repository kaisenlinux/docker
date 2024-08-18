package manifest

import (
	"io/ioutil"
	"testing"

	"github.com/docker/cli/internal/test"
	"gotest.tools/v3/assert"
)

// create two manifest lists and remove them both
func TestRmSeveralManifests(t *testing.T) {
<<<<<<< HEAD
	manifestStore := store.NewStore(t.TempDir())
=======
	store, cleanup := newTempManifestStore(t)
	defer cleanup()
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)

	cli := test.NewFakeCli(nil)
	cli.SetManifestStore(manifestStore)

	list1 := ref(t, "first:1")
	namedRef := ref(t, "alpine:3.0")
	err := manifestStore.Save(list1, namedRef, fullImageManifest(t, namedRef))
	assert.NilError(t, err)
	namedRef = ref(t, "alpine:3.1")
	err = manifestStore.Save(list1, namedRef, fullImageManifest(t, namedRef))
	assert.NilError(t, err)

	list2 := ref(t, "second:2")
	namedRef = ref(t, "alpine:3.2")
	err = manifestStore.Save(list2, namedRef, fullImageManifest(t, namedRef))
	assert.NilError(t, err)

	cmd := newRmManifestListCommand(cli)
	cmd.SetArgs([]string{"example.com/first:1", "example.com/second:2"})
	cmd.SetOut(ioutil.Discard)
	err = cmd.Execute()
	assert.NilError(t, err)

	_, search1 := cli.ManifestStore().GetList(list1)
	_, search2 := cli.ManifestStore().GetList(list2)
	assert.Error(t, search1, "No such manifest: example.com/first:1")
	assert.Error(t, search2, "No such manifest: example.com/second:2")
}

// attempt to remove a manifest list which was never created
func TestRmManifestNotCreated(t *testing.T) {
<<<<<<< HEAD
	manifestStore := store.NewStore(t.TempDir())
=======
	store, cleanup := newTempManifestStore(t)
	defer cleanup()
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)

	cli := test.NewFakeCli(nil)
	cli.SetManifestStore(manifestStore)

	list2 := ref(t, "second:2")
	namedRef := ref(t, "alpine:3.2")
	err := manifestStore.Save(list2, namedRef, fullImageManifest(t, namedRef))
	assert.NilError(t, err)

	cmd := newRmManifestListCommand(cli)
	cmd.SetArgs([]string{"example.com/first:1", "example.com/second:2"})
	cmd.SetOut(ioutil.Discard)
	err = cmd.Execute()
	assert.Error(t, err, "No such manifest: example.com/first:1")

	_, err = cli.ManifestStore().GetList(list2)
	assert.Error(t, err, "No such manifest: example.com/second:2")
}
