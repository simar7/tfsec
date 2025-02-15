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

const AWSRedisClusterBackupRetention = "AWS088"
const AWSRedisClusterBackupRetentionDescription = "Redis cluster should be backup retention turned on"
const AWSRedisClusterBackupRetentionImpact = "Without backups of the redis cluster recovery is made difficult"
const AWSRedisClusterBackupRetentionResolution = "Configure snapshot retention for redis cluster"
const AWSRedisClusterBackupRetentionExplanation = `
Redis clusters should have a snapshot retention time to ensure that they are backed up and can be restored if required.
`
const AWSRedisClusterBackupRetentionBadExample = `
resource "aws_elasticache_cluster" "bad_example" {
	cluster_id           = "cluster-example"
	engine               = "redis"
	node_type            = "cache.m4.large"
	num_cache_nodes      = 1
	parameter_group_name = "default.redis3.2"
	engine_version       = "3.2.10"
	port                 = 6379
}
`
const AWSRedisClusterBackupRetentionGoodExample = `
resource "aws_elasticache_cluster" "good_example" {
	cluster_id           = "cluster-example"
	engine               = "redis"
	node_type            = "cache.m4.large"
	num_cache_nodes      = 1
	parameter_group_name = "default.redis3.2"
	engine_version       = "3.2.10"
	port                 = 6379

	snapshot_retention_limit = 5
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSRedisClusterBackupRetention,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSRedisClusterBackupRetentionDescription,
			Explanation: AWSRedisClusterBackupRetentionExplanation,
			Impact:      AWSRedisClusterBackupRetentionImpact,
			Resolution:  AWSRedisClusterBackupRetentionResolution,
			BadExample:  AWSRedisClusterBackupRetentionBadExample,
			GoodExample: AWSRedisClusterBackupRetentionGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/elasticache_cluster#snapshot_retention_limit",
				"https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/backups-automatic.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_elasticache_cluster"},
		CheckFunc: func(set result.Set, b *block.Block, _ *hclcontext.Context) {

			engineAttr := b.GetAttribute("engine")
			if engineAttr.Equals("redis", block.IgnoreCase) && !b.GetAttribute("node_type").Equals("cache.t1.micro") {
				snapshotRetentionAttr := b.GetAttribute("snapshot_retention_limit")
				if snapshotRetentionAttr == nil {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' should have snapshot retention specified", b.FullName())).
							WithRange(b.Range()).
							WithSeverity(severity.Warning),
					)
				}

				if snapshotRetentionAttr.Equals(0) {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' has snapshot retention set to 0", b.FullName())).
							WithRange(snapshotRetentionAttr.Range()).
							WithAttributeAnnotation(snapshotRetentionAttr).
							WithSeverity(severity.Warning),
					)
				}
			}

		},
	})
}
