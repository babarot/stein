package funcs

import (
	"path/filepath"

	pathshorten "github.com/b4b4r07/go-pathshorten"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// GlobFunc is
var GlobFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "pattern",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.List(cty.String)),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		pattern := args[0].AsString()

		files, err := filepath.Glob(pattern)
		if err != nil {
			return cty.NilVal, err
		}

		vals := make([]cty.Value, len(files))
		for i, file := range files {
			vals[i] = cty.StringVal(file)
		}

		if len(vals) == 0 {
			return cty.ListValEmpty(cty.String), nil
		}
		return cty.ListVal(vals), nil
	},
})

// PathShortenFunc is
var PathShortenFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "path",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		path := args[0].AsString()
		return cty.StringVal(pathshorten.Run(path)), nil
	},
})

// ExtFunc is
var ExtFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "file",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		return cty.StringVal(filepath.Ext(args[0].AsString())), nil
	},
})
