package directory // import "github.com/docker/docker/pkg/directory"

<<<<<<< HEAD
import "context"

// Size walks a directory tree and returns its total size in bytes.
func Size(ctx context.Context, dir string) (int64, error) {
	return calcSize(ctx, dir)
=======
import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// MoveToSubdir moves all contents of a directory to a subdirectory underneath the original path
func MoveToSubdir(oldpath, subdir string) error {

	infos, err := ioutil.ReadDir(oldpath)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.Name() != subdir {
			oldName := filepath.Join(oldpath, info.Name())
			newName := filepath.Join(oldpath, subdir, info.Name())
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}
	return nil
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}
