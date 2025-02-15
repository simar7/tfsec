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

	"github.com/zclconf/go-cty/cty"
)

// AzureUnencryptedDataLakeStore See https://github.com/tfsec/tfsec#included-checks for check info
const AzureUnencryptedDataLakeStore = "AZU004"
const AzureUnencryptedDataLakeStoreDescription = "Unencrypted data lake storage."
const AzureUnencryptedDataLakeStoreImpact = "Data could be read if compromised"
const AzureUnencryptedDataLakeStoreResolution = "Enable encryption of data lake storage"
const AzureUnencryptedDataLakeStoreExplanation = `
Datalake storage encryption defaults to Enabled, it shouldn't be overridden to Disabled.
`
const AzureUnencryptedDataLakeStoreBadExample = `
resource "azurerm_data_lake_store" "bad_example" {
	encryption_state = "Disabled"
}`
const AzureUnencryptedDataLakeStoreGoodExample = `
resource "azurerm_data_lake_store" "good_example" {
	encryption_state = "Enabled"
}`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AzureUnencryptedDataLakeStore,
		Documentation: rule.RuleDocumentation{
			Summary:     AzureUnencryptedDataLakeStoreDescription,
			Impact:      AzureUnencryptedDataLakeStoreImpact,
			Resolution:  AzureUnencryptedDataLakeStoreResolution,
			Explanation: AzureUnencryptedDataLakeStoreExplanation,
			BadExample:  AzureUnencryptedDataLakeStoreBadExample,
			GoodExample: AzureUnencryptedDataLakeStoreGoodExample,
			Links: []string{
				"https://docs.microsoft.com/en-us/azure/data-lake-store/data-lake-store-security-overview",
				"https://www.terraform.io/docs/providers/azurerm/r/data_lake_store.html",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_data_lake_store"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			encryptionStateAttr := block.GetAttribute("encryption_state")
			if encryptionStateAttr != nil && encryptionStateAttr.Type() == cty.String && encryptionStateAttr.Value().AsString() == "Disabled" {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf(
							"Resource '%s' defines an unencrypted data lake store.",
							block.FullName(),
						)).
						WithRange(encryptionStateAttr.Range()).
						WithAttributeAnnotation(encryptionStateAttr).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
