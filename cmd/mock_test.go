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

	getBucketPolicyResult string
	getBucketPolicyErr    error
	listBucketsResult     []string
	listBucketsErr        error
	deleteBucketResult    rootmanager.PolicyDeletionResult
	deleteBucketErr       error

	getQueuePolicyResult string
	getQueuePolicyErr    error
	listQueuesResult     []string
	listQueuesErr        error
	deleteQueueResult    rootmanager.PolicyDeletionResult
	deleteQueueErr       error
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
func (m *mockRootManager) GetS3BucketPolicy(_ context.Context, _, _ string) (string, error) {
	return m.getBucketPolicyResult, m.getBucketPolicyErr
}
func (m *mockRootManager) ListAccountBuckets(_ context.Context, _ string) ([]string, error) {
	return m.listBucketsResult, m.listBucketsErr
}
func (m *mockRootManager) DeleteS3BucketPolicy(_ context.Context, _, _ string) (rootmanager.PolicyDeletionResult, error) {
	return m.deleteBucketResult, m.deleteBucketErr
}
func (m *mockRootManager) GetSQSQueuePolicy(_ context.Context, _, _ string) (string, error) {
	return m.getQueuePolicyResult, m.getQueuePolicyErr
}
func (m *mockRootManager) ListAccountQueues(_ context.Context, _ string) ([]string, error) {
	return m.listQueuesResult, m.listQueuesErr
}
func (m *mockRootManager) DeleteSQSQueuePolicy(_ context.Context, _, _ string) (rootmanager.PolicyDeletionResult, error) {
	return m.deleteQueueResult, m.deleteQueueErr
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
