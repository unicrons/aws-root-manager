package rootmanager

import (
	"context"
	"log/slog"

	"github.com/unicrons/aws-root-manager/internal/aws"
)

const (
	s3UnlockTaskPolicy  = "S3UnlockBucketPolicy"
	sqsUnlockTaskPolicy = "SQSUnlockQueuePolicy"

	resourceTypeS3Bucket = "s3-bucket"
	resourceTypeSqsQueue = "sqs-queue"
)

func listAccountBuckets(ctx context.Context, sts aws.StsClient, factory aws.S3ClientFactory, accountId string) ([]string, error) {
	slog.Debug("listing account buckets", "account_id", accountId)

	cfg, err := sts.GetAssumeRootConfig(ctx, accountId, s3UnlockTaskPolicy)
	if err != nil {
		return nil, err
	}
	return factory.NewS3Client(cfg).ListBuckets(ctx)
}

func deleteS3BucketPolicy(ctx context.Context, sts aws.StsClient, factory aws.S3ClientFactory, accountId, bucketName string) (PolicyDeletionResult, error) {
	slog.Debug("deleting s3 bucket policy", "account_id", accountId, "bucket", bucketName)

	result := PolicyDeletionResult{
		AccountId:    accountId,
		ResourceType: resourceTypeS3Bucket,
		ResourceName: bucketName,
	}

	cfg, err := sts.GetAssumeRootConfig(ctx, accountId, s3UnlockTaskPolicy)
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	if err := factory.NewS3Client(cfg).DeleteBucketPolicy(ctx, bucketName); err != nil {
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	return result, nil
}

func listAccountQueues(ctx context.Context, sts aws.StsClient, factory aws.SqsClientFactory, accountId string) ([]string, error) {
	slog.Debug("listing account queues", "account_id", accountId)

	cfg, err := sts.GetAssumeRootConfig(ctx, accountId, sqsUnlockTaskPolicy)
	if err != nil {
		return nil, err
	}
	return factory.NewSqsClient(cfg).ListQueues(ctx)
}

func deleteSQSQueuePolicy(ctx context.Context, sts aws.StsClient, factory aws.SqsClientFactory, accountId, queueUrl string) (PolicyDeletionResult, error) {
	slog.Debug("deleting sqs queue policy", "account_id", accountId, "queue_url", queueUrl)

	result := PolicyDeletionResult{
		AccountId:    accountId,
		ResourceType: resourceTypeSqsQueue,
		ResourceName: queueUrl,
	}

	cfg, err := sts.GetAssumeRootConfig(ctx, accountId, sqsUnlockTaskPolicy)
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	if err := factory.NewSqsClient(cfg).DeleteQueuePolicy(ctx, queueUrl); err != nil {
		result.Error = err.Error()
		return result, nil
	}

	result.Success = true
	return result, nil
}
