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

const AWSDontUseDefaultAWSVPC = "AWS082"
const AWSDontUseDefaultAWSVPCDescription = "It is AWS best practice to not use the default VPC for workflows"
const AWSDontUseDefaultAWSVPCImpact = "The default VPC does not have critical security features applied"
const AWSDontUseDefaultAWSVPCResolution = "Create a non-default vpc for resources to be created in"
const AWSDontUseDefaultAWSVPCExplanation = `
Default VPC does not have a lot of the critical security features that standard VPC comes with, new resources should not be created in the default VPC and it should not be present in the Terraform.
`
const AWSDontUseDefaultAWSVPCBadExample = `
resource "aws_default_vpc" "default" {
	tags = {
	  Name = "Default VPC"
	}
  }
`
const AWSDontUseDefaultAWSVPCGoodExample = `
# no aws default vpc present
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSDontUseDefaultAWSVPC,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSDontUseDefaultAWSVPCDescription,
			Explanation: AWSDontUseDefaultAWSVPCExplanation,
			Impact:      AWSDontUseDefaultAWSVPCImpact,
			Resolution:  AWSDontUseDefaultAWSVPCResolution,
			BadExample:  AWSDontUseDefaultAWSVPCBadExample,
			GoodExample: AWSDontUseDefaultAWSVPCGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/default_vpc",
				"https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_default_vpc"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			set.Add(
				result.New().
					WithDescription(fmt.Sprintf("Resource '%s' should not exist", block.FullName())).
					WithRange(block.Range()).
					WithSeverity(severity.Error),
			)
		},
	})
}
