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

const AWSUnencryptedAtRestElasticacheReplicationGroup = "AWS035"
const AWSUnencryptedAtRestElasticacheReplicationGroupDescription = "Unencrypted Elasticache Replication Group."
const AWSUnencryptedAtRestElasticacheReplicationGroupImpact = "Data in the replication group could be readable if compromised"
const AWSUnencryptedAtRestElasticacheReplicationGroupResolution = "Enable encryption for replication group"
const AWSUnencryptedAtRestElasticacheReplicationGroupExplanation = `
You should ensure your Elasticache data is encrypted at rest to help prevent sensitive information from being read by unauthorised users.
`
const AWSUnencryptedAtRestElasticacheReplicationGroupBadExample = `
resource "aws_elasticache_replication_group" "bad_example" {
        replication_group_id = "foo"
        replication_group_description = "my foo cluster"

        at_rest_encryption_enabled = false
}
`
const AWSUnencryptedAtRestElasticacheReplicationGroupGoodExample = `
resource "aws_elasticache_replication_group" "good_example" {
        replication_group_id = "foo"
        replication_group_description = "my foo cluster"

        at_rest_encryption_enabled = true
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSUnencryptedAtRestElasticacheReplicationGroup,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSUnencryptedAtRestElasticacheReplicationGroupDescription,
			Impact:      AWSUnencryptedAtRestElasticacheReplicationGroupImpact,
			Resolution:  AWSUnencryptedAtRestElasticacheReplicationGroupResolution,
			Explanation: AWSUnencryptedAtRestElasticacheReplicationGroupExplanation,
			BadExample:  AWSUnencryptedAtRestElasticacheReplicationGroupBadExample,
			GoodExample: AWSUnencryptedAtRestElasticacheReplicationGroupGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_replication_group#at_rest_encryption_enabled",
				"https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/at-rest-encryption.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticache_replication_group"},
		CheckFunc: func(set result.Set, block *block.Block, context *hclcontext.Context) {

			encryptionAttr := block.GetAttribute("at_rest_encryption_enabled")
			if encryptionAttr == nil {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an unencrypted Elasticache Replication Group (missing at_rest_encryption_enabled attribute).", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			} else if !isBooleanOrStringTrue(encryptionAttr) {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines an unencrypted Elasticache Replication Group (at_rest_encryption_enabled set to false).", block.FullName())).
						WithRange(encryptionAttr.Range()).
						WithAttributeAnnotation(encryptionAttr).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
