package rules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/pkg/provider"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/zclconf/go-cty/cty"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

const AWSSqsPolicyWildcardActions = "AWS047"
const AWSSqsPolicyWildcardActionsDescription = "AWS SQS policy document has wildcard action statement."
const AWSSqsPolicyWildcardActionsImpact = "SQS policies with wildcard actions allow more that is required"
const AWSSqsPolicyWildcardActionsResolution = "Keep policy scope to the minimum that is required to be effective"
const AWSSqsPolicyWildcardActionsExplanation = `
SQS Policy actions should always be restricted to a specific set.

This ensures that the queue itself cannot be modified or deleted, and prevents possible future additions to queue actions to be implicitly allowed.
`
const AWSSqsPolicyWildcardActionsBadExample = `
resource "aws_sqs_queue_policy" "bad_example" {
  queue_url = aws_sqs_queue.q.id

  policy = <<POLICY
{
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*"
    }
  ]
}
POLICY
}
`
const AWSSqsPolicyWildcardActionsGoodExample = `
resource "aws_sqs_queue_policy" "good_example" {
  queue_url = aws_sqs_queue.q.id

  policy = <<POLICY
{
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage"
    }
  ]
}
POLICY
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSSqsPolicyWildcardActions,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSSqsPolicyWildcardActionsDescription,
			Impact:      AWSSqsPolicyWildcardActionsImpact,
			Resolution:  AWSSqsPolicyWildcardActionsResolution,
			Explanation: AWSSqsPolicyWildcardActionsExplanation,
			BadExample:  AWSSqsPolicyWildcardActionsBadExample,
			GoodExample: AWSSqsPolicyWildcardActionsGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-security-best-practices.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_sqs_queue_policy"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.GetAttribute("policy").Value().Type() != cty.String {
			}

			rawJSON := []byte(block.GetAttribute("policy").Value().AsString())
			var policy struct {
				Statement []struct {
					Effect string `json:"Effect"`
					Action string `json:"Action"`
				} `json:"Statement"`
			}

			if err := json.Unmarshal(rawJSON, &policy); err == nil {
				for _, statement := range policy.Statement {
					if strings.ToLower(statement.Effect) == "allow" && (statement.Action == "*" || statement.Action == "sqs:*") {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("SQS policy '%s' has a wildcard action specified.", block.FullName())).
								WithRange(block.Range()).
								WithSeverity(severity.Error),
						)
					}
				}
			}

		},
	})
}
