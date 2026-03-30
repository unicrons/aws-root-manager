package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

type organizationsClient struct {
	client *organizations.Client
}

func NewOrganizationsClient(awscfg aws.Config) OrganizationsClient {
	client := organizations.NewFromConfig(awscfg)
	return &organizationsClient{client: client}
}

type OrganizationAccount struct {
	Name      string
	AccountID string
}

// GetNonManagementOrganizationAccounts fetches active organization accounts, excluding the management account.
func GetNonManagementOrganizationAccounts(ctx context.Context, org OrganizationsClient) ([]OrganizationAccount, error) {
	slog.Debug("getting organization accounts")

	mgmAccountId, err := org.DescribeOrganization(ctx)
	if err != nil {
		return nil, err
	}

	allAccounts, err := org.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}

	var nonManagementAccounts []OrganizationAccount
	for _, acc := range allAccounts {
		if acc.AccountID != mgmAccountId {
			nonManagementAccounts = append(nonManagementAccounts, acc)
		}
	}

	return nonManagementAccounts, nil
}

func (c *organizationsClient) DescribeOrganization(ctx context.Context) (string, error) {
	slog.Debug("describing organization")

	organization, err := c.client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return "", fmt.Errorf("failed to describe organization: %w", err)
	}

	return *organization.Organization.MasterAccountId, nil
}

func (c *organizationsClient) ListAccounts(ctx context.Context) ([]OrganizationAccount, error) {
	slog.Debug("listing organization accounts")

	params := &organizations.ListAccountsInput{}
	paginator := organizations.NewListAccountsPaginator(c.client, params)

	var accounts []OrganizationAccount

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list organization accounts: %v", err)
		}
		for _, acc := range page.Accounts {
			if acc.Status == types.AccountStatusActive {
				accounts = append(accounts, OrganizationAccount{
					Name:      aws.ToString(acc.Name),
					AccountID: aws.ToString(acc.Id),
				})
			}
		}
	}

	return accounts, nil
}

func (c *organizationsClient) EnableAWSServiceAccess(ctx context.Context, service string) error {
	slog.Debug("enabling service access", "service", service)

	_, err := c.client.EnableAWSServiceAccess(ctx, &organizations.EnableAWSServiceAccessInput{
		ServicePrincipal: aws.String(service),
	})
	if err != nil {
		return fmt.Errorf("aws.enableAWSServiceAccess: failed to enable service access: %w", err)
	}

	return nil
}
