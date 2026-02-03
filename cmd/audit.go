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

func Audit() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "audit",
		Short:        "Retrieve root user credentials",
		Long:         `Retrieve available root user credentials for all member accounts within an AWS Organization.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Trace("cmd.audit", "audit called")

			ctx := context.Background()
			awscfg, err := aws.LoadAWSConfig(ctx)
			if err != nil {
				logger.Error("cmd.audit", err, "failed to load aws config")
				return err
			}

			auditAccounts, err := ui.SelectTargetAccounts(ctx, accountsFlags)
			if err != nil {
				logger.Error("cmd.audit", err, "failed to get accounts to audit")
				return err
			}
			if len(auditAccounts) == 0 {
				logger.Info("cmd.audit", "no accounts selected")
				return nil
			}
			logger.Debug("cmd.audit", "selected accounts: %s", strings.Join(auditAccounts, ", "))

			rm := service.NewRootManager(aws.NewIamClient(awscfg), aws.NewStsClient(awscfg), aws.NewOrganizationsClient(awscfg))
			audit, err := rm.AuditAccounts(ctx, auditAccounts)
			if err != nil {
				logger.Error("cmd.audit", err, "failed to audit accounts")
				return err
			}

			var skipped int
			headers := []string{"Account", "LoginProfile", "AccessKeys", "MFA Devices", "Signing Certificates"}
			var data [][]any
			for i, acc := range audit {
				if acc.Error != "" {
					skipped++
					continue
				}
				data = append(data, []any{
					auditAccounts[i],
					acc.LoginProfile,
					acc.AccessKeys,
					acc.MfaDevices,
					acc.SigningCertificates,
				})
			}
			output.HandleOutput(outputFlag, headers, data)

			if skipped > 0 {
				return fmt.Errorf("audit skipped for %d account(s)", skipped)
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of AWS account IDs to audit (comma-separated). Use \"all\" to audit all accounts.")
	return cmd
}
