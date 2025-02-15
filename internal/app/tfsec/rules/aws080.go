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

const AWSCodeBuildProjectEncryptionNotDisabled = "AWS080"
const AWSCodeBuildProjectEncryptionNotDisabledDescription = "CodeBuild Project artifacts encryption should not be disabled"
const AWSCodeBuildProjectEncryptionNotDisabledImpact = "CodeBuild project artifacts are unencrypted"
const AWSCodeBuildProjectEncryptionNotDisabledResolution = "Enable encryption for CodeBuild project artifacts"
const AWSCodeBuildProjectEncryptionNotDisabledExplanation = `
All artifacts produced by your CodeBuild project pipeline should always be encrypted
`
const AWSCodeBuildProjectEncryptionNotDisabledBadExample = `
resource "aws_codebuild_project" "bad_example" {
	// other config

	artifacts {
		// other artifacts config

		encryption_disabled = true
	}
}

resource "aws_codebuild_project" "bad_example" {
	// other config including primary artifacts

	secondary_artifacts {
		// other artifacts config
		
		encryption_disabled = false
	}

	secondary_artifacts {
		// other artifacts config

		encryption_disabled = true
	}
}
`
const AWSCodeBuildProjectEncryptionNotDisabledGoodExample = `
resource "aws_codebuild_project" "good_example" {
	// other config

	artifacts {
		// other artifacts config

		encryption_disabled = false
	}
}

resource "aws_codebuild_project" "good_example" {
	// other config

	artifacts {
		// other artifacts config
	}
}

resource "aws_codebuild_project" "codebuild" {
	// other config

	secondary_artifacts {
		// other artifacts config

		encryption_disabled = false
	}

	secondary_artifacts {
		// other artifacts config
	}
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSCodeBuildProjectEncryptionNotDisabled,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSCodeBuildProjectEncryptionNotDisabledDescription,
			Impact:      AWSCodeBuildProjectEncryptionNotDisabledImpact,
			Resolution:  AWSCodeBuildProjectEncryptionNotDisabledResolution,
			Explanation: AWSCodeBuildProjectEncryptionNotDisabledExplanation,
			BadExample:  AWSCodeBuildProjectEncryptionNotDisabledBadExample,
			GoodExample: AWSCodeBuildProjectEncryptionNotDisabledGoodExample,
			Links: []string{
				"https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-codebuild-project.html",
				"https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-artifacts.html",
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/codebuild_project#encryption_disabled",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_codebuild_project"},
		CheckFunc: func(set result.Set, b *block.Block, _ *hclcontext.Context) {

			blocks := b.GetBlocks("secondary_artifacts")
			if artifact := b.GetBlock("artifacts"); artifact != nil {
				blocks = append(blocks, artifact)
			}

			for _, artifactBlock := range blocks {
				if encryptionDisabledAttr := artifactBlock.GetAttribute("encryption_disabled"); encryptionDisabledAttr != nil && encryptionDisabledAttr.IsTrue() {
					artifactTypeAttr := artifactBlock.GetAttribute("type")

					if artifactTypeAttr.Equals("NO_ARTIFACTS", block.IgnoreCase) {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("CodeBuild project '%s' is configured to disable artifact encryption while no artifacts are produced", b.FullName())).
								WithRange(artifactBlock.Range()).
								WithAttributeAnnotation(artifactTypeAttr).
								WithSeverity(severity.Warning),
						)
					} else {
						set.Add(
							result.New().
								WithDescription(fmt.Sprintf("CodeBuild project '%s' does not encrypt produced artifacts", b.FullName())).
								WithRange(artifactBlock.Range()).
								WithAttributeAnnotation(encryptionDisabledAttr).
								WithSeverity(severity.Error),
						)
					}
				}
			}
		},
	})
}
