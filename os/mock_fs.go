package os

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Mock filesystem implementation for testing VirtualOS file operations
type MockFS struct {
	files     map[string]*InMemoryFile
	dirs      map[string]bool
	fileInfos map[string]*mockFileInfo
}

type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fi *mockFileInfo) Name() string       { return fi.name }
func (fi *mockFileInfo) Size() int64        { return fi.size }
func (fi *mockFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *mockFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *mockFileInfo) IsDir() bool        { return fi.isDir }
func (fi *mockFileInfo) Sys() interface{}   { return nil }

type mockDirEntry struct {
	info os.FileInfo
}

func (e *mockDirEntry) Name() string               { return e.info.Name() }
func (e *mockDirEntry) IsDir() bool                { return e.info.IsDir() }
func (e *mockDirEntry) Type() os.FileMode          { return e.info.Mode().Type() }
func (e *mockDirEntry) Info() (os.FileInfo, error) { return e.info, nil }
func (e *mockDirEntry) HasInfo() bool              { return true }

func NewMockFS() *MockFS {
	return &MockFS{
		files:     make(map[string]*InMemoryFile),
		dirs:      make(map[string]bool),
		fileInfos: make(map[string]*mockFileInfo),
	}
}

func (fs *MockFS) Create(name string) (File, error) {
	file := &InMemoryFile{
		name: name,
		data: []byte{}, // Initialize with empty slice instead of nil
	}
	fs.files[name] = file
	fs.fileInfos[name] = &mockFileInfo{
		name:    filepath.Base(name),
		mode:    0644,
		isDir:   false,
		modTime: time.Now(),
	}
	return file, nil
}

func (fs *MockFS) Open(name string) (File, error) {
	file, ok := fs.files[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	// Create a new MemoryFile with the same content but reset position
	newFile := &InMemoryFile{
		name: file.name,
		data: append([]byte{}, file.data...), // Create a copy of the data
		pos:  0,
	}
	return newFile, nil
}

func (fs *MockFS) OpenFile(name string, flag int, perm FileMode) (File, error) {
	return fs.Create(name)
}

func (fs *MockFS) Mkdir(name string, perm FileMode) error {
	fs.dirs[name] = true
	fs.fileInfos[name] = &mockFileInfo{
		name:    filepath.Base(name),
		mode:    perm | os.ModeDir,
		isDir:   true,
		modTime: time.Now(),
	}
	return nil
}

func (fs *MockFS) MkdirAll(path string, perm FileMode) error {
	// Clean the path first
	path = filepath.Clean(path)

	// If path is empty or just a slash, return nil
	if path == "" || path == "/" {
		return nil
	}

	// Create all directories in the path
	var currentPath string
	if path[0] == '/' {
		currentPath = "/"
	}

	segments := strings.Split(filepath.ToSlash(path), "/")
	for _, segment := range segments {
		if segment == "" {
			continue // Skip empty segments
		}

		// Build the path incrementally
		if currentPath == "/" {
			currentPath = "/" + segment
		} else if currentPath == "" {
			currentPath = segment
		} else {
			currentPath = currentPath + "/" + segment
		}

		// Create the directory
		fs.Mkdir(currentPath, perm)
	}

	return nil
}

func (fs *MockFS) Remove(name string) error {
	if _, ok := fs.files[name]; ok {
		delete(fs.files, name)
		delete(fs.fileInfos, name)
		return nil
	}
	if _, ok := fs.dirs[name]; ok {
		delete(fs.dirs, name)
		delete(fs.fileInfos, name)
		return nil
	}
	return &os.PathError{Op: "remove", Path: name, Err: os.ErrNotExist}
}

func (fs *MockFS) RemoveAll(path string) error {
	// Simple implementation that just calls Remove
	return fs.Remove(path)
}

func (fs *MockFS) Rename(oldname, newname string) error {
	file, ok := fs.files[oldname]
	if !ok {
		return &os.PathError{Op: "rename", Path: oldname, Err: os.ErrNotExist}
	}

	newFile := &InMemoryFile{
		name: file.name,
		data: file.data,
	}
	fs.files[newname] = newFile
	delete(fs.files, oldname)

	fileInfo := fs.fileInfos[oldname]
	fs.fileInfos[newname] = &mockFileInfo{
		name:    filepath.Base(newname),
		size:    fileInfo.size,
		mode:    fileInfo.mode,
		isDir:   fileInfo.isDir,
		modTime: fileInfo.modTime,
	}
	delete(fs.fileInfos, oldname)

	return nil
}

func (fs *MockFS) Stat(name string) (os.FileInfo, error) {
	if fileInfo, ok := fs.fileInfos[name]; ok {
		return fileInfo, nil
	}
	return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrNotExist}
}

func (fs *MockFS) Symlink(oldname, newname string) error {
	// Simple mock that just creates a file with the same content
	file, ok := fs.files[oldname]
	if !ok {
		return &os.PathError{Op: "symlink", Path: oldname, Err: os.ErrNotExist}
	}

	newFile := &InMemoryFile{
		name: file.name,
		data: file.data,
	}
	fs.files[newname] = newFile

	fs.fileInfos[newname] = &mockFileInfo{
		name:    filepath.Base(newname),
		mode:    0644 | os.ModeSymlink,
		isDir:   false,
		modTime: time.Now(),
	}

	return nil
}

func (fs *MockFS) WriteFile(name string, data []byte, perm FileMode) error {
	file, err := fs.Create(name)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	fs.fileInfos[name] = &mockFileInfo{
		name:    filepath.Base(name),
		size:    int64(len(data)),
		mode:    perm,
		isDir:   false,
		modTime: time.Now(),
	}
	return nil
}

func (fs *MockFS) ReadFile(name string) ([]byte, error) {
	file, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	memFile, ok := file.(*InMemoryFile)
	if !ok {
		return nil, &os.PathError{Op: "readfile", Path: name, Err: os.ErrInvalid}
	}

	// Since we've fixed MemoryFile.Read to track position and return EOF,
	// we need to make a copy of the contents to avoid changing the position
	// in the original file object
	data := make([]byte, len(memFile.data))
	copy(data, memFile.data)
	return data, nil
}

func (fs *MockFS) ReadDir(name string) ([]DirEntry, error) {
	if _, ok := fs.dirs[name]; !ok {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: os.ErrNotExist}
	}
	var entries []DirEntry
	for path, fileInfo := range fs.fileInfos {
		dir := filepath.Dir(path)
		if dir == name || (name == "/" && path != "" && path[0] == '/') {
			entries = append(entries, &mockDirEntry{info: fileInfo})
		}
	}
	return entries, nil
}

func (fs *MockFS) WalkDir(root string, fn WalkDirFunc) error {
	entries, err := fs.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		path := filepath.Join(root, entry.Name())
		if err := fn(path, entry, nil); err != nil {
			return err
		}
		if info.IsDir() {
			if err := fs.WalkDir(path, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
