package rootmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- hasCredentialsToDelete ---

func TestHasCredentialsToDelete_All(t *testing.T) {
	tests := []struct {
		name  string
		creds RootCredentials
		want  bool
	}{
		{"login profile", RootCredentials{LoginProfile: true}, true},
		{"access keys", RootCredentials{AccessKeys: []string{"key1"}}, true},
		{"mfa devices", RootCredentials{MfaDevices: []string{"mfa1"}}, true},
		{"certificates", RootCredentials{SigningCertificates: []string{"cert1"}}, true},
		{"none", RootCredentials{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, hasCredentialsToDelete(tt.creds, "all"))
		})
	}
}

func TestHasCredentialsToDelete_SpecificTypes(t *testing.T) {
	creds := RootCredentials{
		LoginProfile:        true,
		AccessKeys:          []string{"key1"},
		MfaDevices:          []string{"mfa1"},
		SigningCertificates: []string{"cert1"},
	}

	assert.True(t, hasCredentialsToDelete(creds, "login"))
	assert.True(t, hasCredentialsToDelete(creds, "keys"))
	assert.True(t, hasCredentialsToDelete(creds, "mfa"))
	assert.True(t, hasCredentialsToDelete(creds, "certificate"))
}

func TestHasCredentialsToDelete_UnknownType(t *testing.T) {
	creds := RootCredentials{LoginProfile: true}
	assert.False(t, hasCredentialsToDelete(creds, "unknown"))
}

// --- deleteAccountsCredentials ---

func TestDeleteAccountsCredentials_CheckAccessError(t *testing.T) {
	accessErr := errors.New("root access not enabled")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{accessErr},
	}

	_, err := deleteAccountsCredentials(context.Background(), iam, nil, nil, []RootCredentials{}, "all")
	require.Error(t, err)
	assert.ErrorIs(t, err, accessErr)
}

func TestDeleteAccountsCredentials_NoCredentials(t *testing.T) {
	iam := &mockIamClient{}
	sts := &mockStsClient{assumeRootErr: errors.New("should not be called")}

	creds := []RootCredentials{
		{AccountId: "123456789012"}, // no credentials set
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, nil, creds, "all")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
}

func TestDeleteAccountsCredentials_DeleteLoginProfile(t *testing.T) {
	rootIam := &mockIamClient{}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	creds := []RootCredentials{
		{AccountId: "123456789012", LoginProfile: true},
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, factory, creds, "login")
	require.NoError(t, err)
	assert.True(t, results[0].Success)
}

func TestDeleteAccountsCredentials_DeleteAccessKeys(t *testing.T) {
	rootIam := &mockIamClient{}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	creds := []RootCredentials{
		{AccountId: "123456789012", AccessKeys: []string{"AKIA123"}},
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, factory, creds, "keys")
	require.NoError(t, err)
	assert.True(t, results[0].Success)
}

func TestDeleteAccountsCredentials_DeactivateMFA(t *testing.T) {
	rootIam := &mockIamClient{}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	creds := []RootCredentials{
		{AccountId: "123456789012", MfaDevices: []string{"arn:aws:iam::123456789012:mfa/root"}},
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, factory, creds, "mfa")
	require.NoError(t, err)
	assert.True(t, results[0].Success)
}

func TestDeleteAccountsCredentials_DeleteCertificates(t *testing.T) {
	rootIam := &mockIamClient{}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	creds := []RootCredentials{
		{AccountId: "123456789012", SigningCertificates: []string{"cert-id-1"}},
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, factory, creds, "certificate")
	require.NoError(t, err)
	assert.True(t, results[0].Success)
}

func TestDeleteAccountsCredentials_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	iam := &mockIamClient{}

	creds := []RootCredentials{
		{AccountId: "123456789012", LoginProfile: true},
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, nil, creds, "login")
	require.NoError(t, err)
	assert.False(t, results[0].Success)
	assert.NotEmpty(t, results[0].Error)
}

// --- recoverAccountsRootPassword ---

func TestRecoverAccountsRootPassword_CheckAccessError(t *testing.T) {
	accessErr := errors.New("root access not enabled")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{accessErr},
	}

	_, err := recoverAccountsRootPassword(context.Background(), iam, nil, nil, []string{"123456789012"})
	require.Error(t, err)
	assert.ErrorIs(t, err, accessErr)
}

func TestRecoverAccountsRootPassword_Success(t *testing.T) {
	rootIam := &mockIamClient{} // CreateLoginProfile returns nil
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	results, err := recoverAccountsRootPassword(context.Background(), iam, sts, factory, []string{"123456789012"})
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Success)
}

func TestRecoverAccountsRootPassword_AlreadyExists(t *testing.T) {
	rootIam := &mockIamClient{createLoginProfile: ErrEntityAlreadyExists}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	results, err := recoverAccountsRootPassword(context.Background(), iam, sts, factory, []string{"123456789012"})
	require.NoError(t, err)
	assert.False(t, results[0].Success)
	assert.Empty(t, results[0].Error)
}

func TestRecoverAccountsRootPassword_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	iam := &mockIamClient{}

	results, err := recoverAccountsRootPassword(context.Background(), iam, sts, nil, []string{"123456789012"})
	require.NoError(t, err)
	assert.False(t, results[0].Success)
	assert.NotEmpty(t, results[0].Error)
}
