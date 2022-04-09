//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestReportEndpoint(t *testing.T) {
	t.Run("store aggregated file in s3 and return it", func(t *testing.T) {
		bucket := os.Getenv("BUCKET")
		previousMonth := time.Now().AddDate(0, 0, -time.Now().Day())
		s3Client := s3ClientToMockAws(t)
		newEncryptedS3Bucket(t, s3Client, bucket)
		putTestCsvAtPath(t, s3Client, bucket, fmt.Sprintf("%d/%02d/single/cluster-1.csv", previousMonth.Year(), previousMonth.Month()))
		putTestCsvAtPath(t, s3Client, bucket, fmt.Sprintf("%d/%02d/single/cluster-2.csv", previousMonth.Year(), previousMonth.Month()))

		apiUrl := os.Getenv("API_URL")
		resp, err := http.Get(fmt.Sprintf("%s/reports/single", apiUrl))
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		responseBodyIsAggregatedReport(t, resp)
		aggregatedReportExistsInS3(t, s3Client, bucket, previousMonth)
	})
}

func aggregatedReportExistsInS3(t *testing.T, s3Client *s3.Client, bucket string, previousMonth time.Time) {
	_, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%d/%02d/aggregate/single.csv", previousMonth.Year(), previousMonth.Month())),
	})
	// GetObject() returns error if object not found
	require.NoError(t, err)
}

func responseBodyIsAggregatedReport(t *testing.T, resp *http.Response) {
	reader := csv.NewReader(resp.Body)
	lines, err := reader.ReadAll()
	require.NoError(t, err)
	require.Len(t, lines, 2+1) // 2 content lines + 1 header
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

func s3ClientToMockAws(t *testing.T) *s3.Client {
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           s3Endpoint,
					SigningRegion: region,
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})))
	require.NoError(t, err)

	return s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}
