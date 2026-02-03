package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const defaultRegion = "us-east-1"

type configOption func(*config.LoadOptions) error

func WithCredentials(creds *aws.Credentials) configOption {
	return func(o *config.LoadOptions) error {
		o.Credentials = credentials.NewStaticCredentialsProvider(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
		)
		return nil
	}
}

func WithDefaultRegion(region string) configOption {
	return func(o *config.LoadOptions) error {
		o.DefaultRegion = region
		return nil
	}
}

func LoadAWSConfig(ctx context.Context, opts ...configOption) (aws.Config, error) {
	options := []func(*config.LoadOptions) error{
		WithDefaultRegion(defaultRegion),
	}

	for _, opt := range opts {
		options = append(options, opt)
	}

	awscfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return aws.Config{}, err
	}

	return awscfg, nil
}
