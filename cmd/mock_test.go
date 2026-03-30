package cmd

import (
	"context"

	"github.com/unicrons/aws-root-manager/rootmanager"
)

// mockRootManager implements rootmanager.RootManager for testing.
type mockRootManager struct {
	checkResult    rootmanager.RootAccessStatus
	checkErr       error
	auditResult    []rootmanager.RootCredentials
	auditErr       error
	enableInit     rootmanager.RootAccessStatus
	enableFinal    rootmanager.RootAccessStatus
	enableErr      error
	deleteResult   []rootmanager.DeletionResult
	deleteErr      error
	recoveryResult []rootmanager.RecoveryResult
	recoveryErr    error
}

func (m *mockRootManager) CheckRootAccess(_ context.Context) (rootmanager.RootAccessStatus, error) {
	return m.checkResult, m.checkErr
}
func (m *mockRootManager) AuditAccounts(_ context.Context, _ []string) ([]rootmanager.RootCredentials, error) {
	return m.auditResult, m.auditErr
}
func (m *mockRootManager) EnableRootAccess(_ context.Context, _ bool) (rootmanager.RootAccessStatus, rootmanager.RootAccessStatus, error) {
	return m.enableInit, m.enableFinal, m.enableErr
}
func (m *mockRootManager) DeleteCredentials(_ context.Context, _ []rootmanager.RootCredentials, _ string) ([]rootmanager.DeletionResult, error) {
	return m.deleteResult, m.deleteErr
}
func (m *mockRootManager) RecoverRootPassword(_ context.Context, _ []string) ([]rootmanager.RecoveryResult, error) {
	return m.recoveryResult, m.recoveryErr
}

// newMockFactory returns a factory function that always returns the given mock.
func newMockFactory(mock rootmanager.RootManager) func(context.Context) (rootmanager.RootManager, error) {
	return func(_ context.Context) (rootmanager.RootManager, error) {
		return mock, nil
	}
}

// newFailingFactory returns a factory function that always returns the given error.
func newFailingFactory(err error) func(context.Context) (rootmanager.RootManager, error) {
	return func(_ context.Context) (rootmanager.RootManager, error) {
		return nil, err
	}
}
