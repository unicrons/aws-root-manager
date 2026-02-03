package service

import (
	"context"
	"errors"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

// rootManagerImpl implements rootmanager.RootManager and converts between
// rootmanager types (public API) and aws types (infra) at the boundary.
type rootManagerImpl struct {
	iam *aws.IamClient
	sts *aws.StsClient
	org *aws.OrganizationsClient
}

// NewRootManager returns a RootManager that uses the given AWS clients.
// sts and org may be nil for callers that only use CheckRootAccess.
func NewRootManager(iam *aws.IamClient, sts *aws.StsClient, org *aws.OrganizationsClient) rootmanager.RootManager {
	return &rootManagerImpl{iam: iam, sts: sts, org: org}
}

func (r *rootManagerImpl) AuditAccounts(ctx context.Context, accountIds []string) ([]rootmanager.RootCredentials, error) {
	if r.sts == nil {
		return nil, errors.New("STS client required for audit")
	}
	creds, err := AuditAccounts(ctx, r.iam, r.sts, accountIds)
	if err != nil {
		return nil, err
	}
	return toRootCredentialsSlice(creds), nil
}

func (r *rootManagerImpl) CheckRootAccess(ctx context.Context) (rootmanager.RootAccessStatus, error) {
	status, err := CheckRootAccess(ctx, r.iam)
	if err != nil {
		return rootmanager.RootAccessStatus{}, err
	}
	return toRootAccessStatus(status), nil
}

func (r *rootManagerImpl) EnableRootAccess(ctx context.Context, enableSessions bool) (rootmanager.RootAccessStatus, rootmanager.RootAccessStatus, error) {
	if r.org == nil {
		return rootmanager.RootAccessStatus{}, rootmanager.RootAccessStatus{}, errors.New("Organizations client required for enable")
	}
	initStatus, status, err := EnableRootAccess(ctx, r.iam, r.org, enableSessions)
	if err != nil {
		return toRootAccessStatus(initStatus), toRootAccessStatus(status), err
	}
	return toRootAccessStatus(initStatus), toRootAccessStatus(status), nil
}

func (r *rootManagerImpl) DeleteCredentials(ctx context.Context, creds []rootmanager.RootCredentials, credentialType string) error {
	if r.sts == nil {
		return errors.New("STS client required for delete")
	}
	awsCreds := fromRootCredentialsSlice(creds)
	if err := DeleteAccountsCredentials(ctx, r.iam, r.sts, awsCreds, credentialType); err != nil {
		return err
	}
	return nil
}

func (r *rootManagerImpl) RecoverRootPassword(ctx context.Context, accountIds []string) (map[string]bool, error) {
	if r.sts == nil {
		return nil, errors.New("STS client required for recovery")
	}
	resultMap, err := RecoverAccountsRootPassword(ctx, r.iam, r.sts, accountIds)
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}

func toRootCredentials(c aws.RootCredentials) rootmanager.RootCredentials {
	return rootmanager.RootCredentials{
		AccountId:           c.AccountId,
		LoginProfile:        c.LoginProfile,
		AccessKeys:          append([]string(nil), c.AccessKeys...),
		MfaDevices:          append([]string(nil), c.MfaDevices...),
		SigningCertificates: append([]string(nil), c.SigningCertificates...),
		Error:               c.Error,
	}
}

func toRootCredentialsSlice(creds []aws.RootCredentials) []rootmanager.RootCredentials {
	out := make([]rootmanager.RootCredentials, len(creds))
	for i := range creds {
		out[i] = toRootCredentials(creds[i])
	}
	return out
}

func fromRootCredentials(c rootmanager.RootCredentials) aws.RootCredentials {
	return aws.RootCredentials{
		AccountId:           c.AccountId,
		LoginProfile:        c.LoginProfile,
		AccessKeys:          append([]string(nil), c.AccessKeys...),
		MfaDevices:          append([]string(nil), c.MfaDevices...),
		SigningCertificates: append([]string(nil), c.SigningCertificates...),
		Error:               c.Error,
	}
}

func fromRootCredentialsSlice(creds []rootmanager.RootCredentials) []aws.RootCredentials {
	out := make([]aws.RootCredentials, len(creds))
	for i := range creds {
		out[i] = fromRootCredentials(creds[i])
	}
	return out
}

func toRootAccessStatus(s aws.RootAccessStatus) rootmanager.RootAccessStatus {
	return rootmanager.RootAccessStatus{
		TrustedAccess:             s.TrustedAccess,
		RootCredentialsManagement: s.RootCredentialsManagement,
		RootSessions:              s.RootSessions,
	}
}
