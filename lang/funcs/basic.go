package funcs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// https://godoc.org/github.com/apparentlymart/go-cty/cty/function/stdlib

// HogeFunc is
var HogeFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "str",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		return cty.StringVal(args[0].AsString() + " hoge"), nil
	},
})

// MapHogeFunc is
func MapHogeFunc(ctx *hcl.EvalContext) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "fn",
				Type: cty.String,
			},
			// {
			// 	Name: "params",
			// 	Type: cty.List(cty.String),
			// },
		},
		VarParam: &function.Parameter{
			Name: "params",
			Type: cty.String,
		},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
			fn := args[0].AsString()
			// TODO
			return ctx.Functions[fn].Call(args[1:])
		},
	})
}

// MatchFunc is
var MatchFunc = function.New(&function.Spec{
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
	Type: function.StaticReturnType(cty.Bool),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		pattern := args[0].AsString()
		text := args[1].AsString()

		matched, err := regexp.MatchString(pattern, text)
		if err != nil {
			return cty.NilVal, err
		}

		return cty.BoolVal(matched), nil
	},
})

// ColorFunc is
var ColorFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "text",
			Type: cty.String,
		},
	},
	VarParam: &function.Parameter{
		Name: "attrs",
		Type: cty.String,
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		text := args[0].AsString()
		var attrs []color.Attribute
		for _, arg := range args[1:] {
			attr := strings.ToLower(arg.AsString())
			switch attr {
			// Fg
			case "black", "FgBlack", "fgblack":
				attrs = append(attrs, color.FgBlack)
			case "red", "FgRed", "fgred":
				attrs = append(attrs, color.FgRed)
			case "green", "FgGreen", "fggreen":
				attrs = append(attrs, color.FgGreen)
			case "yellow", "FgYellow", "fgyellow":
				attrs = append(attrs, color.FgYellow)
			case "blue", "FgBlue", "fgblue":
				attrs = append(attrs, color.FgBlue)
			case "magenta", "FgMagenta", "fgmagenta":
				attrs = append(attrs, color.FgMagenta)
			case "cyan", "FgCyan", "fgcyan":
				attrs = append(attrs, color.FgCyan)
			case "white", "FgWhite", "fgwhite":
				attrs = append(attrs, color.FgWhite)
			// Bg
			case "BgBlack", "bgblack":
				attrs = append(attrs, color.BgBlack)
			case "BgRed", "bgred":
				attrs = append(attrs, color.BgRed)
			case "BgGreen", "bggreen":
				attrs = append(attrs, color.BgGreen)
			case "BgYellow", "bgyellow":
				attrs = append(attrs, color.BgYellow)
			case "BgBlue", "bgblue":
				attrs = append(attrs, color.BgBlue)
			case "BgMagenta", "bgmagenta":
				attrs = append(attrs, color.BgMagenta)
			case "BgCyan", "bgcyan":
				attrs = append(attrs, color.BgCyan)
			case "BgWhite", "bgwhite":
				attrs = append(attrs, color.BgWhite)
			// FgHi
			case "FgHiBlack", "fghiblack":
				attrs = append(attrs, color.FgHiBlack)
			case "FgHiRed", "fghired":
				attrs = append(attrs, color.FgHiRed)
			case "FgHiGreen", "fghigreen":
				attrs = append(attrs, color.FgHiGreen)
			case "FgHiYellow", "fghiyellow":
				attrs = append(attrs, color.FgHiYellow)
			case "FgHiBlue", "fghiblue":
				attrs = append(attrs, color.FgHiBlue)
			case "FgHiMagenta", "fghimagenta":
				attrs = append(attrs, color.FgHiMagenta)
			case "FgHiCyan", "fghicyan":
				attrs = append(attrs, color.FgHiCyan)
			case "FgHiWhite", "fghiwhite":
				attrs = append(attrs, color.FgHiWhite)
			// Attr
			case "Reset", "reset":
				attrs = append(attrs, color.Reset)
			case "Bold", "bold":
				attrs = append(attrs, color.Bold)
			case "Faint", "faint":
				attrs = append(attrs, color.Faint)
			case "Italic", "italic":
				attrs = append(attrs, color.Italic)
			case "Underline", "underline":
				attrs = append(attrs, color.Underline)
			case "BlinkSlow", "blinkslow":
				attrs = append(attrs, color.BlinkSlow)
			case "BlinkRapid", "blinkrapid":
				attrs = append(attrs, color.BlinkRapid)
			case "ReverseVideo", "reversevideo":
				attrs = append(attrs, color.ReverseVideo)
			case "Concealed", "concealed":
				attrs = append(attrs, color.Concealed)
			case "CrossedOut", "crossedout":
				attrs = append(attrs, color.CrossedOut)
			default:
				return cty.NilVal, fmt.Errorf("%q: invalid attr name", attr)
			}
		}
		c := color.New(attrs...)
		return cty.StringVal(c.Sprintf(text)), nil
	},
})
