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

const AWSEC2InstanceSensitiveUserdata = "AWS062"
const AWSEC2InstanceSensitiveUserdataDescription = "User data for EC2 instances must not contain sensitive AWS keys"
const AWSEC2InstanceSensitiveUserdataImpact = "User data is visible through the AWS Management console"
const AWSEC2InstanceSensitiveUserdataResolution = "Remove sensitive data from the EC2 instance user-data"
const AWSEC2InstanceSensitiveUserdataExplanation = `
EC2 instance data is used to pass start up information into the EC2 instance. This userdata must not contain access key credentials. Instead use an IAM Instance Profile assigned to the instance to grant access to other AWS Services.
`
const AWSEC2InstanceSensitiveUserdataBadExample = `
resource "aws_instance" "bad_example" {

  ami           = "ami-12345667"
  instance_type = "t2.small"

  user_data = <<EOF
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=us-west-2 
EOF
}
`
const AWSEC2InstanceSensitiveUserdataGoodExample = `
resource "aws_iam_instance_profile" "good_example" {
    // ...
}

resource "aws_instance" "good_example" {
  ami           = "ami-12345667"
  instance_type = "t2.small"

  iam_instance_profile = aws_iam_instance_profile.good_profile.arn

  user_data = <<EOF
  export GREETING=hello
EOF
}
`

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: AWSEC2InstanceSensitiveUserdata,
		Documentation: rule.RuleDocumentation{
			Summary:     AWSEC2InstanceSensitiveUserdataDescription,
			Impact:      AWSEC2InstanceSensitiveUserdataImpact,
			Resolution:  AWSEC2InstanceSensitiveUserdataResolution,
			Explanation: AWSEC2InstanceSensitiveUserdataExplanation,
			BadExample:  AWSEC2InstanceSensitiveUserdataBadExample,
			GoodExample: AWSEC2InstanceSensitiveUserdataGoodExample,
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance#user_data",
				"https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html",
			},
		},
		Provider:       provider.AWSProvider,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_instance"},
		CheckFunc: func(set result.Set, resourceBlock *block.Block, _ *hclcontext.Context) {

			if resourceBlock.MissingChild("user_data") {
			}

			userDataAttr := resourceBlock.GetAttribute("user_data")
			if userDataAttr.Contains("AWS_ACCESS_KEY_ID", block.IgnoreCase) &&
				userDataAttr.RegexMatches("(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has userdata with access key id defined.", resourceBlock.FullName())).
						WithRange(userDataAttr.Range()).
						WithAttributeAnnotation(userDataAttr).
						WithSeverity(severity.Error),
				)
			}

			if userDataAttr.Contains("AWS_SECRET_ACCESS_KEY", block.IgnoreCase) &&
				userDataAttr.RegexMatches("(?i)aws_secre.+[=:]\\s{0,}[A-Za-z0-9\\/+=]{40}.?") {
				set.Add(
					result.New().
						WithDescription(fmt.Sprintf("Resource '%s' has userdata with access secret key defined.", resourceBlock.FullName())).
						WithRange(userDataAttr.Range()).
						WithAttributeAnnotation(userDataAttr).
						WithSeverity(severity.Error),
				)
			}
		},
	})
}
