package rootmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditAccounts_CheckAccessError(t *testing.T) {
	accessErr := errors.New("trusted access not enabled")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{accessErr},
	}

	_, err := auditAccounts(context.Background(), iam, nil, nil, []string{"123456789012"})
	require.Error(t, err)
	assert.ErrorIs(t, err, accessErr)
}

func TestAuditAccounts_Success(t *testing.T) {
	rootIam := &mockIamClient{
		getLoginProfileResult: true,
		listAccessKeysResult:  []string{"AKIA123"},
		listMFADevicesResult:  []string{},
		listCertsResult:       []string{},
	}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	results, err := auditAccounts(context.Background(), iam, sts, factory, []string{"123456789012"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Empty(t, results[0].Error)
	assert.True(t, results[0].LoginProfile)
	assert.Equal(t, []string{"AKIA123"}, results[0].AccessKeys)
}

func TestAuditAccounts_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	iam := &mockIamClient{}

	results, err := auditAccounts(context.Background(), iam, sts, nil, []string{"123456789012"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.NotEmpty(t, results[0].Error)
}

func TestAuditAccounts_MultipleAccounts(t *testing.T) {
	rootIam := &mockIamClient{}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	accounts := []string{"111111111111", "222222222222", "333333333333"}
	results, err := auditAccounts(context.Background(), iam, sts, factory, accounts)
	require.NoError(t, err)
	assert.Len(t, results, len(accounts))
}
