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

const AZUAKSClusterRBACenabled = "AZU007"
const AZUAKSClusterRBACenabledDescription = "Ensure RBAC is enabled on AKS clusters"
const AZUAKSClusterRBACenabledImpact = "No role based access control is in place for the AKS cluster"
const AZUAKSClusterRBACenabledResolution = "Enable RBAC"
const AZUAKSClusterRBACenabledExplanation = `
Using Kubernetes role-based access control (RBAC), you can grant users, groups, and service accounts access to only the resources they need.
`
const AZUAKSClusterRBACenabledBadExample = `
resource "azurerm_kubernetes_cluster" "bad_example" {
	role_based_access_control {
		enabled = false
	}
}
`
const AZUAKSClusterRBACenabledGoodExample = `
resource "azurerm_kubernetes_cluster" "good_example" {
	role_based_access_control {
		enabled = true
	}
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AZUAKSClusterRBACenabled,
		Documentation: rule.RuleDocumentation{
			Summary:     AZUAKSClusterRBACenabledDescription,
			Impact:      AZUAKSClusterRBACenabledImpact,
			Resolution:  AZUAKSClusterRBACenabledResolution,
			Explanation: AZUAKSClusterRBACenabledExplanation,
			BadExample:  AZUAKSClusterRBACenabledBadExample,
			GoodExample: AZUAKSClusterRBACenabledGoodExample,
			Links: []string{
				"https://www.terraform.io/docs/providers/azurerm/r/kubernetes_cluster.html#role_based_access_control",
				"https://docs.microsoft.com/en-us/azure/aks/concepts-identity",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_kubernetes_cluster", "role_based_access_control"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			rbacBlock := block.GetBlock("role_based_access_control")
			if rbacBlock == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines without RBAC", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			enabledAttr := rbacBlock.GetAttribute("enabled")
			if enabledAttr.Type() == cty.Bool && enabledAttr.Value().False() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf(
							"Resource '%s' RBAC disabled.",
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
