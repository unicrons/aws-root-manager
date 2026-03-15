package rootmanager

import (
	"context"
	"errors"

	"github.com/unicrons/aws-root-manager/internal/aws"
)

// manager implements RootManager using AWS clients.
type manager struct {
	iam     aws.IamClient
	sts     aws.StsClient
	org     aws.OrganizationsClient
	factory aws.IamClientFactory
}

// newManager returns a RootManager that uses the given AWS clients and factory.
// sts and org may be nil for callers that only use CheckRootAccess.
func newManager(iam aws.IamClient, sts aws.StsClient, org aws.OrganizationsClient, factory aws.IamClientFactory) RootManager {
	return &manager{iam: iam, sts: sts, org: org, factory: factory}
}

// NewRootManager returns a RootManager configured from the default AWS environment.
// It loads credentials from the standard AWS credential chain (env vars, ~/.aws, IAM role).
func NewRootManager(ctx context.Context) (RootManager, error) {
	cfg, err := aws.LoadAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return newManager(
		aws.NewIamClient(cfg),
		aws.NewStsClient(cfg),
		aws.NewOrganizationsClient(cfg),
		&aws.DefaultIamClientFactory{},
	), nil
}

func (m *manager) AuditAccounts(ctx context.Context, accountIds []string) ([]RootCredentials, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for audit")
	}
	return auditAccounts(ctx, m.iam, m.sts, m.factory, accountIds)
}

func (m *manager) CheckRootAccess(ctx context.Context) (RootAccessStatus, error) {
	return checkRootAccess(ctx, m.iam)
}

func (m *manager) EnableRootAccess(ctx context.Context, enableSessions bool) (RootAccessStatus, RootAccessStatus, error) {
	if m.org == nil {
		return RootAccessStatus{}, RootAccessStatus{}, errors.New("Organizations client required for enable")
	}
	return enableRootAccess(ctx, m.iam, m.org, enableSessions)
}

func (m *manager) DeleteCredentials(ctx context.Context, creds []RootCredentials, credentialType string) ([]DeletionResult, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for delete")
	}
	return deleteAccountsCredentials(ctx, m.iam, m.sts, m.factory, creds, credentialType)
}

func (m *manager) RecoverRootPassword(ctx context.Context, accountIds []string) ([]RecoveryResult, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for recovery")
	}
	return recoverAccountsRootPassword(ctx, m.iam, m.sts, m.factory, accountIds)
}
