package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/jsonpath"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// JSONPathFunc is
func JSONPathFunc(file string) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "query",
				Type:             cty.String,
				AllowDynamicType: true,
			},
		},
		// Type: function.StaticReturnType(cty.String),
		Type: function.StaticReturnType(cty.DynamicPseudoType),
		Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
			query := "{" + args[0].AsString() + "}"

			filecontent, err := ioutil.ReadFile(file)
			if err != nil {
				return cty.NilVal, err
			}

			// []byte(filecontent)
			decode := scheme.Codecs.UniversalDeserializer().Decode
			obj, _, _ := decode(filecontent, nil, nil)
			// o := obj.(type)

			resourceJSON, err := json.Marshal(obj)
			if err != nil {
				return cty.NilVal, err
			}
			j := jsonpath.New("test")
			j.AllowMissingKeys(true)
			err = j.Parse(query)
			if err != nil {
				return cty.NilVal, err
			}
			buf := new(bytes.Buffer)
			var data interface{}
			err = json.Unmarshal(resourceJSON, &data)
			if err != nil {
				return cty.NilVal, err
			}
			// // pp.Println(data)
			// fullResults, _ := j.FindResults(data)
			// for ix := range fullResults {
			// 	results := fullResults[ix]
			// 	for i, r := range results {
			// 		text, err := j.EvalToText(r)
			// 		if err != nil {
			// 		}
			// 		if i != len(results)-1 {
			// 			text = append(text, ' ')
			// 		}
			// 		pp.Println(string(text))
			// 	}
			// }

			err = j.Execute(buf, data)
			if err != nil {
				return cty.NilVal, err
			}
			// res, _ := j.FindResults(data)
			// pp.Println(res)
			return cty.StringVal(buf.String()), nil
		},
	})

}

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
			res, err := getJSON(query, file, data)
			if err != nil {
				return defaultVal.Type(), nil
			}
			ty, err := ctyjson.ImpliedType(res)
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
			res, err := getJSON(query, file, data)
			if err != nil {
				return defaultVal, nil
			}
			// switch res[0] {
			// case '{', '[':
			// 	val, err := ctyjson.Unmarshal(res, retType)
			// 	if err != nil {
			// 		return cty.StringVal(string(res)), nil
			// 	}
			// 	return val, nil
			// }
			// return cty.StringVal(string(res)), nil
			return ctyjson.Unmarshal(res, retType)
		},
	})
}
