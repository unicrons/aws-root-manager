package cmd

import (
	"context"
	"strconv"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/service"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if centralized root access is enabled.",
	Long:  `Retrieve the status of centralized root access settings for an AWS Organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("cmd.check", "check called")

		ctx := context.Background()
		awscfg, err := aws.LoadAWSConfig(ctx)
		if err != nil {
			logger.Error("cmd.check", err, "failed to load aws config")
			return
		}

		iam := aws.NewIamClient(awscfg)
		test, err := service.CheckRootAccess(ctx, iam)
		if err != nil {
			logger.Error("cmd.check", err, "failed to check root access configuration")
			return
		}

		headers := []string{"Name", "Status"}
		data := [][]any{
			{"TrustedAccess", strconv.FormatBool(test.TrustedAccess)},
			{"RootCredentialsManagement", strconv.FormatBool(test.RootCredentialsManagement)},
			{"RootSessions", strconv.FormatBool(test.RootSessions)},
		}
		output.HandleOutput(outputFlag, headers, data)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
