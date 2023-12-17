package s3fs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/risor-io/risor/limits"
	ros "github.com/risor-io/risor/os"
)

var (
	_                    ros.FS = (*Filesystem)(nil)
	ErrMissingClient            = errors.New("missing s3 client")
	ErrMissingBucketName        = errors.New("missing bucket name")
)

type File struct {
	fs       *Filesystem
	key      string
	writeBuf *bytes.Buffer
	reader   io.ReadCloser
	once     sync.Once
	closed   chan bool
}

func (f *File) Key() string {
	return f.key
}

func (f *File) Stat() (ros.FileInfo, error) {
	result, err := f.fs.client.HeadObject(f.fs.ctx, &s3.HeadObjectInput{
		Bucket: &f.fs.bucket,
		Key:    &f.key,
	})
	if err != nil {
		return nil, err
	}
	var lastModified time.Time
	if result.LastModified != nil {
		lastModified = *result.LastModified
	}
	return ros.NewFileInfo(ros.GenericFileInfoOpts{
		Name:    filepath.Base(f.key),
		Size:    result.ContentLength,
		ModTime: lastModified,
		IsDir:   false,
	}), nil
}

func (f *File) get() error {
	if f.reader != nil {
		return nil
	}
	result, err := f.fs.client.GetObject(f.fs.ctx, &s3.GetObjectInput{
		Bucket: &f.fs.bucket,
		Key:    &f.key,
	})
	if err != nil {
		return err
	}
	f.reader = result.Body
	return nil
}

func (f *File) Read(buf []byte) (int, error) {
	if f.reader == nil {
		if err := f.get(); err != nil {
			return 0, err
		}
	}
	return f.reader.Read(buf)
}

func (f *File) Write(p []byte) (n int, err error) {
	if f.writeBuf == nil {
		f.writeBuf = &bytes.Buffer{}
	}
	return f.writeBuf.Write(p)
}

func (f *File) runWaitToClose() {
	// This is used to guarantee that f.Close() is called
	go func() {
		select {
		case <-f.closed:
			// f.Close() was called elsewhere
		case <-f.fs.ctx.Done():
			// The context is done. This goroutine should close the file.
			f.Close()
		}
	}()
}

func (f *File) Close() error {
	var rErr, wErr error
	f.once.Do(func() {
		if f.writeBuf != nil {
			_, wErr = f.fs.client.PutObject(f.fs.ctx, &s3.PutObjectInput{
				Bucket: &f.fs.bucket,
				Key:    &f.key,
				Body:   f.writeBuf,
			})
			f.writeBuf = nil
		}
		if f.reader != nil {
			rErr = f.reader.Close()
			f.reader = nil
		}
		close(f.closed)
	})
	if wErr != nil {
		return wErr
	}
	return rErr
}

func newFile(fs *Filesystem, key string) *File {
	f := &File{
		fs:     fs,
		key:    key,
		closed: make(chan bool),
	}
	f.runWaitToClose()
	return f
}

type Filesystem struct {
	ctx    context.Context
	base   string
	bucket string
	limits limits.Limits
	client *s3.Client
}

// Option is a configuration function for an S3 Filesystem.
type Option func(*Filesystem)

// WithBase sets the base directory for the filesystem.
func WithBase(dir string) Option {
	return func(fs *Filesystem) {
		fs.base = dir
	}
}

// WithClient sets the S3 client for the filesystem.
func WithClient(client *s3.Client) Option {
	return func(fs *Filesystem) {
		fs.client = client
	}
}

// WithBucket sets the S3 bucket name for the filesystem.
func WithBucket(bucket string) Option {
	return func(fs *Filesystem) {
		fs.bucket = bucket
	}
}

// New creates a new S3 filesystem with the given options.
func New(ctx context.Context, opts ...Option) (*Filesystem, error) {
	fs := &Filesystem{ctx: ctx}
	if lim, ok := limits.GetLimits(ctx); ok {
		fs.limits = lim
	} else {
		fs.limits = limits.New()
	}
	for _, opt := range opts {
		opt(fs)
	}
	if fs.client == nil {
		return nil, ErrMissingClient
	}
	if fs.bucket == "" {
		return nil, ErrMissingBucketName
	}
	return fs, nil
}

func (fs *Filesystem) resolvePath(path, op string) (string, error) {
	resolved, err := ros.ResolvePath(fs.base, path, op)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(resolved, "/"), nil
}

func (fs *Filesystem) Create(name string) (ros.File, error) {
	if strings.HasSuffix(name, "/") {
		return nil, fmt.Errorf("cannot create file with trailing slash: %s", name)
	}
	key, err := fs.resolvePath(name, "create")
	if err != nil {
		return nil, err
	}
	return newFile(fs, key), nil
}

func (fs *Filesystem) Mkdir(name string, perm ros.FileMode) error {
	key, err := fs.resolvePath(name, "mkdir")
	if err != nil {
		return err
	}
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	f := newFile(fs, key)
	if _, err := f.Write([]byte{}); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (fs *Filesystem) MkdirAll(path string, perm ros.FileMode) error {
	key, err := fs.resolvePath(path, "mkdir")
	if err != nil {
		return err
	}
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	f := newFile(fs, key)
	if _, err := f.Write([]byte{}); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (fs *Filesystem) Open(name string) (ros.File, error) {
	if strings.HasSuffix(name, "/") {
		return nil, fmt.Errorf("cannot open file with trailing slash: %s", name)
	}
	key, err := fs.resolvePath(name, "open")
	if err != nil {
		return nil, err
	}
	f := newFile(fs, key)
	if err := f.get(); err != nil {
		return nil, err
	}
	return f, nil
}

func (fs *Filesystem) ReadFile(name string) ([]byte, error) {
	if strings.HasSuffix(name, "/") {
		return nil, fmt.Errorf("cannot open file with trailing slash: %s", name)
	}
	key, err := fs.resolvePath(name, "open")
	if err != nil {
		return nil, err
	}
	f := newFile(fs, key)
	defer f.Close()
	if err := f.get(); err != nil {
		return nil, err
	}
	return fs.limits.ReadAll(f)
}

func (fs *Filesystem) Remove(name string) error {
	if strings.HasSuffix(name, "/") {
		return fmt.Errorf("cannot remove file with trailing slash: %s", name)
	}
	key, err := fs.resolvePath(name, "remove")
	if err != nil {
		return err
	}
	if _, err := fs.client.DeleteObject(fs.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}
	return nil
}

func (fs *Filesystem) RemoveAll(path string) error {
	return errors.New("remove all is not yet implemented for s3")
}

func (fs *Filesystem) Rename(oldpath, newpath string) error {
	return errors.New("not implemented")
}

func (fs *Filesystem) Stat(name string) (ros.FileInfo, error) {
	if strings.HasSuffix(name, "/") {
		return nil, fmt.Errorf("cannot open file with trailing slash: %s", name)
	}
	key, err := fs.resolvePath(name, "open")
	if err != nil {
		return nil, err
	}
	f := newFile(fs, key)
	defer f.Close()
	return f.Stat()
}

func (fs *Filesystem) Symlink(oldname, newname string) error {
	return errors.New("not implemented")
}

func (fs *Filesystem) WriteFile(name string, data []byte, perm ros.FileMode) error {
	if strings.HasSuffix(name, "/") {
		return fmt.Errorf("cannot write file with trailing slash: %s", name)
	}
	key, err := fs.resolvePath(name, "write")
	if err != nil {
		return err
	}
	f := newFile(fs, key)
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func (fs *Filesystem) ReadDir(name string) ([]ros.DirEntry, error) {
	name = strings.TrimPrefix(name, "/")
	if name != "" && !strings.HasSuffix(name, "/") {
		name += "/"
	}
	result, err := fs.client.ListObjects(fs.ctx, &s3.ListObjectsInput{
		Bucket:    aws.String(fs.bucket),
		Prefix:    aws.String(name),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}
	var entries []ros.DirEntry
	for _, obj := range result.CommonPrefixes {
		entries = append(entries, ros.NewDirEntry(ros.GenericDirEntryOpts{
			Name: strings.TrimSuffix(strings.TrimPrefix(*obj.Prefix, name), "/"),
			Mode: ros.FileMode(0660) | os.ModeDir | 0110,
		}))
	}
	for _, obj := range result.Contents {
		key := aws.ToString(obj.Key)
		keyName := strings.TrimPrefix(key, name)
		var mod time.Time
		if obj.LastModified != nil {
			mod = *obj.LastModified
		}
		size := obj.Size
		mode := ros.FileMode(0660)
		isDir := strings.HasSuffix(key, "/")
		if isDir {
			mode |= os.ModeDir | 0110
		}
		entries = append(entries, ros.NewDirEntry(ros.GenericDirEntryOpts{
			Name: keyName,
			Mode: mode,
			Info: ros.NewFileInfo(ros.GenericFileInfoOpts{
				Name:    keyName,
				Size:    size,
				Mode:    mode,
				ModTime: mod,
				IsDir:   isDir,
			}),
		}))
	}
	return entries, nil
}

func (fs *Filesystem) WalkDir(root string, fn ros.WalkDirFunc) error {
	entries, err := fs.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := fn(root, entry, nil); err != nil {
			return err
		}
		if entry.IsDir() {
			if err := fs.WalkDir(filepath.Join(root, entry.Name()), fn); err != nil {
				return err
			}
		}
	}
	return nil
}
