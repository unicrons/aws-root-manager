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
