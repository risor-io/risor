package main

import (
	"io"
	"net"
	"os"
	"time"
)

type Stdio struct {
	in  io.ReadCloser
	out io.WriteCloser
}

func NewDefaultStdio() *Stdio {
	return NewStdio(os.Stdin, os.Stdout)
}

func NewStdio(in io.ReadCloser, out io.WriteCloser) *Stdio {
	return &Stdio{
		in:  in,
		out: out,
	}
}

// Read implements io.Reader interface.
func (s *Stdio) Read(b []byte) (int, error) { return s.in.Read(b) }

// Write implements io.Writer interface.
func (s *Stdio) Write(b []byte) (int, error) { return s.out.Write(b) }

// Close implements io.Closer interface.
func (s *Stdio) Close() error {
	if err := s.in.Close(); err != nil {
		return err
	}
	return s.out.Close()
}

// LocalAddr implements net.Conn interface.
func (s Stdio) LocalAddr() net.Addr { return s }

// RemoteAddr implements net.Conn interface.
func (s Stdio) RemoteAddr() net.Addr { return s }

// SetDeadline implements net.Conn interface.
func (Stdio) SetDeadline(t time.Time) error { return nil }

// SetReadDeadline implements net.Conn interface.
func (Stdio) SetReadDeadline(t time.Time) error { return nil }

// SetWriteDeadline implements net.Conn interface.
func (Stdio) SetWriteDeadline(t time.Time) error { return nil }

// Network implements net.Addr interface.
func (Stdio) Network() string { return "Stdio" }

// String implements net.Addr interface.
func (Stdio) String() string { return "Stdio" }
