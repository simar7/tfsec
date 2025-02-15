package rules

import (
	"fmt"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/pkg/provider"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/zclconf/go-cty/cty"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

// AWSEcrImageScanNotEnabled See https://github.com/tfsec/tfsec#included-checks for check info
const AWSEcrImageScanNotEnabled = "AWS023"
const AWSEcrImageScanNotEnabledDescription = "ECR repository has image scans disabled."
const AWSEcrImageScanNotEnabledImpact = "The ability to scan images is not being used and vulnerabilities will not be highlighted"
const AWSEcrImageScanNotEnabledResolution = "Enable ECR image scanning"
const AWSEcrImageScanNotEnabledExplanation = `
Repository image scans should be enabled to ensure vulnerable software can be discovered and remediated as soon as possible.
`
const AWSEcrImageScanNotEnabledBadExample = `
resource "aws_ecr_repository" "bad_example" {
  name                 = "bar"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = false
  }
}
`
const AWSEcrImageScanNotEnabledGoodExample = `
resource "aws_ecr_repository" "good_example" {
  name                 = "bar"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSEcrImageScanNotEnabled,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSEcrImageScanNotEnabledDescription,
			Impact:      AWSEcrImageScanNotEnabledImpact,
			Resolution:  AWSEcrImageScanNotEnabledResolution,
			Explanation: AWSEcrImageScanNotEnabledExplanation,
			BadExample:  AWSEcrImageScanNotEnabledBadExample,
			GoodExample: AWSEcrImageScanNotEnabledGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecr_repository#image_scanning_configuration",
				"https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-scanning.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_ecr_repository"},
		CheckFunc: func(set result.Set, block *block.Block, context *hclcontext.Context) {

			ecrScanStatusBlock := block.GetBlock("image_scanning_configuration")
			ecrScanStatusAttr := ecrScanStatusBlock.GetAttribute("scan_on_push")

			if ecrScanStatusAttr == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a disabled ECR image scan.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			} else if ecrScanStatusAttr.Type() == cty.Bool && ecrScanStatusAttr.Value().False() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a disabled ECR image scan.", block.FullName())).
						WithRange(ecrScanStatusAttr.Range()).
						WithAttributeAnnotation(ecrScanStatusAttr).
						WithSeverity(severity.Error),
				)
			}
		},
	})
}
