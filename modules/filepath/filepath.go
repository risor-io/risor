package filepath

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

func GetOS(ctx context.Context) os.OS {
	return os.GetDefaultOS(ctx)
}

func Abs(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.abs", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	osObj := GetOS(ctx)
	if filepath.IsAbs(path) {
		return object.NewString(filepath.Clean(path))
	}
	wd, wdErr := osObj.Getwd()
	if wdErr != nil {
		return object.NewError(wdErr)
	}
	return object.NewString(filepath.Join(wd, path))
}

func Base(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.base", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	base := filepath.Base(path)
	return object.NewString(base)
}

func Clean(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.clean", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	cleanPath := filepath.Clean(path)
	return object.NewString(cleanPath)
}

func Dir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.dir", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	dirPath := filepath.Dir(path)
	return object.NewString(dirPath)
}

func Ext(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.ext", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	extension := filepath.Ext(path)
	return object.NewString(extension)
}

func IsAbs(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.is_abs", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	isAbs := filepath.IsAbs(path)
	return object.NewBool(isAbs)
}

func Join(ctx context.Context, args ...object.Object) object.Object {
	paths := make([]string, len(args))
	for i, arg := range args {
		path, err := object.AsString(arg)
		if err != nil {
			return err
		}
		paths[i] = path
	}
	return object.NewString(filepath.Join(paths...))
}

func Match(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.match", 2, args); err != nil {
		return err
	}
	pattern, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	name, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	matched, matchErr := filepath.Match(pattern, name)
	if matchErr != nil {
		return object.NewError(matchErr)
	}
	return object.NewBool(matched)
}

func Rel(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.rel", 2, args); err != nil {
		return err
	}
	basepath, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	targpath, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	relativePath, relErr := filepath.Rel(basepath, targpath)
	if relErr != nil {
		return object.NewError(relErr)
	}
	return object.NewString(relativePath)
}

func Split(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.split", 1, args); err != nil {
		return err
	}
	path, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	dir, file := filepath.Split(path)
	return object.NewList([]object.Object{
		object.NewString(dir),
		object.NewString(file),
	})
}

func SplitList(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.split_list", 1, args); err != nil {
		return err
	}
	pathList, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	paths := filepath.SplitList(pathList)
	pathObjs := make([]object.Object, 0, len(paths))
	for _, path := range paths {
		pathObjs = append(pathObjs, object.NewString(path))
	}
	return object.NewList(pathObjs)
}

func WalkDir(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("filepath.walk_dir", 2, args); err != nil {
		return err
	}
	root, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	callFunc, found := object.GetCallFunc(ctx)
	if !found {
		return object.Errorf("eval error: filepath.walk() context did not contain a call function")
	}
	osObj := GetOS(ctx)

	type callable func(path, info, err object.Object) object.Object
	var callback callable

	switch obj := args[1].(type) {
	case *object.Builtin:
		callback = func(path, info, err object.Object) object.Object {
			return obj.Call(ctx, path, info, err)
		}
	case *object.Function:
		callback = func(path, info, err object.Object) object.Object {
			args := []object.Object{path, info, err}
			result, resultErr := callFunc(ctx, obj, args)
			if resultErr != nil {
				return object.NewError(resultErr)
			}
			return result
		}
	default:
		return object.TypeErrorf("type error: filepath.walk() expected a function (%s given)", obj.Type())
	}

	walkFn := func(path string, d fs.DirEntry, err error) error {
		var errObj object.Object
		if err != nil {
			errObj = object.NewError(err)
		} else {
			errObj = object.Nil
		}
		wrapper := os.DirEntryWrapper{DirEntry: d}
		result := callback(object.NewString(path), object.NewDirEntry(&wrapper), errObj)
		switch result := result.(type) {
		case *object.Error:
			return result.Value()
		default:
			return nil
		}
	}
	walkErr := osObj.WalkDir(root, walkFn)
	if walkErr != nil {
		return object.NewError(walkErr)
	}
	return object.Nil
}

func Module() *object.Module {
	return object.NewBuiltinsModule("filepath", map[string]object.Object{
		"abs":        object.NewBuiltin("abs", Abs),
		"base":       object.NewBuiltin("base", Base),
		"clean":      object.NewBuiltin("clean", Clean),
		"dir":        object.NewBuiltin("dir", Dir),
		"ext":        object.NewBuiltin("ext", Ext),
		"is_abs":     object.NewBuiltin("is_abs", IsAbs),
		"join":       object.NewBuiltin("join", Join),
		"match":      object.NewBuiltin("match", Match),
		"rel":        object.NewBuiltin("rel", Rel),
		"split_list": object.NewBuiltin("split_list", SplitList),
		"split":      object.NewBuiltin("split", Split),
		"walk_dir":   object.NewBuiltin("walk_dir", WalkDir),
	})
}
