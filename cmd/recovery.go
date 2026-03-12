package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
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
			rm, err := service.NewRootManagerFromConfig(ctx)
			if err != nil {
				logger.Error("cmd.recovery", err, "failed to initialize root manager")
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

			results, err := rm.RecoverRootPassword(ctx, targetAccounts)
			if err != nil {
				logger.Error("cmd.recovery", err, "failed to recover root password")
				return err
			}

			headers := []string{"Account", "Login Profile", "Error"}
			var data [][]any
			var failureCount int
			for _, result := range results {
				status := "recovered"
				errorMsg := ""
				if !result.Success {
					if result.Error != "" {
						status = "failed"
						errorMsg = result.Error
						failureCount++
					} else {
						status = "already exists"
					}
				}
				data = append(data, []any{result.AccountId, status, errorMsg})
			}

			output.HandleOutput(outputFlag, headers, data)

			if failureCount > 0 {
				return fmt.Errorf("recovery failed for %d account(s)", failureCount)
			}

			return nil
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of tarjet AWS account IDs (comma-separated). Use \"all\" to select all accounts.")
	return cmd
}
