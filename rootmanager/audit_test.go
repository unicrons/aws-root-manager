package rootmanager

import (
	"context"
	"errors"
	"testing"
)

func TestAuditAccounts_CheckAccessError(t *testing.T) {
	accessErr := errors.New("trusted access not enabled")
	iam := &mockIamClient{
		checkOrgRootAccessErrs: []error{accessErr},
	}

	_, err := auditAccounts(context.Background(), iam, nil, nil, []string{"123456789012"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, accessErr) {
		t.Errorf("expected access error, got: %v", err)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error != "" {
		t.Errorf("expected no error in result, got: %s", results[0].Error)
	}
	if !results[0].LoginProfile {
		t.Error("expected LoginProfile=true")
	}
	if len(results[0].AccessKeys) != 1 || results[0].AccessKeys[0] != "AKIA123" {
		t.Errorf("unexpected AccessKeys: %v", results[0].AccessKeys)
	}
}

func TestAuditAccounts_STSError(t *testing.T) {
	stsErr := errors.New("assume root denied")
	sts := &mockStsClient{assumeRootErr: stsErr}
	iam := &mockIamClient{}

	results, err := auditAccounts(context.Background(), iam, sts, nil, []string{"123456789012"})
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error == "" {
		t.Error("expected error in result, got empty string")
	}
}

func TestAuditAccounts_MultipleAccounts(t *testing.T) {
	rootIam := &mockIamClient{}
	sts := &mockStsClient{}
	factory := &mockIamClientFactory{client: rootIam}
	iam := &mockIamClient{}

	accounts := []string{"111111111111", "222222222222", "333333333333"}
	results, err := auditAccounts(context.Background(), iam, sts, factory, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != len(accounts) {
		t.Errorf("expected %d results, got %d", len(accounts), len(results))
	}
}
