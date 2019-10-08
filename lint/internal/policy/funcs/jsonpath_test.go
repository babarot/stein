package funcs

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func toJSON(json []byte) []byte {
	yaml, err := yaml.ToJSON(json)
	if err != nil {
		return json
	}
	return yaml
}

func TestGJSONFunc(t *testing.T) {
	tests := []struct {
		Filename string
		Contents []byte
		Query    cty.Value
		Want     cty.Value
	}{
		{
			"simple.json",
			[]byte(`{"key": "value"}`),
			cty.StringVal("key"),
			cty.StringVal("value"),
		},
		{
			"simple_kubernetes.yaml",
			toJSON([]byte(`---
apiVersion: v1
kind: Test
metadata:
  name: test
  namespace: testns
`)),
			cty.StringVal("metadata.name"),
			cty.StringVal("test"),
		},
		{
			"string_value.yaml",
			toJSON([]byte(`key: value`)),
			cty.StringVal("key"),
			cty.StringVal("value"),
		},
		{
			"number_value.yaml",
			toJSON([]byte(`port: 8080`)),
			cty.StringVal("port"),
			cty.NumberIntVal(8080),
		},
		{
			"object_value.yaml",
			toJSON([]byte(`---
apiVersion: v1
kind: Test
metadata:
  name: test
  namespace: testns
`)),
			cty.StringVal("metadata"),
			cty.ObjectVal(map[string]cty.Value{"name": cty.StringVal("test"), "namespace": cty.StringVal("testns")}),
		},
		{
			"tuple_value.yaml",
			toJSON([]byte(`---
test:
  - a
  - b
`)),
			cty.StringVal("test"),
			cty.TupleVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
		},
		{
			"strict_string.json",
			[]byte(`{"maxSurge": "100%"}`),
			cty.StringVal("maxSurge"),
			cty.StringVal("100%"),
		},
		{
			"strict_boolean.json",
			[]byte(`{"ok": true}`),
			cty.StringVal("ok"),
			cty.BoolVal(true),
		},
		{
			"strict_boolean.json",
			[]byte(`{"ok": "true"}`),
			cty.StringVal("ok"),
			cty.BoolVal(true),
		},
		{
			"strict_boolean.json",
			[]byte(`{"ok": "false"}`),
			cty.StringVal("ok"),
			cty.BoolVal(false),
		},
	}

	for _, test := range tests {
		t.Run(test.Filename, func(t *testing.T) {
			got, err := GJSON(test.Filename, test.Contents, test.Query)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
