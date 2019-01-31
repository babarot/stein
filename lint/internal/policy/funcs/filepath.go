package funcs

import (
	"os"
	"path/filepath"

	pathshorten "github.com/b4b4r07/go-pathshorten"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// GlobFunc returns a list of files matching a given pattern
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

// PathShortenFunc returns the shorten path of given path
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

// ExtFunc returns an extension of given file
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

// ExistFunc returns true if given path exists
var ExistFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "path",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.Bool),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		path := args[0].AsString()
		_, err = os.Stat(path)
		exist := !os.IsNotExist(err)
		return cty.BoolVal(exist), nil
	},
})
