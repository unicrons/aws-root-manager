package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var accountsFlags []string

var rootCmd = &cobra.Command{
	Use:   "aws-root-manager",
	Short: "Manage your AWS Organization root access",
	Long: `aws-root-manager is a CLI tool for managing AWS Centralized Root Access.

This tool allows you to:
- ✅ Check if Centralized Root Access is enabled in your AWS Organization.
- 🔒 Enable Centralized Root Access for better security and control.
- 📊 Audit root access status across all organization accounts.
- 🗑️ Delete root credentials to enforce security best practices.

✨ More features coming soon!

🚀 Made with ❤️ by unicrons.cloud`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
