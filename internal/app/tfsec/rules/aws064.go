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

const AWSCloudtrailLogValidationEnabled = "AWS064"
const AWSCloudtrailLogValidationEnabledDescription = "Cloudtrail log validation should be enabled to prevent tampering of log data"
const AWSCloudtrailLogValidationEnabledImpact = "Illicit activity could be removed from the logs"
const AWSCloudtrailLogValidationEnabledResolution = "Turn on log validation for Cloudtrail"
const AWSCloudtrailLogValidationEnabledExplanation = `
Log validation should be activated on Cloudtrail logs to prevent the tampering of the underlying data in the S3 bucket. It is feasible that a rogue actor compromising an AWS account might want to modify the log data to remove trace of their actions.
`
const AWSCloudtrailLogValidationEnabledBadExample = `
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
const AWSCloudtrailLogValidationEnabledGoodExample = `
resource "aws_cloudtrail" "good_example" {
  is_multi_region_trail = true
  enable_log_file_validation = true

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
		ID: AWSCloudtrailLogValidationEnabled,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSCloudtrailLogValidationEnabledDescription,
			Impact:      AWSCloudtrailLogValidationEnabledImpact,
			Resolution:  AWSCloudtrailLogValidationEnabledResolution,
			Explanation: AWSCloudtrailLogValidationEnabledExplanation,
			BadExample:  AWSCloudtrailLogValidationEnabledBadExample,
			GoodExample: AWSCloudtrailLogValidationEnabledGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudtrail#enable_log_file_validation",
				"https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-log-file-validation-intro.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_cloudtrail"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if block.MissingChild("enable_log_file_validation") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not enable log file validation.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Warning),
				)
			}

			logFileValidationAttr := block.GetAttribute("enable_log_file_validation")
			if logFileValidationAttr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not enable log file validation.", block.FullName())).
						WithRange(logFileValidationAttr.Range()).
						WithAttributeAnnotation(logFileValidationAttr).
						WithSeverity(severity.Warning),
				)
			} /**/
		},
	})
}
