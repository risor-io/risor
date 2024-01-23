package filepath

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

//risor:generate no-module-func

//risor:export
func abs(ctx context.Context, path string) (string, error) {
	osObj := os.GetDefaultOS(ctx)
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	wd, wdErr := osObj.Getwd()
	if wdErr != nil {
		return "", wdErr
	}
	return filepath.Join(wd, path), nil
}

//risor:export
func base(path string) string {
	return filepath.Base(path)
}

//risor:export
func clean(path string) string {
	return filepath.Clean(path)
}

//risor:export
func dir(path string) string {
	return filepath.Dir(path)
}

//risor:export
func ext(path string) string {
	return filepath.Ext(path)
}

//risor:export is_abs
func isAbs(path string) bool {
	return filepath.IsAbs(path)
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

//risor:export
func match(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}

//risor:export
func rel(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

//risor:export
func split(path string) []string {
	dir, file := filepath.Split(path)
	return []string{dir, file}
}

//risor:export split_list
func splitList(pathList string) []string {
	return filepath.SplitList(pathList)
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
	osObj := os.GetDefaultOS(ctx)

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
		return object.Errorf("type error: filepath.walk() expected a function (%s given)", obj.Type())
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
	return object.NewBuiltinsModule("filepath", addGeneratedBuiltins(map[string]object.Object{
		"join":     object.NewBuiltin("join", Join),
		"walk_dir": object.NewBuiltin("walk_dir", WalkDir),
	}))
}
