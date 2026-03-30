package rootmanager

import (
	"context"
	"errors"
	"log/slog"

	"github.com/unicrons/aws-root-manager/internal/aws"
)

func checkRootAccess(ctx context.Context, iam aws.IamClient) (RootAccessStatus, error) {
	var status = RootAccessStatus{
		TrustedAccess:             false,
		RootCredentialsManagement: false,
		RootSessions:              false,
	}

	err := iam.CheckOrganizationRootAccess(ctx, true)
	if err != nil {
		if errors.Is(err, ErrTrustedAccessNotEnabled) {
			return status, nil
		}
		status.TrustedAccess = true

		if errors.Is(err, ErrRootCredentialsManagementNotEnabled) {
			return status, nil
		}
		status.RootCredentialsManagement = true

		if errors.Is(err, ErrRootSessionsNotEnabled) {
			return status, nil
		}

		return RootAccessStatus{}, err
	}

	status = RootAccessStatus{
		TrustedAccess:             true,
		RootCredentialsManagement: true,
		RootSessions:              true,
	}

	return status, nil
}

func enableRootAccess(ctx context.Context, iam aws.IamClient, org aws.OrganizationsClient, enableSessions bool) (RootAccessStatus, RootAccessStatus, error) {
	var initStatus, status RootAccessStatus

	initStatus, err := checkRootAccess(ctx, iam)
	if err != nil {
		return initStatus, status, err
	}

	if !initStatus.TrustedAccess {
		slog.Debug("trusted access is disabled")
		err := org.EnableAWSServiceAccess(ctx, "iam.amazonaws.com")
		if err != nil {
			return initStatus, status, err
		}
	}

	if !initStatus.RootCredentialsManagement {
		slog.Debug("root credentials management is disabled")
		err = iam.EnableOrganizationsRootCredentialsManagement(ctx)
		if err != nil {
			return initStatus, status, err
		}
	}

	if !initStatus.RootSessions && enableSessions {
		slog.Debug("root sessions is disabled")

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
