package funcs

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// EqualFunc is
var EqualFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "a",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowDynamicType: true,
			AllowNull:        true,
		},
		{
			Name:             "b",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowDynamicType: true,
			AllowNull:        true,
		},
	},
	Type: function.StaticReturnType(cty.Bool),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		a := args[0]
		b := args[1]
		if a.Type() != b.Type() {
			return cty.NilVal, fmt.Errorf("type not same: left is %q but right is %q", a.Type().FriendlyName(), b.Type().FriendlyName())
		}
		return a.Equals(b), nil
	},
})

// TypeFunc is
var TypeFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "arg",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowDynamicType: true,
			AllowNull:        true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		return cty.StringVal(args[0].Type().FriendlyName()), nil
	},
})
