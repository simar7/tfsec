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

const AWSOpenAllIngressNetworkACLRule = "AWS050"
const AWSOpenAllIngressNetworkACLRuleDescription = "An ingress Network ACL rule allows ALL ports from /0."
const AWSOpenAllIngressNetworkACLRuleImpact = "All ports exposed for egressing data to the internet"
const AWSOpenAllIngressNetworkACLRuleResolution = "Set a more restrictive cidr range"
const AWSOpenAllIngressNetworkACLRuleExplanation = `
Opening up ACLs to the public internet is potentially dangerous. You should restrict access to IP addresses or ranges that explicitly require it where possible, and ensure that you specify required ports.

`
const AWSOpenAllIngressNetworkACLRuleBadExample = `
resource "aws_network_acl_rule" "bad_example" {
  egress         = false
  protocol       = "all"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}
`
const AWSOpenAllIngressNetworkACLRuleGoodExample = `
resource "aws_network_acl_rule" "good_example" {
  egress         = false
  protocol       = "tcp"
  from_port      = 22
  to_port        = 22
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSOpenAllIngressNetworkACLRule,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSOpenAllIngressNetworkACLRuleDescription,
			Impact:      AWSOpenAllIngressNetworkACLRuleImpact,
			Resolution:  AWSOpenAllIngressNetworkACLRuleResolution,
			Explanation: AWSOpenAllIngressNetworkACLRuleExplanation,
			BadExample:  AWSOpenAllIngressNetworkACLRuleBadExample,
			GoodExample: AWSOpenAllIngressNetworkACLRuleGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/vpc/latest/userguide/vpc-network-acls.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_network_acl_rule"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			egressAttr := block.GetAttribute("egress")
			actionAttr := block.GetAttribute("rule_action")
			protoAttr := block.GetAttribute("protocol")

			if egressAttr.Type() == cty.Bool && egressAttr.Value().True() {
			}

			if actionAttr == nil || actionAttr.Type() != cty.String {
			}

			if actionAttr.Value().AsString() != "allow" {
			}

			if cidrBlockAttr := block.GetAttribute("cidr_block"); cidrBlockAttr != nil {

				if isOpenCidr(cidrBlockAttr) {
					if protoAttr.Value().AsString() == "all" || protoAttr.Value().AsString() == "-1" {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("Resource '%s' defines a fully open ingress Network ACL rule with ALL ports open.", block.FullName())).
								WithRange(cidrBlockAttr.Range()).
								WithAttributeAnnotation(cidrBlockAttr).
								WithSeverity(severity.Error),
						)
					} else {
					}
				}

			}

			if ipv6CidrBlockAttr := block.GetAttribute("ipv6_cidr_block"); ipv6CidrBlockAttr != nil {

				if isOpenCidr(ipv6CidrBlockAttr) {
					if protoAttr.Value().AsString() == "all" || protoAttr.Value().AsString() == "-1" {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("Resource '%s' defines a fully open ingress Network ACL rule with ALL ports open.", block.FullName())).
								WithRange(ipv6CidrBlockAttr.Range()).
								WithAttributeAnnotation(ipv6CidrBlockAttr).
								WithSeverity(severity.Error),
						)
					} else {
					}
				}

			}

		},
	})
}
