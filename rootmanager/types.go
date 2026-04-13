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
	Error     string // Error message if recovery failed (empty if Success=true)
}

// DeletionResult represents the result of a credential deletion operation for an account.
type DeletionResult struct {
	AccountId      string // AWS account ID
	CredentialType string // Type of credential deleted (login, keys, mfa, certificate, all)
	Success        bool   // Whether deletion was successful
	Error          string // Error message if deletion failed (empty if Success=true)
}

// PolicyDeletionResult represents the result of a resource policy deletion operation.
type PolicyDeletionResult struct {
	AccountId    string // AWS account ID
	ResourceType string // Type of resource ("s3-bucket", "sqs-queue")
	ResourceName string // Bucket name or queue URL
	Success      bool   // Whether deletion was successful
	Error        string // Error message if deletion failed (empty if Success=true)
}
