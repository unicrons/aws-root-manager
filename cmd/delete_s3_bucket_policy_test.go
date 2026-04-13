package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestDeleteS3BucketPolicyCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		deleteBucketResult: rootmanager.PolicyDeletionResult{
			AccountId:    "123456789012",
			ResourceType: "s3-bucket",
			ResourceName: "my-bucket",
			Success:      true,
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"s3-bucket-policy", "--account", "123456789012", "--bucket", "my-bucket"})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, buf.String(), "my-bucket")
}

func TestDeleteS3BucketPolicyCommand_WithAccount(t *testing.T) {
	mock := &mockRootManager{
		deleteBucketResult: rootmanager.PolicyDeletionResult{
			AccountId:    "123456789012",
			ResourceType: "s3-bucket",
			ResourceName: "my-bucket",
			Success:      true,
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"s3-bucket-policy", "--account", "123456789012", "--bucket", "my-bucket"})

	require.NoError(t, cmd.Execute())
}

func TestDeleteS3BucketPolicyCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Delete(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"s3-bucket-policy", "--account", "123456789012", "--bucket", "my-bucket"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestDeleteS3BucketPolicyCommand_DeletionFailure(t *testing.T) {
	mock := &mockRootManager{
		deleteBucketResult: rootmanager.PolicyDeletionResult{
			AccountId:    "123456789012",
			ResourceType: "s3-bucket",
			ResourceName: "my-bucket",
			Success:      false,
			Error:        "access denied",
		},
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"s3-bucket-policy", "--account", "123456789012", "--bucket", "my-bucket"})

	require.Error(t, cmd.Execute())
}

func TestDeleteS3BucketPolicyCommand_NoBucketsFoundInTUI(t *testing.T) {
	mock := &mockRootManager{
		listBucketsResult: []string{},
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"s3-bucket-policy", "--account", "123456789012"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no buckets found")
}

func TestDeleteS3BucketPolicyCommand_ListBucketsError(t *testing.T) {
	mock := &mockRootManager{
		listBucketsErr: errors.New("assume root denied"),
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"s3-bucket-policy", "--account", "123456789012"})

	require.Error(t, cmd.Execute())
}
