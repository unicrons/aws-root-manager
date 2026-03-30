package ui

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unicrons/aws-root-manager/internal/aws"
)

type mockOrganizationsClient struct {
	managementAccount string
	accounts          []aws.OrganizationAccount
	describeErr       error
	listErr           error
}

func (m *mockOrganizationsClient) DescribeOrganization(_ context.Context) (string, error) {
	return m.managementAccount, m.describeErr
}

func (m *mockOrganizationsClient) ListAccounts(_ context.Context) ([]aws.OrganizationAccount, error) {
	return m.accounts, m.listErr
}

func (m *mockOrganizationsClient) EnableAWSServiceAccess(_ context.Context, _ string) error {
	return nil
}

func TestSelectTargetAccounts_ExplicitIDs(t *testing.T) {
	// org is nil to prove it's never called on the explicit IDs path
	accounts, err := SelectTargetAccounts(context.Background(), nil, []string{"123456789012", "234567890123"})
	require.NoError(t, err)
	assert.Equal(t, []string{"123456789012", "234567890123"}, accounts)
}

func TestSelectTargetAccounts_AllFlag(t *testing.T) {
	mock := &mockOrganizationsClient{
		managementAccount: "000000000000",
		accounts: []aws.OrganizationAccount{
			{AccountID: "111111111111", Name: "account-a"},
			{AccountID: "000000000000", Name: "management"},
			{AccountID: "222222222222", Name: "account-b"},
		},
	}

	accounts, err := SelectTargetAccounts(context.Background(), mock, []string{"all"})
	require.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.NotContains(t, accounts, "000000000000")
}

func TestSelectTargetAccounts_DescribeOrgError(t *testing.T) {
	orgErr := errors.New("organizations API unavailable")
	mock := &mockOrganizationsClient{describeErr: orgErr}

	_, err := SelectTargetAccounts(context.Background(), mock, []string{"all"})
	require.Error(t, err)
	assert.ErrorIs(t, err, orgErr)
}

func TestSelectTargetAccounts_ListAccountsError(t *testing.T) {
	listErr := errors.New("failed to list accounts")
	mock := &mockOrganizationsClient{
		managementAccount: "000000000000",
		listErr:           listErr,
	}

	_, err := SelectTargetAccounts(context.Background(), mock, []string{"all"})
	require.Error(t, err)
	assert.ErrorIs(t, err, listErr)
}
