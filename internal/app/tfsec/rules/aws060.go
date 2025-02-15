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

const AWSAthenaWorkgroupEnforceConfiguration = "AWS060"
const AWSAthenaWorkgroupEnforceConfigurationDescription = "Athena workgroups should enforce configuration to prevent client disabling encryption"
const AWSAthenaWorkgroupEnforceConfigurationImpact = "Clients can ginore encryption requirements"
const AWSAthenaWorkgroupEnforceConfigurationResolution = "Enforce the configuration to prevent client overrides"
const AWSAthenaWorkgroupEnforceConfigurationExplanation = `
Athena workgroup configuration should be enforced to prevent client side changes to disable encryption settings.
`
const AWSAthenaWorkgroupEnforceConfigurationBadExample = `
resource "aws_athena_workgroup" "bad_example" {
  name = "example"

  configuration {
    enforce_workgroup_configuration    = false
    publish_cloudwatch_metrics_enabled = true

    result_configuration {
      output_location = "s3://${aws_s3_bucket.example.bucket}/output/"

      encryption_configuration {
        encryption_option = "SSE_KMS"
        kms_key_arn       = aws_kms_key.example.arn
      }
    }
  }
}

resource "aws_athena_workgroup" "bad_example" {
  name = "example"

}
`
const AWSAthenaWorkgroupEnforceConfigurationGoodExample = `
resource "aws_athena_workgroup" "good_example" {
  name = "example"

  configuration {
    enforce_workgroup_configuration    = true
    publish_cloudwatch_metrics_enabled = true

    result_configuration {
      output_location = "s3://${aws_s3_bucket.example.bucket}/output/"

      encryption_configuration {
        encryption_option = "SSE_KMS"
        kms_key_arn       = aws_kms_key.example.arn
      }
    }
  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSAthenaWorkgroupEnforceConfiguration,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSAthenaWorkgroupEnforceConfigurationDescription,
			Impact:      AWSAthenaWorkgroupEnforceConfigurationImpact,
			Resolution:  AWSAthenaWorkgroupEnforceConfigurationResolution,
			Explanation: AWSAthenaWorkgroupEnforceConfigurationExplanation,
			BadExample:  AWSAthenaWorkgroupEnforceConfigurationBadExample,
			GoodExample: AWSAthenaWorkgroupEnforceConfigurationGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/athena_workgroup#configuration",
				"https://docs.aws.amazon.com/athena/latest/ug/manage-queries-control-costs-with-workgroups.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_athena_workgroup"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("configuration") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' is missing the configuration block.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			configBlock := block.GetBlock("configuration")
			if configBlock.HasChild("enforce_workgroup_configuration") &&
				configBlock.GetAttribute("enforce_workgroup_configuration").IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has enforce_workgroup_configuration set to false.", block.FullName())).
WithRange(configBlock.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
