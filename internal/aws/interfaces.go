package aws

import (
	"context"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
)

// IamClient defines the interface for IAM operations.
// This interface enables mocking and dependency injection for testing.
type IamClient interface {
	// CheckOrganizationRootAccess verifies if AWS centralized root access is enabled
	CheckOrganizationRootAccess(ctx context.Context, rootSessionsRequired bool) error

	// GetLoginProfile checks if an account has root login profile enabled
	GetLoginProfile(ctx context.Context, accountId string) (bool, error)

	// DeleteLoginProfile deletes root login profile for a specific account
	DeleteLoginProfile(ctx context.Context, accountId string) error

	// ListAccessKeys gets a list of root access keys for a specific account
	ListAccessKeys(ctx context.Context, accountId string) ([]string, error)

	// DeleteAccessKeys deletes a list of root access keys for a specific account
	DeleteAccessKeys(ctx context.Context, accountId string, accessKeyIds []string) error

	// ListMFADevices gets a list of root MFA devices for a specific account
	ListMFADevices(ctx context.Context, accountId string) ([]string, error)

	// DeactivateMFADevices deactivates a list of root MFA devices for a specific account
	DeactivateMFADevices(ctx context.Context, accountId string, mfaSerialNumbers []string) error

	// ListSigningCertificates gets a list of root signing certificates for a specific account
	ListSigningCertificates(ctx context.Context, accountId string) ([]string, error)

	// DeleteSigningCertificates deletes a list of root signing certificates for a specific account
	DeleteSigningCertificates(ctx context.Context, accountId string, certificates []string) error

	// EnableOrganizationsRootCredentialsManagement enables centralized root credentials management
	EnableOrganizationsRootCredentialsManagement(ctx context.Context) error

	// EnableOrganizationsRootSessions enables centralized root sessions
	EnableOrganizationsRootSessions(ctx context.Context) error

	// CreateLoginProfile allows root password recovery
	CreateLoginProfile(ctx context.Context) error
}

// StsClient defines the interface for STS operations.
// This interface enables mocking and dependency injection for testing.
type StsClient interface {
	// GetAssumeRootConfig gets AWS config with assumed root credentials for a specific account and task
	GetAssumeRootConfig(ctx context.Context, accountId, taskPolicyName string) (awssdk.Config, error)
}

// S3Client defines the interface for S3 operations scoped to a single account.
// This interface enables mocking and dependency injection for testing.
type S3Client interface {
	// ListBuckets returns the names of all buckets owned by the caller.
	ListBuckets(ctx context.Context) ([]string, error)
	// GetBucketPolicy returns the bucket policy JSON, or empty string if none exists.
	GetBucketPolicy(ctx context.Context, bucketName string) (string, error)
	// DeleteBucketPolicy deletes the bucket policy attached to the given bucket.
	DeleteBucketPolicy(ctx context.Context, bucketName string) error
}

// SqsClient defines the interface for SQS operations scoped to a single account.
// This interface enables mocking and dependency injection for testing.
type SqsClient interface {
	// ListQueues returns the URLs of all queues owned by the caller.
	ListQueues(ctx context.Context) ([]string, error)
	// GetQueuePolicy returns the queue policy JSON, or empty string if none exists.
	GetQueuePolicy(ctx context.Context, queueUrl string) (string, error)
	// DeleteQueuePolicy clears the access policy attached to the given queue URL.
	DeleteQueuePolicy(ctx context.Context, queueUrl string) error
}

// OrganizationsClient defines the interface for AWS Organizations operations.
// This interface enables mocking and dependency injection for testing.
type OrganizationsClient interface {
	// DescribeOrganization returns the management account ID of the organization
	DescribeOrganization(ctx context.Context) (string, error)

	// ListAccounts returns all accounts in the organization
	ListAccounts(ctx context.Context) ([]OrganizationAccount, error)

	// EnableAWSServiceAccess enables AWS service access for the organization
	EnableAWSServiceAccess(ctx context.Context, service string) error
}
