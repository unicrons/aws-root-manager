package cmd

import (
	"context"
	"strconv"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/internal/service"

	"github.com/spf13/cobra"
)

func Enable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable centralized root access",
		Long:  `Enable centralized root access management in an AWS Organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Trace("cmd.enable", "enable called")

			enableRootSessions, _ := cmd.Flags().GetBool("enableRootSessions")

			ctx := context.Background()
			awscfg, err := aws.LoadAWSConfig(ctx)
			if err != nil {
				logger.Error("cmd.enable", err, "failed to load aws config")
				return err
			}

			rm := service.NewRootManager(aws.NewIamClient(awscfg), nil, aws.NewOrganizationsClient(awscfg))
			initStatus, status, err := rm.EnableRootAccess(ctx, enableRootSessions)
			if err != nil {
				logger.Error("cmd.enable", err, "failed to enable root access")
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
