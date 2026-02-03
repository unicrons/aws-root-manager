package cmd

import (
	"context"
	"strings"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/internal/service"

	"github.com/spf13/cobra"
)

func Recovery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery",
		Short: "Allow root password recovery",
		Long:  `Retrieve the status of centralized root access settings for an AWS Organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Trace("cmd.recovery", "recovery called")

			ctx := context.Background()
			awscfg, err := aws.LoadAWSConfig(ctx)
			if err != nil {
				logger.Error("cmd.recovery", err, "failed to load aws config")
				return err
			}

			targetAccounts, err := ui.SelectTargetAccounts(ctx, accountsFlags)
			if err != nil {
				logger.Error("cmd.recovery", err, "failed to get target accounts")
				return err
			}
			if len(targetAccounts) == 0 {
				logger.Info("cmd.recovery", "no accounts selected")
				return nil
			}
			logger.Debug("cmd.recovery", "selected accounts: %s", strings.Join(targetAccounts, ", "))

			rm := service.NewRootManager(aws.NewIamClient(awscfg), aws.NewStsClient(awscfg), aws.NewOrganizationsClient(awscfg))
			resultMap, err := rm.RecoverRootPassword(ctx, targetAccounts)
			if err != nil {
				logger.Error("cmd.recovery", err, "failed to recover root password")
				return err
			}

			headers := []string{"Account", "Login Profile"}
			var data [][]any
			for acc, success := range resultMap {
				status := "recovered"
				if !success {
					status = "already exists"
				}
				data = append(data, []any{acc, status})
			}

			output.HandleOutput(outputFlag, headers, data)
			return nil
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of tarjet AWS account IDs (comma-separated). Use \"all\" to select all accounts.")
	return cmd
}
