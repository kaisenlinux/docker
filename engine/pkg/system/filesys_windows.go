package system // import "github.com/docker/docker/pkg/system"

import (
	"os"
	"regexp"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// SddlAdministratorsLocalSystem is local administrators plus NT AUTHORITY\System.
const SddlAdministratorsLocalSystem = "D:P(A;OICI;GA;;;BA)(A;OICI;GA;;;SY)"

// volumePath is a regular expression to check if a path is a Windows
// volume path (e.g., "\\?\Volume{4c1b02c1-d990-11dc-99ae-806e6f6e6963}"
// or "\\?\Volume{4c1b02c1-d990-11dc-99ae-806e6f6e6963}\").
var volumePath = regexp.MustCompile(`^\\\\\?\\Volume{[a-z0-9-]+}\\?$`)

// MkdirAllWithACL is a custom version of os.MkdirAll modified for use on Windows
// so that it is both volume path aware, and can create a directory with
// an appropriate SDDL defined ACL.
func MkdirAllWithACL(path string, _ os.FileMode, sddl string) error {
	sa, err := makeSecurityAttributes(sddl)
	if err != nil {
		return &os.PathError{Op: "mkdirall", Path: path, Err: err}
	}
	return mkdirall(path, sa)
}

// MkdirAll is a custom version of os.MkdirAll that is volume path aware for
// Windows. It can be used as a drop-in replacement for os.MkdirAll.
func MkdirAll(path string, _ os.FileMode) error {
	return mkdirall(path, nil)
}

// mkdirall is a custom version of os.MkdirAll modified for use on Windows
// so that it is both volume path aware, and can create a directory with
// a DACL.
func mkdirall(path string, perm *windows.SecurityAttributes) error {
	if volumePath.MatchString(path) {
		return nil
	}

	// The rest of this method is largely copied from os.MkdirAll and should be kept
	// as-is to ensure compatibility.

	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent.
		err = mkdirall(fixRootDirectory(path[:j-1]), perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = mkdirWithACL(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

// mkdirWithACL creates a new directory. If there is an error, it will be of
// type *PathError. .
//
// This is a modified and combined version of os.Mkdir and windows.Mkdir
// in golang to cater for creating a directory am ACL permitting full
// access, with inheritance, to any subfolder/file for Built-in Administrators
// and Local System.
func mkdirWithACL(name string, sa *windows.SecurityAttributes) error {
	if sa == nil {
		return os.Mkdir(name, 0)
	}

	namep, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}

	err = windows.CreateDirectory(namep, sa)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}
	return nil
}

// fixRootDirectory fixes a reference to a drive's root directory to
// have the required trailing slash.
func fixRootDirectory(p string) string {
	if len(p) == len(`\\?\c:`) {
		if os.IsPathSeparator(p[0]) && os.IsPathSeparator(p[1]) && p[2] == '?' && os.IsPathSeparator(p[3]) && p[5] == ':' {
			return p + `\`
		}
	}
	return p
}

func makeSecurityAttributes(sddl string) (*windows.SecurityAttributes, error) {
	var sa windows.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	var err error
	sa.SecurityDescriptor, err = windows.SecurityDescriptorFromString(sddl)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	return &sa, nil
=======
	var access uint32
	switch mode & (windows.O_RDONLY | windows.O_WRONLY | windows.O_RDWR) {
	case windows.O_RDONLY:
		access = windows.GENERIC_READ
	case windows.O_WRONLY:
		access = windows.GENERIC_WRITE
	case windows.O_RDWR:
		access = windows.GENERIC_READ | windows.GENERIC_WRITE
	}
	if mode&windows.O_CREAT != 0 {
		access |= windows.GENERIC_WRITE
	}
	if mode&windows.O_APPEND != 0 {
		access &^= windows.GENERIC_WRITE
		access |= windows.FILE_APPEND_DATA
	}
	sharemode := uint32(windows.FILE_SHARE_READ | windows.FILE_SHARE_WRITE)
	var sa *windows.SecurityAttributes
	if mode&windows.O_CLOEXEC == 0 {
		sa = makeInheritSa()
	}
	var createmode uint32
	switch {
	case mode&(windows.O_CREAT|windows.O_EXCL) == (windows.O_CREAT | windows.O_EXCL):
		createmode = windows.CREATE_NEW
	case mode&(windows.O_CREAT|windows.O_TRUNC) == (windows.O_CREAT | windows.O_TRUNC):
		createmode = windows.CREATE_ALWAYS
	case mode&windows.O_CREAT == windows.O_CREAT:
		createmode = windows.OPEN_ALWAYS
	case mode&windows.O_TRUNC == windows.O_TRUNC:
		createmode = windows.TRUNCATE_EXISTING
	default:
		createmode = windows.OPEN_EXISTING
	}
	// Use FILE_FLAG_SEQUENTIAL_SCAN rather than FILE_ATTRIBUTE_NORMAL as implemented in golang.
	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa363858(v=vs.85).aspx
	const fileFlagSequentialScan = 0x08000000 // FILE_FLAG_SEQUENTIAL_SCAN
	h, e := windows.CreateFile(pathp, access, sharemode, sa, createmode, fileFlagSequentialScan, 0)
	return h, e
}

// Helpers for TempFileSequential
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}
func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFileSequential is a copy of ioutil.TempFile, modified to use sequential
// file access. Below is the original comment from golang:
// TempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, opens the file for reading
// and writing, and returns the resulting *os.File.
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func TempFileSequential(dir, prefix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix())
		f, err = OpenFileSequential(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
}
