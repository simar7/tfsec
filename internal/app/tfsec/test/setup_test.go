package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tfsec/tfsec/pkg/result"
	"github.com/tfsec/tfsec/pkg/severity"

	"github.com/tfsec/tfsec/internal/app/tfsec/hclcontext"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"

	"github.com/tfsec/tfsec/pkg/rule"

	"github.com/stretchr/testify/assert"

	"github.com/tfsec/tfsec/internal/app/tfsec/parser"
	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

const exampleCheckCode = "EXA001"

var excludedChecksList []string

func TestMain(t *testing.M) {

	scanner.RegisterCheckRule(rule.Rule{
		ID: exampleCheckCode,
		Documentation: rule.RuleDocumentation{
			Summary:     "A stupid example check for a test.",
			Impact:      "You will look stupid",
			Resolution:  "Don't do stupid stuff",
			Explanation: "Bad should not be set.",
			BadExample: `
resource "problem" "x" {
	bad = "1"
}
`,
			GoodExample: `
resource "problem" "x" {
	
}
`,
			Links: nil,
		},
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"problem"},
		CheckFunc: func(set result.Set, block *block.Block, _ *hclcontext.Context) {
			if block.GetAttribute("bad") != nil {
				set.Add(
					result.New().WithDescription("example problem").WithRange(block.Range()).WithSeverity(severity.Error),
				)
			}
		},
	})

	os.Exit(t.Run())
}

func scanSource(source string) []result.Result {
	blocks := createBlocksFromSource(source)
	return scanner.New(scanner.OptionExcludeRules(excludedChecksList)).Scan(blocks)
}

func createBlocksFromSource(source string) []*block.Block {
	path := createTestFile("test.tf", source)
	blocks, err := parser.New(filepath.Dir(path), parser.OptionStopOnHCLError()).ParseDirectory()
	if err != nil {
		panic(err)
	}
	return blocks
}

func createTestFile(filename, contents string) string {
	dir, err := ioutil.TempDir(os.TempDir(), "tfsec")
	if err != nil {
		panic(err)
	}
	path := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(path, []byte(contents), 0755); err != nil {
		panic(err)
	}
	return path
}

func assertCheckCode(t *testing.T, includeCode string, excludeCode string, results []result.Result) {

	var foundInclude bool
	var foundExclude bool

	var excludeText string

	for _, res := range results {
		if res.RuleID == excludeCode {
			foundExclude = true
			excludeText = res.Description
		}
		if res.RuleID == includeCode {
			foundInclude = true
		}
	}

	assert.False(t, foundExclude, fmt.Sprintf("res with code '%s' was found but should not have been: %s", excludeCode, excludeText))
	if includeCode != "" {
		assert.True(t, foundInclude, fmt.Sprintf("res with code '%s' was not found but should have been", includeCode))
	}
}

func createTestFileWithModule(contents string, moduleContents string) string {
	dir, err := ioutil.TempDir(os.TempDir(), "tfsec")
	if err != nil {
		panic(err)
	}

	rootPath := filepath.Join(dir, "main")
	modulePath := filepath.Join(dir, "module")

	if err := os.Mkdir(rootPath, 0755); err != nil {
		panic(err)
	}

	if err := os.Mkdir(modulePath, 0755); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(rootPath, "main.tf"), []byte(contents), 0755); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(modulePath, "main.tf"), []byte(moduleContents), 0755); err != nil {
		panic(err)
	}

	return rootPath
}
