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

const AZUTrustedMicrosoftServicesHaveStroageAccountAccess = "AZU013"
const AZUTrustedMicrosoftServicesHaveStroageAccountAccessDescription = "Trusted Microsoft Services should have bypass access to Storage accounts"
const AZUTrustedMicrosoftServicesHaveStroageAccountAccessImpact = "Trusted Microsoft Services won't be able to access storage account unless rules set to allow"
const AZUTrustedMicrosoftServicesHaveStroageAccountAccessResolution = "Allow Trusted Microsoft Services to bypass"
const AZUTrustedMicrosoftServicesHaveStroageAccountAccessExplanation = `
Some Microsoft services that interact with storage accounts operate from networks that can't be granted access through network rules. 

To help this type of service work as intended, allow the set of trusted Microsoft services to bypass the network rules
`
const AZUTrustedMicrosoftServicesHaveStroageAccountAccessBadExample = `
resource "azurerm_storage_account" "bad_example" {
  name                = "storageaccountname"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["100.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.example.id]
	bypass                     = ["Metrics"]
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Allow"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
  bypass                     = ["Metrics"]
}
`
const AZUTrustedMicrosoftServicesHaveStroageAccountAccessGoodExample = `
resource "azurerm_storage_account" "good_example" {
  name                = "storageaccountname"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["100.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.example.id]
    bypass                     = ["Metrics", "AzureServices"]
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Allow"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
  bypass                     = ["Metrics", "AzureServices"]
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AZUTrustedMicrosoftServicesHaveStroageAccountAccess,
		Documentation: rule.RuleDocumentation{
			Summary:     AZUTrustedMicrosoftServicesHaveStroageAccountAccessDescription,
			Impact:      AZUTrustedMicrosoftServicesHaveStroageAccountAccessImpact,
			Resolution:  AZUTrustedMicrosoftServicesHaveStroageAccountAccessResolution,
			Explanation: AZUTrustedMicrosoftServicesHaveStroageAccountAccessExplanation,
			BadExample:  AZUTrustedMicrosoftServicesHaveStroageAccountAccessBadExample,
			GoodExample: AZUTrustedMicrosoftServicesHaveStroageAccountAccessGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account#bypass",
				"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account_network_rules#bypass",
				"https://docs.microsoft.com/en-us/azure/storage/common/storage-network-security#trusted-microsoft-services",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_storage_account_network_rules", "azurerm_storage_account"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.IsResourceType("azurerm_storage_account") {
				if block.MissingChild("network_rules") {
				}
				block = block.GetBlock("network_rules")
			}

			if block.HasChild("bypass") {
				bypass := block.GetAttribute("bypass")
				if !bypass.Contains("AzureServices") {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' defines a network rule that doesn't allow bypass of Microsoft Services.", block.FullName())).
							WithRange(block.Range()).
							WithSeverity(severity.Error),
					)
				}
			}

		},
	})
}
