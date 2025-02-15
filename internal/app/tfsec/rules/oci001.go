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

const OCIComputeIpReservation = "OCI001"
const OCIComputeIpReservationDescription = "Compute instance requests an IP reservation from a public pool"
const OCIComputeIpReservationImpact = "The compute instance has the ability to be reached from outside"
const OCIComputeIpReservationResolution = "Reconsider the use of an public IP"
const OCIComputeIpReservationExplanation = `
Compute instance requests an IP reservation from a public pool

The compute instance has the ability to be reached from outside, you might want to sonder the use of a non public IP.
`
const OCIComputeIpReservationBadExample = `
resource "opc_compute_ip_address_reservation" "my-ip-address" {
	name            = "my-ip-address"
	ip_address_pool = "public-ippool"
  }
`
const OCIComputeIpReservationGoodExample = `
resource "opc_compute_ip_address_reservation" "my-ip-address" {
	name            = "my-ip-address"
	ip_address_pool = "cloud-ippool"
  }
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: OCIComputeIpReservation,
		Documentation: rule.RuleDocumentation{
			Summary:     OCIComputeIpReservationDescription,
			Explanation: OCIComputeIpReservationExplanation,
			Impact:      OCIComputeIpReservationImpact,
			Resolution:  OCIComputeIpReservationResolution,
			BadExample:  OCIComputeIpReservationBadExample,
			GoodExample: OCIComputeIpReservationGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/opc/latest/docs/resources/opc_compute_ip_address_reservation",
				"https://registry.terraform.io/providers/hashicorp/opc/latest/docs/resources/opc_compute_instance",
			},
		},
		Provider:       provider.OracleProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"opc_compute_ip_address_reservation"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if attr := block.GetAttribute("ip_address_pool"); attr != nil {
				if attr.IsAny("public-ippool") {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' is using an IP from a public IP pool", block.FullName())).
							WithRange(attr.Range()).
							WithAttributeAnnotation(attr).
							WithSeverity(severity.Warning),
					)
				}
			}
		},
	})
}
