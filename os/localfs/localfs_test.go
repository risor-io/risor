package localfs

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalFilesystem(t *testing.T) {
	tmp, err := os.MkdirTemp("", "-risor-localfs-test")
	require.Nil(t, err)
	defer os.RemoveAll(tmp)

	ctx := context.Background()
	fs, err := New(ctx, WithBase(tmp))
	require.Nil(t, err)

	_, err = fs.Open("test.txt")
	require.NotNil(t, err)

	f, err := fs.Create("test.txt")
	require.Nil(t, err)
	require.NotNil(t, f)

	n, err := f.Write([]byte("hello world"))
	require.Nil(t, err)
	require.Equal(t, 11, n)

	err = f.Close()
	require.Nil(t, err)

	data, err := os.ReadFile(filepath.Join(tmp, "test.txt"))
	require.Nil(t, err)
	require.Equal(t, "hello world", string(data))
}

func TestLocalFilesystemStat(t *testing.T) {
	tmp, err := os.MkdirTemp("", "-risor-localfs-test")
	require.Nil(t, err)
	defer os.RemoveAll(tmp)

	ctx := context.Background()
	fs, err := New(ctx, WithBase(tmp))
	require.Nil(t, err)

	require.Nil(t, fs.WriteFile("stat.txt", []byte("hmm"), 0644))

	stat, err := fs.Stat("stat.txt")
	require.Nil(t, err)

	require.Equal(t, "stat.txt", stat.Name())
	require.Equal(t, int64(3), stat.Size())
	require.Equal(t, os.FileMode(0644), stat.Mode())
	require.Equal(t, false, stat.IsDir())

	require.Nil(t, fs.Remove("stat.txt"))

	_, err = fs.Stat("stat.txt")
	require.NotNil(t, err)
	require.True(t, os.IsNotExist(err))
	require.Equal(t, "stat /stat.txt: no such file or directory", err.Error())
}

func TestLocalFilesystemRelativePaths(t *testing.T) {
	tmp, err := os.MkdirTemp("", "-risor-localfs-test")
	require.Nil(t, err)
	defer os.RemoveAll(tmp)

	ctx := context.Background()
	lfs, err := New(ctx, WithBase(tmp))
	require.Nil(t, err)

	require.Nil(t, lfs.MkdirAll("foo/bar", 0755))

	f, err := lfs.Create("/foo/bar/test.txt")
	require.Nil(t, err)
	require.NotNil(t, f)
	f.Close()

	// Test with absolute path
	_, err = lfs.Stat("/foo/bar/test.txt")
	require.Nil(t, err)

	// Test with relative path
	_, err = lfs.Stat("foo/bar/test.txt")
	require.Nil(t, err)

	// Confirm that we can't go up a directory
	_, err = lfs.Stat("../foo/bar/test.txt")
	require.Equal(t, &fs.PathError{
		Op:   "stat",
		Path: "../foo/bar/test.txt",
		Err:  fs.ErrInvalid,
	}, err)
}
