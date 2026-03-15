package rootmanager

import internalaws "github.com/unicrons/aws-root-manager/internal/aws"

var (
	// ErrTrustedAccessNotEnabled indicates AWS IAM does not have trusted access to the organization.
	ErrTrustedAccessNotEnabled = internalaws.ErrTrustedAccessNotEnabled

	// ErrRootCredentialsManagementNotEnabled indicates centralized root credentials management is not enabled.
	ErrRootCredentialsManagementNotEnabled = internalaws.ErrRootCredentialsManagementNotEnabled

	// ErrRootSessionsNotEnabled indicates root sessions (AssumeRoot) are not enabled.
	ErrRootSessionsNotEnabled = internalaws.ErrRootSessionsNotEnabled

	// ErrEntityAlreadyExists indicates the requested entity already exists.
	ErrEntityAlreadyExists = internalaws.ErrEntityAlreadyExists
)
