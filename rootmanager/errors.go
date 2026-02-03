package rootmanager

import "errors"

var (
	// ErrTrustedAccessNotEnabled indicates AWS IAM does not have trusted access to the organization.
	ErrTrustedAccessNotEnabled = errors.New("AWS IAM trusted access is not enabled for the organization")

	// ErrRootCredentialsManagementNotEnabled indicates centralized root credentials management is not enabled.
	ErrRootCredentialsManagementNotEnabled = errors.New("centralized root credentials management is not enabled")

	// ErrRootSessionsNotEnabled indicates root sessions (AssumeRoot) are not enabled.
	ErrRootSessionsNotEnabled = errors.New("root sessions are not enabled")

	// ErrEntityAlreadyExists indicates the requested entity already exists.
	ErrEntityAlreadyExists = errors.New("entity already exists")
)
