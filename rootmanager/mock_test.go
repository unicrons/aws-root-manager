package rootmanager

import (
	"context"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/unicrons/aws-root-manager/internal/aws"
)

// mockIamClient implements aws.IamClient for testing.
// checkOrgRootAccessErrs is consumed in order on each call; the last entry repeats.
type mockIamClient struct {
	checkOrgRootAccessErrs []error
	checkOrgRootAccessCall int

	getLoginProfileResult bool
	getLoginProfileErr    error

	listAccessKeysResult []string
	listAccessKeysErr    error

	listMFADevicesResult []string
	listMFADevicesErr    error

	listCertsResult []string
	listCertsErr    error

	deleteLoginProfileErr error
	deleteAccessKeysErr   error
	deactivateMFAErr      error
	deleteCertsErr        error

	enableCredMgmtErr  error
	enableSessionsErr  error
	createLoginProfile error
}

func (m *mockIamClient) CheckOrganizationRootAccess(_ context.Context, _ bool) error {
	if len(m.checkOrgRootAccessErrs) == 0 {
		return nil
	}
	idx := m.checkOrgRootAccessCall
	if idx >= len(m.checkOrgRootAccessErrs) {
		idx = len(m.checkOrgRootAccessErrs) - 1
	}
	m.checkOrgRootAccessCall++
	return m.checkOrgRootAccessErrs[idx]
}
func (m *mockIamClient) GetLoginProfile(_ context.Context, _ string) (bool, error) {
	return m.getLoginProfileResult, m.getLoginProfileErr
}
func (m *mockIamClient) DeleteLoginProfile(_ context.Context, _ string) error {
	return m.deleteLoginProfileErr
}
func (m *mockIamClient) ListAccessKeys(_ context.Context, _ string) ([]string, error) {
	return m.listAccessKeysResult, m.listAccessKeysErr
}
func (m *mockIamClient) DeleteAccessKeys(_ context.Context, _ string, _ []string) error {
	return m.deleteAccessKeysErr
}
func (m *mockIamClient) ListMFADevices(_ context.Context, _ string) ([]string, error) {
	return m.listMFADevicesResult, m.listMFADevicesErr
}
func (m *mockIamClient) DeactivateMFADevices(_ context.Context, _ string, _ []string) error {
	return m.deactivateMFAErr
}
func (m *mockIamClient) ListSigningCertificates(_ context.Context, _ string) ([]string, error) {
	return m.listCertsResult, m.listCertsErr
}
func (m *mockIamClient) DeleteSigningCertificates(_ context.Context, _ string, _ []string) error {
	return m.deleteCertsErr
}
func (m *mockIamClient) EnableOrganizationsRootCredentialsManagement(_ context.Context) error {
	return m.enableCredMgmtErr
}
func (m *mockIamClient) EnableOrganizationsRootSessions(_ context.Context) error {
	return m.enableSessionsErr
}
func (m *mockIamClient) CreateLoginProfile(_ context.Context) error {
	return m.createLoginProfile
}

// mockStsClient implements aws.StsClient for testing.
type mockStsClient struct {
	assumeRootErr error
}

func (m *mockStsClient) GetAssumeRootConfig(_ context.Context, _, _ string) (awssdk.Config, error) {
	return awssdk.Config{}, m.assumeRootErr
}

// mockIamClientFactory implements aws.IamClientFactory for testing.
// It always returns the same IamClient regardless of the config passed.
type mockIamClientFactory struct {
	client aws.IamClient
}

func (f *mockIamClientFactory) NewIamClient(_ awssdk.Config) aws.IamClient {
	return f.client
}
