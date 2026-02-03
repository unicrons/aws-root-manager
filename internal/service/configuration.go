package service

import (
	"context"
	"errors"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func CheckRootAccess(ctx context.Context, iam *aws.IamClient) (aws.RootAccessStatus, error) {
	var status = aws.RootAccessStatus{
		TrustedAccess:             false,
		RootCredentialsManagement: false,
		RootSessions:              false,
	}

	err := iam.CheckOrganizationRootAccess(ctx, true)
	if err != nil {
		if errors.Is(err, rootmanager.ErrTrustedAccessNotEnabled) {
			return status, nil
		}
		status.TrustedAccess = true

		if errors.Is(err, rootmanager.ErrRootCredentialsManagementNotEnabled) {
			return status, nil
		}
		status.RootCredentialsManagement = true

		if errors.Is(err, rootmanager.ErrRootSessionsNotEnabled) {
			return status, nil
		}

		return aws.RootAccessStatus{}, err
	}

	status = aws.RootAccessStatus{
		TrustedAccess:             true,
		RootCredentialsManagement: true,
		RootSessions:              true,
	}

	return status, nil
}

func EnableRootAccess(ctx context.Context, iam *aws.IamClient, org *aws.OrganizationsClient, enableSessions bool) (aws.RootAccessStatus, aws.RootAccessStatus, error) {
	var initStatus, status aws.RootAccessStatus

	initStatus, err := CheckRootAccess(ctx, iam)
	if err != nil {
		return initStatus, status, err
	}

	if !initStatus.TrustedAccess {
		logger.Debug("service.EnableRootAccess", "trusted access is disabled")
		err := org.EnableAWSServiceAccess(ctx, "iam.amazonaws.com")
		if err != nil {
			return initStatus, status, err
		}
	}

	if !initStatus.RootCredentialsManagement {
		logger.Debug("service.EnableRootAccess", "root credentials management is disabled")
		err = iam.EnableOrganizationsRootCredentialsManagement(ctx)
		if err != nil {
			return initStatus, status, err
		}
	}

	if !initStatus.RootSessions && enableSessions {
		logger.Debug("service.EnableRootAccess", "root sessions is disabled")

		err = iam.EnableOrganizationsRootSessions(ctx)
		if err != nil {
			return initStatus, status, err
		}
	}

	status, err = CheckRootAccess(ctx, iam)
	if err != nil {
		return initStatus, status, err
	}

	return initStatus, status, nil
}
