package rootmanager

// RootAccessStatus represents the status of centralized root access features in an AWS Organization.
type RootAccessStatus struct {
	TrustedAccess             bool // Whether AWS IAM has trusted access to the organization
	RootCredentialsManagement bool // Whether centralized root credentials management is enabled
	RootSessions              bool // Whether root sessions (assume root) are enabled
}

// RootCredentials represents the root user credentials for an AWS account.
type RootCredentials struct {
	AccountId           string   // AWS account ID
	LoginProfile        bool     // Whether a root password exists
	AccessKeys          []string // List of root access key IDs
	MfaDevices          []string // List of root MFA device serial numbers
	SigningCertificates []string // List of root signing certificate IDs
	Error               string   // Error message if audit failed for this account
}

// RecoveryResult represents the result of a root password recovery operation for an account.
type RecoveryResult struct {
	AccountId string // AWS account ID
	Success   bool   // Whether recovery email was successfully sent
}
