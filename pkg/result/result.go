package result

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"

	"github.com/tfsec/tfsec/internal/app/tfsec/block"
	"github.com/tfsec/tfsec/pkg/provider"
	"github.com/tfsec/tfsec/pkg/severity"
)

// Result is a positive result for a security check. It encapsulates a code unique to the specific check it was raised
// by, a human-readable description and a range
type Result struct {
	RuleID          string            `json:"rule_id"`
	RuleSummary     string            `json:"rule_description"`
	RuleProvider    provider.Provider `json:"rule_provider"`
	Impact          string            `json:"impact"`
	Resolution      string            `json:"resolution"`
	Links           []string          `json:"links"`
	Range           block.Range       `json:"location"`
	Description     string            `json:"description"`
	RangeAnnotation string            `json:"-"`
	Severity        severity.Severity `json:"severity"`
	Status          Status            `json:"status"`
}

type Status string

const (
	Failed  Status = "failed"
	Passed  Status = "passed"
	Ignored Status = "ignored"
)

func New() *Result {
	return &Result{
		Status: Failed,
	}
}

func (r *Result) Passed() bool {
	return r.Status == Passed
}

func (r *Result) HashCode() string {
	return fmt.Sprintf("%s:%s", r.Range, r.RuleID)
}

func (r *Result) WithRuleID(id string) *Result {
	r.RuleID = id
	return r
}

func (r *Result) WithRuleSummary(description string) *Result {
	r.RuleSummary = description
	return r
}

func (r *Result) WithRuleProvider(provider provider.Provider) *Result {
	r.RuleProvider = provider
	return r
}

func (r *Result) WithImpact(impact string) *Result {
	r.Impact = impact
	return r
}

func (r *Result) WithResolution(resolution string) *Result {
	r.Resolution = resolution
	return r
}

func (r *Result) WithLink(link string) *Result {
	r.Links = append(r.Links, link)
	return r
}

func (r *Result) WithLinks(links []string) *Result {
	r.Links = links
	return r
}

func (r *Result) WithRange(codeRange block.Range) *Result {
	r.Range = codeRange
	return r
}

func (r *Result) WithDescription(description string) *Result {
	r.Description = description
	return r
}

func (r *Result) WithSeverity(sev severity.Severity) *Result {
	r.Severity = sev
	return r
}

func (r *Result) WithStatus(status Status) *Result {
	r.Status = status
	return r
}

func (r *Result) WithAttributeAnnotation(attr *block.Attribute) *Result {

	var raw interface{}

	var typeStr string

	switch attr.Type() {
	case cty.String:
		raw = attr.Value().AsString()
		typeStr = "string"
	case cty.Bool:
		raw = attr.Value().True()
		typeStr = "bool"
	case cty.Number:
		raw, _ = attr.Value().AsBigFloat().Float64()
		typeStr = "number"
	default:
		return r
	}

	r.RangeAnnotation = fmt.Sprintf("[%s] %#v", typeStr, raw)
	return r
}
