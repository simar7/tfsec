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

const AWSCheckLambdaFunctionForSourceARN = "AWS058"
const AWSCheckLambdaFunctionForSourceARNDescription = "Ensure that lambda function permission has a source arn specified"
const AWSCheckLambdaFunctionForSourceARNImpact = "Not providing the source ARN allows any resource from principal, even from other accounts"
const AWSCheckLambdaFunctionForSourceARNResolution = "Always provide a source arn for Lambda permissions"
const AWSCheckLambdaFunctionForSourceARNExplanation = `When the principal is an AWS service, the ARN of the specific resource within that service to grant permission to. 

Without this, any resource from principal will be granted permission – even if that resource is from another account. 

For S3, this should be the ARN of the S3 Bucket. For CloudWatch Events, this should be the ARN of the CloudWatch Events Rule. For API Gateway, this should be the ARN of the API`
const AWSCheckLambdaFunctionForSourceARNBadExample = `
resource "aws_lambda_permission" "bad_example" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.func.function_name
  principal     = "sns.amazonaws.com"
}
`
const AWSCheckLambdaFunctionForSourceARNGoodExample = `
resource "aws_lambda_permission" "good_example" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.func.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.default.arn
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSCheckLambdaFunctionForSourceARN,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSCheckLambdaFunctionForSourceARNDescription,
			Impact:      AWSCheckLambdaFunctionForSourceARNImpact,
			Resolution:  AWSCheckLambdaFunctionForSourceARNResolution,
			Explanation: AWSCheckLambdaFunctionForSourceARNExplanation,
			BadExample:  AWSCheckLambdaFunctionForSourceARNBadExample,
			GoodExample: AWSCheckLambdaFunctionForSourceARNGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_permission",
				"https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_lambda_permission"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.HasChild("principal") {
				if block.GetAttribute("principal").EndsWith("amazonaws.com") {
					if block.MissingChild("source_arn") {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("Resource '%s' missing source ARN but has *.amazonaws.com Principal.", block.FullName())).
								WithRange(block.Range()).
								WithSeverity(severity.Error),
						)
					}
				}
			}

		},
	})
}
