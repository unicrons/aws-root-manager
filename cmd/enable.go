package cmd

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/service"

	"github.com/spf13/cobra"
)

func Enable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable centralized root access",
		Long:  `Enable centralized root access management in an AWS Organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("enable called")

			enableRootSessions, _ := cmd.Flags().GetBool("enableRootSessions")

			ctx := context.Background()
			rm, err := service.NewRootManagerFromConfig(ctx)
			if err != nil {
				slog.Error("failed to initialize root manager", "error", err)
				return err
			}

			initStatus, status, err := rm.EnableRootAccess(ctx, enableRootSessions)
			if err != nil {
				slog.Error("failed to enable root access", "error", err)
				return err
			}

			headers := []string{"Name", "InitialStatus", "CurrentStatus"}
			data := [][]any{
				{"TrustedAccess", strconv.FormatBool(initStatus.TrustedAccess), strconv.FormatBool(status.TrustedAccess)},
				{"RootCredentialsManagement", strconv.FormatBool(initStatus.RootCredentialsManagement), strconv.FormatBool(status.RootCredentialsManagement)},
				{"RootSessions", strconv.FormatBool(initStatus.RootSessions), strconv.FormatBool(status.RootSessions)},
			}
			output.HandleOutput(outputFlag, headers, data)
			return nil
		},
	}
	cmd.PersistentFlags().Bool("enableRootSessions", false, "Enable Root Sessions, required only when working with resource policies.")
	return cmd
}
