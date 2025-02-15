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

const AWSEKSSecretsEncryptionEnabled = "AWS066"
const AWSEKSSecretsEncryptionEnabledDescription = "EKS should have the encryption of secrets enabled"
const AWSEKSSecretsEncryptionEnabledImpact = "EKS secrets could be read if compromised"
const AWSEKSSecretsEncryptionEnabledResolution = "Enable encryption of EKS secrets"
const AWSEKSSecretsEncryptionEnabledExplanation = `
EKS cluster resources should have the encryption_config block set with protection of the secrets resource.
`
const AWSEKSSecretsEncryptionEnabledBadExample = `
resource "aws_eks_cluster" "bad_example" {
    name = "bad_example_cluster"

    role_arn = var.cluster_arn
    vpc_config {
        endpoint_public_access = false
    }
}
`
const AWSEKSSecretsEncryptionEnabledGoodExample = `
resource "aws_eks_cluster" "good_example" {
    encryption_config {
        resources = [ "secrets" ]
        provider {
            key_arn = var.kms_arn
        }
    }

    name = "good_example_cluster"
    role_arn = var.cluster_arn
    vpc_config {
        endpoint_public_access = false
    }
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSEKSSecretsEncryptionEnabled,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSEKSSecretsEncryptionEnabledDescription,
			Impact:      AWSEKSSecretsEncryptionEnabledImpact,
			Resolution:  AWSEKSSecretsEncryptionEnabledResolution,
			Explanation: AWSEKSSecretsEncryptionEnabledExplanation,
			BadExample:  AWSEKSSecretsEncryptionEnabledBadExample,
			GoodExample: AWSEKSSecretsEncryptionEnabledGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster#encryption_config",
				"https://aws.amazon.com/about-aws/whats-new/2020/03/amazon-eks-adds-envelope-encryption-for-secrets-with-aws-kms/",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_eks_cluster"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("encryption_config") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has no encryptionConfigBlock block", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			encryptionConfigBlock := block.GetBlock("encryption_config")
			if encryptionConfigBlock.MissingChild("resources") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has encryptionConfigBlock block with no resourcesAttr attribute specified", block.FullName())).
						WithRange(encryptionConfigBlock.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			resourcesAttr := encryptionConfigBlock.GetAttribute("resources")
			if !resourcesAttr.Contains("secrets") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' does not include secrets in encrypted resources", block.FullName())).
						WithRange(resourcesAttr.Range()).
						WithAttributeAnnotation(resourcesAttr).
						WithSeverity(severity.Error),
				)
			}

			if encryptionConfigBlock.MissingChild("provider") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has encryptionConfigBlock block with no provider block specified", block.FullName())).
						WithRange(encryptionConfigBlock.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			providerBlock := encryptionConfigBlock.GetBlock("provider")
			if providerBlock.MissingChild("key_arn") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has encryptionConfigBlock block with provider block specified missing key arn", block.FullName())).
						WithRange(encryptionConfigBlock.Range()).
						WithSeverity(severity.Error),
				)
				return
			}

			keyArnAttr := providerBlock.GetAttribute("key_arn")
			if keyArnAttr.IsEmpty() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has encryptionConfigBlock block with provider block specified but key_arn is empty", block.FullName())).
						WithRange(keyArnAttr.Range()).
						WithAttributeAnnotation(keyArnAttr).
						WithSeverity(severity.Error),
				)
			}

		},
	})
}
