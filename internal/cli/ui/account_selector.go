package ui

import (
	"context"
	"fmt"

	"github.com/unicrons/aws-root-manager/internal/infra/aws"
	"github.com/unicrons/aws-root-manager/internal/logger"
)

const (
	AllAccountsOption       = "all"
	AllAccountsSelectorText = "all non management accounts"
)

// SelectTargetAccounts handles interactive account selection or returns accounts based on flags.
// Returns account IDs based on flags or TUI prompt.
func SelectTargetAccounts(ctx context.Context, accountsFlag []string) ([]string, error) {
	logger.Trace("ui.SelectTargetAccounts", "processing target accounts: %s", accountsFlag)

	// if accounts are provided and "all" is not specified, return them
	if len(accountsFlag) > 0 && accountsFlag[0] != AllAccountsOption {
		return accountsFlag, nil
	}

	// fetch all non-management accounts
	orgAccounts, err := aws.GetNonManagementOrganizationAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching organization accounts: %w", err)
	}

	// if "all" is specified, return all account IDs
	if len(accountsFlag) > 0 && accountsFlag[0] == AllAccountsOption {
		return convertAccountsToIDs(orgAccounts), nil
	}

	// prompt the user for account selection
	var selectorChoices []string
	selectorChoices = append(selectorChoices, AllAccountsSelectorText)
	for _, account := range orgAccounts {
		selectorChoices = append(selectorChoices, fmt.Sprintf("%s - %s", account.AccountID, account.Name))
	}
	selectedIndexes, err := Prompt("Please select the AWS accounts to audit", selectorChoices)
	if err != nil {
		return nil, err
	}
	if len(selectedIndexes) == 0 {
		return nil, nil
	}

	// Resolve selected accounts
	if allSelected(selectedIndexes) {
		logger.Debug("ui.SelectTargetAccounts", "all accounts selected")
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
