package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/internal/service"

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
	awscfg, err := aws.LoadAWSConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load aws config: %w", err)
	}

	auditAccounts, err := ui.SelectTargetAccounts(ctx, accountsFlags)
	if err != nil {
		return fmt.Errorf("failed to get accounts to audit: %w", err)
	}
	if len(auditAccounts) == 0 {
		logger.Info("cmd.audit", "no accounts selected")
		return nil
	}
	logger.Debug("cmd.audit", "selected accounts: %s", strings.Join(auditAccounts, ", "))

	rm := service.NewRootManager(aws.NewIamClient(awscfg), aws.NewStsClient(awscfg), aws.NewOrganizationsClient(awscfg))
	audit, err := rm.AuditAccounts(ctx, auditAccounts)
	if err != nil {
		return err
	}

	if err = rm.DeleteCredentials(ctx, audit, credentialType); err != nil {
		return err
	}

	headers := []string{"Account", "CredentialType", "Status"}
	var data [][]any
	for _, account := range auditAccounts {
		data = append(data, []any{
			account,
			credentialType,
			"deleted", // TODO: this is not real
		})
	}
	output.HandleOutput(outputFlag, headers, data)

	return nil
}
