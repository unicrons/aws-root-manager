package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/unicrons/aws-root-manager/internal/aws"
	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/spf13/cobra"
)

func Recovery(newRM func(context.Context) (rootmanager.RootManager, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery",
		Short: "Allow root password recovery",
		Long:  `Retrieve the status of centralized root access settings for an AWS Organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("recovery called")

			ctx := context.Background()
			rm, err := newRM(ctx)
			if err != nil {
				slog.Error("failed to initialize root manager", "error", err)
				return err
			}

			awscfg, err := aws.LoadAWSConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to load aws config: %w", err)
			}
			targetAccounts, err := ui.SelectTargetAccounts(ctx, aws.NewOrganizationsClient(awscfg), accountsFlags)
			if err != nil {
				slog.Error("failed to get target accounts", "error", err)
				return err
			}
			if len(targetAccounts) == 0 {
				slog.Info("no accounts selected")
				return nil
			}
			slog.Debug("selected accounts", "accounts", strings.Join(targetAccounts, ", "))

			if !skipFlag {
				confirmed, err := ui.Confirm(fmt.Sprintf("Restore root password for %d account(s)?", len(targetAccounts)))
				if err != nil {
					return err
				}
				if !confirmed {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}

			results, err := rm.RecoverRootPassword(ctx, targetAccounts)
			if err != nil {
				slog.Error("failed to recover root password", "error", err)
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

			output.HandleOutput(cmd.OutOrStdout(), outputFlag, headers, data)

			if failureCount > 0 {
				return fmt.Errorf("recovery failed for %d account(s)", failureCount)
			}

			return nil
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of tarjet AWS account IDs (comma-separated). Use \"all\" to select all accounts.")
	cmd.Flags().BoolVar(&skipFlag, "yes", false, "Skip the confirmation prompt")
	return cmd
}
