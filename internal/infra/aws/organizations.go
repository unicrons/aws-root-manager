package aws

import (
	"context"
	"fmt"

	"github.com/unicrons/aws-root-manager/internal/logger"

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

// Fetches AWS Organization accounts, excluding the management account
func GetNonManagementOrganizationAccounts(ctx context.Context) ([]OrganizationAccount, error) {
	logger.Trace("aws.GetNonManagementOrganizationAccounts", "getting organization accounts")

	awscfg, err := LoadAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	orgs := NewOrganizationsClient(awscfg)

	mgmAccount, err := orgs.(*organizationsClient).describeOrganization(ctx)
	if err != nil {
		return nil, err
	}

	orgAccounts, err := orgs.(*organizationsClient).listOrganizationAccounts()
	if err != nil {
		return nil, err
	}

	var nonManagementOrgAccounts []OrganizationAccount
	for _, acc := range orgAccounts {
		if string(acc.State) == "ACTIVE" && *acc.Id != mgmAccount {
			account := OrganizationAccount{
				Name:      *acc.Name,
				AccountID: *acc.Id,
			}
			nonManagementOrgAccounts = append(nonManagementOrgAccounts, account)
		}
	}

	return nonManagementOrgAccounts, nil
}

func (c *organizationsClient) listOrganizationAccounts() ([]types.Account, error) {
	logger.Trace("aws.listOrganizationAccounts", "listing organization accounts")

	params := &organizations.ListAccountsInput{}
	paginator := organizations.NewListAccountsPaginator(c.client, params)

	var allAccounts []types.Account

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to list organization accounts: %v", err)
		}
		allAccounts = append(allAccounts, page.Accounts...)
	}

	return allAccounts, nil
}

func (c *organizationsClient) describeOrganization(ctx context.Context) (string, error) {
	logger.Trace("aws.describeOrganization", "describing organization")

	organization, err := c.client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return "", fmt.Errorf("failed to describe organization: %w", err)
	}

	return *organization.Organization.MasterAccountId, nil
}

func (c *organizationsClient) EnableAWSServiceAccess(ctx context.Context, service string) error {
	logger.Trace("aws.EnableAWSServiceAccess", "enabling %s service access", service)

	_, err := c.client.EnableAWSServiceAccess(ctx, &organizations.EnableAWSServiceAccessInput{
		ServicePrincipal: aws.String(service),
	})
	if err != nil {
		return fmt.Errorf("aws.enableAWSServiceAccess: failed to enable service access: %w", err)
	}

	return nil
}
