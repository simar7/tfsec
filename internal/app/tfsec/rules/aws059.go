package rules

import (
	"fmt"
	"strings"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/pkg/provider"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

const AWSEnsureAthenaDbEncrypted = "AWS059"
const AWSEnsureAthenaDbEncryptedDescription = "Athena databases and workgroup configurations are created unencrypted at rest by default, they should be encrypted"
const AWSEnsureAthenaDbEncryptedImpact = "Data can be read if the Athena Database is compromised"
const AWSEnsureAthenaDbEncryptedResolution = "Enable encryption at rest for Athena databases and workgroup configurations"
const AWSEnsureAthenaDbEncryptedExplanation = `
Athena databases and workspace result sets should be encrypted at rests. These databases and query sets are generally derived from data in S3 buckets and should have the same level of at rest protection.

`
const AWSEnsureAthenaDbEncryptedBadExample = `
resource "aws_athena_database" "bad_example" {
  name   = "database_name"
  bucket = aws_s3_bucket.hoge.bucket
}

resource "aws_athena_workgroup" "bad_example" {
  name = "example"

  configuration {
    enforce_workgroup_configuration    = true
    publish_cloudwatch_metrics_enabled = true

    result_configuration {
      output_location = "s3://${aws_s3_bucket.example.bucket}/output/"
    }
  }
}
`
const AWSEnsureAthenaDbEncryptedGoodExample = `
resource "aws_athena_database" "good_example" {
  name   = "database_name"
  bucket = aws_s3_bucket.hoge.bucket

  encryption_configuration {
     encryption_option = "SSE_KMS"
     kms_key_arn       = aws_kms_key.example.arn
 }
}

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
		ID: AWSEnsureAthenaDbEncrypted,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSEnsureAthenaDbEncryptedDescription,
			Impact:      AWSEnsureAthenaDbEncryptedImpact,
			Resolution:  AWSEnsureAthenaDbEncryptedResolution,
			Explanation: AWSEnsureAthenaDbEncryptedExplanation,
			BadExample:  AWSEnsureAthenaDbEncryptedBadExample,
			GoodExample: AWSEnsureAthenaDbEncryptedGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/athena_workgroup#encryption_configuration",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/athena_database#encryption_configuration",
				"https://docs.aws.amazon.com/athena/latest/ug/encryption.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_athena_database", "aws_athena_workgroup"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			blockName := block.FullName()

			if strings.EqualFold(block.TypeLabel(), "aws_athena_workgroup") {
				if block.HasChild("configuration") && block.GetBlock("configuration").
					HasChild("result_configuration") {
					block = block.GetBlock("configuration").GetBlock("result_configuration")
				} else {
					return
				}
			}

			if block.MissingChild("encryption_configuration") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' missing encryption configuration block.", blockName)).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
