package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestRecoveryCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		recoveryResult: []rootmanager.RecoveryResult{
			{AccountId: "123456789012", Success: true},
		},
	}

	var buf bytes.Buffer
	cmd := Recovery(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output, got empty buffer")
	}
}

func TestRecoveryCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Recovery(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, factoryErr) {
		t.Errorf("expected factory error, got: %v", err)
	}
}

func TestRecoveryCommand_AlreadyExists(t *testing.T) {
	mock := &mockRootManager{
		recoveryResult: []rootmanager.RecoveryResult{
			{AccountId: "123456789012", Success: false, Error: ""},
		},
	}

	var buf bytes.Buffer
	cmd := Recovery(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error for already-exists case, got: %v", err)
	}
	if !strings.Contains(buf.String(), "already exists") {
		t.Errorf("expected 'already exists' in output, got: %s", buf.String())
	}
}

func TestRecoveryCommand_RecoveryFailure(t *testing.T) {
	mock := &mockRootManager{
		recoveryResult: []rootmanager.RecoveryResult{
			{AccountId: "123456789012", Success: false, Error: "account suspended"},
		},
	}

	cmd := Recovery(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--accounts", "123456789012"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for failed recovery, got nil")
	}
}
