package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/errdefs"
	"gotest.tools/v3/assert"
)

func TestTlsCreateUpdateGetRemove(t *testing.T) {
<<<<<<< HEAD
	testee := tlsStore{root: t.TempDir()}
=======
	testDir, err := ioutil.TempDir("", "TestTlsCreateUpdateGetRemove")
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := tlsStore{root: testDir}
	_, err = testee.getData("test-ctx", "test-ep", "test-data")
	assert.Equal(t, true, IsErrTLSDataDoesNotExist(err))
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)

	const contextName = "test-ctx"

	_, err := testee.getData(contextName, "test-ep", "test-data")
	assert.ErrorType(t, err, errdefs.IsNotFound)

	err = testee.createOrUpdate(contextName, "test-ep", "test-data", []byte("data"))
	assert.NilError(t, err)
	data, err := testee.getData(contextName, "test-ep", "test-data")
	assert.NilError(t, err)
	assert.Equal(t, string(data), "data")
	err = testee.createOrUpdate(contextName, "test-ep", "test-data", []byte("data2"))
	assert.NilError(t, err)
	data, err = testee.getData(contextName, "test-ep", "test-data")
	assert.NilError(t, err)
	assert.Equal(t, string(data), "data2")

	err = testee.removeEndpoint(contextName, "test-ep")
	assert.NilError(t, err)
	_, err = testee.getData(contextName, "test-ep", "test-data")
	assert.ErrorType(t, err, errdefs.IsNotFound)
}

func TestTlsListAndBatchRemove(t *testing.T) {
	testDir, err := ioutil.TempDir("", "TestTlsListAndBatchRemove")
	assert.NilError(t, err)
	defer os.RemoveAll(testDir)
	testee := tlsStore{root: testDir}

	all := map[string]EndpointFiles{
		"ep1": {"f1", "f2", "f3"},
		"ep2": {"f1", "f2", "f3"},
		"ep3": {"f1", "f2", "f3"},
	}

	ep1ep2 := map[string]EndpointFiles{
		"ep1": {"f1", "f2", "f3"},
		"ep2": {"f1", "f2", "f3"},
	}

	const contextName = "test-ctx"
	for name, files := range all {
		for _, file := range files {
<<<<<<< HEAD
			err := testee.createOrUpdate(contextName, name, file, []byte("data"))
=======
			err = testee.createOrUpdate("test-ctx", name, file, []byte("data"))
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
			assert.NilError(t, err)
		}
	}

	resAll, err := testee.listContextData(contextName)
	assert.NilError(t, err)
	assert.DeepEqual(t, resAll, all)

	err = testee.removeEndpoint(contextName, "ep3")
	assert.NilError(t, err)
	resEp1ep2, err := testee.listContextData(contextName)
	assert.NilError(t, err)
	assert.DeepEqual(t, resEp1ep2, ep1ep2)

	err = testee.remove(contextName)
	assert.NilError(t, err)
	resEmpty, err := testee.listContextData(contextName)
	assert.NilError(t, err)
	assert.DeepEqual(t, resEmpty, map[string]EndpointFiles{})
}
