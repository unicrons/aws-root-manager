package cmd

import (
	"bytes"
	"errors"
	"testing"

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

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output, got empty buffer")
	}
}

func TestAuditCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Audit(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, factoryErr) {
		t.Errorf("expected factory error, got: %v", err)
	}
}

func TestAuditCommand_AuditError(t *testing.T) {
	auditErr := errors.New("sts assume root failed")
	mock := &mockRootManager{auditErr: auditErr}

	cmd := Audit(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, auditErr) {
		t.Errorf("expected audit error, got: %v", err)
	}
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

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for skipped accounts, got nil")
	}
}
