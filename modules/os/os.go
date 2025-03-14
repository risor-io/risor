package os

import (
	"bytes"
	"context"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

func GetOS(ctx context.Context) os.OS {
	return os.GetDefaultOS(ctx)
}

func Args(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.args", 0, args); err != nil {
		return err
	}
	argz := GetOS(ctx).Args()
	items := make([]object.Object, len(argz))
	for i, arg := range argz {
		items[i] = object.NewString(arg)
	}
	return object.NewList(items)
}

func Exit(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.ArgsErrorf("args error: exit() expected at most 1 argument (%d given)", nArgs)
	}
	tos := GetOS(ctx)
	if nArgs == 0 {
		tos.Exit(0)
		return object.Nil
	}
	switch obj := args[0].(type) {
	case *object.Int:
		tos.Exit(int(obj.Value()))
		return object.EvalErrorf("eval error: exit(%d)", obj.Value())
	case *object.Error:
		tos.Exit(1)
		return object.EvalErrorf("eval error: exit(%s)", obj.Value().Error())
	}
	return object.TypeErrorf("type error: exit() argument must be an int or error (%s given)", args[0].Type())
}

func Chdir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.chdir", 1, args); err != nil {
		return err
	}
	dir, ok := args[0].(*object.String)
	if !ok {
		return object.TypeErrorf("type error: expected a string (got %v)", args[0].Type())
	}
	if err := GetOS(ctx).Chdir(dir.Value()); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Getwd(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getwd", 0, args); err != nil {
		return err
	}
	dir, err := GetOS(ctx).Getwd()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func Mkdir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("os.mkdir", 1, 2, args); err != nil {
		return err
	}
	dir, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	perm := int64(0o755)
	if len(args) == 2 {
		perm, err = object.AsInt(args[1])
		if err != nil {
			return err
		}
	}
	if err := GetOS(ctx).Mkdir(dir, os.FileMode(perm)); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Remove(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.remove", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	if err := GetOS(ctx).Remove(path); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func RemoveAll(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.remove_all", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	if err := GetOS(ctx).RemoveAll(path); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Open(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.open", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	if file, err := GetOS(ctx).Open(path); err != nil {
		return object.NewError(err)
	} else {
		return object.NewFile(ctx, file, path)
	}
}

func Rename(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.rename", 2, args); err != nil {
		return err
	}
	oldpath, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	newpath, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	if err := GetOS(ctx).Rename(oldpath, newpath); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Stat(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.stat", 1, args); err != nil {
		return err
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	info, ioErr := GetOS(ctx).Stat(name)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewFileInfo(info)
}

func TempDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.temp_dir", 0, args); err != nil {
		return err
	}
	return object.NewString(GetOS(ctx).TempDir())
}

func Getenv(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getenv", 1, args); err != nil {
		return err
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(GetOS(ctx).Getenv(key))
}

func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.create", 1, args); err != nil {
		return err
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	file, ioErr := GetOS(ctx).Create(name)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewFile(ctx, file, name)
}

func Setenv(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.setenv", 2, args); err != nil {
		return err
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	value, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	if err := GetOS(ctx).Setenv(key, value); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Unsetenv(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.unsetenv", 1, args); err != nil {
		return err
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	if err := GetOS(ctx).Unsetenv(key); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func ReadFile(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.read_file", 1, args); err != nil {
		return err
	}
	filename, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	bytes, ioErr := GetOS(ctx).ReadFile(filename)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewByteSlice(bytes)
}

func ReadDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("os.read_dir", 0, 1, args); err != nil {
		return err
	}
	var dirName string
	osObj := GetOS(ctx)
	if len(args) == 0 {
		var err error
		dirName, err = osObj.Getwd()
		if err != nil {
			return object.NewError(err)
		}
	} else {
		var err *object.Error
		dirName, err = object.AsString(args[0])
		if err != nil {
			return err
		}
	}
	entries, ioErr := osObj.ReadDir(dirName)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	items := make([]object.Object, 0, len(entries))
	for _, entry := range entries {
		var infoObj *object.FileInfo
		if entry.HasInfo() {
			info, _ := entry.Info()
			infoObj = object.NewFileInfo(info)
		}
		items = append(items, object.NewDirEntry(entry, infoObj))
	}
	return object.NewList(items)
}

func WriteFile(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("os.write_file", 2, 3, args); err != nil {
		return err
	}
	filename, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var data []byte
	switch arg := args[1].(type) {
	case *object.ByteSlice:
		data = arg.Value()
	case *object.String:
		data = []byte(arg.Value())
	default:
		return object.TypeErrorf("type error: expected byte_slice or string (got %s)", args[1].Type())
	}
	var perm int64 = 0o644
	if len(args) == 3 {
		perm, err = object.AsInt(args[2])
		if err != nil {
			return err
		}
	}
	if err := GetOS(ctx).WriteFile(filename, data, os.FileMode(perm)); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func UserCacheDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.user_cache_dir", 0, args); err != nil {
		return err
	}
	dir, err := GetOS(ctx).UserCacheDir()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func UserConfigDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.user_config_dir", 0, args); err != nil {
		return err
	}
	dir, err := GetOS(ctx).UserConfigDir()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func UserHomeDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.user_home_dir", 0, args); err != nil {
		return err
	}
	dir, err := GetOS(ctx).UserHomeDir()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func Symlink(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.symlink", 2, args); err != nil {
		return err
	}
	oldname, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	newname, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	if err := GetOS(ctx).Symlink(oldname, newname); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func MkdirAll(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("os.mkdir_all", 1, 2, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	perm := 0o755
	if len(args) == 2 {
		givenPerm, err := object.AsInt(args[1])
		if err != nil {
			return err
		}
		perm = int(givenPerm)
	}
	if err := GetOS(ctx).MkdirAll(path, os.FileMode(perm)); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Environ(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.environ", 0, args); err != nil {
		return err
	}
	envVars := GetOS(ctx).Environ()
	items := make([]object.Object, len(envVars))
	for i, envVar := range envVars {
		items[i] = object.NewString(envVar)
	}
	return object.NewList(items)
}

func Getpid(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getpid", 0, args); err != nil {
		return err
	}
	return object.NewInt(int64(GetOS(ctx).Getpid()))
}

func Getuid(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getuid", 0, args); err != nil {
		return err
	}
	return object.NewInt(int64(GetOS(ctx).Getuid()))
}

func Hostname(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.hostname", 0, args); err != nil {
		return err
	}
	hostname, err := GetOS(ctx).Hostname()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(hostname)
}

func MkdirTemp(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.mkdir_temp", 2, args); err != nil {
		return err
	}
	dir, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	pattern, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	tempDir, ioErr := GetOS(ctx).MkdirTemp(dir, pattern)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewString(tempDir)
}

func Copy(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("cp", 2, args); err != nil {
		return err
	}
	src, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	dst, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	os := GetOS(ctx)
	srcData, ioErr := os.ReadFile(src)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	if ioErr := os.WriteFile(dst, srcData, 0o644); ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.Nil
}

func Cat(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("cat", 1, 100, args); err != nil {
		return err
	}
	os := GetOS(ctx)
	var buf bytes.Buffer
	for _, arg := range args {
		// Read the file and append it to the buffer.
		filename, err := object.AsString(arg)
		if err != nil {
			return err
		}
		bytes, ioErr := os.ReadFile(filename)
		if ioErr != nil {
			return object.NewError(ioErr)
		}
		buf.Write(bytes)
	}
	return object.NewString(buf.String())
}

func Module() *object.Module {
	return object.NewBuiltinsModule("os", map[string]object.Object{
		"args":            object.NewBuiltin("args", Args),
		"chdir":           object.NewBuiltin("chdir", Chdir),
		"create":          object.NewBuiltin("create", Create),
		"environ":         object.NewBuiltin("environ", Environ),
		"exit":            object.NewBuiltin("exit", Exit),
		"getenv":          object.NewBuiltin("getenv", Getenv),
		"getpid":          object.NewBuiltin("getpid", Getpid),
		"getuid":          object.NewBuiltin("getuid", Getuid),
		"getwd":           object.NewBuiltin("getwd", Getwd),
		"hostname":        object.NewBuiltin("hostname", Hostname),
		"mkdir_all":       object.NewBuiltin("mkdir_all", MkdirAll),
		"mkdir_temp":      object.NewBuiltin("mkdir_temp", MkdirTemp),
		"mkdir":           object.NewBuiltin("mkdir", Mkdir),
		"open":            object.NewBuiltin("open", Open),
		"read_dir":        object.NewBuiltin("read_dir", ReadDir),
		"read_file":       object.NewBuiltin("read_file", ReadFile),
		"remove":          object.NewBuiltin("remove", Remove),
		"remove_all":      object.NewBuiltin("remove_all", RemoveAll),
		"rename":          object.NewBuiltin("rename", Rename),
		"setenv":          object.NewBuiltin("setenv", Setenv),
		"stat":            object.NewBuiltin("stat", Stat),
		"symlink":         object.NewBuiltin("symlink", Symlink),
		"temp_dir":        object.NewBuiltin("temp_dir", TempDir),
		"unsetenv":        object.NewBuiltin("unsetenv", Unsetenv),
		"user_cache_dir":  object.NewBuiltin("user_cache_dir", UserCacheDir),
		"user_config_dir": object.NewBuiltin("user_config_dir", UserConfigDir),
		"user_home_dir":   object.NewBuiltin("user_home_dir", UserHomeDir),
		"write_file":      object.NewBuiltin("write_file", WriteFile),
		"stdin": object.NewDynamicAttr("stdin", func(ctx context.Context, name string) (object.Object, error) {
			f := GetOS(ctx).Stdin()
			return object.NewFile(ctx, f, "/dev/stdin"), nil
		}),
		"stdout": object.NewDynamicAttr("stdout", func(ctx context.Context, name string) (object.Object, error) {
			f := GetOS(ctx).Stdout()
			return object.NewFile(ctx, f, "/dev/stdout"), nil
		}),
	})
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"cat":      object.NewBuiltin("cat", Cat),
		"cd":       object.NewBuiltin("cd", Chdir),
		"cp":       object.NewBuiltin("cp", Copy),
		"getenv":   object.NewBuiltin("getenv", Getenv),
		"ls":       object.NewBuiltin("ls", ReadDir),
		"setenv":   object.NewBuiltin("setenv", Setenv),
		"unsetenv": object.NewBuiltin("unsetenv", Unsetenv),
		"open":     object.NewBuiltin("open", Open),
	}
}
