package service

import (
	"context"
	"sync"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

// auditAccounts returns root credentials for a list of AWS accounts.
func auditAccounts(ctx context.Context, iam *aws.IamClient, sts *aws.StsClient, accounts []string) ([]rootmanager.RootCredentials, error) {
	logger.Trace("service.auditAccounts", "auditing accounts %s", accounts)

	rootCredentials := make([]rootmanager.RootCredentials, len(accounts))
	var wgAccounts sync.WaitGroup

	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return nil, err
	}

	for i, accountId := range accounts {
		wgAccounts.Add(1)
		go func(idx int, accountId string) {
			defer wgAccounts.Done()
			if accStatus, err := auditAccount(ctx, sts, accountId); err != nil {
				logger.Error("service.auditAccounts", err, "account %s: audit skipped", accountId)
				rootCredentials[idx] = rootmanager.RootCredentials{AccountId: accountId, Error: err.Error()}
			} else {
				rootCredentials[idx] = accStatus
			}
		}(i, accountId)
	}

	wgAccounts.Wait()

	return rootCredentials, nil
}

// Get root credentials for a specific account
func auditAccount(ctx context.Context, sts *aws.StsClient, accountId string) (rootmanager.RootCredentials, error) {
	logger.Trace("service.auditAccount", "auditing account %s", accountId)

	awscfgRoot, err := sts.GetAssumeRootConfig(ctx, accountId, "IAMAuditRootUserCredentials")
	if err != nil {
		return rootmanager.RootCredentials{}, err
	}

	iamRoot := aws.NewIamClient(awscfgRoot)
	var accountRootCredentials rootmanager.RootCredentials

	loginProfile, err := iamRoot.GetLoginProfile(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.auditAccounts", "account %s - login_profile: %t", accountId, loginProfile)

	accessKeys, err := iamRoot.ListAccessKeys(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.auditAccounts", "account %s - access_keys: %s", accountId, accessKeys)

	mfaDevices, err := iamRoot.ListMFADevices(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.auditAccounts", "account %s - mfa_devices: %s", accountId, mfaDevices)

	certificates, err := iamRoot.ListSigningCertificates(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	logger.Debug("service.auditAccounts", "account %s - signing_certificates: %s", accountId, certificates)

	accountRootCredentials = rootmanager.RootCredentials{
		AccountId:           accountId,
		LoginProfile:        loginProfile,
		AccessKeys:          accessKeys,
		MfaDevices:          mfaDevices,
		SigningCertificates: certificates,
	}

	return accountRootCredentials, nil
}
