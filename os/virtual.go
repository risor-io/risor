package os

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/risor-io/risor/limits"
)

var (
	_ OS = (*VirtualOS)(nil)
)

type ExitHandler func(int)

type Mount struct {
	Source FS
	Target string
	Type   string
}

type VirtualOS struct {
	ctx           context.Context
	limits        limits.Limits
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
	exitHandler   ExitHandler
	stdin         File
	stdout        File
	args          []string
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

// set the args passed to the os package for os.args()
func WithArgs(args []string) Option {
	return func(vos *VirtualOS) {
		vos.args = args
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
			vos.mounts[k] = v
		}
	}
}

// WithExitHandler sets the exit handler.
func WithExitHandler(exitHandler ExitHandler) Option {
	return func(vos *VirtualOS) {
		vos.exitHandler = exitHandler
	}
}

// WithStdin sets the stdin.
func WithStdin(stdin File) Option {
	return func(vos *VirtualOS) {
		vos.stdin = stdin
	}
}

// WithStdout sets the stdout.
func WithStdout(stdout File) Option {
	return func(vos *VirtualOS) {
		vos.stdout = stdout
	}
}

// NewVirtualOS creates a new VirtualOS configured with the given options.
func NewVirtualOS(ctx context.Context, opts ...Option) *VirtualOS {
	vos := &VirtualOS{
		ctx:    ctx,
		env:    map[string]string{},
		mounts: map[string]*Mount{},
		cwd:    "/",
		stdin:  &NilFile{},
		stdout: &NilFile{},
	}
	if lim, ok := limits.GetLimits(ctx); ok {
		vos.limits = lim
	} else {
		vos.limits = limits.New()
	}
	for _, opt := range opts {
		opt(vos)
	}
	return vos
}

func (osObj *VirtualOS) Args() []string {
	return osObj.args
}

// a way to override or set the args passed to the os package
// would typically be used when risor is employed in an embedded manner
func (osObj *VirtualOS) SetArgs(args []string) {
	osObj.args = args
}

func (osObj *VirtualOS) Chdir(dir string) error {
	osObj.cwd = dir
	return nil
}

func (osObj *VirtualOS) Create(name string) (File, error) {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return nil, fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.Create(resolvedPath)
}

func (osObj *VirtualOS) Append(name string) (File, error) {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return nil, fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.Append(resolvedPath)
}

func (osObj *VirtualOS) Environ() []string {
	var result []string
	for k, v := range osObj.env {
		result = append(result, k+"="+v)
	}
	return result
}

func (osObj *VirtualOS) Exit(code int) {
	if osObj.exitHandler != nil {
		osObj.exitHandler(code)
	}
}

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
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.Mkdir(resolvedPath, perm)
}

func (osObj *VirtualOS) MkdirAll(path string, perm FileMode) error {
	mount, resolvedPath, found := osObj.findMount(path)
	if !found {
		return fmt.Errorf("no such file or directory: %s", path)
	}
	return mount.Source.MkdirAll(resolvedPath, perm)
}

func (osObj *VirtualOS) MkdirTemp(dir, pattern string) (string, error) {
	if dir != "" {
		return "", errors.New("cannot specify directory")
	}
	if osObj.tmp == "" {
		return "", errors.New("no temporary directory")
	}
	mount, _, found := osObj.findMount(osObj.tmp)
	if !found {
		return "", fmt.Errorf("temporary directory not found: %s", osObj.tmp)
	}
	rint := rand.Int63()
	dirName := fmt.Sprintf("%d-%s", rint, pattern)
	if err := mount.Source.Mkdir(dirName, 0755); err != nil {
		return "", err
	}
	return filepath.Join(osObj.tmp, dirName), nil
}

func (osObj *VirtualOS) Open(name string) (File, error) {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return nil, fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.Open(resolvedPath)
}

func (osObj *VirtualOS) findMount(path string) (*Mount, string, bool) {
	endsWithSlash := strings.HasSuffix(path, "/")
	if !filepath.IsAbs(path) {
		path = filepath.Join(osObj.cwd, path)
	}
	path = filepath.Clean(path)
	if endsWithSlash && path != "/" {
		path += "/"
	}
	var match *Mount
	for k, v := range osObj.mounts {
		if k == path {
			// Exact match
			return v, "/", true
		}
		if strings.HasPrefix(path, k) {
			// Prefix match. Keep looking to confirm this is the longest match.
			if match == nil || len(k) > len(match.Target) {
				match = v
			}
		}
	}
	if match != nil {
		relPath := strings.TrimPrefix(path, match.Target)
		if relPath == "" {
			relPath = "/"
		}
		return match, relPath, true
	}
	return nil, "", false
}

func (osObj *VirtualOS) ReadFile(name string) ([]byte, error) {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return nil, fmt.Errorf("no such file or directory: %s", name)
	}
	file, err := mount.Source.Open(resolvedPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return osObj.limits.ReadAll(file)
}

func (osObj *VirtualOS) Remove(name string) error {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.Remove(resolvedPath)
}

func (osObj *VirtualOS) RemoveAll(path string) error {
	mount, resolvedPath, found := osObj.findMount(path)
	if !found {
		return fmt.Errorf("no such file or directory: %s", path)
	}
	return mount.Source.RemoveAll(resolvedPath)
}

func (osObj *VirtualOS) Rename(oldpath, newpath string) error {
	mountOld, resolvedPathOld, found := osObj.findMount(oldpath)
	if !found {
		return fmt.Errorf("no such file or directory: %s", oldpath)
	}
	mountNew, resolvedPathNew, found := osObj.findMount(newpath)
	if !found {
		return fmt.Errorf("no such file or directory: %s", newpath)
	}
	if mountOld != mountNew {
		return fmt.Errorf("cannot rename across filesystems: %s -> %s", oldpath, newpath)
	}
	return mountOld.Source.Rename(resolvedPathOld, resolvedPathNew)
}

func (osObj *VirtualOS) Setenv(key, value string) error {
	osObj.env[key] = value
	return nil
}

func (osObj *VirtualOS) Stat(name string) (os.FileInfo, error) {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return nil, fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.Stat(resolvedPath)
}

func (osObj *VirtualOS) Symlink(oldname, newname string) error {
	mountOld, resolvedPathOld, found := osObj.findMount(oldname)
	if !found {
		return fmt.Errorf("no such file or directory: %s", oldname)
	}
	mountNew, resolvedPathNew, found := osObj.findMount(newname)
	if !found {
		return fmt.Errorf("no such file or directory: %s", newname)
	}
	if mountOld != mountNew {
		return fmt.Errorf("cannot symlink across filesystems: %s -> %s", oldname, newname)
	}
	return mountOld.Source.Symlink(resolvedPathOld, resolvedPathNew)
}

func (osObj *VirtualOS) TempDir() string {
	return osObj.tmp
}

func (osObj *VirtualOS) Unsetenv(key string) error {
	delete(osObj.env, key)
	return nil
}

func (osObj *VirtualOS) UserCacheDir() (string, error) {
	if osObj.userCacheDir == "" {
		return "", errors.New("no user cache dir configured")
	}
	return osObj.userCacheDir, nil
}

func (osObj *VirtualOS) UserConfigDir() (string, error) {
	if osObj.userConfigDir == "" {
		return "", errors.New("no user config dir configured")
	}
	return osObj.userConfigDir, nil
}

func (osObj *VirtualOS) UserHomeDir() (string, error) {
	if osObj.userHomeDir == "" {
		return "", errors.New("no user home dir configured")
	}
	return osObj.userHomeDir, nil
}

func (osObj *VirtualOS) WriteFile(name string, data []byte, perm FileMode) error {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.WriteFile(resolvedPath, data, perm)
}

func (osObj *VirtualOS) ReadDir(name string) ([]DirEntry, error) {
	mount, resolvedPath, found := osObj.findMount(name)
	if !found {
		return nil, fmt.Errorf("no such file or directory: %s", name)
	}
	return mount.Source.ReadDir(resolvedPath)
}

func (osObj *VirtualOS) WalkDir(root string, fn WalkDirFunc) error {
	mount, resolvedPath, found := osObj.findMount(root)
	if !found {
		return fmt.Errorf("no such file or directory: %s", root)
	}
	return mount.Source.WalkDir(resolvedPath, fn)
}

func (osObj *VirtualOS) Stdin() File {
	return osObj.stdin
}

func (osObj *VirtualOS) Stdout() File {
	return osObj.stdout
}

func (osObj *VirtualOS) PathSeparator() rune {
	return os.PathSeparator
}

func (osObj *VirtualOS) PathListSeparator() rune {
	return os.PathSeparator
}
