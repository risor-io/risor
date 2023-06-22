//go:build aws
// +build aws

package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func TestS3Client(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile("curtis-cmds-dev"))
	require.Nil(t, err)

	cl := s3.NewFromConfig(cfg)
	require.NotNil(t, cl)

	clObj := NewClient("s3", cl)
	fmt.Println(clObj)

	require.True(t, false)
}
