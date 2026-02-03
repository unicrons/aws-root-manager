package service

import (
	"context"
	"errors"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

// manager implements rootmanager.RootManager using AWS clients.
type manager struct {
	iam *aws.IamClient
	sts *aws.StsClient
	org *aws.OrganizationsClient
}

// NewRootManager returns a RootManager that uses the given AWS clients.
// sts and org may be nil for callers that only use CheckRootAccess.
func NewRootManager(iam *aws.IamClient, sts *aws.StsClient, org *aws.OrganizationsClient) rootmanager.RootManager {
	return &manager{iam: iam, sts: sts, org: org}
}

func (m *manager) AuditAccounts(ctx context.Context, accountIds []string) ([]rootmanager.RootCredentials, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for audit")
	}
	return auditAccounts(ctx, m.iam, m.sts, accountIds)
}

func (m *manager) CheckRootAccess(ctx context.Context) (rootmanager.RootAccessStatus, error) {
	return checkRootAccess(ctx, m.iam)
}

func (m *manager) EnableRootAccess(ctx context.Context, enableSessions bool) (rootmanager.RootAccessStatus, rootmanager.RootAccessStatus, error) {
	if m.org == nil {
		return rootmanager.RootAccessStatus{}, rootmanager.RootAccessStatus{}, errors.New("Organizations client required for enable")
	}
	return enableRootAccess(ctx, m.iam, m.org, enableSessions)
}

func (m *manager) DeleteCredentials(ctx context.Context, creds []rootmanager.RootCredentials, credentialType string) error {
	if m.sts == nil {
		return errors.New("STS client required for delete")
	}
	return deleteAccountsCredentials(ctx, m.iam, m.sts, creds, credentialType)
}

func (m *manager) RecoverRootPassword(ctx context.Context, accountIds []string) (map[string]bool, error) {
	if m.sts == nil {
		return nil, errors.New("STS client required for recovery")
	}
	resultMap, err := recoverAccountsRootPassword(ctx, m.iam, m.sts, accountIds)
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}
