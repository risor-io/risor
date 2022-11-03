package evaluator

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/skx/monkey/lexer"
	"github.com/skx/monkey/object"
	"github.com/skx/monkey/parser"
)

// Change a mode of a file - note the second argument is a string
// to emphasise octal.
func chmodFun(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	path := args[0].Inspect()
	mode := ""

	switch args[1].(type) {
	case *object.String:
		mode = args[1].(*object.String).Value
	default:
		return newError("Second argument must be string, got %v", args[1])
	}

	// convert from octal -> decimal
	result, err := strconv.ParseInt(mode, 8, 64)
	if err != nil {
		return &object.Boolean{Value: false}
	}

	// Change the mode.
	err = os.Chmod(path, os.FileMode(result))
	if err != nil {
		return &object.Boolean{Value: false}
	}
	return &object.Boolean{Value: true}
}

// evaluate a string containing monkey-code
func evalFun(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch args[0].(type) {
	case *object.String:
		txt := args[0].(*object.String).Value

		// Lex the input
		l := lexer.New(txt)

		// parse it.
		p := parser.New(l)

		// If there are no errors
		program := p.ParseProgram()
		if len(p.Errors()) == 0 {
			// evaluate it, and return the output.
			return (Eval(program, env))
		}

		// Otherwise abort.  We should have try { } catch
		// to allow this kind of error to be caught in the future!
		fmt.Printf("Error parsing eval-string: %s", txt)
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		os.Exit(1)
	}
	return newError("argument to `eval` not supported, got=%s",
		args[0].Type())
}

// exit a program.
func exitFun(args ...object.Object) object.Object {

	code := 0

	// Optionally an exit-code might be supplied as an argument
	if len(args) > 0 {
		switch arg := args[0].(type) {
		case *object.Integer:
			code = int(arg.Value)
		case *object.Float:
			code = int(arg.Value)
		}
	}

	os.Exit(code)
	return NULL
}

// convert a double/string to an int
func intFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch args[0].(type) {
	case *object.String:
		input := args[0].(*object.String).Value
		i, err := strconv.Atoi(input)
		if err == nil {
			return &object.Integer{Value: int64(i)}
		}
		return newError("Converting string '%s' to int failed %s", input, err.Error())

	case *object.Boolean:
		input := args[0].(*object.Boolean).Value
		if input {
			return &object.Integer{Value: 1}

		}
		return &object.Integer{Value: 0}
	case *object.Integer:
		// nop
		return args[0]
	case *object.Float:
		input := args[0].(*object.Float).Value
		return &object.Integer{Value: int64(input)}
	default:
		return newError("argument to `int` not supported, got=%s",
			args[0].Type())
	}
}

// length of item
func lenFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	case *object.Null:
		return &object.Integer{Value: 0}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got=%s",
			args[0].Type())
	}
}

// regular expression match
func matchFun(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to `match` must be STRING, got %s",
			args[0].Type())
	}
	if args[1].Type() != object.STRING_OBJ {
		return newError("argument to `match` must be STRING, got %s",
			args[1].Type())
	}

	//
	// Compile and match
	//
	reg := regexp.MustCompile(args[0].(*object.String).Value)
	res := reg.FindStringSubmatch(args[1].(*object.String).Value)

	if len(res) > 0 {

		newHash := make(map[object.HashKey]object.HashPair)

		//
		// If we get a match then the output is an array
		// First entry is the match, any additional parts
		// are the capture-groups.
		//
		if len(res) > 1 {
			for i := 1; i < len(res); i++ {

				// Capture groups start at index 0.
				k := &object.Integer{Value: int64(i - 1)}
				v := &object.String{Value: res[i]}

				newHashPair := object.HashPair{Key: k, Value: v}
				newHash[k.HashKey()] = newHashPair

			}
		}

		return &object.Hash{Pairs: newHash}
	}

	// No match
	return NULL
}

// mkdir
func mkdirFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to `mkdir` must be STRING, got %s",
			args[0].Type())
	}

	path := args[0].(*object.String).Value

	// Can't fail?
	mode, err := strconv.ParseInt("755", 8, 64)
	if err != nil {
		return &object.Boolean{Value: false}
	}

	err = os.MkdirAll(path, os.FileMode(mode))
	if err != nil {
		return &object.Boolean{Value: false}
	}
	return &object.Boolean{Value: true}

}

// Open a file
func openFun(args ...object.Object) object.Object {

	path := ""
	mode := "r"

	// We need at least one arg
	if len(args) < 1 {
		return newError("wrong number of arguments. got=%d, want=1+",
			len(args))
	}

	// Get the filename
	switch args[0].(type) {
	case *object.String:
		path = args[0].(*object.String).Value
	default:
		return newError("argument to `file` not supported, got=%s",
			args[0].Type())

	}

	// Get the mode (optiona)
	if len(args) > 1 {
		switch args[1].(type) {
		case *object.String:
			mode = args[1].(*object.String).Value
		default:
			return newError("argument to `file` not supported, got=%s",
				args[0].Type())

		}
	}

	// Create the object
	file := &object.File{Filename: path}
	file.Open(mode)
	return (file)
}

// set a global pragma
func pragmaFun(args ...object.Object) object.Object {

	// If more than one argument that's an error
	if len(args) > 1 {
		return newError("wrong number of arguments. got=%d, want=0|1",
			len(args))
	}

	// If one argument update to enable the given pragma
	if len(args) == 1 {
		switch args[0].(type) {
		case *object.String:
			input := args[0].(*object.String).Value
			input = strings.ToLower(input)

			if strings.HasPrefix(input, "no-") {
				real := strings.TrimPrefix(input, "no-")
				delete(PRAGMAS, real)
			} else {
				PRAGMAS[input] = 1
			}
		default:
			return newError("argument to `pragma` not supported, got=%s",
				args[0].Type())
		}
	}

	// Now return the pragmas that are in-use.
	len := len(PRAGMAS)

	// Create a new array for the results.
	array := make([]object.Object, len)

	i := 0
	for key := range PRAGMAS {
		array[i] = &object.String{Value: key}
		i++

	}
	return &object.Array{Elements: array}
}

// push something onto an array
func pushFun(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got=%s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]
	return &object.Array{Elements: newElements}
}

// output a string to stdout
func putsFun(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect())
	}
	return NULL
}

// printfFun is the implementation of our `printf` function.
func printfFun(args ...object.Object) object.Object {

	// Convert to the formatted version, via our `sprintf`
	// function.
	out := sprintfFun(args...)

	// If that returned a string then we can print it
	if out.Type() == object.STRING_OBJ {
		fmt.Print(out.(*object.String).Value)

	}

	return NULL
}

// sprintfFun is the implementation of our `sprintf` function.
func sprintfFun(args ...object.Object) object.Object {

	// We expect 1+ arguments
	if len(args) < 1 {
		return &object.Null{}
	}

	// Type-check
	if args[0].Type() != object.STRING_OBJ {
		return &object.Null{}
	}

	// Get the format-string.
	fs := args[0].(*object.String).Value

	// Convert the arguments to something go's sprintf
	// code will understand.
	argLen := len(args)
	fmtArgs := make([]interface{}, argLen-1)

	// Here we convert and assign.
	for i, v := range args[1:] {
		fmtArgs[i] = v.ToInterface()
	}

	// Call the helper
	out := fmt.Sprintf(fs, fmtArgs...)

	// And now return the value.
	return &object.String{Value: out}
}

// Get file info.
func statFun(args ...object.Object) object.Object {

	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	path := args[0].Inspect()
	info, err := os.Stat(path)

	res := make(map[object.HashKey]object.HashPair)
	if err != nil {
		// Empty hash as we've not yet set anything
		return &object.Hash{Pairs: res}
	}

	//
	// OK populate the hash
	//

	// size -> int
	sizeData := &object.Integer{Value: info.Size()}
	sizeKey := &object.String{Value: "size"}
	sizeHash := object.HashPair{Key: sizeKey, Value: sizeData}
	res[sizeKey.HashKey()] = sizeHash

	// mod-time -> int
	mtimeData := &object.Integer{Value: info.ModTime().Unix()}
	mtimeKey := &object.String{Value: "mtime"}
	mtimeHash := object.HashPair{Key: mtimeKey, Value: mtimeData}
	res[mtimeKey.HashKey()] = mtimeHash

	// Perm -> string
	permData := &object.String{Value: info.Mode().String()}
	permKey := &object.String{Value: "perm"}
	permHash := object.HashPair{Key: permKey, Value: permData}
	res[permKey.HashKey()] = permHash

	// Mode -> string  (because we want to emphasise the octal nature)
	m := fmt.Sprintf("%04o", info.Mode().Perm())
	modeData := &object.String{Value: m}
	modeKey := &object.String{Value: "mode"}
	modeHash := object.HashPair{Key: modeKey, Value: modeData}
	res[modeKey.HashKey()] = modeHash

	typeStr := "unknown"
	if info.Mode().IsDir() {
		typeStr = "directory"
	}
	if info.Mode().IsRegular() {
		typeStr = "file"
	}

	// type: string
	typeData := &object.String{Value: typeStr}
	typeKey := &object.String{Value: "type"}
	typeHash := object.HashPair{Key: typeKey, Value: typeData}
	res[typeKey.HashKey()] = typeHash

	return &object.Hash{Pairs: res}

}

// Get hash keys
func hashKeys(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `keys` must be HASH, got=%s",
			args[0].Type())
	}

	// The object we're working with
	hash := args[0].(*object.Hash)
	ents := len(hash.Pairs)

	// Create a new array for the results.
	array := make([]object.Object, ents)

	// Now copy the keys into it.
	i := 0
	for _, ent := range hash.Pairs {
		array[i] = ent.Key
		i++
	}

	// Return the array.
	return &object.Array{Elements: array}
}

// Delete a given hash-key
func hashDelete(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `delete` must be HASH, got=%s",
			args[0].Type())
	}

	// The object we're working with
	hash := args[0].(*object.Hash)

	// The key we're going to delete
	key, ok := args[1].(object.Hashable)
	if !ok {
		return newError("key `delete` into HASH must be Hashable, got=%s",
			args[1].Type())
	}

	// Make a new hash
	newHash := make(map[object.HashKey]object.HashPair)

	// Copy the values EXCEPT the one we have.
	for k, v := range hash.Pairs {
		if k != key.HashKey() {
			newHash[k] = v
		}
	}
	return &object.Hash{Pairs: newHash}
}

// set a hash-field
func setFun(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `set` must be HASH, got=%s",
			args[0].Type())
	}
	key, ok := args[1].(object.Hashable)
	if !ok {
		return newError("key `set` into HASH must be Hashable, got=%s",
			args[1].Type())
	}
	newHash := make(map[object.HashKey]object.HashPair)
	hash := args[0].(*object.Hash)
	for k, v := range hash.Pairs {
		newHash[k] = v
	}
	newHashKey := key.HashKey()
	newHashPair := object.HashPair{Key: args[1], Value: args[2]}
	newHash[newHashKey] = newHashPair
	return &object.Hash{Pairs: newHash}
}

func strFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	out := args[0].Inspect()
	return &object.String{Value: out}
}

// type of an item
func typeFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch args[0].(type) {
	case *object.String:
		return &object.String{Value: "string"}
	case *object.Regexp:
		return &object.String{Value: "regexp"}
	case *object.Boolean:
		return &object.String{Value: "bool"}
	case *object.Builtin:
		return &object.String{Value: "builtin"}
	case *object.File:
		return &object.String{Value: "file"}
	case *object.Array:
		return &object.String{Value: "array"}
	case *object.Function:
		return &object.String{Value: "function"}
	case *object.Integer:
		return &object.String{Value: "integer"}
	case *object.Float:
		return &object.String{Value: "float"}
	case *object.Hash:
		return &object.String{Value: "hash"}
	default:
		return newError("argument to `type` not supported, got=%s",
			args[0].Type())
	}
}

// Remove a file/directory.
func unlinkFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	path := args[0].Inspect()

	err := os.Remove(path)
	if err != nil {
		return &object.Boolean{Value: false}
	}
	return &object.Boolean{Value: true}
}

func init() {
	RegisterBuiltin("chmod",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (chmodFun(args...))
		})
	RegisterBuiltin("delete",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (hashDelete(args...))
		})
	RegisterBuiltin("eval",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (evalFun(env, args...))
		})
	RegisterBuiltin("exit",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (exitFun(args...))
		})
	RegisterBuiltin("int",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (intFun(args...))
		})
	RegisterBuiltin("keys",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (hashKeys(args...))
		})
	RegisterBuiltin("len",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (lenFun(args...))
		})
	RegisterBuiltin("match",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (matchFun(args...))
		})
	RegisterBuiltin("mkdir",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (mkdirFun(args...))
		})
	RegisterBuiltin("pragma",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (pragmaFun(args...))
		})
	RegisterBuiltin("open",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (openFun(args...))
		})
	RegisterBuiltin("push",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (pushFun(args...))
		})
	RegisterBuiltin("puts",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (putsFun(args...))
		})
	RegisterBuiltin("printf",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (printfFun(args...))
		})
	RegisterBuiltin("set",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (setFun(args...))
		})
	RegisterBuiltin("sprintf",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (sprintfFun(args...))
		})
	RegisterBuiltin("stat",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (statFun(args...))
		})
	RegisterBuiltin("string",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (strFun(args...))
		})
	RegisterBuiltin("type",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (typeFun(args...))
		})
	RegisterBuiltin("unlink",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (unlinkFun(args...))
		})
}
