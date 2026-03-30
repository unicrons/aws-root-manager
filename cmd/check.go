package cmd

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/spf13/cobra"
)

func Check(newRM func(context.Context) (rootmanager.RootManager, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check if centralized root access is enabled.",
		Long:  `Retrieve the status of centralized root access settings for an AWS Organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("check called")

			ctx := context.Background()
			rm, err := newRM(ctx)
			if err != nil {
				slog.Error("failed to initialize root manager", "error", err)
				return err
			}

			status, err := rm.CheckRootAccess(ctx)
			if err != nil {
				slog.Error("failed to check root access configuration", "error", err)
				return err
			}

			headers := []string{"Name", "Status"}
			data := [][]any{
				{"TrustedAccess", strconv.FormatBool(status.TrustedAccess)},
				{"RootCredentialsManagement", strconv.FormatBool(status.RootCredentialsManagement)},
				{"RootSessions", strconv.FormatBool(status.RootSessions)},
			}
			output.HandleOutput(cmd.OutOrStdout(), outputFlag, headers, data)
			return nil
		},
	}
}
