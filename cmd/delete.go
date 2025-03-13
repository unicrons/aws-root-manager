package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/unicrons/aws-root-manager/pkg/aws"
	"github.com/unicrons/aws-root-manager/pkg/logger"
	"github.com/unicrons/aws-root-manager/pkg/output"
	"github.com/unicrons/aws-root-manager/pkg/service"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete root user credentials",
	Long:  `Delete root user credentials for specific AWS Organization member accounts.`,
}

var deleteAllCmd = &cobra.Command{
	Use:          "all",
	Short:        "Delete all existing root user credentials",
	Long:         `Delete all existing root user credentials for specific AWS Organization member accounts.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Trace("cmd.deleteAll", "delete all called")

		if err := delete(accountsFlags, "all"); err != nil {
			return err
		}

		return nil
	},
}

var deleteLoginCmd = &cobra.Command{
	Use:          "login",
	Short:        "Delete root user Login Profile",
	Long:         `Delete existing root user Login Profile for specific AWS Organization member accounts.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Trace("cmd.deleteLogin", "delete login called")

		if err := delete(accountsFlags, "login"); err != nil {
			return err
		}

		return nil
	},
}

var deleteKeysCmd = &cobra.Command{
	Use:          "keys",
	Short:        "Delete root user Access Keys",
	Long:         `Delete existing root user Access Keys for specific AWS Organization member accounts.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Trace("cmd.deleteKeys", "delete keys called")

		if err := delete(accountsFlags, "keys"); err != nil {
			return err
		}

		return nil
	},
}

var deleteMfaCmd = &cobra.Command{
	Use:          "mfa",
	Short:        "Deactivate root user MFA Devices",
	Long:         `Deactivate existing root user MFA Devices for specific AWS Organization member accounts.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Trace("cmd.deleteMfa", "delete mfa called")

		if err := delete(accountsFlags, "mfa"); err != nil {
			return err
		}

		return nil
	},
}

var deleteCertificateCmd = &cobra.Command{
	Use:          "certificate",
	Short:        "Delete root user Signin Certificates",
	Long:         `Delete existing root user Signing Certificates for specific AWS Organization member accounts.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Trace("cmd.deleteCertificate", "delete certificate called")

		if err := delete(accountsFlags, "certificate"); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteAllCmd)
	deleteCmd.AddCommand(deleteLoginCmd)
	deleteCmd.AddCommand(deleteKeysCmd)
	deleteCmd.AddCommand(deleteMfaCmd)
	deleteCmd.AddCommand(deleteCertificateCmd)
	deleteCmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of AWS account IDs to audit (comma-separated). Use \"all\" to audit all accounts.")
}

func delete(accountsFlags []string, credentialType string) error {
	ctx := context.Background()
	awscfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load aws config: %w", err)
	}

	auditAccounts, err := service.GetTargetAccounts(ctx, accountsFlags)
	if err != nil {
		return fmt.Errorf("failed to get accounts to audit: %w", err)
	}
	if len(auditAccounts) == 0 {
		logger.Info("cmd.audit", "no accounts selected")
		return nil
	}
	logger.Debug("cmd.audit", "selected accounts: %s", strings.Join(auditAccounts, ", "))

	iam := aws.NewIamClient(awscfg)
	sts := aws.NewStsClient(awscfg)

	audit, err := service.AuditAccounts(ctx, iam, sts, auditAccounts)
	if err != nil {
		return err
	}

	if err = service.DeleteAccountsCredentials(ctx, iam, sts, audit, credentialType); err != nil {
		return err
	}

	headers := []string{"Account", "CredentialType", "Status"}
	var data [][]any
	for _, account := range auditAccounts {
		data = append(data, []any{
			account,
			credentialType,
			fmt.Sprintf("deleted"), // TODO: this is not real
		})
	}
	output.HandleOutput(outputFlag, headers, data)

	return nil
}
