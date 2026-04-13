package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type sqsClient struct {
	client *sqs.Client
}

func NewSqsClient(awscfg aws.Config) SqsClient {
	return &sqsClient{client: sqs.NewFromConfig(awscfg)}
}

func (c *sqsClient) ListQueues(ctx context.Context) ([]string, error) {
	slog.Debug("listing sqs queues")

	var queues []string
	var nextToken *string
	for {
		output, err := c.client.ListQueues(ctx, &sqs.ListQueuesInput{NextToken: nextToken})
		if err != nil {
			return nil, fmt.Errorf("error listing sqs queues: %w", err)
		}
		queues = append(queues, output.QueueUrls...)
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}
	return queues, nil
}

func (c *sqsClient) DeleteQueuePolicy(ctx context.Context, queueUrl string) error {
	slog.Debug("deleting sqs queue policy", "queue_url", queueUrl)

	_, err := c.client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(queueUrl),
		Attributes: map[string]string{
			string(types.QueueAttributeNamePolicy): "",
		},
	})
	if err != nil {
		return fmt.Errorf("error deleting queue policy for queue %s: %w", queueUrl, err)
	}
	return nil
}
