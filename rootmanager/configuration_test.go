package rootmanager

import (
	"context"
	"errors"
	"testing"

	internalaws "github.com/unicrons/aws-root-manager/internal/aws"
)

func TestCheckRootAccess_AllEnabled(t *testing.T) {
	iam := &mockIamClient{} // CheckOrganizationRootAccess returns nil → all enabled

	status, err := checkRootAccess(context.Background(), iam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.TrustedAccess || !status.RootCredentialsManagement || !status.RootSessions {
		t.Errorf("expected all enabled, got: %+v", status)
	}
}

func TestCheckRootAccess_TrustedAccessDisabled(t *testing.T) {
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrTrustedAccessNotEnabled},
	}

	status, err := checkRootAccess(context.Background(), iam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.TrustedAccess || status.RootCredentialsManagement || status.RootSessions {
		t.Errorf("expected all disabled, got: %+v", status)
	}
}

func TestCheckRootAccess_CredentialsManagementDisabled(t *testing.T) {
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrRootCredentialsManagementNotEnabled},
	}

	status, err := checkRootAccess(context.Background(), iam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.TrustedAccess {
		t.Error("expected TrustedAccess=true")
	}
	if status.RootCredentialsManagement || status.RootSessions {
		t.Errorf("expected CredMgmt and Sessions disabled, got: %+v", status)
	}
}

func TestCheckRootAccess_RootSessionsDisabled(t *testing.T) {
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrRootSessionsNotEnabled},
	}

	status, err := checkRootAccess(context.Background(), iam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.TrustedAccess || !status.RootCredentialsManagement {
		t.Errorf("expected TrustedAccess and CredMgmt enabled, got: %+v", status)
	}
	if status.RootSessions {
		t.Error("expected RootSessions=false")
	}
}

func TestCheckRootAccess_UnknownError(t *testing.T) {
	unexpected := errors.New("unexpected AWS error")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{unexpected},
	}

	_, err := checkRootAccess(context.Background(), iam)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, unexpected) {
		t.Errorf("expected unexpected error, got: %v", err)
	}
}

func TestEnableRootAccess_AlreadyEnabled(t *testing.T) {
	iam := &mockIamClient{} // both checkRootAccess calls return nil → all enabled
	org := &mockOrganizationsClient{}

	init, final, err := enableRootAccess(context.Background(), iam, org, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !init.TrustedAccess || !init.RootCredentialsManagement {
		t.Errorf("expected init all enabled, got: %+v", init)
	}
	if !final.TrustedAccess || !final.RootCredentialsManagement {
		t.Errorf("expected final all enabled, got: %+v", final)
	}
}

func TestEnableRootAccess_AllDisabled(t *testing.T) {
	// First checkRootAccess → TrustedAccess disabled (init state)
	// Second checkRootAccess (after enabling) → all enabled (final state)
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrTrustedAccessNotEnabled, nil},
	}
	org := &mockOrganizationsClient{}

	init, final, err := enableRootAccess(context.Background(), iam, org, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if init.TrustedAccess || init.RootCredentialsManagement {
		t.Errorf("expected init all disabled, got: %+v", init)
	}
	if !final.TrustedAccess || !final.RootCredentialsManagement {
		t.Errorf("expected final all enabled, got: %+v", final)
	}
}

func TestEnableRootAccess_EnableSessionsTrue(t *testing.T) {
	// Init: sessions disabled, Final: all enabled
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrRootSessionsNotEnabled, nil},
	}
	org := &mockOrganizationsClient{}

	_, final, err := enableRootAccess(context.Background(), iam, org, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !final.RootSessions {
		t.Errorf("expected RootSessions enabled after enabling with enableSessions=true, got: %+v", final)
	}
}

func TestEnableRootAccess_EnableAWSServiceAccessError(t *testing.T) {
	serviceErr := errors.New("org service access denied")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrTrustedAccessNotEnabled},
	}
	org := &mockOrganizationsClient{enableServiceAccessErr: serviceErr}

	_, _, err := enableRootAccess(context.Background(), iam, org, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, serviceErr) {
		t.Errorf("expected service access error, got: %v", err)
	}
}

// mockOrganizationsClient implements aws.OrganizationsClient for rootmanager tests.
type mockOrganizationsClient struct {
	enableServiceAccessErr error
}

func (m *mockOrganizationsClient) DescribeOrganization(_ context.Context) (string, error) {
	return "000000000000", nil
}
func (m *mockOrganizationsClient) ListAccounts(_ context.Context) ([]internalaws.OrganizationAccount, error) {
	return nil, nil
}

func (m *mockOrganizationsClient) EnableAWSServiceAccess(_ context.Context, _ string) error {
	return m.enableServiceAccessErr
}
