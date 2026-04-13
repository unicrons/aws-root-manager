package rootmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- listAccountBuckets ---

func TestListAccountBuckets_Success(t *testing.T) {
	s3 := &mockS3Client{listBucketsResult: []string{"a", "b"}}
	factory := &mockS3ClientFactory{client: s3}
	sts := &mockStsClient{}

	got, err := listAccountBuckets(context.Background(), sts, factory, "123456789012")
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, got)
}

func TestListAccountBuckets_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	factory := &mockS3ClientFactory{client: &mockS3Client{}}

	_, err := listAccountBuckets(context.Background(), sts, factory, "123456789012")
	require.Error(t, err)
	assert.ErrorIs(t, err, stsErr)
}

func TestListAccountBuckets_S3Error(t *testing.T) {
	s3 := &mockS3Client{listBucketsErr: errors.New("access denied")}
	factory := &mockS3ClientFactory{client: s3}
	sts := &mockStsClient{}

	_, err := listAccountBuckets(context.Background(), sts, factory, "123456789012")
	require.Error(t, err)
}

// --- deleteS3BucketPolicy ---

func TestDeleteS3BucketPolicy_Success(t *testing.T) {
	s3 := &mockS3Client{}
	factory := &mockS3ClientFactory{client: s3}
	sts := &mockStsClient{}

	result, err := deleteS3BucketPolicy(context.Background(), sts, factory, "123456789012", "my-bucket")
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "s3-bucket", result.ResourceType)
	assert.Equal(t, "my-bucket", result.ResourceName)
	assert.Empty(t, result.Error)
}

func TestDeleteS3BucketPolicy_STSError(t *testing.T) {
	sts := &mockStsClient{assumeRootErr: errors.New("assume root denied")}
	factory := &mockS3ClientFactory{client: &mockS3Client{}}

	result, err := deleteS3BucketPolicy(context.Background(), sts, factory, "123456789012", "my-bucket")
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}

func TestDeleteS3BucketPolicy_S3Error(t *testing.T) {
	s3 := &mockS3Client{deleteBucketPolErr: errors.New("policy not found")}
	factory := &mockS3ClientFactory{client: s3}
	sts := &mockStsClient{}

	result, err := deleteS3BucketPolicy(context.Background(), sts, factory, "123456789012", "my-bucket")
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}

// --- listAccountQueues ---

func TestListAccountQueues_Success(t *testing.T) {
	sqs := &mockSqsClient{listQueuesResult: []string{"https://sqs/q1", "https://sqs/q2"}}
	factory := &mockSqsClientFactory{client: sqs}
	sts := &mockStsClient{}

	got, err := listAccountQueues(context.Background(), sts, factory, "123456789012")
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestListAccountQueues_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	factory := &mockSqsClientFactory{client: &mockSqsClient{}}

	_, err := listAccountQueues(context.Background(), sts, factory, "123456789012")
	require.Error(t, err)
	assert.ErrorIs(t, err, stsErr)
}

// --- deleteSQSQueuePolicy ---

func TestDeleteSQSQueuePolicy_Success(t *testing.T) {
	sqs := &mockSqsClient{}
	factory := &mockSqsClientFactory{client: sqs}
	sts := &mockStsClient{}

	result, err := deleteSQSQueuePolicy(context.Background(), sts, factory, "123456789012", "https://sqs/q1")
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "sqs-queue", result.ResourceType)
	assert.Equal(t, "https://sqs/q1", result.ResourceName)
}

func TestDeleteSQSQueuePolicy_STSError(t *testing.T) {
	sts := &mockStsClient{assumeRootErr: errors.New("assume root denied")}
	factory := &mockSqsClientFactory{client: &mockSqsClient{}}

	result, err := deleteSQSQueuePolicy(context.Background(), sts, factory, "123456789012", "https://sqs/q1")
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}

func TestDeleteSQSQueuePolicy_SqsError(t *testing.T) {
	sqs := &mockSqsClient{deleteQueuePolErr: errors.New("access denied")}
	factory := &mockSqsClientFactory{client: sqs}
	sts := &mockStsClient{}

	result, err := deleteSQSQueuePolicy(context.Background(), sts, factory, "123456789012", "https://sqs/q1")
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
}
