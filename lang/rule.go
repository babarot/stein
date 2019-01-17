package lang

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// RulePrefix is the prefix of the rule name
const RulePrefix = "rule."

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

	// TODO: Depends on jsonpath's default value
	//   By introducing the default value concept into jsonpath func,
	//   it became not to need to take care of this
	//   if jsonpath doesn't have default value, this case will be failed sometimes
	//   e.g. "${jsonpath("spec.replicas") > 0}"
	//   if kind is Service, jsonpath("spec.replicas") returns ""
	//   so this expression will be "${"" > 0}" as a result
	//   This is the reason why needs default value for jsonpath
	//   Otherwise, operate the LHS which is dependent on RHS in hcl2 library
	//
	//   See also
	//   https://github.com/hashicorp/hcl2/blob/cce5ae6cc5c890122f922573d6bf973eef0509f7/hcl/hclsyntax/expression_ops.go#L123-L196
	//
	// if attr, exists := content.Attributes["expressions"]; exists {
	// }

	if attr, exists := content.Attributes["depends_on"]; exists {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)
		for it := val.ElementIterator(); it.Next(); {
			_, dep := it.Element()
			dependency := strings.TrimPrefix(dep.AsString(), RulePrefix)
			rule.Dependencies = append(rule.Dependencies, dependency)
		}
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

	Description  string
	Dependencies []string
	IgnoreCases  []bool
	Expressions  []bool
	Report       *Report

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
		},
		{
			Type: "report",
		},
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "description",
			Required: true,
		},
		{
			Name: "depends_on",
		},
		{
			Name: "ignore_cases",
		},
		{
			Name:     "expressions",
			Required: true,
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
