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

const AZUQueueStorageAnalyticsTurnedOn = "AZU016"
const AZUQueueStorageAnalyticsTurnedOnDescription = "When using Queue Services for a storage account, logging should be enabled."
const AZUQueueStorageAnalyticsTurnedOnImpact = "Logging provides valuable information about access and usage"
const AZUQueueStorageAnalyticsTurnedOnResolution = "Enable logging for Queue Services"
const AZUQueueStorageAnalyticsTurnedOnExplanation = `
Storage Analytics logs detailed information about successful and failed requests to a storage service. 

This information can be used to monitor individual requests and to diagnose issues with a storage service. 

Requests are logged on a best-effort basis.
`
const AZUQueueStorageAnalyticsTurnedOnBadExample = `
resource "azurerm_storage_account" "bad_example" {
    name                     = "example"
    resource_group_name      = data.azurerm_resource_group.example.name
    location                 = data.azurerm_resource_group.example.location
    account_tier             = "Standard"
    account_replication_type = "GRS"
    queue_properties  {
  }
}
`
const AZUQueueStorageAnalyticsTurnedOnGoodExample = `
resource "azurerm_storage_account" "good_example" {
    name                     = "example"
    resource_group_name      = data.azurerm_resource_group.example.name
    location                 = data.azurerm_resource_group.example.location
    account_tier             = "Standard"
    account_replication_type = "GRS"
    queue_properties  {
    logging {
        delete                = true
        read                  = true
        write                 = true
        version               = "1.0"
        retention_policy_days = 10
    }
  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AZUQueueStorageAnalyticsTurnedOn,
		Documentation: rule.RuleDocumentation{
			Summary:     AZUQueueStorageAnalyticsTurnedOnDescription,
			Impact:      AZUQueueStorageAnalyticsTurnedOnImpact,
			Resolution:  AZUQueueStorageAnalyticsTurnedOnResolution,
			Explanation: AZUQueueStorageAnalyticsTurnedOnExplanation,
			BadExample:  AZUQueueStorageAnalyticsTurnedOnBadExample,
			GoodExample: AZUQueueStorageAnalyticsTurnedOnGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account#logging",
				"https://docs.microsoft.com/en-us/azure/storage/common/storage-analytics-logging?tabs=dotnet",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_storage_account"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.HasChild("queue_properties") {
				queueProps := block.GetBlock("queue_properties")
				if queueProps.MissingChild("logging") {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' defines a Queue Services storage account without Storage Analytics logging.", block.FullName())).
							WithRange(block.Range()).
							WithSeverity(severity.Warning),
					)
				}
			}

		},
	})
}
