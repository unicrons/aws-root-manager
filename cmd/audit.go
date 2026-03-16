package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/rootmanager"

	"github.com/spf13/cobra"
)

func Audit(newRM func(context.Context) (rootmanager.RootManager, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "audit",
		Short:        "Retrieve root user credentials",
		Long:         `Retrieve available root user credentials for all member accounts within an AWS Organization.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("audit called")

			ctx := context.Background()
			rm, err := newRM(ctx)
			if err != nil {
				slog.Error("failed to initialize root manager", "error", err)
				return err
			}

			auditAccounts, err := ui.SelectTargetAccounts(ctx, accountsFlags)
			if err != nil {
				slog.Error("failed to get accounts to audit", "error", err)
				return err
			}
			if len(auditAccounts) == 0 {
				slog.Info("no accounts selected")
				return nil
			}
			slog.Debug("selected accounts", "accounts", strings.Join(auditAccounts, ", "))

			audit, err := rm.AuditAccounts(ctx, auditAccounts)
			if err != nil {
				slog.Error("failed to audit accounts", "error", err)
				return err
			}

			var skipped int
			headers := []string{"Account", "LoginProfile", "AccessKeys", "MFA Devices", "Signing Certificates"}
			var data [][]any
			for i, acc := range audit {
				if acc.Error != "" {
					skipped++
					continue
				}
				data = append(data, []any{
					auditAccounts[i],
					acc.LoginProfile,
					acc.AccessKeys,
					acc.MfaDevices,
					acc.SigningCertificates,
				})
			}
			output.HandleOutput(outputFlag, headers, data)

			if skipped > 0 {
				return fmt.Errorf("audit skipped for %d account(s)", skipped)
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringSliceVarP(&accountsFlags, "accounts", "a", []string{}, "List of AWS account IDs to audit (comma-separated). Use \"all\" to audit all accounts.")
	return cmd
}
