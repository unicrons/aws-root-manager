package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestDeleteCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012"},
		},
		deleteResult: []rootmanager.DeletionResult{
			{AccountId: "123456789012", CredentialType: "all", Success: true},
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"all", "--accounts", "123456789012"})

	require.NoError(t, cmd.Execute())
	assert.NotEmpty(t, buf.String())
}

func TestDeleteCommand_CertificatesSubcommand(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012"},
		},
		deleteResult: []rootmanager.DeletionResult{
			{AccountId: "123456789012", CredentialType: "certificate", Success: true},
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"certificates", "--accounts", "123456789012"})

	require.NoError(t, cmd.Execute())
	assert.NotEmpty(t, buf.String())
}

func TestDeleteCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Delete(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"all", "--accounts", "123456789012"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestDeleteCommand_DeleteFailure(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012"},
		},
		deleteResult: []rootmanager.DeletionResult{
			{AccountId: "123456789012", CredentialType: "all", Success: false, Error: "access denied"},
		},
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"all", "--accounts", "123456789012"})

	require.Error(t, cmd.Execute())
}
