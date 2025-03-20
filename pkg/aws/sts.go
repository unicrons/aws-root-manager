package aws

import (
	"context"
	"fmt"

	"github.com/unicrons/aws-root-manager/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

const rootPolicyPrefix = "arn:aws:iam::aws:policy/root-task/"

type StsClient struct {
	client *sts.Client
}

func NewStsClient(awscfg aws.Config) *StsClient {
	client := sts.NewFromConfig(awscfg)
	return &StsClient{client: client}
}

func (c *StsClient) GetAssumeRootConfig(ctx context.Context, accountId, taskPolicyName string) (aws.Config, error) {
	logger.Trace("aws.GetAssumeRootConfig", "getting root aws.config account %s and task %s", accountId, taskPolicyName)

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

	logger.Debug("aws.GetAssumeRootConfig", "successfully generated assume root credentials for account %s and task %s", accountId, taskPolicyName)

	return awsrootcfg, nil
}

func (c *StsClient) assumeRoot(ctx context.Context, accountId, taskPolicyName string) (types.Credentials, error) {
	logger.Trace("aws.assumeRoot", "assuming root for account %s and task %s", accountId, taskPolicyName)

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
