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

	"github.com/zclconf/go-cty/cty"
)

// AWSOutdatedSSLPolicy See https://github.com/tfsec/tfsec#included-checks for check info
const AWSOutdatedSSLPolicy = "AWS010"
const AWSOutdatedSSLPolicyDescription = "An outdated SSL policy is in use by a load balancer."
const AWSOutdatedSSLPolicyImpact = "The SSL policy is outdated and has known vulnerabilities"
const AWSOutdatedSSLPolicyResolution = "Use a more recent TLS/SSL policy for the load balancer"
const AWSOutdatedSSLPolicyExplanation = `
You should not use outdated/insecure TLS versions for encryption. You should be using TLS v1.2+. 
`
const AWSOutdatedSSLPolicyBadExample = `
resource "aws_alb_listener" "bad_example" {
	ssl_policy = "ELBSecurityPolicy-TLS-1-1-2017-01"
	protocol = "HTTPS"
}
`
const AWSOutdatedSSLPolicyGoodExample = `
resource "aws_alb_listener" "good_example" {
	ssl_policy = "ELBSecurityPolicy-TLS-1-2-2017-01"
	protocol = "HTTPS"
}
`

var outdatedSSLPolicies = []string{
	"ELBSecurityPolicy-2015-05",
	"ELBSecurityPolicy-TLS-1-0-2015-04",
	"ELBSecurityPolicy-2016-08",
	"ELBSecurityPolicy-TLS-1-1-2017-01",
}

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSOutdatedSSLPolicy,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSOutdatedSSLPolicyDescription,
			Impact:      AWSOutdatedSSLPolicyImpact,
			Resolution:  AWSOutdatedSSLPolicyResolution,
			Explanation: AWSOutdatedSSLPolicyExplanation,
			BadExample:  AWSOutdatedSSLPolicyBadExample,
			GoodExample: AWSOutdatedSSLPolicyGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_lb_listener", "aws_alb_listener"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if sslPolicyAttr := block.GetAttribute("ssl_policy"); sslPolicyAttr != nil && sslPolicyAttr.Type() == cty.String {
				for _, policy := range outdatedSSLPolicies {
					if policy == sslPolicyAttr.Value().AsString() {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("Resource '%s' is using an outdated SSL policy.", block.FullName())).
								WithRange(sslPolicyAttr.Range()).
								WithAttributeAnnotation(sslPolicyAttr).
								WithSeverity(severity.Error),
						)
					}
				}
			}

		},
	})
}
