package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestRecoveryCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		recoveryResult: []rootmanager.RecoveryResult{
			{AccountId: "123456789012", Success: true},
		},
	}

	var buf bytes.Buffer
	cmd := Recovery(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--accounts", "123456789012", "--skip"})

	require.NoError(t, cmd.Execute())
	assert.NotEmpty(t, buf.String())
}

func TestRecoveryCommand_AlreadyExists(t *testing.T) {
	mock := &mockRootManager{
		recoveryResult: []rootmanager.RecoveryResult{
			{AccountId: "123456789012", Success: false, Error: ""},
		},
	}

	var buf bytes.Buffer
	cmd := Recovery(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--accounts", "123456789012", "--skip"})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, buf.String(), "already exists")
}

func TestRecoveryCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Recovery(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestRecoveryCommand_RecoveryFailure(t *testing.T) {
	mock := &mockRootManager{
		recoveryResult: []rootmanager.RecoveryResult{
			{AccountId: "123456789012", Success: false, Error: "account suspended"},
		},
	}

	cmd := Recovery(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--accounts", "123456789012", "--skip"})

	require.Error(t, cmd.Execute())
}
