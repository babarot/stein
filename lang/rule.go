package lang

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func decodeRuleBlock(block *hcl.Block) (*Rule, hcl.Diagnostics) {
	rule := &Rule{
		Name: block.Labels[0],
		// Name:      block.Labels[1],
		DeclRange: block.DefRange,
		TypeRange: block.LabelRanges[0],
	}
	content, remain, diags := block.Body.PartialContent(ruleBlockSchema)
	rule.Config = remain

	// TODO
	if !hclsyntax.ValidIdentifier(rule.Name) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid output name",
			Detail:   badIdentifierDetail,
			Subject:  &block.LabelRanges[0],
		})
	}

	for _, block := range content.Blocks {
		switch block.Type {
		case "locals":
			// locals, localsDiags := decodeLocalsBlock(block)
			// diags = append(diags, localsDiags...)
			// if locals != nil {
			// 	rule.Locals = locals
			// }
			// pp.Println(locals)
			// rule.Locals = append(rule.Locals, locals...)
		case "report":
			report, reportDiags := decodeReportBlock(block)
			diags = append(diags, reportDiags...)
			if report != nil {
				rule.Report = report
			}
		default:
			continue
		}
	}

	return rule, diags
}

// Rule isj
type Rule struct {
	Name string

	Config hcl.Body

	TypeRange hcl.Range
	DeclRange hcl.Range

	Description string
	IgnoreCases []bool
	Expressions []bool
	Report      *Report

	Debug []string
}

// // Locals is
// type Locals struct {
// 	Config hcl.Body
//
// 	TypeRange hcl.Range
// 	DeclRange hcl.Range
// }

// Report is
type Report struct {
	Config hcl.Body

	TypeRange hcl.Range
	DeclRange hcl.Range

	Level   string
	Message string
}

var ruleBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "locals",
			// LabelNames: []string{"name"},
		},
		{
			Type: "report",
			// LabelNames: []string{"name"},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "description",
			Required: true,
		},
		{
			Name:     "expressions",
			Required: true,
		},
		{
			Name: "ignore_cases",
		},
	},
}

func decodeReportBlock(block *hcl.Block) (*Report, hcl.Diagnostics) {
	content, config, diags := block.Body.PartialContent(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "level",
				Required: true,
			},
			{
				Name:     "message",
				Required: true,
			},
		},
	})

	report := &Report{
		Config:    config,
		DeclRange: block.DefRange,
	}

	if attr, exists := content.Attributes["level"]; exists {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)
		if diags.HasErrors() {
			return report, diags
		}
		switch val.Type() {
		case cty.String:
			// ok
		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Type string is required here",
				Detail:   fmt.Sprintf("It can take %q as the level", []string{"WARN", "ERROR"}),
				Subject:  &attr.NameRange,
			})
			return report, diags
		}
		level := val.AsString()
		switch level {
		case "ERROR":
		case "WARN":
		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid level type",
				Detail:   fmt.Sprintf("got %q but want %q", level, []string{"WARN", "ERROR"}),
				Subject:  &attr.NameRange,
			})
		}
	} else {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing required argument",
			Detail:   "The argument \"level\" is required, but no definition was found.",
			Subject:  &block.DefRange,
		})
	}

	if _, exists := content.Attributes["message"]; !exists {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing required argument",
			Detail:   "The argument \"message\" is required, but no definition was found.",
			Subject:  &block.DefRange,
		})
	}

	return report, diags
}
