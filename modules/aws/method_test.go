//go:build aws
// +build aws

package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestMethod(t *testing.T) {

	// input := s3.PutObjectAclInput{}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	require.NoError(t, err)

	client := NewClient("s3", s3.NewFromConfig(cfg), NewConfig(cfg))

	obj, ok := client.GetAttr("put_object_acl")
	require.True(t, ok)

	method, ok := obj.(*object.Builtin)
	require.True(t, ok)

	result := method.Call(ctx, object.NewMap(map[string]object.Object{
		"AccessControlPolicy": object.NewString("my-bucket"),
	}))
	if err, ok := result.(*object.Error); ok {
		fmt.Println("ERR:", err.Error())
	}
	require.Equal(t, object.Nil, result)

}
