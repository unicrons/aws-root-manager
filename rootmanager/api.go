// Package rootmanager provides operations for managing AWS root credentials
// across an AWS Organization using centralized root access (sts:AssumeRoot).
//
// Basic usage:
//
//	rm, err := rootmanager.NewRootManager(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	results, err := rm.AuditAccounts(ctx, []string{"123456789012"})
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
	// Returns a slice of DeletionResult showing the outcome for each account.
	DeleteCredentials(ctx context.Context, creds []RootCredentials, credentialType string) ([]DeletionResult, error)

	// RecoverRootPassword initiates root password recovery for the specified accounts.
	// This triggers AWS to send password reset emails to the account's root email address.
	// Returns a slice of RecoveryResult showing the outcome for each account.
	RecoverRootPassword(ctx context.Context, accountIds []string) ([]RecoveryResult, error)
}
