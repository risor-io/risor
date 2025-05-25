package redis

import (
	"context"
	"os"
	"testing"

	"github.com/risor-io/risor/object"
)

var client object.Object

func TestMain(m *testing.M) {
	redisURL := object.NewString("redis://localhost:6379")
	client = Create(context.Background(), redisURL)
	code := m.Run()
	os.Exit(code)
}

func TestRedisModule(t *testing.T) {
	mod := Module()
	if mod.Type() != object.MODULE {
		t.Errorf("expected module type to be MODULe, got %s", mod.Type())
	}

	clientObj, ok := mod.GetAttr("client")
	if !ok {
		t.Fatal("expected 'client' function to exist in module")
	}

	builtinClient, ok := clientObj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'client' to be a builtin function, got %T", clientObj)
	}
	if builtinClient.Name() != "client" {
		t.Errorf("expected function name to be 'client', got %s", builtinClient.Name())
	}
}

func TestRedisCreate(t *testing.T) {
	if client.Type() != REDIS {
		t.Errorf("expected Client client type to be REDIS, got %s", client.Type())
	}
}
