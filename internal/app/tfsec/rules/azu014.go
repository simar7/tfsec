package rules

import (
	"fmt"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/pkg/provider"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

const AZURequireSecureTransferForStorageAccounts = "AZU014"
const AZURequireSecureTransferForStorageAccountsDescription = "Storage accounts should be configured to only accept transfers that are over secure connections"
const AZURequireSecureTransferForStorageAccountsImpact = "Insecure transfer of data into secure accounts could be read if intercepted"
const AZURequireSecureTransferForStorageAccountsResolution = "Only allow secure connection for transfering data into storage accounts"
const AZURequireSecureTransferForStorageAccountsExplanation = `
You can configure your storage account to accept requests from secure connections only by setting the Secure transfer required property for the storage account. 

When you require secure transfer, any requests originating from an insecure connection are rejected. 

Microsoft recommends that you always require secure transfer for all of your storage accounts.
`
const AZURequireSecureTransferForStorageAccountsBadExample = `
resource "azurerm_storage_account" "bad_example" {
  name                      = "storageaccountname"
  resource_group_name       = azurerm_resource_group.example.name
  location                  = azurerm_resource_group.example.location
  account_tier              = "Standard"
  account_replication_type  = "GRS"
  enable_https_traffic_only = false
}
`
const AZURequireSecureTransferForStorageAccountsGoodExample = `
resource "azurerm_storage_account" "good_example" {
  name                      = "storageaccountname"
  resource_group_name       = azurerm_resource_group.example.name
  location                  = azurerm_resource_group.example.location
  account_tier              = "Standard"
  account_replication_type  = "GRS"
  enable_https_traffic_only = true
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AZURequireSecureTransferForStorageAccounts,
		Documentation: rule.RuleDocumentation{
			Summary:     AZURequireSecureTransferForStorageAccountsDescription,
			Impact:      AZURequireSecureTransferForStorageAccountsImpact,
			Resolution:  AZURequireSecureTransferForStorageAccountsResolution,
			Explanation: AZURequireSecureTransferForStorageAccountsExplanation,
			BadExample:  AZURequireSecureTransferForStorageAccountsBadExample,
			GoodExample: AZURequireSecureTransferForStorageAccountsGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account#enable_https_traffic_only",
				"https://docs.microsoft.com/en-us/azure/storage/common/storage-require-secure-transfer",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_storage_account"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.HasChild("enable_https_traffic_only") && block.GetAttribute("enable_https_traffic_only").IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' explicitly turns off secure transfer to storage account.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
