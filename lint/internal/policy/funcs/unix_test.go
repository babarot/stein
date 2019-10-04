package funcs

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestGrep(t *testing.T) {
	tests := map[string]struct {
		Pattern cty.Value
		Text    cty.Value
		Want    cty.Value
	}{
		"Match": {
			cty.StringVal("hoge"),
			cty.StringVal("test\nhogehoge\nfoo\nbar\n"),
			cty.StringVal("hogehoge"),
		},
		"NotMatch": {
			cty.StringVal("buz"),
			cty.StringVal("test\nhogehoge\nfoo\nbar\n"),
			cty.StringVal(""),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Grep(test.Pattern, test.Text)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestWc(t *testing.T) {
	tests := map[string]struct {
		Args []cty.Value
		Want cty.Value
	}{
		"OneLine": {
			[]cty.Value{cty.StringVal("foo\nbar baz")},
			cty.NumberIntVal(1),
		},
		"OneLineWithCharOption": {
			[]cty.Value{cty.StringVal("foo\nbar baz"), cty.StringVal("c")},
			cty.NumberIntVal(11),
		},
		"OneLineWithWordOption": {
			[]cty.Value{cty.StringVal("foo\nbar baz"), cty.StringVal("w")},
			cty.NumberIntVal(3),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Wc(test.Args[0], test.Args[1:]...)

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
