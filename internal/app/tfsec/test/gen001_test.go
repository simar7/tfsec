package test

import (
	"testing"

	"github.com/tfsec/tfsec/internal/app/tfsec/rules"
)

func Test_AWSSensitiveVariables(t *testing.T) {

	var tests = []struct {
		name                  string
		source                string
		mustIncludeResultCode string
		mustExcludeResultCode string
	}{
		{
			name: "check sensitive variable with value",
			source: `
variable "db_password" {
	default = "something"
}`,
			mustIncludeResultCode: rules.GenericSensitiveVariables,
		},
		{
			name: "check sensitive variable without value",
			source: `
variable "db_password" {
	default = ""
}`,
			mustExcludeResultCode: rules.GenericSensitiveVariables,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results := scanSource(test.source)
			assertCheckCode(t, test.mustIncludeResultCode, test.mustExcludeResultCode, results)
		})
	}

}
