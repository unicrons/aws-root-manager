package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/unicrons/aws-root-manager/internal/aws"
	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func DeleteS3BucketPolicy(newRM func(context.Context) (rootmanager.RootManager, error)) *cobra.Command {
	var accountId, bucketName string
	cmd := &cobra.Command{
		Use:          "s3-bucket-policy",
		Short:        "Delete an S3 bucket policy",
		Long:         `Delete the bucket policy attached to an S3 bucket owned by a member account using the S3UnlockBucketPolicy root task policy.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDeleteS3BucketPolicy(newRM, cmd.OutOrStdout(), accountId, bucketName)
		},
	}
	cmd.Flags().StringVar(&accountId, "account", "", "AWS account ID that owns the bucket (optional; if absent, a TUI lists the organization's accounts)")
	cmd.Flags().StringVar(&bucketName, "bucket", "", "Name of the S3 bucket (optional; if absent, a TUI lists the account's buckets)")
	return cmd
}

func runDeleteS3BucketPolicy(newRM func(context.Context) (rootmanager.RootManager, error), w io.Writer, accountId, bucketName string) error {
	ctx := context.Background()
	rm, err := newRM(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize root manager: %w", err)
	}

	accountId, err = selectSingleAccount(ctx, accountId)
	if err != nil {
		return err
	}

	if bucketName == "" {
		buckets, err := rm.ListAccountBuckets(ctx, accountId)
		if err != nil {
			return fmt.Errorf("failed to list buckets for account %s: %w", accountId, err)
		}
		if len(buckets) == 0 {
			return fmt.Errorf("no buckets found in account %s", accountId)
		}
		idx, err := ui.PromptSingle("Select the bucket whose policy will be deleted", buckets)
		if err != nil {
			return err
		}
		if idx < 0 {
			return fmt.Errorf("no bucket selected")
		}
		bucketName = buckets[idx]
	}

	policy, err := rm.GetS3BucketPolicy(ctx, accountId, bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket policy: %w", err)
	}
	if policy == "" {
		fmt.Fprintln(w, "No bucket policy found.")
		return nil
	}
	if outputFlag == "table" {
		fmt.Fprintf(w, "Current bucket policy for %s:\n\n", bucketName)
		output.RenderPolicy(w, policy)
	}

	if !skipFlag {
		confirmed, err := ui.Confirm("Delete this policy?")
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(w, "Aborted.")
			return nil
		}
	}

	result, err := rm.DeleteS3BucketPolicy(ctx, accountId, bucketName)
	if err != nil {
		return err
	}

	if !result.Success {
		slog.Error("failed to delete s3 bucket policy", "account_id", result.AccountId, "bucket", result.ResourceName, "error", result.Error)
		return fmt.Errorf("failed to delete bucket policy for bucket %s", result.ResourceName)
	}

	var headers []string
	var data [][]any
	if outputFlag == "table" {
		headers = []string{"Account", "ResourceType", "Bucket", "Status"}
		data = [][]any{{result.AccountId, result.ResourceType, result.ResourceName, "deleted"}}
	} else {
		headers = []string{"Account", "ResourceType", "Bucket", "Status", "Policy"}
		data = [][]any{{result.AccountId, result.ResourceType, result.ResourceName, "deleted", json.RawMessage(policy)}}
	}
	output.HandleOutput(w, outputFlag, headers, data)
	return nil
}

// selectSingleAccount resolves a single account ID from the --account flag or
// via a single-select TUI (no "all" option).
func selectSingleAccount(ctx context.Context, accountId string) (string, error) {
	var flag []string
	if accountId != "" {
		flag = []string{accountId}
	}
	awscfg, err := aws.LoadAWSConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load aws config: %w", err)
	}
	return ui.SelectSingleTargetAccount(ctx, aws.NewOrganizationsClient(awscfg), flag)
}
