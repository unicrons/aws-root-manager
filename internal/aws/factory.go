package aws

import (
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
)

// IamClientFactory creates IAM clients with a given AWS config.
// This abstraction enables dependency injection of client creation logic,
// which is especially important for goroutines that create clients dynamically
// after AssumeRoot.
type IamClientFactory interface {
	NewIamClient(cfg awssdk.Config) IamClient
}

// DefaultIamClientFactory is the production implementation of IamClientFactory.
type DefaultIamClientFactory struct{}

// NewIamClient creates a new IAM client using the concrete implementation.
func (f *DefaultIamClientFactory) NewIamClient(cfg awssdk.Config) IamClient {
	return NewIamClient(cfg)
}

// S3ClientFactory creates S3 clients with a given AWS config.
type S3ClientFactory interface {
	NewS3Client(cfg awssdk.Config) S3Client
}

// DefaultS3ClientFactory is the production implementation of S3ClientFactory.
type DefaultS3ClientFactory struct{}

func (f *DefaultS3ClientFactory) NewS3Client(cfg awssdk.Config) S3Client {
	return NewS3Client(cfg)
}

// SqsClientFactory creates SQS clients with a given AWS config.
type SqsClientFactory interface {
	NewSqsClient(cfg awssdk.Config) SqsClient
}

// DefaultSqsClientFactory is the production implementation of SqsClientFactory.
type DefaultSqsClientFactory struct{}

func (f *DefaultSqsClientFactory) NewSqsClient(cfg awssdk.Config) SqsClient {
	return NewSqsClient(cfg)
}
