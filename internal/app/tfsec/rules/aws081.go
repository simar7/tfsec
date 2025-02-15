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

const AWSDAXEncryptedAtRest = "AWS081"
const AWSDAXEncryptedAtRestDescription = "DAX Cluster should always encrypt data at rest"
const AWSDAXEncryptedAtRestImpact = "Data can be freely read if compromised"
const AWSDAXEncryptedAtRestResolution = "Enable encryption at rest for DAX Cluster"
const AWSDAXEncryptedAtRestExplanation = `
Amazon DynamoDB Accelerator (DAX) encryption at rest provides an additional layer of data protection by helping secure your data from unauthorized access to the underlying storage.
`
const AWSDAXEncryptedAtRestBadExample = `
resource "aws_dax_cluster" "bad_example" {
	// no server side encryption at all
}

resource "aws_dax_cluster" "bad_example" {
	// other DAX config

	server_side_encryption {
		// empty server side encryption config
	}
}

resource "aws_dax_cluster" "bad_example" {
	// other DAX config

	server_side_encryption {
		enabled = false // disabled server side encryption
	}
}
`
const AWSDAXEncryptedAtRestGoodExample = `
resource "aws_dax_cluster" "good_example" {
	// other DAX config

	server_side_encryption {
		enabled = true // enabled server side encryption
	}
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSDAXEncryptedAtRest,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSDAXEncryptedAtRestDescription,
			Impact:      AWSDAXEncryptedAtRestImpact,
			Resolution:  AWSDAXEncryptedAtRestResolution,
			Explanation: AWSDAXEncryptedAtRestExplanation,
			BadExample:  AWSDAXEncryptedAtRestBadExample,
			GoodExample: AWSDAXEncryptedAtRestGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DAXEncryptionAtRest.html",
				"https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dax-cluster.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dax_cluster#server_side_encryption",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_dax_cluster"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("server_side_encryption") {
				res := result.New().
					WithDescription(fmt.Sprintf("DAX cluster '%s' does not have server side encryption configured. By default it is disabled.", block.FullName())).
					WithRange(block.Range()).
					WithSeverity(severity.Error)
				set.Add(res)
				return
			}

			sseBlock := block.GetBlock("server_side_encryption")
			if sseBlock.MissingChild("enabled") {
				res := result.New().
					WithDescription(fmt.Sprintf("DAX cluster '%s' server side encryption block is empty. By default SSE is disabled.", block.FullName())).
					WithRange(sseBlock.Range()).
					WithSeverity(severity.Error)
				set.Add(res)
				return
			}

			if sseEnabledAttr := sseBlock.GetAttribute("enabled"); sseEnabledAttr.IsFalse() {
				res := result.New().
					WithDescription(fmt.Sprintf("DAX cluster '%s' has disabled server side encryption", block.FullName())).
					WithRange(sseEnabledAttr.Range()).
					WithAttributeAnnotation(sseEnabledAttr).
					WithSeverity(severity.Error)
				set.Add(res)
			}

		},
	})
}
