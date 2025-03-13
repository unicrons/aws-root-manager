package service

import (
	"context"
	"fmt"

	"github.com/unicrons/aws-root-manager/pkg/aws"
	"github.com/unicrons/aws-root-manager/pkg/logger"
	"github.com/unicrons/aws-root-manager/pkg/ui"
)

const (
	AllAccountsOption       = "all"
	AllAccountsSelectorText = "all non management accounts"
)

// Get target AWS accounts based on input flags or user interaction
func GetTargetAccounts(ctx context.Context, accounts []string) ([]string, error) {
	logger.Trace("service.GetTargetAccounts", "processing target accounts: %s", accounts)

	// if accounts are provided and "all" is not specified, return them
	if len(accounts) > 0 && accounts[0] != AllAccountsOption {
		return accounts, nil
	}

	// fetch all non-management accounts
	orgAccounts, err := aws.GetNonManagementOrganizationAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching organization accounts: %w", err)
	}

	// if "all" is specified, return all account IDs
	if len(accounts) > 0 && accounts[0] == AllAccountsOption {
		return convertAccountsToIDs(orgAccounts), nil
	}

	// prompt the user for account selection
	var selectorChoices []string
	selectorChoices = append(selectorChoices, AllAccountsSelectorText)
	for _, account := range orgAccounts {
		selectorChoices = append(selectorChoices, fmt.Sprintf("%s - %s", account.AccountID, account.Name))
	}
	selectedIndexes, err := ui.Prompt("Please select the AWS accounts to audit", selectorChoices)
	if err != nil {
		return nil, err
	}
	if len(selectedIndexes) == 0 {
		return nil, nil
	}

	// Resolve selected accounts
	if allSelected(selectedIndexes) {
		logger.Debug("service.GetTargetAccounts", "all accounts selected")
		return convertAccountsToIDs(orgAccounts), nil
	}

	return extractSelectedAccounts(orgAccounts, selectedIndexes), nil
}

// Checks if all option is selected
func allSelected(selectedIndexes []int) bool {
	for _, index := range selectedIndexes {
		if index == 0 {
			return true
		}
	}
	return false
}

// Get account IDs based on selector indexes
func extractSelectedAccounts(orgAccounts []aws.OrganizationAccount, selectedIndexes []int) []string {
	var selectedAccounts []string
	for _, index := range selectedIndexes {
		if index > 0 && index <= len(orgAccounts) {
			selectedAccounts = append(selectedAccounts, orgAccounts[index-1].AccountID) // -1 because of the all option
		}
	}
	return selectedAccounts
}

// Converts a slice of OrganizationAccount to a slice of account IDs
func convertAccountsToIDs(orgAccounts []aws.OrganizationAccount) []string {
	accountIDs := make([]string, len(orgAccounts))
	for i, account := range orgAccounts {
		accountIDs[i] = account.AccountID
	}
	return accountIDs
}
