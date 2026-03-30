package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestAuditCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012", LoginProfile: false, AccessKeys: []string{}, MfaDevices: []string{}, SigningCertificates: []string{}},
		},
	}

	var buf bytes.Buffer
	cmd := Audit(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	require.NoError(t, cmd.Execute())
	assert.NotEmpty(t, buf.String())
}

func TestAuditCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Audit(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, factoryErr)
}

func TestAuditCommand_AuditError(t *testing.T) {
	auditErr := errors.New("sts assume root failed")
	mock := &mockRootManager{auditErr: auditErr}

	cmd := Audit(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, auditErr)
}

func TestAuditCommand_SkippedAccounts(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012", Error: "access denied"},
		},
	}

	cmd := Audit(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	require.Error(t, cmd.Execute())
}
