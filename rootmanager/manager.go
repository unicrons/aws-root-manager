package rootmanager

import (
	"context"
	"errors"
	"time"

	"github.com/unicrons/aws-root-manager/internal/aws"
)

// manager implements RootManager using AWS clients.
type manager struct {
	iam        aws.IamClient
	sts        aws.StsClient
	org        aws.OrganizationsClient
	factory    aws.IamClientFactory
	s3Factory  aws.S3ClientFactory
	sqsFactory aws.SqsClientFactory
}

// newManager returns a RootManager that uses the given AWS clients and factories.
// sts and org may be nil for callers that only use CheckRootAccess.
func newManager(iam aws.IamClient, sts aws.StsClient, org aws.OrganizationsClient, factory aws.IamClientFactory, s3Factory aws.S3ClientFactory, sqsFactory aws.SqsClientFactory) RootManager {
	return &manager{iam: iam, sts: sts, org: org, factory: factory, s3Factory: s3Factory, sqsFactory: sqsFactory}
}

// NewRootManager returns a RootManager configured from the default AWS environment.
// It loads credentials from the standard AWS credential chain (env vars, ~/.aws, IAM role).
func NewRootManager(ctx context.Context) (RootManager, error) {
	cfg, err := aws.LoadAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	// AssumeRoot is subject to strict AWS rate limits, so the STS client uses
	// a more aggressive retry policy than the default.
	stsCfg, err := aws.LoadAWSConfig(ctx, aws.WithRetry(10, 30*time.Second))
	if err != nil {
		return nil, err
	}
	return newManager(
		aws.NewIamClient(cfg),
		aws.NewStsClient(stsCfg),
		aws.NewOrganizationsClient(cfg),
		&aws.DefaultIamClientFactory{},
		&aws.DefaultS3ClientFactory{},
		&aws.DefaultSqsClientFactory{},
	), nil
}

func (m *manager) GetS3BucketPolicy(ctx context.Context, accountId, bucketName string) (string, error) {
	if m.sts == nil {
		return "", errors.New("STS client required for get")
	}
	return getS3BucketPolicy(ctx, m.sts, m.s3Factory, accountId, bucketName)
}

func (m *manager) ListAccountBuckets(ctx context.Context, accountId string) ([]string, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for listing buckets")
	}
	return listAccountBuckets(ctx, m.sts, m.s3Factory, accountId)
}

func (m *manager) DeleteS3BucketPolicy(ctx context.Context, accountId, bucketName string) (PolicyDeletionResult, error) {
	if m.sts == nil {
		return PolicyDeletionResult{}, errors.New("STS client required for delete")
	}
	return deleteS3BucketPolicy(ctx, m.sts, m.s3Factory, accountId, bucketName)
}

func (m *manager) GetSQSQueuePolicy(ctx context.Context, accountId, queueUrl string) (string, error) {
	if m.sts == nil {
		return "", errors.New("STS client required for get")
	}
	return getSQSQueuePolicy(ctx, m.sts, m.sqsFactory, accountId, queueUrl)
}

func (m *manager) ListAccountQueues(ctx context.Context, accountId string) ([]string, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for listing queues")
	}
	return listAccountQueues(ctx, m.sts, m.sqsFactory, accountId)
}

func (m *manager) DeleteSQSQueuePolicy(ctx context.Context, accountId, queueUrl string) (PolicyDeletionResult, error) {
	if m.sts == nil {
		return PolicyDeletionResult{}, errors.New("STS client required for delete")
	}
	return deleteSQSQueuePolicy(ctx, m.sts, m.sqsFactory, accountId, queueUrl)
}

func (m *manager) AuditAccounts(ctx context.Context, accountIds []string) ([]RootCredentials, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for audit")
	}
	return auditAccounts(ctx, m.iam, m.sts, m.factory, accountIds)
}

func (m *manager) CheckRootAccess(ctx context.Context) (RootAccessStatus, error) {
	return checkRootAccess(ctx, m.iam)
}

func (m *manager) EnableRootAccess(ctx context.Context, enableSessions bool) (RootAccessStatus, RootAccessStatus, error) {
	if m.org == nil {
		return RootAccessStatus{}, RootAccessStatus{}, errors.New("Organizations client required for enable")
	}
	return enableRootAccess(ctx, m.iam, m.org, enableSessions)
}

func (m *manager) DeleteCredentials(ctx context.Context, creds []RootCredentials, credentialType string) ([]DeletionResult, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for delete")
	}
	return deleteAccountsCredentials(ctx, m.iam, m.sts, m.factory, creds, credentialType)
}

func (m *manager) RecoverRootPassword(ctx context.Context, accountIds []string) ([]RecoveryResult, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for recovery")
	}
	return recoverAccountsRootPassword(ctx, m.iam, m.sts, m.factory, accountIds)
}
