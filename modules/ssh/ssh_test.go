package ssh

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	// Test missing arguments
	result := Connect(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires at least 3 arguments")

	// Test invalid host type
	result = Connect(context.Background(), object.NewInt(123), object.NewString("user"), object.NewString("pass"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")

	// Test invalid user type
	result = Connect(context.Background(), object.NewString("host"), object.NewInt(123), object.NewString("pass"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")

	// Test invalid password type
	result = Connect(context.Background(), object.NewString("host"), object.NewString("user"), object.NewInt(123))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")

	// Test invalid timeout type
	result = Connect(context.Background(), object.NewString("host"), object.NewString("user"), object.NewString("pass"), object.NewString("invalid"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected an int")

	// Test invalid port type
	result = Connect(context.Background(), object.NewString("host"), object.NewString("user"), object.NewString("pass"), object.NewInt(30), object.NewString("invalid"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected an int")
}

func TestConnectWithKey(t *testing.T) {
	// Test missing arguments
	result := ConnectWithKey(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires at least 3 arguments")

	// Test invalid host type
	result = ConnectWithKey(context.Background(), object.NewInt(123), object.NewString("user"), object.NewString("key"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")

	// Test invalid user type
	result = ConnectWithKey(context.Background(), object.NewString("host"), object.NewInt(123), object.NewString("key"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")

	// Test invalid private key type
	result = ConnectWithKey(context.Background(), object.NewString("host"), object.NewString("user"), object.NewInt(123))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")

	// Test invalid private key content
	result = ConnectWithKey(context.Background(), object.NewString("host"), object.NewString("user"), object.NewString("invalid-key"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "failed to parse private key")
}

func TestExecute(t *testing.T) {
	// Test missing arguments
	result := Execute(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires 2 arguments")

	// Test invalid client type
	result = Execute(context.Background(), object.NewString("not-a-client"), object.NewString("command"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "first argument must be an SSH client")

	// Test invalid command type
	mockClient := &Client{}
	result = Execute(context.Background(), mockClient, object.NewInt(123))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")
}

func TestNewSession(t *testing.T) {
	// Test missing arguments
	result := NewSession(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires 1 argument")

	// Test invalid client type
	result = NewSession(context.Background(), object.NewString("not-a-client"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "first argument must be an SSH client")
}

func TestSessionRun(t *testing.T) {
	// Test missing arguments
	result := SessionRun(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires 2 arguments")

	// Test invalid session type
	result = SessionRun(context.Background(), object.NewString("not-a-session"), object.NewString("command"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "first argument must be an SSH session")

	// Test invalid command type
	mockSession := &Session{}
	result = SessionRun(context.Background(), mockSession, object.NewInt(123))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "expected a string")
}

func TestSessionClose(t *testing.T) {
	// Test missing arguments
	result := SessionClose(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires 1 argument")

	// Test invalid session type
	result = SessionClose(context.Background(), object.NewString("not-a-session"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "first argument must be an SSH session")
}

func TestClose(t *testing.T) {
	// Test missing arguments
	result := Close(context.Background())
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "requires 1 argument")

	// Test invalid client type
	result = Close(context.Background(), object.NewString("not-a-client"))
	require.IsType(t, &object.Error{}, result)
	require.Contains(t, result.Inspect(), "first argument must be an SSH client")
}

func TestModule(t *testing.T) {
	mod := Module()
	require.NotNil(t, mod)
	require.Equal(t, "ssh", mod.Name().Value())
	
	// Test that all expected functions are present
	functions := []string{"connect", "connect_with_key", "execute", "new_session", "session_run", "session_close", "close"}
	for _, fn := range functions {
		attr, found := mod.GetAttr(fn)
		require.True(t, found, "Function %s not found in module", fn)
		require.NotNil(t, attr, "Function %s is nil", fn)
	}
}