package rules

import (
	"fmt"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/pkg/provider"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/zclconf/go-cty/cty"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

const (
	AWSIAMPasswordRequiresNumber            = "AWS041"
	AWSIAMPasswordRequiresNumberDescription = "IAM Password policy should have requirement for at least one number in the password."
	AWSIAMPasswordRequiresNumberImpact      = "Short, simple passwords are easier to compromise"
	AWSIAMPasswordRequiresNumberResolution  = "Enforce longer, more complex passwords in the policy"
	AWSIAMPasswordRequiresNumberExplanation = `
IAM account password policies should ensure that passwords content including at least one number.
`
	AWSIAMPasswordRequiresNumberBadExample = `
resource "aws_iam_account_password_policy" "bad_example" {
	# ...
	# require_numbers not set
	# ...
}
`
	AWSIAMPasswordRequiresNumberGoodExample = `
resource "aws_iam_account_password_policy" "good_example" {
	# ...
	require_numbers = true
	# ...
}
`
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSIAMPasswordRequiresNumber,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSIAMPasswordRequiresNumberDescription,
			Impact:      AWSIAMPasswordRequiresNumberImpact,
			Resolution:  AWSIAMPasswordRequiresNumberResolution,
			Explanation: AWSIAMPasswordRequiresNumberExplanation,
			BadExample:  AWSIAMPasswordRequiresNumberBadExample,
			GoodExample: AWSIAMPasswordRequiresNumberGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_passwords_account-policy.html#password-policy-details",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_account_password_policy",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_iam_account_password_policy"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if attr := block.GetAttribute("require_numbers"); attr == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not require a number in the password.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Warning),
				)
			} else if attr.Value().Type() == cty.Bool {
				if attr.Value().False() {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' explicitly specifies not requiring at least one number in the password.", block.FullName())).
							WithRange(block.Range()).
							WithSeverity(severity.Warning),
					)
				}
			}
		},
	})
}
