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

const AWSRDSAuroraClusterEncryptionDisabled = "AWS051"
const AWSRDSAuroraClusterEncryptionDisabledDescription = "There is no encryption specified or encryption is disabled on the RDS Cluster."
const AWSRDSAuroraClusterEncryptionDisabledImpact = "Data can be read from the RDS cluster if it is compromised"
const AWSRDSAuroraClusterEncryptionDisabledResolution = "Enable encryption for RDS clusters and instances"
const AWSRDSAuroraClusterEncryptionDisabledExplanation = `
Encryption should be enabled for an RDS Aurora cluster. 

When enabling encryption by setting the kms_key_id, the storage_encrypted must also be set to true. 
`
const AWSRDSAuroraClusterEncryptionDisabledBadExample = `
resource "aws_rds_cluster" "bad_example" {
  name       = "bar"
  kms_key_id = ""
}`
const AWSRDSAuroraClusterEncryptionDisabledGoodExample = `
resource "aws_rds_cluster" "good_example" {
  name              = "bar"
  kms_key_id  = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
  storage_encrypted = true
}`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSRDSAuroraClusterEncryptionDisabled,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSRDSAuroraClusterEncryptionDisabledDescription,
			Impact:      AWSRDSAuroraClusterEncryptionDisabledImpact,
			Resolution:  AWSRDSAuroraClusterEncryptionDisabledResolution,
			Explanation: AWSRDSAuroraClusterEncryptionDisabledExplanation,
			BadExample:  AWSRDSAuroraClusterEncryptionDisabledBadExample,
			GoodExample: AWSRDSAuroraClusterEncryptionDisabledGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_rds_cluster"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			kmsKeyIdAttr := block.GetAttribute("kms_key_id")
			storageEncryptedattr := block.GetAttribute("storage_encrypted")

			if (kmsKeyIdAttr == nil || kmsKeyIdAttr.IsEmpty()) &&
				(storageEncryptedattr == nil || storageEncryptedattr.IsFalse()) {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a disabled RDS Cluster encryption.", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
			} else if kmsKeyIdAttr.Equals("") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a disabled RDS Cluster encryption.", block.FullName())).
						WithRange(kmsKeyIdAttr.Range()).
						WithAttributeAnnotation(kmsKeyIdAttr).
						WithSeverity(severity.Error),
				)
			} else if storageEncryptedattr == nil || storageEncryptedattr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' defines a enabled RDS Cluster encryption but not the required encrypted_storage.", block.FullName())).
						WithRange(kmsKeyIdAttr.Range()).
						WithAttributeAnnotation(kmsKeyIdAttr).
						WithSeverity(severity.Error),
				)
			}
		},
	})
}
