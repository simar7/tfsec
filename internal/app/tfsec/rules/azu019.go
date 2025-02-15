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

const AZUDatabaseAuditingRetention90Days = "AZU019"
const AZUDatabaseAuditingRetention90DaysDescription = "Database auditing rentention period should be longer than 90 days"
const AZUDatabaseAuditingRetention90DaysImpact = "Short logging retention could result in missing valuable historical information"
const AZUDatabaseAuditingRetention90DaysResolution = "Set retention periods of database auditing to greater than 90 days"
const AZUDatabaseAuditingRetention90DaysExplanation = `
When Auditing is configured for a SQL database, if the retention period is not set, the retention will be unlimited.

If the retention period is to be explicitly set, it should be set for no less than 90 days.

`
const AZUDatabaseAuditingRetention90DaysBadExample = `
resource "azurerm_mssql_database_extended_auditing_policy" "bad_example" {
  database_id                             = azurerm_mssql_database.example.id
  storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.example.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
`
const AZUDatabaseAuditingRetention90DaysGoodExample = `
resource "azurerm_mssql_database_extended_auditing_policy" "good_example" {
  database_id                             = azurerm_mssql_database.example.id
  storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.example.primary_access_key
  storage_account_access_key_is_secondary = false
}

resource "azurerm_mssql_database_extended_auditing_policy" "good_example" {
  database_id                             = azurerm_mssql_database.example.id
  storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.example.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 90
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AZUDatabaseAuditingRetention90Days,
		Documentation: rule.RuleDocumentation{
			Summary:     AZUDatabaseAuditingRetention90DaysDescription,
			Impact:      AZUDatabaseAuditingRetention90DaysImpact,
			Resolution:  AZUDatabaseAuditingRetention90DaysResolution,
			Explanation: AZUDatabaseAuditingRetention90DaysExplanation,
			BadExample:  AZUDatabaseAuditingRetention90DaysBadExample,
			GoodExample: AZUDatabaseAuditingRetention90DaysGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_database_extended_auditing_policy",
				"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_server#retention_in_days",
				"https://docs.microsoft.com/en-us/azure/azure-sql/database/auditing-overview",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_sql_server", "azurerm_sql_server", "azurerm_mssql_database_extended_auditing_policy"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if !block.IsResourceType("azurerm_mssql_database_extended_auditing_policy") {
				if block.MissingChild("extended_auditing_policy") {
				}
				block = block.GetBlock("extended_auditing_policy")
			}

			if block.MissingChild("retention_in_days") {
				// using default of unlimited
			}
			if block.GetAttribute("retention_in_days").LessThan(90) {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' specifies a retention period of less than 90 days.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
