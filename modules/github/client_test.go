package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v73/github"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	client := github.NewClient(nil)
	ghClient := New(client)

	require.NotNil(t, ghClient)
	require.IsType(t, &Client{}, ghClient)
	require.Equal(t, CLIENT, ghClient.Type())
}

func TestCreateUnauthenticated(t *testing.T) {
	ctx := context.Background()
	result := Create(ctx)

	require.NotNil(t, result)
	require.IsType(t, &Client{}, result)

	require.Equal(t, CLIENT, result.(*Client).Type())
}

func TestCreateAuthenticated(t *testing.T) {
	ctx := context.Background()
	token := object.NewString("test-token")
	result := Create(ctx, token)

	require.NotNil(t, result)
	require.IsType(t, &Client{}, result)
	require.Equal(t, CLIENT, result.(*Client).Type())
}

func TestModule(t *testing.T) {
	module := Module()

	require.NotNil(t, module)
	require.IsType(t, &object.Module{}, module)

	nameObj := module.Name()
	name, err := object.AsString(nameObj)
	require.NoError(t, err)
	require.Equal(t, "github", name)
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

	require.NotNil(t, result)
	require.IsType(t, &object.Map{}, result)

	mapObj := result.(*object.Map)
	nameObj := mapObj.Get("name")
	require.NotNil(t, nameObj)
	require.IsType(t, &object.String{}, nameObj)
	require.Equal(t, "test", nameObj.(*object.String).Value())
}

func TestGetRepository(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	ghClient := New(client)

	result := ghClient.GetRepo(ctx, object.NewString("octocat"), object.NewString("Hello-World"))

	require.NotNil(t, result)
	require.IsType(t, &object.Map{}, result)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	ghClient := New(client)

	result := ghClient.GetUser(ctx, object.NewString("octocat"))

	require.NotNil(t, result)
	require.IsType(t, &object.Map{}, result)

	m := result.(*object.Map)
	nameObj := m.Get("name")
	require.NotNil(t, nameObj)
	require.IsType(t, &object.String{}, nameObj)
	require.Equal(t, "The Octocat", nameObj.(*object.String).Value())
}
