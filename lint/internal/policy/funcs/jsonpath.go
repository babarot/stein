package funcs

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/tidwall/gjson"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func getJSON(query string, file string, data []byte) ([]byte, error) {
	result := gjson.GetBytes(data, query)
	if !result.Exists() {
		// Even if the given query refers to a field that does not exist,
		// it does not return an error
		// This is because there are cases which a rule defined by a user
		// refers to a field that does not exist and tests it
		//
		// return []byte{}, fmt.Errorf("%q: not found in %q", query, file)
		// return []byte(""), nil
		//
		// TODO: Fix this handling
		//   By introducing the default value,
		//   no problem for now even if the result is not found.
		//   This is because returns default values if not found case
		return []byte{}, fmt.Errorf("%q: not found in %q", query, file)
	}
	return []byte(result.String()), nil
}

// GJSON determines whether a file exists at the given path.
//
// The underlying function implementation works relative to a input file
// and its contents, so this wrapper takes a input file string, etc
// and uses it to onstruct the underlying function before calling it.
func GJSON(file string, data []byte, query cty.Value) (cty.Value, error) {
	fn := GJSONFunc(file, data)
	return fn.Call([]cty.Value{query})
}

// GJSONFunc is
func GJSONFunc(file string, data []byte) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "str",
				Type: cty.String,
			},
		},
		VarParam: &function.Parameter{
			Name: "default",
			Type: cty.DynamicPseudoType,
		},
		Type: func(args []cty.Value) (cty.Type, error) {
			query := args[0].AsString()
			defaultVal := cty.StringVal("")
			if len(args) > 1 {
				defaultVal = args[1]
			}
			b, err := getJSON(query, file, data)
			if err != nil {
				return defaultVal.Type(), nil
			}
			ty, err := ctyjson.ImpliedType(b)
			if err != nil {
				// When the result from getJSON can not be converted to JSON (that is, array or map),
				// treat the return value as a string
				return defaultVal.Type(), nil
			}
			return ty, nil
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			query := args[0].AsString()
			defaultVal := cty.StringVal("")
			if len(args) > 1 {
				defaultVal = args[1]
			}
			b, err := getJSON(query, file, data)
			if err != nil {
				return defaultVal, nil
			}
			switch b[0] {
			case '{', '[':
				val, err := ctyjson.Unmarshal(b, retType)
				if err != nil {
					return cty.StringVal(string(b)), nil
				}
				return val, nil
			}
			isValidNumber := func(b byte) bool {
				return '0' <= b && b <= '9'
			}
			var shouldReturnString bool
			for _, char := range b {
				if !isValidNumber(char) {
					shouldReturnString = true
				}
			}
			if shouldReturnString {
				return cty.StringVal(string(b)), nil
			}
			f64, _ := strconv.ParseFloat(string(b), 64)
			val := big.NewFloat(f64)
			return cty.NumberVal(val), nil
		},
	})
}
