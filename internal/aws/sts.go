package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

const rootPolicyPrefix = "arn:aws:iam::aws:policy/root-task/"

type stsClient struct {
	client *sts.Client
}

func NewStsClient(awscfg aws.Config) StsClient {
	client := sts.NewFromConfig(awscfg)
	return &stsClient{client: client}
}

func (c *stsClient) GetAssumeRootConfig(ctx context.Context, accountId, taskPolicyName string) (aws.Config, error) {
	slog.Debug("getting root aws config", "account_id", accountId, "task", taskPolicyName)

	stsCreds, err := c.assumeRoot(ctx, accountId, taskPolicyName)
	if err != nil {
		return aws.Config{}, err
	}

	// Convert sts.Credentials to aws.Credentials
	awsCreds := aws.Credentials{
		AccessKeyID:     aws.ToString(stsCreds.AccessKeyId),
		SecretAccessKey: aws.ToString(stsCreds.SecretAccessKey),
		SessionToken:    aws.ToString(stsCreds.SessionToken),
	}

	awsrootcfg, err := LoadAWSConfig(ctx, WithCredentials(&awsCreds))

	if err != nil {
		return aws.Config{}, fmt.Errorf("error loading aws root config: %s", err)
	}

	slog.Debug("successfully generated assume root credentials", "account_id", accountId, "task", taskPolicyName)

	return awsrootcfg, nil
}

func (c *stsClient) assumeRoot(ctx context.Context, accountId, taskPolicyName string) (types.Credentials, error) {
	slog.Debug("assuming root", "account_id", accountId, "task", taskPolicyName)

	params := &sts.AssumeRootInput{
		TargetPrincipal: aws.String(accountId),
		TaskPolicyArn: &types.PolicyDescriptorType{
			Arn: aws.String(rootPolicyPrefix + taskPolicyName),
		},
		DurationSeconds: aws.Int32(60),
	}

	output, err := c.client.AssumeRoot(ctx, params)
	if err != nil {
		return types.Credentials{}, fmt.Errorf("assume root failed for account %s and task %s: %w", accountId, taskPolicyName, err)
	}

	return *output.Credentials, err
}
