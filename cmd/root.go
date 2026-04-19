package cmd

import (
	"os"

	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/spf13/cobra"
)

var (
	accountsFlags []string
	outputFlag    string
	skipFlag      bool
)

var rootCmd = &cobra.Command{
	Use:   "aws-root-manager",
	Short: "Manage your AWS Organization root access",
	Long: `aws-root-manager is a CLI tool for managing AWS Centralized Root Access.

This tool allows you to:
- ✅ Check if Centralized Root Access is enabled in your AWS Organization.
- 🔒 Enable Centralized Root Access for better security and control.
- 📊 Audit root access status across all organization accounts.
- 🗑️ Delete root credentials to enforce security best practices.
- 🛠️ Allow root password recovery.

✨ More features coming soon!

Made with ❤️ by unicrons.cloud 🦄`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logger.Configure(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_FORMAT"))

	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "table", "Set the output format (table, json, csv)")
	rootCmd.AddCommand(Audit(rootmanager.NewRootManager))
	rootCmd.AddCommand(Check(rootmanager.NewRootManager))
	rootCmd.AddCommand(Enable(rootmanager.NewRootManager))
	rootCmd.AddCommand(Delete(rootmanager.NewRootManager))
	rootCmd.AddCommand(Recovery(rootmanager.NewRootManager))
	rootCmd.AddCommand(Version())
}
