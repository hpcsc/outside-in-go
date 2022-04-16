package storer

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io/ioutil"
)

var _ Storer = &s3Storer{}

type s3Storer struct {
	bucket string
	client *s3.Client
}

func NewS3Storer(endpoint string, bucket string) (Storer, error) {
	cfg, err := s3Config(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %v", err)
	}

	return &s3Storer{
		bucket: bucket,
		client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		}),
	}, nil
}

func (s *s3Storer) RetrieveIndividualFiles(reportType ReportType, year int, month int) ([][]byte, error) {
	prefix := fmt.Sprintf("%d/%02d/%s", year, month, reportType)
	listResponse, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects at %s: %v", prefix, err)
	}

	var result [][]byte
	for _, f := range listResponse.Contents {
		object, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    f.Key,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get object at %s: %v", *f.Key, err)
		}

		content, err := ioutil.ReadAll(object.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read object content at %s: %v", *f.Key, err)
		}

		result = append(result, content)
	}

	return result, nil
}

func (s *s3Storer) RetrieveAggregated(reportType ReportType, year int, month int) ([]byte, error) {
	key := fmt.Sprintf("%d/%02d/aggregate/%s.csv", year, month, reportType)
	getOutput, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var noSuchKeyErr *types.NoSuchKey
		if errors.As(err, &noSuchKeyErr) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get object at %s: %v", key, err)
	}

	body, err := ioutil.ReadAll(getOutput.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content at %s: %v", key, err)
	}

	return body, nil
}

func (s *s3Storer) StoreAggregated(reportType ReportType, year int, month int, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func s3Config(endpoint string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})))
}
