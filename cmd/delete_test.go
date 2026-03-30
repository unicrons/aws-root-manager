package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestDeleteCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012"},
		},
		deleteResult: []rootmanager.DeletionResult{
			{AccountId: "123456789012", CredentialType: "all", Success: true},
		},
	}

	var buf bytes.Buffer
	cmd := Delete(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"all", "--accounts", "123456789012"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output, got empty buffer")
	}
}

func TestDeleteCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Delete(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"all", "--accounts", "123456789012"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, factoryErr) {
		t.Errorf("expected factory error, got: %v", err)
	}
}

func TestDeleteCommand_DeleteFailure(t *testing.T) {
	mock := &mockRootManager{
		auditResult: []rootmanager.RootCredentials{
			{AccountId: "123456789012"},
		},
		deleteResult: []rootmanager.DeletionResult{
			{AccountId: "123456789012", CredentialType: "all", Success: false, Error: "access denied"},
		},
	}

	cmd := Delete(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"all", "--accounts", "123456789012"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for failed deletion, got nil")
	}
}
