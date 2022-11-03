package object

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// File wraps a file.
type File struct {
	// The path this file object refers to.
	Filename string

	// Reader is a helper for file-reading.
	Reader *bufio.Reader

	// Writer is a helper for file-writing.
	Writer *bufio.Writer

	// Handle contains the filehandle we wrap.
	Handle *os.File
}

// Type returns the type of this object.
func (f *File) Type() Type {
	return FILE_OBJ
}

// Inspect returns a string-representation of the given object.
func (f *File) Inspect() string {
	return fmt.Sprintf("<file:%s>", f.Filename)
}

// Open opens the file - called only from the open-primitive where the
// Filename will have been filled in for us.
func (f *File) Open(mode string) error {

	//
	// Special case STDIN, STDOUT, STDERR.
	// We only need to setup readers/writers for these.
	//
	if f.Filename == "!STDIN!" {
		f.Reader = bufio.NewReader(os.Stdin)
		return nil
	}
	if f.Filename == "!STDOUT!" {
		f.Writer = bufio.NewWriter(os.Stdout)
		return nil
	}
	if f.Filename == "!STDERR!" {
		f.Writer = bufio.NewWriter(os.Stderr)
		return nil
	}

	//
	// Default mode is to read.
	//
	md := os.O_RDONLY

	//
	// Allow writing though ("w"), and write+append ( "wa"| "aw" )
	//
	if mode == "w" {
		md = os.O_WRONLY

		// If we're opening a file just for writing then
		// we should truncate any existing file at the open-point
		// too.
		os.Remove(f.Filename)
	} else {
		// write and append
		if strings.Contains(mode, "w") &&
			strings.Contains(mode, "a") {
			md = os.O_WRONLY
			md |= os.O_APPEND
		}
	}

	//
	// Open the file.
	//
	file, err := os.OpenFile(f.Filename, os.O_CREATE|md, 0644)
	if err != nil {
		return err
	}

	//
	// Save the handle away
	//
	f.Handle = file

	//
	// Setup the reader/writer handles, as appropriate.
	//
	if md == os.O_RDONLY {
		f.Reader = bufio.NewReader(file)
	} else {
		f.Writer = bufio.NewWriter(file)
	}

	// All good.
	return nil
}

// InvokeMethod invokes a method against the object.
// (Built-in methods only.)
func (f *File) InvokeMethod(method string, env Environment, args ...Object) Object {
	if method == "close" {
		f.Handle.Close()
		return &Boolean{Value: true}
	}
	if method == "lines" {

		// Do we not have a reader?
		if f.Reader == nil {
			return (&Null{})
		}

		// Result.
		var lines []string
		for {
			line, err := f.Reader.ReadString('\n')
			if err != nil {
				break
			}
			lines = append(lines, line)
		}

		// make results
		l := len(lines)
		result := make([]Object, l)
		for i, txt := range lines {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	}
	if method == "methods" {
		static := []string{"methods"}
		dynamic := env.Names("file.")

		var names []string
		names = append(names, static...)
		for _, e := range dynamic {
			bits := strings.Split(e, ".")
			names = append(names, bits[1])
		}
		sort.Strings(names)

		result := make([]Object, len(names))
		for i, txt := range names {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	}
	if method == "read" {
		// Check we have a reader.
		if f.Reader == nil {
			return (&String{Value: ""})
		}

		// Read and return a line.
		line, err := f.Reader.ReadString('\n')
		if err != nil {
			return (&String{Value: ""})
		}
		return &String{Value: line}
	}
	if method == "rewind" {
		// Rewind a handle by seeking to the beginning of the file.
		f.Handle.Seek(0, 0)
		return &Boolean{Value: true}
	}
	if method == "write" {

		// check we have an argument to write.
		if len(args) < 1 {
			return &Error{Message: "Missing argument to write()!"}
		}

		// Ensure we have a writer.
		if f.Writer == nil {
			return (&Null{})
		}

		// Write the text - coorcing to a string first.
		txt := args[0].Inspect()
		_, err := f.Writer.Write([]byte(txt))
		if err == nil {
			f.Writer.Flush()
			return &Boolean{Value: true}
		}

		return &Boolean{Value: false}
	}
	return nil
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (f *File) ToInterface() interface{} {
	return "<FILE>"
}
