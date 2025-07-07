package os

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var _ fs.FileInfo = (*GenericFileInfo)(nil)

// Flags to OpenFile wrapping those of the underlying system. Not all
// flags may be implemented on a given system.
const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
)

type FileMode = fs.FileMode

type FileInfo = fs.FileInfo

type ReadDirFile = fs.ReadDirFile

type WalkDirFunc = fs.WalkDirFunc

type DirEntry interface {
	fs.DirEntry
	HasInfo() bool
}

type File interface {
	fs.File
	io.Writer
}

type FS interface {
	Create(name string) (File, error)
	Mkdir(name string, perm FileMode) error
	MkdirAll(path string, perm FileMode) error
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm FileMode) (File, error)
	ReadFile(name string) ([]byte, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (FileInfo, error)
	Symlink(oldname, newname string) error
	WriteFile(name string, data []byte, perm FileMode) error
	ReadDir(name string) ([]DirEntry, error)
	WalkDir(root string, fn WalkDirFunc) error
}

type User interface {
	Uid() string
	Gid() string
	Username() string
	Name() string
	HomeDir() string
}

type Group interface {
	Gid() string
	Name() string
}

type OS interface {
	FS
	Args() []string
	Chdir(dir string) error
	Environ() []string
	Exit(code int)
	Getenv(key string) string
	Getpid() int
	Getuid() int
	Getwd() (dir string, err error)
	Hostname() (name string, err error)
	LookupEnv(key string) (string, bool)
	MkdirTemp(dir, pattern string) (string, error)
	Setenv(key, value string) error
	TempDir() string
	Unsetenv(key string) error
	UserCacheDir() (string, error)
	UserConfigDir() (string, error)
	UserHomeDir() (string, error)
	Stdin() File
	Stdout() File
	Stderr() File
	PathSeparator() rune
	PathListSeparator() rune
	CurrentUser() (User, error)
	LookupUser(name string) (User, error)
	LookupUid(uid string) (User, error)
	LookupGroup(name string) (Group, error)
	LookupGid(gid string) (Group, error)
}

type contextKey string

var globalScriptargs []string

const osKey = contextKey("risor:os")

// WithOS adds an OS to the context. Subsequently, when this context is present
// in the invocation of Risor builtins, this OS will be used for all related
// functionality.
func WithOS(ctx context.Context, osObj OS) context.Context {
	return context.WithValue(ctx, osKey, osObj)
}

// GetOS returns the OS from the context, if it exists.
func GetOS(ctx context.Context) (OS, bool) {
	osObj, ok := ctx.Value(osKey).(OS)
	return osObj, ok
}

// GetDefaultOS returns the OS from the context, if it exists. Otherwise, it
// returns a new SimpleOS.
func GetDefaultOS(ctx context.Context) OS {
	if osObj, found := GetOS(ctx); found {
		return osObj
	}
	return NewSimpleOS(ctx)
}

// if risor is started from the command line and args
// are passed in, this is is how the to tell the os package about them
func SetScriptArgs(args []string) {
	globalScriptargs = args
}

// if risor is started from the command line and args
// are passed in, this is is how the to get them
func GetScriptArgs() []string {
	return globalScriptargs
}

// MassagePathError transforms a fs.PathError into a new one with the base path
// removed from the Path field.
func MassagePathError(basePath string, err error) error {
	switch err := err.(type) {
	case *fs.PathError:
		// Return a new PathError with the prefix removed.
		return &fs.PathError{
			Op:   err.Op,
			Path: strings.TrimPrefix(err.Path, basePath),
			Err:  err.Err,
		}
	}
	return err
}

// ResolvePath resolves a path relative to a base path. An error is returned if
// the path is invalid.
func ResolvePath(base, path, op string) (string, error) {
	path = filepath.Clean(path)
	if strings.HasPrefix(path, "..") {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}
	if base == "" || base == "/" {
		return path, nil
	}
	return filepath.Join(base, path), nil
}

type GenericFileInfo struct {
	name    string
	size    int64
	mode    FileMode
	modTime time.Time
	isDir   bool
}

func (fi *GenericFileInfo) Name() string       { return fi.name }
func (fi *GenericFileInfo) Size() int64        { return fi.size }
func (fi *GenericFileInfo) Mode() FileMode     { return fi.mode }
func (fi *GenericFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *GenericFileInfo) IsDir() bool        { return fi.isDir }
func (fi *GenericFileInfo) Sys() interface{}   { return nil }

type GenericFileInfoOpts struct {
	Name    string
	Size    int64
	Mode    FileMode
	ModTime time.Time
	IsDir   bool
}

func NewFileInfo(opts GenericFileInfoOpts) *GenericFileInfo {
	return &GenericFileInfo{
		name:    opts.Name,
		size:    opts.Size,
		mode:    opts.Mode,
		modTime: opts.ModTime,
		isDir:   opts.IsDir,
	}
}

type GenericDirEntry struct {
	name string
	mode FileMode
	info *GenericFileInfo
}

func (de *GenericDirEntry) Name() string   { return de.name }
func (de *GenericDirEntry) IsDir() bool    { return de.mode.IsDir() }
func (de *GenericDirEntry) Type() FileMode { return de.mode.Type() }
func (de *GenericDirEntry) HasInfo() bool  { return de.info != nil }
func (de *GenericDirEntry) Info() (FileInfo, error) {
	if de.info == nil {
		return nil, errors.New("file info not available")
	}
	return de.info, nil
}

type GenericDirEntryOpts struct {
	Name string
	Mode FileMode
	Info *GenericFileInfo
}

func NewDirEntry(opts GenericDirEntryOpts) *GenericDirEntry {
	return &GenericDirEntry{
		name: opts.Name,
		mode: opts.Mode,
		info: opts.Info,
	}
}

type DirEntryWrapper struct {
	fs.DirEntry
}

func (de *DirEntryWrapper) HasInfo() bool {
	return false
}
