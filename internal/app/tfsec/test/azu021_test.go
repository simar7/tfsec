package test

import (
	"testing"

	"github.com/tfsec/tfsec/internal/app/tfsec/rules"
)

func Test_AZUKeyVaultPurgeProtection(t *testing.T) {

	var tests = []struct {
		name                  string
		source                string
		mustIncludeResultCode string
		mustExcludeResultCode string
	}{
		{
			name: "check if purge_protection_enabled not set, check fails",
			source: `
resource "azurerm_key_vault" "bad_example" {
    name                        = "examplekeyvault"
    location                    = azurerm_resource_group.bad_example.location
    enabled_for_disk_encryption = true
    soft_delete_retention_days  = 7
    purge_protection_enabled    = false
}
`,
			mustIncludeResultCode: rules.AZUKeyVaultPurgeProtection,
		},
		{
			name: "check if purge_protection_enabled is set, check passes",
			source: `
resource "azurerm_key_vault" "good_example" {
    name                        = "examplekeyvault"
    location                    = azurerm_resource_group.good_example.location
    enabled_for_disk_encryption = true
    soft_delete_retention_days  = 7
    purge_protection_enabled    = true
}
`,
			mustExcludeResultCode: rules.AZUKeyVaultPurgeProtection,
		},
		{
			name: "check if purge_protection_enabled and soft_delete_retention_days is not set, check fails",
			source: `
resource "azurerm_key_vault" "bad_example" {
    name                        = "examplekeyvault"
    location                    = azurerm_resource_group.bad_example.location
    enabled_for_disk_encryption = true
    purge_protection_enabled    = false
}
`,
			mustIncludeResultCode: rules.AZUKeyVaultPurgeProtection,
		},
		{
			name: "check if purge_protection_enabled is set but soft_delete_retention_days is not set, check fails",
			source: `
resource "azurerm_key_vault" "bad_example" {
    name                        = "examplekeyvault"
    location                    = azurerm_resource_group.bad_example.location
    enabled_for_disk_encryption = true
    purge_protection_enabled    = true
}
`,
			mustIncludeResultCode: rules.AZUKeyVaultPurgeProtection,
		},
		{
			name: "check if purge_protection_enabled is set but soft_delete_retention_days is not set, check fails",
			source: `
resource "azurerm_key_vault" "bad_example" {
    name                        = "examplekeyvault"
    location                    = azurerm_resource_group.bad_example.location
    enabled_for_disk_encryption = true
	soft_delete_retention_days  = 0
    purge_protection_enabled    = true
}
`,
			mustIncludeResultCode: rules.AZUKeyVaultPurgeProtection,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results := scanSource(test.source)
			assertCheckCode(t, test.mustIncludeResultCode, test.mustExcludeResultCode, results)
		})
	}

}
