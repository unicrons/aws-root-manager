package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestCheckCommand_AllEnabled(t *testing.T) {
	mock := &mockRootManager{
		checkResult: rootmanager.RootAccessStatus{
			TrustedAccess:             true,
			RootCredentialsManagement: true,
			RootSessions:              true,
		},
	}

	var buf bytes.Buffer
	cmd := Check(newMockFactory(mock))
	cmd.SetOut(&buf)

	require.NoError(t, cmd.Execute())
	assert.Contains(t, buf.String(), "true")
}

func TestCheckCommand_AllDisabled(t *testing.T) {
	mock := &mockRootManager{
		checkResult: rootmanager.RootAccessStatus{
			TrustedAccess:             false,
			RootCredentialsManagement: false,
			RootSessions:              false,
		},
	}

	var buf bytes.Buffer
	cmd := Check(newMockFactory(mock))
	cmd.SetOut(&buf)

	require.NoError(t, cmd.Execute())
	assert.Contains(t, buf.String(), "false")
}

func TestCheckCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Check(newFailingFactory(factoryErr))
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()

	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestCheckCommand_CheckRootAccessError(t *testing.T) {
	checkErr := errors.New("organization not found")
	mock := &mockRootManager{checkErr: checkErr}

	cmd := Check(newMockFactory(mock))
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()

	require.Error(t, err)
	assert.ErrorIs(t, err, checkErr)
}
