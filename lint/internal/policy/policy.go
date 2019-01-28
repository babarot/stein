package policy

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

// Policy is
type Policy struct {
	Config    *Config
	Rules     []*Rule
	Variables []*Variable
	Outputs   []*Output
}

// policySchema is the schema for the top-level of a config file. We use
// the low-level HCL API for this level so we can easily deal with each
// block type separately with its own decoding logic.
var policySchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "locals",
		},
		// lint
		{
			Type:       "rule",
			LabelNames: []string{"name"},
		},
		{
			Type: "config",
		},
		{
			Type:       "function",
			LabelNames: []string{"name"},
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
		{
			Type:       "debug",
			LabelNames: []string{"name"},
		},
	},
}

// Decode is
func Decode(body hcl.Body) (*Policy, hcl.Diagnostics) {
	policy := &Policy{}
	// Files: files,
	content, diags := body.Content(policySchema)

	for _, block := range content.Blocks {
		switch block.Type {

		case "variable":
			cfg, cfgDiags := decodeVariableBlock(block, false)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				policy.Variables = append(policy.Variables, cfg)
			}

		case "rule":
			rule, ruleDiags := decodeRuleBlock(block)
			diags = append(diags, ruleDiags...)
			if rule != nil {
				policy.Rules = append(policy.Rules, rule)
			}

		case "output":
			output, outputDiags := decodeOutputBlock(block, false)
			diags = append(diags, outputDiags...)
			if output != nil {
				policy.Outputs = append(policy.Outputs, output)
			}

		case "config":
			config, configDiags := decodeConfigBlock(block)
			diags = append(diags, configDiags...)
			if config != nil {
				policy.Config = config
			}

		case "function":

		}
	}

	diags = append(diags, checkRulesUnique(policy.Rules)...)
	diags = append(diags, checkRulesDependencies(policy.Rules)...)

	return policy, diags
}

func checkRulesUnique(rules []*Rule) hcl.Diagnostics {
	encountered := map[string]*Rule{}
	var diags hcl.Diagnostics
	for _, rule := range rules {
		if existing, exist := encountered[rule.Name]; exist {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate rule definition",
				Detail:   fmt.Sprintf("A rule named %q was already defined at %s. Rule names must be unique within a policy.", existing.Name, existing.DeclRange),
				Subject:  &rule.DeclRange,
			})
		}
		encountered[rule.Name] = rule
	}
	return diags
}

func checkRulesDependencies(rules []*Rule) hcl.Diagnostics {
	var diags hcl.Diagnostics
	for _, rule := range rules {
		for _, dep := range rule.Dependencies {
			exists := func(rules []*Rule) bool {
				for _, rule := range rules {
					if rule.Name == dep {
						return true
					}
				}
				return false
			}(rules)
			if !exists {
				// TODO: Replace more suitable range with rule.DeclRange
				//   "rule.DeclRange" is the declaration range of this rule
				//   however, what we want here is the declaration range of "depends_on"
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid dependency rule",
					Detail:   fmt.Sprintf("A dependency rule %q specified in %q is not defined", dep, rule.Name),
					Subject:  &rule.DeclRange,
				})
			}
		}
	}
	return diags
}
