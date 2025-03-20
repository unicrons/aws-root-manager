package cmd

import (
	"context"
	"strconv"

	"github.com/unicrons/aws-root-manager/pkg/aws"
	"github.com/unicrons/aws-root-manager/pkg/logger"
	"github.com/unicrons/aws-root-manager/pkg/output"
	"github.com/unicrons/aws-root-manager/pkg/service"

	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable centralized root access",
	Long:  `Enable centralized root access management in an AWS Organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("cmd.enable", "enable called")

		enableRootSessions, _ := cmd.Flags().GetBool("enableRootSessions")

		ctx := context.Background()
		awscfg, err := aws.LoadAWSConfig(ctx)
		if err != nil {
			logger.Error("cmd.enable", err, "failed to load aws config")
			return
		}

		iam := aws.NewIamClient(awscfg)
		org := aws.NewOrganizationsClient(awscfg)

		initStatus, status, err := service.EnableRootAccess(ctx, iam, org, enableRootSessions)
		if err != nil {
			logger.Error("cmd.enable", err, "failed to enable root access")
			return
		}

		headers := []string{"Name", "InitialStatus", "CurrentStatus"}
		data := [][]any{
			{"TrustedAccess", strconv.FormatBool(initStatus.TrustedAccess), strconv.FormatBool(status.TrustedAccess)},
			{"RootCredentialsManagement", strconv.FormatBool(initStatus.RootCredentialsManagement), strconv.FormatBool(status.RootCredentialsManagement)},
			{"RootSessions", strconv.FormatBool(initStatus.RootSessions), strconv.FormatBool(status.RootSessions)},
		}
		output.HandleOutput(outputFlag, headers, data)
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.PersistentFlags().Bool("enableRootSessions", false, "Enable Root Sessions, required only when working with resource policies.")
}
