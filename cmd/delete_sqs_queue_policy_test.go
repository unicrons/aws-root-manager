package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// When getQueuePolicyResult is empty, the command prints "No queue policy found." and exits
// without invoking the TUI confirmation — safe to use in tests.

func TestDeleteSQSQueuePolicyCommand_NoPolicyFound(t *testing.T) {
	mock := &mockRootManager{
		getQueuePolicyResult: "",
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012", "--queue", "https://sqs/q1"})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, buf.String(), "No queue policy found.")
}

func TestDeleteSQSQueuePolicyCommand_GetPolicyError(t *testing.T) {
	mock := &mockRootManager{
		getQueuePolicyErr: errors.New("assume root denied"),
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"sqs-queue-policy", "--account", "123456789012", "--queue", "https://sqs/q1"})

	require.Error(t, cmd.Execute())
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
