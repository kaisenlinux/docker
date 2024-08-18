package context

import (
	"testing"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/docker"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestUpdateDescriptionOnly(t *testing.T) {
	cli, cleanup := makeFakeCli(t)
	defer cleanup()
	err := RunCreate(cli, &CreateOptions{
		Name:   "test",
		Docker: map[string]string{},
	})
	assert.NilError(t, err)
	cli.OutBuffer().Reset()
	cli.ErrBuffer().Reset()
	assert.NilError(t, RunUpdate(cli, &UpdateOptions{
		Name:        "test",
		Description: "description",
	}))
	c, err := cli.ContextStore().GetMetadata("test")
	assert.NilError(t, err)
	dc, err := command.GetDockerContext(c)
	assert.NilError(t, err)
	assert.Equal(t, dc.Description, "description")

	assert.Equal(t, "test\n", cli.OutBuffer().String())
	assert.Equal(t, "Successfully updated context \"test\"\n", cli.ErrBuffer().String())
}

func TestUpdateDockerOnly(t *testing.T) {
<<<<<<< HEAD
	cli := makeFakeCli(t)
	createTestContext(t, cli, "test", nil)
=======
	cli, cleanup := makeFakeCli(t)
	defer cleanup()
	createTestContextWithKubeAndSwarm(t, cli, "test", "swarm")
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	assert.NilError(t, RunUpdate(cli, &UpdateOptions{
		Name: "test",
		Docker: map[string]string{
			keyHost: "tcp://some-host",
		},
	}))
	c, err := cli.ContextStore().GetMetadata("test")
	assert.NilError(t, err)
	dc, err := command.GetDockerContext(c)
	assert.NilError(t, err)
	assert.Equal(t, dc.Description, "description of test")
	assert.Check(t, is.Contains(c.Endpoints, docker.DockerEndpoint))
	assert.Equal(t, c.Endpoints[docker.DockerEndpoint].(docker.EndpointMeta).Host, "tcp://some-host")
}

<<<<<<< HEAD
=======
func TestUpdateStackOrchestratorStrategy(t *testing.T) {
	cli, cleanup := makeFakeCli(t)
	defer cleanup()
	err := RunCreate(cli, &CreateOptions{
		Name:                     "test",
		DefaultStackOrchestrator: "swarm",
		Docker:                   map[string]string{},
	})
	assert.NilError(t, err)
	err = RunUpdate(cli, &UpdateOptions{
		Name:                     "test",
		DefaultStackOrchestrator: "kubernetes",
	})
	assert.ErrorContains(t, err, `cannot specify orchestrator "kubernetes" without configuring a Kubernetes endpoint`)
}

func TestUpdateStackOrchestratorStrategyRemoveKubeEndpoint(t *testing.T) {
	cli, cleanup := makeFakeCli(t)
	defer cleanup()
	createTestContextWithKubeAndSwarm(t, cli, "test", "kubernetes")
	err := RunUpdate(cli, &UpdateOptions{
		Name:       "test",
		Kubernetes: map[string]string{},
	})
	assert.ErrorContains(t, err, `cannot specify orchestrator "kubernetes" without configuring a Kubernetes endpoint`)
}

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
func TestUpdateInvalidDockerHost(t *testing.T) {
	cli, cleanup := makeFakeCli(t)
	defer cleanup()
	err := RunCreate(cli, &CreateOptions{
		Name:   "test",
		Docker: map[string]string{},
	})
	assert.NilError(t, err)
	err = RunUpdate(cli, &UpdateOptions{
		Name: "test",
		Docker: map[string]string{
			keyHost: "some///invalid/host",
		},
	})
	assert.ErrorContains(t, err, "unable to parse docker host")
}
