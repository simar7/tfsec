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

const AWSS3DataShouldBeVersioned = "AWS077"
const AWSS3DataShouldBeVersionedDescription = "S3 Data should be versioned"
const AWSS3DataShouldBeVersionedImpact = "Deleted or modified data would not be recoverable"
const AWSS3DataShouldBeVersionedResolution = "Enable versioning to protect against accidental/malicious removal or modification"
const AWSS3DataShouldBeVersionedExplanation = `
Versioning in Amazon S3 is a means of keeping multiple variants of an object in the same bucket. 
You can use the S3 Versioning feature to preserve, retrieve, and restore every version of every object stored in your buckets. 
With versioning you can recover more easily from both unintended user actions and application failures.
`
const AWSS3DataShouldBeVersionedBadExample = `
resource "aws_s3_bucket" "bad_example" {

}
`
const AWSS3DataShouldBeVersionedGoodExample = `
resource "aws_s3_bucket" "good_example" {

	versioning {
		enabled = true
	}
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSS3DataShouldBeVersioned,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSS3DataShouldBeVersionedDescription,
			Impact:      AWSS3DataShouldBeVersionedImpact,
			Resolution:  AWSS3DataShouldBeVersionedResolution,
			Explanation: AWSS3DataShouldBeVersionedExplanation,
			BadExample:  AWSS3DataShouldBeVersionedBadExample,
			GoodExample: AWSS3DataShouldBeVersionedGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AmazonS3/latest/userguide/Versioning.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket#versioning",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_s3_bucket"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("versioning") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not have versioning enabled", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			versioningBlock := block.GetBlock("versioning")
			if versioningBlock.HasChild("enabled") && versioningBlock.GetAttribute("enabled").IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has versioning block but is disabled", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
