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

const AWSALBDropsInvalidHeaders = "AWS083"
const AWSALBDropsInvalidHeadersDescription = "Load balancers should drop invalid headers"
const AWSALBDropsInvalidHeadersImpact = "Invalid headers being passed through to the target of the load balance may exploit vulnerabilities"
const AWSALBDropsInvalidHeadersResolution = "Set drop_invalid_header_fields to true"
const AWSALBDropsInvalidHeadersExplanation = `
Passing unknown or invalid headers through to the target poses a potential risk of compromise. 

By setting drop_invalid_header_fields to true, anything that doe not conform to well known, defined headers will be removed by the load balancer.
`
const AWSALBDropsInvalidHeadersBadExample = `
resource "aws_alb" "bad_example" {
	name               = "bad_alb"
	internal           = false
	load_balancer_type = "application"
	
	access_logs {
	  bucket  = aws_s3_bucket.lb_logs.bucket
	  prefix  = "test-lb"
	  enabled = true
	}
  
	drop_invalid_header_fields = false
  }
`
const AWSALBDropsInvalidHeadersGoodExample = `
resource "aws_alb" "good_example" {
	name               = "good_alb"
	internal           = false
	load_balancer_type = "application"
	
	access_logs {
	  bucket  = aws_s3_bucket.lb_logs.bucket
	  prefix  = "test-lb"
	  enabled = true
	}
  
	drop_invalid_header_fields = true
  }
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSALBDropsInvalidHeaders,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSALBDropsInvalidHeadersDescription,
			Explanation: AWSALBDropsInvalidHeadersExplanation,
			Impact:      AWSALBDropsInvalidHeadersImpact,
			Resolution:  AWSALBDropsInvalidHeadersResolution,
			BadExample:  AWSALBDropsInvalidHeadersBadExample,
			GoodExample: AWSALBDropsInvalidHeadersGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb#drop_invalid_header_fields",
				"https://docs.aws.amazon.com/elasticloadbalancing/latest/application/application-load-balancers.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_alb", "aws_lb"},
		CheckFunc: func(set result.Set, b *block.Block, _ *hclcontext.Context) {

			if b.GetAttribute("load_balancer_type").Equals("application", block.IgnoreCase) {
				if b.MissingChild("drop_invalid_header_fields") {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' does not drop invalid header fields", b.FullName())).
							WithRange(b.Range()).
							WithSeverity(severity.Error),
					)
				}

				attr := b.GetAttribute("drop_invalid_header_fields")
				if attr.IsFalse() {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' sets the drop_invalid_header_fields to false", b.FullName())).
							WithRange(attr.Range()).
							WithAttributeAnnotation(attr).
							WithSeverity(severity.Error),
					)
				}

			}
		},
	})
}
