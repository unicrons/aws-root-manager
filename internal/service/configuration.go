package service

import (
	"context"
	"errors"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func checkRootAccess(ctx context.Context, iam *aws.IamClient) (rootmanager.RootAccessStatus, error) {
	var status = rootmanager.RootAccessStatus{
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

		return rootmanager.RootAccessStatus{}, err
	}

	status = rootmanager.RootAccessStatus{
		TrustedAccess:             true,
		RootCredentialsManagement: true,
		RootSessions:              true,
	}

	return status, nil
}

func enableRootAccess(ctx context.Context, iam *aws.IamClient, org *aws.OrganizationsClient, enableSessions bool) (rootmanager.RootAccessStatus, rootmanager.RootAccessStatus, error) {
	var initStatus, status rootmanager.RootAccessStatus

	initStatus, err := checkRootAccess(ctx, iam)
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

	status, err = checkRootAccess(ctx, iam)
	if err != nil {
		return initStatus, status, err
	}

	return initStatus, status, nil
}
