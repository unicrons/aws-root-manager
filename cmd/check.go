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

func Check() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check if centralized root access is enabled.",
		Long:  `Retrieve the status of centralized root access settings for an AWS Organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Trace("cmd.check", "check called")

			ctx := context.Background()
			awscfg, err := aws.LoadAWSConfig(ctx)
			if err != nil {
				logger.Error("cmd.check", err, "failed to load aws config")
				return err
			}

			rm := service.NewRootManager(aws.NewIamClient(awscfg), nil, nil)
			status, err := rm.CheckRootAccess(ctx)
			if err != nil {
				logger.Error("cmd.check", err, "failed to check root access configuration")
				return err
			}

			headers := []string{"Name", "Status"}
			data := [][]any{
				{"TrustedAccess", strconv.FormatBool(status.TrustedAccess)},
				{"RootCredentialsManagement", strconv.FormatBool(status.RootCredentialsManagement)},
				{"RootSessions", strconv.FormatBool(status.RootSessions)},
			}
			output.HandleOutput(outputFlag, headers, data)
			return nil
		},
	}
}
