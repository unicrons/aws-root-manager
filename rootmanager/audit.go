package rootmanager

import (
	"context"
	"log/slog"
	"sync"

	"github.com/unicrons/aws-root-manager/internal/aws"
)

// auditAccounts returns root credentials for a list of AWS accounts.
func auditAccounts(ctx context.Context, iam aws.IamClient, sts aws.StsClient, factory aws.IamClientFactory, accounts []string) ([]RootCredentials, error) {
	slog.Debug("auditing accounts", "accounts", accounts)

	rootCredentials := make([]RootCredentials, len(accounts))
	var wgAccounts sync.WaitGroup

	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return nil, err
	}

	for i, accountId := range accounts {
		wgAccounts.Add(1)
		go func(idx int, accountId string) {
			defer wgAccounts.Done()
			if accStatus, err := auditAccount(ctx, sts, factory, accountId); err != nil {
				rootCredentials[idx] = RootCredentials{AccountId: accountId, Error: err.Error()}
			} else {
				rootCredentials[idx] = accStatus
			}
		}(i, accountId)
	}

	wgAccounts.Wait()

	return rootCredentials, nil
}

// auditAccount returns root credentials for a specific account.
func auditAccount(ctx context.Context, sts aws.StsClient, factory aws.IamClientFactory, accountId string) (RootCredentials, error) {
	slog.Debug("auditing account", "account_id", accountId)

	awscfgRoot, err := sts.GetAssumeRootConfig(ctx, accountId, "IAMAuditRootUserCredentials")
	if err != nil {
		return RootCredentials{}, err
	}

	iamRoot := factory.NewIamClient(awscfgRoot)
	var accountRootCredentials RootCredentials

	loginProfile, err := iamRoot.GetLoginProfile(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	slog.Debug("audit result", "account_id", accountId, "login_profile", loginProfile)

	accessKeys, err := iamRoot.ListAccessKeys(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	slog.Debug("audit result", "account_id", accountId, "access_keys", accessKeys)

	mfaDevices, err := iamRoot.ListMFADevices(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	slog.Debug("audit result", "account_id", accountId, "mfa_devices", mfaDevices)

	certificates, err := iamRoot.ListSigningCertificates(ctx, accountId)
	if err != nil {
		return accountRootCredentials, err
	}
	slog.Debug("audit result", "account_id", accountId, "signing_certificates", certificates)

	accountRootCredentials = RootCredentials{
		AccountId:           accountId,
		LoginProfile:        loginProfile,
		AccessKeys:          accessKeys,
		MfaDevices:          mfaDevices,
		SigningCertificates: certificates,
	}

	return accountRootCredentials, nil
}
