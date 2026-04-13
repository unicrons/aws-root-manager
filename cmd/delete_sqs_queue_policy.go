package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/unicrons/aws-root-manager/internal/cli/output"
	"github.com/unicrons/aws-root-manager/internal/cli/ui"
	"github.com/unicrons/aws-root-manager/rootmanager"
)

func DeleteSQSQueuePolicy(newRM func(context.Context) (rootmanager.RootManager, error)) *cobra.Command {
	var accountId, queueUrl string
	cmd := &cobra.Command{
		Use:          "sqs-queue-policy",
		Short:        "Delete an SQS queue policy",
		Long:         `Clear the access policy attached to an SQS queue owned by a member account using the SQSUnlockQueuePolicy root task policy.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDeleteSQSQueuePolicy(newRM, cmd.OutOrStdout(), accountId, queueUrl)
		},
	}
	cmd.Flags().StringVar(&accountId, "account", "", "AWS account ID that owns the queue (optional; if absent, a TUI lists the organization's accounts)")
	cmd.Flags().StringVar(&queueUrl, "queue", "", "URL of the SQS queue (optional; if absent, a TUI lists the account's queues)")
	return cmd
}

func runDeleteSQSQueuePolicy(newRM func(context.Context) (rootmanager.RootManager, error), w io.Writer, accountId, queueUrl string) error {
	ctx := context.Background()
	rm, err := newRM(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize root manager: %w", err)
	}

	accountId, err = selectSingleAccount(ctx, accountId)
	if err != nil {
		return err
	}

	if queueUrl == "" {
		queues, err := rm.ListAccountQueues(ctx, accountId)
		if err != nil {
			return fmt.Errorf("failed to list queues for account %s: %w", accountId, err)
		}
		if len(queues) == 0 {
			return fmt.Errorf("no queues found in account %s", accountId)
		}
		idx, err := ui.PromptSingle("Select the queue whose policy will be deleted", queues)
		if err != nil {
			return err
		}
		if idx < 0 {
			return fmt.Errorf("no queue selected")
		}
		queueUrl = queues[idx]
	}

	policy, err := rm.GetSQSQueuePolicy(ctx, accountId, queueUrl)
	if err != nil {
		return fmt.Errorf("failed to get queue policy: %w", err)
	}
	if policy == "" {
		fmt.Fprintln(w, "No queue policy found.")
		return nil
	}
	fmt.Fprintf(w, "Current queue policy for %s:\n\n%s\n\n", queueUrl, policy)

	confirmed, err := ui.Confirm("Delete this policy?")
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Fprintln(w, "Aborted.")
		return nil
	}

	result, err := rm.DeleteSQSQueuePolicy(ctx, accountId, queueUrl)
	if err != nil {
		return err
	}

	if !result.Success {
		slog.Error("failed to delete sqs queue policy", "account_id", result.AccountId, "queue_url", result.ResourceName, "error", result.Error)
		return fmt.Errorf("failed to delete queue policy for queue %s", result.ResourceName)
	}

	headers := []string{"Account", "ResourceType", "Queue", "Status"}
	data := [][]any{{result.AccountId, result.ResourceType, result.ResourceName, "deleted"}}
	output.HandleOutput(w, outputFlag, headers, data)
	return nil
}
