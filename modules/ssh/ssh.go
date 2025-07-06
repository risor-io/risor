package ssh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

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
	switch name {
	case "execute":
		return object.NewBuiltin("execute", c.executeMethod), true
	case "close":
		return object.NewBuiltin("close", c.closeMethod), true
	default:
		return nil, false
	}
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

// executeMethod is the method implementation for client.execute()
func (c *Client) executeMethod(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.ArgsErrorf("execute() requires 1 argument: command")
	}

	command, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	// Create a session
	session, sessionErr := c.client.NewSession()
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

// closeMethod is the method implementation for client.close()
func (c *Client) closeMethod(ctx context.Context, args ...object.Object) object.Object {
	closeErr := c.client.Close()
	if closeErr != nil {
		return object.NewError(fmt.Errorf("failed to close connection: %w", closeErr))
	}

	return object.Nil
}

// createHostKeyCallback creates a HostKeyCallback using the system's known_hosts file
func createHostKeyCallback() (ssh.HostKeyCallback, error) {
	// Use system known_hosts file
	knownHostsPath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")

	// Create the callback
	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, err
	}

	return hostKeyCallback, nil
}

// Connect establishes an SSH connection with consolidated authentication methods
func Connect(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 3 {
		return object.ArgsErrorf("ssh.connect() requires 3 arguments: user, host, options")
	}

	user, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	host, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	optionsObj, ok := args[2].(*object.Map)
	if !ok {
		return object.NewError(fmt.Errorf("third argument must be a map of options"))
	}

	// Parse options
	options := optionsObj.Value()

	// Default values
	timeout := 30 * time.Second
	port := 22
	insecure := false

	// Parse timeout if provided
	if timeoutObj, exists := options["timeout"]; exists {
		timeoutSecs, err := object.AsInt(timeoutObj)
		if err != nil {
			return err
		}
		timeout = time.Duration(timeoutSecs) * time.Second
	}

	// Parse port if provided
	if portObj, exists := options["port"]; exists {
		portArg, err := object.AsInt(portObj)
		if err != nil {
			return err
		}
		port = int(portArg)
	}

	// Parse insecure flag if provided
	if insecureObj, exists := options["insecure"]; exists {
		insecure = insecureObj.IsTruthy()
	}

	// Set up authentication methods
	var authMethods []ssh.AuthMethod

	// Check for password authentication
	if passwordObj, exists := options["password"]; exists {
		password, err := object.AsString(passwordObj)
		if err != nil {
			return err
		}
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Check for private key authentication
	if privateKeyObj, exists := options["private_key"]; exists {
		privateKey, err := object.AsString(privateKeyObj)
		if err != nil {
			return err
		}

		// Parse private key
		signer, parseErr := ssh.ParsePrivateKey([]byte(privateKey))
		if parseErr != nil {
			return object.NewError(fmt.Errorf("failed to parse private key: %w", parseErr))
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Ensure at least one authentication method is provided
	if len(authMethods) == 0 {
		return object.NewError(fmt.Errorf("at least one authentication method (password or private_key) must be provided"))
	}

	// Set up host key callback
	var hostKeyCallback ssh.HostKeyCallback
	if insecure {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		var err error
		hostKeyCallback, err = createHostKeyCallback()
		if err != nil {
			return object.NewError(fmt.Errorf("failed to create host key callback: %w", err))
		}
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         timeout,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	client, connectErr := ssh.Dial("tcp", address, config)
	if connectErr != nil {
		return object.NewError(fmt.Errorf("failed to connect to %s: %w", address, connectErr))
	}

	return &Client{client: client}
}

// Module returns the SSH module
func Module() *object.Module {
	return object.NewBuiltinsModule("ssh", map[string]object.Object{
		"connect": object.NewBuiltin("connect", Connect),
	}, Connect)
}
