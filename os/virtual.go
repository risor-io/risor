package os

import (
	"context"
	"errors"
	"os"
)

type Mount struct {
	Source FS
	Target string
	Type   string
}

type VirtualOS struct {
	ctx           context.Context
	userCacheDir  string
	userConfigDir string
	userHomeDir   string
	env           map[string]string
	cwd           string
	tmp           string
	hostname      string
	pid           int
	uid           int
	mounts        map[string]*Mount
}

// Option is a configuration function for a Virtual Machine.
type Option func(*VirtualOS)

// WithUserCacheDir sets the user cache directory.
func WithUserCacheDir(dir string) Option {
	return func(vos *VirtualOS) {
		vos.userCacheDir = dir
	}
}

// WithUserConfigDir sets the user config directory.
func WithUserConfigDir(dir string) Option {
	return func(vos *VirtualOS) {
		vos.userConfigDir = dir
	}
}

// WithUserHomeDir sets the user home directory.
func WithUserHomeDir(dir string) Option {
	return func(vos *VirtualOS) {
		vos.userHomeDir = dir
	}
}

// WithEnvironment sets the user home directory.
func WithEnvironment(env map[string]string) Option {
	return func(vos *VirtualOS) {
		for k, v := range env {
			vos.env[k] = v
		}
	}
}

// WithCwd sets the current working directory.
func WithCwd(cwd string) Option {
	return func(vos *VirtualOS) {
		vos.cwd = cwd
	}
}

// WithTmp sets the path to the temporary directory.
func WithTmp(tmp string) Option {
	return func(vos *VirtualOS) {
		vos.tmp = tmp
	}
}

// WithPid sets the process ID.
func WithPid(pid int) Option {
	return func(vos *VirtualOS) {
		vos.pid = pid
	}
}

// WithUid sets the user ID.
func WithUid(uid int) Option {
	return func(vos *VirtualOS) {
		vos.uid = uid
	}
}

// WithHostname sets the hostname.
func WithHostname(hostname string) Option {
	return func(vos *VirtualOS) {
		vos.hostname = hostname
	}
}

// WithMounts sets the mounts.
func WithMounts(mounts map[string]*Mount) Option {
	return func(vos *VirtualOS) {
		for k, v := range mounts {
			mounts[k] = v
		}
	}
}

// NewVirtualOS creates a new VirtualOS configured with the given options.
func NewVirtualOS(ctx context.Context, opts ...Option) *VirtualOS {
	vos := &VirtualOS{
		ctx:    ctx,
		env:    map[string]string{},
		mounts: map[string]*Mount{},
	}
	for _, opt := range opts {
		opt(vos)
	}
	return vos
}

func (osObj *VirtualOS) Chdir(dir string) error {
	osObj.cwd = dir
	return nil
}

func (osObj *VirtualOS) Create(name string) (File, error) {
	return nil, errors.New("not implemented")
}

func (osObj *VirtualOS) Environ() []string {
	var result []string
	for k, v := range osObj.env {
		result = append(result, k+"="+v)
	}
	return result
}

func (osObj *VirtualOS) Exit(code int) {}

func (osObj *VirtualOS) Getenv(key string) string {
	return osObj.env[key]
}

func (osObj *VirtualOS) Getpid() int {
	return osObj.pid
}

func (osObj *VirtualOS) Getuid() int {
	return osObj.uid
}

func (osObj *VirtualOS) Getwd() (string, error) {
	return osObj.cwd, nil
}

func (osObj *VirtualOS) Hostname() (string, error) {
	return osObj.hostname, nil
}

func (osObj *VirtualOS) LookupEnv(key string) (string, bool) {
	value, found := osObj.env[key]
	return value, found
}

func (osObj *VirtualOS) Mkdir(name string, perm FileMode) error {
	return errors.New("not implemented")
}

func (osObj *VirtualOS) MkdirAll(path string, perm FileMode) error {
	return errors.New("not implemented")
}

func (osObj *VirtualOS) MkdirTemp(dir, pattern string) (string, error) {
	return "", errors.New("not implemented")
}

func (osObj *VirtualOS) Open(name string) (File, error) {
	return nil, errors.New("not implemented")
}

func (osObj *VirtualOS) OpenFile(name string, flag int, perm FileMode) (File, error) {
	return nil, errors.New("not implemented")
}

func (osObj *VirtualOS) ReadFile(name string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (osObj *VirtualOS) Remove(name string) error {
	return errors.New("not implemented")
}

func (osObj *VirtualOS) Rename(oldpath, newpath string) error {
	return errors.New("not implemented")
}

func (osObj *VirtualOS) Setenv(key, value string) error {
	osObj.env[key] = value
	return nil
}

func (osObj *VirtualOS) Stat(name string) (os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (osObj *VirtualOS) Symlink(oldname, newname string) error {
	return errors.New("not implemented")
}

func (osObj *VirtualOS) TempDir() string {
	return osObj.tmp
}

func (osObj *VirtualOS) Unsetenv(key string) error {
	delete(osObj.env, key)
	return nil
}

func (osObj *VirtualOS) UserCacheDir() (string, error) {
	return osObj.userCacheDir, nil
}

func (osObj *VirtualOS) UserConfigDir() (string, error) {
	return osObj.userConfigDir, nil
}

func (osObj *VirtualOS) UserHomeDir() (string, error) {
	return osObj.userHomeDir, nil
}

func (osObj *VirtualOS) WriteFile(name string, data []byte, perm FileMode) error {
	return errors.New("not implemented")
}
