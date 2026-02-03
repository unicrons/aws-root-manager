package service

import (
	"context"
	"sync"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
)

// Get root credentials for a list of AWS accounts.
func AuditAccounts(ctx context.Context, iam *aws.IamClient, sts *aws.StsClient, accounts []string) ([]aws.RootCredentials, error) {
	logger.Trace("service.AuditAccounts", "auditing accounts %s", accounts)

	rootCredentials := make([]aws.RootCredentials, len(accounts))
	var wgAccounts sync.WaitGroup

	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return nil, err
	}

	for i, accountId := range accounts {
		wgAccounts.Add(1)
		go func(idx int, accountId string) {
			defer wgAccounts.Done()
			if accStatus, err := auditAccount(ctx, sts, accountId); err != nil {
				logger.Error("service.AuditAccounts", err, "account %s: audit skipped", accountId)
				rootCredentials[idx] = aws.RootCredentials{AccountId: accountId, Error: err.Error()}
			} else {
				rootCredentials[idx] = accStatus
			}
		}(i, accountId)
	}

	wgAccounts.Wait()

	return rootCredentials, nil
}

// Get root credentials for a specific account
func auditAccount(ctx context.Context, sts *aws.StsClient, accountId string) (aws.RootCredentials, error) {
	logger.Trace("service.auditAccount", "auditing account %s", accountId)

	awscfgRoot, err := sts.GetAssumeRootConfig(ctx, accountId, "IAMAuditRootUserCredentials")
	if err != nil {
		return aws.RootCredentials{}, err
	}

	iamRoot := aws.NewIamClient(awscfgRoot)
	var accountRootCredentials aws.RootCredentials

	loginProfile, err := iamRoot.GetLoginProfile(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.AuditAccounts", "account %s - login_profile: %t", accountId, loginProfile)

	accessKeys, err := iamRoot.ListAccessKeys(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.AuditAccounts", "account %s - access_keys: %s", accountId, accessKeys)

	mfaDevices, err := iamRoot.ListMFADevices(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.AuditAccounts", "account %s - mfa_devices: %s", accountId, mfaDevices)

	certificates, err := iamRoot.ListSigningCertificates(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.AuditAccounts", "account %s - signing_certificates: %s", accountId, certificates)

	accountRootCredentials = aws.RootCredentials{
		AccountId:           accountId,
		LoginProfile:        loginProfile,
		AccessKeys:          accessKeys,
		MfaDevices:          mfaDevices,
		SigningCertificates: certificates,
	}

	return accountRootCredentials, nil
}
