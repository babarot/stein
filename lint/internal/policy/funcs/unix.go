package funcs

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// GrepFunc is
var GrepFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "pattern",
			Type: cty.String,
		},
		{
			Name: "text",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		pattern := args[0].AsString()
		text := args[1].AsString()
		var matches []string
		in := strings.NewReader(text)
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			matched, err := regexp.MatchString(pattern, scanner.Text())
			if err != nil {
				return cty.NilVal, err
			}
			if matched {
				matches = append(matches, scanner.Text())
			}
		}
		return cty.StringVal(strings.Join(matches, "\n")), nil
	},
})

// WcFunc is
var WcFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "text",
			Type: cty.String,
		},
	},
	VarParam: &function.Parameter{
		Name: "opts",
		Type: cty.String,
	},
	Type: function.StaticReturnType(cty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		text := args[0].AsString()
		opts := args[1:]
		var (
			lines = int64(strings.Count(text, "\n"))
			chars = int64(len(text))
			words = int64(len(strings.Fields(text)))
		)
		for _, opt := range opts {
			switch opt.AsString() {
			case "l":
				return cty.NumberIntVal(lines), nil
			case "c":
				return cty.NumberIntVal(chars), nil
			case "w":
				return cty.NumberIntVal(words), nil
			default:
				return cty.NilVal, fmt.Errorf("%v: not supported option", opt.AsString())
			}
		}
		// default option is l
		return cty.NumberIntVal(lines), nil
	},
})

// Grep returns the matched line from given text by using regex
func Grep(pattern cty.Value, text cty.Value) (cty.Value, error) {
	return GrepFunc.Call([]cty.Value{pattern, text})
}

// Wc counts the characters, words, or lines from given text
func Wc(text cty.Value, opts ...cty.Value) (cty.Value, error) {
	args := make([]cty.Value, len(opts)+1)
	args[0] = text
	copy(args[1:], opts)
	return WcFunc.Call(args)
}
