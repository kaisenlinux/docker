//go:build !windows

package chrootarchive // import "github.com/docker/docker/pkg/chrootarchive"

import (
	"io"
<<<<<<< HEAD
	"net"
	"os/user"
=======
	"io/ioutil"
	"os"
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/archive"
	"github.com/pkg/errors"
)

func init() {
	// initialize nss libraries in Glibc so that the dynamic libraries are loaded in the host
	// environment not in the chroot from untrusted files.
	_, _ = user.Lookup("docker")
	_, _ = net.LookupHost("localhost")
}

func invokeUnpack(decompressedArchive io.Reader, dest string, options *archive.TarOptions, root string) error {
	relDest, err := resolvePathInChroot(root, dest)
	if err != nil {
		return err
	}

<<<<<<< HEAD
	return doUnpack(decompressedArchive, relDest, root, options)
=======
	if root != "" {
		relDest, err := filepath.Rel(root, dest)
		if err != nil {
			return err
		}
		if relDest == "." {
			relDest = "/"
		}
		if relDest[0] != '/' {
			relDest = "/" + relDest
		}
		dest = relDest
	}

	cmd := reexec.Command("docker-untar", dest, root)
	cmd.Stdin = decompressedArchive

	cmd.ExtraFiles = append(cmd.ExtraFiles, r)
	output := bytes.NewBuffer(nil)
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Start(); err != nil {
		w.Close()
		return fmt.Errorf("Untar error on re-exec cmd: %v", err)
	}

	// write the options to the pipe for the untar exec to read
	if err := json.NewEncoder(w).Encode(options); err != nil {
		w.Close()
		return fmt.Errorf("Untar json encode to pipe failed: %v", err)
	}
	w.Close()

	if err := cmd.Wait(); err != nil {
		// when `xz -d -c -q | docker-untar ...` failed on docker-untar side,
		// we need to exhaust `xz`'s output, otherwise the `xz` side will be
		// pending on write pipe forever
		io.Copy(ioutil.Discard, decompressedArchive)

		return fmt.Errorf("Error processing tar file(%v): %s", err, output)
	}
	return nil
}

func tar() {
	runtime.LockOSThread()
	flag.Parse()

	src := flag.Arg(0)
	var root string
	if len(flag.Args()) > 1 {
		root = flag.Arg(1)
	}

	if root == "" {
		root = src
	}

	if err := realChroot(root); err != nil {
		fatal(err)
	}

	var options archive.TarOptions
	if err := json.NewDecoder(os.Stdin).Decode(&options); err != nil {
		fatal(err)
	}

	rdr, err := archive.TarWithOptions(src, &options)
	if err != nil {
		fatal(err)
	}
	defer rdr.Close()

	if _, err := io.Copy(os.Stdout, rdr); err != nil {
		fatal(err)
	}

	os.Exit(0)
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}

func invokePack(srcPath string, options *archive.TarOptions, root string) (io.ReadCloser, error) {
	relSrc, err := resolvePathInChroot(root, srcPath)
	if err != nil {
		return nil, err
	}

	// make sure we didn't trim a trailing slash with the call to `resolvePathInChroot`
	if strings.HasSuffix(srcPath, "/") && !strings.HasSuffix(relSrc, "/") {
		relSrc += "/"
	}

	return doPack(relSrc, root, options)
}

// resolvePathInChroot returns the equivalent to path inside a chroot rooted at root.
// The returned path always begins with '/'.
//
//   - resolvePathInChroot("/a/b", "/a/b/c/d") -> "/c/d"
//   - resolvePathInChroot("/a/b", "/a/b")     -> "/"
//
// The implementation is buggy, and some bugs may be load-bearing.
// Here be dragons.
func resolvePathInChroot(root, path string) (string, error) {
	if root == "" {
		return "", errors.New("root path must not be empty")
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	if rel == "." {
		rel = "/"
	}
	if rel[0] != '/' {
		rel = "/" + rel
	}
	return rel, nil
}
