package rules

import (
	"fmt"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/pkg/provider"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/tfsec/tfsec/internal/app/tfsec/security"

	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"

	"github.com/zclconf/go-cty/cty"
)

// GenericSensitiveAttributes See https://github.com/tfsec/tfsec#included-checks for check info
const GenericSensitiveAttributes = "GEN003"
const GenericSensitiveAttributesDescription = "Potentially sensitive data stored in block attribute."
const GenericSensitiveAttributesImpact = "Block attribute could be leaking secrets"
const GenericSensitiveAttributesResolution = "Don't include sensitive data in blocks"
const GenericSensitiveAttributesExplanation = `
Sensitive attributes such as passwords and API tokens should not be available in your templates, especially in a plaintext form. You can declare variables to hold the secrets, assuming you can provide values for those variables in a secure fashion. Alternatively, you can store these secrets in a secure secret store, such as AWS KMS.

*NOTE: It is also recommended to store your Terraform state in an encrypted form.*
`
const GenericSensitiveAttributesBadExample = `
resource "evil_corp" "bad_example" {
	root_password = "p4ssw0rd"
}
`
const GenericSensitiveAttributesGoodExample = `
variable "password" {
  description = "The root password for our VM"
  type        = string
}

resource "evil_corp" "good_example" {
	root_password = var.password
}
`

var sensitiveWhitelist = []struct {
	Resource  string
	Attribute string
}{
	{
		Resource:  "aws_efs_file_system",
		Attribute: "creation_token",
	},
	{
		Resource:  "aws_instance",
		Attribute: "get_password_data",
	},
	{
		Resource:  "github_actions_secret",
		Attribute: "secret_name",
	},
	{
		Resource:  "github_actions_organization_secret",
		Attribute: "secret_name",
	},
	{
		Resource:  "google_secret_manager_secret",
		Attribute: "secret_id",
	},
}

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		ID: GenericSensitiveAttributes,
		Documentation: rule.RuleDocumentation{
			Summary:     GenericSensitiveAttributesDescription,
			Impact:      GenericSensitiveAttributesImpact,
			Resolution:  GenericSensitiveAttributesResolution,
			Explanation: GenericSensitiveAttributesExplanation,
			BadExample:  GenericSensitiveAttributesBadExample,
			GoodExample: GenericSensitiveAttributesGoodExample,
			Links: []string{
				"https://www.terraform.io/docs/state/sensitive-data.html",
			},
		},
		Provider:      provider.GeneralProvider,
		RequiredTypes: []string{"resource", "provider", "module"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {

			attributes := block.GetAttributes()

		SKIP:
			for _, attribute := range attributes {
				for _, whitelisted := range sensitiveWhitelist {
					if whitelisted.Resource == block.TypeLabel() && whitelisted.Attribute == attribute.Name() {
						continue SKIP
					}
				}
				if security.IsSensitiveAttribute(attribute.Name()) {
					if attribute.Type() == cty.String && attribute.Value().AsString() != "" {
						set.Add(result.New().
							WithDescription(fmt.Sprintf("Block '%s' includes a potentially sensitive attribute which is defined within the project.", block.FullName())).
							WithRange(attribute.Range()).
							WithAttributeAnnotation(attribute).
							WithSeverity(severity.Warning),
						)
					}

				}
			}

		},
	})
}
