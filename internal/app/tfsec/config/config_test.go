package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfsec/tfsec/internal/app/tfsec/config"
)

func TestExcludesElementsFromYAML(t *testing.T) {
	content := `
severity_overrides:
  AWS018: INFO

exclude:
  - DP001
`
	c := load(t, "config.yaml", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	assert.Contains(t, c.ExcludedChecks, "DP001")
}

func TestExcludesElementsFromYML(t *testing.T) {
	content := `
severity_overrides:
  AWS018: INFO

exclude:
  - DP001
`
	c := load(t, "config.yml", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	assert.Contains(t, c.ExcludedChecks, "DP001")
}

func TestExcludesElementsFromJSON(t *testing.T) {
	content := `{
  "severity_overrides": {
    "AWS018": "INFO"
  },
  "exclude": [
    "DP001"
  ]
}
`
	c := load(t, "config.json", content)

	assert.Contains(t, c.SeverityOverrides, "AWS018")
	assert.Contains(t, c.ExcludedChecks, "DP001")
}

func load(t *testing.T, filename, content string) *config.Config {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	configFileName := fmt.Sprintf("%s/%s", dir, filename)

	err = ioutil.WriteFile(configFileName, []byte(content), os.ModePerm)
	require.NoError(t, err)

	c, err := config.LoadConfig(configFileName)
	require.NoError(t, err)

	return c
}
