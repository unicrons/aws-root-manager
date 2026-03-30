package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/unicrons/aws-root-manager/rootmanager"
)

func TestCheckCommand_AllEnabled(t *testing.T) {
	mock := &mockRootManager{
		checkResult: rootmanager.RootAccessStatus{
			TrustedAccess:             true,
			RootCredentialsManagement: true,
			RootSessions:              true,
		},
	}

	var buf bytes.Buffer
	cmd := Check(newMockFactory(mock))
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCheckCommand_AllDisabled(t *testing.T) {
	mock := &mockRootManager{
		checkResult: rootmanager.RootAccessStatus{
			TrustedAccess:             false,
			RootCredentialsManagement: false,
			RootSessions:              false,
		},
	}

	var buf bytes.Buffer
	cmd := Check(newMockFactory(mock))
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCheckCommand_FactoryError(t *testing.T) {
	factoryErr := errors.New("failed to load AWS config")

	cmd := Check(newFailingFactory(factoryErr))
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, factoryErr) {
		t.Errorf("expected factory error, got: %v", err)
	}
}

func TestCheckCommand_CheckRootAccessError(t *testing.T) {
	checkErr := errors.New("organization not found")
	mock := &mockRootManager{checkErr: checkErr}

	cmd := Check(newMockFactory(mock))
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, checkErr) {
		t.Errorf("expected check error, got: %v", err)
	}
}
