package os

import (
	"context"
	"testing"
)

// TestVirtualOSBasics tests basic functionality of VirtualOS
func TestVirtualOSBasics(t *testing.T) {
	ctx := context.Background()

	// Setup basic VirtualOS instance
	vos := NewVirtualOS(ctx,
		WithCwd("/home/test"),
		WithHostname("testhost"),
		WithPid(12345),
		WithUid(1000),
		WithTmp("/tmp"),
		WithUserHomeDir("/home/test"),
		WithUserCacheDir("/home/test/.cache"),
		WithUserConfigDir("/home/test/.config"),
		WithEnvironment(map[string]string{
			"HOME": "/home/test",
			"PATH": "/usr/bin:/bin",
		}),
	)

	// Test GetWd
	cwd, err := vos.Getwd()
	if err != nil {
		t.Errorf("Getwd failed: %v", err)
	}
	if cwd != "/home/test" {
		t.Errorf("Expected cwd to be /home/test, got %s", cwd)
	}

	// Test Chdir
	err = vos.Chdir("/tmp")
	if err != nil {
		t.Errorf("Chdir failed: %v", err)
	}
	cwd, _ = vos.Getwd()
	if cwd != "/tmp" {
		t.Errorf("Expected cwd to be /tmp after Chdir, got %s", cwd)
	}

	// Test Hostname
	hostname, err := vos.Hostname()
	if err != nil {
		t.Errorf("Hostname failed: %v", err)
	}
	if hostname != "testhost" {
		t.Errorf("Expected hostname to be testhost, got %s", hostname)
	}

	// Test Getpid
	pid := vos.Getpid()
	if pid != 12345 {
		t.Errorf("Expected pid to be 12345, got %d", pid)
	}

	// Test Getuid
	uid := vos.Getuid()
	if uid != 1000 {
		t.Errorf("Expected uid to be 1000, got %d", uid)
	}

	// Test TempDir
	tempDir := vos.TempDir()
	if tempDir != "/tmp" {
		t.Errorf("Expected tempDir to be /tmp, got %s", tempDir)
	}

	// Test UserHomeDir
	homeDir, err := vos.UserHomeDir()
	if err != nil {
		t.Errorf("UserHomeDir failed: %v", err)
	}
	if homeDir != "/home/test" {
		t.Errorf("Expected homeDir to be /home/test, got %s", homeDir)
	}

	// Test UserCacheDir
	cacheDir, err := vos.UserCacheDir()
	if err != nil {
		t.Errorf("UserCacheDir failed: %v", err)
	}
	if cacheDir != "/home/test/.cache" {
		t.Errorf("Expected cacheDir to be /home/test/.cache, got %s", cacheDir)
	}

	// Test UserConfigDir
	configDir, err := vos.UserConfigDir()
	if err != nil {
		t.Errorf("UserConfigDir failed: %v", err)
	}
	if configDir != "/home/test/.config" {
		t.Errorf("Expected configDir to be /home/test/.config, got %s", configDir)
	}
}

// TestVirtualOSEnvironment tests environment variable handling
func TestVirtualOSEnvironment(t *testing.T) {
	ctx := context.Background()
	vos := NewVirtualOS(ctx,
		WithEnvironment(map[string]string{
			"HOME": "/home/test",
			"PATH": "/usr/bin:/bin",
		}),
	)

	// Test Getenv
	home := vos.Getenv("HOME")
	if home != "/home/test" {
		t.Errorf("Expected HOME to be /home/test, got %s", home)
	}

	// Test Environ
	env := vos.Environ()
	if len(env) != 2 {
		t.Errorf("Expected 2 environment variables, got %d", len(env))
	}

	// Test Setenv
	err := vos.Setenv("USER", "testuser")
	if err != nil {
		t.Errorf("Setenv failed: %v", err)
	}
	user := vos.Getenv("USER")
	if user != "testuser" {
		t.Errorf("Expected USER to be testuser, got %s", user)
	}

	// Test LookupEnv
	value, found := vos.LookupEnv("PATH")
	if !found {
		t.Errorf("Expected PATH to be found")
	}
	if value != "/usr/bin:/bin" {
		t.Errorf("Expected PATH to be /usr/bin:/bin, got %s", value)
	}

	_, found = vos.LookupEnv("NONEXISTENT")
	if found {
		t.Errorf("Expected NONEXISTENT to not be found")
	}

	// Test Unsetenv
	err = vos.Unsetenv("USER")
	if err != nil {
		t.Errorf("Unsetenv failed: %v", err)
	}
	user = vos.Getenv("USER")
	if user != "" {
		t.Errorf("Expected USER to be empty after Unsetenv, got %s", user)
	}
}

// TestVirtualOSStdIO tests standard input/output handling
func TestVirtualOSStdIO(t *testing.T) {
	ctx := context.Background()
	stdin := &InMemoryFile{name: "stdin", data: []byte("test input")}
	stdout := &InMemoryFile{name: "stdout"}

	vos := NewVirtualOS(ctx,
		WithStdin(stdin),
		WithStdout(stdout),
	)

	// Instead of comparing the File objects directly, we verify stdin by reading from it
	stdinData := make([]byte, 100)
	n, err := vos.Stdin().Read(stdinData)
	if err != nil {
		t.Errorf("Reading from stdin failed: %v", err)
	}
	if string(stdinData[:n]) != "test input" {
		t.Errorf("Expected stdin contents to be 'test input', got '%s'", string(stdinData[:n]))
	}

	// Verify stdout by writing to it and checking the underlying memory file
	stdoutFile := vos.Stdout()
	n, err = stdoutFile.Write([]byte("test output"))
	if err != nil {
		t.Errorf("Writing to stdout failed: %v", err)
	}
	if n != 11 {
		t.Errorf("Expected to write 11 bytes to stdout, wrote %d", n)
	}

	// Verify stdout contents
	stdoutContents := stdout.data
	if string(stdoutContents) != "test output" {
		t.Errorf("Expected stdout contents to be 'test output', got '%s'", string(stdoutContents))
	}
}

// TestVirtualOSUsers tests user and group handling
func TestVirtualOSUsers(t *testing.T) {
	ctx := context.Background()

	testUser := &VirtualUser{
		uid:      "1000",
		gid:      "1000",
		username: "testuser",
		name:     "Test User",
		homeDir:  "/home/testuser",
	}

	testGroup := &VirtualGroup{
		gid:  "1000",
		name: "testgroup",
	}

	vos := NewVirtualOS(ctx,
		WithCurrentUser(testUser),
		WithGroup(testGroup),
	)

	// Test CurrentUser
	currentUser, err := vos.CurrentUser()
	if err != nil {
		t.Errorf("CurrentUser failed: %v", err)
	}

	if currentUser.Username() != "testuser" {
		t.Errorf("Expected current username to be testuser, got %s", currentUser.Username())
	}

	if currentUser.Uid() != "1000" {
		t.Errorf("Expected current uid to be 1000, got %s", currentUser.Uid())
	}

	// Test LookupUser
	user, err := vos.LookupUser("testuser")
	if err != nil {
		t.Errorf("LookupUser failed: %v", err)
	}

	if user.Username() != "testuser" {
		t.Errorf("Expected username to be testuser, got %s", user.Username())
	}

	// Test LookupUid
	user, err = vos.LookupUid("1000")
	if err != nil {
		t.Errorf("LookupUid failed: %v", err)
	}

	if user.Name() != "Test User" {
		t.Errorf("Expected name to be Test User, got %s", user.Name())
	}

	// Test LookupGroup
	group, err := vos.LookupGroup("testgroup")
	if err != nil {
		t.Errorf("LookupGroup failed: %v", err)
	}

	if group.Name() != "testgroup" {
		t.Errorf("Expected group name to be testgroup, got %s", group.Name())
	}

	// Test LookupGid
	group, err = vos.LookupGid("1000")
	if err != nil {
		t.Errorf("LookupGid failed: %v", err)
	}

	if group.Gid() != "1000" {
		t.Errorf("Expected gid to be 1000, got %s", group.Gid())
	}

	// Test lookup failures
	_, err = vos.LookupUser("nonexistent")
	if err == nil {
		t.Errorf("Expected LookupUser to fail for nonexistent user")
	}

	_, err = vos.LookupUid("999")
	if err == nil {
		t.Errorf("Expected LookupUid to fail for nonexistent uid")
	}

	_, err = vos.LookupGroup("nonexistent")
	if err == nil {
		t.Errorf("Expected LookupGroup to fail for nonexistent group")
	}

	_, err = vos.LookupGid("999")
	if err == nil {
		t.Errorf("Expected LookupGid to fail for nonexistent gid")
	}
}

// TestVirtualOSFileOperations tests file operations with mock filesystems
func TestVirtualOSFileOperations(t *testing.T) {
	ctx := context.Background()

	// Create mock filesystem
	mockFs := NewMockFS()
	mounts := map[string]*Mount{
		"/": {
			Source: mockFs,
			Target: "/",
			Type:   "mock",
		},
	}

	vos := NewVirtualOS(ctx,
		WithCwd("/"),
		WithMounts(mounts),
	)

	// Test WriteFile and ReadFile
	err := vos.WriteFile("/test.txt", []byte("hello world"), 0o644)
	if err != nil {
		t.Errorf("WriteFile failed: %v", err)
	}

	data, err := vos.ReadFile("/test.txt")
	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
	}

	if string(data) != "hello world" {
		t.Errorf("Expected file contents to be 'hello world', got '%s'", string(data))
	}

	// Test Create and Open
	file, err := vos.Create("/created.txt")
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	_, err = file.Write([]byte("created file"))
	if err != nil {
		t.Errorf("Writing to created file failed: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Errorf("Closing created file failed: %v", err)
	}

	file, err = vos.Open("/created.txt")
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}

	buffer := make([]byte, 100)
	n, err := file.Read(buffer)
	if err != nil {
		t.Errorf("Reading from opened file failed: %v", err)
	}

	if string(buffer[:n]) != "created file" {
		t.Errorf("Expected opened file contents to be 'created file', got '%s'", string(buffer[:n]))
	}

	// Test Mkdir and ReadDir
	err = vos.Mkdir("/testdir", 0o755)
	if err != nil {
		t.Errorf("Mkdir failed: %v", err)
	}

	err = vos.WriteFile("/testdir/file1.txt", []byte("file1"), 0o644)
	if err != nil {
		t.Errorf("WriteFile in directory failed: %v", err)
	}

	entries, err := vos.ReadDir("/testdir")
	if err != nil {
		t.Errorf("ReadDir failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry in directory, got %d", len(entries))
	}

	if entries[0].Name() != "file1.txt" {
		t.Errorf("Expected directory entry name to be file1.txt, got %s", entries[0].Name())
	}

	// Test MkdirAll
	err = vos.MkdirAll("/a/b/c", 0o755)
	if err != nil {
		t.Errorf("MkdirAll failed: %v", err)
	}

	_, err = vos.Stat("/a/b/c")
	if err != nil {
		t.Errorf("Stat on directory created with MkdirAll failed: %v", err)
	}

	// Test Remove
	err = vos.Remove("/test.txt")
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}

	_, err = vos.Stat("/test.txt")
	if err == nil {
		t.Errorf("Expected Stat to fail after Remove")
	}

	// Test Rename
	err = vos.Rename("/created.txt", "/renamed.txt")
	if err != nil {
		t.Errorf("Rename failed: %v", err)
	}

	_, err = vos.Stat("/created.txt")
	if err == nil {
		t.Errorf("Expected Stat to fail for original file after Rename")
	}

	_, err = vos.Stat("/renamed.txt")
	if err != nil {
		t.Errorf("Stat on renamed file failed: %v", err)
	}
}

// TestVirtualOSExitHandler tests the exit handler functionality
func TestVirtualOSExitHandler(t *testing.T) {
	ctx := context.Background()

	exitCode := -1
	exitHandler := func(code int) {
		exitCode = code
	}

	vos := NewVirtualOS(ctx, WithExitHandler(exitHandler))

	vos.Exit(42)

	if exitCode != 42 {
		t.Errorf("Expected exit code to be 42, got %d", exitCode)
	}
}

// TestVirtualOSArgs tests the command line arguments functionality
func TestVirtualOSArgs(t *testing.T) {
	ctx := context.Background()

	args := []string{"program", "-flag", "value"}
	vos := NewVirtualOS(ctx, WithArgs(args))

	osArgs := vos.Args()
	if len(osArgs) != 3 {
		t.Errorf("Expected 3 arguments, got %d", len(osArgs))
	}

	if osArgs[0] != "program" || osArgs[1] != "-flag" || osArgs[2] != "value" {
		t.Errorf("Arguments don't match expected values")
	}

	// Test SetArgs
	newArgs := []string{"newprogram", "-newflag"}
	vos.SetArgs(newArgs)

	osArgs = vos.Args()
	if len(osArgs) != 2 {
		t.Errorf("Expected 2 arguments after SetArgs, got %d", len(osArgs))
	}

	if osArgs[0] != "newprogram" || osArgs[1] != "-newflag" {
		t.Errorf("Arguments don't match expected values after SetArgs")
	}
}
