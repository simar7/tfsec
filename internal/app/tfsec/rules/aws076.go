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

const AWSBlockPublicPolicyS3 = "AWS076"
const AWSBlockPublicPolicyS3Description = "S3 Access block should block public policy"
const AWSBlockPublicPolicyS3Impact = "Users could put a policy that allows public access"
const AWSBlockPublicPolicyS3Resolution = "Prevent policies that allow public access being PUT"
const AWSBlockPublicPolicyS3Explanation = `
S3 bucket policy should have block public policy to prevent users from PUTing a policy that enable public access.
`
const AWSBlockPublicPolicyS3BadExample = `
resource "aws_s3_bucket_public_access_block" "bad_example" {
	bucket = aws_s3_bucket.example.id
}

resource "aws_s3_bucket_public_access_block" "bad_example" {
	bucket = aws_s3_bucket.example.id
  
	block_public_policy = false
}
`
const AWSBlockPublicPolicyS3GoodExample = `
resource "aws_s3_bucket_public_access_block" "good_example" {
	bucket = aws_s3_bucket.example.id
  
	block_public_policy = true
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSBlockPublicPolicyS3,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSBlockPublicPolicyS3Description,
			Impact:      AWSBlockPublicPolicyS3Impact,
			Resolution:  AWSBlockPublicPolicyS3Resolution,
			Explanation: AWSBlockPublicPolicyS3Explanation,
			BadExample:  AWSBlockPublicPolicyS3BadExample,
			GoodExample: AWSBlockPublicPolicyS3GoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AmazonS3/latest/dev-retired/access-control-block-public-access.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_public_access_block#block_public_policy",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_s3_bucket_public_access_block"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if block.MissingChild("block_public_policy") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not specify block_public_policy, defaults to false", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			attr := block.GetAttribute("block_public_policy")
			if attr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' sets block_public_policy explicitly to false", block.FullName())).
						WithRange(attr.Range()).
						WithAttributeAnnotation(attr).
						WithSeverity(severity.Error),
				)
			}
		},
	})
}
