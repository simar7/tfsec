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

const AWSCloudtrailEncryptedAtRest = "AWS065"
const AWSCloudtrailEncryptedAtRestDescription = "Cloudtrail should be encrypted at rest to secure access to sensitive trail data"
const AWSCloudtrailEncryptedAtRestImpact = "Data can be freely read if compromised"
const AWSCloudtrailEncryptedAtRestResolution = "Enable encryption at rest"
const AWSCloudtrailEncryptedAtRestExplanation = `
Cloudtrail logs should be encrypted at rest to secure the sensitive data. Cloudtrail logs record all activity that occurs in the the account through API calls and would be one of the first places to look when reacting to a breach.
`
const AWSCloudtrailEncryptedAtRestBadExample = `
resource "aws_cloudtrail" "bad_example" {
  is_multi_region_trail = true

  event_selector {
    read_write_type           = "All"
    include_management_events = true

    data_resource {
      type = "AWS::S3::Object"
      values = ["${data.aws_s3_bucket.important-bucket.arn}/"]
    }
  }
}
`
const AWSCloudtrailEncryptedAtRestGoodExample = `
resource "aws_cloudtrail" "good_example" {
  is_multi_region_trail = true
  enable_log_file_validation = true
  kms_key_id = var.kms_id

  event_selector {
    read_write_type           = "All"
    include_management_events = true

    data_resource {
      type = "AWS::S3::Object"
      values = ["${data.aws_s3_bucket.important-bucket.arn}/"]
    }
  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSCloudtrailEncryptedAtRest,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSCloudtrailEncryptedAtRestDescription,
			Impact:      AWSCloudtrailEncryptedAtRestImpact,
			Resolution:  AWSCloudtrailEncryptedAtRestResolution,
			Explanation: AWSCloudtrailEncryptedAtRestExplanation,
			BadExample:  AWSCloudtrailEncryptedAtRestBadExample,
			GoodExample: AWSCloudtrailEncryptedAtRestGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudtrail#kms_key_id",
				"https://docs.aws.amazon.com/awscloudtrail/latest/userguide/encrypting-cloudtrail-log-files-with-aws-kms.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_cloudtrail"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("kms_key_id") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not have a kms_key_id set.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			kmsKeyIdAttr := block.GetAttribute("kms_key_id")
			if kmsKeyIdAttr.IsEmpty() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has a kms_key_id but it is not set.", block.FullName())).
						WithRange(kmsKeyIdAttr.Range()).
						WithAttributeAnnotation(kmsKeyIdAttr).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
