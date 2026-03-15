package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete root user credentials",
		Long:  `Delete root user credentials for specific AWS Organization member accounts.`,
	}
	cmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of AWS account IDs to audit (comma-separated). Use \"all\" to audit all accounts.")
	cmd.AddCommand(deleteSubcommand("all", "Delete all existing root user credentials", "Delete all existing root user credentials for specific AWS Organization member accounts."))
	cmd.AddCommand(deleteSubcommand("login", "Delete root user Login Profile", "Delete existing root user Login Profile for specific AWS Organization member accounts."))
	cmd.AddCommand(deleteSubcommand("keys", "Delete root user Access Keys", "Delete existing root user Access Keys for specific AWS Organization member accounts."))
	cmd.AddCommand(deleteSubcommand("mfa", "Deactivate root user MFA Devices", "Deactivate existing root user MFA Devices for specific AWS Organization member accounts."))
	cmd.AddCommand(deleteSubcommand("certificates", "Delete root user Signin Certificates", "Delete existing root user Signing Certificates for specific AWS Organization member accounts."))
	return cmd
}

func deleteSubcommand(use, short, long string) *cobra.Command {
	credentialType := use
	if use == "certificates" {
		credentialType = "certificate"
	}
	return &cobra.Command{
		Use:          use,
		Short:        short,
		Long:         long,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(accountsFlags, credentialType)
		},
	}
}

func runDelete(accountsFlags []string, credentialType string) error {
	ctx := context.Background()
	rm, err := rootmanager.NewRootManager(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize root manager: %w", err)
	}

	auditAccounts, err := ui.SelectTargetAccounts(ctx, accountsFlags)
	if err != nil {
		return fmt.Errorf("failed to get accounts to audit: %w", err)
	}
	if len(auditAccounts) == 0 {
		slog.Info("no accounts selected")
		return nil
	}
	slog.Debug("selected accounts", "accounts", strings.Join(auditAccounts, ", "))

	audit, err := rm.AuditAccounts(ctx, auditAccounts)
	if err != nil {
		return err
	}

	results, err := rm.DeleteCredentials(ctx, audit, credentialType)
	if err != nil {
		return err
	}

	headers := []string{"Account", "CredentialType", "Status", "Error"}
	var data [][]any
	var failureCount int
	for _, result := range results {
		status := "deleted"
		errorMsg := ""
		if !result.Success {
			status = "failed"
			errorMsg = result.Error
			failureCount++
		}
		data = append(data, []any{
			result.AccountId,
			result.CredentialType,
			status,
			errorMsg,
		})
	}
	output.HandleOutput(outputFlag, headers, data)

	if failureCount > 0 {
		return fmt.Errorf("deletion failed for %d account(s)", failureCount)
	}

	return nil
}
