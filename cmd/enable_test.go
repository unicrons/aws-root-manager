package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestEnableCommand_Success(t *testing.T) {
	mock := &mockRootManager{
		enableInit:  rootmanager.RootAccessStatus{TrustedAccess: false, RootCredentialsManagement: false, RootSessions: false},
		enableFinal: rootmanager.RootAccessStatus{TrustedAccess: true, RootCredentialsManagement: true, RootSessions: false},
	}

	var buf bytes.Buffer
	cmd := Enable(newMockFactory(mock))
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output, got empty buffer")
	}
}

func TestEnableCommand_WithRootSessions(t *testing.T) {
	mock := &mockRootManager{
		enableInit:  rootmanager.RootAccessStatus{TrustedAccess: true, RootCredentialsManagement: true, RootSessions: false},
		enableFinal: rootmanager.RootAccessStatus{TrustedAccess: true, RootCredentialsManagement: true, RootSessions: true},
	}

	var buf bytes.Buffer
	cmd := Enable(newMockFactory(mock))
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--enableRootSessions=true"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output, got empty buffer")
	}
}

func TestEnableCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Enable(newFailingFactory(factoryErr))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, factoryErr) {
		t.Errorf("expected factory error, got: %v", err)
	}
}

func TestEnableCommand_EnableError(t *testing.T) {
	enableErr := errors.New("trusted access already enabled")
	mock := &mockRootManager{enableErr: enableErr}

	cmd := Enable(newMockFactory(mock))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, enableErr) {
		t.Errorf("expected enable error, got: %v", err)
	}
}
