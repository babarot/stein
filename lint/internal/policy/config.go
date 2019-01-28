package policy

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

// Config is
type Config struct {
	Name string

	Config hcl.Body

	TypeRange hcl.Range
	DeclRange hcl.Range

	// can ignore (TODO)
	Report *ReportConfig
}

// ReportConfig is
type ReportConfig struct {
	Name string

	Config hcl.Body

	TypeRange hcl.Range
	DeclRange hcl.Range
}

var configBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "report",
			// LabelNames: []string{"name"},
		},
	},
}

func decodeConfigBlock(block *hcl.Block) (*Config, hcl.Diagnostics) {
	// attrs, diags := block.Body.JustAttributes()
	// if len(attrs) == 0 {
	// 	return nil, diags
	// }
	//
	// // locals := make([]*Local, 0, len(attrs))
	// var config *Config
	// for name, attr := range attrs {
	// 	if !hclsyntax.ValidIdentifier(name) {
	// 		diags = append(diags, &hcl.Diagnostic{
	// 			Severity: hcl.DiagError,
	// 			Summary:  "Invalid config value name",
	// 			Detail:   badIdentifierDetail,
	// 			Subject:  &attr.NameRange,
	// 		})
	// 	}
	//
	// 	config = &Config{
	// 		Name: name,
	// 		// Expr:      attr.Expr,
	// 		DeclRange: attr.Range,
	// 	}
	// }
	//
	// return config, diags

	cfg := &Config{
		DeclRange: block.DefRange,
		// TypeRange: block.LabelRanges[0],
	}
	content, remain, diags := block.Body.PartialContent(configBlockSchema)
	cfg.Config = remain

	for _, block := range content.Blocks {
		switch block.Type {
		case "report":
			report, reportDiags := decodeReportConfigBlock(block)
			diags = append(diags, reportDiags...)
			if report != nil {
				cfg.Report = report
			}
		default:
			continue
		}
	}

	return cfg, diags
}

func decodeReportConfigBlock(block *hcl.Block) (*ReportConfig, hcl.Diagnostics) {
	content, config, diags := block.Body.PartialContent(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "format",
				// Required: true,
			},
			{
				Name: "style",
				// Required: true,
			},
			{
				Name: "color",
				// Required: true,
			},
		},
	})
	rc := &ReportConfig{
		Config:    config,
		DeclRange: block.DefRange,
	}

	if attr, exists := content.Attributes["style"]; exists {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)
		style := val.AsString()
		switch style {
		case "console":
		default:
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid style value",
				Detail:   fmt.Sprintf("got %q but want %q", style, []string{"console"}),
				Subject:  &attr.NameRange,
			})
		}
	}

	return rc, diags
}
