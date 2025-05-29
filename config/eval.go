package config

import (
	"fmt"
	"io/fs"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func newEvalContext(fsys fs.FS) *hcl.EvalContext {
	return &hcl.EvalContext{
		Functions: map[string]function.Function{
			"file": newFileFunction(fsys),
		},
	}
}

func newFileFunction(fsys fs.FS) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{{
			Name: "path",
			Type: cty.String,
		}},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, _ cty.Type) (cty.Value, error) {
			path := args[0].AsString()
			data, err := fs.ReadFile(fsys, path)
			if err != nil {
				return cty.UnknownVal(cty.String), function.NewArgError(0, err)
			}
			if !utf8.Valid(data) {
				return cty.UnknownVal(cty.String), fmt.Errorf("contents of %s are not valid UTF-8", path)
			}
			return cty.StringVal(string(data)), nil
		},
	})
}
