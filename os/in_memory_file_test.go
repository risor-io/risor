package os

import (
	"io"
	"testing"
)

func TestInMemoryFile(t *testing.T) {
	t.Run("Write and Bytes", func(t *testing.T) {
		file := &InMemoryFile{}
		file.Write([]byte("123"))
		file.Rewind()

		if got := string(file.Bytes()); got != "123" {
			t.Errorf("Bytes() = %q, want %q", got, "123")
		}

		file.Seek(file.Len())
		file.Write([]byte("456"))
		file.Rewind()

		if got := string(file.Bytes()); got != "123456" {
			t.Errorf("Bytes() = %q, want %q", got, "123456")
		}

		file.Seek(3)
		file.Write([]byte("321"))
		file.Rewind()

		if got := string(file.Bytes()); got != "123321" {
			t.Errorf("Bytes() = %q, want %q", got, "123321")
		}

		file.Seek(3)
		file.Write([]byte("456789"))
		file.Rewind()

		if got := string(file.Bytes()); got != "123456789" {
			t.Errorf("Bytes() = %q, want %q", got, "123456789")
		}
	})

	t.Run("Read", func(t *testing.T) {
		file := NewInMemoryFile([]byte("hello world"))

		buf := make([]byte, 5)
		n, err := file.Read(buf)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if n != 5 {
			t.Errorf("Read() n = %v, want %v", n, 5)
		}
		if string(buf) != "hello" {
			t.Errorf("Read() buf = %q, want %q", buf, "hello")
		}

		n, err = file.Read(buf)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if n != 5 {
			t.Errorf("Read() n = %v, want %v", n, 5)
		}
		if string(buf) != " worl" {
			t.Errorf("Read() buf = %q, want %q", buf, " worl")
		}

		n, err = file.Read(buf)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if n != 1 {
			t.Errorf("Read() n = %v, want %v", n, 1)
		}
		if string(buf[:1]) != "d" {
			t.Errorf("Read() buf = %q, want %q", buf[:1], "d")
		}

		n, err = file.Read(buf)
		if err != io.EOF {
			t.Errorf("Read() error = %v, want EOF", err)
		}
		if n != 0 {
			t.Errorf("Read() n = %v, want %v", n, 0)
		}
	})

	t.Run("ReadAt", func(t *testing.T) {
		file := NewInMemoryFile([]byte("hello world"))

		buf := make([]byte, 5)
		n, err := file.ReadAt(buf, 0)
		if err != nil {
			t.Fatalf("ReadAt() error = %v", err)
		}
		if n != 5 {
			t.Errorf("ReadAt() n = %v, want %v", n, 5)
		}
		if string(buf) != "hello" {
			t.Errorf("ReadAt() buf = %q, want %q", buf, "hello")
		}

		n, err = file.ReadAt(buf, 6)
		if err != nil {
			t.Fatalf("ReadAt() error = %v", err)
		}
		if n != 5 {
			t.Errorf("ReadAt() n = %v, want %v", n, 5)
		}
		if string(buf) != "world" {
			t.Errorf("ReadAt() buf = %q, want %q", buf, "world")
		}

		// Check if ReadAt doesn't move the position
		buf = make([]byte, 5)
		_, err = file.Read(buf)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if string(buf) != "hello" {
			t.Errorf("Read() buf = %q, want %q", buf, "hello")
		}

		// Test negative offset
		_, err = file.ReadAt(buf, -1)
		if err == nil || err.Error() != "negative offset" {
			t.Errorf("ReadAt() with negative offset error = %v, want 'negative offset'", err)
		}

		// Test EOF
		_, err = file.ReadAt(buf, 20)
		if err != io.EOF {
			t.Errorf("ReadAt() with offset beyond file length error = %v, want EOF", err)
		}
	})

	t.Run("Seek and Len", func(t *testing.T) {
		file := NewInMemoryFile([]byte("hello world"))

		if file.Len() != 11 {
			t.Errorf("Len() = %v, want %v", file.Len(), 11)
		}

		file.Seek(6)
		if file.Len() != 5 {
			t.Errorf("Len() after Seek(6) = %v, want %v", file.Len(), 5)
		}

		buf := make([]byte, 5)
		n, err := file.Read(buf)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if n != 5 {
			t.Errorf("Read() n = %v, want %v", n, 5)
		}
		if string(buf) != "world" {
			t.Errorf("Read() after Seek(6) buf = %q, want %q", buf, "world")
		}

		// Test seeking to negative position
		file.Seek(-1)
		if file.pos != 0 {
			t.Errorf("Seek(-1) pos = %v, want %v", file.pos, 0)
		}

		// Test seeking beyond file length
		file.Seek(100)
		if file.pos != 100 {
			t.Errorf("Seek(100) pos = %v, want %v", file.pos, 100)
		}

		// Write after seeking beyond file length
		file.Write([]byte("!"))
		file.Rewind()
		if got := string(file.Bytes()); got != "hello world!" {
			t.Errorf("Bytes() after write beyond length = %q, want %q", got, "hello world!")
		}
	})

	t.Run("Rewind and Close", func(t *testing.T) {
		file := NewInMemoryFile([]byte("hello world"))

		file.Seek(6)
		file.Rewind()
		if file.pos != 0 {
			t.Errorf("Rewind() pos = %v, want %v", file.pos, 0)
		}

		file.Seek(6)
		file.Close()
		if file.pos != 0 {
			t.Errorf("Close() pos = %v, want %v", file.pos, 0)
		}
	})

	t.Run("Stat", func(t *testing.T) {
		file := NewInMemoryFile([]byte("hello world"))

		info, err := file.Stat()
		if err != nil {
			t.Fatalf("Stat() error = %v", err)
		}
		if info.Size() != 11 {
			t.Errorf("Stat().Size() = %v, want %v", info.Size(), 11)
		}
	})

	t.Run("NewInMemoryFile", func(t *testing.T) {
		data := []byte("hello world")
		file := NewInMemoryFile(data)

		if string(file.data) != "hello world" {
			t.Errorf("NewInMemoryFile() data = %q, want %q", file.data, "hello world")
		}

		if file.pos != 0 {
			t.Errorf("NewInMemoryFile() pos = %v, want %v", file.pos, 0)
		}
	})
}
