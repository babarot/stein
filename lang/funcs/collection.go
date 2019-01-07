package funcs

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
)

// LookupListFunc contructs a function that performs dynamic lookups of map types.
var LookupListFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "inputMap",
			Type: cty.DynamicPseudoType,
		},
		{
			Name: "key",
			Type: cty.String,
		},
	},
	Type: func(args []cty.Value) (ret cty.Type, err error) {
		mapVar := args[0]
		lookupKey := args[1].AsString()
		if !mapVar.IsWhollyKnown() {
			return cty.NilType, nil
		}
		m := mapVar.AsValueMap()
		val, ok := m[lookupKey]
		if !ok {
			return cty.List(cty.String), nil
		}
		i := 0
		types := make([]cty.Type, val.LengthInt())
		for it := val.ElementIterator(); it.Next(); {
			_, av := it.Element()
			types[i] = av.Type()
			i++
		}
		retType, _ := convert.UnifyUnsafe(types)
		if retType == cty.NilType {
			return cty.NilType, fmt.Errorf("all arguments must have the same type")
		}

		return cty.List(retType), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		mapVar := args[0]
		lookupKey := args[1].AsString()
		if !mapVar.IsWhollyKnown() {
			return cty.UnknownVal(retType), nil
		}
		// mapVar.Type().IsMapType()
		m := mapVar.AsValueMap()
		val, ok := m[lookupKey]
		if !ok {
			return cty.ListValEmpty(cty.String), nil
		}
		list := make([]cty.Value, 0, val.LengthInt())
		for it := val.ElementIterator(); it.Next(); {
			_, av := it.Element()
			av, _ = convert.Convert(av, retType.ElementType())
			list = append(list, av)
		}
		return cty.ListVal(list), nil
	},
})
