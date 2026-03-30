package ui

import (
	"context"
	"errors"
	"testing"

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accounts) != 2 {
		t.Errorf("expected 2 accounts, got %d", len(accounts))
	}
	if accounts[0] != "123456789012" || accounts[1] != "234567890123" {
		t.Errorf("unexpected accounts: %v", accounts)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// management account must be excluded
	if len(accounts) != 2 {
		t.Errorf("expected 2 non-management accounts, got %d: %v", len(accounts), accounts)
	}
	for _, id := range accounts {
		if id == "000000000000" {
			t.Error("management account must not be in the result")
		}
	}
}

func TestSelectTargetAccounts_DescribeOrgError(t *testing.T) {
	orgErr := errors.New("organizations API unavailable")
	mock := &mockOrganizationsClient{describeErr: orgErr}

	_, err := SelectTargetAccounts(context.Background(), mock, []string{"all"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, orgErr) {
		t.Errorf("expected org error, got: %v", err)
	}
}

func TestSelectTargetAccounts_ListAccountsError(t *testing.T) {
	listErr := errors.New("failed to list accounts")
	mock := &mockOrganizationsClient{
		managementAccount: "000000000000",
		listErr:           listErr,
	}

	_, err := SelectTargetAccounts(context.Background(), mock, []string{"all"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, listErr) {
		t.Errorf("expected list error, got: %v", err)
	}
}
