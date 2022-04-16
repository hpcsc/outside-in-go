//go:build integration

package storer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestS3Storer_RetrieveIndividualFiles(t *testing.T) {
	t.Run("return empty and no error when no files found", func(t *testing.T) {
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		bucket := randomName("bucket")

		client := s3ClientToMockAws(t, s3Endpoint)
		newEncryptedS3Bucket(t, client, bucket)

		s, err := NewS3Storer(s3Endpoint, bucket)
		require.NoError(t, err)
		data, err := s.RetrieveIndividualFiles(SingleReportType, 2022, 4)

		require.NoError(t, err)
		require.Empty(t, data)
	})

	t.Run("return files when found", func(t *testing.T) {
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		bucket := randomName("bucket")

		client := s3ClientToMockAws(t, s3Endpoint)
		newEncryptedS3Bucket(t, client, bucket)
		putTestCsvAtPath(t, client, bucket, "2022/04/single/cluster-1.csv")
		putTestCsvAtPath(t, client, bucket, "2022/04/single/cluster-2.csv")

		s, err := NewS3Storer(s3Endpoint, bucket)
		require.NoError(t, err)
		data, err := s.RetrieveIndividualFiles(SingleReportType, 2022, 4)

		require.NoError(t, err)
		require.Len(t, data, 2)
		expected := [][]byte{
			[]byte(fmt.Sprintf("BUCKET,PATH\n%s,2022/04/single/cluster-1.csv", bucket)),
			[]byte(fmt.Sprintf("BUCKET,PATH\n%s,2022/04/single/cluster-2.csv", bucket)),
		}
		require.Equal(t, expected, data)
	})
}

func TestS3Storer_RetrieveAggregated(t *testing.T) {
	t.Run("return nil and no error when no aggregated file found", func(t *testing.T) {
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		bucket := randomName("bucket")

		client := s3ClientToMockAws(t, s3Endpoint)
		newEncryptedS3Bucket(t, client, bucket)

		s, err := NewS3Storer(s3Endpoint, bucket)
		require.NoError(t, err)

		data, err := s.RetrieveAggregated(SingleReportType, 2022, 4)

		require.NoError(t, err)
		require.Nil(t, data)
	})

	t.Run("return error for any other error than file not found", func(t *testing.T) {
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		bucket := randomName("bucket")

		s, err := NewS3Storer(s3Endpoint, bucket)
		require.NoError(t, err)

		_, err = s.RetrieveAggregated(SingleReportType, 2022, 4)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get object at")
	})

	t.Run("return aggregated file when found", func(t *testing.T) {
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		bucket := randomName("bucket")

		client := s3ClientToMockAws(t, s3Endpoint)
		newEncryptedS3Bucket(t, client, bucket)
		putTestCsvAtPath(t, client, bucket, "2022/04/aggregate/single.csv")

		s, err := NewS3Storer(s3Endpoint, bucket)
		require.NoError(t, err)

		data, err := s.RetrieveAggregated(SingleReportType, 2022, 4)

		require.NoError(t, err)
		expected := fmt.Sprintf("BUCKET,PATH\n%s,2022/04/aggregate/single.csv", bucket)
		require.Equal(t, []byte(expected), data)
	})
}

func putTestCsvAtPath(t *testing.T, s3Client *s3.Client, bucket string, path string) {
	csv := fmt.Sprintf("BUCKET,PATH\n%s,%s", bucket, path)
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		ContentType: aws.String("text/csv"),
		Body:        bytes.NewReader([]byte(csv)),
	})
	require.NoError(t, err)
}

func newEncryptedS3Bucket(t *testing.T, s3Client *s3.Client, bucket string) {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    types.BucketCannedACLPrivate,
	})
	require.NoError(t, err)

	_, err = s3Client.PutBucketEncryption(context.TODO(), &s3.PutBucketEncryptionInput{
		Bucket: aws.String(bucket),
		ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
			Rules: []types.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
						SSEAlgorithm: types.ServerSideEncryptionAes256,
					},
				},
			},
		},
	})
	require.NoError(t, err)
}

func randomName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func s3ClientToMockAws(t *testing.T, endpoint string) *s3.Client {
	awsConfig, err := s3Config(endpoint)
	require.NoError(t, err)

	return s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}
