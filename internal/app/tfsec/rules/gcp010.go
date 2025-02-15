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

// GkeShieldedNodesDisabled See https://github.com/tfsec/tfsec#included-checks for check info
const GkeShieldedNodesDisabled = "GCP010"
const GkeShieldedNodesDisabledDescription = "Shielded GKE nodes not enabled."
const GkeShieldedNodesDisabledImpact = "Node identity and integrity can't be verified without shielded GKE nodes"
const GkeShieldedNodesDisabledResolution = "Enable node shielding"
const GkeShieldedNodesDisabledExplanation = `
CIS GKE Benchmark Recommendation: 6.5.5. Ensure Shielded GKE Nodes are Enabled

Shielded GKE Nodes provide strong, verifiable node identity and integrity to increase the security of GKE nodes and should be enabled on all GKE clusters.
`
const GkeShieldedNodesDisabledBadExample = `
resource "google_container_cluster" "bad_example" {
	enable_shielded_nodes = "false"
}`
const GkeShieldedNodesDisabledGoodExample = `
resource "google_container_cluster" "good_example" {
	enable_shielded_nodes = "true"
}`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: GkeShieldedNodesDisabled,
		Documentation: rule.RuleDocumentation{
			Summary:     GkeShieldedNodesDisabledDescription,
			Impact:      GkeShieldedNodesDisabledImpact,
			Resolution:  GkeShieldedNodesDisabledResolution,
			Explanation: GkeShieldedNodesDisabledExplanation,
			BadExample:  GkeShieldedNodesDisabledBadExample,
			GoodExample: GkeShieldedNodesDisabledGoodExample,
			Links: []string{
				"https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster#shielded_nodes",
				"https://www.terraform.io/docs/providers/google/r/container_cluster.html#enable_shielded_nodes",
			},
		},
		Provider:       provider.GCPProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"google_container_cluster"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("enable_shielded_nodes") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a cluster with shielded nodes disabled. Shielded GKE Nodes provide strong, verifiable node identity and integrity to increase the security of GKE nodes and should be enabled on all GKE clusters.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			enableShieldedNodesAttr := block.GetAttribute("enable_shielded_nodes")

			if enableShieldedNodesAttr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a cluster with shielded nodes disabled. Shielded GKE Nodes provide strong, verifiable node identity and integrity to increase the security of GKE nodes and should be enabled on all GKE clusters.", block.FullName())).
						WithRange(enableShieldedNodesAttr.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
