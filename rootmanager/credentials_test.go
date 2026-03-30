package rootmanager

import (
	"context"
	"errors"
	"testing"
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
			if got := hasCredentialsToDelete(tt.creds, "all"); got != tt.want {
				t.Errorf("hasCredentialsToDelete(%q) = %v, want %v", tt.name, got, tt.want)
			}
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

	if !hasCredentialsToDelete(creds, "login") {
		t.Error("expected login=true")
	}
	if !hasCredentialsToDelete(creds, "keys") {
		t.Error("expected keys=true")
	}
	if !hasCredentialsToDelete(creds, "mfa") {
		t.Error("expected mfa=true")
	}
	if !hasCredentialsToDelete(creds, "certificate") {
		t.Error("expected certificate=true")
	}
}

func TestHasCredentialsToDelete_UnknownType(t *testing.T) {
	creds := RootCredentials{LoginProfile: true}
	if hasCredentialsToDelete(creds, "unknown") {
		t.Error("expected false for unknown credential type")
	}
}

// --- deleteAccountsCredentials ---

func TestDeleteAccountsCredentials_CheckAccessError(t *testing.T) {
	accessErr := errors.New("root access not enabled")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{accessErr},
	}

	_, err := deleteAccountsCredentials(context.Background(), iam, nil, nil, []RootCredentials{}, "all")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, accessErr) {
		t.Errorf("expected access error, got: %v", err)
	}
}

func TestDeleteAccountsCredentials_NoCredentials(t *testing.T) {
	iam := &mockIamClient{}
	sts := &mockStsClient{assumeRootErr: errors.New("should not be called")}

	creds := []RootCredentials{
		{AccountId: "123456789012"}, // no credentials set
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, nil, creds, "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Success {
		t.Errorf("expected success when no credentials to delete, got error: %s", results[0].Error)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Success {
		t.Errorf("expected success, got error: %s", results[0].Error)
	}
}

func TestDeleteAccountsCredentials_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	iam := &mockIamClient{}

	creds := []RootCredentials{
		{AccountId: "123456789012", LoginProfile: true},
	}

	results, err := deleteAccountsCredentials(context.Background(), iam, sts, nil, creds, "login")
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if results[0].Success {
		t.Error("expected failure for STS error")
	}
	if results[0].Error == "" {
		t.Error("expected error message in result")
	}
}

// --- recoverAccountsRootPassword ---

func TestRecoverAccountsRootPassword_CheckAccessError(t *testing.T) {
	accessErr := errors.New("root access not enabled")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{accessErr},
	}

	_, err := recoverAccountsRootPassword(context.Background(), iam, nil, nil, []string{"123456789012"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, accessErr) {
		t.Errorf("expected access error, got: %v", err)
	}
}

func TestRecoverAccountsRootPassword_Success(t *testing.T) {
	rootIam := &mockIamClient{} // CreateLoginProfile returns nil
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	results, err := recoverAccountsRootPassword(context.Background(), iam, sts, factory, []string{"123456789012"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Success {
		t.Errorf("expected success, got: %+v", results[0])
	}
}

func TestRecoverAccountsRootPassword_AlreadyExists(t *testing.T) {
	rootIam := &mockIamClient{createLoginProfile: ErrEntityAlreadyExists}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	results, err := recoverAccountsRootPassword(context.Background(), iam, sts, factory, []string{"123456789012"})
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if results[0].Success {
		t.Error("expected success=false when login profile already exists")
	}
	if results[0].Error != "" {
		t.Errorf("expected no error string when profile already exists, got: %s", results[0].Error)
	}
}

func TestRecoverAccountsRootPassword_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	iam := &mockIamClient{}

	results, err := recoverAccountsRootPassword(context.Background(), iam, sts, nil, []string{"123456789012"})
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if results[0].Success {
		t.Error("expected failure for STS error")
	}
	if results[0].Error == "" {
		t.Error("expected error message in result")
	}
}
