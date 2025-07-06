package ssh

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// Client represents an SSH client connection and implements object.Object
type Client struct {
	client *ssh.Client
}

func (c *Client) Type() object.Type {
	return "ssh.client"
}

func (c *Client) Inspect() string {
	return "ssh.client"
}

func (c *Client) Interface() interface{} {
	return c.client
}

func (c *Client) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Client) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (c *Client) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: ssh.client object has no settable attributes")
}

func (c *Client) IsTruthy() bool {
	return c.client != nil
}

func (c *Client) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for ssh.client: %v", opType)
}

func (c *Client) Cost() int {
	return 0
}

func (c *Client) String() string {
	return c.Inspect()
}

// Session represents an SSH session and implements object.Object
type Session struct {
	session *ssh.Session
}

func (s *Session) Type() object.Type {
	return "ssh.session"
}

func (s *Session) Inspect() string {
	return "ssh.session"
}

func (s *Session) Interface() interface{} {
	return s.session
}

func (s *Session) Equals(other object.Object) object.Object {
	if s == other {
		return object.True
	}
	return object.False
}

func (s *Session) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (s *Session) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: ssh.session object has no settable attributes")
}

func (s *Session) IsTruthy() bool {
	return s.session != nil
}

func (s *Session) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for ssh.session: %v", opType)
}

func (s *Session) Cost() int {
	return 0
}

func (s *Session) String() string {
	return s.Inspect()
}

// Connect establishes an SSH connection
func Connect(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 3 {
		return object.ArgsErrorf("ssh.connect() requires at least 3 arguments: host, user, password")
	}

	host, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	user, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	password, err := object.AsString(args[2])
	if err != nil {
		return err
	}

	// Optional timeout parameter
	timeout := 30 * time.Second
	if len(args) > 3 {
		timeoutSecs, err := object.AsInt(args[3])
		if err != nil {
			return err
		}
		timeout = time.Duration(timeoutSecs) * time.Second
	}

	// Optional port parameter
	port := 22
	if len(args) > 4 {
		portArg, err := object.AsInt(args[4])
		if err != nil {
			return err
		}
		port = int(portArg)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Use proper validation in production
		Timeout:         timeout,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	client, connectErr := ssh.Dial("tcp", address, config)
	if connectErr != nil {
		return object.NewError(fmt.Errorf("failed to connect to %s: %w", address, connectErr))
	}

	return &Client{client: client}
}

// ConnectWithKey establishes an SSH connection using private key authentication
func ConnectWithKey(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 3 {
		return object.ArgsErrorf("ssh.connect_with_key() requires at least 3 arguments: host, user, private_key")
	}

	host, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	user, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	privateKey, err := object.AsString(args[2])
	if err != nil {
		return err
	}

	// Optional timeout parameter
	timeout := 30 * time.Second
	if len(args) > 3 {
		timeoutSecs, err := object.AsInt(args[3])
		if err != nil {
			return err
		}
		timeout = time.Duration(timeoutSecs) * time.Second
	}

	// Optional port parameter
	port := 22
	if len(args) > 4 {
		portArg, err := object.AsInt(args[4])
		if err != nil {
			return err
		}
		port = int(portArg)
	}

	// Parse private key
	signer, parseErr := ssh.ParsePrivateKey([]byte(privateKey))
	if parseErr != nil {
		return object.NewError(fmt.Errorf("failed to parse private key: %w", parseErr))
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Use proper validation in production
		Timeout:         timeout,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	client, connectErr := ssh.Dial("tcp", address, config)
	if connectErr != nil {
		return object.NewError(fmt.Errorf("failed to connect to %s: %w", address, connectErr))
	}

	return &Client{client: client}
}

// Execute runs a command on the SSH connection
func Execute(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.ArgsErrorf("ssh.execute() requires 2 arguments: client, command")
	}

	client, ok := args[0].(*Client)
	if !ok {
		return object.NewError(fmt.Errorf("first argument must be an SSH client"))
	}

	command, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	// Create a session
	session, sessionErr := client.client.NewSession()
	if sessionErr != nil {
		return object.NewError(fmt.Errorf("failed to create session: %w", sessionErr))
	}
	defer session.Close()

	// Run the command
	output, outputErr := session.Output(command)
	if outputErr != nil {
		return object.NewError(fmt.Errorf("failed to execute command: %w", outputErr))
	}

	return object.NewString(string(output))
}

// NewSession creates a new SSH session
func NewSession(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.ArgsErrorf("ssh.new_session() requires 1 argument: client")
	}

	client, ok := args[0].(*Client)
	if !ok {
		return object.NewError(fmt.Errorf("first argument must be an SSH client"))
	}

	session, sessionErr := client.client.NewSession()
	if sessionErr != nil {
		return object.NewError(fmt.Errorf("failed to create session: %w", sessionErr))
	}

	return &Session{session: session}
}

// SessionRun runs a command on a session
func SessionRun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.ArgsErrorf("ssh.session_run() requires 2 arguments: session, command")
	}

	session, ok := args[0].(*Session)
	if !ok {
		return object.NewError(fmt.Errorf("first argument must be an SSH session"))
	}

	command, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	output, outputErr := session.session.Output(command)
	if outputErr != nil {
		return object.NewError(fmt.Errorf("failed to execute command: %w", outputErr))
	}

	return object.NewString(string(output))
}

// SessionClose closes a session
func SessionClose(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.ArgsErrorf("ssh.session_close() requires 1 argument: session")
	}

	session, ok := args[0].(*Session)
	if !ok {
		return object.NewError(fmt.Errorf("first argument must be an SSH session"))
	}

	closeErr := session.session.Close()
	if closeErr != nil {
		return object.NewError(fmt.Errorf("failed to close session: %w", closeErr))
	}

	return object.Nil
}

// Close closes the SSH connection
func Close(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.ArgsErrorf("ssh.close() requires 1 argument: client")
	}

	client, ok := args[0].(*Client)
	if !ok {
		return object.NewError(fmt.Errorf("first argument must be an SSH client"))
	}

	closeErr := client.client.Close()
	if closeErr != nil {
		return object.NewError(fmt.Errorf("failed to close connection: %w", closeErr))
	}

	return object.Nil
}

// Module returns the SSH module
func Module() *object.Module {
	return object.NewBuiltinsModule("ssh", map[string]object.Object{
		"connect":         object.NewBuiltin("connect", Connect),
		"connect_with_key": object.NewBuiltin("connect_with_key", ConnectWithKey),
		"execute":         object.NewBuiltin("execute", Execute),
		"new_session":     object.NewBuiltin("new_session", NewSession),
		"session_run":     object.NewBuiltin("session_run", SessionRun),
		"session_close":   object.NewBuiltin("session_close", SessionClose),
		"close":           object.NewBuiltin("close", Close),
	})
}