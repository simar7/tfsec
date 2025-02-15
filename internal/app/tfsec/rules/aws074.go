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

const AWSBlockPublicAclS3 = "AWS074"
const AWSBlockPublicAclS3Description = "S3 Access block should block public ACL"
const AWSBlockPublicAclS3Impact = "PUT calls with public ACLs specified can make objects public"
const AWSBlockPublicAclS3Resolution = "Enable blocking any PUT calls with a public ACL specified"
const AWSBlockPublicAclS3Explanation = `
S3 buckets should block public ACLs on buckets and any objects they contain. By blocking, PUTs with fail if the object has any public ACL a.
`
const AWSBlockPublicAclS3BadExample = `
resource "aws_s3_bucket_public_access_block" "bad_example" {
	bucket = aws_s3_bucket.example.id
}

resource "aws_s3_bucket_public_access_block" "bad_example" {
	bucket = aws_s3_bucket.example.id
  
	block_public_acls = false
}
`
const AWSBlockPublicAclS3GoodExample = `
resource "aws_s3_bucket_public_access_block" "good_example" {
	bucket = aws_s3_bucket.example.id
  
	block_public_acls = true
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSBlockPublicAclS3,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSBlockPublicAclS3Description,
			Impact:      AWSBlockPublicAclS3Impact,
			Resolution:  AWSBlockPublicAclS3Resolution,
			Explanation: AWSBlockPublicAclS3Explanation,
			BadExample:  AWSBlockPublicAclS3BadExample,
			GoodExample: AWSBlockPublicAclS3GoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-block-public-access.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_public_access_block#block_public_acls",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_s3_bucket_public_access_block"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if block.MissingChild("block_public_acls") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not specify block_public_acls, defaults to false", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			attr := block.GetAttribute("block_public_acls")
			if attr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' sets block_public_acls explicitly to false", block.FullName())).
						WithRange(attr.Range()).
						WithAttributeAnnotation(attr).
						WithSeverity(severity.Error),
				)
			}
		},
	})
}
