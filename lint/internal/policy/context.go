package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/b4b4r07/stein/lint/internal/policy/funcs"
	"github.com/b4b4r07/stein/lint/internal/policy/terraform"
	"github.com/hashicorp/hcl2/ext/userfunc"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
)

// BuildContext is
func (p Policy) BuildContext(body hcl.Body, filename string, filedata []byte) (*hcl.EvalContext, hcl.Diagnostics) {
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"filename": cty.StringVal(filename), // alias of path.filename
			"path": cty.ObjectVal(map[string]cty.Value{
				"file": cty.StringVal(filename),
				"dir":  cty.StringVal(filepath.Base(filename)),
			}),
		},
		Functions: map[string]function.Function{
			// jsonpath
			"jsonpath": funcs.GJSONFunc(filename, filedata),
			// filepath
			"glob":        funcs.GlobFunc,
			"pathshorten": funcs.PathShortenFunc,
			"ext":         funcs.ExtFunc,
			"exist":       funcs.ExistFunc,
			// basic
			"match": funcs.MatchFunc,
			"color": funcs.ColorFunc,
			"hoge":  funcs.HogeFunc,
			// assertion
			// "equal": funcs.EqualFunc, // Disable
			// "type":  funcs.TypeFunc,  // Disable
			// unix
			"grep": funcs.GrepFunc,
			"wc":   funcs.WcFunc,
			// collection
			"lookuplist": funcs.LookupListFunc,
		},
	}

	functions, body, diags := userfunc.DecodeUserFunctions(body, "function", func() *hcl.EvalContext {
		return ctx
	})

	wantType := cty.DynamicPseudoType

	variableMap := map[string]cty.Value{}
	for _, variable := range p.Variables {
		val, err := convert.Convert(variable.Default, wantType)
		if err != nil {
			// We should never get here because this problem should've been caught
			// during earlier validation, but we'll do something reasonable anyway.
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  `Incorrect variable type`,
				Detail:   fmt.Sprintf(`The resolved value of variable %q is not appropriate: %s.`, "", err),
				Subject:  &variable.DeclRange,
			})
			// Stub out our return value so that the semantic checker doesn't
			// produce redundant downstream errors.
			val = cty.UnknownVal(wantType)
		}
		variableMap[variable.Name] = val
	}
	ctx.Variables["var"] = cty.ObjectVal(variableMap)

	envs := make(map[string]cty.Value)
	for _, env := range os.Environ() {
		key := strings.Split(env, "=")[0]
		val, _ := os.LookupEnv(key)
		envs[key] = cty.StringVal(val)
	}
	ctx.Variables["env"] = cty.ObjectVal(envs)

	for name, f := range functions {
		ctx.Functions[name] = f
	}

	// TODO
	for name, f := range terraform.Functions(os.Getenv("PWD")) {
		ctx.Functions[name] = f
	}

	// expandFuncs := map[string]function.Function{
	// 		"maphoge":    funcs.MapHogeFunc,
	// 	}
	// 	for name, f := range expandFuncs{
	// 	ctx.Functions[name] = f
	// 	}
	ctx.Functions["maphoge"] = funcs.MapHogeFunc(ctx)
	return ctx, diags
}
