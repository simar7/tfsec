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

const AZUAKSClusterNetworkPolicy = "AZU006"
const AZUAKSClusterNetworkPolicyDescription = "Ensure AKS cluster has Network Policy configured"
const AZUAKSClusterNetworkPolicyImpact = "No network policy is protecting the AKS cluster"
const AZUAKSClusterNetworkPolicyResolution = "Configure a network policy"
const AZUAKSClusterNetworkPolicyExplanation = `
The Kubernetes object type NetworkPolicy should be defined to have opportunity allow or block traffic to pods, as in a Kubernetes cluster configured with default settings, all pods can discover and communicate with each other without any restrictions.
`
const AZUAKSClusterNetworkPolicyBadExample = `
resource "azurerm_kubernetes_cluster" "bad_example" {
	network_profile {
	  }
}
`
const AZUAKSClusterNetworkPolicyGoodExample = `
resource "azurerm_kubernetes_cluster" "good_example" {
	network_profile {
	  network_policy = "calico"
	  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AZUAKSClusterNetworkPolicy,
		Documentation: rule.RuleDocumentation{
			Summary:     AZUAKSClusterNetworkPolicyDescription,
			Impact:      AZUAKSClusterNetworkPolicyImpact,
			Resolution:  AZUAKSClusterNetworkPolicyResolution,
			Explanation: AZUAKSClusterNetworkPolicyExplanation,
			BadExample:  AZUAKSClusterNetworkPolicyBadExample,
			GoodExample: AZUAKSClusterNetworkPolicyGoodExample,
			Links: []string{
				"https://www.terraform.io/docs/providers/azurerm/r/kubernetes_cluster.html#network_policy",
				"https://docs.microsoft.com/en-us/azure/aks/use-network-policies",
				"https://kubernetes.io/docs/concepts/services-networking/network-policies",
			},
		},
		Provider:       provider.AzureProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_kubernetes_cluster"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if networkProfileBlock := block.GetBlock("network_profile"); networkProfileBlock != nil {
				if networkProfileBlock.GetAttribute("network_policy") == nil {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' do not have network_policy define. network_policy should be defined to have opportunity allow or block traffic to pods", block.FullName())).
							WithRange(networkProfileBlock.Range()).
							WithSeverity(severity.Error),
					)
				}
			}

		},
	})
}
