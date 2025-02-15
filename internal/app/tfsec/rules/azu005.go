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

// AzureVMWithPasswordAuthentication See https://github.com/tfsec/tfsec#included-checks for check info
const AzureVMWithPasswordAuthentication = "AZU005"
const AzureVMWithPasswordAuthenticationDescription = "Password authentication in use instead of SSH keys."
const AzureVMWithPasswordAuthenticationImpact = "Passwords are potentially easier to compromise than SSH Keys"
const AzureVMWithPasswordAuthenticationResolution = "Use SSH keys for authentication"
const AzureVMWithPasswordAuthenticationExplanation = `
Access to instances should be authenticated using SSH keys. Removing the option of password authentication enforces more secure methods while removing the risks inherent with passwords.
`
const AzureVMWithPasswordAuthenticationBadExample = `
resource "azurerm_virtual_machine" "bad_example" {
	os_profile_linux_config {
		disable_password_authentication = false
	}
}`
const AzureVMWithPasswordAuthenticationGoodExample = `
resource "azurerm_virtual_machine" "good_example" {
	os_profile_linux_config {
		disable_password_authentication = true
	}
}`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AzureVMWithPasswordAuthentication,
		Documentation: rule.RuleDocumentation{
			Summary:     AzureVMWithPasswordAuthenticationDescription,
			Impact:      AzureVMWithPasswordAuthenticationImpact,
			Resolution:  AzureVMWithPasswordAuthenticationResolution,
			Explanation: AzureVMWithPasswordAuthenticationExplanation,
			BadExample:  AzureVMWithPasswordAuthenticationBadExample,
			GoodExample: AzureVMWithPasswordAuthenticationGoodExample,
			Links: []string{
				"https://docs.microsoft.com/en-us/azure/virtual-machines/linux/create-ssh-keys-detailed",
				"https://www.terraform.io/docs/providers/azurerm/r/virtual_machine.html",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_virtual_machine"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if linuxConfigBlock := block.GetBlock("os_profile_linux_config"); linuxConfigBlock != nil {
				passwordAuthDisabledAttr := linuxConfigBlock.GetAttribute("disable_password_authentication")
				if passwordAuthDisabledAttr != nil && passwordAuthDisabledAttr.Type() == cty.Bool && passwordAuthDisabledAttr.Value().False() {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf(
								"Resource '%s' has password authentication enabled. Use SSH keys instead.",
								block.FullName(),
							)).
							WithRange(passwordAuthDisabledAttr.Range()).
							WithAttributeAnnotation(passwordAuthDisabledAttr).
							WithSeverity(severity.Error),
					)
				}
			}

		},
	})
}
