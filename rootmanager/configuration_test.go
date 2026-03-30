package rootmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	internalaws "github.com/unicrons/aws-root-manager/internal/aws"
)

func TestCheckRootAccess_AllEnabled(t *testing.T) {
	iam := &mockIamClient{} // CheckOrganizationRootAccess returns nil → all enabled

	status, err := checkRootAccess(context.Background(), iam)
	require.NoError(t, err)
	assert.True(t, status.TrustedAccess)
	assert.True(t, status.RootCredentialsManagement)
	assert.True(t, status.RootSessions)
}

func TestCheckRootAccess_TrustedAccessDisabled(t *testing.T) {
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrTrustedAccessNotEnabled},
	}

	status, err := checkRootAccess(context.Background(), iam)
	require.NoError(t, err)
	assert.False(t, status.TrustedAccess)
	assert.False(t, status.RootCredentialsManagement)
	assert.False(t, status.RootSessions)
}

func TestCheckRootAccess_CredentialsManagementDisabled(t *testing.T) {
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrRootCredentialsManagementNotEnabled},
	}

	status, err := checkRootAccess(context.Background(), iam)
	require.NoError(t, err)
	assert.True(t, status.TrustedAccess)
	assert.False(t, status.RootCredentialsManagement)
	assert.False(t, status.RootSessions)
}

func TestCheckRootAccess_RootSessionsDisabled(t *testing.T) {
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrRootSessionsNotEnabled},
	}

	status, err := checkRootAccess(context.Background(), iam)
	require.NoError(t, err)
	assert.True(t, status.TrustedAccess)
	assert.True(t, status.RootCredentialsManagement)
	assert.False(t, status.RootSessions)
}

func TestCheckRootAccess_UnknownError(t *testing.T) {
	unexpected := errors.New("unexpected AWS error")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{unexpected},
	}

	_, err := checkRootAccess(context.Background(), iam)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpected)
}

func TestEnableRootAccess_AlreadyEnabled(t *testing.T) {
	iam := &mockIamClient{} // both checkRootAccess calls return nil → all enabled
	org := &mockOrganizationsClient{}

	init, final, err := enableRootAccess(context.Background(), iam, org, false)
	require.NoError(t, err)
	assert.True(t, init.TrustedAccess)
	assert.True(t, init.RootCredentialsManagement)
	assert.True(t, final.TrustedAccess)
	assert.True(t, final.RootCredentialsManagement)
}

func TestEnableRootAccess_AllDisabled(t *testing.T) {
	// First checkRootAccess → TrustedAccess disabled (init state)
	// Second checkRootAccess (after enabling) → all enabled (final state)
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrTrustedAccessNotEnabled, nil},
	}
	org := &mockOrganizationsClient{}

	init, final, err := enableRootAccess(context.Background(), iam, org, false)
	require.NoError(t, err)
	assert.False(t, init.TrustedAccess)
	assert.False(t, init.RootCredentialsManagement)
	assert.True(t, final.TrustedAccess)
	assert.True(t, final.RootCredentialsManagement)
}

func TestEnableRootAccess_EnableSessionsTrue(t *testing.T) {
	// Init: sessions disabled, Final: all enabled
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrRootSessionsNotEnabled, nil},
	}
	org := &mockOrganizationsClient{}

	_, final, err := enableRootAccess(context.Background(), iam, org, true)
	require.NoError(t, err)
	assert.True(t, final.RootSessions)
}

func TestEnableRootAccess_EnableAWSServiceAccessError(t *testing.T) {
	serviceErr := errors.New("org service access denied")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{ErrTrustedAccessNotEnabled},
	}
	org := &mockOrganizationsClient{enableServiceAccessErr: serviceErr}

	_, _, err := enableRootAccess(context.Background(), iam, org, false)
	require.Error(t, err)
	assert.ErrorIs(t, err, serviceErr)
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
