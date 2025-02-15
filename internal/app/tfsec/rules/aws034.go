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

const AWSOutdatedTLSPolicyElasticsearchDomainEndpoint = "AWS034"
const AWSOutdatedTLSPolicyElasticsearchDomainEndpointDescription = "Elasticsearch domain endpoint is using outdated TLS policy."
const AWSOutdatedTLSPolicyElasticsearchDomainEndpointImpact = "Outdated SSL policies increase exposure to known vulnerabilites"
const AWSOutdatedTLSPolicyElasticsearchDomainEndpointResolution = "Use the most modern TLS/SSL policies available"
const AWSOutdatedTLSPolicyElasticsearchDomainEndpointExplanation = `
You should not use outdated/insecure TLS versions for encryption. You should be using TLS v1.2+.
`
const AWSOutdatedTLSPolicyElasticsearchDomainEndpointBadExample = `
resource "aws_elasticsearch_domain" "bad_example" {
  domain_name = "domain-foo"

  domain_endpoint_options {
    enforce_https = true
    tls_security_policy = "Policy-Min-TLS-1-0-2019-07"
  }
}
`
const AWSOutdatedTLSPolicyElasticsearchDomainEndpointGoodExample = `
resource "aws_elasticsearch_domain" "good_example" {
  domain_name = "domain-foo"

  domain_endpoint_options {
    enforce_https = true
    tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSOutdatedTLSPolicyElasticsearchDomainEndpoint,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSOutdatedTLSPolicyElasticsearchDomainEndpointDescription,
			Impact:      AWSOutdatedTLSPolicyElasticsearchDomainEndpointImpact,
			Resolution:  AWSOutdatedTLSPolicyElasticsearchDomainEndpointResolution,
			Explanation: AWSOutdatedTLSPolicyElasticsearchDomainEndpointExplanation,
			BadExample:  AWSOutdatedTLSPolicyElasticsearchDomainEndpointBadExample,
			GoodExample: AWSOutdatedTLSPolicyElasticsearchDomainEndpointGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticsearch_domain#tls_security_policy",
				"https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-data-protection.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticsearch_domain"},
		CheckFunc: func(set result.Set, block *block.Block, context *hclcontext.Context) {

			endpointBlock := block.GetBlock("domain_endpoint_options")
			if endpointBlock == nil {
				// Rule AWS033 covers this case.
				return
			}

			tlsPolicyAttr := endpointBlock.GetAttribute("tls_security_policy")
			if tlsPolicyAttr == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an Elasticsearch domain with an outdated TLS policy (defaults to Policy-Min-TLS-1-0-2019-07).", block.FullName())).
						WithRange(endpointBlock.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			if tlsPolicyAttr.Value().Equals(cty.StringVal("Policy-Min-TLS-1-0-2019-07")).True() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an Elasticsearch domain with an outdated TLS policy (set to Policy-Min-TLS-1-0-2019-07).", block.FullName())).
						WithRange(tlsPolicyAttr.Range()).
						WithAttributeAnnotation(tlsPolicyAttr).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
