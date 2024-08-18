//go:build !windows

package chrootarchive // import "github.com/docker/docker/pkg/chrootarchive"

import (
	"io"
<<<<<<< HEAD
=======
	"io/ioutil"
	"os"
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	"path/filepath"

	"github.com/containerd/containerd/pkg/userns"
	"github.com/docker/docker/pkg/archive"
)

<<<<<<< HEAD
=======
type applyLayerResponse struct {
	LayerSize int64 `json:"layerSize"`
}

// applyLayer is the entry-point for docker-applylayer on re-exec. This is not
// used on Windows as it does not support chroot, hence no point sandboxing
// through chroot and rexec.
func applyLayer() {

	var (
		tmpDir  string
		err     error
		options *archive.TarOptions
	)
	runtime.LockOSThread()
	flag.Parse()

	inUserns := sys.RunningInUserNS()
	if err := chroot(flag.Arg(0)); err != nil {
		fatal(err)
	}

	// We need to be able to set any perms
	oldmask, err := system.Umask(0)
	defer system.Umask(oldmask)
	if err != nil {
		fatal(err)
	}

	if err := json.Unmarshal([]byte(os.Getenv("OPT")), &options); err != nil {
		fatal(err)
	}

	if inUserns {
		options.InUserNS = true
	}

	if tmpDir, err = ioutil.TempDir("/", "temp-docker-extract"); err != nil {
		fatal(err)
	}

	os.Setenv("TMPDIR", tmpDir)
	size, err := archive.UnpackLayer("/", os.Stdin, options)
	os.RemoveAll(tmpDir)
	if err != nil {
		fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(applyLayerResponse{size}); err != nil {
		fatal(fmt.Errorf("unable to encode layerSize JSON: %s", err))
	}

	if _, err := flush(os.Stdin); err != nil {
		fatal(err)
	}

	os.Exit(0)
}

>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
// applyLayerHandler parses a diff in the standard layer format from `layer`, and
// applies it to the directory `dest`. Returns the size in bytes of the
// contents of the layer.
func applyLayerHandler(dest string, layer io.Reader, options *archive.TarOptions, decompress bool) (size int64, err error) {
	dest = filepath.Clean(dest)
	if decompress {
		decompressed, err := archive.DecompressStream(layer)
		if err != nil {
			return 0, err
		}
		defer decompressed.Close()

		layer = decompressed
	}
	if options == nil {
		options = &archive.TarOptions{}
	}
	if userns.RunningInUserNS() {
		options.InUserNS = true
	}
	if options.ExcludePatterns == nil {
		options.ExcludePatterns = []string{}
	}
	return doUnpackLayer(dest, layer, options)
}
