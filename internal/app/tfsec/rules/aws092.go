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

const AWSDynamoDBTableEncryption = "AWS092"
const AWSDynamoDBTableEncryptionDescription = "DynamoDB tables should use at rest encyption with a Customer Managed Key"
const AWSDynamoDBTableEncryptionImpact = "Using AWS managed keys does not allow for fine grained control"
const AWSDynamoDBTableEncryptionResolution = "Enable server side encrytion with a customer managed key"
const AWSDynamoDBTableEncryptionExplanation = `
DynamoDB tables are encrypted by default using AWS managed encryption keys. To increase control of the encryption and control the management of factors like key rotation, use a Customer Managed Key.
`
const AWSDynamoDBTableEncryptionBadExample = `
resource "aws_dynamodb_table" "bad_example" {
	name             = "example"
	hash_key         = "TestTableHashKey"
	billing_mode     = "PAY_PER_REQUEST"
	stream_enabled   = true
	stream_view_type = "NEW_AND_OLD_IMAGES"
  
	attribute {
	  name = "TestTableHashKey"
	  type = "S"
	}
  
	replica {
	  region_name = "us-east-2"
	}
  
	replica {
	  region_name = "us-west-2"
	}
  }
`
const AWSDynamoDBTableEncryptionGoodExample = `
resource "aws_kms_key" "dynamo_db_kms" {
	enable_key_rotation = true
}

resource "aws_dynamodb_table" "good_example" {
	name             = "example"
	hash_key         = "TestTableHashKey"
	billing_mode     = "PAY_PER_REQUEST"
	stream_enabled   = true
	stream_view_type = "NEW_AND_OLD_IMAGES"
  
	attribute {
	  name = "TestTableHashKey"
	  type = "S"
	}
  
	replica {
	  region_name = "us-east-2"
	}
  
	replica {
	  region_name = "us-west-2"
	}

	server_side_encryption {
		enabled     = true
		kms_key_arn = aws_kms_key.dynamo_db_kms.key_id
	}
  }
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSDynamoDBTableEncryption,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSDynamoDBTableEncryptionDescription,
			Explanation: AWSDynamoDBTableEncryptionExplanation,
			Impact:      AWSDynamoDBTableEncryptionImpact,
			Resolution:  AWSDynamoDBTableEncryptionResolution,
			BadExample:  AWSDynamoDBTableEncryptionBadExample,
			GoodExample: AWSDynamoDBTableEncryptionGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table#server_side_encryption",
				"https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/EncryptionAtRest.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_dynamodb_table"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			if block.MissingChild("server_side_encryption") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' is not using KMS CMK for encryption", block.FullName())).
						WithRange(block.Range()).
						WithSeverity(severity.Info),
				)
			}

			sseBlock := block.GetBlock("server_side_encryption")
			enabledAttr := sseBlock.GetAttribute("enabled")
			if enabledAttr.IsFalse() {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has server side encryption configured but disabled", block.FullName())).
						WithRange(enabledAttr.Range()).
						WithAttributeAnnotation(enabledAttr).
						WithSeverity(severity.Info),
				)
			}

			if sseBlock.HasChild("kms_key_arn") {
				keyIdAttr := sseBlock.GetAttribute("kms_key_arn")
				if keyIdAttr.Equals("alias/aws/dynamodb") {
					set.Add(
						result.New().
							WithDescription(fmt.Sprintf("Resource '%s' has KMS encryption configured but is using the default aws key", block.FullName())).
							WithRange(keyIdAttr.Range()).
							WithAttributeAnnotation(keyIdAttr).
							WithSeverity(severity.Info),
					)
				}
			}

		},
	})
}
