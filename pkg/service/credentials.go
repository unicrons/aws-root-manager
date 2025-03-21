package service

import (
	"context"
	"sync"

	"github.com/unicrons/aws-root-manager/pkg/aws"
	"github.com/unicrons/aws-root-manager/pkg/logger"
)

// Delete root credentials for a list of AWS accounts
func DeleteAccountsCredentials(ctx context.Context, iam *aws.IamClient, sts *aws.StsClient, creds []aws.RootCredentials, credentialType string) error {
	var (
		wgAccounts sync.WaitGroup
		errChan    = make(chan error, len(creds))
	)

	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return err
	}

	for _, accountCredentials := range creds {
		wgAccounts.Add(1)
		go func(accountId aws.RootCredentials) {
			defer wgAccounts.Done()
			if err := deleteAccountCrendentials(ctx, sts, accountCredentials, credentialType); err != nil {
				errChan <- err
			}
		}(accountCredentials)
	}

	wgAccounts.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

// Delete root credentials for a specific account
func deleteAccountCrendentials(ctx context.Context, sts *aws.StsClient, creds aws.RootCredentials, credentialType string) error {
	logger.Trace("service.deleteAccountCrendentials", "checking if account %s has %s credentials to delete", credentialType, credentialType)

	// Check if there are credentials to delete before assuming root
	if !hasCredentialsToDelete(creds, credentialType) {
		logger.Info("service.deleteAccountCrendentials", "no %s credentials found for account %s", credentialType, creds.AccountId)
		return nil
	}

	awscfgDeleteRoot, err := sts.GetAssumeRootConfig(ctx, creds.AccountId, "IAMDeleteRootUserCredentials")
	if err != nil {
		return err
	}
	iamDeleteRoot := aws.NewIamClient(awscfgDeleteRoot)

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
func hasCredentialsToDelete(creds aws.RootCredentials, credentialType string) bool {
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

// Enable the recovery process for root passwords for a list of AWS accounts
func RecoverAccountsRootPassword(ctx context.Context, iam *aws.IamClient, sts *aws.StsClient, accountIds []string) (map[string]bool, error) {
	var (
		wgAccounts sync.WaitGroup
		results    = sync.Map{}
		errChan    = make(chan error, len(accountIds))
	)

	if err := iam.CheckOrganizationRootAccess(ctx, false); err != nil {
		return nil, err
	}

	for _, acc := range accountIds {
		wgAccounts.Add(1)
		go func(accountId string) {
			defer wgAccounts.Done()
			success, err := recoverAccountRootPassowrd(ctx, sts, acc)
			results.Store(accountId, success)
			if err != nil {
				errChan <- err
			}
		}(acc)
	}

	wgAccounts.Wait()
	close(errChan)

	resultMap := make(map[string]bool)
	results.Range(func(key, value any) bool {
		resultMap[key.(string)] = value.(bool)
		return true
	})

	if len(errChan) > 0 {
		return resultMap, <-errChan
	}

	return resultMap, nil
}

// Enable the recovery process for root passwords for a specific account
func recoverAccountRootPassowrd(ctx context.Context, sts *aws.StsClient, accountId string) (bool, error) {
	logger.Trace("service.recoverAccountRootPassowrd", "trying to recover root password for account %s ", accountId)

	awscfgRecoverRoot, err := sts.GetAssumeRootConfig(ctx, accountId, "IAMCreateRootUserPassword")
	if err != nil {
		return false, err
	}
	iamRecoverRoot := aws.NewIamClient(awscfgRecoverRoot)

	err = iamRecoverRoot.CreateLoginProfile(ctx)
	if err != nil {
		if err == aws.ErrEntityAlreadyExists {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
