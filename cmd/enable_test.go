package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestEnableCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		enableInit:  rootmanager.RootAccessStatus{TrustedAccess: false, RootCredentialsManagement: false, RootSessions: false},
		enableFinal: rootmanager.RootAccessStatus{TrustedAccess: true, RootCredentialsManagement: true, RootSessions: false},
	}

	var buf bytes.Buffer
	cmd := Enable(newMockFactory(mock))
	cmd.SetOut(&buf)

	require.NoError(t, cmd.Execute())
	assert.NotEmpty(t, buf.String())
}

func TestEnableCommand_WithRootSessions(t *testing.T) {
	mock := &mockRootManager{
		enableInit:  rootmanager.RootAccessStatus{TrustedAccess: true, RootCredentialsManagement: true, RootSessions: false},
		enableFinal: rootmanager.RootAccessStatus{TrustedAccess: true, RootCredentialsManagement: true, RootSessions: true},
	}

	var buf bytes.Buffer
	cmd := Enable(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--enableRootSessions=true"})

	require.NoError(t, cmd.Execute())
	assert.NotEmpty(t, buf.String())
}

func TestEnableCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Enable(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestEnableCommand_EnableError(t *testing.T) {
	enableErr := errors.New("trusted access already enabled")
	mock := &mockRootManager{enableErr: enableErr}

	cmd := Enable(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, enableErr)
}
