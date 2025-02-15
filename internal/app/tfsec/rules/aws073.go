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

const AWSIngorePublicAclS3 = "AWS073"
const AWSIngorePublicAclS3Description = "S3 Access Block should Ignore Public Acl"
const AWSIngorePublicAclS3Impact = "PUT calls with public ACLs specified can make objects public"
const AWSIngorePublicAclS3Resolution = "Enable ignoring the application of public ACLs in PUT calls"
const AWSIngorePublicAclS3Explanation = `
S3 buckets should ignore public ACLs on buckets and any objects they contain. By ignoring rather than blocking, PUT calls with public ACLs will still be applied but the ACL will be ignored.
`
const AWSIngorePublicAclS3BadExample = `
resource "aws_s3_bucket_public_access_block" "bad_example" {
	bucket = aws_s3_bucket.example.id
}

resource "aws_s3_bucket_public_access_block" "bad_example" {
	bucket = aws_s3_bucket.example.id
  
	ignore_public_acls = false
}
`
const AWSIngorePublicAclS3GoodExample = `
resource "aws_s3_bucket_public_access_block" "good_example" {
	bucket = aws_s3_bucket.example.id
  
	ignore_public_acls = true
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSIngorePublicAclS3,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSIngorePublicAclS3Description,
			Impact:      AWSIngorePublicAclS3Impact,
			Resolution:  AWSIngorePublicAclS3Resolution,
			Explanation: AWSIngorePublicAclS3Explanation,
			BadExample:  AWSIngorePublicAclS3BadExample,
			GoodExample: AWSIngorePublicAclS3GoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-block-public-access.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_public_access_block#ignore_public_acls",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_s3_bucket_public_access_block"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("ignore_public_acls") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not specify ignore_public_acls, defaults to false", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			attr := block.GetAttribute("ignore_public_acls")
			if attr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' sets ignore_public_acls explicitly to false", block.FullName())).
						WithRange(attr.Range()).
						WithAttributeAnnotation(attr).
						WithSeverity(severity.Error),
				)
			}
		},
	})
}
