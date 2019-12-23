package hclconvert

import (
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	ctyconvert "github.com/zclconf/go-cty/cty/convert"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// This file is attributed to https://github.com/tmccombs/hcl2json.
// convertBlock() is customized to create the same data structure as by previous stein version

type jsonObj map[string]interface{}

// Convert an hcl File to a json serializable object
// This assumes that the body is a hclsyntax.Body
func ConvertFile(file *hcl.File) (jsonObj, error) {
	c := converter{bytes: file.Bytes}
	body := file.Body.(*hclsyntax.Body)
	return c.convertBody(body)
}

type converter struct {
	bytes []byte
}

func (c *converter) rangeSource(r hcl.Range) string {
	return string(c.bytes[r.Start.Byte:r.End.Byte])
}

func (c *converter) convertBody(body *hclsyntax.Body) (jsonObj, error) {
	var err error
	out := make(jsonObj)
	for key, value := range body.Attributes {
		out[key], err = c.convertExpression(value.Expr)
		if err != nil {
			return nil, err
		}
	}

	for _, block := range body.Blocks {
		err = c.convertBlock(block, out)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (c *converter) convertBlock(block *hclsyntax.Block, out jsonObj) error {
	var key string = block.Type

	value, err := c.convertBody(block.Body)
	if err != nil {
		return err
	}

	// This implementation is different from hcl2json's
	root := NewBlockTree(key)
	current := root
	for _, label := range block.Labels {
		inner := NewBlockTree(label)
		current.appendSubTree(inner)
		current = inner
	}
	current.appendValue(value)
	result := root.out()

	if _, ok := out[key]; ok {
		tmp := map[string][]interface{}{}
		tmp[key] = out[key].([]interface{})
		result = mergeSliceMap(tmp, result)
	}
	out[key] = result[key]

	return nil
}

func (c *converter) convertExpression(expr hclsyntax.Expression) (interface{}, error) {
	// assume it is hcl syntax (because, um, it is)
	switch value := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return ctyjson.SimpleJSONValue{Value: value.Val}, nil
	case *hclsyntax.TemplateExpr:
		return c.convertTemplate(value)
	case *hclsyntax.TemplateWrapExpr:
		return c.convertExpression(value.Wrapped)
	case *hclsyntax.TupleConsExpr:
		var list []interface{}
		for _, ex := range value.Exprs {
			elem, err := c.convertExpression(ex)
			if err != nil {
				return nil, err
			}
			list = append(list, elem)
		}
		return list, nil
	case *hclsyntax.ObjectConsExpr:
		m := make(jsonObj)
		for _, item := range value.Items {
			key, err := c.convertKey(item.KeyExpr)
			if err != nil {
				return nil, err
			}
			m[key], err = c.convertExpression(item.ValueExpr)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	default:
		return c.wrapExpr(expr), nil
	}
}

func (c *converter) convertTemplate(t *hclsyntax.TemplateExpr) (string, error) {
	if t.IsStringLiteral() {
		// safe because the value is just the string
		v, err := t.Value(nil)
		if err != nil {
			return "", err
		}
		return v.AsString(), nil
	}
	var builder strings.Builder
	for _, part := range t.Parts {
		s, err := c.convertStringPart(part)
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
	}
	return builder.String(), nil
}

func (c *converter) convertStringPart(expr hclsyntax.Expression) (string, error) {
	switch v := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		s, err := ctyconvert.Convert(v.Val, cty.String)
		if err != nil {
			return "", err
		}
		return s.AsString(), nil
	case *hclsyntax.TemplateExpr:
		return c.convertTemplate(v)
	case *hclsyntax.TemplateWrapExpr:
		return c.convertStringPart(v.Wrapped)
	case *hclsyntax.ConditionalExpr:
		return c.convertTemplateConditional(v)
	case *hclsyntax.TemplateJoinExpr:
		return c.convertTemplateFor(v.Tuple.(*hclsyntax.ForExpr))
	default:
		// treating as an embedded expression
		return c.wrapExpr(expr), nil
	}
}

func (c *converter) convertKey(keyExpr hclsyntax.Expression) (string, error) {
	// a key should never have dynamic input
	if k, isKeyExpr := keyExpr.(*hclsyntax.ObjectConsKeyExpr); isKeyExpr {
		keyExpr = k.Wrapped
		if _, isTraversal := keyExpr.(*hclsyntax.ScopeTraversalExpr); isTraversal {
			return c.rangeSource(keyExpr.Range()), nil
		}
	}
	return c.convertStringPart(keyExpr)
}

func (c *converter) convertTemplateConditional(expr *hclsyntax.ConditionalExpr) (string, error) {
	var builder strings.Builder
	builder.WriteString("%{if ")
	builder.WriteString(c.rangeSource(expr.Condition.Range()))
	builder.WriteString("}")
	trueResult, err := c.convertStringPart(expr.TrueResult)
	if err != nil {
		return "", err
	}
	builder.WriteString(trueResult)
	falseResult, err := c.convertStringPart(expr.FalseResult)
	if err != nil {
		return "", err
	}
	if len(falseResult) > 0 {
		builder.WriteString("%{else}")
		builder.WriteString(falseResult)
	}
	builder.WriteString("%{endif}")

	return builder.String(), nil
}

func (c *converter) convertTemplateFor(expr *hclsyntax.ForExpr) (string, error) {
	var builder strings.Builder
	builder.WriteString("%{for ")
	if len(expr.KeyVar) > 0 {
		builder.WriteString(expr.KeyVar)
		builder.WriteString(", ")
	}
	builder.WriteString(expr.ValVar)
	builder.WriteString(" in ")
	builder.WriteString(c.rangeSource(expr.CollExpr.Range()))
	builder.WriteString("}")
	templ, err := c.convertStringPart(expr.ValExpr)
	if err != nil {
		return "", err
	}
	builder.WriteString(templ)
	builder.WriteString("%{endfor}")

	return builder.String(), nil
}

func (c *converter) wrapExpr(expr hclsyntax.Expression) string {
	return "${" + c.rangeSource(expr.Range()) + "}"
}
