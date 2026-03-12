package aws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type iamClient struct {
	client *iam.Client
}

func NewIamClient(awscfg aws.Config) IamClient {
	client := iam.NewFromConfig(awscfg)
	return &iamClient{client: client}
}

// Verifies if AWS centralized root access is enabled
func (c *iamClient) CheckOrganizationRootAccess(ctx context.Context, rootSessionsRequired bool) error {
	slog.Debug("checking if organization root access is enabled")

	features, err := c.client.ListOrganizationsFeatures(ctx, &iam.ListOrganizationsFeaturesInput{})
	if err != nil {
		var serviceAccessNotEnabledErr *types.ServiceAccessNotEnabledException
		if errors.As(err, &serviceAccessNotEnabledErr) {
			return rootmanager.ErrTrustedAccessNotEnabled
		}
		return fmt.Errorf("aws.CheckOrganizationRootAccess: failed to list organization features: %w", err)
	}

	rootCredentialsManagement := slices.Contains(features.EnabledFeatures, "RootCredentialsManagement")
	if !rootCredentialsManagement {
		return rootmanager.ErrRootCredentialsManagementNotEnabled
	}

	if !rootSessionsRequired {
		return nil
	}

	rootSessions := slices.Contains(features.EnabledFeatures, "RootSessions")
	if !rootSessions {
		return rootmanager.ErrRootSessionsNotEnabled
	}

	return nil
}

// Check if an account has root login profile enabled
func (c *iamClient) GetLoginProfile(ctx context.Context, accountId string) (bool, error) {
	slog.Debug("getting login profile", "account_id", accountId)

	_, err := c.client.GetLoginProfile(ctx, &iam.GetLoginProfileInput{})
	if err != nil {
		var notFoundErr *types.NoSuchEntityException
		if errors.As(err, &notFoundErr) {
			slog.Debug("account does not have a root login profile", "account_id", accountId)
			return false, nil
		}
		return true, fmt.Errorf("error getting root login profile for account %s: %w", accountId, err)
	}

	return true, nil
}

// Delete root login profile for a specific account
func (c *iamClient) DeleteLoginProfile(ctx context.Context, accountId string) error {
	slog.Debug("deleting login profile", "account_id", accountId)

	_, err := c.client.DeleteLoginProfile(ctx, &iam.DeleteLoginProfileInput{})
	if err != nil {
		return fmt.Errorf("error deleting root login profile for account %s: %w", accountId, err)
	}

	return nil
}

// Get a list of root access keys for a specific account
func (c *iamClient) ListAccessKeys(ctx context.Context, accountId string) ([]string, error) {
	slog.Debug("listing access keys", "account_id", accountId)

	accessKeys, err := c.client.ListAccessKeys(ctx, &iam.ListAccessKeysInput{})
	if err != nil {
		return nil, fmt.Errorf("error listing root access keys for account %s: %w", accountId, err)
	}

	// convert []AccessKeyMetadata to []string
	var accessKeyIDs []string
	for _, key := range accessKeys.AccessKeyMetadata {
		accessKeyIDs = append(accessKeyIDs, *key.AccessKeyId)
	}
	return accessKeyIDs, nil
}

// Delete a list of root access for a specific account
func (c *iamClient) DeleteAccessKeys(ctx context.Context, accountId string, accessKeyIds []string) error {
	slog.Debug("deleting root access keys", "account_id", accountId, "access_key_ids", accessKeyIds)

	for _, accessKeyId := range accessKeyIds {
		_, err := c.client.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(accessKeyId),
		})
		if err != nil {
			return fmt.Errorf("error deleting access key %s for account %s: %w", accessKeyId, accountId, err)
		}
	}

	return nil
}

// Get a list of root MFA devices for a specific account
func (c *iamClient) ListMFADevices(ctx context.Context, accountId string) ([]string, error) {
	mfaDevices, err := c.client.ListMFADevices(ctx, &iam.ListMFADevicesInput{})
	if err != nil {
		return nil, fmt.Errorf("error listing root mfa devices for account %s: %w", accountId, err)
	}

	// convert []MFADevices to []string
	var serialNumbers []string
	for _, device := range mfaDevices.MFADevices {
		serialNumbers = append(serialNumbers, *device.SerialNumber)
	}
	return serialNumbers, nil
}

// Deactivate a list of root MFA devices for a specific account
func (c *iamClient) DeactivateMFADevices(ctx context.Context, accountId string, mfaSerialNumbers []string) error {
	slog.Debug("deactivating root mfa devices", "account_id", accountId, "mfa_serial_numbers", mfaSerialNumbers)

	for _, mfaSerialNumber := range mfaSerialNumbers {
		_, err := c.client.DeactivateMFADevice(ctx, &iam.DeactivateMFADeviceInput{
			SerialNumber: aws.String(mfaSerialNumber),
		})
		if err != nil {
			return fmt.Errorf("error deleting root mfa device %s for account %s: %w", mfaSerialNumber, accountId, err)
		}
	}

	return nil
}

// Get a list of root signing certificates devices for a specific account
func (c *iamClient) ListSigningCertificates(ctx context.Context, accountId string) ([]string, error) {
	certificates, err := c.client.ListSigningCertificates(ctx, &iam.ListSigningCertificatesInput{})
	if err != nil {
		return nil, fmt.Errorf("error listing signing certificates for account %s: %w", accountId, err)
	}

	// convert []SigningCertificate to []string
	var certificateIDs []string
	for _, certificate := range certificates.Certificates {
		certificateIDs = append(certificateIDs, *certificate.CertificateId)
	}
	return certificateIDs, nil
}

// Delete a list of root signing certificates for a specific account
func (c *iamClient) DeleteSigningCertificates(ctx context.Context, accountId string, certificates []string) error {
	slog.Debug("deleting signing certificates", "account_id", accountId, "certificates", certificates)

	for _, certificate := range certificates {
		if _, err := c.client.DeleteSigningCertificate(ctx, &iam.DeleteSigningCertificateInput{
			CertificateId: aws.String(certificate),
		}); err != nil {
			return err
		}
	}

	return nil
}

// Enable centralized root credentials management
func (c *iamClient) EnableOrganizationsRootCredentialsManagement(ctx context.Context) error {
	slog.Debug("enabling organization root credentials management")

	_, err := c.client.EnableOrganizationsRootCredentialsManagement(ctx, &iam.EnableOrganizationsRootCredentialsManagementInput{})
	if err != nil {
		return fmt.Errorf("error enabling organization root credentials management: %w", err)
	}

	return nil
}

// Enable centralized root sessions
func (c *iamClient) EnableOrganizationsRootSessions(ctx context.Context) error {
	slog.Debug("enabling organization root sessions")

	_, err := c.client.EnableOrganizationsRootSessions(ctx, &iam.EnableOrganizationsRootSessionsInput{})
	if err != nil {
		return fmt.Errorf("error enabling organization root sessions: %w", err)
	}

	return nil
}

// Allow root password recovery
func (c *iamClient) CreateLoginProfile(ctx context.Context) error {
	slog.Debug("creating login profile")

	_, err := c.client.CreateLoginProfile(ctx, &iam.CreateLoginProfileInput{})
	if err != nil {
		var entityAlreadyExistsErr *types.EntityAlreadyExistsException
		if errors.As(err, &entityAlreadyExistsErr) {
			slog.Debug("login profile already exists")
			return rootmanager.ErrEntityAlreadyExists
		}
		return fmt.Errorf("error creating login profile: %w", err)
	}

	return nil
}
