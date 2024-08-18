// FIXME(thaJeztah): remove once we are a module; the go:build directive prevents go from downgrading language version to go1.16:
//go:build go1.19

package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/errdefs"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func testMetadata(name string) Metadata {
	return Metadata{
		Endpoints: map[string]any{
			"ep1": endpoint{Foo: "bar"},
		},
		Metadata: context{Bar: "baz"},
		Name:     name,
	}
}

func TestMetadataGetNotExisting(t *testing.T) {
<<<<<<< HEAD
	testee := metadataStore{root: t.TempDir(), config: testCfg}
	_, err := testee.get("noexist")
	assert.ErrorType(t, err, errdefs.IsNotFound)
=======
	testDir, err := ioutil.TempDir("", t.Name())
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := metadataStore{root: testDir, config: testCfg}
	_, err = testee.get("noexist")
	assert.Assert(t, IsErrContextDoesNotExist(err))
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}

func TestMetadataCreateGetRemove(t *testing.T) {
	testDir, err := ioutil.TempDir("", t.Name())
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := metadataStore{root: testDir, config: testCfg}
	expected2 := Metadata{
		Endpoints: map[string]any{
			"ep1": endpoint{Foo: "baz"},
			"ep2": endpoint{Foo: "bee"},
		},
		Metadata: context{Bar: "foo"},
		Name:     "test-context",
	}
	testMeta := testMetadata("test-context")
	err = testee.createOrUpdate(testMeta)
	assert.NilError(t, err)
	// create a new instance to check it does not depend on some sort of state
	testee = metadataStore{root: testDir, config: testCfg}
	meta, err := testee.get("test-context")
	assert.NilError(t, err)
	assert.DeepEqual(t, meta, testMeta)

	// update

	err = testee.createOrUpdate(expected2)
	assert.NilError(t, err)
	meta, err = testee.get("test-context")
	assert.NilError(t, err)
	assert.DeepEqual(t, meta, expected2)

	assert.NilError(t, testee.remove("test-context"))
	assert.NilError(t, testee.remove("test-context")) // support duplicate remove
	_, err = testee.get("test-context")
	assert.ErrorType(t, err, errdefs.IsNotFound)
}

func TestMetadataRespectJsonAnnotation(t *testing.T) {
	testDir, err := ioutil.TempDir("", t.Name())
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := metadataStore{root: testDir, config: testCfg}
	assert.NilError(t, testee.createOrUpdate(testMetadata("test")))
	bytes, err := ioutil.ReadFile(filepath.Join(testDir, string(contextdirOf("test")), "meta.json"))
	assert.NilError(t, err)
	assert.Assert(t, cmp.Contains(string(bytes), "a_very_recognizable_field_name"))
	assert.Assert(t, cmp.Contains(string(bytes), "another_very_recognizable_field_name"))
}

func TestMetadataList(t *testing.T) {
	testDir, err := ioutil.TempDir("", t.Name())
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := metadataStore{root: testDir, config: testCfg}
	wholeData := []Metadata{
		testMetadata("context1"),
		testMetadata("context2"),
		testMetadata("context3"),
	}

	for _, s := range wholeData {
		err = testee.createOrUpdate(s)
		assert.NilError(t, err)
	}

	data, err := testee.list()
	assert.NilError(t, err)
	assert.DeepEqual(t, data, wholeData)
}

func TestEmptyConfig(t *testing.T) {
	testDir, err := ioutil.TempDir("", t.Name())
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := metadataStore{root: testDir}
	wholeData := []Metadata{
		testMetadata("context1"),
		testMetadata("context2"),
		testMetadata("context3"),
	}

	for _, s := range wholeData {
		err = testee.createOrUpdate(s)
		assert.NilError(t, err)
	}

	data, err := testee.list()
	assert.NilError(t, err)
	assert.Equal(t, len(data), len(wholeData))
}

type contextWithEmbedding struct {
	embeddedStruct
}
type embeddedStruct struct {
	Val string
}

func TestWithEmbedding(t *testing.T) {
<<<<<<< HEAD
	testee := metadataStore{
		root:   t.TempDir(),
		config: NewConfig(func() any { return &contextWithEmbedding{} }),
	}
=======
	testDir, err := ioutil.TempDir("", t.Name())
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := metadataStore{root: testDir, config: NewConfig(func() interface{} { return &contextWithEmbedding{} })}
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	testCtxMeta := contextWithEmbedding{
		embeddedStruct: embeddedStruct{
			Val: "Hello",
		},
	}
	assert.NilError(t, testee.createOrUpdate(Metadata{Metadata: testCtxMeta, Name: "test"}))
	res, err := testee.get("test")
	assert.NilError(t, err)
	assert.Equal(t, testCtxMeta, res.Metadata)
}
