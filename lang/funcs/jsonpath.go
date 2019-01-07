package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"sigs.k8s.io/yaml"

	"github.com/hashicorp/hcl"
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

func getJSON(file, query string) ([]byte, error) {
	var data []byte
	// var err error
	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return data, nil
		}
		data, err = yaml.YAMLToJSON(contents)
		if err != nil {
			return data, nil
		}
	case ".json":
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return data, nil
		}
		data = contents
	case ".hcl", ".tf":
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return data, nil
		}
		var v interface{}
		err = hcl.Unmarshal(contents, &v)
		if err != nil {
			return data, fmt.Errorf("unable to parse HCL: %s", err)
		}
		jd, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return data, fmt.Errorf("unable to marshal json: %s", err)
		}
		data = jd
	default:
		return data, fmt.Errorf("invalid file format: %q (%s)", file, filepath.Ext(file))
	}
	result := gjson.GetBytes(data, query)
	if !result.Exists() {
		// return []byte{}, fmt.Errorf("%q: not found in %q", query, file)
		return []byte(""), nil
	}
	return []byte(result.String()), nil

}

// GJSONFunc is
func GJSONFunc(file string) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "str",
				Type: cty.String,
			},
		},
		Type: func(args []cty.Value) (cty.Type, error) {
			query := args[0].AsString()
			res, err := getJSON(file, query)
			if err != nil {
				return cty.NilType, err
			}
			ty, err := ctyjson.ImpliedType(res)
			if err != nil {
				// When the result from getJSON can not be converted to JSON (that is, array or map),
				// treat the return value as a string
				return cty.String, nil
			}
			return ty, nil
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			query := args[0].AsString()
			res, err := getJSON(file, query)
			if err != nil {
				return cty.NilVal, err
			}
			val, err := ctyjson.Unmarshal(res, retType)
			if err != nil {
				return cty.StringVal(string(res)), nil
			}
			return val, nil
		},
	})
}
