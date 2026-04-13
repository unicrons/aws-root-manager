package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Client struct {
	client *s3.Client
}

func NewS3Client(awscfg aws.Config) S3Client {
	return &s3Client{client: s3.NewFromConfig(awscfg)}
}

func (c *s3Client) ListBuckets(ctx context.Context) ([]string, error) {
	slog.Debug("listing s3 buckets")

	output, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("error listing s3 buckets: %w", err)
	}

	buckets := make([]string, 0, len(output.Buckets))
	for _, b := range output.Buckets {
		buckets = append(buckets, aws.ToString(b.Name))
	}
	return buckets, nil
}

func (c *s3Client) DeleteBucketPolicy(ctx context.Context, bucketName string) error {
	slog.Debug("deleting s3 bucket policy", "bucket", bucketName)

	_, err := c.client.DeleteBucketPolicy(ctx, &s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("error deleting bucket policy for bucket %s: %w", bucketName, err)
	}
	return nil
}
