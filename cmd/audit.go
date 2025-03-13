package cmd

import (
	"context"
	"strings"

	"github.com/unicrons/aws-root-manager/pkg/aws"
	"github.com/unicrons/aws-root-manager/pkg/logger"
	"github.com/unicrons/aws-root-manager/pkg/output"
	"github.com/unicrons/aws-root-manager/pkg/service"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:          "audit",
	Short:        "Retrieve root user credentials",
	Long:         `Retrieve available root user credentials for all member accounts within an AWS Organization.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("cmd.audit", "audit called")

		ctx := context.Background()
		awscfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			logger.Error("cmd.audit", err, "failed to load aws config")
			return
		}

		auditAccounts, err := service.GetTargetAccounts(ctx, accountsFlags)
		if err != nil {
			logger.Error("cmd.audit", err, "failed to get accounts to audit")
		}
		if len(auditAccounts) == 0 {
			logger.Info("cmd.audit", "no accounts selected")
			return
		}
		logger.Debug("cmd.audit", "selected accounts: %s", strings.Join(auditAccounts, ", "))

		iam := aws.NewIamClient(awscfg)
		sts := aws.NewStsClient(awscfg)
		audit, err := service.AuditAccounts(ctx, iam, sts, auditAccounts)
		if err != nil {
			logger.Error("cmd.audit", err, "failed to audit accounts")
			return
		}

		headers := []string{"Account", "LoginProfile", "AccessKeys", "MFA Devices", "Signing Certificates"}
		var data [][]any
		for i, acc := range audit {
			data = append(data, []any{
				auditAccounts[i],
				acc.LoginProfile,
				acc.AccessKeys,
				acc.MfaDevices,
				acc.SigningCertificates,
			})
		}
		output.HandleOutput(outputFlag, headers, data)

		return
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of AWS account IDs to audit (comma-separated). Use \"all\" to audit all accounts.")
}
