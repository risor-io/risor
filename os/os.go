package os

import (
	"context"
	"io/fs"
)

type FS = fs.FS

type File = fs.File

type FileMode = fs.FileMode

type FileInfo = fs.FileInfo

type ReadDirFile = fs.ReadDirFile

type DirEntry = fs.DirEntry

type OS interface {
	Chdir(dir string) error
	Create(name string) (File, error)
	Environ() []string
	Exit(code int)
	Getenv(key string) string
	Getpid() int
	Getuid() int
	Getwd() (dir string, err error)
	Hostname() (name string, err error)
	LookupEnv(key string) (string, bool)
	Mkdir(name string, perm FileMode) error
	MkdirAll(path string, perm FileMode) error
	MkdirTemp(dir, pattern string) (string, error)
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm FileMode) (File, error)
	ReadFile(name string) ([]byte, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Setenv(key, value string) error
	Stat(name string) (FileInfo, error)
	Symlink(oldname, newname string) error
	TempDir() string
	Unsetenv(key string) error
	UserCacheDir() (string, error)
	UserConfigDir() (string, error)
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm FileMode) error
}

type contextKey string

const osKey = contextKey("tamarin:os")

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
