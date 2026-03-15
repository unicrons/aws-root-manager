package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/unicrons/aws-root-manager/internal/aws"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

// deleteAccountsCredentials deletes root credentials for a list of AWS accounts.
// Returns a slice of DeletionResult containing the outcome for each account.
func deleteAccountsCredentials(ctx context.Context, iam aws.IamClient, sts aws.StsClient, factory aws.IamClientFactory, creds []rootmanager.RootCredentials, credentialType string) ([]rootmanager.DeletionResult, error) {
	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return nil, err
	}

	results := make([]rootmanager.DeletionResult, len(creds))
	var wgAccounts sync.WaitGroup

	for i, accountCredentials := range creds {
		wgAccounts.Add(1)
		go func(idx int, accountCreds rootmanager.RootCredentials) {
			defer wgAccounts.Done()
			if err := deleteAccountCredentials(ctx, sts, factory, accountCreds, credentialType); err != nil {
				results[idx] = rootmanager.DeletionResult{
					AccountId:      accountCreds.AccountId,
					CredentialType: credentialType,
					Success:        false,
					Error:          err.Error(),
				}
			} else {
				results[idx] = rootmanager.DeletionResult{
					AccountId:      accountCreds.AccountId,
					CredentialType: credentialType,
					Success:        true,
					Error:          "",
				}
			}
		}(i, accountCredentials)
	}

	wgAccounts.Wait()

	return results, nil
}

// deleteAccountCredentials deletes root credentials for a specific account.
func deleteAccountCredentials(ctx context.Context, sts aws.StsClient, factory aws.IamClientFactory, creds rootmanager.RootCredentials, credentialType string) error {
	slog.Debug("checking credentials to delete", "account_id", creds.AccountId, "credential_type", credentialType)

	// Check if there are credentials to delete before assuming root
	if !hasCredentialsToDelete(creds, credentialType) {
		return nil
	}

	awscfgDeleteRoot, err := sts.GetAssumeRootConfig(ctx, creds.AccountId, "IAMDeleteRootUserCredentials")
	if err != nil {
		return err
	}
	iamDeleteRoot := factory.NewIamClient(awscfgDeleteRoot)

	if creds.LoginProfile && (credentialType == "all" || credentialType == "login") {
		err = iamDeleteRoot.DeleteLoginProfile(ctx, creds.AccountId)
		if err != nil {
			return err
		}
	}

	if len(creds.AccessKeys) > 0 && (credentialType == "all" || credentialType == "keys") {
		err = iamDeleteRoot.DeleteAccessKeys(ctx, creds.AccountId, creds.AccessKeys)
		if err != nil {
			return err
		}
	}

	if len(creds.MfaDevices) > 0 && (credentialType == "all" || credentialType == "mfa") {
		err = iamDeleteRoot.DeactivateMFADevices(ctx, creds.AccountId, creds.MfaDevices)
		if err != nil {
			return err
		}
	}

	if len(creds.SigningCertificates) > 0 && (credentialType == "all" || credentialType == "certificate") {
		err = iamDeleteRoot.DeleteSigningCertificates(ctx, creds.AccountId, creds.SigningCertificates)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check if the account has credentials to delete based on the credential type
func hasCredentialsToDelete(creds rootmanager.RootCredentials, credentialType string) bool {
	switch credentialType {
	case "all":
		return creds.LoginProfile || len(creds.AccessKeys) > 0 || len(creds.MfaDevices) > 0 || len(creds.SigningCertificates) > 0
	case "login":
		return creds.LoginProfile
	case "keys":
		return len(creds.AccessKeys) > 0
	case "mfa":
		return len(creds.MfaDevices) > 0
	case "certificate":
		return len(creds.SigningCertificates) > 0
	default:
		return false
	}
}

// recoverAccountsRootPassword initiates root password recovery for a list of AWS accounts.
// Returns a slice of RecoveryResult containing the outcome for each account.
func recoverAccountsRootPassword(ctx context.Context, iam aws.IamClient, sts aws.StsClient, factory aws.IamClientFactory, accountIds []string) ([]rootmanager.RecoveryResult, error) {
	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return nil, err
	}

	results := make([]rootmanager.RecoveryResult, len(accountIds))
	var wgAccounts sync.WaitGroup

	for i, accountId := range accountIds {
		wgAccounts.Add(1)
		go func(idx int, accId string) {
			defer wgAccounts.Done()
			success, err := recoverAccountRootPassowrd(ctx, sts, factory, accId)
			if err != nil {
				results[idx] = rootmanager.RecoveryResult{
					AccountId: accId,
					Success:   false,
					Error:     err.Error(),
				}
			} else {
				results[idx] = rootmanager.RecoveryResult{
					AccountId: accId,
					Success:   success,
					Error:     "",
				}
			}
		}(i, accountId)
	}

	wgAccounts.Wait()

	return results, nil
}

// Enable the recovery process for root passwords for a specific account
func recoverAccountRootPassowrd(ctx context.Context, sts aws.StsClient, factory aws.IamClientFactory, accountId string) (bool, error) {
	slog.Debug("trying to recover root password", "account_id", accountId)

	awscfgRecoverRoot, err := sts.GetAssumeRootConfig(ctx, accountId, "IAMCreateRootUserPassword")
	if err != nil {
		return false, err
	}
	iamRecoverRoot := factory.NewIamClient(awscfgRecoverRoot)

	err = iamRecoverRoot.CreateLoginProfile(ctx)
	if err != nil {
		if errors.Is(err, rootmanager.ErrEntityAlreadyExists) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
