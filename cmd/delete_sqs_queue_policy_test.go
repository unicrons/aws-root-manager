package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestDeleteSQSQueuePolicyCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		deleteQueueResult: rootmanager.PolicyDeletionResult{
			AccountId:    "123456789012",
			ResourceType: "sqs-queue",
			ResourceName: "https://sqs/q1",
			Success:      true,
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012", "--queue", "https://sqs/q1"})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, buf.String(), "https://sqs/q1")
}

func TestDeleteSQSQueuePolicyCommand_WithAccount(t *testing.T) {
	mock := &mockRootManager{
		deleteQueueResult: rootmanager.PolicyDeletionResult{
			AccountId:    "123456789012",
			ResourceType: "sqs-queue",
			ResourceName: "https://sqs/q1",
			Success:      true,
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012", "--queue", "https://sqs/q1"})

	require.NoError(t, cmd.Execute())
}

func TestDeleteSQSQueuePolicyCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Delete(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012", "--queue", "https://sqs/q1"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestDeleteSQSQueuePolicyCommand_DeletionFailure(t *testing.T) {
	mock := &mockRootManager{
		deleteQueueResult: rootmanager.PolicyDeletionResult{
			AccountId:    "123456789012",
			ResourceType: "sqs-queue",
			ResourceName: "https://sqs/q1",
			Success:      false,
			Error:        "access denied",
		},
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012", "--queue", "https://sqs/q1"})

	require.Error(t, cmd.Execute())
}

func TestDeleteSQSQueuePolicyCommand_NoQueuesFoundInTUI(t *testing.T) {
	mock := &mockRootManager{
		listQueuesResult: []string{},
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no queues found")
}

func TestDeleteSQSQueuePolicyCommand_ListQueuesError(t *testing.T) {
	mock := &mockRootManager{
		listQueuesErr: errors.New("assume root denied"),
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012"})

	require.Error(t, cmd.Execute())
}
