//go:build awstests
// +build awstests

package s3fs

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func TestS3Filesystem(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile(os.Getenv("RISOR_TEST_AWS_PROFILE")),
	)
	require.Nil(t, err)

	s3Client := s3.NewFromConfig(cfg)

	fs, err := New(ctx,
		WithBase("/some/base/"),
		WithClient(s3Client),
		WithBucket(os.Getenv("RISOR_TEST_S3_BUCKET")),
	)
	require.Nil(t, err)

	f, err := fs.Create("test.txt")
	require.Nil(t, err)

	n, err := f.Write([]byte("hello world"))
	require.Nil(t, err)
	require.Equal(t, 11, n)

	err = f.Close()
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	require.Nil(t, err, errStr)
}

func TestS3Mkdir(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile(os.Getenv("RISOR_TEST_AWS_PROFILE")),
	)
	require.Nil(t, err)

	s3Client := s3.NewFromConfig(cfg)

	fs, err := New(ctx,
		WithBase("/foo"),
		WithClient(s3Client),
		WithBucket(os.Getenv("RISOR_TEST_S3_BUCKET")),
	)
	require.Nil(t, err)

	require.Nil(t, fs.Mkdir("bar", 0o755))
	require.Nil(t, fs.MkdirAll("1/2/3", 0o755))
}

func TestS3WriteFileAndStat(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile(os.Getenv("RISOR_TEST_AWS_PROFILE")),
	)
	require.Nil(t, err)

	s3Client := s3.NewFromConfig(cfg)

	fs, err := New(ctx,
		WithBase("/stat"),
		WithClient(s3Client),
		WithBucket(os.Getenv("RISOR_TEST_S3_BUCKET")),
	)
	require.Nil(t, err)

	require.Nil(t, fs.WriteFile("stat.txt", []byte("yup!"), 0o644))

	stat, err := fs.Stat("stat.txt")
	require.Nil(t, err)
	require.Equal(t, "stat.txt", stat.Name())
	require.Equal(t, int64(4), stat.Size())
	require.Equal(t, false, stat.IsDir())
}
