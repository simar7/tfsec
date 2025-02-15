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

const AWSUnencryptedInTransitElasticacheReplicationGroup = "AWS036"
const AWSUnencryptedInTransitElasticacheReplicationGroupDescription = "Elasticache Replication Group uses unencrypted traffic."
const AWSUnencryptedInTransitElasticacheReplicationGroupImpact = "In transit data in the Replication Group could be read if intercepted"
const AWSUnencryptedInTransitElasticacheReplicationGroupResolution = "Enable in transit encryptuon for replication group"
const AWSUnencryptedInTransitElasticacheReplicationGroupExplanation = `
Traffic flowing between Elasticache replication nodes should be encrypted to ensure sensitive data is kept private.
`
const AWSUnencryptedInTransitElasticacheReplicationGroupBadExample = `
resource "aws_elasticache_replication_group" "bad_example" {
        replication_group_id = "foo"
        replication_group_description = "my foo cluster"

        transit_encryption_enabled = false
}
`
const AWSUnencryptedInTransitElasticacheReplicationGroupGoodExample = `
resource "aws_elasticache_replication_group" "good_example" {
        replication_group_id = "foo"
        replication_group_description = "my foo cluster"

        transit_encryption_enabled = true
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSUnencryptedInTransitElasticacheReplicationGroup,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSUnencryptedInTransitElasticacheReplicationGroupDescription,
			Impact:      AWSUnencryptedInTransitElasticacheReplicationGroupImpact,
			Resolution:  AWSUnencryptedInTransitElasticacheReplicationGroupResolution,
			Explanation: AWSUnencryptedInTransitElasticacheReplicationGroupExplanation,
			BadExample:  AWSUnencryptedInTransitElasticacheReplicationGroupBadExample,
			GoodExample: AWSUnencryptedInTransitElasticacheReplicationGroupGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_replication_group#transit_encryption_enabled",
				"https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/in-transit-encryption.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticache_replication_group"},
		CheckFunc: func(set result.Set, block *block.Block, context *hclcontext.Context) {

			encryptionAttr := block.GetAttribute("transit_encryption_enabled")
			if encryptionAttr == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an unencrypted Elasticache Replication Group (missing transit_encryption_enabled attribute).", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			} else if !isBooleanOrStringTrue(encryptionAttr) {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an unencrypted Elasticache Replication Group (transit_encryption_enabled set to false).", block.FullName())).
						WithRange(encryptionAttr.Range()).
						WithAttributeAnnotation(encryptionAttr).
						WithSeverity(severity.Error),
				)

			}

		},
	})
}
