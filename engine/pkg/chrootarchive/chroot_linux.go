package chrootarchive // import "github.com/docker/docker/pkg/chrootarchive"

import (
<<<<<<< HEAD
	"github.com/docker/docker/internal/mounttree"
	"github.com/docker/docker/internal/unshare"
=======
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/sys"
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
	"github.com/moby/sys/mount"
	"golang.org/x/sys/unix"
)

<<<<<<< HEAD
// goInChroot starts fn in a goroutine where the root directory, current working
// directory and umask are unshared from other goroutines and the root directory
// has been changed to path. These changes are only visible to the goroutine in
// which fn is executed. Any other goroutines, including ones started from fn,
// will see the same root directory and file system attributes as the rest of
// the process.
func goInChroot(path string, fn func()) error {
	return unshare.Go(
		unix.CLONE_FS|unix.CLONE_NEWNS,
		func() error {
			// Make everything in new ns slave.
			// Don't use `private` here as this could race where the mountns gets a
			//   reference to a mount and an unmount from the host does not propagate,
			//   which could potentially cause transient errors for other operations,
			//   even though this should be relatively small window here `slave` should
			//   not cause any problems.
			if err := mount.MakeRSlave("/"); err != nil {
				return err
=======
// chroot on linux uses pivot_root instead of chroot
// pivot_root takes a new root and an old root.
// Old root must be a sub-dir of new root, it is where the current rootfs will reside after the call to pivot_root.
// New root is where the new rootfs is set to.
// Old root is removed after the call to pivot_root so it is no longer available under the new root.
// This is similar to how libcontainer sets up a container's rootfs
func chroot(path string) (err error) {
	// if the engine is running in a user namespace we need to use actual chroot
	if sys.RunningInUserNS() {
		return realChroot(path)
	}
	if err := unix.Unshare(unix.CLONE_NEWNS); err != nil {
		return fmt.Errorf("Error creating mount namespace before pivot: %v", err)
	}

	// Make everything in new ns slave.
	// Don't use `private` here as this could race where the mountns gets a
	//   reference to a mount and an unmount from the host does not propagate,
	//   which could potentially cause transient errors for other operations,
	//   even though this should be relatively small window here `slave` should
	//   not cause any problems.
	if err := mount.MakeRSlave("/"); err != nil {
		return err
	}

	if mounted, _ := mountinfo.Mounted(path); !mounted {
		if err := mount.Mount(path, path, "bind", "rbind,rw"); err != nil {
			return realChroot(path)
		}
	}

	// setup oldRoot for pivot_root
	pivotDir, err := ioutil.TempDir(path, ".pivot_root")
	if err != nil {
		return fmt.Errorf("Error setting up pivot dir: %v", err)
	}

	var mounted bool
	defer func() {
		if mounted {
			// make sure pivotDir is not mounted before we try to remove it
			if errCleanup := unix.Unmount(pivotDir, unix.MNT_DETACH); errCleanup != nil {
				if err == nil {
					err = errCleanup
				}
				return
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
			}

			return mounttree.SwitchRoot(path)
		},
		fn,
	)
}
