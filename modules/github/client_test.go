package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v73/github"
	"github.com/risor-io/risor/object"
)

func TestNew(t *testing.T) {
	client := github.NewClient(nil)
	ghClient := New(client)

	if ghClient == nil {
		t.Fatal("Expected non-nil client")
	}

	if ghClient.Type() != CLIENT {
		t.Errorf("Expected type %s, got %s", CLIENT, ghClient.Type())
	}
}

func TestCreateUnauthenticated(t *testing.T) {
	ctx := context.Background()
	result := Create(ctx)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	client, ok := result.(*Client)
	if !ok {
		t.Fatalf("Expected *Client, got %T", result)
	}

	if client.Type() != CLIENT {
		t.Errorf("Expected type %s, got %s", CLIENT, client.Type())
	}
}

func TestCreateAuthenticated(t *testing.T) {
	ctx := context.Background()
	token := object.NewString("test-token")
	result := Create(ctx, token)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	client, ok := result.(*Client)
	if !ok {
		t.Fatalf("Expected *Client, got %T", result)
	}

	if client.Type() != CLIENT {
		t.Errorf("Expected type %s, got %s", CLIENT, client.Type())
	}
}

func TestModule(t *testing.T) {
	module := Module()

	if module == nil {
		t.Fatal("Expected non-nil module")
	}

	nameObj := module.Name()
	name, err := object.AsString(nameObj)
	if err != nil {
		t.Fatalf("Expected string name, got error: %v", err)
	}
	if name != "github" {
		t.Errorf("Expected module name 'github', got %s", name)
	}
}

func TestAsMap(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	test := TestStruct{
		Name:  "test",
		Count: 42,
	}

	result := asMap(test)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	mapObj, ok := result.(*object.Map)
	if !ok {
		t.Fatalf("Expected *object.Map, got %T", result)
	}

	nameObj := mapObj.Get("name")
	if nameObj == nil {
		t.Fatal("Expected name field")
	}

	name, err := object.AsString(nameObj)
	if err != nil {
		t.Fatalf("Expected string, got error: %v", err)
	}

	if name != "test" {
		t.Errorf("Expected 'test', got %s", name)
	}
}
