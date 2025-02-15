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

// AWSUnenforcedHTTPSElasticsearchDomainEndpoint See
// https://github.com/tfsec/tfsec#included-checks for check info
const AWSUnenforcedHTTPSElasticsearchDomainEndpoint = "AWS033"
const AWSUnenforcedHTTPSElasticsearchDomainEndpointDescription = "Elasticsearch doesn't enforce HTTPS traffic."
const AWSUnenforcedHTTPSElasticsearchDomainEndpointImpact = "HTTP traffic can be intercepted and the contents read"
const AWSUnenforcedHTTPSElasticsearchDomainEndpointResolution = "Enforce the use of HTTPS for ElasticSearch"
const AWSUnenforcedHTTPSElasticsearchDomainEndpointExplanation = `
Plain HTTP is unencrypted and human-readable. This means that if a malicious actor was to eavesdrop on your connection, they would be able to see all of your data flowing back and forth.

You should use HTTPS, which is HTTP over an encrypted (TLS) connection, meaning eavesdroppers cannot read your traffic.
`
const AWSUnenforcedHTTPSElasticsearchDomainEndpointBadExample = `
resource "aws_elasticsearch_domain" "bad_example" {
  domain_name = "domain-foo"

  domain_endpoint_options {
    enforce_https = false
  }
}
`
const AWSUnenforcedHTTPSElasticsearchDomainEndpointGoodExample = `
resource "aws_elasticsearch_domain" "good_example" {
  domain_name = "domain-foo"

  domain_endpoint_options {
    enforce_https = true
  }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSUnenforcedHTTPSElasticsearchDomainEndpoint,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSUnenforcedHTTPSElasticsearchDomainEndpointDescription,
			Impact:      AWSUnenforcedHTTPSElasticsearchDomainEndpointImpact,
			Resolution:  AWSUnenforcedHTTPSElasticsearchDomainEndpointResolution,
			Explanation: AWSUnenforcedHTTPSElasticsearchDomainEndpointExplanation,
			BadExample:  AWSUnenforcedHTTPSElasticsearchDomainEndpointBadExample,
			GoodExample: AWSUnenforcedHTTPSElasticsearchDomainEndpointGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticsearch_domain#enforce_https",
				"https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-data-protection.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticsearch_domain"},
		CheckFunc: func(set result.Set, block *block.Block, context *hclcontext.Context) {

			endpointBlock := block.GetBlock("domain_endpoint_options")
			if endpointBlock == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an Elasticsearch domain with plaintext traffic (missing domain_endpoint_options block).", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			}

			enforceHTTPSAttr := endpointBlock.GetAttribute("enforce_https")
			if enforceHTTPSAttr == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an Elasticsearch domain with plaintext traffic (missing enforce_https attribute).", block.FullName())).
						WithRange(endpointBlock.Range()).
						WithSeverity(severity.Error),
				)
			}

			isTrueBool := enforceHTTPSAttr.Type() == cty.Bool && enforceHTTPSAttr.Value().True()
			isTrueString := enforceHTTPSAttr.Type() == cty.String &&
				enforceHTTPSAttr.Value().Equals(cty.StringVal("true")).True()
			enforcedHTTPS := isTrueBool || isTrueString
			if !enforcedHTTPS {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an Elasticsearch domain with plaintext traffic (enabled attribute set to false).", block.FullName())).
						WithRange(endpointBlock.Range()).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
