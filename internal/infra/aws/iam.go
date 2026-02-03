package aws

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IamClient struct {
	client *iam.Client
}

func NewIamClient(awscfg aws.Config) *IamClient {
	client := iam.NewFromConfig(awscfg)
	return &IamClient{client: client}
}

// Verifies if AWS centralized root access is enabled
func (c *IamClient) CheckOrganizationRootAccess(ctx context.Context, rootSessionsRequired bool) error {
	logger.Trace("aws.CheckOrganizationRootAccess", "checking if organization root access is enabled")

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
func (c *IamClient) GetLoginProfile(ctx context.Context, accountId string) (bool, error) {
	logger.Debug("aws.GetLoginProfile", "getting login profile for account %s", accountId)

	_, err := c.client.GetLoginProfile(ctx, &iam.GetLoginProfileInput{})
	if err != nil {
		var notFoundErr *types.NoSuchEntityException
		if errors.As(err, &notFoundErr) {
			logger.Debug("aws.GetLoginProfile", "account %s does not have a root login profile", accountId)
			return false, nil
		}
		return true, fmt.Errorf("error getting root login profile for account %s: %w", accountId, err)
	}

	return true, nil
}

// Delete root login profile for a specific account
func (c *IamClient) DeleteLoginProfile(ctx context.Context, accountId string) error {
	logger.Debug("aws.DeleteLoginProfile", "deleting login profile for account %s", accountId)

	_, err := c.client.DeleteLoginProfile(ctx, &iam.DeleteLoginProfileInput{})
	if err != nil {
		return fmt.Errorf("error deleting root login profile for account %s: %w", accountId, err)
	}

	logger.Info("aws.DeleteLoginProfile", "successfully deleted login profile for account %s", accountId)

	return nil
}

// Get a list of root access keys for a specific account
func (c *IamClient) ListAccessKeys(ctx context.Context, accountId string) ([]string, error) {
	logger.Debug("aws.ListAccessKeys", "listing access keys for account %s", accountId)

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
func (c *IamClient) DeleteAccessKeys(ctx context.Context, accountId string, accessKeyIds []string) error {
	logger.Debug("aws.DeleteAccessKeys", "deleting root access key %s for account %s", accessKeyIds, accountId)

	for _, accessKeyId := range accessKeyIds {
		_, err := c.client.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(accessKeyId),
		})
		if err != nil {
			return fmt.Errorf("error deleting access key %s for account %s: %w", accessKeyId, accountId, err)
		}
	}

	logger.Info("aws.DeleteAccessKeys", "successfully deleted access keys for account %s", accountId)

	return nil
}

// Get a list of root MFA devices for a specific account
func (c *IamClient) ListMFADevices(ctx context.Context, accountId string) ([]string, error) {
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
func (c *IamClient) DeactivateMFADevices(ctx context.Context, accountId string, mfaSerialNumbers []string) error {
	logger.Debug("aws.DeactivateMFADevices", "deleting root mfa device %s for account %s", mfaSerialNumbers, accountId)

	for _, mfaSerialNumber := range mfaSerialNumbers {
		_, err := c.client.DeactivateMFADevice(ctx, &iam.DeactivateMFADeviceInput{
			SerialNumber: aws.String(mfaSerialNumber),
		})
		if err != nil {
			return fmt.Errorf("error deleting root mfa device %s for account %s: %w", mfaSerialNumber, accountId, err)
		}
	}

	logger.Info("aws.DeactivateMFADevices", "successfully deactivated mfa devices for account %s", accountId)

	return nil
}

// Get a list of root signing certificates devices for a specific account
func (c *IamClient) ListSigningCertificates(ctx context.Context, accountId string) ([]string, error) {
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
func (c *IamClient) DeleteSigningCertificates(ctx context.Context, accountId string, certificates []string) error {
	logger.Debug("aws.DeleteSigningCertificates", "deleting singin certificates %s for account %s", certificates, accountId)

	for _, certificate := range certificates {
		logger.Debug("aws.DeleteSigningCertificates", "deleting root Signing Certificate %s", certificate)

		if _, err := c.client.DeleteSigningCertificate(ctx, &iam.DeleteSigningCertificateInput{
			CertificateId: aws.String(certificate),
		}); err != nil {
			return err
		}
	}

	logger.Info("aws.DeleteSigningCertificates", "successfully deleted signing certificates for account %s", accountId)

	return nil
}

// Enable centralized root credentials management
func (c *IamClient) EnableOrganizationsRootCredentialsManagement(ctx context.Context) error {
	logger.Debug("aws.EnableOrganizationsRootCredentialsManagement", "enabling organization root credentials management")

	_, err := c.client.EnableOrganizationsRootCredentialsManagement(ctx, &iam.EnableOrganizationsRootCredentialsManagementInput{})
	if err != nil {
		return fmt.Errorf("error enabling organization root credentials management: %w", err)
	}

	logger.Info("aws.EnableOrganizationsRootCredentialsManagement", "successfully enabled organization root credentials management")

	return nil
}

// Enable centralized root sessions
func (c *IamClient) EnableOrganizationsRootSessions(ctx context.Context) error {
	logger.Debug("aws.EnableOrganizationsRootSessions", "enabling organization root sessions")

	_, err := c.client.EnableOrganizationsRootSessions(ctx, &iam.EnableOrganizationsRootSessionsInput{})
	if err != nil {
		return fmt.Errorf("error enabling organization root sessions: %w", err)
	}

	logger.Info("aws.EnableOrganizationsRootSessions", "successfully enabled organization root sessions management")

	return nil
}

// Allow root password recovery
func (c *IamClient) CreateLoginProfile(ctx context.Context) error {
	logger.Debug("aws.createLoginProfile", "creating loggin profile")

	_, err := c.client.CreateLoginProfile(ctx, &iam.CreateLoginProfileInput{})
	if err != nil {
		var entityAlreadyExistsErr *types.EntityAlreadyExistsException
		if errors.As(err, &entityAlreadyExistsErr) {
			logger.Debug("aws.createLoginProfile", "login profile already exists")
			return rootmanager.ErrEntityAlreadyExists
		}
		return fmt.Errorf("error creating login profile: %w", err)
	}

	logger.Info("aws.createLoginProfile", "successfully created login profile")

	return nil
}
