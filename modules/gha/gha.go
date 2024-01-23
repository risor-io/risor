package gha

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

//risor:generate no-module-func

//risor:export is_debug
func isDebug(ctx context.Context) bool {
	return os.GetDefaultOS(ctx).Getenv("RUNNER_DEBUG") == "1"
}

//risor:export log_debug
func logDebug(ctx context.Context, msg object.Object) object.Object {
	message := printableValue(msg)
	stdout := os.GetDefaultOS(ctx).Stdout()
	return runWorkflowCommand(stdout, "debug", message, nil)
}

func LogNotice(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("gha.log_notice", 1, 2, args); err != nil {
		return err
	}
	message := printableValue(args[0])
	var props map[string]any
	if len(args) > 1 {
		m, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		props = asCommandProps(m)
	}
	stdout := os.GetDefaultOS(ctx).Stdout()
	return runWorkflowCommand(stdout, "notice", message, props)
}

func LogWarning(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("gha.log_warning", 1, 2, args); err != nil {
		return err
	}
	message := printableValue(args[0])
	var props map[string]any
	if len(args) > 1 {
		m, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		props = asCommandProps(m)
	}
	stdout := os.GetDefaultOS(ctx).Stdout()
	return runWorkflowCommand(stdout, "warning", message, props)
}

func LogError(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("gha.log_error", 1, 2, args); err != nil {
		return err
	}
	message := printableValue(args[0])
	var props map[string]any
	if len(args) > 1 {
		m, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		props = asCommandProps(m)
	}
	stdout := os.GetDefaultOS(ctx).Stdout()
	return runWorkflowCommand(stdout, "error", message, props)
}

//risor:export start_group
func startGroup(ctx context.Context, msg string) object.Object {
	stdout := os.GetDefaultOS(ctx).Stdout()
	return runWorkflowCommand(stdout, "group", msg, nil)
}

//risor:export end_group
func endGroup(ctx context.Context) object.Object {
	stdout := os.GetDefaultOS(ctx).Stdout()
	return runWorkflowCommand(stdout, "endgroup", "", nil)
}

//risor:export set_output
func setOutput(ctx context.Context, key string, value object.Object) object.Object {
	printableValue := printableValue(value)

	ros := os.GetDefaultOS(ctx)
	outputFile := ros.Getenv("GITHUB_OUTPUT")
	if outputFile != "" {
		return appendWorkflowFile(ros, outputFile, workflowFileKeyValue(key, printableValue))
	}

	stdout := ros.Stdout()
	// Using "::set-output::" command is deprecated, but it's a good enough fallback
	return runWorkflowCommand(stdout, "set-output", value, map[string]any{"name": key})
}

//risor:export set_env
func setEnv(ctx context.Context, key string, value object.Object) object.Object {
	valueStr := fmt.Sprint(printableValue(value))

	ros := os.GetDefaultOS(ctx)
	ros.Setenv(key, valueStr)

	envFile := ros.Getenv("GITHUB_ENV")
	if envFile != "" {
		return appendWorkflowFile(ros, envFile, workflowFileKeyValue(key, valueStr))
	}

	stdout := ros.Stdout()
	// Using "::set-env::" command is deprecated, but it's a good enough fallback
	return runWorkflowCommand(stdout, "set-env", value, map[string]any{"name": key})
}

//risor:export add_path
func addPath(ctx context.Context, path string) object.Object {
	ros := os.GetDefaultOS(ctx)
	oldPath := ros.Getenv("PATH")
	ros.Setenv("PATH", fmt.Sprintf("%s%c%s", path, ros.PathListSeparator(), oldPath))

	pathFile := ros.Getenv("GITHUB_PATH")
	if pathFile != "" {
		return appendWorkflowFile(ros, pathFile, path)
	}

	stdout := ros.Stdout()
	// Using "::add-path::" command is deprecated, but it's a good enough fallback
	return runWorkflowCommand(stdout, "add-path", path, nil)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("gha", addGeneratedBuiltins(map[string]object.Object{
		"log_notice":  object.NewBuiltin("gha.log_notice", LogNotice),
		"log_warning": object.NewBuiltin("gha.log_warning", LogWarning),
		"log_error":   object.NewBuiltin("gha.log_error", LogError),
	}))
}
