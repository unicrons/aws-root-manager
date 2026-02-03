package rootmanager

import "context"

// RootManager provides operations for managing AWS root credentials for an AWS organization.
type RootManager interface {
	// AuditAccounts audits root credentials across the specified AWS accounts.
	// It checks for login profiles, access keys, MFA devices, and signing certificates.
	AuditAccounts(ctx context.Context, accountIds []string) ([]RootCredentials, error)

	// CheckRootAccess checks the status of centralized root access features in the organization.
	// It verifies whether trusted access, root credentials management, and root sessions are enabled.
	CheckRootAccess(ctx context.Context) (RootAccessStatus, error)

	// EnableRootAccess enables centralized root access features in the organization.
	// The enableSessions parameter controls whether to enable root sessions (AssumeRoot).
	// Returns the initial status, final status after enabling, and any error encountered.
	EnableRootAccess(ctx context.Context, enableSessions bool) (RootAccessStatus, RootAccessStatus, error)

	// DeleteCredentials deletes root credentials for the specified accounts.
	// The creds parameter should contain audit results identifying what credentials exist.
	// The credentialType parameter specifies what to delete: "all", "login", "keys", "mfa", or "certificate".
	DeleteCredentials(ctx context.Context, creds []RootCredentials, credentialType string) error

	// RecoverRootPassword initiates root password recovery for the specified accounts.
	// This triggers AWS to send password reset emails to the account's root email address.
	// Returns a map of account ID to success status, and any error encountered.
	RecoverRootPassword(ctx context.Context, accountIds []string) (map[string]bool, error)
}
