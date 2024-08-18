package main

import (
	"context"
	"encoding/json"
	"fmt"
<<<<<<< HEAD
=======
	"io"
	"io/ioutil"
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/docker/docker/internal/testutils/specialimage"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/icmd"
	"gotest.tools/v3/skip"
)

type DockerCLISaveLoadSuite struct {
	ds *DockerSuite
}

func (s *DockerCLISaveLoadSuite) TearDownTest(ctx context.Context, c *testing.T) {
	s.ds.TearDownTest(ctx, c)
}

func (s *DockerCLISaveLoadSuite) OnTimeout(c *testing.T) {
	s.ds.OnTimeout(c)
}

// save a repo using gz compression and try to load it using stdout
func (s *DockerCLISaveLoadSuite) TestSaveXzAndLoadRepoStdout(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	name := "test-save-xz-and-load-repo-stdout"
	cli.DockerCmd(c, "run", "--name", name, "busybox", "true")

	imgRepoName := "foobar-save-load-test-xz-gz"
	out := cli.DockerCmd(c, "commit", name, imgRepoName).Combined()

	cli.DockerCmd(c, "inspect", imgRepoName)

	repoTarball, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", imgRepoName),
		exec.Command("xz", "-c"),
		exec.Command("gzip", "-c"))
	assert.NilError(c, err, "failed to save repo: %v %v", out, err)
	deleteImages(imgRepoName)

	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "load"},
		Stdin:   strings.NewReader(repoTarball),
	}).Assert(c, icmd.Expected{
		ExitCode: 1,
	})

	after, _, err := dockerCmdWithError("inspect", imgRepoName)
	assert.ErrorContains(c, err, "", "the repo should not exist: %v", after)
}

// save a repo using xz+gz compression and try to load it using stdout
func (s *DockerCLISaveLoadSuite) TestSaveXzGzAndLoadRepoStdout(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	name := "test-save-xz-gz-and-load-repo-stdout"
	cli.DockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test-xz-gz"
	cli.DockerCmd(c, "commit", name, repoName)

	cli.DockerCmd(c, "inspect", repoName)

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("xz", "-c"),
		exec.Command("gzip", "-c"))
	assert.NilError(c, err, "failed to save repo: %v %v", out, err)

	deleteImages(repoName)

	icmd.RunCmd(icmd.Cmd{
		Command: []string{dockerBinary, "load"},
		Stdin:   strings.NewReader(out),
	}).Assert(c, icmd.Expected{
		ExitCode: 1,
	})

	after, _, err := dockerCmdWithError("inspect", repoName)
	assert.ErrorContains(c, err, "", "the repo should not exist: %v", after)
}

func (s *DockerCLISaveLoadSuite) TestSaveSingleTag(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	imgRepoName := "foobar-save-single-tag-test"
	cli.DockerCmd(c, "tag", "busybox:latest", fmt.Sprintf("%v:latest", imgRepoName))

	out := cli.DockerCmd(c, "images", "-q", "--no-trunc", imgRepoName).Stdout()
	cleanedImageID := strings.TrimSpace(out)

	filesFilter := fmt.Sprintf("(^manifest.json$|%v)", cleanedImageID)
	if testEnv.UsingSnapshotter() {
		filesFilter = fmt.Sprintf("(^index.json$|^manifest.json$|%v)", cleanedImageID)
	}
	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", fmt.Sprintf("%v:latest", imgRepoName)),
		exec.Command("tar", "t"),
		exec.Command("grep", "-E", filesFilter))
	assert.NilError(c, err, "failed to save repo with image ID and index files: %s, %v", out, err)
}

func (s *DockerCLISaveLoadSuite) TestSaveImageId(c *testing.T) {
	testRequires(c, DaemonIsLinux)

	emptyFSImage := loadSpecialImage(c, specialimage.EmptyFS)

	imgRepoName := "foobar-save-image-id-test"
	cli.DockerCmd(c, "tag", emptyFSImage, fmt.Sprintf("%v:latest", imgRepoName))

	out := cli.DockerCmd(c, "images", "-q", "--no-trunc", imgRepoName).Stdout()
	cleanedLongImageID := strings.TrimPrefix(strings.TrimSpace(out), "sha256:")

	out = cli.DockerCmd(c, "images", "-q", imgRepoName).Stdout()
	cleanedShortImageID := strings.TrimSpace(out)

	// Make sure IDs are not empty
	assert.Assert(c, cleanedLongImageID != "", "Id should not be empty.")
	assert.Assert(c, cleanedShortImageID != "", "Id should not be empty.")

	saveCmd := exec.Command(dockerBinary, "save", cleanedShortImageID)
	tarCmd := exec.Command("tar", "t")

	var err error
	tarCmd.Stdin, err = saveCmd.StdoutPipe()
	assert.Assert(c, err == nil, "cannot set stdout pipe for tar: %v", err)
	grepCmd := exec.Command("grep", cleanedLongImageID)
	grepCmd.Stdin, err = tarCmd.StdoutPipe()
	assert.Assert(c, err == nil, "cannot set stdout pipe for grep: %v", err)

	assert.Assert(c, tarCmd.Start() == nil, "tar failed with error: %v", err)
	assert.Assert(c, saveCmd.Start() == nil, "docker save failed with error: %v", err)
	defer func() {
		saveCmd.Wait()
		tarCmd.Wait()
		cli.DockerCmd(c, "rmi", imgRepoName)
	}()

	out, _, err = runCommandWithOutput(grepCmd)

	assert.Assert(c, err == nil, "failed to save repo with image ID: %s, %v", out, err)
}

// save a repo and try to load it using flags
func (s *DockerCLISaveLoadSuite) TestSaveAndLoadRepoFlags(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	const name = "test-save-and-load-repo-flags"
	cli.DockerCmd(c, "run", "--name", name, "busybox", "true")

	const imgRepoName = "foobar-save-load-test"

	deleteImages(imgRepoName)
	cli.DockerCmd(c, "commit", name, imgRepoName)

	beforeStr := cli.DockerCmd(c, "inspect", imgRepoName).Stdout()

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", imgRepoName),
		exec.Command(dockerBinary, "load"))
	assert.NilError(c, err, "failed to save and load repo: %s, %v", out, err)

	afterStr := cli.DockerCmd(c, "inspect", imgRepoName).Stdout()

	var before, after []types.ImageInspect
	err = json.Unmarshal([]byte(beforeStr), &before)
	assert.NilError(c, err, "failed to parse inspect 'before' output")
	err = json.Unmarshal([]byte(afterStr), &after)
	assert.NilError(c, err, "failed to parse inspect 'after' output")

	assert.Assert(c, is.Len(before, 1))
	assert.Assert(c, is.Len(after, 1))

	if testEnv.UsingSnapshotter() {
		// Ignore LastTagTime difference with c8d.
		// It is not stored in the image archive, but in the imageStore
		// which is a graphdrivers implementation detail.
		//
		// It works because we load the image into the same daemon which saved
		// the image. It would still fail with the graphdrivers if the image
		// was loaded into a different daemon (which should be the case in a
		// real-world scenario).
		before[0].Metadata.LastTagTime = after[0].Metadata.LastTagTime
	}

	assert.Check(c, is.DeepEqual(before, after), "inspect is not the same after a save / load")
}

func (s *DockerCLISaveLoadSuite) TestSaveWithNoExistImage(c *testing.T) {
	testRequires(c, DaemonIsLinux)

	imgName := "foobar-non-existing-image"

	out, _, err := dockerCmdWithError("save", "-o", "test-img.tar", imgName)
	assert.ErrorContains(c, err, "", "save image should fail for non-existing image")
	assert.Assert(c, strings.Contains(out, fmt.Sprintf("No such image: %s", imgName)))
}

func (s *DockerCLISaveLoadSuite) TestSaveMultipleNames(c *testing.T) {
	testRequires(c, DaemonIsLinux)

	emptyFSImage := loadSpecialImage(c, specialimage.EmptyFS)

	const imgRepoName = "foobar-save-multi-name-test"

	oneTag := fmt.Sprintf("%v-one:latest", imgRepoName)
	twoTag := fmt.Sprintf("%v-two:latest", imgRepoName)

	cli.DockerCmd(c, "tag", emptyFSImage, oneTag)
	cli.DockerCmd(c, "tag", emptyFSImage, twoTag)

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", strings.TrimSuffix(oneTag, ":latest"), twoTag),
		exec.Command("tar", "xO", "index.json"),
	)
	assert.NilError(c, err, "failed to save multiple repos: %s, %v", out, err)

<<<<<<< HEAD
	assert.Check(c, is.Contains(out, oneTag))
	assert.Check(c, is.Contains(out, twoTag))
=======
func (s *DockerSuite) TestSaveRepoWithMultipleImages(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	makeImage := func(from string, tag string) string {
		var (
			out string
		)
		out, _ = dockerCmd(c, "run", "-d", from, "true")
		cleanedContainerID := strings.TrimSpace(out)

		out, _ = dockerCmd(c, "commit", cleanedContainerID, tag)
		imageID := strings.TrimSpace(out)
		return imageID
	}

	repoName := "foobar-save-multi-images-test"
	tagFoo := repoName + ":foo"
	tagBar := repoName + ":bar"

	idFoo := makeImage("busybox:latest", tagFoo)
	idBar := makeImage("busybox:latest", tagBar)

	deleteImages(repoName)

	// create the archive
	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName, "busybox:latest"),
		exec.Command("tar", "t"))
	assert.NilError(c, err, "failed to save multiple images: %s, %v", out, err)

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var actual []string
	for _, l := range lines {
		if regexp.MustCompile(`^[a-f0-9]{64}\.json$`).Match([]byte(l)) {
			actual = append(actual, strings.TrimSuffix(l, ".json"))
		}
	}

	// make the list of expected layers
	out = inspectField(c, "busybox:latest", "Id")
	expected := []string{strings.TrimSpace(out), idFoo, idBar}

	// prefixes are not in tar
	for i := range expected {
		expected[i] = digest.Digest(expected[i]).Hex()
	}

	sort.Strings(actual)
	sort.Strings(expected)
	assert.Assert(c, is.DeepEqual(actual, expected), "archive does not contains the right layers: got %v, expected %v, output: %q", actual, expected, out)
}

// Issue #6722 #5892 ensure directories are included in changes
func (s *DockerSuite) TestSaveDirectoryPermissions(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	layerEntries := []string{"opt/", "opt/a/", "opt/a/b/", "opt/a/b/c"}
	layerEntriesAUFS := []string{"./", ".wh..wh.aufs", ".wh..wh.orph/", ".wh..wh.plnk/", "opt/", "opt/a/", "opt/a/b/", "opt/a/b/c"}

	name := "save-directory-permissions"
	tmpDir, err := ioutil.TempDir("", "save-layers-with-directories")
	assert.Assert(c, err == nil, "failed to create temporary directory: %s", err)
	extractionDirectory := filepath.Join(tmpDir, "image-extraction-dir")
	os.Mkdir(extractionDirectory, 0777)

	defer os.RemoveAll(tmpDir)
	buildImageSuccessfully(c, name, build.WithDockerfile(`FROM busybox
	RUN adduser -D user && mkdir -p /opt/a/b && chown -R user:user /opt/a
	RUN touch /opt/a/b/c && chown user:user /opt/a/b/c`))

	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", name),
		exec.Command("tar", "-xf", "-", "-C", extractionDirectory),
	)
	assert.NilError(c, err, "failed to save and extract image: %s", out)

	dirs, err := ioutil.ReadDir(extractionDirectory)
	assert.NilError(c, err, "failed to get a listing of the layer directories: %s", err)

	found := false
	for _, entry := range dirs {
		var entriesSansDev []string
		if entry.IsDir() {
			layerPath := filepath.Join(extractionDirectory, entry.Name(), "layer.tar")

			f, err := os.Open(layerPath)
			assert.NilError(c, err, "failed to open %s: %s", layerPath, err)

			defer f.Close()

			entries, err := listTar(f)
			for _, e := range entries {
				if !strings.Contains(e, "dev/") {
					entriesSansDev = append(entriesSansDev, e)
				}
			}
			assert.NilError(c, err, "encountered error while listing tar entries: %s", err)

			if reflect.DeepEqual(entriesSansDev, layerEntries) || reflect.DeepEqual(entriesSansDev, layerEntriesAUFS) {
				found = true
				break
			}
		}
	}

	assert.Assert(c, found, "failed to find the layer with the right content listing")
}

func listTar(f io.Reader) ([]string, error) {
	tr := tar.NewReader(f)
	var entries []string

	for {
		th, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			return entries, nil
		}
		if err != nil {
			return entries, err
		}
		entries = append(entries, th.Name)
	}
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}

// Test loading a weird image where one of the layers is of zero size.
// The layer.tar file is actually zero bytes, no padding or anything else.
// See issue: 18170
func (s *DockerCLISaveLoadSuite) TestLoadZeroSizeLayer(c *testing.T) {
	// TODO(vvoland): Create an OCI image with 0 bytes layer.
	skip.If(c, testEnv.UsingSnapshotter(), "input archive is not OCI compatible")

	// this will definitely not work if using remote daemon
	// very weird test
	testRequires(c, DaemonIsLinux, testEnv.IsLocalDaemon)

	cli.DockerCmd(c, "load", "-i", "testdata/emptyLayer.tar")
}

func (s *DockerCLISaveLoadSuite) TestSaveLoadParents(c *testing.T) {
	testRequires(c, DaemonIsLinux)
	skip.If(c, testEnv.UsingSnapshotter(), "Parent image property is not supported with containerd")

	makeImage := func(from string, addfile string) string {
		id := cli.DockerCmd(c, "run", "-d", from, "touch", addfile).Stdout()
		id = strings.TrimSpace(id)

		imageID := cli.DockerCmd(c, "commit", id).Stdout()
		imageID = strings.TrimSpace(imageID)

		cli.DockerCmd(c, "rm", "-f", id)
		return imageID
	}

	idFoo := makeImage("busybox", "foo")
	idBar := makeImage(idFoo, "bar")

	tmpDir, err := ioutil.TempDir("", "save-load-parents")
	assert.NilError(c, err)
	defer os.RemoveAll(tmpDir)

	c.Log("tmpdir", tmpDir)

	outfile := filepath.Join(tmpDir, "out.tar")

	cli.DockerCmd(c, "save", "-o", outfile, idBar, idFoo)
	cli.DockerCmd(c, "rmi", idBar)
	cli.DockerCmd(c, "load", "-i", outfile)

	inspectOut := inspectField(c, idBar, "Parent")
	assert.Equal(c, inspectOut, idFoo)

	inspectOut = inspectField(c, idFoo, "Parent")
	assert.Equal(c, inspectOut, "")
}

func (s *DockerCLISaveLoadSuite) TestSaveLoadNoTag(c *testing.T) {
	testRequires(c, DaemonIsLinux)

	name := "saveloadnotag"

	buildImageSuccessfully(c, name, build.WithDockerfile("FROM busybox\nENV foo=bar"))
	id := inspectField(c, name, "Id")

	// Test to make sure that save w/o name just shows imageID during load
	out, err := RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", id),
		exec.Command(dockerBinary, "load"))
	assert.NilError(c, err, "failed to save and load repo: %s, %v", out, err)

	// Should not show 'name' but should show the image ID during the load
	assert.Assert(c, !strings.Contains(out, "Loaded image: "))
	assert.Assert(c, strings.Contains(out, "Loaded image ID:"))
	assert.Assert(c, strings.Contains(out, id))
	// Test to make sure that save by name shows that name during load
	out, err = RunCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", name),
		exec.Command(dockerBinary, "load"))
	assert.NilError(c, err, "failed to save and load repo: %s, %v", out, err)

	assert.Assert(c, strings.Contains(out, "Loaded image: "+name+":latest"))
	assert.Assert(c, !strings.Contains(out, "Loaded image ID:"))
}
