package os

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudcmds/tamarin/v2/internal/arg"
	"github.com/cloudcmds/tamarin/v2/object"
)

func Exit(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: exit() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		os.Exit(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		os.Exit(int(obj.Value()))
	case *object.Error:
		os.Exit(1)
	}
	return object.Errorf("type error: exit() argument must be an int or error (%s given)", args[0].Type())
}

func Chdir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.chdir", 1, args); err != nil {
		return err
	}
	dir, ok := args[0].(*object.String)
	if !ok {
		return object.Errorf("type error: expected a string (got %v)", args[0].Type())
	}
	if err := os.Chdir(dir.Value()); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Getwd(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getwd", 0, args); err != nil {
		return err
	}
	dir, err := os.Getwd()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func Mkdir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.mkdir", 2, args); err != nil {
		return err
	}
	dir, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	perm, err := object.AsInt(args[1])
	if err != nil {
		return err
	}
	if err := os.Mkdir(dir, os.FileMode(perm)); err != nil {
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
	if err := os.Remove(path); err != nil {
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
	if file, err := os.Open(path); err != nil {
		return object.NewError(err)
	} else {
		return object.NewFile(ctx, file)
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
	if err := os.Rename(oldpath, newpath); err != nil {
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
	info, ioErr := os.Stat(name)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewMap(map[string]object.Object{
		"name":     object.NewString(info.Name()),
		"size":     object.NewInt(info.Size()),
		"mode":     object.NewInt(int64(info.Mode())),
		"mod_time": object.NewInt(info.ModTime().Unix()),
		"is_dir":   object.NewBool(info.IsDir()),
	})
}

func TempDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.temp_dir", 0, args); err != nil {
		return err
	}
	return object.NewString(os.TempDir())
}

func Getenv(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getenv", 1, args); err != nil {
		return err
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(os.Getenv(key))
}

func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.create", 1, args); err != nil {
		return err
	}
	name, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	file, ioErr := os.Create(name)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewFile(ctx, file)
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
	if err := os.Setenv(key, value); err != nil {
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
	if err := os.Unsetenv(key); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func LookupEnv(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.lookup_env", 1, args); err != nil {
		return err
	}
	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	value, ok := os.LookupEnv(key)
	return object.NewMap(map[string]object.Object{
		"value":  object.NewString(value),
		"exists": object.NewBool(ok),
	})
}

func ReadFile(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.read_file", 1, args); err != nil {
		return err
	}
	filename, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	bytes, ioErr := os.ReadFile(filename)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewByteSlice(bytes)
}

func WriteFile(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.write_file", 3, args); err != nil {
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
		return object.NewError(fmt.Errorf("type error: expected byte_slice or string (got %s)", args[1].Type()))
	}
	perm, err := object.AsInt(args[2])
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, data, os.FileMode(perm)); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func UserCacheDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.user_cache_dir", 0, args); err != nil {
		return err
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func UserConfigDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.user_config_dir", 0, args); err != nil {
		return err
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(dir)
}

func UserHomeDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.user_home_dir", 0, args); err != nil {
		return err
	}
	dir, err := os.UserHomeDir()
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
	if err := os.Symlink(oldname, newname); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func MkdirAll(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.mkdir_all", 2, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	perm, err := object.AsInt(args[1])
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path, os.FileMode(perm)); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func Environ(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.environ", 0, args); err != nil {
		return err
	}
	envVars := os.Environ()
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
	return object.NewInt(int64(os.Getpid()))
}

func Getuid(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.getuid", 0, args); err != nil {
		return err
	}
	return object.NewInt(int64(os.Getuid()))
}

func Hostname(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("os.hostname", 0, args); err != nil {
		return err
	}
	hostname, err := os.Hostname()
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
	tempDir, ioErr := os.MkdirTemp(dir, pattern)
	if ioErr != nil {
		return object.NewError(ioErr)
	}
	return object.NewString(tempDir)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("os", map[string]object.Object{
		"chdir":           object.NewBuiltin("chdir", Chdir),
		"create":          object.NewBuiltin("create", Create),
		"environ":         object.NewBuiltin("environ", Environ),
		"exit":            object.NewBuiltin("exit", Exit),
		"getenv":          object.NewBuiltin("getenv", Getenv),
		"getpid":          object.NewBuiltin("getpid", Getpid),
		"getuid":          object.NewBuiltin("getuid", Getuid),
		"getwd":           object.NewBuiltin("getwd", Getwd),
		"hostname":        object.NewBuiltin("hostname", Hostname),
		"lookup_env":      object.NewBuiltin("lookup_env", LookupEnv),
		"mkdir_all":       object.NewBuiltin("mkdir_all", MkdirAll),
		"mkdir_temp":      object.NewBuiltin("mkdir_temp", MkdirTemp),
		"mkdir":           object.NewBuiltin("mkdir", Mkdir),
		"open":            object.NewBuiltin("open", Open),
		"read_file":       object.NewBuiltin("read_file", ReadFile),
		"remove":          object.NewBuiltin("remove", Remove),
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
	})
}
