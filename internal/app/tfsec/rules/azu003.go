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

// AzureUnencryptedManagedDisk See https://github.com/tfsec/tfsec#included-checks for check info
const AzureUnencryptedManagedDisk = "AZU003"
const AzureUnencryptedManagedDiskDescription = "Unencrypted managed disk."
const AzureUnencryptedManagedDiskImpact = "Data could be read if compromised"
const AzureUnencryptedManagedDiskResolution = "Enable encryption on managed disks"
const AzureUnencryptedManagedDiskExplanation = `
Manage disks should be encrypted at rest. When specifying the <code>encryption_settings</code> block, the enabled attribute should be set to <code>true</code>.
`
const AzureUnencryptedManagedDiskBadExample = `
resource "azurerm_managed_disk" "bad_example" {
	encryption_settings {
		enabled = false
	}
}`
const AzureUnencryptedManagedDiskGoodExample = `
resource "azurerm_managed_disk" "good_example" {
	encryption_settings {
		enabled = true
	}
}`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AzureUnencryptedManagedDisk,
		Documentation: rule.RuleDocumentation{
			Summary:     AzureUnencryptedManagedDiskDescription,
			Impact:      AzureUnencryptedManagedDiskImpact,
			Resolution:  AzureUnencryptedManagedDiskResolution,
			Explanation: AzureUnencryptedManagedDiskExplanation,
			BadExample:  AzureUnencryptedManagedDiskBadExample,
			GoodExample: AzureUnencryptedManagedDiskGoodExample,
			Links: []string{
				"https://docs.microsoft.com/en-us/azure/virtual-machines/linux/disk-encryption",
				"https://www.terraform.io/docs/providers/azurerm/r/managed_disk.html",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_managed_disk"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			encryptionSettingsBlock := block.GetBlock("encryption_settings")
			if encryptionSettingsBlock == nil {
				return // encryption is by default now, so this is fine
			}

			enabledAttr := encryptionSettingsBlock.GetAttribute("enabled")
			if enabledAttr != nil && enabledAttr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf(
							"Resource '%s' defines an unencrypted managed disk.",
							block.FullName(),
						)).
						WithRange(enabledAttr.Range()).
						WithAttributeAnnotation(enabledAttr).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
